package service

import (
	"context"
	"fmt"
	
	"mall-go/services/cart-service/application/dto"
	"mall-go/services/cart-service/domain/model"
	"mall-go/services/cart-service/domain/repository"
)

// CartServiceImpl implements the CartService interface
type CartServiceImpl struct {
	cartRepo  repository.CartRepository
	cartCache repository.CartCache
}

// NewCartService creates a new cart service
func NewCartService(cartRepo repository.CartRepository, cartCache repository.CartCache) CartService {
	return &CartServiceImpl{
		cartRepo:  cartRepo,
		cartCache: cartCache,
	}
}

// AddItem adds a new item to the cart
func (s *CartServiceImpl) AddItem(ctx context.Context, userID string, req dto.AddCartItemRequest) (*dto.CartItemDTO, error) {
	// Check if the item already exists in the cart
	existingItem, err := s.cartRepo.CheckExistingItem(ctx, userID, req.ProductID, req.ProductSKUID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing item: %w", err)
	}

	// If the item exists, update the quantity
	if existingItem != nil {
		err := existingItem.IncreaseQuantity(req.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to increase quantity: %w", err)
		}

		if err := s.cartRepo.UpdateItem(ctx, existingItem); err != nil {
			return nil, fmt.Errorf("failed to update item: %w", err)
		}
		
		// Invalidate cache
		if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
			// Log error but continue
			fmt.Printf("failed to delete cart cache: %v", err)
		}
		
		return s.mapToCartItemDTO(existingItem), nil
	}

	// Create a new cart item
	cartItem, err := model.NewCartItem(
		userID,
		req.ProductID,
		req.ProductSKUID,
		req.ProductName,
		req.ProductPic,
		req.ProductSkuCode,
		req.ProductPrice,
		req.Quantity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart item: %w", err)
	}

	// Set additional fields
	cartItem.ProductSubTitle = req.ProductSubTitle
	cartItem.ProductAttr = req.ProductAttr

	// Add item to repository
	if err := s.cartRepo.AddItem(ctx, cartItem); err != nil {
		return nil, fmt.Errorf("failed to add item to repository: %w", err)
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	// Increment cart count in cache
	if err := s.cartCache.IncrementCartCount(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to increment cart count: %v", err)
	}
	
	return s.mapToCartItemDTO(cartItem), nil
}

// GetCart returns the user's cart
func (s *CartServiceImpl) GetCart(ctx context.Context, userID string) (*dto.GetCartResponse, error) {
	// Try to get items from cache first
	cachedItems, err := s.cartCache.GetCartItems(ctx, userID)
	if err != nil {
		// Log error but continue to fetch from DB
		fmt.Printf("failed to get items from cache: %v", err)
	}
	
	var items []*model.CartItem
	if cachedItems != nil && len(cachedItems) > 0 {
		items = cachedItems
	} else {
		// Fetch from repository if not in cache
		dbItems, err := s.cartRepo.GetItemsByUserID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get items from repository: %w", err)
		}
		
		items = dbItems
		
		// Update cache
		if len(dbItems) > 0 {
			if err := s.cartCache.SetCartItems(ctx, userID, dbItems); err != nil {
				// Log error but continue
				fmt.Printf("failed to set items in cache: %v", err)
			}
			
			// Update cart count
			if err := s.cartCache.SetCartCount(ctx, userID, len(dbItems)); err != nil {
				// Log error but continue
				fmt.Printf("failed to set cart count: %v", err)
			}
		}
	}
	
	// Calculate summary
	var summary dto.CartSummaryDTO
	summary.ItemCount = len(items)
	
	for _, item := range items {
		if item.CheckStatus {
			summary.TotalAmount += item.ProductPrice * float64(item.ProductQuantity)
			summary.PromotionAmount += item.PromotionAmount
			summary.CouponAmount += item.CouponAmount
			summary.RealAmount += item.RealAmount
		}
	}
	
	// Map to DTO
	itemDTOs := make([]dto.CartItemDTO, len(items))
	for i, item := range items {
		itemDTOs[i] = *s.mapToCartItemDTO(item)
	}
	
	return &dto.GetCartResponse{
		Items:   itemDTOs,
		Summary: summary,
	}, nil
}

