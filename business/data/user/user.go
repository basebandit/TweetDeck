package user

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/api/global"
	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrNotFound is used when a specific User is requested but does not exist.
	ErrNotFound = errors.New("not found")

	//ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	//ErrAuthenticationFailure occurs when a user attempts to authenticate but
	//something goes wrong.
	ErrAuthenticationFailure = errors.New("authentication failed")
)

//Create inserts a new user into the database.
func Create(ctx context.Context, db *sqlx.DB, nu NewUser) (User, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.user.create")
	defer span.End()

	var u User
	if nu.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*nu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, errors.Wrap(err, "generating password hash")
		}

		u.Password = hash
	}

	if nu.Email != nil {
		u.Email = *nu.Email
	}

	u.ID = uuid.New().String()
	u.Firstname = nu.Firstname
	u.Lastname = nu.Lastname
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()

	const q = `INSERT INTO users
	(id,firstname,lastname,email,password,created_at,updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7)`

	if _, err := db.ExecContext(ctx, q, u.ID, u.Firstname, u.Lastname, u.Password, u.CreatedAt, u.UpdatedAt); err != nil {
		return User{}, errors.Wrap(err, "inserting user")
	}
	return u, nil
}
