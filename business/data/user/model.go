package user

import (
	"fmt"
	"time"

	b64 "encoding/base64"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

//User represents someone with access to our system.
type User struct {
	ID           uuid.UUID      `db:"id" json:"-"`
	UID          string         `json:"id"` //encoded/shortened id string
	Firstname    string         `db:"firstname" json:"firstname"`
	Lastname     string         `db:"lastname" json:"lastname"`
	Email        string         `db:"email" json:"email"`
	PasswordHash []byte         `db:"password_hash" json:"-"` //password hash
	Roles        pq.StringArray `db:"roles" json:"roles"`
	Active       bool           `db:"active" json:"active"`
	CreatedAt    time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updatedAt"`
}

//NewUser contains information needed to create a new user.
type NewUser struct {
	Firstname       string   `json:"firstname" validate:"required"`
	Lastname        string   `json:"lastname" validate:"required"`
	Password        *string  `json:"password"` //optional because we will be adding users who cannot login.
	Roles           []string `json:"roles" validate:"required"`
	PasswordConfirm *string  `json:"passwordConfirm" validate:"eqfield=Password"`
	Email           *string  `json:"email"` //optional because some users don't have this.
}

//UpdateUser defines what information may be provided to modify an existing
//User. All fields are optional so client can send just the fields they want
//changed. It uses a pointer fields so that we can differentiate between a field that
//was not provided and a field that was provided as explicitly blank.
type UpdateUser struct {
	Firstname       *string  `json:"firstname"`
	Lastname        *string  `json:"lastname"`
	Email           *string  `json:"email"`
	Roles           []string `json:"roles"`
	Password        *string  `json:"password"`
	PasswordConfirm *string  `json:"passwordConfirm" validate:"omitempty,eqfield=Password"`
}

//Encode encodes the id to a browser friendly short format
func Encode(id uuid.UUID) string {
	_id, err := id.MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}
	return b64.RawURLEncoding.EncodeToString(_id)
}

//Decode decodes the user ID back to the original internal format
func Decode(id string) (*uuid.UUID, error) {
	dec, err := b64.RawURLEncoding.DecodeString(id)

	if err != nil {
		return nil, err
	}

	decoded, err := uuid.FromBytes(dec)
	if err != nil {
		return nil, err
	}

	return &decoded, nil
}
