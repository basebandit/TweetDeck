package jwt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

const expireOffSet = 3600

//Service defines how a jwt service should behave
type Service interface {
	Create(userID primitive.ObjectID) (token string, err error)
	Verify(token string) (userID primitive.ObjectID, issuedAt time.Time, err error)
	Refresh(oldToken string) (newToken string, err error)
}

type service struct {
	secret []byte
}

//Claims user claims
type Claims struct {
	jwt.StandardClaims
	UserID *primitive.ObjectID `json:"_uid"`
}

//NewService instantiates a new jwt object and returns it.
func NewService(secret string) (Service, error) {
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode jwt secret from hex: %s", err)
	}

	if len(secretBytes) < 32 {
		return nil, errors.New("jwt: secret too short")
	}

	return &service{secret: secretBytes}, nil
}

//Create creates a JWT signing using the given secret key
func (s *service) Create(userID primitive.ObjectID) (string, error) {
	c := Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
		UserID: &userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", errors.New("jwt: token signing failed: " + err.Error())
	}
	return tokenString, nil
}

//Verify verifies the JWT string using the given secret key.
//On success it returns the user ID and the time the token was issued.
func (s *service) Verify(
	tokenString string) (primitive.ObjectID, time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt: unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return primitive.ObjectID{}, time.Time{}, errors.New("jwt: ParseWithClaims failed: " + err.Error())
	}
	if !token.Valid {
		return primitive.ObjectID{}, time.Time{}, errors.New("jwt: token is not valid")
	}
	c, ok := token.Claims.(*Claims)
	if !ok {
		return primitive.ObjectID{}, time.Time{}, errors.New("jwt: failed to get token claims")
	}
	if c.UserID == nil {
		return primitive.ObjectID{}, time.Time{}, errors.New("jwt: UserID claim is not valid")
	}

	if c.IssuedAt == 0 {
		return primitive.ObjectID{}, time.Time{}, errors.New("jwt: IssuedAt claim is not valid")
	}
	return *c.UserID, time.Unix(c.IssuedAt, 0), nil
}

//Refesh returns a new JWT signing using the old token's secret key
func (s *service) Refresh(
	tokenString string) (newToken string, err error) {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt: unexpected signing method")
		}
		return s.secret, nil
	})

	c, ok := token.Claims.(*Claims)
	if !ok {
		return newToken, errors.New("jwt: failed to get token claims")
	}
	if c.UserID == nil {
		return newToken, errors.New("jwt: UserID claim is not valid")
	}

	if c.IssuedAt == 0 {
		return newToken, errors.New("jwt: IssuedAt claim is not valid")
	}
	nTkn, err := s.Create(*c.UserID)
	if err != nil {
		return newToken, err
	}
	return nTkn, nil
}

//GetTokenRemainingValidity calculates and return the remaining time to expiry of a token
func GetTokenRemainingValidity(tokenString string, secret string) int {

	secretBytes, err := hex.DecodeString(secret)

	if err != nil {
		log.Fatal(err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt: unexpected signing method")
		}
		return secretBytes, nil
	})

	if err != nil {
		log.Fatal("GetTokenRemainingValidity", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		tm := time.Unix(claims.ExpiresAt, 0)
		remainder := tm.Sub(time.Now())
		if remainder > 0 {
			return int(remainder.Seconds() + expireOffSet)
		}
	}
	return expireOffSet
}
