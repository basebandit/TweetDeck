package avatar

import "time"

//Avatar is a twitter account profile.
type Avatar struct {
	ID        string    `db:"id" json:"id"`                //Unique Identifier.
	Username  string    `db:"username" json:"username"`    //The twitter handle.
	UserID    string    `db:"user_id" json:"user_id"`      //The user who manages/runs this twitter account.
	Active    bool      `db:"active" json:"active"`        //Use this flag to perform soft deletes.
	CreatedAt time.Time `db:"created_at" json:"createdAt"` //When the record was added.
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"` //When the record was last modified.
}

//NewAvatar is what we require from clients when adding an Avatar.
type NewAvatar struct {
	Username string `json:"username" validate:"required"`
}

//UpdateAvatar defines what information may be provided to modify an
//existing Avatar.All fields are optional so clients can send only
//thos fields they wish to modify.
type UpdateAvatar struct {
	Username *string `json:"username"`
	UserID   *string `json:"userID"`
}
