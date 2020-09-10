package profile_test

import (
	"testing"
	"time"

	"ekraal.org/avatarlysis/business/data/auth"
	"ekraal.org/avatarlysis/business/data/profile"

	"ekraal.org/avatarlysis/business/data/avatar"
	"ekraal.org/avatarlysis/business/data/tests"
	"ekraal.org/avatarlysis/business/data/user"
	"github.com/google/uuid"
)

func TestProfile(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	t.Log("Given the need to work with Profile table.")
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single Profile.", testID)
		{
			now := time.Date(2020, time.September, 9, 0, 0, 0, 0, time.UTC)
			ctx := tests.Context()

			//we need to first create an avatar record because we need avatar_id in profiles table
			//an avatar record in turns needs a user_id.
			//we will need this user's id
			nu := user.NewUser{
				Firstname:       "Evanson",
				Lastname:        "Mwangi",
				Roles:           []string{auth.RoleAdmin},
				Password:        tests.StringPointer("gophers"),
				PasswordConfirm: tests.StringPointer("gophers"),
			}

			u, err := user.Create(ctx, db, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create user : %s", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

			na := avatar.NewAvatar{
				Username: "the_basebandit",
			}

			a, err := avatar.Create(ctx, db, na, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create an Avatar : %s.", tests.Failed, testID, err)
			}

			updA := avatar.UpdateAvatar{
				UserID: tests.StringPointer(u.ID),
			}

			updatedTime := time.Date(2020, time.September, 9, 1, 1, 1, 0, time.UTC)

			if err := avatar.Update(ctx, db, a.ID, updA, updatedTime); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update Avatar.", tests.Success, testID)

			np := profile.NewProfile{
				ID:              uuid.New().String(),
				Bio:             tests.StringPointer("Cuppycake\n\nLiving large"),
				Followers:       tests.IntPointer(450),
				Following:       tests.IntPointer(600),
				Tweets:          tests.IntPointer(2304),
				Likes:           tests.IntPointer(3042),
				Name:            tests.StringPointer("Jean Wangari"),
				JoinDate:        tests.StringPointer("2020-06-02 09:59:17"),
				ProfileImageURL: tests.StringPointer("https://pbs.twimg.com/profile_images/1288401307204804608/0s5DK5ej.jpg"),
				LastTweetTime:   tests.StringPointer("2020-08-22 18:55:03"),
				TwitterID:       tests.StringPointer("1267757177999101953"),
			}

			if err := profile.Create(ctx, db, a.ID, &np, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create Profile : %s", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create Profile.", tests.Success, testID)
		}
	}
}
