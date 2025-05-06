package main

import (
	"context"
	"fmt"
	"log"
	"mall-go/services/user-service/api/handler"
	"mall-go/services/user-service/api/middleware"
	"mall-go/services/user-service/api/router"
	"mall-go/services/user-service/application/service"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
	"mall-go/services/user-service/infrastructure/config"
	persistence "mall-go/services/user-service/infrastructure/persistence"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// UserRepositoryAdapter adapts the persistence implementation to the domain interface
type UserRepositoryAdapter struct {
	repo *persistence.UserRepository
}

func NewUserRepositoryAdapter(db *sqlx.DB) repository.UserRepository {
	return &UserRepositoryAdapter{
		repo: persistence.NewUserRepository(db),
	}
}

// Implement all methods required by the repository.UserRepository interface
func (a *UserRepositoryAdapter) FindByID(ctx context.Context, id string) (*model.User, error) {
	return a.repo.FindByID(ctx, id)
}

func (a *UserRepositoryAdapter) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return a.repo.FindByUsername(ctx, username)
}

func (a *UserRepositoryAdapter) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return a.repo.FindByEmail(ctx, email)
}

func (a *UserRepositoryAdapter) Save(ctx context.Context, user *model.User) error {
	return a.repo.Create(ctx, user)
}

func (a *UserRepositoryAdapter) Update(ctx context.Context, user *model.User) error {
	return a.repo.Update(ctx, user)
}

func (a *UserRepositoryAdapter) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}

func (a *UserRepositoryAdapter) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return a.repo.FindAll(ctx, page, pageSize)
}

func (a *UserRepositoryAdapter) Search(ctx context.Context, query string, page, pageSize int) ([]*model.User, int64, error) {
	return a.repo.Search(ctx, query, page, pageSize)
}

// MockRoleRepository is a temporary implementation for the role repository interface
type MockRoleRepository struct{}

func NewMockRoleRepository() repository.RoleRepository {
	return &MockRoleRepository{}
}

// Implement stub methods for the RoleRepository interface
func (r *MockRoleRepository) Save(ctx context.Context, role *model.Role) error {
	return nil
}

func (r *MockRoleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	return nil, nil
}

func (r *MockRoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	return nil, nil
}

func (r *MockRoleRepository) Update(ctx context.Context, role *model.Role) error {
	return nil
}

func (r *MockRoleRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *MockRoleRepository) List(ctx context.Context) ([]*model.Role, error) {
	return nil, nil
}

func (r *MockRoleRepository) GetUserRoles(ctx context.Context, userID string) ([]*model.Role, error) {
	return nil, nil
}

func (r *MockRoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	return nil
}

func (r *MockRoleRepository) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	return nil
}

func main() {
	// 加载配置
	configPath := "./services/user-service/configs/config.yaml"
	// 检查环境变量
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 注册服务到Consul
	var serviceID string
	if cfg.Registry.Type == "consul" && cfg.Registry.Address != "" {
		serviceID, err = registerService(cfg.Registry, cfg.Server.Port)
		if err != nil {
			log.Printf("警告: 注册服务到Consul失败: %v", err)
		} else {
			log.Printf("服务已注册到Consul (ID: %s)", serviceID)
			// 服务关闭时注销
			defer deregisterService(cfg.Registry.Address, serviceID)
		}
	}

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if cfg.Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	}

	// 初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// 使用sqlx连接数据库
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("初始化数据库连接失败: %v", err)
	}
	defer db.Close()

	// 设置连接池参数
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.Database.GetConnMaxLifetime())

	// 初始化仓库 - 使用适配器模式来解决接口不匹配问题
	userRepo := NewUserRepositoryAdapter(db)
	roleRepo := NewMockRoleRepository()

	// 初始化服务
	userService := service.NewUserService(userRepo, roleRepo, cfg.Auth.JWTSecret, cfg.Auth.GetTokenExpiry())
	roleService := service.NewRoleService(roleRepo)

	// 初始化路由处理器
	userHandler := handler.NewUserHandler(userService, roleService)

	// 创建Gin引擎
	engine := gin.New()

	// 使用自定义中间件
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.RecoveryWithZap())
	engine.Use(middleware.Cors())

	// 注册路由
	userRouter := router.NewUserRouter(userHandler, cfg.Auth.JWTSecret)
	userRouter.Register(engine)

	// 启动HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: engine,
	}

	// 优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听端口失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭服务器...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已安全关闭")
}

// registerService registers the service with Consul
func registerService(cfg config.RegistryConfig, port int) (string, error) {
	// Create a default Consul client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Address
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create Consul client: %w", err)
	}

	// Get hostname for service registration
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	// Create a unique service ID
	serviceID := fmt.Sprintf("%s-%s", cfg.ServiceName, uuid.New().String())

	// Create service registration
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    cfg.ServiceName,
		Tags:    cfg.Tags,
		Port:    port,
		Address: hostname,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", hostname, port),
			Timeout:                        "5s",
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	// Register the service
	if err := client.Agent().ServiceRegister(registration); err != nil {
		return "", err
	}

	return serviceID, nil
}

// deregisterService deregisters the service from Consul
func deregisterService(address, serviceID string) {
	// Create a default Consul client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = address
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Printf("Failed to create Consul client for deregistration: %v", err)
		return
	}

	// Deregister the service
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
		return
	}

	log.Printf("Service deregistered from Consul (ID: %s)", serviceID)
}
