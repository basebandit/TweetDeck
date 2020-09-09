package profile

import (
	"context"

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
