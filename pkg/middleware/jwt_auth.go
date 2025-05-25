package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/auth"
	"mall-go/pkg/response"
)

// JWTConfig defines the configuration for the JWT middleware
type JWTConfig struct {
	// Secret is the key used to sign the JWT
	Secret string
	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string
	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string
	// ContextKey is the key used to store user information in the context
	ContextKey string
}

// DefaultJWTConfig returns a default JWT configuration
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		ContextKey:    "user",
	}
}

// JWTAuth returns a middleware that authorizes JWT tokens
func JWTAuth(config *JWTConfig) gin.HandlerFunc {
	// Use default config if none provided
	if config == nil {
		config = DefaultJWTConfig()
	}

	// Create JWT service
	jwtConfig := &auth.Config{
		Secret: config.Secret,
	}
	jwtService := auth.NewService(jwtConfig)

	// Parse token lookup configuration
	parts := strings.Split(config.TokenLookup, ":")
	extractor := jwtFromHeader
	param := "Authorization"

	if len(parts) == 2 {
		switch parts[0] {
		case "header":
			extractor = jwtFromHeader
			param = parts[1]
		case "query":
			extractor = jwtFromQuery
			param = parts[1]
		case "cookie":
			extractor = jwtFromCookie
			param = parts[1]
		}
	}

	// Return the middleware
	return func(c *gin.Context) {
		// Extract token
		tokenString, err := extractor(c, param)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// Remove token head if present
		if len(config.TokenHeadName) > 0 {
			tokenPrefix := config.TokenHeadName + " "
			if strings.HasPrefix(tokenString, tokenPrefix) {
				tokenString = tokenString[len(tokenPrefix):]
			}
		}

		// Validate token
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "无效的认证令牌")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

var (
	ErrMissingToken = errors.New("authorization token is missing")
)

// jwtFromHeader returns a function that extracts the JWT from the Authorization header
func jwtFromHeader(c *gin.Context, param string) (string, error) {
	token := c.GetHeader(param)
	if token == "" {
		return "", ErrMissingToken
	}
	return token, nil
}

// jwtFromQuery returns a function that extracts the JWT from the query string
func jwtFromQuery(c *gin.Context, param string) (string, error) {
	token := c.Query(param)
	if token == "" {
		return "", ErrMissingToken
	}
	return token, nil
}

// jwtFromCookie returns a function that extracts the JWT from the cookie
func jwtFromCookie(c *gin.Context, param string) (string, error) {
	token, err := c.Cookie(param)
	if err != nil {
		return "", ErrMissingToken
	}
	return token, nil
}
