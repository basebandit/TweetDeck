package api

const (
	avatars      = "/api/avatar/"
	avatarUpload = "/api/avatar/upload"

	people = "/api/people"

	signup = "/api/signup"
	login  = "/api/token"
	ping   = "/api/ping"
)

func (s *Server) routes() {
	logger := Logger(s.ctx, s.log)

	//middlewares
	s.router.Use(logger)

	//Auth
	s.router.Post(signup, s.handleSignup)
	s.router.Get(ping, s.handlePing)
	s.router.Get(login, s.handleToken)
	//Avatars
	// s.router.Get(avatarPing, s.handleAvatarPing)
	s.router.Post(avatarUpload, s.handleAvatarUpload)
	s.router.Get(avatars, s.handleAvatars)
	// s.router.Get(avatars, s.handleGetAvatar)

	//People
	s.router.Get(people, s.handlePeople)
}
