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

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	a, err := GetByID(ctx, db, id)
	if err != nil {
		return err
	}

	if ua.Username != nil {
		a.Username = *ua.Username
	}

	if ua.UserID != nil {

		if _, err := uuid.Parse(*ua.UserID); err != nil {
			return ErrInvalidID
		}

		a.UserID = ua.UserID
	}

	a.UpdatedAt = now

	const q = `UPDATE avatars SET
	"user_id" = $2,
	"username" = $3,
	"updated_at" = $4
	WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, a.ID, *a.UserID, a.Username, a.UpdatedAt); err != nil {
		return errors.Wrap(err, "updating avatar")
	}

	return nil
}

//Delete removes the avatar identified by a given ID.
func Delete(ctx context.Context, db *sqlx.DB, id string, now time.Time) error {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.delete")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `UPDATE avatars SET
	active = $2,
	updated_at = $3
	WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, id, false, now); err != nil {
		return errors.Wrapf(err, "deleting avatar %s", id)
	}

	return nil
}

//Get retrieves all Avatars from the database.
func Get(ctx context.Context, db *sqlx.DB) ([]Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.get")
	defer span.End()

	const q = `with allp as (
	SELECT 
	a.id,
	a.username,
	a.user_id,
	p.followers,
	p.following,
	p.tweets,
	p.join_date,
	p.likes,
	p.bio,
	row_number() over (
		partition by 
		a.id,
		a.user_id ,
		a.username order by p.created_at desc,
		p.id desc) as priority_number from 
		avatars a LEFT JOIN profiles p ON
		a.id = p.avatar_id
		) 
		select 
		allp.id,
		allp.username,
		allp.user_id,
		allp.followers,
		allp.following,
		allp.tweets,
		allp.join_date,
		allp.likes,
		allp.bio from allp where priority_number = 1;`

	avatars := []Avatar{}
	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, errors.Wrap(err, "selecting avatars")
	}

	return avatars, nil
}

//GetByID finds the avatar identified by a given ID.
func GetByID(ctx context.Context, db *sqlx.DB, id string) (Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getbyid")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return Avatar{}, ErrInvalidID
	}

	const q = `SELECT
	  a.id,
		a.username,
		a.user_id,
		a.created_at,
		a.updated_at,
		p.bio,p.profile_image_url,p.twitter_id,
		p.followers,p.following, p.likes,p.tweets, p.join_date,p.last_tweet_time from avatars a LEFT JOIN
	 profiles p on a.id = p.avatar_id
		WHERE a.id=$1 AND active=TRUE ORDER BY p.created_at DESC LIMIT 1;
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

	const q = `with allp as (
		SELECT 
		a.id,
		a.username,
		a.user_id,
		p.followers,
		p.following,
		p.tweets,
		p.join_date,
		p.likes,
		p.bio,
		row_number() over (
			partition by 
			a.id,
			a.user_id ,
			a.username order by p.created_at desc,
			p.id desc) as priority_number from 
			avatars a LEFT JOIN profiles p ON
			a.id = p.avatar_id WHERE a.user_id=$1
			) 
			select 
			allp.id,
			allp.username,
			allp.user_id,
			allp.followers,
			allp.following,
			allp.tweets,
			allp.join_date,
			allp.likes,
			allp.bio from allp where priority_number = 1;
	`
	avatars := []Avatar{}

	if err := db.SelectContext(ctx, &avatars, q, userID); err != nil {
		return nil, errors.Wrap(err, "selecting avatars")
	}

	return avatars, nil
}
