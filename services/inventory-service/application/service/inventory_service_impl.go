package service

import (
	"context"
	"fmt"
	"time"

	"mall-go/services/inventory-service/application/dto"
	"mall-go/services/inventory-service/domain/model"
	"mall-go/services/inventory-service/domain/repository"
)

// InventoryServiceImpl implements the InventoryService interface
type InventoryServiceImpl struct {
	inventoryRepo       repository.InventoryRepository
	inventoryHistoryRepo repository.InventoryHistoryRepository
	warehouseRepo       repository.WarehouseRepository
	inventoryCache      repository.InventoryCache
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	inventoryRepo repository.InventoryRepository,
	inventoryHistoryRepo repository.InventoryHistoryRepository,
	warehouseRepo repository.WarehouseRepository,
	inventoryCache repository.InventoryCache,
) InventoryService {
	return &InventoryServiceImpl{
		inventoryRepo:       inventoryRepo,
		inventoryHistoryRepo: inventoryHistoryRepo,
		warehouseRepo:       warehouseRepo,
		inventoryCache:      inventoryCache,
	}
}

// CreateInventory creates a new inventory item
func (s *InventoryServiceImpl) CreateInventory(ctx context.Context, req dto.CreateInventoryRequest) (*dto.InventoryItemDTO, error) {
	// Validate warehouse exists
	_, err := s.warehouseRepo.GetWarehouse(ctx, req.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate warehouse: %w", err)
	}

	// Create inventory item
	inventory, err := model.NewInventoryItem(req.ProductID, req.SkuID, req.SkuCode, req.WarehouseID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory item: %w", err)
	}

	// Save to repository
	if err := s.inventoryRepo.CreateInventory(ctx, inventory); err != nil {
		return nil, fmt.Errorf("failed to save inventory item: %w", err)
	}

	// Create history record if initial quantity > 0
	if req.Quantity > 0 {
		history := model.NewInventoryHistory(
			inventory.ID,
			inventory.ProductID,
			inventory.SkuID,
			inventory.WarehouseID,
			"add",
			req.Quantity,
			0,
			req.Quantity,
			"system", // Default operator for initial creation
			"",       // No order ID for initial creation
			"Initial stock",
		)

		if err := s.inventoryHistoryRepo.CreateHistory(ctx, history); err != nil {
			// Log error but continue - history creation should not fail the main operation
			fmt.Printf("failed to create inventory history: %v", err)
		}
	}

	// Cache the new inventory
	if err := s.inventoryCache.SetInventory(ctx, inventory); err != nil {
		// Log error but continue
		fmt.Printf("failed to cache inventory: %v", err)
	}

	return s.mapToInventoryItemDTO(inventory), nil
}

// GetInventory retrieves an inventory item by ID
func (s *InventoryServiceImpl) GetInventory(ctx context.Context, id string) (*dto.InventoryItemDTO, error) {
	// Try to get from cache first
	cached, err := s.inventoryCache.GetInventory(ctx, id)
	if err != nil {
		// Log error but continue to repository
		fmt.Printf("failed to get inventory from cache: %v", err)
	}

	if cached != nil {
		return s.mapToInventoryItemDTO(cached), nil
	}

	// Fetch from repository
	inventory, err := s.inventoryRepo.GetInventory(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	// Update cache
	if err := s.inventoryCache.SetInventory(ctx, inventory); err != nil {
		// Log error but continue
		fmt.Printf("failed to update inventory cache: %v", err)
	}

	return s.mapToInventoryItemDTO(inventory), nil
}

// GetInventoryByProduct retrieves an inventory item by product ID and SKU ID
func (s *InventoryServiceImpl) GetInventoryByProduct(ctx context.Context, productID, skuID string) (*dto.InventoryItemDTO, error) {
	// No direct cache for this query pattern
	inventory, err := s.inventoryRepo.GetInventoryByProduct(ctx, productID, skuID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory by product: %w", err)
	}

	// Update cache
	if err := s.inventoryCache.SetInventory(ctx, inventory); err != nil {
		// Log error but continue
		fmt.Printf("failed to update inventory cache: %v", err)
	}

	return s.mapToInventoryItemDTO(inventory), nil
}

// GetInventoriesByProductID retrieves all inventory items for a product
func (s *InventoryServiceImpl) GetInventoriesByProductID(ctx context.Context, productID string) ([]dto.InventoryItemDTO, error) {
	inventories, err := s.inventoryRepo.GetInventoriesByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventories by product ID: %w", err)
	}

	// Map to DTOs
	dtos := make([]dto.InventoryItemDTO, len(inventories))
	for i, inventory := range inventories {
		dtos[i] = *s.mapToInventoryItemDTO(inventory)
	}

	return dtos, nil
}

