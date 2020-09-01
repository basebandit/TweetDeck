package user

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"ekraal.org/avatarlysis/database"
)

func testService(t *testing.T) (*Service, func(col string)) {
	ctx := context.Background()

	db, err := database.Connect(ctx, "localhost:27017", "avatar", "avatar", "avatars")

	if err != nil {
		t.Fatal(err)
	}

	tearDown := func(col string) {
		err := db.Collection(col).Drop(ctx)
		if err != nil {
			t.Fatalf("failed to drop collection of the test mongo database: %s", err)
		}
	}

	return NewService(ctx, db), tearDown
}

func TestUser(t *testing.T) {
	us, tearDown := testService(t)
	defer tearDown("users")

	u1 := User{Name: "testuser1", Email: "testuser1@gmail.com"}

	id, err := us.Insert(u1.Name, u1.Email)

	if err != nil {
		t.Fatal(err)
	}

	idStr := id.Hex()

	u, err := us.GetByID(idStr)
	if err != nil {
		t.Fatalf("failed to get user by id: %s err: %s", id, err)
	}

	sinceCreated := time.Since(u.CreatedAt)
	if sinceCreated > 3*time.Second || sinceCreated < 0 {
		t.Fatalf("bad user.CreatedAt: %v\n", u.CreatedAt)
	}

	//test for duplicate email
	u3 := User{Name: "testuser1", Email: "testuser1@gmail.com"}

	_, err = us.Insert(u3.Name, u3.Email)

	if err != nil {
		if _, ok := err.(mongo.WriteException); !ok {

			t.Fatalf("expected a mongo.WriteException exception got %T", err)
		}
	}

	want := User{
		ID:        id,
		Email:     "testuser1@gmail.com",
		Name:      "testuser1",
		Password:  u.Password,
		Active:    true,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	//Test GetByID method
	user, err := us.GetByID(idStr)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(user, want) {
		t.Fatalf("got %v expected %v\n", user, want)
	}

	//Test GetByEmail method
	user, err = us.GetByEmail("testuser1@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(user, want) {
		t.Fatalf("got %v expected %v\n", user, want)
	}

	u2 := User{Name: "testuser2", Email: "testuser2@gmail.com"}

	id, err = us.Insert(u2.Name, u2.Email)
	if err != nil {
		t.Fatal(err)
	}

	expect := "testuser3@gmail.com"

	if err := us.UpdateEmail(id.Hex(), expect); err != nil {
		t.Fatal(err)
	}

	u, err = us.GetByID(id.Hex())

	if err != nil {
		t.Fatal(err)
	}

	if u.Email != expect {
		t.Fatalf("got %v expected %v\n", u.Email, expect)
	}

}
