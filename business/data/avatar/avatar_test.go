package avatar_test

import (
	"context"
	"testing"
	"time"

	"ekraal.org/avatarlysis/business/data/avatar"
	"ekraal.org/avatarlysis/business/data/tests"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

func TestAvatar(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	t.Cleanup(teardown)

	t.Log("Given the need to work with Avatar records.")
	{
		testID := 0
		// testUserID := "45b5fbd3-755f-4379-8f07-a58d4a30fa2f"
		t.Logf("\tTest %d:\tWhen handling a single Avatar.", testID)
		{
			na := avatar.NewAvatar{
				Username: "the_basebandit",
			}

			now := time.Date(2020, time.September, 7, 0, 0, 0, 0, time.UTC)
			ctx := context.Background()

			a, err := avatar.Create(ctx, db, na, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create an Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create an Avatar.", tests.Success, testID)

			saved, err := avatar.GetByID(ctx, db, a.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve an Avatar by ID: %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve an Avatar by ID.", tests.Success, testID)

			if diff := cmp.Diff(a, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same Avatar.Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%sTest\t %d:\tShould get back the same Avatar.", tests.Success, testID)

			updA := avatar.UpdateAvatar{
				Username: tests.StringPointer("Basebandit"),
				UserID:   tests.StringPointer("6ba7b812-9dad-11d1-80b4-00c04fd430c8"),
			}

			updatedTime := time.Date(2020, time.September, 7, 1, 1, 1, 0, time.UTC)

			if err := avatar.Update(ctx, db, a.ID, updA, updatedTime); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update Avatar.", tests.Success, testID)

			saved, err = avatar.GetByID(ctx, db, a.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve updated Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve updated Avatar.", tests.Success, testID)

			want := a
			want.Username = *updA.Username
			want.CreatedAt = updatedTime

			if diff := cmp.Diff(want, saved); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same Avatar. Diff:\n%s", tests.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same Avatar.", tests.Success, testID)

			updA = avatar.UpdateAvatar{
				Username: tests.StringPointer("Lure_Strings"),
			}

			if err := avatar.Update(ctx, db, a.ID, updA, updatedTime); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update just some of fields of Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update just some of the fields of Avatar.", tests.Success, testID)

			saved, err = avatar.GetByID(ctx, db, a.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve updated Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve updated Avatar.", tests.Success, testID)

			if saved.Username != *updA.Username {
				t.Fatalf("\t%s\tTest %d:\tShould be able to see updated Username field : got %q want %q.", tests.Failed, testID, saved.Username, *updA.Username)
			} else {
				t.Logf("\t%s\tTest %d:\tShould be able to see updated Username field.", tests.Success, testID)
			}

			if err := avatar.Delete(ctx, db, a.ID, now); err != nil {
				t.Fatalf("\t%s\t %d:\tShould be able to delete Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%sTest %d:\tShould be able to delete Avatar.", tests.Success, testID)

			_, err = avatar.GetByID(ctx, db, a.ID)
			if errors.Cause(err) != avatar.ErrNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve deleted Avatar : %s.", tests.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould NOT be able to retrieve deleted Avatar.", tests.Success, testID)
		}
	}
}
