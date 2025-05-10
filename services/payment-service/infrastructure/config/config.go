package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/spf13/viper"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Registry RegistryConfig `mapstructure:"registry"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Services ServicesConfig `mapstructure:"services"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Tracing  TracingConfig  `mapstructure:"tracing"`
	Payment  PaymentConfig  `mapstructure:"payment"`
}

// ServerConfig holds server specific configuration
type ServerConfig struct {
	HTTPPort int    `mapstructure:"http_port"`
	GRPCPort int    `mapstructure:"grpc_port"`
	Mode     string `mapstructure:"mode"`
	Timeout  int    `mapstructure:"timeout"`
}

// DatabaseConfig holds database specific configuration
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

// RedisConfig holds Redis specific configuration
type RedisConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Password      string `mapstructure:"password"`
	DB            int    `mapstructure:"db"`
	PoolSize      int    `mapstructure:"pool_size"`
	MinIdleConns  int    `mapstructure:"min_idle_conns"`
}

// RegistryConfig holds service registry configuration
type RegistryConfig struct {
	Type        string   `mapstructure:"type"`
	Address     string   `mapstructure:"address"`
	ServiceName string   `mapstructure:"service_name"`
	Tags        []string `mapstructure:"tags"`
}

// RabbitMQConfig holds RabbitMQ specific configuration
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

// ServicesConfig holds dependencies on other services
type ServicesConfig struct {
	Order OrderServiceConfig `mapstructure:"order"`
}

// OrderServiceConfig holds order service specific configuration
type OrderServiceConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
}

// LoggingConfig holds logging specific configuration
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// TracingConfig holds tracing specific configuration
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Type        string `mapstructure:"type"`
	ServiceName string `mapstructure:"service_name"`
	Endpoint    string `mapstructure:"endpoint"`
}

// PaymentConfig holds payment specific configuration
type PaymentConfig struct {
	Providers []PaymentProviderConfig `mapstructure:"providers"`
}

// PaymentProviderConfig holds payment provider specific configuration
type PaymentProviderConfig struct {
	Name              string `mapstructure:"name"`
	AppID             string `mapstructure:"app_id"`
	MerchantPrivateKey string `mapstructure:"merchant_private_key"`
	AlipayPublicKey   string `mapstructure:"alipay_public_key"`
	MerchantID        string `mapstructure:"mch_id"`
	APIKey            string `mapstructure:"api_key"`
	NotifyURL         string `mapstructure:"notify_url"`
	ReturnURL         string `mapstructure:"return_url"`
	Sandbox           bool   `mapstructure:"sandbox"`
}

// LoadConfig loads the configuration from config.yaml
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// Try to find the config file in various locations
	configPaths := []string{
		".",
		"./configs",
		"../configs",
		"../../configs",
	}
	
	// Get executable directory and add it to the search paths
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		configPaths = append(configPaths, execDir, filepath.Join(execDir, "configs"))
	}
	
	// Add the paths to viper's search paths
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
	
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &config, nil
}
