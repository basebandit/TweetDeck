package avatar

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Service defines dependencie required for creating a user service.
type Service struct {
	db  *mongo.Database
	ctx context.Context
}

//NewService returns an initialized user service for use with the database.
func NewService(ctx context.Context, db *mongo.Database) *Service {

	//create a unique index on user email field.
	if _, err := db.Collection("avatars").Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)}); err != nil {
		log.Fatalf("Users: failed to create unique index: %s", err)
	}

	return &Service{
		db:  db,
		ctx: ctx,
	}
}

//Insert inserts a new user record into the database
func (s Service) Insert(userID string, username string) (primitive.ObjectID, error) {
	var insertedID primitive.ObjectID

	avatar := Avatar{
		ID:        primitive.NewObjectID(),
		Username:  username,
		CreatedAt: time.Now(),
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return insertedID, err
	}

	avatar.UserID = id

	result, err := s.db.Collection("avatars").InsertOne(s.ctx, avatar)

	if err != nil {
		return insertedID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

//GetByID retrieves an avatar by id
func (s Service) GetByID(id string) (*Avatar, error) {
	avatar := Avatar{}

	aid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := s.db.Collection("avatars").FindOne(s.ctx, bson.M{"_id": aid}).Decode(&avatar); err != nil {
		return nil, err
	}

	return &avatar, nil
}

//GetByUsername retrieves an avatar by username
func (s Service) GetByUsername(username string) (*Avatar, error) {
	avatar := Avatar{}

	if err := s.db.Collection("avatars").FindOne(s.ctx, bson.M{"username": username}).Decode(&avatar); err != nil {
		return nil, err
	}

	return &avatar, nil
}

//GetAll retrieves all avatars
func (s Service) GetAll() ([]*Avatar, error) {

	avatars := []*Avatar{}

	cursor, err := s.db.Collection("avatars").Find(s.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(s.ctx)
	for cursor.Next(s.ctx) {
		var avatar Avatar
		if err = cursor.Decode(&avatar); err != nil {
			return nil, err
		}

		avatars = append(avatars, &avatar)
	}

	return avatars, nil
}

//UpdateUsername updates an existing avatar's name field value
func (s Service) UpdateUsername(id, username string) error {
	result, err := s.db.Collection("avatars").UpdateMany(s.ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"username": username, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}

	fmt.Printf("Updated %v documents(s)!\n", result.ModifiedCount)

	return nil
}
