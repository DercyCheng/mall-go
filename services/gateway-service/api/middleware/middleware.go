package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"mall-go/services/gateway-service/infrastructure/config"
)

// SetupMiddleware configures all middleware for the gateway
func SetupMiddleware(router *gin.Engine, cfg *config.Config) {
	// Setup CORS
	setupCORS(router, cfg.CORS)

	// Setup rate limiting if enabled
	if cfg.RateLimit.Enabled {
		setupRateLimit(router, cfg.RateLimit)
	}

	// Add request logging middleware
	router.Use(requestLogger())
}

// setupCORS configures CORS for the gateway
func setupCORS(router *gin.Engine, cfg config.CORSConfig) {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	}
	router.Use(cors.New(corsConfig))
}

// setupRateLimit configures rate limiting for the gateway
func setupRateLimit(router *gin.Engine, cfg config.RateLimitConfig) {
	// Create a rate limiter with the configured values
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  int64(cfg.RequestsPerSecond),
	}

	// Create a memory store with the configured burst size
	store := memory.NewStore()

	// Create the rate limiter
	instance := limiter.New(store, rate)

	// Add the rate limiting middleware
	router.Use(func(c *gin.Context) {
		// Get the client IP address as the rate limiting key
		key := c.ClientIP()

		// Check if the request should be limited
		context, err := instance.Get(c, key)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error": "Error checking rate limit",
			})
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", "100")
		c.Header("X-RateLimit-Reset", "60")

		// If the request exceeds the rate limit, return 429 Too Many Requests
		if context.Reached {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		c.Next()
	})
}

// requestLogger creates a middleware for logging requests
func requestLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	})
}
