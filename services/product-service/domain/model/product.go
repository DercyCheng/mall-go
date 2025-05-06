package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ProductStatus 商品状态枚举
type ProductStatus string

const (
	ProductStatusDraft    ProductStatus = "draft"    // 草稿
	ProductStatusActive   ProductStatus = "active"   // 上架
	ProductStatusInactive ProductStatus = "inactive" // 下架
	ProductStatusDeleted  ProductStatus = "deleted"  // 已删除
)

// PromotionType 促销类型枚举
type PromotionType string

const (
	PromotionTypeNone        PromotionType = "none"         // 无促销
	PromotionTypePercentage  PromotionType = "percentage"   // 折扣
	PromotionTypeFixedAmount PromotionType = "fixed_amount" // 满减
	PromotionTypeBundle      PromotionType = "bundle"       // 套餐
)

// Money 金额值对象
type Money struct {
	Amount   float64
	Currency string
}

// Inventory 库存值对象
type Inventory struct {
	AvailableQuantity int
	ReservedQuantity  int
	LowStockThreshold int
}

// Attribute 属性值对象
type Attribute struct {
	Name  string
	Value string
}

// Product 商品聚合根
type Product struct {
	ID                 string
	Name               string
	Description        string
	Price              Money
	OriginalPrice      Money
	Status             ProductStatus
	Inventory          Inventory
	Brand              Brand
	Category           Category
	Attributes         []Attribute
	Promotion          *Promotion
	ProductSn          string
	SubTitle           string
	Keywords           string
	Note               string
	Unit               string
	Weight             float64
	Sort               int
	Sale               int
	PublishStatus      int // 0->下架；1->上架
	NewStatus          int // 0->不是新品；1->新品
	RecommendStatus    int // 0->不推荐；1->推荐
	VerifyStatus       int // 0->未审核；1->审核通过
	PreviewStatus      int // 0->不是预告商品；1->是预告商品
	DeleteStatus       int // 0->未删除；1->已删除
	ServiceIDs         []string
	Pic                string
	AlbumPics          []string
	DetailTitle        string
	DetailDesc         string
	DetailHTML         string
	DetailMobileHTML   string
	PromotionStartTime time.Time
	PromotionEndTime   time.Time
	PromotionPerLimit  int
	PromotionType      int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewProduct 创建新商品的工厂方法
func NewProduct(name string, price float64, desc string, brandID string, categoryID string, stock int) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if price <= 0 {
		return nil, errors.New("product price must be greater than zero")
	}

	if stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}

	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: desc,
		Price: Money{
			Amount:   price,
			Currency: "CNY",
		},
		Status: ProductStatusDraft,
		Inventory: Inventory{
			AvailableQuantity: stock,
			ReservedQuantity:  0,
			LowStockThreshold: 10,
		},
		Brand: Brand{
			ID: brandID,
		},
		Category: Category{
			ID: categoryID,
		},
		Attributes:    []Attribute{},
		PublishStatus: 0,
		NewStatus:     1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// Publish 发布商品
func (p *Product) Publish() error {
	if p.Name == "" || p.Price.Amount <= 0 {
		return errors.New("product must have a name and a valid price to be published")
	}

	if p.Inventory.AvailableQuantity <= 0 {
		return errors.New("cannot publish product with zero inventory")
	}

	p.Status = ProductStatusActive
	p.PublishStatus = 1
	p.UpdatedAt = time.Now()
	return nil
}

// Unpublish 下架商品
func (p *Product) Unpublish() error {
	p.Status = ProductStatusInactive
	p.PublishStatus = 0
	p.UpdatedAt = time.Now()
	return nil
}

// UpdateInventory 更新库存
func (p *Product) UpdateInventory(quantity int) error {
	if quantity < 0 && p.Inventory.AvailableQuantity+quantity < 0 {
		return errors.New("insufficient inventory")
	}

	p.Inventory.AvailableQuantity += quantity
	p.UpdatedAt = time.Now()
	return nil
}

// ReserveStock 预留库存
func (p *Product) ReserveStock(quantity int) error {
	if p.Inventory.AvailableQuantity < quantity {
		return errors.New("insufficient inventory to reserve")
	}

	p.Inventory.AvailableQuantity -= quantity
	p.Inventory.ReservedQuantity += quantity
	p.UpdatedAt = time.Now()
	return nil
}

// ConfirmReservation 确认库存预留（扣减库存）
func (p *Product) ConfirmReservation(quantity int) error {
	if p.Inventory.ReservedQuantity < quantity {
		return errors.New("reservation quantity exceeds reserved amount")
	}

	p.Inventory.ReservedQuantity -= quantity
	p.Sale += quantity
	p.UpdatedAt = time.Now()
	return nil
}

// CancelReservation 取消库存预留
func (p *Product) CancelReservation(quantity int) error {
	if p.Inventory.ReservedQuantity < quantity {
		return errors.New("cancellation quantity exceeds reserved amount")
	}

	p.Inventory.ReservedQuantity -= quantity
	p.Inventory.AvailableQuantity += quantity
	p.UpdatedAt = time.Now()
	return nil
}

// ApplyPromotion 应用促销
func (p *Product) ApplyPromotion(promotion *Promotion) error {
	if promotion.EndTime.Before(time.Now()) {
		return errors.New("cannot apply expired promotion")
	}

	p.Promotion = promotion
	p.UpdatedAt = time.Now()
	return nil
}

// SetRecommended 设置商品为推荐状态
func (p *Product) SetRecommended(status bool) {
	if status {
		p.RecommendStatus = 1
	} else {
		p.RecommendStatus = 0
	}
	p.UpdatedAt = time.Now()
}

// SetNewProduct 设置商品为新品状态
func (p *Product) SetNewProduct(status bool) {
	if status {
		p.NewStatus = 1
	} else {
		p.NewStatus = 0
	}
	p.UpdatedAt = time.Now()
}

// AddAttribute 添加商品属性
func (p *Product) AddAttribute(name, value string) {
	p.Attributes = append(p.Attributes, Attribute{
		Name:  name,
		Value: value,
	})
	p.UpdatedAt = time.Now()
}
