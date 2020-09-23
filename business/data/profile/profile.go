package profile

import (
	"context"
	"fmt"

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

//Create adds a new Profile to the database.It returns the created profile.
func Create(ctx context.Context, db *sqlx.DB, avatarID string, np *NewProfile, now time.Time) error {

	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.create")
	defer span.End()

	np.ID = uuid.New().String()
	np.CreatedAt = now
	np.UpdatedAt = now

	const q = `INSERT INTO profiles (id,avatar_id,followers,following,tweets,likes,bio,name,twitter_id,profile_image_url,last_tweet_time,join_date,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)`

	if _, err := db.ExecContext(ctx, q, np.ID, np.AvatarID, np.Followers, np.Following, np.Tweets, np.Likes, np.Bio, np.Name, np.TwitterID, np.ProfileImageURL, np.LastTweetTime, np.JoinDate, np.CreatedAt, np.UpdatedAt); err != nil {
		return errors.Wrap(err, "inserting profile")
	}

	return nil
}

//CreateMultiple adds multiple Profile records to the database with one query.It returns an error if not succesful.
func CreateMultiple(ctx context.Context, db *sqlx.DB, np []NewProfile, now time.Time) error {

	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.createmultiple")
	defer span.End()

	q := `INSERT INTO profiles(id,avatar_id,followers,"following",tweets,likes,bio,"name",twitter_id,profile_image_url,last_tweet_time,join_date,created_at,updated_at) VALUES `

	insertParams := []interface{}{}

	if len(np) > 0 {

		var (
			avatarID        string
			following       int
			followers       int
			tweets          int
			likes           int
			bio             string
			name            string
			twitterID       string
			profileImageURL string
			lastTweetTime   string
			joinDate        string
		)

		for i, p := range np {
			p1 := i * 14

			if p.Following != nil {
				following = *p.Following
			}
			if p.Followers != nil {
				followers = *p.Followers
			}

			if p.Tweets != nil {
				tweets = *p.Tweets
			}

			if p.Bio != nil {
				bio = *p.Bio
			}

			if p.TwitterID != nil {
				twitterID = *p.TwitterID
			}

			if p.Likes != nil {
				likes = *p.Likes
			}

			if p.JoinDate != nil {
				joinDate = *p.JoinDate
			}

			if p.Name != nil {
				name = *p.Name
			}

			if p.AvatarID != nil {
				avatarID = *p.AvatarID
			}
			if p.LastTweetTime != nil {
				lastTweetTime = *p.LastTweetTime
			}

			if p.ProfileImageURL != nil {
				profileImageURL = *p.ProfileImageURL
			}

			p.CreatedAt = now
			p.UpdatedAt = now

			q += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", p1+1, p1+2, p1+3, p1+4, p1+5, p1+6, p1+7, p1+8, p1+9, p1+10, p1+11, p1+12, p1+13, p1+14)

			insertParams = append(insertParams, p.ID, avatarID, followers, following, tweets, likes, bio, name, twitterID, profileImageURL, lastTweetTime, joinDate, p.CreatedAt, p.UpdatedAt)

		}

		q = q[:len(q)-1] //remove trailing ","

		if _, err := db.ExecContext(ctx, q, insertParams...); err != nil {
			return errors.Wrap(err, "inserting multiple profiles")
		}
	}
	return nil
}
