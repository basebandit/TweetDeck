package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

//Server is the api server.Our dependencies hang off this struct.
type Server struct {
	db     *mongo.Database
	log    *log.Logger
	router *chi.Mux
}

//NewServer boostraps the api server and returns an initialized instance.
func NewServer(db *mongo.Database, log *log.Logger, router *chi.Mux) *Server {
	s := Server{
		db:     db,
		log:    log,
		router: router,
	}

	s.routes() //register handlers

	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "PATCH,GET,OPTIONS,PUT,DELETE,POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if r.Method == "OPTIONS" {
		return
	}

	//let chi-router do the routing now
	s.router.ServeHTTP(w, r)
}
