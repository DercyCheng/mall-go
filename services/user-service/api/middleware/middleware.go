package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/auth"
	"mall-go/pkg/errors"
	"mall-go/pkg/response"
)

// JWTAuthMiddleware 创建JWT认证中间件
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	// 创建JWT服务
	jwtConfig := &auth.Config{
		Secret: secret,
	}
	jwtService := auth.NewService(jwtConfig)

	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		// 处理Bearer令牌格式
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Unauthorized(c, "认证令牌格式不正确")
			c.Abort()
			return
		}

		tokenString := authHeader[len(bearerPrefix):]

		// 验证令牌
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			appErr, ok := errors.As(err)
			if ok && appErr.Type == errors.TypeUnauthorized {
				response.Error(c, appErr)
			} else {
				response.Unauthorized(c, "无效的认证令牌")
			}
			c.Abort()
			return
		}

		// 在上下文中设置用户信息
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

// Cors 处理跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			// 设置允许的域名
			c.Header("Access-Control-Allow-Origin", origin)
			// 允许的请求方法
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// 允许的请求头
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			// 允许暴露的响应头
			c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin")
			// 允许携带凭证
			c.Header("Access-Control-Allow-Credentials", "true")

			// 处理预检请求
			if method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

// RateLimiter 限制请求频率
func RateLimiter() gin.HandlerFunc {
	// 这里可以使用redis或内存实现一个限流器
	// 为简化示例，这里只实现一个基本框架
	return func(c *gin.Context) {
		// TODO: 实现真正的限流逻辑
		c.Next()
	}
}

// RequestLogger 记录请求日志
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取状态码
		status := c.Writer.Status()

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取请求方法和路径
		method := c.Request.Method
		path := c.Request.URL.Path

		// 记录日志
		log.Printf("[GIN] %v | %3d | %13v | %15s | %s %s",
			end.Format("2006/01/02 - 15:04:05"),
			status,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

// RecoveryWithZap 使用Zap记录崩溃恢复中间件
func RecoveryWithZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				stack := string(debug.Stack())
				log.Printf("[PANIC RECOVER] %v\n%s", err, stack)

				// 返回500错误
				response.InternalServerError(c, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}

// RoleRequired 检查角色权限
func RoleRequired(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取角色
		rolesValue, exists := c.Get("roles")
		if !exists {
			response.Unauthorized(c, "未获取到用户角色信息")
			c.Abort()
			return
		}

		roles, ok := rolesValue.([]string)
		if !ok {
			response.InternalServerError(c, "获取角色信息格式错误")
			c.Abort()
			return
		}

		// 检查角色权限
		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			for _, role := range roles {
				if role == requiredRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			response.Forbidden(c, "权限不足，需要以下角色之一: "+strings.Join(requiredRoles, ", "))
			c.Abort()
			return
		}

		c.Next()
	}
}
