package model

import (
	"errors"
	"time"
)

// Product 商品聚合根
type Product struct {
	ID                  string
	Name                string
	Description         string
	Price               Money
	Status              ProductStatus
	Inventory           Inventory
	Brand               Brand
	Category            Category
	Attributes          []Attribute
	Promotion           *Promotion
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// 值对象
type Money struct {
	Amount   float64
	Currency string
}

// 值对象
type Inventory struct {
	AvailableQuantity int
	ReservedQuantity  int
	LowStockThreshold int
}

// 值对象
type Attribute struct {
	Name  string
	Value string
}

// 实体
type Brand struct {
	ID   string
	Name string
	Logo string
}

// 实体
type Category struct {
	ID       string
	Name     string
	ParentID string
	Level    int
}

// 实体
type Promotion struct {
	ID           string
	Type         PromotionType
	Discount     float64
	StartTime    time.Time
	EndTime      time.Time
	Requirements map[string]interface{}
}

// 枚举
type ProductStatus string
const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusActive    ProductStatus = "active"
	ProductStatusInactive  ProductStatus = "inactive"
	ProductStatusDeleted   ProductStatus = "deleted"
)

type PromotionType string
const (
	PromotionTypePercentage PromotionType = "percentage"
	PromotionTypeFixedAmount PromotionType = "fixed_amount"
	PromotionTypeBundle     PromotionType = "bundle"
)

// 领域行为/业务规则
func (p *Product) Publish() error {
	if p.Name == "" || p.Price.Amount <= 0 {
		return errors.New("product must have a name and a valid price to be published")
	}

	if p.Inventory.AvailableQuantity <= 0 {
		return errors.New("cannot publish product with zero inventory")
	}

	p.Status = ProductStatusActive
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) UpdateInventory(quantity int) error {
	if quantity < 0 && p.Inventory.AvailableQuantity + quantity < 0 {
		return errors.New("insufficient inventory")
	}

	p.Inventory.AvailableQuantity += quantity
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) ApplyPromotion(promotion *Promotion) error {
	if promotion.EndTime.Before(time.Now()) {
		return errors.New("cannot apply expired promotion")
	}

	p.Promotion = promotion
	p.UpdatedAt = time.Now()
	return nil
}