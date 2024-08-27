package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	TokenIssuer  = "cdb"
	validMethods = []string{"HS512"}
)

// TokenSet is a set of JWT tokens for use as Authentication and Authorisation.
type TokenSet struct {
	IDToken      string
	AccessToken  string
	RefreshToken string
}

func GenerateTokens(signingKey []byte, user User) (TokenSet, error) {
	idToken, idErr := GenerateIdToken(signingKey, user)
	accessToken, accessErr := GenerateAccessToken(signingKey, user)
	refreshToken, refreshErr := GenerateRefreshToken(signingKey)
	return TokenSet{
		IDToken:      idToken,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, errors.Join(idErr, accessErr, refreshErr)
}

func ValidateRefreshToken(signingKey []byte, refreshToken string) (bool, error) {
	_, err := jwt.ParseWithClaims(
		refreshToken,
		nil,
		func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		jwt.WithValidMethods(validMethods),
		jwt.WithIssuer(TokenIssuer),
		jwt.WithExpirationRequired(),
	)
	return err == nil, err
}

func GenerateRefreshToken(signingKey []byte) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer: TokenIssuer,
		// Expires in two days
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 2 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		claims,
	)

	return token.SignedString(signingKey)
}

type AccessClaims struct {
	Email string
	jwt.RegisteredClaims
}

func ValidateAccessToken(signingKey []byte, accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&IDClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		jwt.WithValidMethods(validMethods),
		jwt.WithIssuer(TokenIssuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*AccessClaims); ok {
		return claims.Email, nil
	}

	return "", errors.New("unrecognized token claims")
}

func GenerateAccessToken(signingKey []byte, user User) (string, error) {
	claims := AccessClaims{
		Email: user.Email,
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

	return token.SignedString(signingKey)
}

type IDClaims struct {
	User
	jwt.RegisteredClaims
}

func ValidateIdToken(signingKey []byte, idToken string) (User, error) {
	token, err := jwt.ParseWithClaims(
		idToken,
		&IDClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		jwt.WithValidMethods(validMethods),
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

func GenerateIdToken(signingKey []byte, user User) (string, error) {
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

	return token.SignedString(signingKey)
}