// UpdateItemQuantity updates the quantity of a cart item
func (s *CartServiceImpl) UpdateItemQuantity(ctx context.Context, userID string, req dto.UpdateCartItemRequest) error {
	// Get the cart item
	item, err := s.cartRepo.GetItem(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}
	
	// Check ownership
	if item.UserID != userID {
		return fmt.Errorf("user does not own this cart item")
	}
	
	// Update quantity
	if err := item.UpdateQuantity(req.Quantity); err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}
	
	// Update item in repository
	if err := s.cartRepo.UpdateItem(ctx, item); err != nil {
		return fmt.Errorf("failed to update item in repository: %w", err)
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	return nil
}

// DeleteItem deletes a cart item
func (s *CartServiceImpl) DeleteItem(ctx context.Context, userID string, req dto.DeleteCartItemRequest) error {
	// Get the cart item
	item, err := s.cartRepo.GetItem(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}
	
	// Check ownership
	if item.UserID != userID {
		return fmt.Errorf("user does not own this cart item")
	}
	
	// Delete item from repository
	if err := s.cartRepo.DeleteItem(ctx, req.ID); err != nil {
		return fmt.Errorf("failed to delete item from repository: %w", err)
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	// Decrement cart count in cache
	if err := s.cartCache.DecrementCartCount(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to decrement cart count: %v", err)
	}
	
	return nil
}

// ClearCart clears all items from the user's cart
func (s *CartServiceImpl) ClearCart(ctx context.Context, userID string) error {
	// Delete items from repository
	if err := s.cartRepo.DeleteItemsByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete items from repository: %w", err)
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	// Reset cart count
	if err := s.cartCache.SetCartCount(ctx, userID, 0); err != nil {
		// Log error but continue
		fmt.Printf("failed to reset cart count: %v", err)
	}
	
	return nil
}

// UpdateCheckStatus updates the check status of a cart item
func (s *CartServiceImpl) UpdateCheckStatus(ctx context.Context, userID string, req dto.UpdateCheckStatusRequest) error {
	// Get the cart item
	item, err := s.cartRepo.GetItem(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}
	
	// Check ownership
	if item.UserID != userID {
		return fmt.Errorf("user does not own this cart item")
	}
	
	// Update check status
	if err := s.cartRepo.UpdateCheckStatus(ctx, req.ID, req.Status); err != nil {
		return fmt.Errorf("failed to update check status: %w", err)
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	return nil
}

// UpdateAllCheckStatus updates the check status of all cart items for a user
func (s *CartServiceImpl) UpdateAllCheckStatus(ctx context.Context, userID string, req dto.UpdateAllCheckStatusRequest) error {
	// Update check status
	if err := s.cartRepo.UpdateAllCheckStatus(ctx, userID, req.Status); err != nil {
		return fmt.Errorf("failed to update all check status: %w", err)
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	return nil
}

// GetCartCount returns the number of items in the user's cart
func (s *CartServiceImpl) GetCartCount(ctx context.Context, userID string) (int, error) {
	// Try to get count from cache first
	count, err := s.cartCache.GetCartCount(ctx, userID)
	if err != nil {
		// Log error but continue to fetch from DB
		fmt.Printf("failed to get cart count from cache: %v", err)
	} else if count > 0 {
		return count, nil
	}
	
	// Fetch from repository if not in cache or count is 0
	count, err = s.cartRepo.GetItemCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get item count from repository: %w", err)
	}
	
	// Update cache
	if err := s.cartCache.SetCartCount(ctx, userID, count); err != nil {
		// Log error but continue
		fmt.Printf("failed to set cart count in cache: %v", err)
	}
	
	return count, nil
}

// GetCartTotalAmount returns the total amount for checked items in the user's cart
func (s *CartServiceImpl) GetCartTotalAmount(ctx context.Context, userID string) (float64, error) {
	amount, err := s.cartRepo.GetCartTotalAmount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total amount: %w", err)
	}
	
	return amount, nil
}

// GetCheckedItems returns all checked items in the user's cart
func (s *CartServiceImpl) GetCheckedItems(ctx context.Context, userID string) ([]dto.CartItemDTO, error) {
	items, err := s.cartRepo.GetCheckedItemsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get checked items: %w", err)
	}
	
	// Map to DTO
	itemDTOs := make([]dto.CartItemDTO, len(items))
	for i, item := range items {
		itemDTOs[i] = *s.mapToCartItemDTO(item)
	}
	
	return itemDTOs, nil
}

