package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims represents JWT claims with user information
type CustomClaims struct {
	UserID   string            `json:"user_id"`
	Username string            `json:"username"`
	Roles    []string          `json:"roles,omitempty"`
	CustomClaims map[string]interface{} `json:"custom_claims,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token
func GenerateToken(userID, username string, roles []string, secret string, expiresAt time.Time) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		CustomClaims: make(map[string]interface{}),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mall-go",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken validates and parses a JWT token
func ParseToken(tokenString, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// TokenClaims contains data to be included in token
type TokenClaims struct {
	UserID   string
	Username string
	Roles    []string
	ExpiresAt time.Time
}

// GenerateTokenWithCustomClaims creates a JWT with additional custom claims
func GenerateTokenWithCustomClaims(tokenClaims TokenClaims, customClaims map[string]interface{}, secret string) (string, error) {
	claims := CustomClaims{
		UserID:       tokenClaims.UserID,
		Username:     tokenClaims.Username,
		Roles:        tokenClaims.Roles,
		CustomClaims: customClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenClaims.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mall-go",
			Subject:   tokenClaims.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
