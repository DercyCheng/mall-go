package grpc

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	grpcclient "mall-go/pkg/grpc"
	productpb "mall-go/services/product-service/proto"
	userpb "mall-go/services/user-service/proto"
)

// ClientManager 管理订单服务所需的所有 gRPC 客户端
type ClientManager struct {
	logger        *zap.Logger
	factory       *grpcclient.ClientFactory
	userClient    userpb.UserServiceClient
	productClient productpb.ProductServiceClient
	mu            sync.RWMutex
}

// NewClientManager 创建新的 gRPC 客户端管理器
func NewClientManager(logger *zap.Logger) *ClientManager {
	return &ClientManager{
		logger:  logger,
		factory: grpcclient.NewClientFactory(),
	}
}

// GetUserClient 获取用户服务的 gRPC 客户端
func (m *ClientManager) GetUserClient(ctx context.Context) (userpb.UserServiceClient, error) {
	m.mu.RLock()
	if m.userClient != nil {
		client := m.userClient
		m.mu.RUnlock()
		return client, nil
	}
	m.mu.RUnlock()

	// 双重检查锁，避免重复创建客户端
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.userClient != nil {
		return m.userClient, nil
	}

	// 从配置中读取用户服务地址
	target := fmt.Sprintf("%s:%d",
		viper.GetString("services.user.host"),
		viper.GetInt("services.user.port"))

	conn, err := m.factory.GetConnection(ctx, "user-service", target)
	if err != nil {
		m.logger.Error("Failed to connect to user service", zap.Error(err))
		return nil, err
	}

	m.userClient = userpb.NewUserServiceClient(conn)
	m.logger.Info("Connected to user service", zap.String("target", target))
	return m.userClient, nil
}

// GetProductClient 获取产品服务的 gRPC 客户端
func (m *ClientManager) GetProductClient(ctx context.Context) (productpb.ProductServiceClient, error) {
	m.mu.RLock()
	if m.productClient != nil {
		client := m.productClient
		m.mu.RUnlock()
		return client, nil
	}
	m.mu.RUnlock()

	// 双重检查锁，避免重复创建客户端
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.productClient != nil {
		return m.productClient, nil
	}

	// 从配置中读取产品服务地址
	target := fmt.Sprintf("%s:%d",
		viper.GetString("services.product.host"),
		viper.GetInt("services.product.port"))

	conn, err := m.factory.GetConnection(ctx, "product-service", target)
	if err != nil {
		m.logger.Error("Failed to connect to product service", zap.Error(err))
		return nil, err
	}

	m.productClient = productpb.NewProductServiceClient(conn)
	m.logger.Info("Connected to product service", zap.String("target", target))
	return m.productClient, nil
}

// Close 关闭所有客户端连接
func (m *ClientManager) Close() {
	m.factory.CloseAll()
	m.logger.Info("All gRPC client connections closed")
}
