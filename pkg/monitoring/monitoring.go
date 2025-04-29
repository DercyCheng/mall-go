// Package monitoring provides instrumentations for metrics and tracing
package monitoring

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// 全局指标
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	databaseQueryTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_query_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)

	databaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	// 全局跟踪器
	tracer = otel.Tracer("mall-go")
)

// Config 监控配置
type Config struct {
	// 服务名称
	ServiceName string
	// 是否启用Prometheus
	EnablePrometheus bool
	// 是否启用跟踪
	EnableTracing bool
	// Jaeger采集器端点
	JaegerEndpoint string
}

// SetupMonitoring 设置监控
func SetupMonitoring(engine *gin.Engine, cfg Config) error {
	// 设置Prometheus端点
	if cfg.EnablePrometheus {
		engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// 设置跟踪
	if cfg.EnableTracing {
		if err := setupTracing(cfg); err != nil {
			return err
		}
	}

	return nil
}

// 设置链路跟踪
func setupTracing(cfg Config) error {
	// 创建Jaeger导出器
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)))
	if err != nil {
		return err
	}

	// 创建资源
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.ServiceName),
	)

	// 创建跟踪提供程序
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource),
	)

	// 设置全局跟踪提供程序
	otel.SetTracerProvider(tp)

	return nil
}

// MetricsMiddleware 创建Gin中间件，用于收集HTTP请求指标
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// 处理请求
		c.Next()

		// 记录指标
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, string(rune(status))).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

// TracingMiddleware 创建Gin中间件，用于跟踪HTTP请求
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// 开始新的跟踪Span
		ctx, span := tracer.Start(c.Request.Context(), "http_request",
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.path", path),
			),
		)
		defer span.End()

		// 将Span上下文传递到请求中
		c.Request = c.Request.WithContext(ctx)
		c.Set("tracingContext", ctx)

		// 处理请求
		c.Next()

		// 添加响应信息
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)
	}
}

// StartDatabaseSegment 开始数据库操作的跟踪段
func StartDatabaseSegment(ctx context.Context, operation, table string) (context.Context, trace.Span, time.Time) {
	start := time.Now()
	ctx, span := tracer.Start(ctx, "database_operation",
		trace.WithAttributes(
			attribute.String("db.operation", operation),
			attribute.String("db.table", table),
		),
	)
	return ctx, span, start
}

// EndDatabaseSegment 结束数据库操作的跟踪段
func EndDatabaseSegment(span trace.Span, start time.Time, operation, table string, err error) {
	// 记录跟踪信息
	duration := time.Since(start)
	span.SetAttributes(attribute.Float64("db.duration_ms", float64(duration.Milliseconds())))
	if err != nil {
		span.RecordError(err)
	}
	span.End()

	// 记录指标
	databaseQueryTotal.WithLabelValues(operation, table).Inc()
	databaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}
