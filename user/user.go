package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

//Service defines dependencies required for creating a user service.
type Service struct {
	db  *mongo.Database
	ctx context.Context
}

//NewService returns an initialized user service for use with the database.
func NewService(ctx context.Context, db *mongo.Database) *Service {

	//create a unique index on user email field.
	if _, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)}); err != nil {
		log.Fatalf("Users: failed to create unique index: %s", err)
	}

	return &Service{
		db:  db,
		ctx: ctx,
	}
}

//Insert inserts a new user record into the database
func (s Service) Insert(name, email, password string) (primitive.ObjectID, error) {
	var insertedID primitive.ObjectID

	user := User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Active:    true,
		Password:  password,
		CreatedAt: time.Now(),
	}

	if err := user.hashPassword(); err != nil {
		return insertedID, err
	}

	result, err := s.db.Collection("users").InsertOne(s.ctx, user)

	if err != nil {
		return insertedID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

//GetByID retrieves a user by id
func (s Service) GetByID(id string) (User, error) {
	user := User{}

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	if err := s.db.Collection("users").FindOne(s.ctx, bson.M{"_id": docID}).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

//GetByEmail retrieves a user by email
func (s Service) GetByEmail(email string) (User, error) {
	user := User{}

	if err := s.db.Collection("users").FindOne(s.ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return user, err
	}

	return user, nil
}

//UpdateEmail updates an existing email's name field value
func (s Service) UpdateEmail(id, email string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := s.db.Collection("users").UpdateOne(
		s.ctx,
		bson.M{"_id": _id},
		bson.D{{"$set", bson.D{{"email", email}}}},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Updated %v document(s)!\n", result.ModifiedCount)

	return nil
}

//UpdateName updates an existing record's name field value
func (s Service) UpdateName(id, name string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := s.db.Collection("users").UpdateOne(
		s.ctx,
		bson.M{"_id": _id},
		bson.D{{"$set", bson.D{{"name", name}}}},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Updated %v document(s)!\n", result.ModifiedCount)

	return nil
}

func (s Service) updateActive(id string, active bool) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := s.db.Collection("users").UpdateOne(s.ctx, bson.M{"_id": _id},
		bson.D{{"$set", bson.D{{"active", active}}}},
	)

	if err != nil {
		return err
	}

	fmt.Printf("Updated %v document(s)!\n", result.ModifiedCount)

	return nil
}
