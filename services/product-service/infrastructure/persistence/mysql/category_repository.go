package mysql

import (
	"context"

	"gorm.io/gorm"

	"mall-go/services/product-service/domain/model"
	"mall-go/services/product-service/domain/repository"
	"mall-go/services/product-service/infrastructure/util"
)

// CategoryEntity 分类数据库实体
type CategoryEntity struct {
	ID           string `gorm:"primaryKey;column:id;type:varchar(36)"`
	ParentID     string `gorm:"column:parent_id;type:varchar(36)"`
	Name         string `gorm:"column:name;type:varchar(64)"`
	Level        int    `gorm:"column:level;type:int"`
	ProductCount int    `gorm:"column:product_count;type:int"`
	ProductUnit  string `gorm:"column:product_unit;type:varchar(64)"`
	NavStatus    int    `gorm:"column:nav_status;type:int"`
	ShowStatus   int    `gorm:"column:show_status;type:int"`
	Sort         int    `gorm:"column:sort;type:int"`
	Icon         string `gorm:"column:icon;type:varchar(255)"`
	Keywords     string `gorm:"column:keywords;type:varchar(255)"`
	Description  string `gorm:"column:description;text"`
	CreatedAt    int64  `gorm:"column:created_at;type:bigint"`
	UpdatedAt    int64  `gorm:"column:updated_at;type:bigint"`
}

// TableName 返回表名
func (CategoryEntity) TableName() string {
	return "pms_product_category"
}

// 分类仓储实现
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储实现
func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

// Save 保存分类
func (r *categoryRepository) Save(ctx context.Context, category *model.Category) error {
	entity := mapCategoryToEntity(category)
	return r.db.WithContext(ctx).Create(entity).Error
}

// FindByID 根据ID查找分类
func (r *categoryRepository) FindByID(ctx context.Context, id string) (*model.Category, error) {
	var entity CategoryEntity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return mapEntityToCategory(&entity), nil
}

// FindByParentID 根据父ID查找子分类
func (r *categoryRepository) FindByParentID(ctx context.Context, parentID string) ([]*model.Category, error) {
	var entities []CategoryEntity
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	categories := make([]*model.Category, len(entities))
	for i, entity := range entities {
		categories[i] = mapEntityToCategory(&entity)
	}

	return categories, nil
}

// FindAll 获取所有分类
func (r *categoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	var entities []CategoryEntity
	err := r.db.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	categories := make([]*model.Category, len(entities))
	for i, entity := range entities {
		categories[i] = mapEntityToCategory(&entity)
	}

	return categories, nil
}

// FindTree 查询分类树形结构
func (r *categoryRepository) FindTree(ctx context.Context) ([]*model.Category, error) {
	// 先获取所有分类
	categories, err := r.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 构建分类树
	categoryMap := make(map[string]*model.Category)
	var roots []*model.Category

	// 第一次遍历，存储所有分类到 map 中
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	// 第二次遍历，构建父子关系
	for _, category := range categories {
		if category.ParentID == "" || category.ParentID == "0" {
			// 根分类
			roots = append(roots, category)
		} else {
			// 子分类
			if parent, ok := categoryMap[category.ParentID]; ok {
				parent.Children = append(parent.Children, category)
			}
		}
	}

	return roots, nil
}

// Update 更新分类
func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	entity := mapCategoryToEntity(category)
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 删除分类
func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&CategoryEntity{}, "id = ?", id).Error
}

// UpdateNavStatus 批量更新导航栏显示状态
func (r *categoryRepository) UpdateNavStatus(ctx context.Context, ids []string, navStatus int) error {
	return r.db.WithContext(ctx).Model(&CategoryEntity{}).
		Where("id IN ?", ids).
		Update("nav_status", navStatus).Error
}

// UpdateShowStatus 批量更新显示状态
func (r *categoryRepository) UpdateShowStatus(ctx context.Context, ids []string, showStatus int) error {
	return r.db.WithContext(ctx).Model(&CategoryEntity{}).
		Where("id IN ?", ids).
		Update("show_status", showStatus).Error
}

// FindWithChildren 查询分类及其子分类
func (r *categoryRepository) FindWithChildren(ctx context.Context, id string) (*model.Category, []*model.Category, error) {
	// 先查询当前分类
	category, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	// 查询所有子分类
	childCategories, err := r.FindByParentID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return category, childCategories, nil
}

// 辅助函数: 领域模型转换为数据库实体
func mapCategoryToEntity(category *model.Category) *CategoryEntity {
	if category == nil {
		return nil
	}

	return &CategoryEntity{
		ID:           category.ID,
		ParentID:     category.ParentID,
		Name:         category.Name,
		Level:        category.Level,
		ProductCount: category.ProductCount,
		ProductUnit:  category.ProductUnit,
		NavStatus:    category.NavStatus,
		ShowStatus:   category.ShowStatus,
		Sort:         category.Sort,
		Icon:         category.Icon,
		Keywords:     category.Keywords,
		Description:  category.Description,
		CreatedAt:    util.TimeToMilliseconds(category.CreatedAt),
		UpdatedAt:    util.TimeToMilliseconds(category.UpdatedAt),
	}
}

// 辅助函数: 数据库实体转换为领域模型
func mapEntityToCategory(entity *CategoryEntity) *model.Category {
	if entity == nil {
		return nil
	}

	return &model.Category{
		ID:           entity.ID,
		ParentID:     entity.ParentID,
		Name:         entity.Name,
		Level:        entity.Level,
		ProductCount: entity.ProductCount,
		ProductUnit:  entity.ProductUnit,
		NavStatus:    entity.NavStatus,
		ShowStatus:   entity.ShowStatus,
		Sort:         entity.Sort,
		Icon:         entity.Icon,
		Keywords:     entity.Keywords,
		Description:  entity.Description,
		Children:     []*model.Category{},
		CreatedAt:    util.MillisecondsToTime(entity.CreatedAt),
		UpdatedAt:    util.MillisecondsToTime(entity.UpdatedAt),
	}
}