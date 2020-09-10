package api

const (
	avatars      = "/api/avatar/"
	avatarUpload = "/api/avatar/upload"

	signup = "/api/signup"
	login  = "/api/login"
	ping   = "/api/ping"
)

func (s *Server) routes() {
	//Auth
	s.router.Post(signup, s.handleSignup)
	s.router.Get(ping, s.handlePing)
	// s.router.Post(login, s.handleLogin)
	//Avatars
	// s.router.Get(avatarPing, s.handleAvatarPing)
	// s.router.Post(avatarUpload, s.handleAvatarUpload)
	// s.router.Get(avatars, s.handleGetAvatar)
}
