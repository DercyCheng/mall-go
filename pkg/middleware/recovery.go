package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/response"
)

// RecoveryConfig defines the configuration options for recovery middleware
type RecoveryConfig struct {
	// Log is a function to log the panic
	Log func(v ...interface{})
	// StackAll is a flag to log the full stack trace
	StackAll bool
	// StackSize is the size of the stack to be printed
	StackSize int
}

// DefaultRecoveryConfig returns a default recovery configuration
func DefaultRecoveryConfig() *RecoveryConfig {
	return &RecoveryConfig{
		Log:       log.Println,
		StackAll:  false,
		StackSize: 2048,
	}
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one
func Recovery(config *RecoveryConfig) gin.HandlerFunc {
	// Use default config if none provided
	if config == nil {
		config = DefaultRecoveryConfig()
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				stack := debug.Stack()
				if config.StackAll {
					config.Log("[PANIC RECOVER]", err, string(stack))
				} else {
					if len(stack) > config.StackSize {
						stack = stack[:config.StackSize]
					}
					config.Log("[PANIC RECOVER]", err, string(stack))
				}

				// Return 500 error
				response.InternalServerError(c, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}