// DeleteCheckedItems deletes all checked items from the user's cart
func (s *CartServiceImpl) DeleteCheckedItems(ctx context.Context, userID string) error {
	// Delete checked items from repository
	if err := s.cartRepo.DeleteCheckedItemsByUserID(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete checked items: %w", err)
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	// Update cart count
	count, err := s.cartRepo.GetItemCount(ctx, userID)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to get updated item count: %v", err)
	} else {
		if err := s.cartCache.SetCartCount(ctx, userID, count); err != nil {
			// Log error but continue
			fmt.Printf("failed to update cart count: %v", err)
		}
	}
	
	return nil
}

// ApplyPromotion applies a promotion to cart items
func (s *CartServiceImpl) ApplyPromotion(ctx context.Context, userID string, req dto.ApplyPromotionRequest) error {
	for _, itemID := range req.ItemIDs {
		// Get the cart item
		item, err := s.cartRepo.GetItem(ctx, itemID)
		if err != nil {
			return fmt.Errorf("failed to get item %s: %w", itemID, err)
		}
		
		// Check ownership
		if item.UserID != userID {
			return fmt.Errorf("user does not own cart item %s", itemID)
		}
		
		// Apply promotion
		item.ApplyPromotion(req.PromotionInfo, req.PromotionAmount)
		
		// Update item in repository
		if err := s.cartRepo.UpdateItem(ctx, item); err != nil {
			return fmt.Errorf("failed to update item %s: %w", itemID, err)
		}
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	return nil
}

// ApplyCoupon applies a coupon to cart items
func (s *CartServiceImpl) ApplyCoupon(ctx context.Context, userID string, req dto.ApplyCouponRequest) error {
	for _, itemID := range req.ItemIDs {
		// Get the cart item
		item, err := s.cartRepo.GetItem(ctx, itemID)
		if err != nil {
			return fmt.Errorf("failed to get item %s: %w", itemID, err)
		}
		
		// Check ownership
		if item.UserID != userID {
			return fmt.Errorf("user does not own cart item %s", itemID)
		}
		
		// Apply coupon
		item.ApplyCoupon(req.CouponAmount)
		
		// Update item in repository
		if err := s.cartRepo.UpdateItem(ctx, item); err != nil {
			return fmt.Errorf("failed to update item %s: %w", itemID, err)
		}
	}
	
	// Invalidate cache
	if err := s.cartCache.DeleteCartCache(ctx, userID); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete cart cache: %v", err)
	}
	
	return nil
}

// mapToCartItemDTO maps a model.CartItem to dto.CartItemDTO
func (s *CartServiceImpl) mapToCartItemDTO(item *model.CartItem) *dto.CartItemDTO {
	return &dto.CartItemDTO{
		ID:               item.ID,
		ProductID:        item.ProductID,
		ProductSKUID:     item.ProductSKUID,
		ProductName:      item.ProductName,
		ProductSubTitle:  item.ProductSubTitle,
		ProductPic:       item.ProductPic,
		ProductPrice:     item.ProductPrice,
		ProductQuantity:  item.ProductQuantity,
		ProductSkuCode:   item.ProductSkuCode,
		ProductCategoryID: item.ProductCategoryID,
		ProductBrand:     item.ProductBrand,
		ProductSn:        item.ProductSn,
		ProductAttr:      item.ProductAttr,
		PromotionInfo:    item.PromotionInfo,
		PromotionAmount:  item.PromotionAmount,
		CouponAmount:     item.CouponAmount,
		RealAmount:       item.RealAmount,
		CheckStatus:      item.CheckStatus,
	}
}
