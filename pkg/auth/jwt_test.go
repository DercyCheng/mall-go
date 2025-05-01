package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	// 创建JWT服务配置
	config := &Config{
		Secret:                "test-secret-key-for-jwt-testing",
		AccessTokenExpiry:     1 * time.Hour,
		RefreshTokenExpiry:    24 * time.Hour,
		AccessTokenCookieName: "access_token",
		RefreshTokenCookieName: "refresh_token",
		EnableRefreshToken:    true,
		IssuerName:            "mall-go-test",
	}

	// 初始化JWT服务
	jwtService := NewService(config)
	
	// 测试用户数据
	userID := "user-123"
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user", "admin"}
	
	t.Run("GenerateAccessToken", func(t *testing.T) {
		// 生成访问令牌
		token, expiryTime, err := jwtService.GenerateAccessToken(userID, username, email, roles)
		
		// 断言
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, expiryTime.After(time.Now()))
		assert.True(t, expiryTime.Before(time.Now().Add(config.AccessTokenExpiry+time.Second)))
	})
	
	t.Run("ValidateAccessToken", func(t *testing.T) {
		// 生成令牌
		token, _, err := jwtService.GenerateAccessToken(userID, username, email, roles)
		assert.NoError(t, err)
		
		// 验证令牌
		claims, err := jwtService.ValidateAccessToken(token)
		
		// 断言
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, username, claims.Username)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, roles, claims.Roles)
	})
	
	t.Run("GenerateRefreshToken", func(t *testing.T) {
		// 生成刷新令牌
		token, expiryTime, err := jwtService.GenerateRefreshToken(userID)
		
		// 断言
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, expiryTime.After(time.Now()))
		assert.True(t, expiryTime.Before(time.Now().Add(config.RefreshTokenExpiry+time.Second)))
	})
	
	t.Run("ValidateRefreshToken", func(t *testing.T) {
		// 生成刷新令牌
		token, _, err := jwtService.GenerateRefreshToken(userID)
		assert.NoError(t, err)
		
		// 验证刷新令牌
		claims, err := jwtService.ValidateRefreshToken(token)
		
		// 断言
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})
	
	t.Run("InvalidToken", func(t *testing.T) {
		// 验证无效令牌
		_, err := jwtService.ValidateAccessToken("invalid.token.string")
		assert.Error(t, err)
	})
	
	t.Run("ExtractClaims", func(t *testing.T) {
		// 生成令牌
		token, _, err := jwtService.GenerateAccessToken(userID, username, email, roles)
		assert.NoError(t, err)
		
		// 提取声明
		claims, err := jwtService.ExtractClaims(token)
		
		// 断言
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})
}

func TestJWTCompatibility(t *testing.T) {
	// 测试旧版JWT函数
	jwtSecret := "old-jwt-secret"
	InitJWTSecret(jwtSecret)
	
	userID := "old-user-123"
	
	t.Run("GenerateToken", func(t *testing.T) {
		token, err := GenerateToken(userID)
		
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
	
	t.Run("ParseToken", func(t *testing.T) {
		token, err := GenerateToken(userID)
		assert.NoError(t, err)
		
		claims, err := ParseToken(token)
		
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})
}