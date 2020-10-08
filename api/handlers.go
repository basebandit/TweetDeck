package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-chi/chi"

	"ekraal.org/avatarlysis/business/data/auth"
	"ekraal.org/avatarlysis/business/data/avatar"

	"go.opentelemetry.io/otel/api/global"

	"net/http"

	"ekraal.org/avatarlysis/business/data/user"

	"github.com/go-chi/render"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	secret = "9f33e0f0086e439ebb41190167c9c83f62db6da68cac8ee95506855788b4abe9"
)

// func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
// 	signupRequest := struct {
// 		Name     string `json:"name"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}{}

// 	if err := json.NewDecoder(r.Body).Decode(&signupRequest); err != nil {
// 		if len(signupRequest.Email) == 0 || len(signupRequest.Password) == 0 {
// 			s.log.Printf("api: %v\n", errors.New("missing email of password field"))
// 			render.Render(w, r, ErrInvalidRequest(errors.New("missing email of password field")))
// 			return
// 		}
// 	}

// 	userService := user.NewService(r.Context(), s.db)

// 	_, err := userService.Insert(signupRequest.Name, signupRequest.Email, signupRequest.Password)
// 	if err != nil {
// 		var e mongo.WriteException
// 		if errors.As(err, &e) {
// 			//If it's aunique key violation
// 			for _, we := range e.WriteErrors {
// 				if we.Code == 11000 {
// 					s.log.Printf("api: %v\n", err)
// 					render.Render(w, r, ErrDuplicateField(ErrEmailTaken))
// 					return
// 				}
// 			}
// 		}
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	render.Status(r, http.StatusCreated)
// 	render.Respond(w, r, http.NoBody)
// }

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(r.Context(), "handlers.ping")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	fmt.Fprintf(w, "pong\n")
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(r.Context(), "handlers.signup")
	defer span.End()

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	var nu user.NewUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		s.log.Println(errors.Wrap(err, "Signup: decoding user"))
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	usr, err := user.Create(s.ctx, s.db, nu, v.Now)
	if err != nil {
		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode("23505") {
				s.log.Println(err)
				render.Render(w, r, ErrDuplicateField(ErrEmailTaken))
				return
			}
		}
		s.log.Println(errors.Wrapf(err, "User: %+v", &usr))
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, usr)
}

func (s *Server) handleToken(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.login")
	defer span.End()

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		render.Render(w, r, ErrUnauthorized(err))
		return
	}

	claims, err := user.Authenticate(ctx, s.db, v.Now, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			render.Render(w, r, ErrUnauthorized(err))
			return
		default:
			err := errors.Wrap(err, "authenticating")
			s.log.Println(err)
			render.Render(w, r, ErrInternalServerError)
			return
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	tkn.Token, err = s.auth.GenerateToken(claims)
	if err != nil {
		err := errors.Wrap(err, "generating token")
		s.log.Println(err)
		render.Render(w, r, ErrInternalServerError)
	}

	render.Respond(w, r, tkn)
}

