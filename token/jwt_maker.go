package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// JSON web token maker
type JWTMaker struct {
	secretKey string
}

// create a new JWT
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

// creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", payload, err
	}

	// create jwt
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	token, err := jwtToken.SignedString([]byte(maker.secretKey))

	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		// does not math with our signin algorithm
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		// check for specific validation errors
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			fmt.Println("token expired")
			return nil, errors.New("token expired")

		case errors.Is(err, jwt.ErrTokenNotValidYet):
			fmt.Println("Token not valid yet")
			return nil, errors.New("token not valid yet")

		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			fmt.Println("Invalid signature")
			return nil, errors.New("invalid signature")

		case errors.Is(err, jwt.ErrTokenMalformed):
			fmt.Println("malformed token")
			return nil, errors.New("malformed token")

		default:
			fmt.Println("other error:", err)
		}
		return nil, ErrInvalidToken
	}

	// the token is good, now get it's payload data
	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
