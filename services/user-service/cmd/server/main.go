package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mall-go/pkg/database"
	"mall-go/pkg/registry"
	"mall-go/services/user-service/api/handler"
	"mall-go/services/user-service/api/router"
	"mall-go/services/user-service/application/service"
	"mall-go/services/user-service/infrastructure/config"
	"mall-go/services/user-service/infrastructure/persistence/mysql"
)

func main() {
	// 1. 加载配置文件
	configPath := filepath.Join("configs", "config.yml")
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. 初始化数据库连接
	if err := database.InitMySQL(
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.Database,
	); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 4. 自动迁移数据库结构
	if err := migrateDB(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 5. 初始化仓储
	userRepo := mysql.NewUserRepository()
	roleRepo := mysql.NewRoleRepository()
	permissionRepo := mysql.NewPermissionRepository()

	// 6. 初始化应用服务
	userService := service.NewUserService(userRepo, roleRepo, permissionRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	permissionService := service.NewPermissionService(permissionRepo, roleRepo)

	// 7. 初始化API处理器
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// 8. 创建Gin引擎
	r := gin.Default()

	// 9. 设置路由
	router.SetupRouter(r, userHandler, roleHandler, permissionHandler)

	// 10. 注册到Consul服务注册中心
	var serviceRegistrationLifecycle *registry.ServiceRegistrationLifecycle
	if cfg.Consul.Address != "" {
		// 创建Consul服务注册
		consulRegistry, err := registry.NewConsulServiceRegistry(cfg.Consul.Address)
		if err != nil {
			log.Printf("创建Consul客户端失败: %v, 将不会注册到服务发现中心", err)
		} else {
			// 获取本机内网IP
			host, err := getLocalIP()
			if err != nil {
				log.Printf("获取本机IP失败: %v, 使用localhost代替", err)
				host = "localhost"
			}

			// 创建服务实例
			instanceID := fmt.Sprintf("%s-%s", cfg.Consul.ServiceName, uuid.New().String())
			serviceInstance := &registry.ServiceInstance{
				ID:      instanceID,
				Name:    cfg.Consul.ServiceName,
				Address: host,
				Port:    cfg.Server.Port,
				Tags:    cfg.Consul.Tags,
				Meta:    cfg.Consul.Meta,
			}

			// 创建服务注册生命周期管理
			serviceRegistrationLifecycle = registry.NewServiceRegistrationLifecycle(consulRegistry, serviceInstance)

			// 启动服务注册
			if err := serviceRegistrationLifecycle.Start(); err != nil {
				log.Printf("服务注册失败: %v", err)
			} else {
				log.Printf("服务已注册到Consul: %s", cfg.Consul.Address)
			}
		}
	}

	// 11. 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// 12. 启动服务器（非阻塞）
	go func() {
		log.Printf("Server is running on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 13. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 从Consul注销服务
	if serviceRegistrationLifecycle != nil {
		if err := serviceRegistrationLifecycle.Stop(); err != nil {
			log.Printf("从Consul注销服务失败: %v", err)
		} else {
			log.Println("服务已从Consul注销")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// 数据库迁移
func migrateDB() error {
	// 迁移数据库表结构
	if err := database.DB.AutoMigrate(
		&mysql.UserEntity{},
		&mysql.RoleEntity{},
		&mysql.PermissionEntity{},
		&mysql.UserRoleRelation{},
		&mysql.RolePermissionRelation{},
	); err != nil {
		return err
	}
	return nil
}

// 获取本机内网IP
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid local IP address found")
}
