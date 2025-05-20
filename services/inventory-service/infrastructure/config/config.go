package config

import (
	"os"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Services ServicesConfig `mapstructure:"services"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	Issuer    string `mapstructure:"issuer"`
	Expiry    int    `mapstructure:"expiry"` // 过期时间（分钟）
}

// ServicesConfig 微服务配置
type ServicesConfig struct {
	OrderService   string `mapstructure:"order_service"`
	ProductService string `mapstructure:"product_service"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 设置配置文件路径
	configPath := getConfigPath()

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	// 从环境变量覆盖特定配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	return &config, nil
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	// 默认配置文件路径
	configPath := "configs/config.yaml"

	// 如果通过环境变量指定了配置文件路径，则使用环境变量中的路径
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	}

	return configPath
}
