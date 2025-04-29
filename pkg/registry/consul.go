// pkg/registry/consul.go
package registry

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// ServiceRegistry 服务注册接口
type ServiceRegistry interface {
	Register(serviceName, serviceID, serviceAddr string, servicePort int, tags []string) error
	Deregister(serviceID string) error
	GetService(serviceName string) ([]*api.AgentService, error)
}

// ConsulRegistry Consul服务注册实现
type ConsulRegistry struct {
	client *api.Client
}

// NewConsulRegistry 创建Consul注册中心客户端
func NewConsulRegistry(addr string) (*ConsulRegistry, error) {
	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulRegistry{
		client: client,
	}, nil
}

// Register 注册服务
func (r *ConsulRegistry) Register(serviceName, serviceID, serviceAddr string, servicePort int, tags []string) error {
	// 创建健康检查
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", serviceAddr, servicePort),
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "30s",
	}

	// 创建服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Tags:    tags,
		Port:    servicePort,
		Address: serviceAddr,
		Check:   check,
	}

	// 注册服务
	err := r.client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Printf("服务注册失败: %v", err)
		return err
	}

	log.Printf("服务[%s]注册成功，地址: %s:%d", serviceName, serviceAddr, servicePort)
	return nil
}

// Deregister 注销服务
func (r *ConsulRegistry) Deregister(serviceID string) error {
	err := r.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		log.Printf("服务注销失败: %v", err)
		return err
	}

	log.Printf("服务[%s]注销成功", serviceID)
	return nil
}

// GetService 获取服务实例
func (r *ConsulRegistry) GetService(serviceName string) ([]*api.AgentService, error) {
	services, err := r.client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", serviceName))
	if err != nil {
		log.Printf("获取服务列表失败: %v", err)
		return nil, err
	}

	var result []*api.AgentService
	for _, service := range services {
		result = append(result, service)
	}

	return result, nil
}
