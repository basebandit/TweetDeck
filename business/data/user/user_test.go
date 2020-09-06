package user_test

import (
	"testing"
	"time"

	"ekraal.org/avatarlysis/business/data/schema"
	"ekraal.org/avatarlysis/business/data/user"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	"ekraal.org/avatarlysis/business/data/tests"
)

func TestUser(t *testing.T) {
	db, tearDown := tests.NewUnit(t)
	t.Cleanup(tearDown)

	t.Log("Given the need to work with User records.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := tests.Context()
			now := time.Date(2020, time.September, 6, 0, 0, 0, 0, time.UTC)

			nu := user.NewUser{
				Firstname:       "Evanson",
				Lastname:        "Mwangi",
				Password:        tests.StringPointer("gophers"),
				PasswordConfirm: tests.StringPointer("gophers"),
			}

			if err := schema.DeleteAll(db); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete all data : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete all data.", tests.Success, testID)

			u, err := user.Create(ctx, db, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			savedU, err := user.GetByID(ctx, db, u.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user by ID.", tests.Success, testID)

			if diff := cmp.Diff(u, savedU); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same user.", tests.Success, testID)

			updU := user.UpdateUser{
				Firstname: tests.StringPointer("Mr"),
				Lastname:  tests.StringPointer("Parish"),
				Email:     tests.StringPointer("parish@nsynclabs.com"),
			}

			if err := user.Update(ctx, db, u.ID, updU, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update user.", tests.Success, testID)

			savedU, err = user.GetByID(ctx, db, u.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve user by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve user.", tests.Success, testID)

			if savedU.Firstname != *updU.Firstname {
				t.Errorf("\t%s\tTest %d:\tShould be able to see update to Firstname.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, savedU.Firstname)
				t.Logf("\t\tTest %d:\tWant: %v", testID, *updU.Firstname)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see update to Firstname.", tests.Success, testID)
			}

			if savedU.Lastname != *updU.Lastname {
				t.Errorf("\t%s\tTest %d:\tShould be able to see update to Lastname.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, savedU.Lastname)
				t.Logf("\t\tTest %d:\tWant: %v", testID, *updU.Lastname)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see update to Lastname.", tests.Success, testID)
			}

			if savedU.Email != *updU.Email {
				t.Errorf("\t%s\tTest %d:\tShould be able to see update to Email.", tests.Failed, testID)
				t.Logf("\t\tTest %d:\tGot: %v", testID, savedU.Email)
				t.Logf("\t\tTest %d:\tWant: %v", testID, *updU.Email)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see update to Email.", tests.Success, testID)
			}

			if err := user.Delete(ctx, db, u.ID); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete user.", tests.Success, testID)

			_, err = user.GetByID(ctx, db, u.ID)
			if errors.Cause(err) != user.ErrNotFound {
				t.Fatalf("\t%s\tTests %d:\tShould NOT be able to retrieve user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve user.", tests.Success, testID)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	t.Log("Given the need to authenticate users")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single User.", testID)
		{
			ctx := tests.Context()

			nu := user.NewUser{
				Firstname:       "Adidja",
				Lastname:        "Palmer",
				Email:           tests.StringPointer("adi@worldboss.org"),
				Password:        tests.StringPointer("goroutines"),
				PasswordConfirm: tests.StringPointer("goroutines"),
			}

			now := time.Date(2020, time.September, 6, 0, 0, 0, 0, time.UTC)

			if err := schema.DeleteAll(db); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to delete all data : %s", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to delete all data.", tests.Success, testID)

			_, err := user.Create(ctx, db, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			if err := user.Authenticate(ctx, db, "adi@worldboss.org", "goroutines"); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to authenticate given password : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to authenticate given password.", tests.Success, testID)
		}
	}
}
