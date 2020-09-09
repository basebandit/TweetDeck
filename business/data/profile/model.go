//Package profile is a twitter account profile.
package profile

import "time"

//NewProfile states the fields that are required when creating a new profile.
type NewProfile struct {
	ID              string    `db:"id" json:"-"`                              //Unique identifier
	AvatarID        *string   `db:"avatar_id" json:"-"`                       //Avatar(twitter account) the owner of this profile
	Followers       *int      `db:"followers" json:"followers"`               //Twitter follower count
	Following       *int      `db:"following" json:"following"`               //Twitter following count
	Tweets          *int      `db:"tweets" json:"tweets"`                     //Twitter tweets count
	Likes           *int      `db:"likes" json:"likes"`                       //Twitter likes count
	Bio             *string   `db:"bio" json:"bio"`                           //Twitter account bio
	Name            *string   `db:"name" json:"name"`                         //Twitter account name
	TwitterID       *string   `db:"twitter_id" json:"twitterID"`              //Twitter string unique identifier
	ProfileImageURL *string   `db:"profile_image_url" json:"profileImageURL"` //Twitter account profile image
	LastTweetTime   *string   `db:"last_tweet_time" json:"lastTweetTime"`     //The last time this profile tweeted
	JoinDate        *string   `db:"join_date" json:"joinDate"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"` //Time this record was inserted into the db
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"` //Time this record was updated in the db
}

//Note we do not update twitter profiles in the database we insert new copies of the existing copies. Only the timestamps and unique IDs will differ, well along with other values but it is not a guarantee.