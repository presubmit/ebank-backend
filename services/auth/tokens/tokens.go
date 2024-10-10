package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type AccessToken struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}

type RefreshToken struct {
	ID            string
	UserID        string
	AccessTokenID string
	ExpiresAt     time.Time
}

func (at *AccessToken) GenerateJWT(secret []byte) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": at.ID,
		"exp":  at.ExpiresAt.Unix(),
	}).SignedString(secret)
}

func (rt *RefreshToken) GenerateJWT(secret []byte) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid":        rt.ID,
		"access_uuid": rt.AccessTokenID,
		"exp":         rt.ExpiresAt.Unix(),
	}).SignedString(secret)
}

func parse(tkn string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		// make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	// extract claims from jwt
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

func ParseRefreshToken(tkn string, secret []byte) (*RefreshToken, error) {
	claims, err := parse(tkn, secret)
	if err != nil {
		return nil, err
	}
	var found bool
	rt := &RefreshToken{}
	// extract token's uuid
	if rt.ID, found = claims["uuid"].(string); !found {
		return nil, errors.New("invalid token uuid")
	}
	// extract token's access uuid
	if rt.AccessTokenID, found = claims["access_uuid"].(string); !found {
		return nil, errors.New("invalid token uuid")
	}
	return rt, nil
}

func ParseAccessToken(tkn string, secret []byte) (*AccessToken, error) {
	claims, err := parse(tkn, secret)
	if err != nil {
		return nil, err
	}
	var found bool
	at := &AccessToken{}
	// extract token's uuid
	if at.ID, found = claims["uuid"].(string); !found {
		return nil, errors.New("invalid token uuid")
	}
	return at, nil
}
