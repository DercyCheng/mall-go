// Package middleware provides common middleware components for HTTP services
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// CORSConfig defines the configuration options for CORS middleware
type CORSConfig struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from
	AllowedOrigins []string
	// AllowedMethods is a list of methods the client is allowed to use with cross-domain requests
	AllowedMethods []string
	// AllowedHeaders is a list of headers the client is allowed to use with cross-domain requests
	AllowedHeaders []string
	// AllowCredentials indicates whether the request can include user credentials
	AllowCredentials bool
	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached
	MaxAge time.Duration
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           86400 * time.Second, // 24 hours
	}
}

// CORS returns a middleware that handles CORS requests
func CORS(config *CORSConfig) gin.HandlerFunc {
	// Use default config if none provided
	if config == nil {
		config = DefaultCORSConfig()
	}

	// Convert string slice to string for header
	allowOrigins := "*"
	if len(config.AllowedOrigins) > 0 && config.AllowedOrigins[0] != "*" {
		allowOrigins = config.AllowedOrigins[0]
		for i := 1; i < len(config.AllowedOrigins); i++ {
			allowOrigins += ", " + config.AllowedOrigins[i]
		}
	}

	allowMethods := "GET, POST, PUT, DELETE, OPTIONS"
	if len(config.AllowedMethods) > 0 {
		allowMethods = config.AllowedMethods[0]
		for i := 1; i < len(config.AllowedMethods); i++ {
			allowMethods += ", " + config.AllowedMethods[i]
		}
	}

	allowHeaders := "Origin, Authorization, Content-Type, X-Requested-With"
	if len(config.AllowedHeaders) > 0 {
		allowHeaders = config.AllowedHeaders[0]
		for i := 1; i < len(config.AllowedHeaders); i++ {
			allowHeaders += ", " + config.AllowedHeaders[i]
		}
	}

	// Return the middleware
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowOrigins)
		c.Header("Access-Control-Allow-Methods", allowMethods)
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
