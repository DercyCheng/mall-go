package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// ServiceConfig 服务配置接口
type ServiceConfig interface {
	// 根据需要添加配置相关方法
}

// LoadServiceConfig 为特定服务加载配置
func LoadServiceConfig(serviceName string, configStruct interface{}) error {
	// 获取环境变量或使用默认值
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	
	fmt.Printf("加载服务 [%s] 的配置，环境: %s\n", serviceName, env)

	// 确定配置目录路径
	// 尝试多种可能的路径
	var configDir string
	possiblePaths := []string{
		filepath.Join(".", "configs"),                  // 当前目录下的configs
		filepath.Join("..", "..", "..", "configs"),     // 服务目录三级向上
		filepath.Join("..", "..", "configs"),           // 服务目录两级向上
		filepath.Join("..", "configs"),                 // 服务目录一级向上
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			configDir = path
			fmt.Printf("找到配置目录: %s\n", path)
			break
		}
	}

	if configDir == "" {
		// 如果找不到配置目录，返回错误
		return fmt.Errorf("找不到配置目录")
	}

	// 创建配置加载器
	loader := NewConfigLoader(configDir, serviceName, env)
	
	// 加载配置
	if err := loader.Load(); err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 将配置解析到结构体
	if err := loader.UnmarshalAll(configStruct); err != nil {
		return fmt.Errorf("解析配置失败: %w", err)
	}

	return nil
}
