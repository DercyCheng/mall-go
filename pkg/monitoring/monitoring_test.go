package monitoring

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// 测试监控设置
func TestSetupMonitoring(t *testing.T) {
	// 创建Gin引擎
	r := gin.New()
	
	// 测试仅Prometheus配置
	t.Run("OnlyPrometheus", func(t *testing.T) {
		cfg := Config{
			ServiceName:      "test-service",
			EnablePrometheus: true,
			EnableTracing:    false,
		}
		
		err := SetupMonitoring(r, cfg)
		assert.NoError(t, err)
		
		// 验证Prometheus端点已注册
		routes := r.Routes()
		found := false
		for _, route := range routes {
			if route.Path == "/metrics" && route.Method == "GET" {
				found = true
				break
			}
		}
		assert.True(t, found, "Prometheus指标端点应当注册")
	})
	
	// 注意：跟踪模块需要Jaeger实例，但我们可以跳过相关测试
}

// 测试指标中间件
func TestMetricsMiddleware(t *testing.T) {
	// 重置指标计数器，避免之前测试的影响
	httpRequestsTotal.Reset()
	httpRequestDuration.Reset()
	
	// 创建Gin引擎并设置中间件
	r := gin.New()
	r.Use(MetricsMiddleware())
	
	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		// 模拟一些处理时间
		time.Sleep(10 * time.Millisecond)
		c.JSON(200, gin.H{"message": "ok"})
	})
	
	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	// 执行请求
	r.ServeHTTP(w, req)
	
	// 验证请求成功
	assert.Equal(t, 200, w.Code)
	
	// 验证指标已记录
	// 注意：这里我们不直接从收集器中读取数据，因为那会涉及到内部API
	// 在实际测试中，您可能需要进一步验证指标数据
	assert.Greater(t, testutil.CollectAndCount(httpRequestsTotal), 0, "请求计数器应增加")
	assert.Greater(t, testutil.CollectAndCount(httpRequestDuration), 0, "请求时长应记录")
}

// 测试链路追踪中间件
func TestTracingMiddleware(t *testing.T) {
	// 创建一个模拟的跟踪提供程序
	// 注意：在完整测试中，您应该创建一个适当的模拟，此处简化处理
	originalTracer := tracer
	defer func() { tracer = originalTracer }()
	
	// 创建Gin引擎并设置中间件
	r := gin.New()
	r.Use(TracingMiddleware())
	
	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		// 验证跟踪上下文已传递
		ctx, exists := c.Get("tracingContext")
		assert.True(t, exists)
		assert.NotNil(t, ctx)
		
		c.JSON(200, gin.H{"message": "ok"})
	})
	
	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	// 执行请求
	r.ServeHTTP(w, req)
	
	// 验证请求成功
	assert.Equal(t, 200, w.Code)
}

// 测试数据库监控功能
func TestDatabaseMonitoring(t *testing.T) {
	// 重置指标计数器
	databaseQueryTotal.Reset()
	databaseQueryDuration.Reset()
	
	// 创建上下文
	ctx := context.Background()
	
	// 开始数据库操作
	ctx, span, start := StartDatabaseSegment(ctx, "SELECT", "users")
	
	// 模拟一些处理时间
	time.Sleep(10 * time.Millisecond)
	
	// 结束操作
	EndDatabaseSegment(span, start, "SELECT", "users", nil)
	
	// 验证指标已记录
	assert.Greater(t, testutil.CollectAndCount(databaseQueryTotal), 0, "数据库查询计数器应增加")
	assert.Greater(t, testutil.CollectAndCount(databaseQueryDuration), 0, "数据库查询时长应记录")
}