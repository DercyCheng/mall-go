package main

import (
	"context"
	"flag"
	"fmt"
	"mall-go/services/product-service/application/service"
	domainService "mall-go/services/product-service/domain/service"
	mysqlx "mall-go/services/product-service/infrastructure/persistence/mysql"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mall-go/services/product-service/api/handler"
	"mall-go/services/product-service/api/router"
)

// 获取服务根目录
func getServiceRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	// 获取cmd/server的上层目录
	return filepath.Dir(filepath.Dir(filepath.Dir(b)))
}

func main() {
	// 解析命令行参数
	configFile := flag.String("config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	// 初始化配置
	initConfig(*configFile)

	// 初始化日志
	initLogger()

	// 初始化数据库连接
	db := initDB()

	// 初始化仓储
	productRepo := mysqlx.NewProductRepository(db)
	brandRepo := mysqlx.NewBrandRepository(db)
	categoryRepo := mysqlx.NewCategoryRepository(db)

	// 初始化领域服务
	productDomainService := domainService.NewProductDomainService(
		productRepo,
		brandRepo,
		categoryRepo,
	)

	// 初始化应用服务
	productService := service.NewProductService(
		productRepo,
		brandRepo,
		categoryRepo,
		productDomainService,
		nil, // TODO: 实现事件发布器
	)

	// 初始化 HTTP 处理器
	productHandler := handler.NewProductHandler(productService)

	// 初始化 Gin 引擎
	gin.SetMode(viper.GetString("server.mode"))
	r := gin.Default()

	// 设置路由
	router.SetupRouter(r, productHandler)

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("server.port")),
		Handler: r,
	}

	// 优雅关闭服务
	go func() {
		logrus.Infof("服务启动，监听端口: %d", viper.GetInt("server.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("服务启动失败: %s", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("服务关闭失败: %s", err)
	}
	logrus.Info("服务已关闭")
}

// 初始化配置
func initConfig(configFile string) {
	// 如果配置文件路径是相对路径，则基于服务根目录解析
	if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(getServiceRootDir(), configFile)
	}
	
	logrus.Infof("加载配置文件: %s", configFile)
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}
}

// 初始化日志
func initLogger() {
	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// 初始化数据库连接
func initDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name"),
	)

	var gormLogger logger.Interface
	if viper.GetBool("database.debug") {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("获取底层数据库连接失败: %w", err))
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(viper.GetInt("database.maxIdleConns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("database.maxOpenConns"))
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("database.connMaxLifetime")) * time.Second)

	return db
}
