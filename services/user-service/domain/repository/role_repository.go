package repository

import (
	"context"

	"mall-go/services/user-service/domain/model"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// 基本CRUD操作
	Save(ctx context.Context, role *model.Role) error
	FindByID(ctx context.Context, id string) (*model.Role, error)
	FindByName(ctx context.Context, name string) (*model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id string) error
	
	// 分页查询
	FindAll(ctx context.Context, page, size int) ([]*model.Role, int64, error)
	
	// 角色权限相关
	AddRolePermission(ctx context.Context, roleID, permissionID string) error
	RemoveRolePermission(ctx context.Context, roleID, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]model.Permission, error)
	
	// 查询用户的所有角色
	FindByUserID(ctx context.Context, userID string) ([]*model.Role, error)
	
	// 更新角色状态
	UpdateStatus(ctx context.Context, id string, status model.RoleStatus) error
}