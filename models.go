package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type record struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Avatars   []avatar           `bson:"avatars,omitempty"`
	ForDate   time.Time          `bson:"forDate,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
}

type avatar struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Bio             string `json:"bio"`
	Location        string `json:"location"`
	URL             string `json:"url"`
	JoinDate        string `json:"join_date"`
	JoinTime        string `json:"join_time"`
	Tweets          string `json:"tweets"`
	Following       string `json:"following"`
	Followers       string `json:"followers"`
	Likes           string `json:"likes"`
	Media           string `json:"media"`
	Private         string `json:"private"`
	Verified        string `json:"verified"`
	ProfileImage    string `json:"profile_image"`
	BackgroundImage string `json:"background_image"`
}
