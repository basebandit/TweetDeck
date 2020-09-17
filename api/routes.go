package api

import "github.com/go-chi/chi"

const (
	avatars         = "/api/avatar/"
	avatarUpload    = "/api/avatar/upload"
	avatarAssign    = "/api/avatar/assign"
	avatarsByUserID = "/api/avatar/{id}"

	peopleAssigned   = "/api/people/assigned"
	peopleUnassigned = "/api/people/unassigned"
	new              = "/api/people/new"

	signup = "/api/signup"
	login  = "/api/token"
	ping   = "/api/ping"
)

func (s *Server) routes() {
	logger := Logger(s.ctx, s.log)

	authenticate := Authenticate(s.auth)
	//middlewares
	s.router.Use(logger)

	s.router.Group(func(r chi.Router) {
		//===============Authenticated Routes==============
		r.Use(authenticate)
		//Avatars
		r.Post(avatarUpload, s.handleAvatarUpload)
		r.Post(avatarAssign, s.handleAssignAvatars)
		r.Get(avatarsByUserID, s.handleAvatarsByUserID)
		r.Get(avatars, s.handleAvatars)

		//People
		r.Get(peopleAssigned, s.handleAssignedPeople)
		r.Get(peopleUnassigned, s.handleUnassignedPeople)
		r.Post(new, s.handleAddMember)
	})

	//==============Unauthenticated Routes===============
	//Auth
	s.router.Post(signup, s.handleSignup)
	s.router.Get(ping, s.handlePing)
	s.router.Get(login, s.handleToken)

}