// GetInventoriesByWarehouseID retrieves all inventory items in a warehouse
func (s *InventoryServiceImpl) GetInventoriesByWarehouseID(ctx context.Context, warehouseID string, page, pageSize int) ([]dto.InventoryItemDTO, int, error) {
	inventories, err := s.inventoryRepo.GetInventoriesByWarehouseID(ctx, warehouseID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get inventories by warehouse ID: %w", err)
	}

	// Calculate total and apply pagination
	total := len(inventories)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= total {
		return []dto.InventoryItemDTO{}, total, nil
	}

	if end > total {
		end = total
	}

	paginatedInventories := inventories[start:end]

	// Map to DTOs
	dtos := make([]dto.InventoryItemDTO, len(paginatedInventories))
	for i, inventory := range paginatedInventories {
		dtos[i] = *s.mapToInventoryItemDTO(inventory)
	}

	return dtos, total, nil
}

// ListInventories lists all inventory items with pagination
func (s *InventoryServiceImpl) ListInventories(ctx context.Context, page, pageSize int) ([]dto.InventoryItemDTO, int, error) {
	inventories, total, err := s.inventoryRepo.ListInventories(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list inventories: %w", err)
	}

	// Map to DTOs
	dtos := make([]dto.InventoryItemDTO, len(inventories))
	for i, inventory := range inventories {
		dtos[i] = *s.mapToInventoryItemDTO(inventory)
	}

	return dtos, total, nil
}

// AddStock adds stock to an inventory item
func (s *InventoryServiceImpl) AddStock(ctx context.Context, req dto.AddStockRequest) error {
	// Get inventory item first
	inventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Record before quantity
	beforeQty := inventory.Quantity

	// Add stock in repository
	if err := s.inventoryRepo.AddStock(ctx, req.ID, req.Quantity); err != nil {
		return fmt.Errorf("failed to add stock: %w", err)
	}

	// Get updated inventory
	updatedInventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to get updated inventory: %v", err)
	} else {
		// Create history record
		history := model.NewInventoryHistory(
			updatedInventory.ID,
			updatedInventory.ProductID,
			updatedInventory.SkuID,
			updatedInventory.WarehouseID,
			"add",
			req.Quantity,
			beforeQty,
			updatedInventory.Quantity,
			req.OperatorID,
			"",
			req.Reason,
		)

		if err := s.inventoryHistoryRepo.CreateHistory(ctx, history); err != nil {
			// Log error but continue
			fmt.Printf("failed to create inventory history: %v", err)
		}

		// Update cache
		if err := s.inventoryCache.DeleteInventory(ctx, req.ID); err != nil {
			// Log error but continue
			fmt.Printf("failed to update inventory cache: %v", err)
		}
	}

	return nil
}

// RemoveStock removes stock from an inventory item
func (s *InventoryServiceImpl) RemoveStock(ctx context.Context, req dto.RemoveStockRequest) error {
	// Get inventory item first
	inventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Record before quantity
	beforeQty := inventory.Quantity

	// Remove stock in repository
	if err := s.inventoryRepo.RemoveStock(ctx, req.ID, req.Quantity); err != nil {
		return fmt.Errorf("failed to remove stock: %w", err)
	}

	// Get updated inventory
	updatedInventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to get updated inventory: %v", err)
	} else {
		// Create history record
		history := model.NewInventoryHistory(
			updatedInventory.ID,
			updatedInventory.ProductID,
			updatedInventory.SkuID,
			updatedInventory.WarehouseID,
			"remove",
			req.Quantity,
			beforeQty,
			updatedInventory.Quantity,
			req.OperatorID,
			"",
			req.Reason,
		)

		if err := s.inventoryHistoryRepo.CreateHistory(ctx, history); err != nil {
			// Log error but continue
			fmt.Printf("failed to create inventory history: %v", err)
		}

		// Update cache
		if err := s.inventoryCache.DeleteInventory(ctx, req.ID); err != nil {
			// Log error but continue
			fmt.Printf("failed to update inventory cache: %v", err)
		}
	}

	return nil
}

