package api

const (
	avatars      = "/api/avatar/"
	avatarPing   = "/api/avatar/ping"
	avatarUpload = "/api/avatar/upload"

	signup = "/api/signup"
	login  = "/api/login"
)

func (s *Server) routes() {
	//Auth
	s.router.Post(signup, s.handleSignup)
	s.router.Post(login, s.handleLogin)
	//Avatars
	s.router.Get(avatarPing, s.handleAvatarPing)
	s.router.Post(avatarUpload, s.handleAvatarUpload)
	s.router.Get(avatars, s.handleGetAvatar)
}
