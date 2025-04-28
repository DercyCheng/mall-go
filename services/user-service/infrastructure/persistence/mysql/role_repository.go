package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mall-go/pkg/database"
	"mall-go/services/user-service/domain/model"
	"mall-go/services/user-service/domain/repository"
)

// RoleEntity 角色实体，用于ORM映射
type RoleEntity struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	Description string    `gorm:"type:varchar(255)"`
	Status      string    `gorm:"type:varchar(10);default:'active'"`
	Sort        int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// RolePermissionRelation 角色权限关联实体
type RolePermissionRelation struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	RoleID       string    `gorm:"type:varchar(36);not null;index"`
	PermissionID string    `gorm:"type:varchar(36);not null;index"`
	CreatedAt    time.Time `gorm:"not null"`
}

// TableName 设置表名
func (RoleEntity) TableName() string {
	return "role"
}

func (RolePermissionRelation) TableName() string {
	return "role_permission_relation"
}

// RoleRepositoryImpl 角色仓储MySQL实现
type RoleRepositoryImpl struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository() repository.RoleRepository {
	return &RoleRepositoryImpl{
		db: database.DB,
	}
}

// Save 保存角色
func (r *RoleRepositoryImpl) Save(ctx context.Context, role *model.Role) error {
	// 将领域模型转换为持久化实体
	roleEntity := &RoleEntity{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Status:      string(role.Status),
		Sort:        role.Sort,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	return r.db.Create(roleEntity).Error
}

// FindByID 根据ID查找角色
func (r *RoleRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Role, error) {
	var roleEntity RoleEntity
	err := r.db.Where("id = ?", id).First(&roleEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role not found: %s", id)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&roleEntity)
}

// FindByName 根据名称查找角色
func (r *RoleRepositoryImpl) FindByName(ctx context.Context, name string) (*model.Role, error) {
	var roleEntity RoleEntity
	err := r.db.Where("name = ?", name).First(&roleEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role not found by name: %s", name)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&roleEntity)
}

// Update 更新角色
func (r *RoleRepositoryImpl) Update(ctx context.Context, role *model.Role) error {
	// 将领域模型转换为持久化实体
	roleEntity := &RoleEntity{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Status:      string(role.Status),
		Sort:        role.Sort,
		UpdatedAt:   role.UpdatedAt,
	}

	return r.db.Model(&RoleEntity{}).Where("id = ?", role.ID).Updates(roleEntity).Error
}

// Delete 删除角色
func (r *RoleRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 开始事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err := tx.Where("role_id = ?", id).Delete(&RolePermissionRelation{}).Error; err != nil {
			return err
		}

		// 删除用户角色关联
		if err := tx.Where("role_id = ?", id).Delete(&UserRoleRelation{}).Error; err != nil {
			return err
		}

		// 删除角色
		if err := tx.Where("id = ?", id).Delete(&RoleEntity{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindAll 获取角色列表
func (r *RoleRepositoryImpl) FindAll(ctx context.Context, page, size int) ([]*model.Role, int64, error) {
	var roleEntities []RoleEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&RoleEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Offset(offset).Limit(size).Order("sort desc, created_at desc").Find(&roleEntities).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	roles := make([]*model.Role, 0, len(roleEntities))
	for _, entity := range roleEntities {
		role, err := r.toDomainModel(&entity)
		if err != nil {
			return nil, 0, err
		}
		roles = append(roles, role)
	}

	return roles, total, nil
}

// AddRolePermission 添加角色权限
func (r *RoleRepositoryImpl) AddRolePermission(ctx context.Context, roleID, permissionID string) error {
	relation := RolePermissionRelation{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now(),
	}

	return r.db.Create(&relation).Error
}

// RemoveRolePermission 移除角色权限
func (r *RoleRepositoryImpl) RemoveRolePermission(ctx context.Context, roleID, permissionID string) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&RolePermissionRelation{}).Error
}

// GetRolePermissions 获取角色权限
func (r *RoleRepositoryImpl) GetRolePermissions(ctx context.Context, roleID string) ([]model.Permission, error) {
	// 查询角色权限关联
	var relations []RolePermissionRelation
	if err := r.db.Where("role_id = ?", roleID).Find(&relations).Error; err != nil {
		return nil, err
	}

	// 如果没有权限，返回空数组
	if len(relations) == 0 {
		return []model.Permission{}, nil
	}

	// 提取权限ID
	permIDs := make([]string, 0, len(relations))
	for _, relation := range relations {
		permIDs = append(permIDs, relation.PermissionID)
	}

	// 查询权限
	var permEntities []PermissionEntity
	if err := r.db.Where("id IN ?", permIDs).Find(&permEntities).Error; err != nil {
		return nil, err
	}

	// 转换为领域模型
	permissions := make([]model.Permission, 0, len(permEntities))
	for _, entity := range permEntities {
		permissions = append(permissions, model.Permission{
			ID:          entity.ID,
			Name:        entity.Name,
			Value:       entity.Value,
			Type:        model.PermissionType(entity.Type),
			Status:      model.PermissionStatus(entity.Status),
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		})
	}

	return permissions, nil
}

// FindByUserID 查询用户的所有角色
func (r *RoleRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*model.Role, error) {
	// 查询用户角色关联
	var relations []UserRoleRelation
	if err := r.db.Where("user_id = ?", userID).Find(&relations).Error; err != nil {
		return nil, err
	}

	// 如果没有角色，返回空数组
	if len(relations) == 0 {
		return []*model.Role{}, nil
	}

	// 提取角色ID
	roleIDs := make([]string, 0, len(relations))
	for _, relation := range relations {
		roleIDs = append(roleIDs, relation.RoleID)
	}

	// 查询角色
	var roleEntities []RoleEntity
	if err := r.db.Where("id IN ?", roleIDs).Find(&roleEntities).Error; err != nil {
		return nil, err
	}

	// 转换为领域模型
	roles := make([]*model.Role, 0, len(roleEntities))
	for _, entity := range roleEntities {
		role, err := r.toDomainModel(&entity)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// UpdateStatus 更新角色状态
func (r *RoleRepositoryImpl) UpdateStatus(ctx context.Context, id string, status model.RoleStatus) error {
	return r.db.Model(&RoleEntity{}).Where("id = ?", id).Update("status", status).Error
}

// 辅助方法: 将持久化实体转换为领域模型
func (r *RoleRepositoryImpl) toDomainModel(entity *RoleEntity) (*model.Role, error) {
	return &model.Role{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Status:      model.RoleStatus(entity.Status),
		Sort:        entity.Sort,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}