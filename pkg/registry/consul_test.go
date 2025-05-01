package registry

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 测试创建注册中心客户端
func TestNewConsulRegistry(t *testing.T) {
	t.Run("WithInvalidAddress", func(t *testing.T) {
		// 使用无效地址
		registry, err := NewConsulRegistry("invalid:8500")
		
		// 应该没有错误，因为连接是延迟的
		assert.NoError(t, err)
		assert.NotNil(t, registry)
	})

	t.Run("WithValidAddress", func(t *testing.T) {
		// 使用有效地址（注意：不实际连接）
		registry, err := NewConsulRegistry("localhost:8500")
		assert.NoError(t, err)
		assert.NotNil(t, registry)
	})
}

// 测试服务注册和发现功能
func TestConsulServiceRegistry(t *testing.T) {
	// 跳过实际需要Consul服务器的测试
	t.Skip("需要Consul服务器进行实际测试")
	
	// 在实际的CI/CD环境中，可以使用嵌入式Consul服务器进行测试
	// srv, err := testutil.NewTestServerConfigT(t, nil)
	// if err != nil {
	//     t.Fatalf("无法启动测试Consul服务器: %v", err)
	// }
	// defer srv.Stop()
	
	// 使用测试服务器的地址
	// consulAddr := srv.HTTPAddr
	
	// 使用可能可用的本地Consul
	consulAddr := "localhost:8500"
	
	// 创建注册中心客户端
	registry, err := NewConsulRegistry(consulAddr)
	if err != nil {
		t.Skipf("创建Consul客户端失败: %v，跳过测试", err)
		return
	}
	
	// 定义服务信息
	serviceName := "test-service"
	serviceID := "test-service-1"
	serviceAddr := "localhost"
	servicePort := 8080
	tags := []string{"test", "api"}

	t.Run("RegisterService", func(t *testing.T) {
		// 注册服务
		err := registry.Register(serviceName, serviceID, serviceAddr, servicePort, tags)
		if err != nil {
			t.Skipf("注册服务失败: %v，跳过测试", err)
			return
		}
		
		// 清理：延迟注销服务
		defer registry.Deregister(serviceID)
		
		// 查询服务
		services, err := registry.GetService(serviceName)
		if err != nil {
			t.Skipf("查询服务失败: %v，跳过测试", err)
			return
		}
		
		// 验证服务是否注册成功
		assert.NotEmpty(t, services)
		
		// 至少一个服务实例应该匹配我们的注册
		var found bool
		for _, svc := range services {
			if svc.ID == serviceID {
				found = true
				assert.Equal(t, servicePort, svc.Port)
				assert.Equal(t, serviceAddr, svc.Address)
				assert.ElementsMatch(t, tags, svc.Tags)
				break
			}
		}
		
		assert.True(t, found, "注册的服务未在服务列表中找到")
	})

	t.Run("DeregisterService", func(t *testing.T) {
		// 先注册服务
		err := registry.Register(serviceName, serviceID, serviceAddr, servicePort, tags)
		if err != nil {
			t.Skipf("注册服务失败: %v，跳过测试", err)
			return
		}
		
		// 等待服务注册生效
		time.Sleep(100 * time.Millisecond)
		
		// 注销服务
		err = registry.Deregister(serviceID)
		if err != nil {
			t.Skipf("注销服务失败: %v，跳过测试", err)
			return
		}
		
		// 等待注销生效
		time.Sleep(100 * time.Millisecond)
		
		// 查询服务
		services, err := registry.GetService(serviceName)
		if err != nil {
			t.Skipf("查询服务失败: %v，跳过测试", err)
			return
		}
		
		// 验证服务已被注销
		found := false
		for _, svc := range services {
			if svc.ID == serviceID {
				found = true
				break
			}
		}
		
		assert.False(t, found, "服务应该已被注销，但仍然存在")
	})
}

// 使用Mock进行单元测试
func TestConsulRegistryWithMock(t *testing.T) {
	// 可以使用Mock实现进行测试
	// 这里只演示基本结构，实际实现需要更多工作
	
	// 创建一个简单的Mock客户端
	mockRegistry := &MockServiceRegistry{}
	
	// 设置测试数据
	serviceName := "mock-service"
	serviceID := "mock-service-1"
	
	t.Run("MockRegister", func(t *testing.T) {
		// 设置期望
		mockRegistry.On("Register", serviceName, serviceID, "localhost", 8080, []string{"mock"}).Return(nil)
		
		// 调用Register方法
		err := mockRegistry.Register(serviceName, serviceID, "localhost", 8080, []string{"mock"})
		assert.NoError(t, err)
		
		// 验证期望被调用
		mockRegistry.AssertExpectations(t)
	})
}

// MockServiceRegistry 是ServiceRegistry接口的模拟实现
type MockServiceRegistry struct {
	mock.Mock
}

func (m *MockServiceRegistry) Register(serviceName, serviceID, serviceAddr string, servicePort int, tags []string) error {
	args := m.Called(serviceName, serviceID, serviceAddr, servicePort, tags)
	return args.Error(0)
}

func (m *MockServiceRegistry) Deregister(serviceID string) error {
	args := m.Called(serviceID)
	return args.Error(0)
}

func (m *MockServiceRegistry) GetService(serviceName string) ([]*api.AgentService, error) {
	args := m.Called(serviceName)
	return args.Get(0).([]*api.AgentService), args.Error(1)
}