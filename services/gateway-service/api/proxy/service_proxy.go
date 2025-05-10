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

// ProxyProductService handles routing requests to the product service
func (p *ServiceProxy) ProxyProductService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Product.Name, p.config.Services.Product.GetTimeout())
}

// ProxyOrderService handles routing requests to the order service
func (p *ServiceProxy) ProxyOrderService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Order.Name, p.config.Services.Order.GetTimeout())
}

// ProxyCartService handles routing requests to the cart service
func (p *ServiceProxy) ProxyCartService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Cart.Name, p.config.Services.Cart.GetTimeout())
}

// ProxyInventoryService handles routing requests to the inventory service
func (p *ServiceProxy) ProxyInventoryService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Inventory.Name, p.config.Services.Inventory.GetTimeout())
}

// ProxyPaymentService handles routing requests to the payment service
func (p *ServiceProxy) ProxyPaymentService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Payment.Name, p.config.Services.Payment.GetTimeout())
}

// ProxySearchService handles routing requests to the search service
func (p *ServiceProxy) ProxySearchService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Search.Name, p.config.Services.Search.GetTimeout())
}

// ProxyPromotionService handles routing requests to the promotion service
func (p *ServiceProxy) ProxyPromotionService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Promotion.Name, p.config.Services.Promotion.GetTimeout())
}

// ProxyContentService handles routing requests to the content service
func (p *ServiceProxy) ProxyContentService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Content.Name, p.config.Services.Content.GetTimeout())
}

// ProxyNotificationService handles routing requests to the notification service
func (p *ServiceProxy) ProxyNotificationService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Notification.Name, p.config.Services.Notification.GetTimeout())
}

// ProxyRecommendationService handles routing requests to the recommendation service
func (p *ServiceProxy) ProxyRecommendationService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Recommendation.Name, p.config.Services.Recommendation.GetTimeout())
}

// ProxyAdminService handles routing requests to the admin service
func (p *ServiceProxy) ProxyAdminService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Admin.Name, p.config.Services.Admin.GetTimeout())
}

// ProxyPortalService handles routing requests to the portal service
func (p *ServiceProxy) ProxyPortalService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Portal.Name, p.config.Services.Portal.GetTimeout())
}

// ProxyAuthService handles routing requests to the auth service
func (p *ServiceProxy) ProxyAuthService(c *gin.Context) {
	p.proxyRequest(c, p.config.Services.Auth.Name, p.config.Services.Auth.GetTimeout())
}

// RegisterRoutes registers all proxy routes
func (p *ServiceProxy) RegisterRoutes(r *gin.Engine) {
	// User service routes
	userRoutes := r.Group("/api/user")
	userRoutes.Any("/*path", p.ProxyUserService)

	// Product service routes
	productRoutes := r.Group("/api/product")
	productRoutes.Any("/*path", p.ProxyProductService)
	
	// Order service routes
	orderRoutes := r.Group("/api/order")
	orderRoutes.Any("/*path", p.ProxyOrderService)
	
	// Cart service routes
	cartRoutes := r.Group("/api/cart")
	cartRoutes.Any("/*path", p.ProxyCartService)
	
	// Inventory service routes
	inventoryRoutes := r.Group("/api/inventory")
	inventoryRoutes.Any("/*path", p.ProxyInventoryService)
	
	// Payment service routes
	paymentRoutes := r.Group("/api/payment")
	paymentRoutes.Any("/*path", p.ProxyPaymentService)
	
	// Search service routes
	searchRoutes := r.Group("/api/search")
	searchRoutes.Any("/*path", p.ProxySearchService)
	
	// Promotion service routes
	promotionRoutes := r.Group("/api/promotion")
	promotionRoutes.Any("/*path", p.ProxyPromotionService)
	
	// Content service routes
	contentRoutes := r.Group("/api/content")
	contentRoutes.Any("/*path", p.ProxyContentService)
	
	// Notification service routes
	notificationRoutes := r.Group("/api/notification")
	notificationRoutes.Any("/*path", p.ProxyNotificationService)
	
	// Recommendation service routes
	recommendationRoutes := r.Group("/api/recommendation")
	recommendationRoutes.Any("/*path", p.ProxyRecommendationService)
	
	// Admin service routes
	adminRoutes := r.Group("/api/admin")
	adminRoutes.Any("/*path", p.ProxyAdminService)
	
	// Portal service routes
	portalRoutes := r.Group("/api/portal")
	portalRoutes.Any("/*path", p.ProxyPortalService)
	
	// Auth service routes
	authRoutes := r.Group("/api/auth")
	authRoutes.Any("/*path", p.ProxyAuthService)

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
