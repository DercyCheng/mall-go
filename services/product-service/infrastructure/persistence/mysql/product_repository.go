package mysql

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"mall-go/services/product-service/domain/model"
	"mall-go/services/product-service/domain/repository"
	"mall-go/services/product-service/infrastructure/util"
)

// ProductEntity 商品数据库实体
type ProductEntity struct {
	ID                 string  `gorm:"primaryKey;column:id;type:varchar(36)"`
	BrandID            string  `gorm:"column:brand_id;type:varchar(36)"`
	BrandName          string  `gorm:"column:brand_name;type:varchar(255)"`
	ProductCategoryID  string  `gorm:"column:product_category_id;type:varchar(36)"`
	CategoryName       string  `gorm:"column:category_name;type:varchar(255)"`
	Name               string  `gorm:"column:name;type:varchar(255)"`
	Pic                string  `gorm:"column:pic;type:varchar(255)"`
	ProductSn          string  `gorm:"column:product_sn;type:varchar(64)"`
	DeleteStatus       int     `gorm:"column:delete_status;type:int;default:0"` // 0->未删除；1->已删除
	PublishStatus      int     `gorm:"column:publish_status;type:int;default:0"` // 0->下架；1->上架
	NewStatus          int     `gorm:"column:new_status;type:int;default:0"` // 0->不是新品；1->新品
	RecommendStatus    int     `gorm:"column:recommend_status;type:int;default:0"` // 0->不推荐；1->推荐
	VerifyStatus       int     `gorm:"column:verify_status;type:int;default:0"` // 0->未审核；1->审核通过
	Sort               int     `gorm:"column:sort;type:int;default:0"`
	Sale               int     `gorm:"column:sale;type:int;default:0"` // 销量
	Price              float64 `gorm:"column:price;type:decimal(10,2)"`
	PromotionPrice     float64 `gorm:"column:promotion_price;type:decimal(10,2)"`
	OriginalPrice      float64 `gorm:"column:original_price;type:decimal(10,2)"`
	Stock              int     `gorm:"column:stock;type:int"`
	LowStock           int     `gorm:"column:low_stock;type:int"`
	Unit               string  `gorm:"column:unit;type:varchar(16)"`
	Weight             float64 `gorm:"column:weight;type:decimal(10,2)"`
	PreviewStatus      int     `gorm:"column:preview_status;type:int;default:0"`
	ServiceIDs         string  `gorm:"column:service_ids;type:varchar(255)"`
	Keywords           string  `gorm:"column:keywords;type:varchar(255)"`
	Note               string  `gorm:"column:note;type:varchar(255)"`
	AlbumPics          string  `gorm:"column:album_pics;type:varchar(1000)"`
	DetailTitle        string  `gorm:"column:detail_title;type:varchar(255)"`
	DetailDesc         string  `gorm:"column:detail_desc;type:text"`
	DetailHTML         string  `gorm:"column:detail_html;type:text"`
	DetailMobileHTML   string  `gorm:"column:detail_mobile_html;type:text"`
	PromotionStartTime int64   `gorm:"column:promotion_start_time;type:bigint"`
	PromotionEndTime   int64   `gorm:"column:promotion_end_time;type:bigint"`
	PromotionPerLimit  int     `gorm:"column:promotion_per_limit;type:int"`
	PromotionType      int     `gorm:"column:promotion_type;type:int;default:0"`
	SubTitle           string  `gorm:"column:sub_title;type:varchar(255)"`
	Description        string  `gorm:"column:description;type:text"`
	CreatedAt          int64   `gorm:"column:created_at;type:bigint"`
	UpdatedAt          int64   `gorm:"column:updated_at;type:bigint"`
}

// TableName 返回表名
func (ProductEntity) TableName() string {
	return "pms_product"
}

// AttributeEntity 商品属性数据库实体
type AttributeEntity struct {
	ID        string `gorm:"primaryKey;column:id;type:varchar(36)"`
	ProductID string `gorm:"column:product_id;type:varchar(36)"`
	Name      string `gorm:"column:name;type:varchar(64)"`
	Value     string `gorm:"column:value;type:varchar(255)"`
	CreatedAt int64  `gorm:"column:created_at;type:bigint"`
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint"`
}

// TableName 返回表名
func (AttributeEntity) TableName() string {
	return "pms_product_attribute_value"
}

