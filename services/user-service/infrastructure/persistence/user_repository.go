package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"mall-go/pkg/concurrency"
	"mall-go/pkg/monitoring"
	"mall-go/services/user-service/domain/model"
)

// UserRepository is the implementation of the user repository interface
type UserRepository struct {
	db *sqlx.DB
	// 并发限制器
	limiter *concurrency.ConcurrencyLimiter
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	// 设置最大并发数为100，可根据系统负载调整
	return &UserRepository{
		db:      db,
		limiter: concurrency.NewConcurrencyLimiter(100),
	}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "SELECT", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "SELECT", "users", nil)

	// 使用并发限制
	var user model.User
	err := r.limiter.ExecuteWithLimit(ctx, func() error {
		// 查询包含版本号的用户记录
		query := `SELECT id, username, password, email, nick_name, phone, icon, status, version, note, 
                  created_at, updated_at, last_login FROM users WHERE id = ? AND deleted_at IS NULL`
		return r.db.GetContext(ctx, &user, query, id)
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "SELECT", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "SELECT", "users", nil)

	// 使用并发限制
	var user model.User
	err := r.limiter.ExecuteWithLimit(ctx, func() error {
		query := `SELECT id, username, password, email, nick_name, phone, icon, status, version, note, 
                  created_at, updated_at, last_login FROM users WHERE username = ? AND deleted_at IS NULL`
		return r.db.GetContext(ctx, &user, query, username)
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "SELECT", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "SELECT", "users", nil)

	// 使用并发限制
	var user model.User
	err := r.limiter.ExecuteWithLimit(ctx, func() error {
		query := `SELECT id, username, password, email, nick_name, phone, icon, status, version, note, 
                  created_at, updated_at, last_login FROM users WHERE email = ? AND deleted_at IS NULL`
		return r.db.GetContext(ctx, &user, query, email)
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "INSERT", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "INSERT", "users", nil)

	// 使用并发限制
	return r.limiter.ExecuteWithLimit(ctx, func() error {
		// 设置初始版本号和时间
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now
		user.Version = 1

		query := `INSERT INTO users (id, username, password, email, nick_name, phone, icon, status, version, note, 
                 created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Password, user.Email,
			user.NickName, user.Phone, user.Icon, user.Status, user.Version,
			user.Note, user.CreatedAt, user.UpdatedAt)
		return err
	})
}

// Update updates an existing user with optimistic locking
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "UPDATE", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "UPDATE", "users", nil)

	// 使用并发限制
	return r.limiter.ExecuteWithLimit(ctx, func() error {
		// 使用乐观锁更新用户
		user.UpdatedAt = time.Now()
		expectedVersion := user.Version
		user.Version = expectedVersion + 1

		query := `UPDATE users SET email = ?, nick_name = ?, phone = ?, icon = ?, status = ?, 
                  note = ?, updated_at = ?, version = ? 
                  WHERE id = ? AND version = ? AND deleted_at IS NULL`
		result, err := r.db.ExecContext(ctx, query, user.Email, user.NickName, user.Phone,
			user.Icon, user.Status, user.Note, user.UpdatedAt,
			user.Version, user.ID, expectedVersion)
		if err != nil {
			return err
		}

		// 检查是否有记录被更新
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		// 如果没有更新记录，说明可能发生了并发冲突
		if rowsAffected == 0 {
			return concurrency.ErrVersionConflict
		}

		return nil
	})
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "DELETE", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "DELETE", "users", nil)

	// 使用并发限制和软删除
	return r.limiter.ExecuteWithLimit(ctx, func() error {
		now := time.Now()
		query := `UPDATE users SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL`
		result, err := r.db.ExecContext(ctx, query, now, now, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("user with ID %s not found", id)
		}

		return nil
	})
}

// UpdateLastLogin updates a user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string, lastLogin time.Time) error {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "UPDATE", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "UPDATE", "users", nil)

	// 使用并发限制
	return r.limiter.ExecuteWithLimit(ctx, func() error {
		query := `UPDATE users SET last_login = ? WHERE id = ? AND deleted_at IS NULL`
		_, err := r.db.ExecContext(ctx, query, lastLogin, id)
		return err
	})
}

// FindAll finds all users with pagination
func (r *UserRepository) FindAll(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "SELECT", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "SELECT", "users", nil)

	// 使用并发限制
	var users []*model.User
	var total int64

	err := r.limiter.ExecuteWithLimit(ctx, func() error {
		// 获取总记录数
		countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
		if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
			return err
		}

		// 计算偏移量
		offset := (page - 1) * pageSize

		// 获取用户列表
		query := `SELECT id, username, email, nick_name, phone, icon, status, version, note, 
                 created_at, updated_at, last_login FROM users 
                 WHERE deleted_at IS NULL
                 ORDER BY created_at DESC LIMIT ? OFFSET ?`
		return r.db.SelectContext(ctx, &users, query, pageSize, offset)
	})

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Search searches for users by keyword
func (r *UserRepository) Search(ctx context.Context, keyword string, page, pageSize int) ([]*model.User, int64, error) {
	// 使用监控创建数据库操作段
	ctx, span, startTime := monitoring.StartDatabaseSegment(ctx, "SELECT", "users")
	defer monitoring.EndDatabaseSegment(span, startTime, "SELECT", "users", nil)

	// 使用并发限制
	var users []*model.User
	var total int64

	err := r.limiter.ExecuteWithLimit(ctx, func() error {
		// 构造LIKE参数
		likeKeyword := "%" + keyword + "%"

		// 获取总记录数
		countQuery := `SELECT COUNT(*) FROM users 
                      WHERE (username LIKE ? OR email LIKE ? OR nick_name LIKE ?) 
                      AND deleted_at IS NULL`
		if err := r.db.GetContext(ctx, &total, countQuery, likeKeyword, likeKeyword, likeKeyword); err != nil {
			return err
		}

		// 计算偏移量
		offset := (page - 1) * pageSize

		// 获取用户列表
		query := `SELECT id, username, email, nick_name, phone, icon, status, version, note, 
                 created_at, updated_at, last_login FROM users 
                 WHERE (username LIKE ? OR email LIKE ? OR nick_name LIKE ?) 
                 AND deleted_at IS NULL
                 ORDER BY created_at DESC LIMIT ? OFFSET ?`
		return r.db.SelectContext(ctx, &users, query, likeKeyword, likeKeyword, likeKeyword, pageSize, offset)
	})

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
