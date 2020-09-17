package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-chi/chi"

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
		if err := avatar.Update(ctx, s.db, id, a, time.Now()); err != nil {
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

// func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
// 	loginRequest := struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}{}

// 	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
// 		if len(loginRequest.Email) == 0 || len(loginRequest.Password) == 0 {
// 			s.log.Printf("api: %v\n", errors.New("missing email of password field"))
// 			render.Render(w, r, ErrInvalidRequest(errors.New("missing email of password field")))
// 			return
// 		}
// 	}

// 	userService := user.NewService(r.Context(), s.db)
// 	u, err := userService.GetByEmail(loginRequest.Email)
// 	if err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrUnauthorized(ErrInvalidEmailOrPassword))
// 		return
// 	}

// 	if _, err := u.Compare(u.Password, loginRequest.Password); err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrUnauthorized(ErrInvalidEmailOrPassword))
// 		return
// 	}

// 	jwtService, err := jwt.NewService(hex.EncodeToString([]byte(secret)))
// 	if err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	token, err := jwtService.Create(u.ID)
// 	if err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	res := struct {
// 		Email  string `json:"email"`
// 		Name   string `json:"name"`
// 		UserID string `json:"id"`
// 		Token  string `json:"token"`
// 	}{
// 		Email:  u.Email,
// 		Name:   u.Name,
// 		UserID: u.ID.Hex(),
// 		Token:  token,
// 	}

// 	render.Respond(w, r, &res)
// }

// func (s *Server) handleAvatarPing(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "pong\n")
// }

// func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
// 	file, _, err := r.FormFile("avatars")
// 	if err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	usernames, err := readFile(file)
// 	if err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	twitterHandles := struct {
// 		Username []string `json:"handles"`
// 	}{
// 		Username: usernames,
// 	}

// 	render.Respond(w, r, twitterHandles)
// }

// func (s *Server) handleAvatarLookup(w http.ResponseWriter, r *http.Request) {

// 	//Retrieve avatars from twitter api.(Will be automated - cronjob after every 24 hours)

// }

// func (s *Server) handleAddAvatar(w http.ResponseWriter, r *http.Request) {
// 	avatarService := avatar.NewService(r.Context(), s.db)

// 	avatarRequest := struct {
// 		Username string `json:"username"`
// 	}{}

// 	if err := json.NewDecoder(r.Body).Decode(&avatarRequest); err != nil {
// 		if len(avatarRequest.Username) == 0 {
// 			render.Render(w, r, ErrInvalidRequest(ErrMissingUsername))
// 			return
// 		}
// 	}

// 	_, err := avatarService.Insert("5d0575344d9f7ff15e989174", avatarRequest.Username)
// 	if err != nil {
// 		s.log.Printf("api: %v\n", err)
// 		render.Render(w, r, ErrInternalServerError)
// 		return
// 	}

// 	render.Respond(w, r, http.NoBody)

// }

// func (s *Server) handleGetAvatar(w http.ResponseWriter, r *http.Request) {

// 	avatarService := avatar.NewService(r.Context(), s.db)

// 	avatars, err := avatarService.GetAll()
// 	if err != nil {
// 		s.log.Println(err)
// 		return
// 	}
// 	fmt.Println(avatars)
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
