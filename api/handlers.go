package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"ekraal.org/avatarlysis/business/data/auth"
	"ekraal.org/avatarlysis/business/data/avatar"
	"ekraal.org/avatarlysis/business/data/profile"
	service "ekraal.org/avatarlysis/foundation/service/twitter"

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

func (s *Server) handleTotalDailyTweets(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.totalavatars")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	total, err := avatar.GetDailyTotalBy(ctx, s.db, "tweets")
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

	tweets, err := avatar.GetDailyTotalBy(ctx, s.db, "tweets")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	followers, err := avatar.GetDailyTotalBy(ctx, s.db, "followers")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	likes, err := avatar.GetDailyTotalBy(ctx, s.db, "likes")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	following, err := avatar.GetDailyTotalBy(ctx, s.db, "following")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	newAccts, err := avatar.GetNewAccounts(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	// var (
	// 	tweets    int
	// 	following int
	// 	followers int
	// 	avatars   int
	// 	likes     int
	// )

	// for _, av := range avs {
	// 	if av.Tweets != nil {
	// 		tweets += *av.Tweets //total tweets
	// 	}

	// 	if av.Following != nil {
	// 		following += *av.Following //total following
	// 	}

	// 	if av.Followers != nil {
	// 		followers += *av.Followers //total followers
	// 	}

	// 	if av.Likes != nil {
	// 		likes += *av.Likes //total likes
	// 	}

	// 	avatars++ //total avatars
	// }

	// var mgByFollowers *avatar.Avatar

	//TODO:
	// mgByFollowers, err = avatar.GetHighestGainedBy(s.ctx, s.db, "Followers", 1)

	// if err != nil {
	// 	s.log.Printf("api: %v\n", err)
	// 	render.Render(w, r, ErrInternalServerError)
	// 	return
	// }

	totals := struct {
		Avatars     int `json:"avatars"`
		Tweets      int `json:"tweets"`
		Likes       int `json:"likes"`
		Followers   int `json:"followers"`
		Following   int `json:"following"`
		NewAccounts int `json:"newAccounts"`
	}{
		len(avs),
		tweets,
		likes,
		followers,
		following,
		newAccts,
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

	following, err := avatar.GetTopFiveDailyBy(s.ctx, s.db, "following")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	followers, err := avatar.GetTopFiveDailyBy(s.ctx, s.db, "followers")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	tweets, err := avatar.GetTopFiveDailyBy(s.ctx, s.db, "tweets")
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	likes, err := avatar.GetTopFiveDailyBy(s.ctx, s.db, "likes")
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

func (s *Server) handleTwitterLookup(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	// log.Println(s.cfg.TwitterConsumerKey)
	if err := s.TwitterLookup(); err != nil {
		s.log.Println(errors.Wrap(err, "api: twitter lookup"))
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, http.NoBody)
}

func (s *Server) handleGetProfileInitialCreateAt(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.handlegetprofielinitialCreateat")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	initialDate, err := profile.GetInitialDate(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	endDate, err := profile.GetMostRecentCreateDate(ctx, s.db)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	res := struct {
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}{
		StartDate: initialDate.Format("2006-1-02"),
		EndDate:   endDate.Format("2006-1-02"),
	}

	render.Respond(w, r, res)
}

func (s *Server) handleWeeklyStats(w http.ResponseWriter, r *http.Request) {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "handlers.people")
	defer span.End()

	_, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		s.log.Println("web value missing from context")
		render.Render(w, r, ErrInternalServerError)
		return
	}

	var start, end string

	//Get the query params
	//http://localhost:8880/api/totals/weekly?start=""&end=""
	if r.URL.Query().Get("start") == "" || r.URL.Query().Get("end") == "" {
		s.log.Println("invalid request.missing query parameter(s)")
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid request.missing query parameter(s)")))
		return
	}

	start = r.URL.Query().Get("start")
	end = r.URL.Query().Get("end")

	fmt.Println("start: ", start, "end: ", end)

	render.Respond(w, r, http.NoBody)
}

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

//TwitterLookup perform twitter user lookup
func (s *Server) TwitterLookup() error {
	ctx, span := global.Tracer("avatarlysis").Start(s.ctx, "api.twitterlookup")
	defer span.End()

	avatars, err := avatar.GetUsernames(s.ctx, s.db)
	if err != nil {
		return err
	}

	var usernames []string
	dict := make(map[string]string)

	if len(avatars) > 0 {
		for _, avatar := range avatars {
			usernames = append(usernames, avatar.Username)
			dict[strings.ToLower(avatar.Username)] = avatar.ID
		}

		twitter := service.NewTwitter(s.cfg.TwitterConsumerKey, s.cfg.TwitterConsumerSecret, s.cfg.TwitterAccessToken, s.cfg.TwitterTokenURL)

		chunks := chunkBy(usernames, 100)
		var nps []profile.NewProfile
		now := time.Now()
		for _, chunk := range chunks {

			tusers := twitter.Lookup(ctx, chunk)

			for _, user := range tusers {
				np := profile.NewProfile{}
				np.ID = uuid.New().String()
				np.CreatedAt = now
				np.UpdatedAt = now
				avatarID, ok := dict[strings.ToLower(user.ScreenName)] //We just being paranoid just incase
				//because tusers only contains the usernames whose accounts actually exist in twitter which in any case
				// they already exist in our database
				if !ok {
					continue //skip to the next iteration
				}
				np.AvatarID = stringPointer(avatarID)
				np.Name = stringPointer(user.Name)
				np.Followers = intPointer(user.FollowersCount)
				np.Following = intPointer(user.FriendsCount)
				np.Likes = intPointer(user.FavouritesCount)
				np.Tweets = intPointer(user.StatusesCount)
				np.ProfileImageURL = stringPointer(strings.ReplaceAll(user.ProfileImageURLHttps, "_normal", ""))
				np.Bio = stringPointer(user.Description)
				np.TwitterID = stringPointer(user.IDStr)
				np.JoinDate = stringPointer(user.CreatedAt)
				np.LastTweetTime = stringPointer(twitter.LastTweetTime(user.ID))
				nps = append(nps, np)
				// if err := profile.Create(s.ctx, s.db, &np, now); err != nil {
				// 	fmt.Printf("profiler: %v\n", err)
				// 	return err
				// }
			}
		}
		if err := profile.CreateMultiple(s.ctx, s.db, nps, now); err != nil {
			fmt.Printf("profiler: %v\n", err)
			return err
		}
	}
	return nil
}

func stringPointer(s string) *string {
	return &s
}

func intPointer(n int) *int {
	return &n
}

func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
