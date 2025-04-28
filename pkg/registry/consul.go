// pkg/registry/consul.go
package registry

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
	// 注册服务实例
	Register(serviceInstance *ServiceInstance) error
	// 注销服务实例
	Deregister(serviceInstance *ServiceInstance) error
	// 根据服务名称获取服务实例列表
	GetServiceInstances(serviceName string) ([]*ServiceInstance, error)
	// 服务发现
	DiscoverServices(serviceName string, tags ...string) ([]*ServiceInstance, error)
}

// ServiceInstance 服务实例
type ServiceInstance struct {
	ID                string            // 服务实例ID
	Name              string            // 服务名称
	Address           string            // 服务地址
	Port              int               // 服务端口
	Tags              []string          // 服务标签
	Meta              map[string]string // 服务元数据
	EnableTagOverride bool              // 是否允许标签覆盖
}

// ConsulServiceRegistry Consul服务注册实现
type ConsulServiceRegistry struct {
	client     *api.Client
	instanceID string
}

// NewConsulServiceRegistry 创建Consul服务注册表
func NewConsulServiceRegistry(consulAddr string) (ServiceRegistry, error) {
	// 创建Consul客户端配置
	config := api.DefaultConfig()
	config.Address = consulAddr

	// 创建Consul客户端
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建Consul客户端失败: %w", err)
	}

	return &ConsulServiceRegistry{
		client: client,
	}, nil
}

// Register 注册服务实例
func (c *ConsulServiceRegistry) Register(serviceInstance *ServiceInstance) error {
	// 创建健康检查
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", serviceInstance.Address, serviceInstance.Port),
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s",
	}

	// 创建服务注册对象
	registration := &api.AgentServiceRegistration{
		ID:                serviceInstance.ID,
		Name:              serviceInstance.Name,
		Tags:              serviceInstance.Tags,
		Port:              serviceInstance.Port,
		Address:           serviceInstance.Address,
		Meta:              serviceInstance.Meta,
		EnableTagOverride: serviceInstance.EnableTagOverride,
		Check:             check,
	}

	// 注册服务
	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("服务注册失败: %w", err)
	}

	// 记录日志
	log.Printf("服务[%s]实例[%s]注册成功, 地址: %s:%d",
		serviceInstance.Name, serviceInstance.ID, serviceInstance.Address, serviceInstance.Port)
	return nil
}

// Deregister 注销服务实例
func (c *ConsulServiceRegistry) Deregister(serviceInstance *ServiceInstance) error {
	if err := c.client.Agent().ServiceDeregister(serviceInstance.ID); err != nil {
		return fmt.Errorf("服务注销失败: %w", err)
	}

	log.Printf("服务[%s]实例[%s]注销成功", serviceInstance.Name, serviceInstance.ID)
	return nil
}

// GetServiceInstances 根据服务名称获取服务实例列表
func (c *ConsulServiceRegistry) GetServiceInstances(serviceName string) ([]*ServiceInstance, error) {
	// 查询服务
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("获取服务实例失败: %w", err)
	}

	// 转换为服务实例对象
	instances := make([]*ServiceInstance, 0, len(services))
	for _, service := range services {
		instance := &ServiceInstance{
			ID:                service.Service.ID,
			Name:              service.Service.Service,
			Address:           service.Service.Address,
			Port:              service.Service.Port,
			Tags:              service.Service.Tags,
			Meta:              service.Service.Meta,
			EnableTagOverride: service.Service.EnableTagOverride,
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// DiscoverServices 服务发现
func (c *ConsulServiceRegistry) DiscoverServices(serviceName string, tags ...string) ([]*ServiceInstance, error) {
	// 使用健康检查API查询健康的服务实例
	var queryOpts api.QueryOptions
	var tag string
	if len(tags) > 0 {
		// 如果提供了标签，使用标签过滤
		tag = tags[0]
	}

	services, _, err := c.client.Health().Service(serviceName, tag, true, &queryOpts)
	if err != nil {
		return nil, fmt.Errorf("服务发现失败: %w", err)
	}

	// 转换为服务实例对象
	instances := make([]*ServiceInstance, 0, len(services))
	for _, service := range services {
		instance := &ServiceInstance{
			ID:                service.Service.ID,
			Name:              service.Service.Service,
			Address:           service.Service.Address,
			Port:              service.Service.Port,
			Tags:              service.Service.Tags,
			Meta:              service.Service.Meta,
			EnableTagOverride: service.Service.EnableTagOverride,
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// ServiceRegistrationLifecycle 服务注册生命周期管理
type ServiceRegistrationLifecycle struct {
	registry        ServiceRegistry
	serviceInstance *ServiceInstance
	running         bool
	stopCh          chan struct{}
}

// NewServiceRegistrationLifecycle 创建服务注册生命周期管理
func NewServiceRegistrationLifecycle(registry ServiceRegistry, serviceInstance *ServiceInstance) *ServiceRegistrationLifecycle {
	return &ServiceRegistrationLifecycle{
		registry:        registry,
		serviceInstance: serviceInstance,
		stopCh:          make(chan struct{}),
	}
}

// Start 开始服务注册
func (s *ServiceRegistrationLifecycle) Start() error {
	if s.running {
		return nil
	}

	if err := s.registry.Register(s.serviceInstance); err != nil {
		return err
	}

	s.running = true

	// 启动一个后台协程定期检查和重新注册服务
	go s.registrationLoop()

	return nil
}

// Stop 停止服务注册
func (s *ServiceRegistrationLifecycle) Stop() error {
	if !s.running {
		return nil
	}

	// 发送停止信号
	close(s.stopCh)

	// 注销服务
	err := s.registry.Deregister(s.serviceInstance)
	s.running = false

	return err
}

// registrationLoop 周期性重新注册，保持服务活跃
func (s *ServiceRegistrationLifecycle) registrationLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 重新注册服务
			if err := s.registry.Register(s.serviceInstance); err != nil {
				log.Printf("服务[%s]实例[%s]重新注册失败: %v",
					s.serviceInstance.Name, s.serviceInstance.ID, err)
			}
		case <-s.stopCh:
			return
		}
	}
}
