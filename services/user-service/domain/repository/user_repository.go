package repository

import (
	"context"

	"mall-go/services/user-service/domain/model"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基本CRUD操作
	Save(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	
	// 分页查询
	FindAll(ctx context.Context, page, size int) ([]*model.User, int64, error)
	
	// 用户角色相关
	AddUserRole(ctx context.Context, userID, roleID string) error
	RemoveUserRole(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]model.Role, error)
	
	// 更新用户状态
	UpdateStatus(ctx context.Context, id string, status model.UserStatus) error
	
	// 更新用户登录时间
	UpdateLastLoginTime(ctx context.Context, id string) error
	
	// 搜索用户
	Search(ctx context.Context, keyword string, page, size int) ([]*model.User, int64, error)
}