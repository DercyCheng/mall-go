package config

import (
	"fmt"
	mallConfig "mall-go/pkg/config"
)

// Config represents the application configuration
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Registry  RegistryConfig  `mapstructure:"registry"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Services  ServicesConfig  `mapstructure:"services"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Tracing   TracingConfig   `mapstructure:"tracing"`
}

// ServerConfig holds server-related configurations
type ServerConfig struct {
	HTTPPort int    `mapstructure:"http_port"`
	GRPCPort int    `mapstructure:"grpc_port"`
	Mode     string `mapstructure:"mode"`
	Timeout  int    `mapstructure:"timeout"`
}

// DatabaseConfig holds database-related configurations
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
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

// RegistryConfig holds service registry configurations
type RegistryConfig struct {
	Type        string   `mapstructure:"type"`
	Address     string   `mapstructure:"address"`
	ServiceName string   `mapstructure:"service_name"`
	Tags        []string `mapstructure:"tags"`
}

// AuthConfig holds authentication service configurations
type AuthConfig struct {
	JWTSecret         string `mapstructure:"jwt_secret"`
	TokenExpiry       int    `mapstructure:"token_expiry"`
	RefreshTokenExpiry int   `mapstructure:"refresh_token_expiry"`
}

// ServiceConfig holds configuration for a dependent service
type ServiceConfig struct {
	Name    string `mapstructure:"name"`
	Timeout int    `mapstructure:"timeout"`
}

// ServicesConfig holds configurations for dependent services
type ServicesConfig struct {
	Product    ServiceConfig `mapstructure:"product"`
	Inventory  ServiceConfig `mapstructure:"inventory"`
	Promotion  ServiceConfig `mapstructure:"promotion"`
}

// LoggingConfig holds logging configurations
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// TracingConfig holds distributed tracing configurations
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Type        string `mapstructure:"type"`
	ServiceName string `mapstructure:"service_name"`
	Endpoint    string `mapstructure:"endpoint"`
}

var config *Config

// LoadConfig loads the configuration for cart service
func LoadConfig() error {
	config = new(Config)
	err := mallConfig.LoadServiceConfig("cart-service", config)
	if err != nil {
		return fmt.Errorf("failed to load cart service config: %w", err)
	}
	return nil
}

// GetConfig returns the loaded configuration
func GetConfig() *Config {
	return config
}
