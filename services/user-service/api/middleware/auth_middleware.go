package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/auth"
	"mall-go/pkg/errors"
	"mall-go/pkg/response"
)

// Auth 认证配置
type Auth struct {
	jwtService       *auth.Service
	skipPaths        map[string]bool
	enableRateLimit  bool
	rateLimiter      AuthRateLimiter
	blockListChecker TokenBlockListChecker
}

// RateLimiter 限流器接口
type AuthRateLimiter interface {
	Allow(key string) bool
}

// TokenBlockListChecker 令牌黑名单检查器接口
type TokenBlockListChecker interface {
	IsBlocked(jti string) bool
}

// NewAuth 创建认证中间件
func NewAuth(jwtService *auth.Service) *Auth {
	return &Auth{
		jwtService: jwtService,
		skipPaths:  make(map[string]bool),
	}
}

// Middleware 创建认证中间件
func (a *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过认证
		if a.skipPaths[c.FullPath()] {
			c.Next()
			return
		}

		// 从请求中获取令牌
		tokenString := a.extractToken(c)
		if tokenString == "" {
			response.Unauthorized(c, "认证令牌缺失")
			c.Abort()
			return
		}

		// 启用限流保护
		if a.enableRateLimit && a.rateLimiter != nil {
			clientIP := c.ClientIP()
			if !a.rateLimiter.Allow(clientIP) {
				response.ErrorWithCode(c, 429, "请求过于频繁，请稍后再试")
				c.Abort()
				return
			}
		}

		// 验证令牌
		claims, err := a.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			if appErr, ok := errors.As(err); ok {
				response.Error(c, appErr)
			} else {
				response.Unauthorized(c, "无效的令牌")
			}
			c.Abort()
			return
		}

		// 检查令牌是否在黑名单中
		if a.blockListChecker != nil && a.blockListChecker.IsBlocked(claims.ID) {
			response.Unauthorized(c, "令牌已被撤销")
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

// RequireRoles 创建角色验证中间件
func (a *Auth) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			response.Unauthorized(c, "用户未认证")
			c.Abort()
			return
		}

		// 检查用户是否具有所需角色
		hasRequiredRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoles.([]string) {
				if requiredRole == userRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// SkipAuth 设置跳过认证的路径
func (a *Auth) SkipAuth(paths ...string) {
	for _, path := range paths {
		a.skipPaths[path] = true
	}
}

// EnableRateLimit 启用限流保护
func (a *Auth) EnableRateLimit(limiter AuthRateLimiter) {
	a.enableRateLimit = true
	a.rateLimiter = limiter
}

// SetBlockListChecker 设置令牌黑名单检查器
func (a *Auth) SetBlockListChecker(checker TokenBlockListChecker) {
	a.blockListChecker = checker
}

// extractToken 从请求中提取令牌
func (a *Auth) extractToken(c *gin.Context) string {
	// 尝试从Authorization头提取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// 尝试从查询参数提取
	token := c.Query("token")
	if token != "" {
		return token
	}

	// 尝试从Cookie提取
	cookie, err := c.Cookie("access_token")
	if err == nil {
		return cookie
	}

	return ""
}

// CorsMiddleware 创建CORS中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SimpleRateLimiter 简单的令牌桶限流器
type SimpleRateLimiter struct {
	tokens map[string]int64
	limit  int64
}

// NewSimpleRateLimiter 创建简单限流器
func NewSimpleRateLimiter(limit int64) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		tokens: make(map[string]int64),
		limit:  limit,
	}
}

// Allow 检查是否允许请求
func (r *SimpleRateLimiter) Allow(key string) bool {
	// 简单实现，实际应用中可以使用Redis等分布式方案
	return true
}

// LegacyJWTAuthMiddleware 旧版认证中间件（保留向后兼容性）
func LegacyJWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	auth.InitJWTSecret(jwtSecret)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
