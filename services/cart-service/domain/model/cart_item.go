package model

import (
	"errors"
	"time"
)

// CartItem represents a cart item domain entity
type CartItem struct {
	ID               string    `json:"id"`
	UserID           string    `json:"userId"`
	ProductID        string    `json:"productId"`
	ProductSKUID     string    `json:"productSkuId"`
	ProductName      string    `json:"productName"`
	ProductSubTitle  string    `json:"productSubTitle"`
	ProductPic       string    `json:"productPic"`
	ProductPrice     float64   `json:"productPrice"`
	ProductQuantity  int       `json:"productQuantity"`
	ProductSkuCode   string    `json:"productSkuCode"`
	ProductCategoryID string    `json:"productCategoryId"`
	ProductBrand     string    `json:"productBrand"`
	ProductSn        string    `json:"productSn"`
	ProductAttr      string    `json:"productAttr"`
	PromotionInfo    string    `json:"promotionInfo"`
	PromotionAmount  float64   `json:"promotionAmount"`
	CouponAmount     float64   `json:"couponAmount"`
	RealAmount       float64   `json:"realAmount"`
	CheckStatus      bool      `json:"checkStatus"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// NewCartItem creates a new cart item
func NewCartItem(
	userID, productID, skuID, productName, productPic, skuCode string,
	productPrice float64, quantity int,
) (*CartItem, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if productID == "" {
		return nil, errors.New("product ID is required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	now := time.Now()
	
	return &CartItem{
		ID:              generateID(),
		UserID:          userID,
		ProductID:       productID,
		ProductSKUID:    skuID,
		ProductName:     productName,
		ProductPic:      productPic,
		ProductPrice:    productPrice,
		ProductQuantity: quantity,
		ProductSkuCode:  skuCode,
		RealAmount:      productPrice * float64(quantity),
		CheckStatus:     false,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// IncreaseQuantity increases the quantity of the cart item
func (ci *CartItem) IncreaseQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	ci.ProductQuantity += quantity
	ci.UpdatedAt = time.Now()
	ci.updateRealAmount()
	return nil
}

// DecreaseQuantity decreases the quantity of the cart item
func (ci *CartItem) DecreaseQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if ci.ProductQuantity < quantity {
		return errors.New("cannot decrease quantity below 0")
	}
	ci.ProductQuantity -= quantity
	ci.UpdatedAt = time.Now()
	ci.updateRealAmount()
	return nil
}

// UpdateQuantity updates the quantity of the cart item
func (ci *CartItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	ci.ProductQuantity = quantity
	ci.UpdatedAt = time.Now()
	ci.updateRealAmount()
	return nil
}

// ToggleCheckStatus toggles the check status of the cart item
func (ci *CartItem) ToggleCheckStatus() {
	ci.CheckStatus = !ci.CheckStatus
	ci.UpdatedAt = time.Now()
}

// ApplyPromotion applies a promotion to the cart item
func (ci *CartItem) ApplyPromotion(promotionInfo string, promotionAmount float64) {
	ci.PromotionInfo = promotionInfo
	ci.PromotionAmount = promotionAmount
	ci.UpdatedAt = time.Now()
	ci.updateRealAmount()
}

// ApplyCoupon applies a coupon to the cart item
func (ci *CartItem) ApplyCoupon(couponAmount float64) {
	ci.CouponAmount = couponAmount
	ci.UpdatedAt = time.Now()
	ci.updateRealAmount()
}

// updateRealAmount recalculates the real amount
func (ci *CartItem) updateRealAmount() {
	baseAmount := ci.ProductPrice * float64(ci.ProductQuantity)
	ci.RealAmount = baseAmount - ci.PromotionAmount - ci.CouponAmount
	if ci.RealAmount < 0 {
		ci.RealAmount = 0
	}
}

// generateID generates a unique ID (placeholder implementation)
func generateID() string {
	return "cart_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString generates a random string of the specified length
func randomString(length int) string {
	// In a real implementation, use a proper random string generator
	// This is just a placeholder
	return "12345678"[:length]
}
