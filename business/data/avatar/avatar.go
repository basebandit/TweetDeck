package avatar

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/api/global"
)

var (
	//ErrNotFound is used when a specified Avatar record is requested but does not exist.
	ErrNotFound = errors.New("not found")

	//ErrInvalidID used when the provided ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its valid form")
)

//Create adds an Avatar record to the database.It returns the created Avatar.
func Create(ctx context.Context, db *sqlx.DB, na NewAvatar, now time.Time) (Avatar, error) {

	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.create")
	defer span.End()

	a := Avatar{
		ID:        uuid.New().String(),
		Username:  na.Username,
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
	}

	const q = `INSERT INTO avatars 
	(id,username,created_at,updated_at) VALUES
	($1,$2,$3,$4)`

	if _, err := db.ExecContext(ctx, q, a.ID, a.Username, a.CreatedAt, a.UpdatedAt); err != nil {
		return Avatar{}, errors.Wrap(err, "inserting avatar")
	}

	return a, nil
}

//Update modifies data about an existing Avatar.It will error if the specified Id is
//invalid or does not reference an existing Avatar.
func Update(ctx context.Context, db *sqlx.DB, id string, ua UpdateAvatar, now time.Time) error {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.update")
	defer span.End()

	a, err := GetByID(ctx, db, id)
	if err != nil {
		return err
	}

	if ua.Username != nil {
		a.Username = *ua.Username
	}

	if ua.UserID != nil {
		a.UserID = *ua.UserID
	}

	a.UpdatedAt = now

	const q = `UPDATE avatars SET
	"username" = $2,
	"user_id" = $3,
	"updated_at" = $4
	WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, id, a.Username, a.UserID); err != nil {
		return errors.Wrap(err, "updating avatar")
	}

	return nil
}

//GetByID finds the avatar identified by a given ID.
func GetByID(ctx context.Context, db *sqlx.DB, id string) (Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getbyid")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return Avatar{}, ErrInvalidID
	}

	const q = `SELECT
	a.username AS username,
	a.user_id AS user_id,
	p.followers AS followers,
	p.following AS following,
	p.tweets AS tweets,
	p.join_date AS join_date,
	p.likes AS likes,
	p.bio AS bio from avatars a LEFT JOIN
 profiles p on a.id = p.avatar_id 
	WHERE a.id=$1 ORDER BY p.created_at DESC LIMIT 1;
	`
	var a Avatar
	if err := db.GetContext(ctx, &a, q, id); err != nil {
		if err == sql.ErrNoRows {
			return Avatar{}, ErrNotFound
		}
		return Avatar{}, errors.Wrap(err, "selecting single avatar")
	}
	//TODO: Check if user_id field is null using sqlx then set a.Assigned to 0 if it is null otherwise set to 1.

	return a, nil
}

//GetByUserID finds the avatars assigned to the given userID.
func GetByUserID(ctx context.Context, db *sqlx.DB, userID string) ([]Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getbyuserid")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return nil, ErrInvalidID
	}

	avatars := []Avatar{}

	const q = `with allp as (
		SELECT a.username,
		a.user_id,
		p.followers,
		p.following,
		p.tweets,
		p.join_date,
		p.likes,
		p.bio,
		row_number() over (
			partition by a.user_id ,
			a.username order by p.created_at desc,
			p.id desc) as priority_number from 
			avatars a LEFT JOIN profiles p ON
			a.id = p.avatar_id WHERE a.user_id='45b5fbd3-755f-4379-8f07-a58d4a30fa2f'
			) 
			select allp.username,
			allp.user_id,
			allp.followers,
			allp.following,
			allp.tweets,
			allp.join_date,
			allp.likes,
			allp.bio from allp where priority_number = 1;
	`
	return avatars, nil
}
