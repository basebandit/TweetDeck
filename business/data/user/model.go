package user

import "time"

//User represents someone with access to our system.
type User struct {
	ID           string    `db:"id" json:"id"`
	Firstname    string    `db:"firstname" json:"firstname"`
	Lastname     string    `db:"lastname" json:"lastname"`
	Email        string    `db:"email" json:"email"`
	PasswordHash []byte    `db:"password_hash" json:"-"` //password hash
	Active       bool      `db:"active" json:"active"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

//NewUser contains information needed to create a new user.
type NewUser struct {
	Firstname       string  `json:"firstname" validate:"required"`
	Lastname        string  `json:"lastname" validate:"required"`
	Password        *string `json:"password"` //optional because we will be adding users who cannot login.
	PasswordConfirm *string `json:"passwordConfirm" validate:"eqfield=Password"`
	Email           *string `json:"email"` //optional because some users don't have this.
}

//UpdateUser defines what information may be provided to modify an existing
//User. All fields are optional so client can send just the fields they want
//changed. It uses a pointer fields so that we can differentiate between a field that
//was not provided and a field that was provided as explicitly blank.
type UpdateUser struct {
	Firstname       *string `json:"firstname"`
	Lastname        *string `json:"lastname"`
	Email           *string `json:"email"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
}
