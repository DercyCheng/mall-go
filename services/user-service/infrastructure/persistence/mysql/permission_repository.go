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

// PermissionEntity 权限实体，用于ORM映射
type PermissionEntity struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	Value       string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Type        string    `gorm:"type:varchar(20);default:'api'"`
	Status      string    `gorm:"type:varchar(10);default:'active'"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// TableName 设置表名
func (PermissionEntity) TableName() string {
	return "permission"
}

// PermissionRepositoryImpl 权限仓储MySQL实现
type PermissionRepositoryImpl struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓储实例
func NewPermissionRepository() repository.PermissionRepository {
	return &PermissionRepositoryImpl{
		db: database.DB,
	}
}

// Save 保存权限
func (r *PermissionRepositoryImpl) Save(ctx context.Context, permission *model.Permission) error {
	// 将领域模型转换为持久化实体
	permEntity := &PermissionEntity{
		ID:          permission.ID,
		Name:        permission.Name,
		Value:       permission.Value,
		Type:        string(permission.Type),
		Status:      string(permission.Status),
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}

	return r.db.Create(permEntity).Error
}

// FindByID 根据ID查找权限
func (r *PermissionRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Permission, error) {
	var permEntity PermissionEntity
	err := r.db.Where("id = ?", id).First(&permEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("permission not found: %s", id)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&permEntity)
}

// FindByValue 根据Value查找权限
func (r *PermissionRepositoryImpl) FindByValue(ctx context.Context, value string) (*model.Permission, error) {
	var permEntity PermissionEntity
	err := r.db.Where("value = ?", value).First(&permEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("permission not found by value: %s", value)
		}
		return nil, err
	}

	// 转换为领域模型
	return r.toDomainModel(&permEntity)
}

// Update 更新权限
func (r *PermissionRepositoryImpl) Update(ctx context.Context, permission *model.Permission) error {
	// 将领域模型转换为持久化实体
	permEntity := &PermissionEntity{
		ID:        permission.ID,
		Name:      permission.Name,
		Value:     permission.Value,
		Type:      string(permission.Type),
		Status:    string(permission.Status),
		UpdatedAt: permission.UpdatedAt,
	}

	return r.db.Model(&PermissionEntity{}).Where("id = ?", permission.ID).Updates(permEntity).Error
}

// Delete 删除权限
func (r *PermissionRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 开始事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除角色权限关联
		if err := tx.Where("permission_id = ?", id).Delete(&RolePermissionRelation{}).Error; err != nil {
			return err
		}

		// 删除权限
		if err := tx.Where("id = ?", id).Delete(&PermissionEntity{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// FindAll 获取权限列表
func (r *PermissionRepositoryImpl) FindAll(ctx context.Context, page, size int) ([]*model.Permission, int64, error) {
	var permEntities []PermissionEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&PermissionEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Offset(offset).Limit(size).Order("created_at desc").Find(&permEntities).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	permissions := make([]*model.Permission, 0, len(permEntities))
	for _, entity := range permEntities {
		perm, err := r.toDomainModel(&entity)
		if err != nil {
			return nil, 0, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, total, nil
}

// FindByType 根据类型查找权限
func (r *PermissionRepositoryImpl) FindByType(ctx context.Context, permType model.PermissionType, page, size int) ([]*model.Permission, int64, error) {
	var permEntities []PermissionEntity
	var total int64

	query := r.db.Model(&PermissionEntity{}).Where("type = ?", permType)

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at desc").Find(&permEntities).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	permissions := make([]*model.Permission, 0, len(permEntities))
	for _, entity := range permEntities {
		perm, err := r.toDomainModel(&entity)
		if err != nil {
			return nil, 0, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, total, nil
}

// FindByRoleID 查询角色的所有权限
func (r *PermissionRepositoryImpl) FindByRoleID(ctx context.Context, roleID string) ([]*model.Permission, error) {
	// 查询角色权限关联
	var relations []RolePermissionRelation
	if err := r.db.Where("role_id = ?", roleID).Find(&relations).Error; err != nil {
		return nil, err
	}

	// 如果没有权限，返回空数组
	if len(relations) == 0 {
		return []*model.Permission{}, nil
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
	permissions := make([]*model.Permission, 0, len(permEntities))
	for _, entity := range permEntities {
		perm, err := r.toDomainModel(&entity)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// UpdateStatus 更新权限状态
func (r *PermissionRepositoryImpl) UpdateStatus(ctx context.Context, id string, status model.PermissionStatus) error {
	return r.db.Model(&PermissionEntity{}).Where("id = ?", id).Update("status", status).Error
}

// 辅助方法: 将持久化实体转换为领域模型
func (r *PermissionRepositoryImpl) toDomainModel(entity *PermissionEntity) (*model.Permission, error) {
	return &model.Permission{
		ID:        entity.ID,
		Name:      entity.Name,
		Value:     entity.Value,
		Type:      model.PermissionType(entity.Type),
		Status:    model.PermissionStatus(entity.Status),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}