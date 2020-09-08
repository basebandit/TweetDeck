package avatar

import (
	"time"
)

//Avatar is a twitter account profile.
type Avatar struct {
	ID              string    `db:"id" json:"id"`                             //Unique Identifier.
	UserID          *string   `db:"user_id" json:"user_id"`                   //The user who manages/runs this twitter account.
	Username        string    `db:"username" json:"username"`                 //The twitter handle.
	Active          bool      `db:"active" json:"active"`                     //Use this flag to perform soft deletes.
	Assigned        *int      `json:"assigned"`                               //Shows whether this avatar profile is already assigned to an existing user or not. 0=not assigned 1= assigned
	Followers       *int      `db:"followers" json:"followers"`               //Total followers count on twitter
	Following       *int      `db:"following" json:"following"`               //Total following count on twitter
	Tweets          *int      `db:"tweets" json:"tweets"`                     //Total tweets count on twitter
	Likes           *int      `db:"likes" json:"likes"`                       //Total favorite count on twitter
	JoinDate        *string   `db:"join_date" json:"joinDate"`                //Day twitter account was created.
	ProfileImageURL *string   `db:"profile_image_url" json:"profileImageURL"` //Twitter profile image location
	Bio             *string   `db:"bio" json:"bio"`                           //Twitter account profile short description.
	TwitterID       *string   `db:"twitter_id" json:"twitterID"`              //Twitter ID
	LastTweetTime   *string   `db:"last_tweet_time" json:"lastTweetTime"`     //When last did this account tweet
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`              //When the record was added.
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`              //When the record was last modified.
}

//NewAvatar is what we require from clients when adding an Avatar.
type NewAvatar struct {
	Username string `json:"username" validate:"required"`
}

//UpdateAvatar defines what information may be provided to modify an
//existing Avatar.All fields are optional so clients can send only
//thos fields they wish to modify.
type UpdateAvatar struct {
	Username *string `json:"username"`
	UserID   *string `json:"userID"`
}
