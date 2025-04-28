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

// BrandEntity 品牌实体，用于ORM映射
type BrandEntity struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Name      string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Logo      string    `gorm:"type:varchar(255)"`
	Sort      int       `gorm:"default:0"`
	Status    int       `gorm:"default:1"` // 1: 启用, 0: 禁用
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TableName 设置表名
func (BrandEntity) TableName() string {
	return "brand"
}

// BrandRepositoryImpl 品牌仓储MySQL实现
type BrandRepositoryImpl struct {
	db *gorm.DB
}

// NewBrandRepository 创建品牌仓储实例
func NewBrandRepository() repository.BrandRepository {
	return &BrandRepositoryImpl{
		db: database.DB,
	}
}

// FindByID 根据ID查找品牌
func (r *BrandRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Brand, error) {
	var brandEntity BrandEntity
	if err := r.db.Where("id = ?", id).First(&brandEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("brand not found: %s", id)
		}
		return nil, err
	}

	// 转换为领域模型
	return &model.Brand{
		ID:   brandEntity.ID,
		Name: brandEntity.Name,
		Logo: brandEntity.Logo,
	}, nil
}

// FindAll 获取品牌列表
func (r *BrandRepositoryImpl) FindAll(ctx context.Context, page, size int) ([]*model.Brand, int64, error) {
	var brandEntities []BrandEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&BrandEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Order("sort desc, id asc").Offset(offset).Limit(size).Find(&brandEntities).Error; err != nil {
		return nil, 0, err
	}

	// 转换为领域模型
	brands := make([]*model.Brand, 0, len(brandEntities))
	for _, entity := range brandEntities {
		brands = append(brands, &model.Brand{
			ID:   entity.ID,
			Name: entity.Name,
			Logo: entity.Logo,
		})
	}

	return brands, total, nil
}

// Save 保存品牌
func (r *BrandRepositoryImpl) Save(ctx context.Context, brand *model.Brand) error {
	brandEntity := BrandEntity{
		ID:        brand.ID,
		Name:      brand.Name,
		Logo:      brand.Logo,
		Sort:      0,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.db.Create(&brandEntity).Error
}

// Update 更新品牌
func (r *BrandRepositoryImpl) Update(ctx context.Context, brand *model.Brand) error {
	return r.db.Model(&BrandEntity{}).Where("id = ?", brand.ID).Updates(map[string]interface{}{
		"name":       brand.Name,
		"logo":       brand.Logo,
		"updated_at": time.Now(),
	}).Error
}

// Delete 删除品牌
func (r *BrandRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.Where("id = ?", id).Delete(&BrandEntity{}).Error
}