// LockStock locks stock for a pending order
func (s *InventoryServiceImpl) LockStock(ctx context.Context, req dto.LockStockRequest) error {
	// Get inventory item first
	inventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Record before locked count and available quantity
	beforeLocked := inventory.LockedCount
	beforeAvailable := inventory.AvailableQty

	// Lock stock in repository
	if err := s.inventoryRepo.LockStock(ctx, req.ID, req.Quantity); err != nil {
		return fmt.Errorf("failed to lock stock: %w", err)
	}

	// Get updated inventory
	updatedInventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to get updated inventory: %v", err)
	} else {
		// Create history record
		history := model.NewInventoryHistory(
			updatedInventory.ID,
			updatedInventory.ProductID,
			updatedInventory.SkuID,
			updatedInventory.WarehouseID,
			"lock",
			req.Quantity,
			beforeAvailable,
			updatedInventory.AvailableQty,
			req.OperatorID,
			req.OrderID,
			"Lock stock for order",
		)

		if err := s.inventoryHistoryRepo.CreateHistory(ctx, history); err != nil {
			// Log error but continue
			fmt.Printf("failed to create inventory history: %v", err)
		}

		// Update cache
		if err := s.inventoryCache.DeleteInventory(ctx, req.ID); err != nil {
			// Log error but continue
			fmt.Printf("failed to update inventory cache: %v", err)
		}
	}

	return nil
}

// UnlockStock unlocks previously locked stock
func (s *InventoryServiceImpl) UnlockStock(ctx context.Context, req dto.UnlockStockRequest) error {
	// Get inventory item first
	inventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Record before locked count and available quantity
	beforeLocked := inventory.LockedCount
	beforeAvailable := inventory.AvailableQty

	// Unlock stock in repository
	if err := s.inventoryRepo.UnlockStock(ctx, req.ID, req.Quantity); err != nil {
		return fmt.Errorf("failed to unlock stock: %w", err)
	}

	// Get updated inventory
	updatedInventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to get updated inventory: %v", err)
	} else {
		// Create history record
		history := model.NewInventoryHistory(
			updatedInventory.ID,
			updatedInventory.ProductID,
			updatedInventory.SkuID,
			updatedInventory.WarehouseID,
			"unlock",
			req.Quantity,
			beforeAvailable,
			updatedInventory.AvailableQty,
			req.OperatorID,
			req.OrderID,
			req.Reason,
		)

		if err := s.inventoryHistoryRepo.CreateHistory(ctx, history); err != nil {
			// Log error but continue
			fmt.Printf("failed to create inventory history: %v", err)
		}

		// Update cache
		if err := s.inventoryCache.DeleteInventory(ctx, req.ID); err != nil {
			// Log error but continue
			fmt.Printf("failed to update inventory cache: %v", err)
		}
	}

	return nil
}

// ConfirmLock confirms a locked stock (e.g., when order is paid)
func (s *InventoryServiceImpl) ConfirmLock(ctx context.Context, req dto.ConfirmLockRequest) error {
	// Get inventory item first
	inventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get inventory: %w", err)
	}

	// Record before quantity and locked count
	beforeQty := inventory.Quantity
	beforeLocked := inventory.LockedCount

	// Confirm lock in repository
	if err := s.inventoryRepo.ConfirmLock(ctx, req.ID, req.Quantity); err != nil {
		return fmt.Errorf("failed to confirm lock: %w", err)
	}

	// Get updated inventory
	updatedInventory, err := s.inventoryRepo.GetInventory(ctx, req.ID)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to get updated inventory: %v", err)
	} else {
		// Create history record
		history := model.NewInventoryHistory(
			updatedInventory.ID,
			updatedInventory.ProductID,
			updatedInventory.SkuID,
			updatedInventory.WarehouseID,
			"confirm",
			req.Quantity,
			beforeQty,
			updatedInventory.Quantity,
			req.OperatorID,
			req.OrderID,
			"Confirm locked stock for order",
		)

		if err := s.inventoryHistoryRepo.CreateHistory(ctx, history); err != nil {
			// Log error but continue
			fmt.Printf("failed to create inventory history: %v", err)
		}

		// Update cache
		if err := s.inventoryCache.DeleteInventory(ctx, req.ID); err != nil {
			// Log error but continue
			fmt.Printf("failed to update inventory cache: %v", err)
		}
	}

	return nil
}

