package config

import (
	"fmt"
	"os"
	"path/filepath"

	mallconfig "mall-go/pkg/config"
)

// Config represents the application configuration for the user service
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Registry RegistryConfig `mapstructure:"registry"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Tracing  TracingConfig  `mapstructure:"tracing"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	HTTPPort int    `mapstructure:"http_port"`
	GRPCPort int    `mapstructure:"grpc_port"`
	Mode     string `mapstructure:"mode"`
	Timeout  int    `mapstructure:"timeout"`
}

// GRPCConfig represents the gRPC server configuration
type GRPCConfig struct {
	MaxConcurrentStreams int `mapstructure:"max_concurrent_streams"`
	ConnectionTimeout    int `mapstructure:"connection_timeout"`
	KeepaliveTime        int `mapstructure:"keepalive_time"`
	KeepaliveTimeout     int `mapstructure:"keepalive_timeout"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds Redis-related configurations
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// AuthConfig represents the authentication configuration
type AuthConfig struct {
	JWTSecret         string `mapstructure:"jwt_secret"`
	TokenExpiry       int    `mapstructure:"token_expiry"`
	RefreshTokenExpiry int   `mapstructure:"refresh_token_expiry"`
}

// RegistryConfig represents the service registry configuration
type RegistryConfig struct {
	Type        string   `mapstructure:"type"`
	Address     string   `mapstructure:"address"`
	ServiceName string   `mapstructure:"service_name"`
	Tags        []string `mapstructure:"tags"`
}

// LoggingConfig represents the logging configuration
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// TracingConfig represents the distributed tracing configuration
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Type        string `mapstructure:"type"`
	ServiceName string `mapstructure:"service_name"`
	Endpoint    string `mapstructure:"endpoint"`
}

// LoadConfig loads the configuration for the service
func LoadConfig() (*Config, error) {
	// 确定服务名称，从当前目录结构推断
	serviceName := "user-service" // 默认为user-service

	// 尝试从环境变量获取服务名称
	if envName := os.Getenv("SERVICE_NAME"); envName != "" {
		serviceName = envName
	} else {
		// 尝试从当前目录路径推断服务名称
		currentDir, err := os.Getwd()
		if err == nil {
			dirParts := filepath.SplitList(currentDir)
			for _, part := range dirParts {
				if part == "user-service" || part == "auth-service" || part == "cart-service" ||
					part == "gateway-service" || part == "inventory-service" ||
					part == "order-service" || part == "payment-service" || part == "product-service" {
					serviceName = part
					break
				}
			}
		}
	}

	// 加载配置
	var config Config
	err := mallconfig.LoadServiceConfig(serviceName, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load config for service %s: %w", serviceName, err)
	}

	return &config, nil
}
