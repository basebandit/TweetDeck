package avatar

import (
	"context"
	"database/sql"
	"fmt"

	"ekraal.org/avatarlysis/business/data/user"

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

//CreateMultiple adds multiple Avatar records to the database with one query.It returns an error if not successful.
func CreateMultiple(ctx context.Context, db *sqlx.DB, na []NewAvatar, now time.Time) error {

	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.createmultiple")
	defer span.End()

	q := `INSERT INTO avatars 
	(id,username,created_at,updated_at) VALUES `

	insertParams := []interface{}{}

	for i, av := range na {
		p1 := i * 4
		a := Avatar{
			ID:        uuid.New().String(),
			Username:  av.Username,
			CreatedAt: now.UTC(),
			UpdatedAt: now.UTC(),
		}

		q += fmt.Sprintf("($%d,$%d,$%d,$%d),", p1+1, p1+2, p1+3, p1+4)
		insertParams = append(insertParams, a.ID, a.Username, a.CreatedAt, a.UpdatedAt)
	}

	q = q[:len(q)-1] //remove trailing ","
	q += "ON CONFLICT DO NOTHING"

	if _, err := db.ExecContext(ctx, q, insertParams...); err != nil {
		return errors.Wrap(err, "inserting multiple avatars")
	}

	return nil
}

//UpdateUserID modifies data about an existing Avatar.It will error if the specified Id is
//invalid or does not reference an existing Avatar.
func UpdateUserID(ctx context.Context, db *sqlx.DB, id string, ua UpdateAvatar, now time.Time) error {
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

		userID, err := user.Decode(*ua.UserID)
		if err != nil {
			return err
		}

		if _, err := uuid.Parse(userID.String()); err != nil {
			return ErrInvalidID
		}

		a.UserID = stringPointer(userID.String())
	}

	a.UpdatedAt = now

	const q = `UPDATE avatars SET
	"user_id" = $2,
	"updated_at" = $3
	WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, a.ID, *a.UserID, a.UpdatedAt); err != nil {
		return errors.Wrap(err, "updating avatar")
	}

	return nil
}

//UpdateUsername modifies data about an existing Avatar.It will error if the specified Id is
//invalid or does not reference an existing Avatar.
func UpdateUsername(ctx context.Context, db *sqlx.DB, id string, ua UpdateAvatar, now time.Time) error {
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

		userID, err := user.Decode(*ua.UserID)
		if err != nil {
			return err
		}

		if _, err := uuid.Parse(userID.String()); err != nil {
			return ErrInvalidID
		}

		a.UserID = stringPointer(userID.String())
	}

	a.UpdatedAt = now

	const q = `UPDATE avatars SET
	"username" = $2,
	"updated_at" = $3
	WHERE id = $1`

	if _, err := db.ExecContext(ctx, q, a.ID, a.Username, a.UpdatedAt); err != nil {
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

//GetSuspendedAccounts retrieves all suspended accounts
func GetSuspendedAccounts(ctx context.Context, db *sqlx.DB) ([]*Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getsuspendedaccounts")
	defer span.End()

	const q = `select followers,following,tweets,likes,profile_image_url, join_date,bio,a.username,(select concat(firstname,' ',lastname) as username from users where id=user_id) as person from avatars a left join profiles p on p.avatar_id=a.id and p.created_at=current_date where p.avatar_id is null`

	avatars := []*Avatar{}
	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, errors.Wrap(err, "selecting suspended accounts/avatars")
	}

	for _, avatar := range avatars {
		if avatar.UserID != nil {
			avatar.Assigned = intPointer(1)
		} else {
			avatar.Assigned = intPointer(0)
		}
	}

	return avatars, nil
}

//GetTotalAccounts retrieves all Avatars from the database.
func GetTotalAccounts(ctx context.Context, db *sqlx.DB) (int, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.get")
	defer span.End()

	const q = `select count(*) from avatars`

	var total int
	if err := db.SelectContext(ctx, &total, q); err != nil {
		return total, errors.Wrap(err, "selecting avatars")
	}

	return total, nil
}

//Get retrieves the most recent Avatars from the database.
// func Get(ctx context.Context, db *sqlx.DB) ([]*Avatar, error) {
// 	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.get")
// 	defer span.End()

// 	const q = `with allp as (
// 	SELECT
// 	a.id,
// 	a.username,
// 	a.user_id,
// 	p.followers,
// 	p.following,
// 	p.tweets,
// 	p.profile_image_url,
// 	p.join_date,
// 	p.likes,
// 	p.bio,
// 	p.created_at,
// 	row_number() over (
// 		partition by
// 		a.id,
// 		a.user_id ,
// 		a.username order by p.created_at desc,
// 		p.id desc) as priority_number from
// 		avatars a LEFT JOIN profiles p ON
// 		a.id = p.avatar_id
// 		)
// 		select
// 		allp.id,
// 		allp.username,
// 		allp.user_id,
// 		(select concat(firstname,' ',lastname) as username from users where id=allp.user_id) as person,
// 		allp.followers,
// 		allp.following,
// 		allp.tweets,
// 		allp.profile_image_url,
// 		allp.join_date,
// 		allp.likes,
// 		allp.bio from allp where priority_number = 1;`

// 	avatars := []*Avatar{}
// 	if err := db.SelectContext(ctx, &avatars, q); err != nil {
// 		return nil, errors.Wrap(err, "selecting avatars")
// 	}

// 	for _, avatar := range avatars {
// 		if avatar.UserID != nil {
// 			avatar.Assigned = intPointer(1)
// 		} else {
// 			avatar.Assigned = intPointer(0)
// 		}
// 	}

// 	return avatars, nil
// }

//Get retrieves the most recent Avatars from the database.
func Get(ctx context.Context, db *sqlx.DB) ([]*Avatar, error) {
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
	p.profile_image_url,
	p.join_date,
	p.likes,
	p.bio,
	p.created_at,
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
		(select concat(firstname,' ',lastname) as username from users where id=allp.user_id) as person,
		allp.followers,
		allp.following,
		allp.tweets,
		allp.profile_image_url,
		allp.join_date,
		allp.likes,
		allp.bio from allp where priority_number = 1 and allp.created_at=current_date;`

	avatars := []*Avatar{}
	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, errors.Wrap(err, "selecting avatars")
	}

	for _, avatar := range avatars {
		if avatar.UserID != nil {
			avatar.Assigned = intPointer(1)
		} else {
			avatar.Assigned = intPointer(0)
		}
	}

	return avatars, nil
}

//GetUsernames returns a map of usernames with their ids as key
func GetUsernames(ctx context.Context, db *sqlx.DB) ([]Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getusernames")
	defer span.End()

	const q = `SELECT id,username from avatars`

	avatars := []Avatar{}
	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, errors.Wrap(err, "selecting usernames")
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

	const q = `with allp as (
		SELECT 
		a.id,
		a.username,
		a.user_id,
		p.followers,
		p.following,
		p.profile_image_url,
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
			allp.profile_image_url,
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

//AggregateAvatarByUserID finds the avatars assigned to the given userID and adds up the totals of likes,following
//followers, tweets fields.
func AggregateAvatarByUserID(ctx context.Context, db *sqlx.DB, userID string) (Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getbyuserid")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return Avatar{}, ErrInvalidID
	}

	const q = `with allp as (
		SELECT 
		a.id,
		a.username,
		a.user_id,
		p.followers,
		p.following,
		p.profile_image_url,
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
			allp.profile_image_url,
			allp.join_date,
			allp.likes,
			allp.bio from allp where priority_number = 1;
	`
	avatars := []Avatar{}

	if err := db.SelectContext(ctx, &avatars, q, userID); err != nil {
		return Avatar{}, errors.Wrap(err, "selecting avatars")
	}

	avatar := Avatar{}

	var (
		tweets    int
		likes     int
		followers int
		following int
	)

	for _, a := range avatars {
		if a.Tweets != nil {
			tweets += *a.Tweets
		}

		if a.Likes != nil {
			likes += *a.Likes
		}
		if a.Following != nil {
			following += *a.Following
		}
		if a.Followers != nil {
			followers += *a.Followers
		}
	}

	avatar.Likes = intPointer(likes)
	avatar.Tweets = intPointer(tweets)
	avatar.Followers = intPointer(followers)
	avatar.Following = intPointer(following)

	fmt.Printf("Username: %s,Following: %d, Followers%d\n", avatar.Username, *avatar.Following, *avatar.Followers)

	return avatar, nil
}

//GetTopFiveByFollowers returns the top five avatars with highest followers in descending order.
func GetTopFiveByFollowers(ctx context.Context, db *sqlx.DB) ([]*Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.gettopfivebyfollowers")
	defer span.End()

	const q = `with allp as (select a.user_id, a.username, p.id, p.avatar_id, p.followers,
		p.following,
		p.tweets,
		p.likes,
		p.bio,
		p.created_at from profiles p left join avatars a on p.avatar_id = a.id where p.created_at=current_date group by p.id,a.user_id,a.username) 
	select allp.username,allp.followers,allp.following,allp.tweets,allp.likes,(select concat(firstname,' ',lastname) as username from users where id=allp.user_id) as person,allp.created_at from allp order by followers desc limit 5`

	var avatars []*Avatar

	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, err
	}

	return avatars, nil
}

//GetTopFiveByFollowing returns the top five avatars with highest following in descending order.
func GetTopFiveByFollowing(ctx context.Context, db *sqlx.DB) ([]*Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.gettopfivebyfollowers")
	defer span.End()

	const q = `with allp as (select a.user_id, a.username, p.id, p.avatar_id, p.followers,
		p.following,
		p.tweets,
		p.likes,
		p.bio,
		p.created_at from profiles p left join avatars a on p.avatar_id = a.id where p.created_at=current_date group by p.id,a.user_id,a.username) 
	select allp.username,allp.followers,allp.following,allp.tweets,allp.likes,(select concat(firstname,' ',lastname) as username from users where id=allp.user_id) as person,allp.created_at from allp order by following desc limit 5`

	var avatars []*Avatar

	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, err
	}

	return avatars, nil
}

//GetTopFiveByTweets returns the top five avatars with highest tweets in descending order.
func GetTopFiveByTweets(ctx context.Context, db *sqlx.DB) ([]*Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.gettopfivebyfollowers")
	defer span.End()

	const q = `with allp as (select a.user_id, a.username, p.id, p.avatar_id, p.followers,
		p.following,
		p.tweets,
		p.likes,
		p.bio,
		p.created_at from profiles p left join avatars a on p.avatar_id = a.id where p.created_at=current_date group by p.id,a.user_id,a.username)
	select allp.username,allp.followers,allp.following,allp.tweets,allp.likes,(select concat(firstname,' ',lastname) as username from users where id=allp.user_id) as person,allp.created_at from allp order by tweets desc limit 5`

	var avatars []*Avatar

	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, err
	}

	return avatars, nil
}

