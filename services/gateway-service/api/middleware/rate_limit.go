package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"

	"mall-go/pkg/response"
	"mall-go/services/gateway-service/infrastructure/config"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
}

// MemoryRateLimiter 基于内存的限流器
type MemoryRateLimiter struct {
	limiters map[string]*rate.Limiter
	config   config.RateLimitConfig
}

// NewMemoryRateLimiter 创建内存限流器
func NewMemoryRateLimiter(config config.RateLimitConfig) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// Allow 判断请求是否允许通过
func (l *MemoryRateLimiter) Allow(key string) bool {
	// 检查是否存在指定路径的限流配置
	limit := l.config.Limit
	burst := l.config.Burst

	// 检查端点特定配置
	for _, endpoint := range l.config.Endpoints {
		if key == endpoint.Path {
			limit = endpoint.Limit
			burst = endpoint.Burst
			break
		}
	}

	// 获取或创建限流器
	limiter, exists := l.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(limit)/60.0, burst) // 每分钟的请求数
		l.limiters[key] = limiter
	}

	return limiter.Allow()
}

// RedisRateLimiter 基于Redis的限流器
type RedisRateLimiter struct {
	client *redis.Client
	config config.RateLimitConfig
}

// NewRedisRateLimiter 创建Redis限流器
func NewRedisRateLimiter(client *redis.Client, config config.RateLimitConfig) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		config: config,
	}
}

// Allow 判断请求是否允许通过
func (l *RedisRateLimiter) Allow(key string) bool {
	// 检查是否存在指定路径的限流配置
	limit := l.config.Limit
	burst := l.config.Burst

	// 检查端点特定配置
	for _, endpoint := range l.config.Endpoints {
		if key == endpoint.Path {
			limit = endpoint.Limit
			burst = endpoint.Burst
			break
		}
	}

	// 构建Redis键
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	
	// 使用Redis的计数器来实现限流
	// 这里使用简单的计数器实现，生产环境可以使用令牌桶或滑动窗口算法
	ctx := l.client.Context()
	
	// 获取当前计数
	val, err := l.client.Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		// 如果出错，允许请求通过（降级策略）
		return true
	}

	var count int
	if val != "" {
		count, _ = strconv.Atoi(val)
	}

	// 如果超过限制，拒绝请求
	if count >= limit {
		return false
	}

	// 增加计数
	l.client.Incr(ctx, redisKey)
	
	// 如果是新的计数周期，设置过期时间（1分钟）
	if count == 0 {
		l.client.Expire(ctx, redisKey, time.Minute)
	}

	return true
}

// 全局限流器
var globalRateLimiter RateLimiter

// InitRateLimiter 初始化限流器
func InitRateLimiter(client *redis.Client) {
	rateLimitConfig := config.GlobalConfig.RateLimit

	if !rateLimitConfig.Enabled {
		return
	}

	if rateLimitConfig.Type == "redis" && client != nil {
		globalRateLimiter = NewRedisRateLimiter(client, rateLimitConfig)
	} else {
		globalRateLimiter = NewMemoryRateLimiter(rateLimitConfig)
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果限流功能未启用或未初始化限流器，直接通过
		if !config.GlobalConfig.RateLimit.Enabled || globalRateLimiter == nil {
			c.Next()
			return
		}

		// 获取限流键（使用请求路径作为键）
		key := c.Request.URL.Path

		// 判断是否允许请求通过
		if !globalRateLimiter.Allow(key) {
			response.Fail(c, http.StatusTooManyRequests, "请求频率过高，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}