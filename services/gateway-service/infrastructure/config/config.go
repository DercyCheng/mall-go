package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Auth      AuthConfig      `yaml:"auth"`
	Registry  RegistryConfig  `yaml:"registry"`
	Services  ServicesConfig  `yaml:"services"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	CORS      CORSConfig      `yaml:"cors"`
	Logging   LoggingConfig   `yaml:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// AuthConfig holds authentication-related configuration
type AuthConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	TokenExpiry int    `yaml:"token_expiry"`
}

// GetTokenExpiry returns the token expiry as time.Duration
func (c *AuthConfig) GetTokenExpiry() time.Duration {
	return time.Duration(c.TokenExpiry) * time.Second
}

// RegistryConfig holds service registry configuration
type RegistryConfig struct {
	Type        string   `yaml:"type"`
	Address     string   `yaml:"address"`
	ServiceName string   `yaml:"service_name"`
	Tags        []string `yaml:"tags"`
}

// ServicesConfig holds configuration for backend services
type ServicesConfig struct {
	User ServiceConfig `yaml:"user"`
	// Add other services as they are created
	// Product ServiceConfig `yaml:"product"`
	// Order   ServiceConfig `yaml:"order"`
}

// ServiceConfig holds configuration for an individual service
type ServiceConfig struct {
	Name    string `yaml:"name"`
	Timeout int    `yaml:"timeout"`
}

// GetTimeout returns the service timeout as time.Duration
func (c *ServiceConfig) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled"`
	RequestsPerSecond int  `yaml:"requests_per_second"`
	BurstSize         int  `yaml:"burst_size"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path"`
}

// LoadConfig loads the configuration from the specified file path
func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
