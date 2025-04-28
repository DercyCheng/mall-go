package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"mall-go/pkg/response"
	"mall-go/services/gateway-service/infrastructure/config"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, http.StatusUnauthorized, "未提供token")
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Fail(c, http.StatusUnauthorized, "token格式错误")
			c.Abort()
			return
		}

		// 解析token
		token := parts[1]
		claims, err := parseToken(token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "无效的token: "+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		
		// 将Token信息传递给后端服务(透传)
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Username", claims.Username)

		c.Next()
	}
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// parseToken 解析JWT Token
func parseToken(tokenString string) (*JWTClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名算法")
		}
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token有效性
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// SkipJWTAuth 判断是否跳过JWT认证
func SkipJWTAuth(path string) bool {
	// 不需要认证的路径
	noAuthPaths := []string{
		"/health",
		"/api/users/login",
		"/api/users/register",
	}

	// 检查是否为无需认证的路径
	for _, p := range noAuthPaths {
		if p == path {
			return true
		}
	}

	return false
}