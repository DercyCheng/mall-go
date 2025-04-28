package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte

// Claims JWT声明结构
type Claims struct {
	Username string `json:"username"`
	UserID   uint   `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(username string, userId uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)

	claims := Claims{
		Username: username,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "mall-go",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 解析JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, errors.New("invalid token")
}

// InitJWTSecret 初始化JWT密钥
func InitJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}