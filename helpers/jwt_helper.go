package helpers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// JWTHelper is an implementation of the TokenHelper interface that uses JSON Web Tokens.
type JWTHelper struct{}

// SessionToken is a custom JWT claim that includes a session id.
type SessionToken struct {
	SID uuid.UUID
	jwt.StandardClaims
}

// CreateSessionToken creates a JWT with a session id to be used for authenticating requests.
func (JWTHelper) CreateSessionToken(sID uuid.UUID) (string, error) {
	//TODO: load key from config
	signingKey := []byte("key")

	//create the session token
	sToken := SessionToken{
		SID: sID,
	}

	//create the jwt and sign it
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sToken)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", ChainError("error signing session token", err)
	}

	return tokenString, nil
}
