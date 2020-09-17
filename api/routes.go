package api

const (
	avatars         = "/api/avatar/"
	avatarUpload    = "/api/avatar/upload"
	avatarAssign    = "/api/avatar/assign"
	avatarsByUserID = "/api/avatar/{id}"

	people = "/api/people"

	signup = "/api/signup"
	login  = "/api/token"
	ping   = "/api/ping"
)

func (s *Server) routes() {
	logger := Logger(s.ctx, s.log)

	authenticate := Authenticate(s.auth)
	//middlewares
	s.router.Use(logger)
	s.router.Use(authenticate)
	//Auth
	s.router.Post(signup, s.handleSignup)
	s.router.Get(ping, s.handlePing)
	s.router.Get(login, s.handleToken)
	//Avatars
	// s.router.Get(avatarPing, s.handleAvatarPing)
	s.router.Post(avatarUpload, s.handleAvatarUpload)
	s.router.Post(avatarAssign, s.handleAssignAvatars)
	s.router.Get(avatarsByUserID, s.handleAvatarsByUserID)
	s.router.Get(avatars, s.handleAvatars)
	// s.router.Get(avatars, s.handleGetAvatar)

	//People
	s.router.Get(people, s.handlePeople)
}
