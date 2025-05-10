package dto

// CartItemDTO represents a cart item data transfer object
type CartItemDTO struct {
	ID               string  `json:"id"`
	ProductID        string  `json:"productId"`
	ProductSKUID     string  `json:"productSkuId"`
	ProductName      string  `json:"productName"`
	ProductSubTitle  string  `json:"productSubTitle,omitempty"`
	ProductPic       string  `json:"productPic,omitempty"`
	ProductPrice     float64 `json:"productPrice"`
	ProductQuantity  int     `json:"productQuantity"`
	ProductSkuCode   string  `json:"productSkuCode,omitempty"`
	ProductCategoryID string  `json:"productCategoryId,omitempty"`
	ProductBrand     string  `json:"productBrand,omitempty"`
	ProductSn        string  `json:"productSn,omitempty"`
	ProductAttr      string  `json:"productAttr,omitempty"`
	PromotionInfo    string  `json:"promotionInfo,omitempty"`
	PromotionAmount  float64 `json:"promotionAmount,omitempty"`
	CouponAmount     float64 `json:"couponAmount,omitempty"`
	RealAmount       float64 `json:"realAmount"`
	CheckStatus      bool    `json:"checkStatus"`
}

// CartSummaryDTO represents a cart summary data transfer object
type CartSummaryDTO struct {
	TotalAmount      float64 `json:"totalAmount"`
	PromotionAmount  float64 `json:"promotionAmount"`
	CouponAmount     float64 `json:"couponAmount"`
	RealAmount       float64 `json:"realAmount"`
	ItemCount        int     `json:"itemCount"`
}

// AddCartItemRequest represents a request to add an item to the cart
type AddCartItemRequest struct {
	ProductID       string  `json:"productId" binding:"required"`
	ProductSKUID    string  `json:"productSkuId" binding:"required"`
	ProductName     string  `json:"productName" binding:"required"`
	ProductSubTitle string  `json:"productSubTitle"`
	ProductPic      string  `json:"productPic"`
	ProductPrice    float64 `json:"productPrice" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required,min=1"`
	ProductSkuCode  string  `json:"productSkuCode"`
	ProductAttr     string  `json:"productAttr"`
}

// UpdateCartItemRequest represents a request to update a cart item
type UpdateCartItemRequest struct {
	ID        string `json:"id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// DeleteCartItemRequest represents a request to delete a cart item
type DeleteCartItemRequest struct {
	ID string `json:"id" binding:"required"`
}

// UpdateCheckStatusRequest represents a request to update the check status of a cart item
type UpdateCheckStatusRequest struct {
	ID     string `json:"id" binding:"required"`
	Status bool   `json:"status"`
}

// UpdateAllCheckStatusRequest represents a request to update the check status of all cart items
type UpdateAllCheckStatusRequest struct {
	Status bool `json:"status"`
}

// GetCartResponse represents a response containing cart items
type GetCartResponse struct {
	Items   []CartItemDTO `json:"items"`
	Summary CartSummaryDTO `json:"summary"`
}

// GenericResponse represents a generic response
type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CartCountResponse represents a response containing the cart count
type CartCountResponse struct {
	Count int `json:"count"`
}

// CartItemResponse represents a response containing a cart item
type CartItemResponse struct {
	Item CartItemDTO `json:"item"`
}

// ApplyPromotionRequest represents a request to apply a promotion to cart items
type ApplyPromotionRequest struct {
	ItemIDs        []string `json:"itemIds" binding:"required"`
	PromotionInfo  string   `json:"promotionInfo" binding:"required"`
	PromotionAmount float64  `json:"promotionAmount" binding:"required"`
}

// ApplyCouponRequest represents a request to apply a coupon to cart items
type ApplyCouponRequest struct {
	ItemIDs      []string `json:"itemIds" binding:"required"`
	CouponAmount float64  `json:"couponAmount" binding:"required"`
}
