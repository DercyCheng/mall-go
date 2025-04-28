package config

import (
	"github.com/spf13/viper"
)

// Config 网关配置
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Consul    ConsulConfig    `mapstructure:"consul"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Services  ServicesConfig  `mapstructure:"services"`
	RateLimit RateLimitConfig `mapstructure:"rateLimit"`
	CORS      CORSConfig      `mapstructure:"cors"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"readTimeout"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
	Mode         string `mapstructure:"mode"`
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Address     string            `mapstructure:"address"`
	ServiceName string            `mapstructure:"serviceName"`
	Tags        []string          `mapstructure:"tags"`
	Meta        map[string]string `mapstructure:"meta"`
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

// ServiceConfig 单个服务配置
type ServiceConfig struct {
	Path        string `mapstructure:"path"`
	StripPrefix bool   `mapstructure:"stripPrefix"`
	Retries     int    `mapstructure:"retries"`
	Timeout     int    `mapstructure:"timeout"` // 超时时间，单位秒
}

// ServicesConfig 服务配置集合
type ServicesConfig struct {
	UserService    ServiceConfig `mapstructure:"user-service"`
	ProductService ServiceConfig `mapstructure:"product-service"`
}

// RateLimitEndpointConfig 限流端点配置
type RateLimitEndpointConfig struct {
	Path  string `mapstructure:"path"`
	Limit int    `mapstructure:"limit"`
	Burst int    `mapstructure:"burst"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled   bool                     `mapstructure:"enabled"`
	Type      string                   `mapstructure:"type"` // memory 或 redis
	Limit     int                      `mapstructure:"limit"`
	Burst     int                      `mapstructure:"burst"`
	Endpoints []RateLimitEndpointConfig `mapstructure:"endpoints"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allowOrigins"`
	AllowMethods     []string `mapstructure:"allowMethods"`
	AllowHeaders     []string `mapstructure:"allowHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"` // 预检请求有效期，单位秒
}

// 全局配置
var GlobalConfig Config

// LoadConfig 加载配置
func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	return nil
}