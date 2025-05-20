package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"mall-go/services/inventory-service/infrastructure/config"
)

// NewDBConnection 创建数据库连接
func NewDBConnection(config config.DatabaseConfig) (*sql.DB, error) {
	// 构建连接字符串
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// 连接数据库
	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	return db, nil
}
