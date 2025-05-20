package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查格式 "Bearer token"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// 获取token字符串
		tokenString := parts[1]

		// 在实际项目中，这里应该使用JWT库验证token
		// 使用token变量，避免编译器警告
		_ = tokenString

		// 模拟解析token获取用户ID
		// 实际项目中应该使用JWT库解析token
		userID := "user123" // 假设从token中获取的用户ID

		// 将用户ID保存到上下文中，供后续处理器使用
		c.Set("userID", userID)

		c.Next()
	}
}