func (s *Server) handleAssignAvatars(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.avatarupload")
	defer span.End()
	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	assignRequest := struct {
		UserID  string   `json:"userID"`
		Avatars []string `json:"avatars"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&assignRequest); err != nil {
		s.log.Println(errors.Wrap(err, "Signup: decoding user"))
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	s.log.Printf("%+v\n", assignRequest)

	var a avatar.UpdateAvatar
	for _, id := range assignRequest.Avatars {
		a.UserID = stringPointer(assignRequest.UserID)
		if err := avatar.UpdateUserID(ctx, s.db, id, a, time.Now()); err != nil {
			if err == avatar.ErrNotFound {
				s.log.Println(errors.Wrap(err, "Assign: assigining avatars to user"))
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
			s.log.Println(errors.Wrap(err, "Assign: assigining avatars to user"))
			render.Render(w, r, ErrInternalServerError)
			return
		}
	}

	render.Respond(w, r, http.NoBody)
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.avatarupload")
	defer span.End()

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	file, _, err := r.FormFile("avatars")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}
	usernames, err := readFile(file)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	var na avatar.NewAvatar
	var nas []avatar.NewAvatar

	for _, username := range usernames {
		na.Username = username
		nas = append(nas, na)
	}

	if err := avatar.CreateMultiple(ctx, s.db, nas, v.Now); err != nil {
		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode("23505") {
				s.log.Println(err)
				render.Render(w, r, ErrDuplicateField(ErrUsernameTaken))
				return
			}
		}
		s.log.Println(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, http.NoBody)
}

func (s *Server) handleAvatars(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.avatars")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	a, err := avatar.Get(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, a)
}

func (s *Server) handleTotalAvatars(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.totalavatars")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	total, err := avatar.GetTotalAccounts(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	t := struct {
		Count int `json:"count"`
	}{
		Count: total,
	}

	render.Respond(w, r, t)
}

func (s *Server) handleAvatarsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	uid := chi.URLParam(r, "id")

	id, err := user.Decode(uid)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	avatars, err := avatar.GetByUserID(ctx, s.db, id.String())
	if err != nil {
		if err == avatar.ErrNotFound {
			s.log.Printf("api: %v\n", err)
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, avatars)
}

func (s *Server) handlePeople(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	us, err := user.Get(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	usr := []struct {
		Avatars   int       `json:"avatars"`
		UserID    string    `json:"id"`
		Firstname string    `json:"firstname"`
		Lastname  string    `json:"lastname"`
		Tweets    int       `json:"tweets"`
		Followers int       `json:"followers"`
		Likes     int       `json:"likes"`
		Following int       `json:"following"`
		CreatedAt time.Time `json:"createdAt"`
	}{}
	//Make the user IDS url friendly and short

	for _, u := range us {
		u.UID = user.Encode(u.ID)
		count, err := avatar.GetAvatarCountByUserID(ctx, s.db, u.ID.String())
		if err != nil {
			s.log.Printf("api: %v\n", err)
			render.Render(w, r, ErrInternalServerError)
			return
		}

		avatar, err := avatar.AggregateAvatarByUserID(ctx, s.db, u.ID.String())
		if err != nil {
			s.log.Printf("api: %v\n", err)
			render.Render(w, r, ErrInternalServerError)
			return
		}

		usr = append(usr, struct {
			Avatars   int       `json:"avatars"`
			UserID    string    `json:"id"`
			Firstname string    `json:"firstname"`
			Lastname  string    `json:"lastname"`
			Tweets    int       `json:"tweets"`
			Followers int       `json:"followers"`
			Likes     int       `json:"likes"`
			Following int       `json:"following"`
			CreatedAt time.Time `json:"createdAt"`
		}{
			Avatars:   count,
			UserID:    u.UID,
			Tweets:    *avatar.Tweets,
			Following: *avatar.Following,
			Followers: *avatar.Followers,
			Likes:     *avatar.Likes,
			Firstname: u.Firstname,
			Lastname:  u.Lastname,
			CreatedAt: u.CreatedAt,
		})

	}

	render.Respond(w, r, usr)
}

func (s *Server) handleTotals(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	avs, err := avatar.Get(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	var (
		tweets    int
		following int
		followers int
		avatars   int
		likes     int
	)

	for _, av := range avs {
		if av.Tweets != nil {
			tweets += *av.Tweets //total tweets
		}

		if av.Following != nil {
			following += *av.Following //total following
		}

		if av.Followers != nil {
			followers += *av.Followers //total followers
		}

		if av.Likes != nil {
			likes += *av.Likes //total likes
		}

		avatars++ //total avatars
	}

	totals := struct {
		Avatars   int `json:"avatars"`
		Tweets    int `json:"tweets"`
		Likes     int `json:"likes"`
		Followers int `json:"followers"`
		Following int `json:"following"`
	}{
		avatars,
		tweets,
		likes,
		followers,
		following,
	}

	render.Respond(w, r, totals)
}

func (s *Server) handleSuspendedAvatars(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	avs, err := avatar.GetSuspendedAccounts(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, avs)
}

func (s *Server) handleTopFives(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	following, err := avatar.GetTopFiveByFollowing(s.ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	followers, err := avatar.GetTopFiveByFollowers(s.ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	tweets, err := avatar.GetTopFiveByTweets(s.ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	likes, err := avatar.GetTopFiveByLikes(s.ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	topFives := struct {
		Following []*avatar.Avatar `json:"following"`
		Followers []*avatar.Avatar `json:"followers"`
		Tweets    []*avatar.Avatar `json:"tweets"`
		Likes     []*avatar.Avatar `json:"likes"`
	}{
		Following: following,
		Followers: followers,
		Tweets:    tweets,
		Likes:     likes,
	}

	render.Respond(w, r, topFives)
}

func (s *Server) handleAddMember(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	var nu user.NewUser
	if err := json.NewDecoder(r.Body).Decode(&nu); err != nil {
		s.log.Println(errors.Wrap(err, "Signup: decoding user"))
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	usr, err := user.Create(s.ctx, s.db, nu, v.Now)
	if err != nil {
		if pqErr, ok := errors.Cause(err).(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode("23505") {
				s.log.Println(err)
				render.Render(w, r, ErrDuplicateField(ErrEmailTaken))
				return
			}
		}
		s.log.Println(errors.Wrapf(err, "User: %+v", &usr))
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, usr)
}

func (s *Server) handleUpdateMember(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	var uu user.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&uu); err != nil {
		s.log.Println(errors.Wrap(err, "updating: decoding user"))
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	claims, ok := r.Context().Value(auth.Key).(auth.Claims)
	if !ok {
		s.log.Println("missing claims")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	uid := chi.URLParam(r, "id")
	id, err := user.Decode(uid)
	if err != nil {
		s.log.Println(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := user.Update(ctx, claims, s.db, id.String(), uu, v.Now); err != nil {
		s.log.Println(err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, http.NoBody)
}

// func (s *Server) handleDailyStats(w http.ResponseWriter, r *http.Request) {
// 	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
// 	defer span.End()

// 	_, ok := ctx.Value(KeyValues).(*Values)
// 	if !ok {
// 		s.log.Println("web value missing from context")
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	var days int

// 	day := time.Now().Weekday()

// 	switch day {
// 	case time.Monday:
// 		days = 1
// 	case time.Tuesday:
// 		days = 2
// 	case time.Wednesday:
// 		days = 3
// 	case time.Thursday:
// 		days = 4
// 	case time.Friday:
// 		days = 5
// 	case time.Saturday:
// 		days = 6
// 	case time.Sunday:
// 		days = 7
// 	}

// 	avs, err := avatar.GetThePastDays(s.ctx, s.db, days)
// 	if err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	// fmt.Printf("%+v\n", avs)
// 	var (
// 		tweets    int
// 		likes     int
// 		followers int
// 		following int
// 	)
// 	result := make(map[string]map[string]int)
// 	for _, av := range avs {
// 		switch *av.Day {
// 		case time.Wednesday.String():
// 			tweets += *av.Tweets
// 			following += *av.Following
// 			followers += *av.Followers
// 			likes += *av.Likes

// 			if result[*av.Day] == nil {
// 				result[*av.Day] = make(map[string]int)
// 				result[*av.Day]["tweets"] = tweets
// 				result[*av.Day]["following"] = following
// 				result[*av.Day]["followers"] = followers
// 				result[*av.Day]["likes"] = likes
// 			}

// 		case time.Thursday.String():
// 			tweets += *av.Tweets
// 			following += *av.Following
// 			followers += *av.Followers
// 			likes += *av.Likes
// 			if result[*av.Day] == nil {
// 				result[*av.Day] = make(map[string]int)
// 				result[*av.Day]["tweets"] = tweets
// 				result[*av.Day]["following"] = following
// 				result[*av.Day]["followers"] = followers
// 				result[*av.Day]["likes"] = likes
// 			}

// 		}
// 	}

// 	fmt.Printf("%+v\n", result)
// 	// // total := struct {
// 	// // 	Day       string `json:"day"`
// 	// // 	Tweets    int    `json:"tweets"`
// 	// // 	Following int    `json:"following"`
// 	// // 	Followers int    `json:"followers"`
// 	// // 	Likes     int    `json:"likes"`
// 	// // }{}
// 	// totals := make(map[string]map[string]int)
// 	// // totals := []struct {
// 	// // 	Day       string `json:"day"`
// 	// // 	Tweets    int    `json:"tweets"`
// 	// // 	Following int    `json:"following"`
// 	// // 	Followers int    `json:"followers"`
// 	// // 	Likes     int    `json:"likes"`
// 	// // }{}

// 	// for _, a := range avs {
// 	// 	switch *a.Day {
// 	// 	case time.Monday.String():
// 	// 		tweets += *a.Tweets
// 	// 		following += *a.Following
// 	// 		followers += *a.Followers
// 	// 		likes += *a.Likes
// 	// 		if totals[*a.Day] == nil {
// 	// 			totals[*a.Day] = make(map[string]int)
// 	// 			totals[*a.Day]["tweets"] = tweets
// 	// 			totals[*a.Day]["following"] = following
// 	// 			totals[*a.Day]["followers"] = followers
// 	// 			totals[*a.Day]["likes"] = likes
// 	// 		}
// 	// 		// totals = append(totals, total)
// 	// 	case time.Tuesday.String():
// 	// 		tweets += *a.Tweets
// 	// 		following += *a.Following
// 	// 		followers += *a.Followers
// 	// 		likes += *a.Likes
// 	// 		if totals[*a.Day] == nil {
// 	// 			totals[*a.Day] = make(map[string]int)
// 	// 			totals[*a.Day]["tweets"] = tweets
// 	// 			totals[*a.Day]["following"] = following
// 	// 			totals[*a.Day]["followers"] = followers
// 	// 			totals[*a.Day]["likes"] = likes
// 	// 		}
// 	// 		// totals = append(totals, total)
// 	// 	case time.Wednesday.String():
// 	// 		tweets += *a.Tweets
// 	// 		following += *a.Following
// 	// 		followers += *a.Followers
// 	// 		likes += *a.Likes
// 	// 		if totals[*a.Day] == nil {
// 	// 			totals[*a.Day] = make(map[string]int)
// 	// 			totals[*a.Day]["tweets"] = tweets
// 	// 			totals[*a.Day]["following"] = following
// 	// 			totals[*a.Day]["followers"] = followers
// 	// 			totals[*a.Day]["likes"] = likes
// 	// 		}
// 	// 		// totals = append(totals, total)
// 	// 	case time.Thursday.String():
// 	// 		tweets += *a.Tweets
// 	// 		following += *a.Following
// 	// 		followers += *a.Followers
// 	// 		likes += *a.Likes

// 	// 		if totals[*a.Day] == nil {
// 	// 			totals[*a.Day] = make(map[string]int)
// 	// 			totals[*a.Day]["tweets"] = tweets
// 	// 			totals[*a.Day]["following"] = following
// 	// 			totals[*a.Day]["followers"] = followers
// 	// 			totals[*a.Day]["likes"] = likes
// 	// 		}
// 	// 		// totals = append(totals, total)
// 	// 	case time.Friday.String():
// 	// 		tweets += *a.Tweets
// 	// 		following += *a.Following
// 	// 		followers += *a.Followers
// 	// 		likes += *a.Likes
// 	// 		if totals[*a.Day] == nil {
// 	// 			totals[*a.Day] = make(map[string]int)
// 	// 			totals[*a.Day]["tweets"] = tweets
// 	// 			totals[*a.Day]["following"] = following
// 	// 			totals[*a.Day]["followers"] = followers
// 	// 			totals[*a.Day]["likes"] = likes
// 	// 		}
// 	// 		// totals = append(totals, total)
// 	// 	case time.Saturday.String():
// 	// 		tweets += *a.Tweets
// 	// 		following += *a.Following
// 	// 		followers += *a.Followers
// 	// 		likes += *a.Likes
// 	// 		if totals[*a.Day] == nil {
// 	// 			totals[*a.Day] = make(map[string]int)
// 	// 			totals[*a.Day]["tweets"] = tweets
// 	// 			totals[*a.Day]["following"] = following
// 	// 			totals[*a.Day]["followers"] = followers
// 	// 			totals[*a.Day]["likes"] = likes
// 	// 		}
// 	// 		// totals = append(totals, total)
// 	// 	case time.Sunday.String():
// 	// 		tweets += *a.Tweets
// 	// 		following += *a.Following
// 	// 		followers += *a.Followers
// 	// 		likes += *a.Likes
// 	// 		if totals[*a.Day] == nil {
// 	// 			totals[*a.Day] = make(map[string]int)
// 	// 			totals[*a.Day]["tweets"] = tweets
// 	// 			totals[*a.Day]["following"] = following
// 	// 			totals[*a.Day]["followers"] = followers
// 	// 			totals[*a.Day]["likes"] = likes
// 	// 		}
// 	// 		// totals = append(totals, total)
// 	// 	}
// 	// }

// 	render.Respond(w, r, avs)
// }

func readFile(reader io.Reader) ([]string, error) {

	r := csv.NewReader(reader)

	var usernames []string
	var header []string

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return usernames, err
		}

		if header == nil {
			header = line
			continue
		}

		usernames = append(usernames, line[0])

	}
	return usernames, nil
}

func stringPointer(s string) *string {
	str := s
	return &str
}