// SetInventoryStatus sets the inventory status
func (s *InventoryServiceImpl) SetInventoryStatus(ctx context.Context, req dto.SetInventoryStatusRequest) error {
	// Map status string to enum
	var status model.InventoryStatus
	switch req.Status {
	case "normal":
		status = model.InventoryStatusNormal
	case "locked":
		status = model.InventoryStatusLocked
	case "sold":
		status = model.InventoryStatusSold
	case "defective":
		status = model.InventoryStatusDefective
	default:
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	// Update status in repository
	if err := s.inventoryRepo.SetStatus(ctx, req.ID, status); err != nil {
		return fmt.Errorf("failed to set status: %w", err)
	}

	// Update cache
	if err := s.inventoryCache.DeleteInventory(ctx, req.ID); err != nil {
		// Log error but continue
		fmt.Printf("failed to update inventory cache: %v", err)
	}

	return nil
}

// GetInventoryHistory retrieves history records for an inventory item
func (s *InventoryServiceImpl) GetInventoryHistory(ctx context.Context, inventoryID string, page, pageSize int) ([]dto.InventoryHistoryDTO, int, error) {
	records, total, err := s.inventoryHistoryRepo.GetHistoryByInventoryID(ctx, inventoryID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get inventory history: %w", err)
	}

	// Map to DTOs
	dtos := make([]dto.InventoryHistoryDTO, len(records))
	for i, record := range records {
		dtos[i] = s.mapToInventoryHistoryDTO(record)
	}

	return dtos, total, nil
}

// GetProductInventoryHistory retrieves history records for a product
func (s *InventoryServiceImpl) GetProductInventoryHistory(ctx context.Context, productID string, page, pageSize int) ([]dto.InventoryHistoryDTO, int, error) {
	records, total, err := s.inventoryHistoryRepo.GetHistoryByProductID(ctx, productID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get product inventory history: %w", err)
	}

	// Map to DTOs
	dtos := make([]dto.InventoryHistoryDTO, len(records))
	for i, record := range records {
		dtos[i] = s.mapToInventoryHistoryDTO(record)
	}

	return dtos, total, nil
}

// GetWarehouseInventoryHistory retrieves history records for a warehouse
func (s *InventoryServiceImpl) GetWarehouseInventoryHistory(ctx context.Context, warehouseID string, page, pageSize int) ([]dto.InventoryHistoryDTO, int, error) {
	records, total, err := s.inventoryHistoryRepo.GetHistoryByWarehouseID(ctx, warehouseID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get warehouse inventory history: %w", err)
	}

	// Map to DTOs
	dtos := make([]dto.InventoryHistoryDTO, len(records))
	for i, record := range records {
		dtos[i] = s.mapToInventoryHistoryDTO(record)
	}

	return dtos, total, nil
}

// GetOrderInventoryHistory retrieves history records for an order
func (s *InventoryServiceImpl) GetOrderInventoryHistory(ctx context.Context, orderID string) ([]dto.InventoryHistoryDTO, error) {
	records, err := s.inventoryHistoryRepo.GetHistoryByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order inventory history: %w", err)
	}

	// Map to DTOs
	dtos := make([]dto.InventoryHistoryDTO, len(records))
	for i, record := range records {
		dtos[i] = s.mapToInventoryHistoryDTO(record)
	}

	return dtos, nil
}

// CheckStock checks if stock is available for a product
func (s *InventoryServiceImpl) CheckStock(ctx context.Context, req dto.CheckStockRequest) (*dto.CheckStockResponse, error) {
	// Try to get inventory by product and SKU
	inventory, err := s.inventoryRepo.GetInventoryByProduct(ctx, req.ProductID, req.SkuID)
	if err != nil {
		return &dto.CheckStockResponse{
			Success:      true,
			Available:    false,
			AvailableQty: 0,
			Message:      "Product not found in inventory",
		}, nil
	}

	// Check if available quantity is sufficient
	isAvailable := inventory.AvailableQty >= req.Quantity

	return &dto.CheckStockResponse{
		Success:      true,
		Available:    isAvailable,
		AvailableQty: inventory.AvailableQty,
		Message:      fmt.Sprintf("Available quantity: %d", inventory.AvailableQty),
	}, nil
}

// mapToInventoryItemDTO maps a model.InventoryItem to dto.InventoryItemDTO
func (s *InventoryServiceImpl) mapToInventoryItemDTO(item *model.InventoryItem) *dto.InventoryItemDTO {
	return &dto.InventoryItemDTO{
		ID:           item.ID,
		ProductID:    item.ProductID,
		SkuID:        item.SkuID,
		SkuCode:      item.SkuCode,
		WarehouseID:  item.WarehouseID,
		Quantity:     item.Quantity,
		LockedCount:  item.LockedCount,
		AvailableQty: item.AvailableQty,
		Status:       item.Status.String(),
	}
}

// mapToInventoryHistoryDTO maps a model.InventoryHistory to dto.InventoryHistoryDTO
func (s *InventoryServiceImpl) mapToInventoryHistoryDTO(history *model.InventoryHistory) dto.InventoryHistoryDTO {
	return dto.InventoryHistoryDTO{
		ID:            history.ID,
		InventoryID:   history.InventoryID,
		ProductID:     history.ProductID,
		SkuID:         history.SkuID,
		WarehouseID:   history.WarehouseID,
		OperationType: history.OperationType,
		Quantity:      history.Quantity,
		BeforeQty:     history.BeforeQty,
		AfterQty:      history.AfterQty,
		Operator:      history.Operator,
		OrderID:       history.OrderID,
		Reason:        history.Reason,
		CreatedAt:     history.CreatedAt.Format(time.RFC3339),
	}
}
