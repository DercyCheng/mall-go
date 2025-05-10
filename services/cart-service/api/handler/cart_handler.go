package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/cart-service/application/dto"
	"mall-go/services/cart-service/application/service"
	"mall-go/services/cart-service/proto/cartpb"
)

// CartHandler handles cart-related gRPC requests
type CartHandler struct {
	cartService service.CartService
	cartpb.UnimplementedCartServiceServer
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// AddItem adds a new item to the cart
func (h *CartHandler) AddItem(ctx context.Context, req *cartpb.AddItemRequest) (*cartpb.CartItemResponse, error) {
	// Validate request
	if req.UserId == "" || req.ProductId == "" || req.ProductSkuId == "" || req.ProductName == "" {
		return nil, status.Error(codes.InvalidArgument, "missing required fields")
	}

	// Map request to DTO
	addReq := dto.AddCartItemRequest{
		ProductID:       req.ProductId,
		ProductSKUID:    req.ProductSkuId,
		ProductName:     req.ProductName,
		ProductSubTitle: req.ProductSubTitle,
		ProductPic:      req.ProductPic,
		ProductPrice:    req.ProductPrice,
		Quantity:        int(req.Quantity),
		ProductSkuCode:  req.ProductSkuCode,
		ProductAttr:     req.ProductAttr,
	}

	// Call service
	cartItem, err := h.cartService.AddItem(ctx, req.UserId, addReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add item: %v", err)
	}

	// Map response
	return &cartpb.CartItemResponse{
		Success: true,
		Message: "Item added to cart successfully",
		Item:    mapToCartItemProto(cartItem),
	}, nil
}

// GetCart retrieves the user's cart
func (h *CartHandler) GetCart(ctx context.Context, req *cartpb.GetCartRequest) (*cartpb.GetCartResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Call service
	cart, err := h.cartService.GetCart(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	// Map response
	items := make([]*cartpb.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = mapToCartItemProto(&item)
	}

	return &cartpb.GetCartResponse{
		Success: true,
		Message: "Cart retrieved successfully",
		Items:   items,
		Summary: &cartpb.CartSummary{
			TotalAmount:     cart.Summary.TotalAmount,
			PromotionAmount: cart.Summary.PromotionAmount,
			CouponAmount:    cart.Summary.CouponAmount,
			RealAmount:      cart.Summary.RealAmount,
			ItemCount:       int32(cart.Summary.ItemCount),
		},
	}, nil
}

// UpdateItemQuantity updates the quantity of a cart item
func (h *CartHandler) UpdateItemQuantity(ctx context.Context, req *cartpb.UpdateItemQuantityRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and item ID are required")
	}

	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than 0")
	}

	// Map request to DTO
	updateReq := dto.UpdateCartItemRequest{
		ID:       req.Id,
		Quantity: int(req.Quantity),
	}

	// Call service
	err := h.cartService.UpdateItemQuantity(ctx, req.UserId, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update item quantity: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "Item quantity updated successfully",
	}, nil
}

// DeleteItem deletes a cart item
func (h *CartHandler) DeleteItem(ctx context.Context, req *cartpb.DeleteItemRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and item ID are required")
	}

	// Map request to DTO
	deleteReq := dto.DeleteCartItemRequest{
		ID: req.Id,
	}

	// Call service
	err := h.cartService.DeleteItem(ctx, req.UserId, deleteReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete item: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "Item deleted successfully",
	}, nil
}

// ClearCart clears all items from the user's cart
func (h *CartHandler) ClearCart(ctx context.Context, req *cartpb.ClearCartRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Call service
	err := h.cartService.ClearCart(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clear cart: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "Cart cleared successfully",
	}, nil
}

// UpdateCheckStatus updates the check status of a cart item
func (h *CartHandler) UpdateCheckStatus(ctx context.Context, req *cartpb.UpdateCheckStatusRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID and item ID are required")
	}

	// Map request to DTO
	updateReq := dto.UpdateCheckStatusRequest{
		ID:     req.Id,
		Status: req.Status,
	}

	// Call service
	err := h.cartService.UpdateCheckStatus(ctx, req.UserId, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update check status: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "Check status updated successfully",
	}, nil
}

