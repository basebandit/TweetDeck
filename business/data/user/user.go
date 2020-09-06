package user

import (
	"context"
	"database/sql"
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
func Create(ctx context.Context, db *sqlx.DB, nu NewUser, now time.Time) (User, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.user.create")
	defer span.End()

	var u User
	if nu.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*nu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, errors.Wrap(err, "generating password hash")
		}

		u.PasswordHash = hash
	}

	if nu.Email != nil {
		u.Email = *nu.Email
	}

	u.ID = uuid.New().String()
	u.Firstname = nu.Firstname
	u.Lastname = nu.Lastname
	u.Active = true
	u.CreatedAt = now
	u.UpdatedAt = now

	const q = `INSERT INTO users
	(id,firstname,lastname,email,password_hash,active,created_at,updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	if _, err := db.ExecContext(ctx, q, u.ID, u.Firstname, u.Lastname, u.Email, u.PasswordHash, u.Active, u.CreatedAt, u.UpdatedAt); err != nil {
		return User{}, errors.Wrap(err, "inserting user")
	}
	return u, nil
}

//Update replaces an existing user record in the database.
func Update(ctx context.Context, db *sqlx.DB, userID string, uu UpdateUser, now time.Time) error {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.user.update")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidID
	}

	var u User

	if uu.Firstname != nil {
		u.Firstname = *uu.Firstname
	}
	if uu.Lastname != nil {
		u.Lastname = *uu.Lastname
	}

	if uu.Email != nil {
		u.Email = *uu.Email
	}

	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "generating password hash")
		}
		u.PasswordHash = pw
	}

	u.UpdatedAt = now

	const q = `UPDATE users SET
	 "firstname" = $2,
	 "lastname" = $3,
	 "email" = $4,
	 "password_hash" = $5,
	 "updated_at" = $6
	 WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, userID, u.Firstname, u.Lastname, u.Email, u.PasswordHash, u.UpdatedAt); err != nil {
		return errors.Wrap(err, "updating user")
	}

	return nil
}

//Delete marks an existing user record as inactive. It thus performs a soft delete.
func Delete(ctx context.Context, db *sqlx.DB, userID string) error {
	ctx, span := global.Tracer("avartalysis").Start(ctx, "business.data.user.delete")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidID
	}

	u := User{
		Active: false,
	}

	const q = `UPDATE users SET
	"active" = $2
	WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, userID, u.Active); err != nil {
		return errors.Wrapf(err, "deleting user %s", userID)
	}

	return nil
}

//Get retrieves a list of all existing user records from the database.
func Get(ctx context.Context, db *sqlx.DB) ([]User, error) {
	ctx, span := global.Tracer("avartalysis").Start(ctx, "business.data.user.get")
	defer span.End()

	const q = `SELECT * FROM users WHERE active = TRUE`

	users := []User{}
	if err := db.SelectContext(ctx, &users, q); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}

	return users, nil
}

//GetByID retrieves the specified user from the database.
func GetByID(ctx context.Context, db *sqlx.DB, userID string) (User, error) {
	ctx, span := global.Tracer("avartalysis").Start(ctx, "business.data.user.getbyid")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return User{}, ErrInvalidID
	}

	const q = `SELECT * FROM users WHERE id = $1 AND active = TRUE`

	var u User
	if err := db.GetContext(ctx, &u, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrNotFound
		}
		return User{}, errors.Wrapf(err, "selecting user %q", userID)
	}
	return u, nil
}

//Authenticate finds a user by their email and verifies their password against the stored hash.On success it returns nil otherwise it returns the error.
func Authenticate(ctx context.Context, db *sqlx.DB, email, password string) error {
	ctx, span := global.Tracer("avartalysis").Start(ctx, "business.data.user.authenticate")
	defer span.End()

	const q = `SELECT * FROM users WHERE email = $1`

	var u User
	if err := db.GetContext(ctx, &u, q, email); err != nil {
		if err == sql.ErrNoRows {
			return ErrAuthenticationFailure
		}
		return errors.Wrapf(err, "selecting single user email %s", email)
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return ErrAuthenticationFailure
	}

	return nil
}