// ProductRepositoryImpl 商品仓储实现
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓储实现
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

// Save 保存商品
func (r *productRepository) Save(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// FindByID 根据ID查找商品
func (r *productRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindByCategory 根据分类查找商品
func (r *productRepository) FindByCategory(ctx context.Context, categoryID string, page, size int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64
	
	offset := (page - 1) * size
	
	err := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("product_category_id = ?", categoryID).
		Count(&total).
		Offset(offset).
		Limit(size).
		Find(&products).Error
	
	return products, total, err
}

// FindByBrand 根据品牌查找商品
func (r *productRepository) FindByBrand(ctx context.Context, brandID string, page, size int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64
	
	offset := (page - 1) * size
	
	err := r.db.WithContext(ctx).Model(&model.Product{}).
		Where("brand_id = ?", brandID).
		Count(&total).
		Offset(offset).
		Limit(size).
		Find(&products).Error
	
	return products, total, err
}

// Search 搜索商品
func (r *productRepository) Search(ctx context.Context, query string, filters map[string]interface{}, page, size int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64
	
	offset := (page - 1) * size
	
	db := r.db.WithContext(ctx).Model(&model.Product{})
	
	if query != "" {
		db = db.Where("name LIKE ? OR keywords LIKE ?", "%"+query+"%", "%"+query+"%")
	}
	
	for key, value := range filters {
		db = db.Where(key+" = ?", value)
	}
	
	err := db.Count(&total).
		Offset(offset).
		Limit(size).
		Find(&products).Error
	
	return products, total, err
}

// Update 更新商品
func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete 删除商品
func (r *productRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", id).Error
}

// List 查询商品列表
func (r *productRepository) List(ctx context.Context, pageNum, pageSize int, name, productSn string, publishStatus, verifyStatus int, brandId, productCategoryId string) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64
	
	offset := (pageNum - 1) * pageSize
	
	db := r.db.WithContext(ctx).Model(&model.Product{})
	
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	
	if productSn != "" {
		db = db.Where("product_sn = ?", productSn)
	}
	
	if publishStatus != -1 {
		db = db.Where("publish_status = ?", publishStatus)
	}
	
	if verifyStatus != -1 {
		db = db.Where("verify_status = ?", verifyStatus)
	}
	
	if brandId != "" {
		db = db.Where("brand_id = ?", brandId)
	}
	
	if productCategoryId != "" {
		db = db.Where("product_category_id = ?", productCategoryId)
	}
	
	err := db.Count(&total).
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error
	
	return products, total, err
}

// UpdatePublishStatus 批量更新上架状态
func (r *productRepository) UpdatePublishStatus(ctx context.Context, ids []string, publishStatus int) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id IN ?", ids).
		Update("publish_status", publishStatus).Error
}

// UpdateRecommendStatus 批量更新推荐状态
func (r *productRepository) UpdateRecommendStatus(ctx context.Context, ids []string, recommendStatus int) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id IN ?", ids).
		Update("recommend_status", recommendStatus).Error
}

// UpdateNewStatus 批量更新新品状态
func (r *productRepository) UpdateNewStatus(ctx context.Context, ids []string, newStatus int) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).
		Where("id IN ?", ids).
		Update("new_status", newStatus).Error
}

// DeleteBatch 批量删除
func (r *productRepository) DeleteBatch(ctx context.Context, ids []string) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, "id IN ?", ids).Error
}

