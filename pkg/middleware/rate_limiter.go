package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig defines the configuration for the rate limiter middleware
type RateLimiterConfig struct {
	// Limit is the maximum number of requests allowed in the period
	Limit int
	// Period is the time window for the limit
	Period time.Duration
	// KeyFunc is a function that returns a key for rate limiting
	KeyFunc func(*gin.Context) string
}

// DefaultRateLimiterConfig returns a default rate limiter configuration
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		Limit:  100,         // 100 requests
		Period: time.Minute, // per minute
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP() // Use client IP as key
		},
	}
}

// SimpleTokenBucket implements a simple token bucket algorithm
type SimpleTokenBucket struct {
	mu         sync.Mutex
	tokens     int
	capacity   int
	fillRate   float64
	lastFilled time.Time
}

// NewSimpleTokenBucket creates a new token bucket
func NewSimpleTokenBucket(capacity int, period time.Duration) *SimpleTokenBucket {
	return &SimpleTokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		fillRate:   float64(capacity) / float64(period.Seconds()),
		lastFilled: time.Now(),
	}
}

// Allow checks if a request is allowed and decrements the token count
func (tb *SimpleTokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(tb.lastFilled).Seconds()
	tb.lastFilled = now

	// Calculate tokens to add
	tokensToAdd := elapsed * tb.fillRate
	if tokensToAdd > 0 {
		tb.tokens += int(tokensToAdd)
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
	}

	// Check if we have tokens
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// SimpleRateLimiter implements a basic in-memory rate limiter
type SimpleRateLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*SimpleTokenBucket
	config  *RateLimiterConfig
}

// NewSimpleRateLimiter creates a new rate limiter
func NewSimpleRateLimiter(config *RateLimiterConfig) *SimpleRateLimiter {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	return &SimpleRateLimiter{
		buckets: make(map[string]*SimpleTokenBucket),
		config:  config,
	}
}

// Allow checks if a request is allowed
func (rl *SimpleRateLimiter) Allow(key string) bool {
	rl.mu.RLock()
	bucket, exists := rl.buckets[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Check again in case another goroutine created it
		bucket, exists = rl.buckets[key]
		if !exists {
			bucket = NewSimpleTokenBucket(rl.config.Limit, rl.config.Period)
			rl.buckets[key] = bucket
		}
		rl.mu.Unlock()
	}

	return bucket.Allow()
}

// RateLimiter returns a middleware that limits the number of requests
func RateLimiter(config *RateLimiterConfig) gin.HandlerFunc {
	// Use default config if none provided
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	// Create rate limiter
	limiter := NewSimpleRateLimiter(config)

	return func(c *gin.Context) {
		// Get key for rate limiting
		key := config.KeyFunc(c)

		// Check if the request is allowed
		if !limiter.Allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
