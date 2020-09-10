package auth

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//expected values for Claims.Roles.
const (
	//RoleAdmin represents the admin role.
	RoleAdmin = "ADMIN"
	//RoleUser represents the user role.
	RoleUser = "USER"
)

//ctxKey represents the type of value for the context key.
type ctxKey int

//Key is used to store/retrieve a Claims value from a context.Context.
const Key ctxKey = 1

//Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

//Valid is called during the parsing of a token.
func (c Claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return errors.Wrap(err, "validating standard claims")
	}
	return nil
}

//HasRole returns true if the claims has at least one of the provided roles.
func (c Claims) HasRole(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}

//KeyLookupFunc defines the signature of a function to lookup public keys.
//A key lookup function is required for creating an Authenticator.
// * Private keys should be rotated.
type KeyLookupFunc func(publicKID string) (*rsa.PublicKey, error)

//Auth is used to authenticate clients. It can generate a token for a
// set of user claims and recreate the claims by parsing the token.
type Auth struct {
	privateKey       *rsa.PrivateKey
	publicKID        string
	algorithm        string
	pubKeyLookupFunc KeyLookupFunc
	parser           *jwt.Parser
}

//New creates an *Authenticator for use. It will error if:
// - The private key is nil.
// - The public key ID is empty.
// - The specified algorithm is unsupported.
// - The public key function is nil.
func New(privateKey *rsa.PrivateKey, publicKID string, algorithm string, publicKeyLookupFunc KeyLookupFunc) (*Auth, error) {
	if privateKey == nil {
		return nil, errors.New("private key cannot be nil")
	}

	if publicKID == "" {
		return nil, errors.New("public kid cannot be blank")
	}

	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}

	if publicKeyLookupFunc == nil {
		return nil, errors.New("public key function cannot be nil")
	}

	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	a := Auth{
		privateKey:       privateKey,
		publicKID:        publicKID,
		algorithm:        algorithm,
		pubKeyLookupFunc: publicKeyLookupFunc,
		parser:           &parser,
	}

	return &a, nil
}

//GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(claims Claims) (string, error) {
	method := jwt.GetSigningMethod(a.algorithm)

	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = a.publicKID

	str, err := tkn.SignedString(a.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return str, nil
}

//ValidateToken recreates the Claims that were used to generate a token. It verifies that the token was signed using our key.
func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id (kid) in token header")
		}

		publicKID, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id (kid) must be string")
		}

		return a.pubKeyLookupFunc(publicKID)
	}

	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &claims, keyFunc)

	if err != nil {
		return Claims{}, errors.Wrap(err, "parsing token")
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	return claims, nil
}
