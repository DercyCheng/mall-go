package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	Issuer     string
	ExpireTime int // 过期时间（分钟）
}

// CustomClaims 自定义JWT Claims
type CustomClaims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTMiddleware JWT中间件
type JWTMiddleware struct {
	Config JWTConfig
}

// NewJWTMiddleware 创建JWT中间件
func NewJWTMiddleware(config JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{
		Config: config,
	}
}

// Authenticate JWT认证中间件
func (m *JWTMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// 解析token
		tokenString := parts[1]
		claims, err := m.parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// 解析JWT Token
func (m *JWTMiddleware) parseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.Config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// 检查令牌是否已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// GenerateToken 生成JWT令牌
func (m *JWTMiddleware) GenerateToken(userID, username, role string) (string, error) {
	now := time.Now()
	expireTime := now.Add(time.Duration(m.Config.ExpireTime) * time.Minute)

	claims := &CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.Config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.Config.Secret))
}

// AdminRequired 管理员权限检查中间件
func (m *JWTMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized"})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Admin permission required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RefreshToken 刷新令牌
func (m *JWTMiddleware) RefreshToken(tokenString string) (string, error) {
	claims, err := m.parseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新的令牌
	return m.GenerateToken(claims.UserID, claims.Username, claims.Role)
}
