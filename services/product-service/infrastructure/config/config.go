package config

import (
	"fmt"
	mallconfig "mall-go/pkg/config"
)

// Config represents the application configuration for the product service
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
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
	Debug           bool   `mapstructure:"debug"`
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

// LoadConfig loads the configuration for the product service
func LoadConfig() (*Config, error) {
	config := &Config{}
	err := mallconfig.LoadServiceConfig("product-service", config)
	if err != nil {
		return nil, fmt.Errorf("failed to load product service config: %w", err)
	}
	return config, nil
}