// 辅助函数: 领域模型转换为数据库实体
func mapDomainToEntity(product *model.Product) *ProductEntity {
	var entity ProductEntity
	entity.ID = product.ID
	entity.Name = product.Name
	entity.SubTitle = product.SubTitle
	entity.Description = product.Description
	entity.ProductSn = product.ProductSn
	entity.Price = product.Price.Amount
	entity.OriginalPrice = product.OriginalPrice.Amount
	entity.Stock = product.Inventory.AvailableQuantity
	entity.LowStock = product.Inventory.LowStockThreshold
	entity.Unit = product.Unit
	entity.Weight = product.Weight
	entity.Sort = product.Sort
	entity.BrandID = product.Brand.ID
	entity.BrandName = product.Brand.Name
	entity.ProductCategoryID = product.Category.ID
	entity.CategoryName = product.Category.Name
	entity.Pic = product.Pic
	entity.AlbumPics = joinStrings(product.AlbumPics, ",")
	entity.DetailTitle = product.DetailTitle
	entity.DetailDesc = product.DetailDesc
	entity.DetailHTML = product.DetailHTML
	entity.DetailMobileHTML = product.DetailMobileHTML
	entity.Keywords = product.Keywords
	entity.Note = product.Note
	entity.ServiceIDs = joinStrings(product.ServiceIDs, ",")
	entity.DeleteStatus = 0 // 默认未删除
	entity.PublishStatus = product.PublishStatus
	entity.NewStatus = product.NewStatus
	entity.RecommendStatus = product.RecommendStatus
	entity.VerifyStatus = product.VerifyStatus
	entity.PreviewStatus = product.PreviewStatus
	entity.Sale = product.Sale
	entity.PromotionPrice = product.Price.Amount // 默认设置为与价格相同
	entity.PromotionType = product.PromotionType
	entity.CreatedAt = util.TimeToMilliseconds(product.CreatedAt)
	entity.UpdatedAt = util.TimeToMilliseconds(product.UpdatedAt)

	return &entity
}

// 辅助函数: 数据库实体转换为领域模型
func mapEntityToDomain(entity *ProductEntity, attributeEntities []AttributeEntity) *model.Product {
	attributes := make([]model.Attribute, len(attributeEntities))
	for i, attr := range attributeEntities {
		attributes[i] = model.Attribute{
			Name:  attr.Name,
			Value: attr.Value,
		}
	}

	albumPics := splitString(entity.AlbumPics, ",")
	serviceIDs := splitString(entity.ServiceIDs, ",")

	return &model.Product{
		ID:               entity.ID,
		Name:             entity.Name,
		SubTitle:         entity.SubTitle,
		Description:      entity.Description,
		ProductSn:        entity.ProductSn,
		Price: model.Money{
			Amount:   entity.Price,
			Currency: "CNY",
		},
		OriginalPrice: model.Money{
			Amount:   entity.OriginalPrice,
			Currency: "CNY",
		},
		Inventory: model.Inventory{
			AvailableQuantity: entity.Stock,
			ReservedQuantity:  0, // 数据库中没有该字段，默认为 0
			LowStockThreshold: entity.LowStock,
		},
		Unit:             entity.Unit,
		Weight:           entity.Weight,
		Sort:             entity.Sort,
		Brand: model.Brand{
			ID:   entity.BrandID,
			Name: entity.BrandName,
		},
		Category: model.Category{
			ID:   entity.ProductCategoryID,
			Name: entity.CategoryName,
		},
		Pic:                entity.Pic,
		AlbumPics:          albumPics,
		DetailTitle:        entity.DetailTitle,
		DetailDesc:         entity.DetailDesc,
		DetailHTML:         entity.DetailHTML,
		DetailMobileHTML:   entity.DetailMobileHTML,
		Keywords:           entity.Keywords,
		Note:               entity.Note,
		ServiceIDs:         serviceIDs,
		Status:             mapStatusToDomain(entity.PublishStatus),
		PublishStatus:      entity.PublishStatus,
		NewStatus:          entity.NewStatus,
		RecommendStatus:    entity.RecommendStatus,
		VerifyStatus:       entity.VerifyStatus,
		PreviewStatus:      entity.PreviewStatus,
		Sale:               entity.Sale,
		PromotionType:      entity.PromotionType,
		Attributes:         attributes,
		CreatedAt:          util.MillisecondsToTime(entity.CreatedAt),
		UpdatedAt:          util.MillisecondsToTime(entity.UpdatedAt),
	}
}

// 辅助函数: 将上架状态映射为领域模型状态
func mapStatusToDomain(publishStatus int) model.ProductStatus {
	if (publishStatus == 1) {
		return model.ProductStatusActive
	}
	return model.ProductStatusInactive
}

// 辅助函数: 字符串数组转字符串
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	return strings.Join(strs, sep)
}

// 辅助函数: 字符串转字符串数组
func splitString(str string, sep string) []string {
	if str == "" {
		return []string{}
	}
	return strings.Split(str, sep)
}

// 辅助函数: 生成唯一ID
func generateID() string {
	// 实际项目中应该使用更健壮的 ID 生成方法
	// 例如 github.com/google/uuid 的 uuid.New().String()
	return "id-placeholder"
}

// 辅助函数: 获取当前时间的毫秒数
func now() int64 {
	return util.TimeToMilliseconds(time.Now())
}