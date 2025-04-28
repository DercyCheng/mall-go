package router

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/registry"
	"mall-go/services/gateway-service/api/middleware"
	"mall-go/services/gateway-service/infrastructure/config"
)

// Gateway API网关路由
type Gateway struct {
	engine         *gin.Engine
	consulRegistry registry.ServiceRegistry
}

// NewGateway 创建API网关
func NewGateway(engine *gin.Engine, consulRegistry registry.ServiceRegistry) *Gateway {
	return &Gateway{
		engine:         engine,
		consulRegistry: consulRegistry,
	}
}

// SetupRoutes 设置路由
func (g *Gateway) SetupRoutes() {
	// 设置中间件
	g.engine.Use(middleware.CORSMiddleware())
	g.engine.Use(middleware.RateLimitMiddleware())

	// 健康检查接口
	g.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// 用户服务路由
	g.setupServiceRoutes(config.GlobalConfig.Services.UserService, "user-service")

	// 商品服务路由
	g.setupServiceRoutes(config.GlobalConfig.Services.ProductService, "product-service")

	// 未匹配到的路由
	g.engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "API not found",
		})
	})
}

// setupServiceRoutes 设置服务路由
func (g *Gateway) setupServiceRoutes(serviceConfig config.ServiceConfig, serviceName string) {
	if serviceConfig.Path == "" {
		log.Printf("服务[%s]的路径配置为空，跳过路由设置", serviceName)
		return
	}

	// 构建路由路径
	routePath := serviceConfig.Path
	if !strings.HasSuffix(routePath, "/*") {
		if strings.HasSuffix(routePath, "/") {
			routePath = routePath + "*path"
		} else {
			routePath = routePath + "/*path"
		}
	}

	// 注册路由
	g.engine.Any(routePath, func(c *gin.Context) {
		// 获取服务实例
		instances, err := g.consulRegistry.DiscoverServices(serviceName)
		if err != nil || len(instances) == 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    503,
				"message": fmt.Sprintf("服务[%s]不可用", serviceName),
			})
			return
		}

		// 简单负载均衡 - 随机选择一个实例
		instance := instances[0] // 此处可以实现更复杂的负载均衡策略

		// 构建目标URL
		targetURL := &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", instance.Address, instance.Port),
		}

		// 获取匹配到的路径参数
		path := c.Param("path")
		if path == "" && c.Request.URL.Path != serviceConfig.Path {
			path = strings.TrimPrefix(c.Request.URL.Path, serviceConfig.Path)
		}

		// 创建反向代理
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		
		// 设置转发请求的修改器
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host

			// 处理路径前缀
			if serviceConfig.StripPrefix {
				req.URL.Path = path
			} else {
				req.URL.Path = serviceConfig.Path + path
			}

			// 处理查询参数
			req.URL.RawQuery = c.Request.URL.RawQuery

			// 保留原始主机头
			req.Host = targetURL.Host

			// 添加X-Forwarded-*头
			req.Header.Set("X-Forwarded-Host", c.Request.Host)
			req.Header.Set("X-Forwarded-Proto", c.Request.URL.Scheme)
			req.Header.Set("X-Forwarded-For", c.ClientIP())
			
			// 添加网关版本信息
			req.Header.Set("X-Gateway-Version", "v1.0")
			
			// 如果需要认证但路径不在跳过列表中，则检查JWT
			if !middleware.SkipJWTAuth(c.FullPath()) {
				authHeader := c.GetHeader("Authorization")
				if authHeader != "" {
					req.Header.Set("Authorization", authHeader)
				}
			}
		}

		// 设置错误处理
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("代理请求错误: %v", err)
			c.JSON(http.StatusBadGateway, gin.H{
				"code":    502,
				"message": "网关错误: " + err.Error(),
			})
		}

		// 设置修改响应的处理器
		proxy.ModifyResponse = func(resp *http.Response) error {
			// 添加CORS头
			if config.GlobalConfig.CORS.Enabled {
				resp.Header.Set("Access-Control-Allow-Origin", "*")
			}
			return nil
		}

		// 执行代理请求
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("已为服务[%s]设置路由: %s", serviceName, routePath)
}