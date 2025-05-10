package service

import (
	"context"
	
	"mall-go/services/cart-service/application/dto"
)

// CartService defines the interface for cart application service
type CartService interface {
	// AddItem adds a new item to the cart
	AddItem(ctx context.Context, userID string, req dto.AddCartItemRequest) (*dto.CartItemDTO, error)
	
	// GetCart returns the user's cart
	GetCart(ctx context.Context, userID string) (*dto.GetCartResponse, error)
	
	// UpdateItemQuantity updates the quantity of a cart item
	UpdateItemQuantity(ctx context.Context, userID string, req dto.UpdateCartItemRequest) error
	
	// DeleteItem deletes a cart item
	DeleteItem(ctx context.Context, userID string, req dto.DeleteCartItemRequest) error
	
	// ClearCart clears all items from the user's cart
	ClearCart(ctx context.Context, userID string) error
	
	// UpdateCheckStatus updates the check status of a cart item
	UpdateCheckStatus(ctx context.Context, userID string, req dto.UpdateCheckStatusRequest) error
	
	// UpdateAllCheckStatus updates the check status of all cart items for a user
	UpdateAllCheckStatus(ctx context.Context, userID string, req dto.UpdateAllCheckStatusRequest) error
	
	// GetCartCount returns the number of items in the user's cart
	GetCartCount(ctx context.Context, userID string) (int, error)
	
	// GetCartTotalAmount returns the total amount for checked items in the user's cart
	GetCartTotalAmount(ctx context.Context, userID string) (float64, error)
	
	// GetCheckedItems returns all checked items in the user's cart
	GetCheckedItems(ctx context.Context, userID string) ([]dto.CartItemDTO, error)
	
	// DeleteCheckedItems deletes all checked items from the user's cart
	DeleteCheckedItems(ctx context.Context, userID string) error
	
	// ApplyPromotion applies a promotion to cart items
	ApplyPromotion(ctx context.Context, userID string, req dto.ApplyPromotionRequest) error
	
	// ApplyCoupon applies a coupon to cart items
	ApplyCoupon(ctx context.Context, userID string, req dto.ApplyCouponRequest) error
}
