package mysql

import (
	"context"

	"gorm.io/gorm"

	"mall-go/services/product-service/domain/model"
	"mall-go/services/product-service/domain/repository"
	"mall-go/services/product-service/infrastructure/util"
)

// BrandEntity 品牌数据库实体
type BrandEntity struct {
	ID              string  `gorm:"primaryKey;column:id;type:varchar(36)"`
	Name            string  `gorm:"column:name;type:varchar(64)"`
	FirstLetter     string  `gorm:"column:first_letter;type:varchar(8)"`
	Sort            int     `gorm:"column:sort;type:int"`
	FactoryStatus   int     `gorm:"column:factory_status;type:int"`
	ShowStatus      int     `gorm:"column:show_status;type:int"`
	ProductCount    int     `gorm:"column:product_count;type:int"`
	ProductCommentCount int `gorm:"column:product_comment_count;type:int"`
	Logo            string  `gorm:"column:logo;type:varchar(255)"`
	BigPic          string  `gorm:"column:big_pic;type:varchar(255)"`
	BrandStory      string  `gorm:"column:brand_story;type:text"`
	CreatedAt       int64   `gorm:"column:created_at;type:bigint"`
	UpdatedAt       int64   `gorm:"column:updated_at;type:bigint"`
}

// TableName 返回表名
func (BrandEntity) TableName() string {
	return "pms_brand"
}

// 品牌仓储实现
type brandRepository struct {
	db *gorm.DB
}

// NewBrandRepository 创建品牌仓储实现
func NewBrandRepository(db *gorm.DB) repository.BrandRepository {
	return &brandRepository{db: db}
}

// Save 保存品牌
func (r *brandRepository) Save(ctx context.Context, brand *model.Brand) error {
	entity := mapBrandToEntity(brand)
	return r.db.WithContext(ctx).Create(entity).Error
}

// FindByID 根据ID查找品牌
func (r *brandRepository) FindByID(ctx context.Context, id string) (*model.Brand, error) {
	var entity BrandEntity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return mapEntityToBrand(&entity), nil
}

// FindByName 根据名称查找品牌
func (r *brandRepository) FindByName(ctx context.Context, name string) (*model.Brand, error) {
	var entity BrandEntity
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return mapEntityToBrand(&entity), nil
}

// FindAll 获取所有品牌
func (r *brandRepository) FindAll(ctx context.Context) ([]*model.Brand, error) {
	var entities []BrandEntity
	err := r.db.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	brands := make([]*model.Brand, len(entities))
	for i, entity := range entities {
		brands[i] = mapEntityToBrand(&entity)
	}

	return brands, nil
}

// List 分页查询品牌列表
func (r *brandRepository) List(ctx context.Context, pageNum, pageSize int, name string) ([]*model.Brand, int64, error) {
	var entities []BrandEntity
	var total int64
	
	offset := (pageNum - 1) * pageSize
	
	db := r.db.WithContext(ctx).Model(&BrandEntity{})
	
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	
	err := db.Count(&total).
		Offset(offset).
		Limit(pageSize).
		Find(&entities).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	brands := make([]*model.Brand, len(entities))
	for i, entity := range entities {
		brands[i] = mapEntityToBrand(&entity)
	}
	
	return brands, total, nil
}

// Update 更新品牌
func (r *brandRepository) Update(ctx context.Context, brand *model.Brand) error {
	entity := mapBrandToEntity(brand)
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 删除品牌
func (r *brandRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&BrandEntity{}, "id = ?", id).Error
}

// UpdateShowStatus 批量更新显示状态
func (r *brandRepository) UpdateShowStatus(ctx context.Context, ids []string, showStatus int) error {
	return r.db.WithContext(ctx).Model(&BrandEntity{}).
		Where("id IN ?", ids).
		Update("show_status", showStatus).Error
}

// UpdateFactoryStatus 批量更新厂商状态
func (r *brandRepository) UpdateFactoryStatus(ctx context.Context, ids []string, factoryStatus int) error {
	return r.db.WithContext(ctx).Model(&BrandEntity{}).
		Where("id IN ?", ids).
		Update("factory_status", factoryStatus).Error
}

// 辅助函数: 领域模型转换为数据库实体
func mapBrandToEntity(brand *model.Brand) *BrandEntity {
	if brand == nil {
		return nil
	}

	return &BrandEntity{
		ID:                brand.ID,
		Name:              brand.Name,
		FirstLetter:       brand.FirstLetter,
		Sort:              brand.Sort,
		FactoryStatus:     brand.FactoryStatus,
		ShowStatus:        brand.ShowStatus,
		ProductCount:      brand.ProductCount,
		ProductCommentCount: brand.ProductCommentCount,
		Logo:              brand.Logo,
		BigPic:            brand.BigPic,
		BrandStory:        brand.BrandStory,
		CreatedAt:         util.TimeToMilliseconds(brand.CreatedAt),
		UpdatedAt:         util.TimeToMilliseconds(brand.UpdatedAt),
	}
}

// 辅助函数: 数据库实体转换为领域模型
func mapEntityToBrand(entity *BrandEntity) *model.Brand {
	if entity == nil {
		return nil
	}

	return &model.Brand{
		ID:                entity.ID,
		Name:              entity.Name,
		FirstLetter:       entity.FirstLetter,
		Sort:              entity.Sort,
		FactoryStatus:     entity.FactoryStatus,
		ShowStatus:        entity.ShowStatus,
		ProductCount:      entity.ProductCount,
		ProductCommentCount: entity.ProductCommentCount,
		Logo:              entity.Logo,
		BigPic:            entity.BigPic,
		BrandStory:        entity.BrandStory,
		CreatedAt:         util.MillisecondsToTime(entity.CreatedAt),
		UpdatedAt:         util.MillisecondsToTime(entity.UpdatedAt),
	}
}