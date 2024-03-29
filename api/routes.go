package api

import "github.com/go-chi/chi"

const (
	avatars          = "/api/avatar/"
	avatarsSuspended = "/api/avatar/suspended"
	avatarUpload     = "/api/avatar/upload"
	avatarAssign     = "/api/avatar/assign"
	avatarsByUserID  = "/api/avatar/{id}"

	people       = "/api/people"
	personUpdate = "/api/people/{id}/edit"
	new          = "/api/people/new"

	totals       = "/api/totals"
	totalAvatars = "/api/totals/avatars"

	dailyStats  = "/api/totals/daily"
	weeklyStats = "/api/totals/weekly"

	initialCreatedAt = "/api/totals/weekly/mindate"
	// topFiveByFollowers = "/api/top/followers"
	// topFiveByTweets    = "/api/top/tweets"
	topFives         = "/api/totals/top"
	totalDailyTweets = "/api/totals/daily/tweets"

	lookup = "/api/lookup/twitter"

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
		r.Get(avatarsSuspended, s.handleSuspendedAvatars)

		//People
		r.Get(people, s.handlePeople)
		r.Post(new, s.handleAddMember)
		r.Put(personUpdate, s.handleUpdateMember)

		//Totals
		r.Get(totals, s.handleTotals)
		r.Get(weeklyStats, s.handleWeeklyStats)
		r.Get(totalAvatars, s.handleTotalAvatars)

		//Top In Each category (followers,following,tweets)
		r.Get(topFives, s.handleTopFives)

		//Perform twitter lookup
		r.Get(lookup, s.handleTwitterLookup)
		r.Get(totalDailyTweets, s.handleTotalDailyTweets)

		//retrieve inital date (used to limit the date picker for weekly reports).To be used as the minimum limit
		r.Get(initialCreatedAt, s.handleGetProfileInitialCreateAt)
	})

	//==============Unauthenticated Routes===============
	//Auth
	s.router.Post(signup, s.handleSignup)
	s.router.Get(ping, s.handlePing)
	s.router.Get(login, s.handleToken)
}
