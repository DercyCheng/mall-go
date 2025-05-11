package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// ConfigLoader 配置加载器
type ConfigLoader struct {
	viper            *viper.Viper
	configDir        string
	environment      string
	serviceName      string
	loadedConfigPath []string
}

// NewConfigLoader 创建配置加载器
func NewConfigLoader(configDir string, serviceName string, environment string) *ConfigLoader {
	if environment == "" {
		environment = "dev" // 默认使用开发环境配置
	}

	return &ConfigLoader{
		viper:       viper.New(),
		configDir:   configDir,
		environment: environment,
		serviceName: serviceName,
	}
}

// Load 加载配置
func (c *ConfigLoader) Load() error {
	// 设置Viper配置
	c.viper.SetConfigType("yaml")
	c.viper.AutomaticEnv()
	c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 加载统一的环境配置文件
	configFile := fmt.Sprintf("%s.yaml", c.environment)
	if err := c.loadFile(configFile); err != nil {
		// 如果没有找到统一的环境配置，尝试使用老的配置加载方式
		// 1. 首先加载公共配置
		if err := c.loadFile("common.yaml"); err != nil {
			return fmt.Errorf("failed to load common config: %w", err)
		}
		
		// 2. 加载环境特定配置
		envConfigFile := fmt.Sprintf("config.%s.yaml", c.environment)
		if err := c.loadFile(envConfigFile); err != nil {
			return fmt.Errorf("failed to load environment config: %w", err)
		}
		
		// 3. 加载服务特定配置
		serviceConfigFile := fmt.Sprintf("%s.yaml", c.serviceName)
		if err := c.loadFile(serviceConfigFile); err != nil {
			return fmt.Errorf("failed to load service config: %w", err)
		}
	} else {
		// 统一配置文件已加载成功
		// 如果使用统一配置，需要正确获取服务特定配置
		c.viper.SetDefault("server", c.viper.Get(fmt.Sprintf("services.%s.server", c.serviceName)))
		c.viper.SetDefault("grpc", c.viper.Get(fmt.Sprintf("services.%s.grpc", c.serviceName)))
		c.viper.SetDefault("registry.service_name", c.viper.GetString(fmt.Sprintf("services.%s.registry.service_name", c.serviceName)))
		c.viper.SetDefault("registry.tags", c.viper.Get(fmt.Sprintf("services.%s.registry.tags", c.serviceName)))
		c.viper.SetDefault("logging.file_path", c.viper.GetString(fmt.Sprintf("services.%s.logging.file_path", c.serviceName)))
		c.viper.SetDefault("tracing.service_name", c.viper.GetString(fmt.Sprintf("services.%s.tracing.service_name", c.serviceName)))
	}

	return nil
}

// loadFile 加载指定配置文件
func (c *ConfigLoader) loadFile(filename string) error {
	configPath := filepath.Join(c.configDir, filename)
	
	// 检查文件是否存在
	_, err := os.Stat(configPath)
	if err != nil {
		return fmt.Errorf("config file %s not found: %w", configPath, err)
	}

	// 读取配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("cannot open config file %s: %w", configPath, err)
	}
	defer file.Close()

	// 加载配置
	c.viper.MergeConfig(file)
	c.loadedConfigPath = append(c.loadedConfigPath, configPath)
	
	return nil
}

// GetViper 获取Viper实例
func (c *ConfigLoader) GetViper() *viper.Viper {
	return c.viper
}

// GetString 获取字符串配置值
func (c *ConfigLoader) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetBool 获取布尔配置值
func (c *ConfigLoader) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetInt 获取整数配置值
func (c *ConfigLoader) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetFloat64 获取浮点数配置值
func (c *ConfigLoader) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

// GetStringSlice 获取字符串切片配置值
func (c *ConfigLoader) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

// GetStringMap 获取map配置值
func (c *ConfigLoader) GetStringMap(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}

// Unmarshal 将配置解析到结构体
func (c *ConfigLoader) Unmarshal(key string, val interface{}) error {
	return c.viper.UnmarshalKey(key, val)
}

// UnmarshalAll 将整个配置解析到结构体
func (c *ConfigLoader) UnmarshalAll(val interface{}) error {
	return c.viper.Unmarshal(val)
}

// GetLoadedConfigPaths 获取已加载的配置文件路径
func (c *ConfigLoader) GetLoadedConfigPaths() []string {
	return c.loadedConfigPath
}
