package avatar

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Avatar defines the fields avatar model
type Avatar struct {
	ID               primitive.ObjectID `bson:"_id" json:"-"`
	TwitterID        string             `bson:"tid" json:"-"`
	Name             string             `bson:"name" json:"name"`
	Username         string             `bson:"username" json:"username"` //twitter handle
	Bio              string             `bson:"bio" json:"bio"`
	Location         string             `bson:"location" json:"location"`
	URL              string             `bson:"url" json:"url"`
	JoinDate         string             `bson:"joinDate" json:"joinDate"`
	JoinTime         string             `bson:"joinTime" json:"joinTime"`
	Tweets           int                `bson:"tweets" json:"tweets"`
	Following        int                `bson:"following" json:"following"`
	Followers        int                `bson:"followers" json:"followers"`
	Likes            int                `bson:"likes" json:"likes"`
	Media            int                `bson:"media" json:"-"`
	Private          bool               `bson:"private" json:"-"`
	Verified         bool               `bson:"verified" json:"-"`
	ProfileImage     string             `bson:"profileImage" json:"profileImage"`
	TwitterCreatedAt string             `bson:"twitterCreatedAt" json:"tCreatedAt"`
	BackgroundImage  string             `bson:"backgroundImage"`
	UserID           primitive.ObjectID `bson:"user" json:"userID"`
	CreatedAt        time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"-"`
}