//GetTopFiveByLikes returns the top five avatars with highest likes in descending order.
func GetTopFiveByLikes(ctx context.Context, db *sqlx.DB) ([]*Avatar, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.gettopfivebyfollowers")
	defer span.End()

	const q = `with allp as (select a.user_id, a.username, p.id, p.avatar_id, p.followers,
		p.following,
		p.tweets,
		p.likes,
		p.bio,
		p.created_at from profiles p left join avatars a on p.avatar_id = a.id where p.created_at=current_date group by p.id,a.user_id,a.username)
	select allp.username,allp.followers,allp.following,allp.tweets,allp.likes,(select concat(firstname,' ',lastname) as username from users where id=allp.user_id) as person,allp.created_at from allp order by likes desc limit 5`

	var avatars []*Avatar

	if err := db.SelectContext(ctx, &avatars, q); err != nil {
		return nil, err
	}

	return avatars, nil
}

//GetAvatarCountByUserID returns the number of avatars assigned to the user with the
//given id. Returns -1 if error was encountered.
func GetAvatarCountByUserID(ctx context.Context, db *sqlx.DB, userID string) (int, error) {
	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getbyuserid")
	defer span.End()

	if _, err := uuid.Parse(userID); err != nil {
		return -1, ErrInvalidID
	}
	const q = `select count(*) from avatars where user_id=$1`
	var count int
	if err := db.GetContext(ctx, &count, q, userID); err != nil {
		return -1, errors.Wrap(err, "selecting avatars")
	}

	return count, nil
}

