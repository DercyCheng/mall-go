package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MongoDBConfig struct {
	URI      string        `mapstructure:"uri"`
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type StorageConfig struct {
	MaxUploadSize int64    `mapstructure:"max_upload_size"`
	AllowedTypes  []string `mapstructure:"allowed_types"`
}

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	Storage StorageConfig `mapstructure:"storage"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AutomaticEnv()

	// 设置默认值
	v.SetDefault("server.port", "8083")
	v.SetDefault("server.mode", "debug")
	v.SetDefault("mongodb.timeout", 10)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	zap.L().Info("Configuration loaded successfully")
	return &cfg, nil
}
