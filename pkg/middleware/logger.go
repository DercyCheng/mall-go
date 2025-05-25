package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func logAdapter(msg string) {
	log.Printf("%s", msg)
}

// LoggerConfig defines the configuration for the logger middleware
type LoggerConfig struct {
	// SkipPaths is a list of paths to skip when logging
	SkipPaths []string
	// TimeFormat is the format to use for the time
	TimeFormat string
	// UTC indicates whether to use UTC time
	UTC bool
	// Output is a function to output the log
	Output func(string)
}

// DefaultLoggerConfig returns a default logger configuration
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		SkipPaths:  []string{"/health", "/metrics"},
		TimeFormat: "2006/01/02 - 15:04:05",
		UTC:        false,
		Output:     logAdapter,
	}
}

// Logger returns a middleware that logs request information
func Logger(config *LoggerConfig) gin.HandlerFunc {
	// Use default config if none provided
	if config == nil {
		config = DefaultLoggerConfig()
	}

	// Create a map for faster lookup of skip paths
	skipPaths := make(map[string]bool, len(config.SkipPaths))
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}

	// Return the middleware
	return func(c *gin.Context) {
		// Skip logging for certain paths
		path := c.Request.URL.Path
		if _, exists := skipPaths[path]; exists {
			c.Next()
			return
		}

		// Start time
		start := time.Now()
		if config.UTC {
			start = start.UTC()
		}

		// Process request
		c.Next()

		// End time
		end := time.Now()
		if config.UTC {
			end = end.UTC()
		}
		latency := end.Sub(start)

		// Log format
		logLine := fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %s %s",
			end.Format(config.TimeFormat),
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
		)

		// Output log
		config.Output(logLine)
	}
}
