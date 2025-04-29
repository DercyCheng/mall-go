package auth

import (
	"errors"
	"fmt"
	"time"

	appErrors "mall-go/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Service JWT服务
type Service struct {
	config *Config
}

// Config JWT配置
type Config struct {
	// Secret 密钥
	Secret string
	// AccessTokenExpiry 访问令牌过期时间
	AccessTokenExpiry time.Duration
	// RefreshTokenExpiry 刷新令牌过期时间
	RefreshTokenExpiry time.Duration
	// AccessTokenCookieName 访问令牌Cookie名称
	AccessTokenCookieName string
	// RefreshTokenCookieName 刷新令牌Cookie名称
	RefreshTokenCookieName string
	// EnableRefreshToken 是否启用刷新令牌
	EnableRefreshToken bool
	// IssuerName 令牌颁发者
	IssuerName string
}

// Claims JWT声明
type Claims struct {
	jwt.RegisteredClaims
	// UserID 用户ID
	UserID string `json:"uid,omitempty"`
	// Username 用户名
	Username string `json:"username,omitempty"`
	// Email 邮箱
	Email string `json:"email,omitempty"`
	// Roles 角色
	Roles []string `json:"roles,omitempty"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		AccessTokenExpiry:      24 * time.Hour,
		RefreshTokenExpiry:     7 * 24 * time.Hour,
		AccessTokenCookieName:  "access_token",
		RefreshTokenCookieName: "refresh_token",
		EnableRefreshToken:     true,
		IssuerName:             "mall-go",
	}
}

// NewService 创建JWT服务
func NewService(config *Config) *Service {
	if config == nil {
		config = DefaultConfig()
	}

	if config.Secret == "" {
		panic("JWT secret cannot be empty")
	}

	return &Service{
		config: config,
	}
}

// GenerateAccessToken 生成访问令牌
func (s *Service) GenerateAccessToken(userID, username, email string, roles []string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.AccessTokenExpiry)

	// 创建唯一的令牌ID
	tokenID := uuid.New().String()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Issuer:    s.config.IssuerName,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken 生成刷新令牌
func (s *Service) GenerateRefreshToken(userID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.RefreshTokenExpiry)

	// 创建唯一的令牌ID
	tokenID := uuid.New().String()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Issuer:    s.config.IssuerName,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateAccessToken 验证访问令牌
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, appErrors.UnauthorizedError(appErrors.CodeTokenExpired, "令牌已过期")
		}
		return nil, appErrors.UnauthorizedError(appErrors.CodeInvalidToken, "无效的令牌")
	}

	// 获取声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, appErrors.UnauthorizedError(appErrors.CodeInvalidToken, "无效的令牌声明")
}

// ValidateRefreshToken 验证刷新令牌
func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, appErrors.UnauthorizedError(appErrors.CodeTokenExpired, "刷新令牌已过期")
		}
		return nil, appErrors.UnauthorizedError(appErrors.CodeInvalidToken, "无效的刷新令牌")
	}

	// 获取声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, appErrors.UnauthorizedError(appErrors.CodeInvalidToken, "无效的刷新令牌声明")
}

// ExtractClaims 从令牌中提取声明
func (s *Service) ExtractClaims(tokenString string) (*Claims, error) {
	// 解析令牌但不验证签名，用于快速提取声明
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// GetTokenID 获取令牌ID
func (s *Service) GetTokenID(tokenString string) (string, error) {
	claims, err := s.ExtractClaims(tokenString)
	if err != nil {
		return "", err
	}

	return claims.ID, nil
}

// --------------------------------------------------------------
// 以下为兼容旧版本的函数

var jwtSecret string

// InitJWTSecret 初始化JWT密钥
func InitJWTSecret(secret string) {
	jwtSecret = secret
}

// OldClaims 旧版JWT声明
type OldClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成令牌
func GenerateToken(userID string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour)

	claims := OldClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "mall-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*OldClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OldClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*OldClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
