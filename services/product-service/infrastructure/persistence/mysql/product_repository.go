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

// ProductEntity 商品实体，用于ORM映射
type ProductEntity struct {
	ID                  string    `gorm:"primaryKey;type:varchar(36)"`
	Name                string    `gorm:"type:varchar(200);not null"`
	Description         string    `gorm:"type:text"`
	PriceAmount         float64   `gorm:"column:price_amount;not null"`
	PriceCurrency       string    `gorm:"column:price_currency;type:varchar(3);not null;default:'CNY'"`
	Status              string    `gorm:"type:varchar(20);not null;default:'draft'"`
	AvailableQuantity   int       `gorm:"not null;default:0"`
	ReservedQuantity    int       `gorm:"not null;default:0"`
	LowStockThreshold   int       `gorm:"default:0"`
	BrandID             string    `gorm:"type:varchar(36);not null"`
	BrandName           string    `gorm:"type:varchar(100)"`
	BrandLogo           string    `gorm:"type:varchar(255)"`
	CategoryID          string    `gorm:"type:varchar(36);not null"`
	CategoryName        string    `gorm:"type:varchar(100)"`
	CategoryParentID    string    `gorm:"type:varchar(36)"`
	CategoryLevel       int       `gorm:"default:0"`
	CreatedAt           time.Time `gorm:"not null"`
	UpdatedAt           time.Time `gorm:"not null"`
}

// ProductAttributeEntity 商品属性实体
type ProductAttributeEntity struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ProductID string `gorm:"type:varchar(36);not null;index"`
	Name      string `gorm:"type:varchar(100);not null"`
	Value     string `gorm:"type:varchar(255);not null"`
}

// ProductPromotionEntity 商品促销实体
type ProductPromotionEntity struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)"`
	ProductID    string    `gorm:"type:varchar(36);not null;uniqueIndex"`
	Type         string    `gorm:"type:varchar(20);not null"`
	Discount     float64   `gorm:"not null"`
	StartTime    time.Time `gorm:"not null"`
	EndTime      time.Time `gorm:"not null"`
	Requirements string    `gorm:"type:text"` // JSON 格式存储
}

// TableName 设置表名
func (ProductEntity) TableName() string {
	return "product"
}

func (ProductAttributeEntity) TableName() string {
	return "product_attribute"
}

func (ProductPromotionEntity) TableName() string {
	return "product_promotion"
}

// ProductRepositoryImpl 商品仓储MySQL实现
type ProductRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓储实例
func NewProductRepository() repository.ProductRepository {
	return &ProductRepositoryImpl{
		db: database.DB,
	}
}

// Save 保存商品
func (r *ProductRepositoryImpl) Save(ctx context.Context, product *model.Product) error {
	// 1. 将领域模型转换为持久化实体
	productEntity := mapToProductEntity(product)

	// 2. 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 3. 保存商品主体
		if err := tx.Create(&productEntity).Error; err != nil {
			return err
		}

		// 4. 保存商品属性
		if len(product.Attributes) > 0 {
			attributeEntities := make([]ProductAttributeEntity, 0, len(product.Attributes))
			for _, attr := range product.Attributes {
				attributeEntities = append(attributeEntities, ProductAttributeEntity{
					ProductID: product.ID,
					Name:      attr.Name,
					Value:     attr.Value,
				})
			}
			if err := tx.Create(&attributeEntities).Error; err != nil {
				return err
			}
		}

		// 5. 保存商品促销信息
		if product.Promotion != nil {
			promotionEntity := ProductPromotionEntity{
				ID:        product.Promotion.ID,
				ProductID: product.ID,
				Type:      string(product.Promotion.Type),
				Discount:  product.Promotion.Discount,
				StartTime: product.Promotion.StartTime,
				EndTime:   product.Promotion.EndTime,
				// 在实际实现中需要将map转为JSON字符串
				Requirements: "{}", // 简化处理，实际应序列化为JSON
			}
			if err := tx.Create(&promotionEntity).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindByID 根据ID查找商品
func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Product, error) {
	var productEntity ProductEntity
	if err := r.db.Where("id = ?", id).First(&productEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product not found: %s", id)
		}
		return nil, err
	}

	// 获取商品属性
	var attributeEntities []ProductAttributeEntity
	if err := r.db.Where("product_id = ?", id).Find(&attributeEntities).Error; err != nil {
		return nil, err
	}

	// 获取商品促销信息
	var promotionEntity ProductPromotionEntity
	hasPromotion := true
	if err := r.db.Where("product_id = ?", id).First(&promotionEntity).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		hasPromotion = false
	}

	// 将持久化实体转换为领域模型
	return mapToDomainModel(&productEntity, attributeEntities, hasPromotion, &promotionEntity), nil
}

