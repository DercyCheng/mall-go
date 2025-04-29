package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"

	"mall-go/services/gateway-service/infrastructure/config"
)

// ServiceProxy handles the routing of requests to appropriate backend services
type ServiceProxy struct {
	consulClient *api.Client
	config       *config.Config
}

// NewServiceProxy creates a new instance of ServiceProxy
func NewServiceProxy(cfg *config.Config) (*ServiceProxy, error) {
	// Create a default Consul client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Registry.Address
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &ServiceProxy{
		consulClient: client,
		config:       cfg,
	}, nil
}

// ProxyUserService handles routing requests to the user service
func (p *ServiceProxy) ProxyUserService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.User.Name, p.config.Services.User.GetTimeout())
}

// RegisterRoutes registers all proxy routes
func (p *ServiceProxy) RegisterRoutes(r *gin.Engine) {
	// User service routes
	userRoutes := r.Group("/api/user")
	userRoutes.Any("/*path", p.ProxyUserService)

	// Add more service routes as they are implemented
	// productRoutes := r.Group("/api/product")
	// productRoutes.Any("/*path", p.ProxyProductService)

	// Add a health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// proxyRequest handles the forwarding of requests to backend services
func (p *ServiceProxy) proxyRequest(c *gin.Context, serviceName string, timeout time.Duration) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
	defer cancel()

	// Get service instance from Consul
	serviceURL, err := p.getServiceURL(serviceName)
	if err != nil {
		log.Printf("Error getting service URL for %s: %v", serviceName, err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Service %s is unavailable", serviceName),
		})
		return
	}

	// Parse the target URL
	target, err := url.Parse(serviceURL)
	if err != nil {
		log.Printf("Error parsing service URL: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Update the request URL path
	originalPath := c.Request.URL.Path
	// Remove the service prefix (like /api/user)
	pathParts := strings.SplitN(originalPath, "/", 4)
	var newPath string
	if len(pathParts) >= 4 {
		// /api/user/something -> /something
		newPath = "/" + pathParts[3]
	} else {
		newPath = "/"
	}
	c.Request.URL.Path = newPath

	// Update request Host and headers
	c.Request.Host = target.Host
	c.Request.URL.Host = target.Host
	c.Request.URL.Scheme = target.Scheme
	c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
	c.Request.Header.Set("X-Forwarded-Proto", c.Request.URL.Scheme)
	c.Request.Header.Set("X-Forwarded-Path", originalPath)

	// Use the context with timeout
	c.Request = c.Request.WithContext(ctx)

	// Forward the request to the service
	proxy.ServeHTTP(c.Writer, c.Request)
}

// getServiceURL retrieves a service URL from Consul
func (p *ServiceProxy) getServiceURL(serviceName string) (string, error) {
	// Get healthy service instances
	services, _, err := p.consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get service instances: %w", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	// Simple round-robin for now - could be enhanced with proper load balancing
	service := services[0]

	// Construct the service URL
	serviceURL := fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)
	return serviceURL, nil
}
