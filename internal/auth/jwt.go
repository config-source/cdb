package auth

import (
	"errors"
	"time"

	"github.com/config-source/cdb/internal/settings"
	"github.com/golang-jwt/jwt/v5"
)

const TokenIssuer = "cdb"

// TokenSet is a set of JWT tokens for use as Authentication and Authorisation.
type TokenSet struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
}

func GenerateTokens(user User) (TokenSet, error) {
	return TokenSet{}, nil
}

type IDClaims struct {
	User
	jwt.RegisteredClaims
}

func ValidateIdToken(idToken string) (User, error) {
	token, err := jwt.ParseWithClaims(
		idToken,
		&IDClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return settings.JWTSigningKey(), nil
		},
		jwt.WithValidMethods([]string{"HS512"}),
		jwt.WithIssuer(TokenIssuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return User{}, err
	}

	if claims, ok := token.Claims.(*IDClaims); ok {
		return claims.User, nil
	}

	return User{}, errors.New("unrecognized token claims")
}

func GenerateIdToken(user User) (string, error) {
	claims := IDClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		claims,
	)

	return token.SignedString(settings.JWTSigningKey())
}