// //GetThePastDays retrieves avatar profile records for the past given days inclusive of today's records.
// func GetThePastDays(ctx context.Context, db *sqlx.DB, days int) ([]*Avatar, error) {
// 	ctx, span := global.Tracer("avatarlysis").Start(ctx, "business.data.avatar.getbyuserid")
// 	defer span.End()

// 	var q = fmt.Sprintf("select p.id,a.username,p.tweets, p.followers, p.following, p.likes, date(p.created_at) as created_at from profiles p left join avatars a on p.avatar_id=a.id where date(p.created_at) > NOW()- interval '%d day' order by date(p.created_at) asc", days)

// 	var avatars []*Avatar

// 	if err := db.SelectContext(ctx, &avatars, q); err != nil {
// 		return nil, err
// 	}

// 	// var (
// 	// 	tweets    int
// 	// 	likes     int
// 	// 	followers int
// 	// 	following int
// 	// )

// 	// totals := make(map[time.Weekday]map[string]int)

// 	for _, a := range avatars {
// 		switch day := a.CreatedAt.Weekday(); day {
// 		case time.Monday:
// 			// tweets += *a.Tweets
// 			// following += *a.Following
// 			// followers += *a.Followers
// 			// likes += *a.Likes
// 			a.Day = stringPointer(time.Monday.String())
// 			// if totals[time.Monday] == nil {
// 			// 	totals[time.Monday] = make(map[string]int)
// 			// 	totals[time.Monday]["tweets"] = tweets
// 			// 	totals[time.Monday]["following"] = following
// 			// 	totals[time.Monday]["follower"] = followers
// 			// 	totals[time.Monday]["likes"] = likes
// 			// }
// 		case time.Tuesday:
// 			// tweets += *a.Tweets
// 			// following += *a.Following
// 			// followers += *a.Followers
// 			// likes += *a.Likes
// 			a.Day = stringPointer(time.Tuesday.String())
// 			// if totals[time.Tuesday] == nil {
// 			// 	totals[time.Tuesday] = make(map[string]int)
// 			// 	totals[time.Tuesday]["tweets"] = tweets
// 			// 	totals[time.Tuesday]["following"] = following
// 			// 	totals[time.Tuesday]["follower"] = followers
// 			// 	totals[time.Tuesday]["likes"] = likes
// 			// }
// 		case time.Wednesday:
// 			// tweets += *a.Tweets
// 			// following += *a.Following
// 			// followers += *a.Followers
// 			// likes += *a.Likes
// 			a.Day = stringPointer(time.Wednesday.String())
// 			// if totals[time.Wednesday] == nil {
// 			// 	totals[time.Wednesday] = make(map[string]int)
// 			// 	totals[time.Wednesday]["tweets"] = tweets
// 			// 	totals[time.Wednesday]["following"] = following
// 			// 	totals[time.Wednesday]["follower"] = followers
// 			// 	totals[time.Wednesday]["likes"] = likes
// 			// }
// 		case time.Thursday:
// 			// tweets += *a.Tweets
// 			// following += *a.Following
// 			// followers += *a.Followers
// 			// likes += *a.Likes
// 			a.Day = stringPointer(time.Thursday.String())
// 			// if totals[time.Thursday] == nil {
// 			// 	totals[time.Thursday] = make(map[string]int)
// 			// 	totals[time.Thursday]["tweets"] = tweets
// 			// 	totals[time.Thursday]["following"] = following
// 			// 	totals[time.Thursday]["follower"] = followers
// 			// 	totals[time.Thursday]["likes"] = likes
// 			// }
// 		case time.Friday:
// 			// tweets += *a.Tweets
// 			// following += *a.Following
// 			// followers += *a.Followers
// 			// likes += *a.Likes
// 			a.Day = stringPointer(time.Friday.String())
// 			// if totals[time.Friday] == nil {
// 			// 	totals[time.Friday] = make(map[string]int)
// 			// 	totals[time.Friday]["tweets"] = tweets
// 			// 	totals[time.Friday]["following"] = following
// 			// 	totals[time.Friday]["follower"] = followers
// 			// 	totals[time.Friday]["likes"] = likes
// 			// }
// 		case time.Saturday:
// 			// tweets += *a.Tweets
// 			// following += *a.Following
// 			// followers += *a.Followers
// 			// likes += *a.Likes
// 			a.Day = stringPointer(time.Saturday.String())
// 			// if totals[time.Saturday] == nil {
// 			// 	totals[time.Saturday] = make(map[string]int)
// 			// 	totals[time.Saturday]["tweets"] = tweets
// 			// 	totals[time.Saturday]["following"] = following
// 			// 	totals[time.Saturday]["follower"] = followers
// 			// 	totals[time.Saturday]["likes"] = likes
// 			// }
// 		case time.Sunday:

// 			// tweets += *a.Tweets
// 			// following += *a.Following
// 			// followers += *a.Followers
// 			// likes += *a.Likes
// 			a.Day = stringPointer(time.Sunday.String())
// 			// if totals[time.Sunday] == nil {
// 			// 	totals[time.Sunday] = make(map[string]int)
// 			// 	totals[time.Sunday]["tweets"] = tweets
// 			// 	totals[time.Sunday]["following"] = following
// 			// 	totals[time.Sunday]["follower"] = followers
// 			// 	totals[time.Sunday]["likes"] = likes
// 			// }
// 		}
// 	}

// 	// fmt.Printf("Totals %+v\n", totals)

// 	return avatars, nil
// }

// intPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func intPointer(i int) *int {
	return &i
}

// intPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func stringPointer(s string) *string {
	return &s
}
