package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RequestID 请求ID中间件，为每个请求添加唯一ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取请求ID，如果没有则生成一个
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将请求ID添加到上下文中
		c.Set("X-Request-ID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算请求处理时间
		latency := time.Since(startTime)

		// 获取请求ID
		requestID, _ := c.Get("X-Request-ID")

		// 记录日志
		method := c.Request.Method
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path
		// 这里可以使用项目的日志库记录请求信息
		// 例如：zap或logrus
		// logger.Info("API请求",
		//    "method", method,
		//    "path", path,
		//    "status", statusCode,
		//    "latency", latency,
		//    "requestID", requestID,
		// )
		
		// 暂时使用gin默认的日志记录
		gin.DefaultWriter.Write([]byte(
			"[API] " + method + " " + path + " " + 
			time.Now().Format("2006/01/02 - 15:04:05") + " " + 
			"| " + string(statusCode) + " | " + 
			latency.String() + " | " + 
			requestID.(string) + "\n"))
	}
}

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}