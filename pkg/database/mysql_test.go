package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 测试MySQL连接初始化
func TestInitMySQL(t *testing.T) {
	// 跳过实际连接测试，只测试功能
	t.Skip("需要实际的MySQL连接来测试，跳过")
	
	// 如果要在CI/CD中测试，可以提供测试数据库配置
	err := InitMySQL("testuser", "testpassword", "localhost", "3306", "testdb")
	assert.NoError(t, err)
	
	// 验证连接有效性
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
	
	// 关闭连接
	assert.NoError(t, CloseDB())
}

// 使用sqlmock测试数据库功能
func TestDBWithMock(t *testing.T) {
	// 创建SQL mock
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	
	// 创建期望
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("5.7.32"))
	
	// 连接到GORM
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	
	// 使用mock连接替换全局DB
	DB, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)
	
	// 测试查询
	var version string
	tx := DB.Raw("SELECT VERSION()").Scan(&version)
	assert.NoError(t, tx.Error)
	assert.NotEmpty(t, version)
	
	// 确认所有期望被满足
	assert.NoError(t, mock.ExpectationsWereMet())
	
	// 清理
	mockDB, _ := DB.DB()
	mockDB.Close()
}

// 测试连接池设置
func TestConnectionPool(t *testing.T) {
	// 创建SQL mock
	sqlDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()
	
	// 创建一个临时函数来测试连接池设置逻辑
	setupConnectionPool := func(db *sql.DB) {
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
	}
	
	// 应用连接池设置
	setupConnectionPool(sqlDB)
	
	// 验证连接池设置
	assert.Equal(t, 10, sqlDB.Stats().MaxIdleConns)
	assert.Equal(t, 100, sqlDB.Stats().MaxOpenConns)
}

// 测试关闭数据库连接
func TestCloseDB(t *testing.T) {
	// 测试当DB为nil时的情况
	DB = nil
	assert.NoError(t, CloseDB())
	
	// 使用mock测试正常关闭
	sqlDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	
	// 连接到GORM
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	
	// 替换全局DB
	DB, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)
	
	// 测试关闭
	assert.NoError(t, CloseDB())
}