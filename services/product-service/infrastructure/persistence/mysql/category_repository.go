package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mall-go/pkg/database"
	"mall-go/services/product-service/domain/model"
	"mall-go/services/product-service/domain/repository"
)

// CategoryEntity 分类实体，用于ORM映射
type CategoryEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Name      string    `gorm:"type:varchar(100);not null"`
	ParentID  string    `gorm:"type:varchar(36);index"`
	Level     int       `gorm:"default:0"`
	Sort      int       `gorm:"default:0"`
	Icon      string    `gorm:"type:varchar(255)"`
	Keywords  string    `gorm:"type:varchar(255)"`
	Status    int       `gorm:"default:1"` // 1: 启用, 0: 禁用
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TableName 设置表名
func (CategoryEntity) TableName() string {
	return "category"
}

// CategoryRepositoryImpl 分类仓储MySQL实现
type CategoryRepositoryImpl struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储实例
func NewCategoryRepository() repository.CategoryRepository {
	return &CategoryRepositoryImpl{
		db: database.DB,
	}
}

// FindByID 根据ID查找分类
func (r *CategoryRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Category, error) {
	var categoryEntity CategoryEntity
	if err := r.db.Where("id = ?", id).First(&categoryEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category not found: %s", id)
		}
		return nil, err
	}

	// 转换为领域模型
	return &model.Category{
		ID:       categoryEntity.ID,
		Name:     categoryEntity.Name,
		ParentID: categoryEntity.ParentID,
		Level:    categoryEntity.Level,
	}, nil
}

// FindByParentID 根据父ID查找分类
func (r *CategoryRepositoryImpl) FindByParentID(ctx context.Context, parentID string) ([]*model.Category, error) {
	var categoryEntities []CategoryEntity
	if err := r.db.Where("parent_id = ?", parentID).Order("sort desc, id asc").Find(&categoryEntities).Error; err != nil {
		return nil, err
	}

	// 转换为领域模型
	categories := make([]*model.Category, 0, len(categoryEntities))
	for _, entity := range categoryEntities {
		categories = append(categories, &model.Category{
			ID:       entity.ID,
			Name:     entity.Name,
			ParentID: entity.ParentID,
			Level:    entity.Level,
		})
	}

	return categories, nil
}

// FindAll 获取分类列表
func (r *CategoryRepositoryImpl) FindAll(ctx context.Context, page, size int) ([]*model.Category, int64, error) {
	var categoryEntities []CategoryEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&CategoryEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Order("sort desc, id asc").Offset(offset).Limit(size).Find(&categoryEntities).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	categories := make([]*model.Category, 0, len(categoryEntities))
	for _, entity := range categoryEntities {
		categories = append(categories, &model.Category{
			ID:       entity.ID,
			Name:     entity.Name,
			ParentID: entity.ParentID,
			Level:    entity.Level,
		})
	}

	return categories, total, nil
}

// Save 保存分类
func (r *CategoryRepositoryImpl) Save(ctx context.Context, category *model.Category) error {
	categoryEntity := CategoryEntity{
		ID:        category.ID,
		Name:      category.Name,
		ParentID:  category.ParentID,
		Level:     category.Level,
		Sort:      0,
		Icon:      "",
		Keywords:  "",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.db.Create(&categoryEntity).Error
}

// Update 更新分类
func (r *CategoryRepositoryImpl) Update(ctx context.Context, category *model.Category) error {
	return r.db.Model(&CategoryEntity{}).Where("id = ?", category.ID).Updates(map[string]interface{}{
		"name":       category.Name,
		"parent_id":  category.ParentID,
		"level":      category.Level,
		"updated_at": time.Now(),
	}).Error
}

// Delete 删除分类
func (r *CategoryRepositoryImpl) Delete(ctx context.Context, id string) error {
	// 检查是否有子分类
	var count int64
	if err := r.db.Model(&CategoryEntity{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("cannot delete category with children: %s", id)
	}

	return r.db.Where("id = ?", id).Delete(&CategoryEntity{}).Error
}