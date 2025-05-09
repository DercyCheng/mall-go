package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientFactory 用于创建和管理 gRPC 客户端连接
type ClientFactory struct {
	mu          sync.RWMutex
	connections map[string]*grpc.ClientConn
	options     map[string][]grpc.DialOption
}

// NewClientFactory 创建一个新的客户端工厂实例
func NewClientFactory() *ClientFactory {
	return &ClientFactory{
		connections: make(map[string]*grpc.ClientConn),
		options:     make(map[string][]grpc.DialOption),
	}
}

// SetOptions 为指定服务设置连接选项
func (f *ClientFactory) SetOptions(serviceName string, options ...grpc.DialOption) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.options[serviceName] = options
}

// GetConnection 获取指定服务的连接，如果不存在则创建新连接
func (f *ClientFactory) GetConnection(ctx context.Context, serviceName, target string) (*grpc.ClientConn, error) {
	// 首先检查是否已有现有连接
	f.mu.RLock()
	conn, exists := f.connections[serviceName]
	f.mu.RUnlock()

	if exists {
		// 如果连接已存在且有效，则直接返回
		if conn.GetState() != connectivity.Shutdown {
			return conn, nil
		}
		// 连接已关闭，移除并重建
		f.mu.Lock()
		delete(f.connections, serviceName)
		f.mu.Unlock()
	}

	// 创建新连接
	return f.createConnection(ctx, serviceName, target)
}

// createConnection 创建一个新的 gRPC 连接
func (f *ClientFactory) createConnection(ctx context.Context, serviceName, target string) (*grpc.ClientConn, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 如果在加锁后发现连接已被其他协程创建，则直接返回
	if conn, exists := f.connections[serviceName]; exists {
		return conn, nil
	}

	// 设置默认选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 使用不安全连接（生产环境应使用TLS）
		grpc.WithBlock(), // 阻塞直到连接建立
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 负载均衡策略
	}

	// 应用自定义选项
	if customOpts, ok := f.options[serviceName]; ok {
		opts = append(opts, customOpts...)
	}

	// 创建连接
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", target, err)
	}

	// 存储连接
	f.connections[serviceName] = conn

	return conn, nil
}

// CloseConnection 关闭指定服务的连接
func (f *ClientFactory) CloseConnection(serviceName string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	conn, exists := f.connections[serviceName]
	if !exists {
		return nil
	}

	delete(f.connections, serviceName)
	return conn.Close()
}

// CloseAll 关闭所有连接
func (f *ClientFactory) CloseAll() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for name, conn := range f.connections {
		_ = conn.Close() // 忽略关闭错误
		delete(f.connections, name)
	}
}
