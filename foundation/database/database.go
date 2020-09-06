package database

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"go.opentelemetry.io/otel/api/global"
)

//Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

//Open opens a database connection based on the configuration
func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

//Ping returns nil if it can successfully talk to the database. It
//returns a non-nil error otherwise.
func Ping(ctx context.Context, db *sqlx.DB) error {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "foundation.database.ping")
	defer span.End()

	//Run a simple query to determine connectivity. The db has a "Ping" method
	//but it can false-positive when it was previously able to talk to the
	//database but the database has since gone away. Running this query forces a
	//round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
