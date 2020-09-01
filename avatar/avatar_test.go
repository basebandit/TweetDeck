package avatar

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"ekraal.org/avatarlysis/database"
)

func testService(t *testing.T) (*Service, func(col string)) {
	ctx := context.Background()

	db, err := database.Connect(ctx, "localhost:27017", "avatar", "avatar", "avatars")
	if err != nil {
		t.Fatal(err)
	}

	tearDown := func(col string) {
		err := db.Collection("avatars").Drop(ctx)
		if err != nil {
			t.Fatalf("failed to drop collection of the test mongo database: %s", err)
		}
	}
	return NewService(ctx, db), tearDown
}

func TestAvatar(t *testing.T) {
	as, tearDown := testService(t)
	defer tearDown("avatars")

	// a1 := Avatar{
	// 	TwitterCreatedAt: "Tue Jun 02 10:32:52 +0000 2020",
	// 	Bio:              "Gíkúyú ní wendo",
	// 	Likes:            87,
	// 	Followers:        238,
	// 	Following:        626,
	// 	TwitterID:        "1267766034179776512",
	// 	Location:         "Nakuru, Kenya",
	// 	Name:             "DK Jnr.",
	// 	ProfileImage:     "http://pbs.twimg.com/profile_images/1293103745648201729/CMwG39AN_normal.jpg",
	// 	Username:         "DKJnr3",
	// 	Tweets:           486,
	// 	CreatedAt:        time.Now(),
	// }

	uid := primitive.NewObjectID()

	id, err := as.Insert(uid.Hex(), "DKJnr3")
	if err != nil {
		t.Fatal(err)
	}

	idStr := id.Hex()

	a, err := as.GetByID(idStr)
	if err != nil {
		t.Fatal(err)
	}

	sinceCreated := time.Since(a.CreatedAt)
	if sinceCreated > 3*time.Second || sinceCreated < 0 {
		t.Fatalf("bad user.CreatedAt: %v\n", a.CreatedAt)
	}

	all, err := as.GetAll()

	if len(all) == 0 {
		t.Fatal("no avatars retrieved,expected length to be greater than 1")
	}

}
