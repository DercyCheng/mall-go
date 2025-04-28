package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config 项目配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Consul   ConsulConfig   `mapstructure:"consul"`
	Nacos    NacosConfig    `mapstructure:"nacos"`
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
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Address     string            `mapstructure:"address"`
	ServiceName string            `mapstructure:"serviceName"`
	Tags        []string          `mapstructure:"tags"`
	Meta        map[string]string `mapstructure:"meta"`
}

// NacosConfig Nacos配置
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	Group     string `mapstructure:"group"`
	DataID    string `mapstructure:"dataID"`
}

var GlobalConfig Config

// InitConfig 初始化配置
func InitConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return err
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Printf("解析配置文件失败: %v", err)
		return err
	}

	log.Printf("成功加载配置文件: %s", viper.ConfigFileUsed())
	return nil
}

// GetDatabaseDSN 获取数据库连接字符串
func GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		GlobalConfig.Database.Username,
		GlobalConfig.Database.Password,
		GlobalConfig.Database.Host,
		GlobalConfig.Database.Port,
		GlobalConfig.Database.Database)
}

// GetServerAddress 获取服务器地址
func GetServerAddress() string {
	return fmt.Sprintf(":%d", GlobalConfig.Server.Port)
}

// GetConsulAddress 获取Consul地址
func GetConsulAddress() string {
	return GlobalConfig.Consul.Address
}

// GetRedisAddress 获取Redis地址
func GetRedisAddress() string {
	return fmt.Sprintf("%s:%d", GlobalConfig.Redis.Host, GlobalConfig.Redis.Port)
}