package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/api/global"
)

//Server is the api server.Our dependencies hang off this struct.
type Server struct {
	ctx    context.Context
	db     *sqlx.DB
	log    *log.Logger
	router *chi.Mux
}

//NewServer boostraps the api server and returns an initialized instance.
func NewServer(ctx context.Context, db *sqlx.DB, log *log.Logger, router *chi.Mux) *Server {
	// Start or expand a distributed trace.
	apiCtx, span := global.Tracer("service").Start(ctx, "api.server")
	defer span.End()

	// Set the context with the required values to
	// process the request.
	v := Values{
		TraceID: span.SpanContext().TraceID.String(),
		Now:     time.Now(),
	}
	apiCtx = context.WithValue(ctx, KeyValues, &v)

	s := Server{
		ctx:    apiCtx,
		db:     db,
		log:    log,
		router: router,
	}

	s.routes() //register handlers

	return &s
}

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values are stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
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
