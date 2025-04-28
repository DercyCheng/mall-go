package repository

import (
	"context"

	"mall-go/services/user-service/domain/model"
)

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	// 基本CRUD操作
	Save(ctx context.Context, permission *model.Permission) error
	FindByID(ctx context.Context, id string) (*model.Permission, error)
	FindByValue(ctx context.Context, value string) (*model.Permission, error)
	Update(ctx context.Context, permission *model.Permission) error
	Delete(ctx context.Context, id string) error
	
	// 分页查询
	FindAll(ctx context.Context, page, size int) ([]*model.Permission, int64, error)
	
	// 根据类型查询权限
	FindByType(ctx context.Context, permType model.PermissionType) ([]*model.Permission, error)
	
	// 查询角色的所有权限
	FindByRoleID(ctx context.Context, roleID string) ([]*model.Permission, error)
	
	// 查询用户的所有权限
	FindByUserID(ctx context.Context, userID string) ([]*model.Permission, error)
	
	// 更新权限状态
	UpdateStatus(ctx context.Context, id string, status model.PermissionStatus) error
}