package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

//User defines a user model
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Password  string             `bson:"password"`
	Active    bool               `bson:"active"`
	Email     string             `bson:"email"`
	CreatedAt time.Time          `bson:"createAt"`
	UpdatedAt time.Time          `bson:"updateat"`
}

//hashes the user passwor using bcrypt hash function
func (u *User) hashPassword() error {
	pwd := []byte(u.Password)

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

//Compare compares the password hash against the passed in password string
func (u User) Compare(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
