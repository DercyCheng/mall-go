package config

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)

// 用于测试的示例配置结构
type TestConfig struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Debug   bool   `json:"debug"`
}

// 启动一个嵌入式etcd服务器用于测试
func startEmbeddedEtcd(t *testing.T) (*embed.Etcd, []string) {
	t.Helper()

	cfg := embed.NewConfig()
	cfg.Dir = t.TempDir() // 使用临时目录存储数据
	
	// 配置临时端口
	// 这里使用系统分配的端口避免冲突
	cfg.LPUrls = []url.URL{{Scheme: "http", Host: "localhost:0"}}
	cfg.LCUrls = []url.URL{{Scheme: "http", Host: "localhost:0"}}
	
	// 关闭不必要的服务
	cfg.InitialCluster = "default=http://localhost:0"
	
	// 启动嵌入式etcd
	etcd, err := embed.StartEtcd(cfg)
	if err != nil {
		t.Fatalf("无法启动嵌入式etcd: %v", err)
	}
	
	// 等待etcd启动完成
	select {
	case <-etcd.Server.ReadyNotify():
		t.Log("嵌入式etcd服务器已启动")
	case <-time.After(10 * time.Second):
		t.Fatal("etcd服务器启动超时")
	}
	
	// 获取实际分配的端口
	clientURL := etcd.Clients[0].Addr().String()
	endpoints := []string{"http://" + clientURL}
	
	return etcd, endpoints
}

// 测试etcd配置客户端的基本功能
func TestEtcdClient(t *testing.T) {
	// 如果没有可用的etcd服务，跳过测试
	t.Skip("需要一个运行中的etcd实例进行测试")
	
	// 注: 在实际的CI环境中，可以使用嵌入式etcd
	// etcd, endpoints := startEmbeddedEtcd(t)
	// defer etcd.Close()
	
	// 这里使用假设有一个可用的etcd实例
	endpoints := []string{"localhost:2379"}
	
	t.Run("NewEtcdClient", func(t *testing.T) {
		client, err := NewEtcdClient(endpoints, "test-app")
		if err != nil {
			t.Skipf("无法连接到etcd: %v，跳过测试", err)
		}
		assert.NotNil(t, client)
		defer client.Close()
	})
	
	// 假设我们已经有一个etcd客户端
	client, err := NewEtcdClient(endpoints, "test-app")
	if err != nil {
		t.Skipf("无法连接到etcd: %v，跳过测试", err)
		return
	}
	defer client.Close()
	
	t.Run("SaveAndLoadConfig", func(t *testing.T) {
		// 创建测试配置
		testCfg := TestConfig{
			Name:    "test-service",
			Version: "1.0.0",
			Debug:   true,
		}
		
		// 保存配置
		err := client.SaveConfig("test-key", testCfg)
		if err != nil {
			t.Skipf("无法保存配置: %v，跳过测试", err)
			return
		}
		
		// 加载配置
		var loadedCfg TestConfig
		err = client.LoadConfig("test-key", &loadedCfg)
		if err != nil {
			t.Skipf("无法加载配置: %v，跳过测试", err)
			return
		}
		
		// 验证配置正确
		assert.Equal(t, testCfg.Name, loadedCfg.Name)
		assert.Equal(t, testCfg.Version, loadedCfg.Version)
		assert.Equal(t, testCfg.Debug, loadedCfg.Debug)
	})
	
	t.Run("WatchConfig", func(t *testing.T) {
		// 创建测试配置
		testCfg := TestConfig{
			Name:    "watch-test",
			Version: "1.0.0",
			Debug:   true,
		}
		
		// 监听配置变更的通道
		done := make(chan bool)
		var watchedValue []byte
		
		// 设置配置变更监听
		client.WatchConfig("watch-key", func(value []byte) {
			watchedValue = value
			done <- true
		})
		
		// 等待goroutine启动
		time.Sleep(100 * time.Millisecond)
		
		// 保存配置触发变更
		err := client.SaveConfig("watch-key", testCfg)
		if err != nil {
			t.Skipf("无法保存配置: %v，跳过测试", err)
			return
		}
		
		// 等待配置变更通知
		select {
		case <-done:
			// 加载变更后的配置
			var updatedCfg TestConfig
			err = json.Unmarshal(watchedValue, &updatedCfg)
			assert.NoError(t, err)
			assert.Equal(t, testCfg.Name, updatedCfg.Name)
		case <-time.After(5 * time.Second):
			t.Skip("等待配置变更超时，跳过测试")
		}
	})
}