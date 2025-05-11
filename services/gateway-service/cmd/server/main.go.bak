package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"

	"mall-go/services/gateway-service/api/middleware"
	"mall-go/services/gateway-service/api/proxy"
	"mall-go/services/gateway-service/infrastructure/config"
)

func main() {
	// Load configuration
	configPath := getConfigPath()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(getGinMode(cfg.Server.Mode))

	// Initialize router
	router := gin.Default()

	// Setup middleware
	middleware.SetupMiddleware(router, cfg)

	// Initialize service proxy
	serviceProxy, err := proxy.NewServiceProxy(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize service proxy: %v", err)
	}

	// Register routes
	serviceProxy.RegisterRoutes(router)

	// Register service with Consul if enabled
	var serviceID string
	if cfg.Registry.Type == "consul" && cfg.Registry.Address != "" {
		serviceID, err = registerService(cfg.Registry, cfg.Server.Port)
		if err != nil {
			log.Printf("Warning: Failed to register service with Consul: %v", err)
		} else {
			log.Printf("Service registered with Consul (ID: %s)", serviceID)
			// Deregister service when shutting down
			defer deregisterService(cfg.Registry.Address, serviceID)
		}
	}

	// Start the server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a separate goroutine
	go func() {
		log.Printf("Starting gateway server on port %d...", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gateway server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Gateway server exited gracefully")
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	// Check if config path is provided via environment variable
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		return configPath
	}

	// Default to services/gateway-service/configs/config.yaml
	return filepath.Join("services", "gateway-service", "configs", "config.yaml")
}

// getGinMode converts our app mode to Gin mode
func getGinMode(mode string) string {
	switch mode {
	case "debug":
		return gin.DebugMode
	case "release":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
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
