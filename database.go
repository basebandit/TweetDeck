package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//connects to a mongodb and returns a connection handle
//address: localhost:27017
func connect(ctx context.Context, address, username, password, database string) (*mongo.Database, error) {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s/?authSource=%s&connectTimoutMS=300000", username, password, address, database)

	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(database)

	return db, nil
}

func insertRecord(ctx context.Context, db *mongo.Database, avatars []avatar, forDate time.Time) (interface{}, error) {
	r := record{
		ID:        primitive.NewObjectID(),
		ForDate:   forDate,
		CreatedAt: time.Now(),
		Avatars:   avatars,
	}

	result, err := db.Collection("records").InsertOne(ctx, r)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
