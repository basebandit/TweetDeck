package api

import (
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"ekraal.org/avatarlysis/jwt"

	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"ekraal.org/avatarlysis/avatar"
	"ekraal.org/avatarlysis/user"
	"github.com/go-chi/render"
)

const (
	secret = "9f33e0f0086e439ebb41190167c9c83f62db6da68cac8ee95506855788b4abe9"
)

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	signupRequest := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&signupRequest); err != nil {
		if len(signupRequest.Email) == 0 || len(signupRequest.Password) == 0 {
			s.log.Printf("api: %v\n", errors.New("missing email of password field"))
			render.Render(w, r, ErrInvalidRequest(errors.New("missing email of password field")))
			return
		}
	}

	userService := user.NewService(r.Context(), s.db)

	_, err := userService.Insert(signupRequest.Name, signupRequest.Email, signupRequest.Password)
	if err != nil {
		var e mongo.WriteException
		if errors.As(err, &e) {
			//If it's aunique key violation
			for _, we := range e.WriteErrors {
				if we.Code == 11000 {
					s.log.Printf("api: %v\n", err)
					render.Render(w, r, ErrDuplicateField(ErrEmailTaken))
					return
				}
			}
		}
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.Respond(w, r, http.NoBody)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	loginRequest := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		if len(loginRequest.Email) == 0 || len(loginRequest.Password) == 0 {
			s.log.Printf("api: %v\n", errors.New("missing email of password field"))
			render.Render(w, r, ErrInvalidRequest(errors.New("missing email of password field")))
			return
		}
	}

	userService := user.NewService(r.Context(), s.db)
	u, err := userService.GetByEmail(loginRequest.Email)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidEmailOrPassword))
		return
	}

	if _, err := u.Compare(u.Password, loginRequest.Password); err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrUnauthorized(ErrInvalidEmailOrPassword))
		return
	}

	jwtService, err := jwt.NewService(hex.EncodeToString([]byte(secret)))
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	token, err := jwtService.Create(u.ID)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	res := struct {
		Email  string `json:"email"`
		Name   string `json:"name"`
		UserID string `json:"id"`
		Token  string `json:"token"`
	}{
		Email:  u.Email,
		Name:   u.Name,
		UserID: u.ID.Hex(),
		Token:  token,
	}

	render.Respond(w, r, &res)
}

func (s *Server) handleAvatarPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

func (s *Server) handleAvatarUpload(w http.ResponseWriter, r *http.Request) {
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

	twitterHandles := struct {
		Username []string `json:"handles"`
	}{
		Username: usernames,
	}

	render.Respond(w, r, twitterHandles)
}

func (s *Server) handleAvatarLookup(w http.ResponseWriter, r *http.Request) {

	//Retrieve avatars from twitter api.(Will be automated - cronjob after every 24 hours)

}

func (s *Server) handleAddAvatar(w http.ResponseWriter, r *http.Request) {
	avatarService := avatar.NewService(r.Context(), s.db)

	avatarRequest := struct {
		Username string `json:"username"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&avatarRequest); err != nil {
		if len(avatarRequest.Username) == 0 {
			render.Render(w, r, ErrInvalidRequest(ErrMissingUsername))
			return
		}
	}

	_, err := avatarService.Insert("5d0575344d9f7ff15e989174", avatarRequest.Username)
	if err != nil {
		s.log.Printf("api: %v\n", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.Respond(w, r, http.NoBody)

}

func (s *Server) handleGetAvatar(w http.ResponseWriter, r *http.Request) {

	avatarService := avatar.NewService(r.Context(), s.db)

	avatars, err := avatarService.GetAll()
	if err != nil {
		s.log.Println(err)
		return
	}
	fmt.Println(avatars)
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
