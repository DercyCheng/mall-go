package config

import (
	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Consul   ConsulConfig   `mapstructure:"consul"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"readTimeout"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
	Mode         string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	LogLevel string `mapstructure:"logLevel"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"` // 过期时间，单位分钟
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Address     string            `mapstructure:"address"`
	ServiceName string            `mapstructure:"serviceName"`
	Tags        []string          `mapstructure:"tags"`
	Meta        map[string]string `mapstructure:"meta"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}