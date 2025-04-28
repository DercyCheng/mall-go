package middleware

import (
	"github.com/gin-gonic/gin"

	"mall-go/services/gateway-service/infrastructure/config"
)

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		corsConfig := config.GlobalConfig.CORS
		
		// 如果CORS不启用，直接跳过
		if !corsConfig.Enabled {
			c.Next()
			return
		}

		// 设置允许的源
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// 检查是否允许该源访问
			allowOrigin := "*"
			for _, o := range corsConfig.AllowOrigins {
				if o == origin || o == "*" {
					allowOrigin = origin
					break
				}
			}
			c.Header("Access-Control-Allow-Origin", allowOrigin)
		}

		// 设置允许的方法
		c.Header("Access-Control-Allow-Methods", joinStrings(corsConfig.AllowMethods, ", "))
		
		// 设置允许的头部
		c.Header("Access-Control-Allow-Headers", joinStrings(corsConfig.AllowHeaders, ", "))
		
		// 设置是否允许发送Cookie
		if corsConfig.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		
		// 设置预检请求缓存时间
		if corsConfig.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", "600") // 缓存10分钟
		}
		
		// 处理OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// joinStrings 使用分隔符拼接字符串数组
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	
	return result
}