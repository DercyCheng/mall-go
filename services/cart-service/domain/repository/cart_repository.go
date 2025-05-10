package repository

import (
	"context"

	"mall-go/services/cart-service/domain/model"
)

// CartRepository defines the interface for cart data access
type CartRepository interface {
	// AddItem adds a new item to the cart
	AddItem(ctx context.Context, item *model.CartItem) error

	// GetItem retrieves a cart item by ID
	GetItem(ctx context.Context, id string) (*model.CartItem, error)

	// GetItemsByUserID retrieves all cart items for a user
	GetItemsByUserID(ctx context.Context, userID string) ([]*model.CartItem, error)

	// GetCheckedItemsByUserID retrieves all checked cart items for a user
	GetCheckedItemsByUserID(ctx context.Context, userID string) ([]*model.CartItem, error)

	// UpdateItem updates an existing cart item
	UpdateItem(ctx context.Context, item *model.CartItem) error

	// DeleteItem deletes a cart item by ID
	DeleteItem(ctx context.Context, id string) error

	// DeleteItemsByUserID deletes all cart items for a user
	DeleteItemsByUserID(ctx context.Context, userID string) error

	// DeleteCheckedItemsByUserID deletes all checked cart items for a user
	DeleteCheckedItemsByUserID(ctx context.Context, userID string) error

	// GetItemCount returns the number of cart items for a user
	GetItemCount(ctx context.Context, userID string) (int, error)

	// GetCartTotalAmount returns the total amount of all checked items in a user's cart
	GetCartTotalAmount(ctx context.Context, userID string) (float64, error)

	// CheckExistingItem checks if an item already exists in the cart
	CheckExistingItem(ctx context.Context, userID, productID, skuID string) (*model.CartItem, error)

	// UpdateCheckStatus updates the check status of a cart item
	UpdateCheckStatus(ctx context.Context, id string, status bool) error

	// UpdateAllCheckStatus updates the check status of all cart items for a user
	UpdateAllCheckStatus(ctx context.Context, userID string, status bool) error

	// Clear clears all cart items in the database (for testing)
	Clear(ctx context.Context) error
}

// CartCache defines the interface for cart data caching
type CartCache interface {
	// GetCartItems gets cart items from cache
	GetCartItems(ctx context.Context, userID string) ([]*model.CartItem, error)

	// SetCartItems sets cart items in cache
	SetCartItems(ctx context.Context, userID string, items []*model.CartItem) error

	// DeleteCartCache deletes cart cache for a user
	DeleteCartCache(ctx context.Context, userID string) error

	// IncrementCartCount increments the cart count for a user
	IncrementCartCount(ctx context.Context, userID string) error

	// DecrementCartCount decrements the cart count for a user
	DecrementCartCount(ctx context.Context, userID string) error

	// GetCartCount gets the cart count for a user
	GetCartCount(ctx context.Context, userID string) (int, error)

	// SetCartCount sets the cart count for a user
	SetCartCount(ctx context.Context, userID string, count int) error
}