// UpdateAllCheckStatus updates the check status of all cart items for a user
func (h *CartHandler) UpdateAllCheckStatus(ctx context.Context, req *cartpb.UpdateAllCheckStatusRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Map request to DTO
	updateReq := dto.UpdateAllCheckStatusRequest{
		Status: req.Status,
	}

	// Call service
	err := h.cartService.UpdateAllCheckStatus(ctx, req.UserId, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update all check status: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "All check status updated successfully",
	}, nil
}

// GetCartCount returns the number of items in the user's cart
func (h *CartHandler) GetCartCount(ctx context.Context, req *cartpb.GetCartCountRequest) (*cartpb.CartCountResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Call service
	count, err := h.cartService.GetCartCount(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart count: %v", err)
	}

	// Return response
	return &cartpb.CartCountResponse{
		Success: true,
		Message: "Cart count retrieved successfully",
		Count:   int32(count),
	}, nil
}

// GetCartTotalAmount returns the total amount for checked items in the user's cart
func (h *CartHandler) GetCartTotalAmount(ctx context.Context, req *cartpb.GetCartTotalAmountRequest) (*cartpb.CartTotalAmountResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Call service
	totalAmount, err := h.cartService.GetCartTotalAmount(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart total amount: %v", err)
	}

	// Return response
	return &cartpb.CartTotalAmountResponse{
		Success:     true,
		Message:     "Cart total amount retrieved successfully",
		TotalAmount: totalAmount,
	}, nil
}

// GetCheckedItems returns all checked items in the user's cart
func (h *CartHandler) GetCheckedItems(ctx context.Context, req *cartpb.GetCheckedItemsRequest) (*cartpb.GetCheckedItemsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Call service
	items, err := h.cartService.GetCheckedItems(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get checked items: %v", err)
	}

	// Map response
	checkedItems := make([]*cartpb.CartItem, len(items))
	for i, item := range items {
		checkedItems[i] = mapToCartItemProto(&item)
	}

	// Return response
	return &cartpb.GetCheckedItemsResponse{
		Success: true,
		Message: "Checked items retrieved successfully",
		Items:   checkedItems,
	}, nil
}

// DeleteCheckedItems deletes all checked items from the user's cart
func (h *CartHandler) DeleteCheckedItems(ctx context.Context, req *cartpb.DeleteCheckedItemsRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	// Call service
	err := h.cartService.DeleteCheckedItems(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete checked items: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "Checked items deleted successfully",
	}, nil
}

// ApplyPromotion applies a promotion to cart items
func (h *CartHandler) ApplyPromotion(ctx context.Context, req *cartpb.ApplyPromotionRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" || len(req.ItemIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID and item IDs are required")
	}

	// Map request to DTO
	applyReq := dto.ApplyPromotionRequest{
		ItemIDs:         req.ItemIds,
		PromotionInfo:   req.PromotionInfo,
		PromotionAmount: req.PromotionAmount,
	}

	// Call service
	err := h.cartService.ApplyPromotion(ctx, req.UserId, applyReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to apply promotion: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "Promotion applied successfully",
	}, nil
}

// ApplyCoupon applies a coupon to cart items
func (h *CartHandler) ApplyCoupon(ctx context.Context, req *cartpb.ApplyCouponRequest) (*cartpb.GenericResponse, error) {
	// Validate request
	if req.UserId == "" || len(req.ItemIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID and item IDs are required")
	}

	// Map request to DTO
	applyReq := dto.ApplyCouponRequest{
		ItemIDs:      req.ItemIds,
		CouponAmount: req.CouponAmount,
	}

	// Call service
	err := h.cartService.ApplyCoupon(ctx, req.UserId, applyReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to apply coupon: %v", err)
	}

	// Return response
	return &cartpb.GenericResponse{
		Success: true,
		Message: "Coupon applied successfully",
	}, nil
}

// mapToCartItemProto maps a dto.CartItemDTO to cartpb.CartItem
func mapToCartItemProto(item *dto.CartItemDTO) *cartpb.CartItem {
	return &cartpb.CartItem{
		Id:               item.ID,
		ProductId:        item.ProductID,
		ProductSkuId:     item.ProductSKUID,
		ProductName:      item.ProductName,
		ProductSubTitle:  item.ProductSubTitle,
		ProductPic:       item.ProductPic,
		ProductPrice:     item.ProductPrice,
		ProductQuantity:  int32(item.ProductQuantity),
		ProductSkuCode:   item.ProductSkuCode,
		ProductCategoryId: item.ProductCategoryID,
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
