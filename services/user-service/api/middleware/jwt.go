package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/auth"
	"mall-go/pkg/response"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, 401, "无效的Token格式")
			c.Abort()
			return
		}

		// 解析Token
		token := parts[1]
		claims, err := auth.ParseToken(token)
		if err != nil {
			response.Error(c, 401, "Token验证失败: "+err.Error())
			c.Abort()
			return
		}

		// 将信息存入上下文
		c.Set("userId", claims.ID)
		c.Set("username", claims.Username)

		c.Next()
	}
}