// Update 更新商品
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *model.Product) error {
	productEntity := mapToProductEntity(product)

	return r.db.Transaction(func(tx *gorm.DB) error {
		// 更新商品主体
		if err := tx.Model(&ProductEntity{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
			"name":                productEntity.Name,
			"description":         productEntity.Description,
			"price_amount":        productEntity.PriceAmount,
			"price_currency":      productEntity.PriceCurrency,
			"status":              productEntity.Status,
			"available_quantity":  productEntity.AvailableQuantity,
			"reserved_quantity":   productEntity.ReservedQuantity,
			"low_stock_threshold": productEntity.LowStockThreshold,
			"brand_id":            productEntity.BrandID,
			"brand_name":          productEntity.BrandName,
			"brand_logo":          productEntity.BrandLogo,
			"category_id":         productEntity.CategoryID,
			"category_name":       productEntity.CategoryName,
			"category_parent_id":  productEntity.CategoryParentID,
			"category_level":      productEntity.CategoryLevel,
			"updated_at":          productEntity.UpdatedAt,
		}).Error; err != nil {
			return err
		}

		// 更新商品属性：先删除旧的，再添加新的
		if err := tx.Where("product_id = ?", product.ID).Delete(&ProductAttributeEntity{}).Error; err != nil {
			return err
		}

		if len(product.Attributes) > 0 {
			attributeEntities := make([]ProductAttributeEntity, 0, len(product.Attributes))
			for _, attr := range product.Attributes {
				attributeEntities = append(attributeEntities, ProductAttributeEntity{
					ProductID: product.ID,
					Name:      attr.Name,
					Value:     attr.Value,
				})
			}
			if err := tx.Create(&attributeEntities).Error; err != nil {
				return err
			}
		}

		// 更新商品促销信息
		if product.Promotion != nil {
			// 先检查是否存在
			var count int64
			tx.Model(&ProductPromotionEntity{}).Where("product_id = ?", product.ID).Count(&count)

			promotionEntity := ProductPromotionEntity{
				ID:        product.Promotion.ID,
				ProductID: product.ID,
				Type:      string(product.Promotion.Type),
				Discount:  product.Promotion.Discount,
				StartTime: product.Promotion.StartTime,
				EndTime:   product.Promotion.EndTime,
				// 在实际实现中需要将map转为JSON字符串
				Requirements: "{}", // 简化处理，实际应序列化为JSON
			}

			if count > 0 {
				// 更新已有促销
				if err := tx.Model(&ProductPromotionEntity{}).Where("product_id = ?", product.ID).Updates(map[string]interface{}{
					"type":         promotionEntity.Type,
					"discount":     promotionEntity.Discount,
					"start_time":   promotionEntity.StartTime,
					"end_time":     promotionEntity.EndTime,
					"requirements": promotionEntity.Requirements,
				}).Error; err != nil {
					return err
				}
			} else {
				// 创建新促销
				if err := tx.Create(&promotionEntity).Error; err != nil {
					return err
				}
			}
		} else {
			// 删除促销
			if err := tx.Where("product_id = ?", product.ID).Delete(&ProductPromotionEntity{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete 删除商品
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除商品属性
		if err := tx.Where("product_id = ?", id).Delete(&ProductAttributeEntity{}).Error; err != nil {
			return err
		}

		// 删除商品促销
		if err := tx.Where("product_id = ?", id).Delete(&ProductPromotionEntity{}).Error; err != nil {
			return err
		}

		// 删除商品主体
		if err := tx.Where("id = ?", id).Delete(&ProductEntity{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateStatus 更新商品状态
func (r *ProductRepositoryImpl) UpdateStatus(ctx context.Context, id string, status model.ProductStatus) error {
	return r.db.Model(&ProductEntity{}).Where("id = ?", id).Update("status", status).Error
}

// FindByCategory 根据分类查找商品
func (r *ProductRepositoryImpl) FindByCategory(ctx context.Context, categoryID string, page, size int) ([]*model.Product, int64, error) {
	var productEntities []ProductEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&ProductEntity{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Where("category_id = ?", categoryID).Offset(offset).Limit(size).Find(&productEntities).Error; err != nil {
		return nil, 0, err
	}

	// 查询对应的属性和促销信息
	products := make([]*model.Product, 0, len(productEntities))
	for _, entity := range productEntities {
		product, err := r.FindByID(ctx, entity.ID)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

// FindByBrand 根据品牌查找商品
func (r *ProductRepositoryImpl) FindByBrand(ctx context.Context, brandID string, page, size int) ([]*model.Product, int64, error) {
	var productEntities []ProductEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&ProductEntity{}).Where("brand_id = ?", brandID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Where("brand_id = ?", brandID).Offset(offset).Limit(size).Find(&productEntities).Error; err != nil {
		return nil, 0, err
	}

	// 查询对应的属性和促销信息
	products := make([]*model.Product, 0, len(productEntities))
	for _, entity := range productEntities {
		product, err := r.FindByID(ctx, entity.ID)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

// Search 搜索商品
func (r *ProductRepositoryImpl) Search(ctx context.Context, query string, filters map[string]interface{}, page, size int) ([]*model.Product, int64, error) {
	var productEntities []ProductEntity
	var total int64

	db := r.db.Model(&ProductEntity{})

	// 添加关键词搜索条件
	if query != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 添加过滤条件
	for key, value := range filters {
		db = db.Where(key+" = ?", value)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := db.Offset(offset).Limit(size).Find(&productEntities).Error; err != nil {
		return nil, 0, err
	}

	// 查询对应的属性和促销信息
	products := make([]*model.Product, 0, len(productEntities))
	for _, entity := range productEntities {
		product, err := r.FindByID(ctx, entity.ID)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

// FindAll 获取所有商品
func (r *ProductRepositoryImpl) FindAll(ctx context.Context, page, size int) ([]*model.Product, int64, error) {
	var productEntities []ProductEntity
	var total int64

	// 查询总数
	if err := r.db.Model(&ProductEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.db.Offset(offset).Limit(size).Find(&productEntities).Error; err != nil {
		return nil, 0, err
	}

	// 查询对应的属性和促销信息
	products := make([]*model.Product, 0, len(productEntities))
	for _, entity := range productEntities {
		product, err := r.FindByID(ctx, entity.ID)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, total, nil
}

// 辅助函数: 将领域模型映射为持久化实体
func mapToProductEntity(product *model.Product) ProductEntity {
	return ProductEntity{
		ID:                product.ID,
		Name:              product.Name,
		Description:       product.Description,
		PriceAmount:       product.Price.Amount,
		PriceCurrency:     product.Price.Currency,
		Status:            string(product.Status),
		AvailableQuantity: product.Inventory.AvailableQuantity,
		ReservedQuantity:  product.Inventory.ReservedQuantity,
		LowStockThreshold: product.Inventory.LowStockThreshold,
		BrandID:           product.Brand.ID,
		BrandName:         product.Brand.Name,
		BrandLogo:         product.Brand.Logo,
		CategoryID:        product.Category.ID,
		CategoryName:      product.Category.Name,
		CategoryParentID:  product.Category.ParentID,
		CategoryLevel:     product.Category.Level,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}
}

// 辅助函数: 将持久化实体映射为领域模型
func mapToDomainModel(entity *ProductEntity, attributeEntities []ProductAttributeEntity, hasPromotion bool, promotionEntity *ProductPromotionEntity) *model.Product {
	// 映射属性
	attributes := make([]model.Attribute, 0, len(attributeEntities))
	for _, attrEntity := range attributeEntities {
		attributes = append(attributes, model.Attribute{
			Name:  attrEntity.Name,
			Value: attrEntity.Value,
		})
	}

	// 创建商品领域模型
	product := &model.Product{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Price: model.Money{
			Amount:   entity.PriceAmount,
			Currency: entity.PriceCurrency,
		},
		Status: model.ProductStatus(entity.Status),
		Inventory: model.Inventory{
			AvailableQuantity: entity.AvailableQuantity,
			ReservedQuantity:  entity.ReservedQuantity,
			LowStockThreshold: entity.LowStockThreshold,
		},
		Brand: model.Brand{
			ID:   entity.BrandID,
			Name: entity.BrandName,
			Logo: entity.BrandLogo,
		},
		Category: model.Category{
			ID:       entity.CategoryID,
			Name:     entity.CategoryName,
			ParentID: entity.CategoryParentID,
			Level:    entity.CategoryLevel,
		},
		Attributes: attributes,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}

	// 添加促销信息
	if hasPromotion && promotionEntity != nil {
		product.Promotion = &model.Promotion{
			ID:           promotionEntity.ID,
			Type:         model.PromotionType(promotionEntity.Type),
			Discount:     promotionEntity.Discount,
			StartTime:    promotionEntity.StartTime,
			EndTime:      promotionEntity.EndTime,
			Requirements: make(map[string]interface{}), // 简化处理
		}
	}

	return product
}