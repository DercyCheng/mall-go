package service

import (
	"context"
	"fmt"

	"mall-go/services/inventory-service/application/dto"
	"mall-go/services/inventory-service/domain/model"
	"mall-go/services/inventory-service/domain/repository"
)

// WarehouseServiceImpl implements the WarehouseService interface
type WarehouseServiceImpl struct {
	warehouseRepo repository.WarehouseRepository
	cache         repository.InventoryCache
}

// NewWarehouseService creates a new warehouse service
func NewWarehouseService(
	warehouseRepo repository.WarehouseRepository,
	cache repository.InventoryCache,
) WarehouseService {
	return &WarehouseServiceImpl{
		warehouseRepo: warehouseRepo,
		cache:         cache,
	}
}

// CreateWarehouse creates a new warehouse
func (s *WarehouseServiceImpl) CreateWarehouse(ctx context.Context, req dto.CreateWarehouseRequest) (*dto.WarehouseDTO, error) {
	// Create warehouse domain model
	warehouse, err := model.NewWarehouse(req.Name, req.Code, req.Address, req.ContactName, req.ContactPhone, req.IsDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}

	// Save to repository
	if err := s.warehouseRepo.CreateWarehouse(ctx, warehouse); err != nil {
		return nil, fmt.Errorf("failed to save warehouse: %w", err)
	}

	// Cache the new warehouse
	if err := s.cache.SetWarehouse(ctx, warehouse); err != nil {
		// Log error but continue
		fmt.Printf("failed to cache warehouse: %v", err)
	}

	return s.mapToWarehouseDTO(warehouse), nil
}

// GetWarehouse retrieves a warehouse by ID
func (s *WarehouseServiceImpl) GetWarehouse(ctx context.Context, id string) (*dto.WarehouseDTO, error) {
	// Try to get from cache first
	cached, err := s.cache.GetWarehouse(ctx, id)
	if err != nil {
		// Log error but continue to repository
		fmt.Printf("failed to get warehouse from cache: %v", err)
	}

	if cached != nil {
		return s.mapToWarehouseDTO(cached), nil
	}

	// Fetch from repository
	warehouse, err := s.warehouseRepo.GetWarehouse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse: %w", err)
	}

	// Update cache
	if err := s.cache.SetWarehouse(ctx, warehouse); err != nil {
		// Log error but continue
		fmt.Printf("failed to update warehouse cache: %v", err)
	}

	return s.mapToWarehouseDTO(warehouse), nil
}

// GetWarehouseByCode retrieves a warehouse by code
func (s *WarehouseServiceImpl) GetWarehouseByCode(ctx context.Context, code string) (*dto.WarehouseDTO, error) {
	// No direct cache for this query pattern
	warehouse, err := s.warehouseRepo.GetWarehouseByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse by code: %w", err)
	}

	// Update cache by ID
	if err := s.cache.SetWarehouse(ctx, warehouse); err != nil {
		// Log error but continue
		fmt.Printf("failed to update warehouse cache: %v", err)
	}

	return s.mapToWarehouseDTO(warehouse), nil
}

// GetDefaultWarehouse retrieves the default warehouse
func (s *WarehouseServiceImpl) GetDefaultWarehouse(ctx context.Context) (*dto.WarehouseDTO, error) {
	// No cache for this query pattern
	warehouse, err := s.warehouseRepo.GetDefaultWarehouse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get default warehouse: %w", err)
	}

	// Update cache by ID
	if err := s.cache.SetWarehouse(ctx, warehouse); err != nil {
		// Log error but continue
		fmt.Printf("failed to update warehouse cache: %v", err)
	}

	return s.mapToWarehouseDTO(warehouse), nil
}

// UpdateWarehouse updates an existing warehouse
func (s *WarehouseServiceImpl) UpdateWarehouse(ctx context.Context, req dto.UpdateWarehouseRequest) error {
	// Get existing warehouse
	warehouse, err := s.warehouseRepo.GetWarehouse(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to get warehouse: %w", err)
	}

	// Update fields
	if req.Name != "" {
		warehouse.Name = req.Name
	}
	if req.Code != "" {
		warehouse.Code = req.Code
	}
	if req.Address != "" {
		warehouse.Address = req.Address
	}
	if req.ContactName != "" {
		warehouse.ContactName = req.ContactName
	}
	if req.ContactPhone != "" {
		warehouse.ContactPhone = req.ContactPhone
	}
	if req.Status != warehouse.Status {
		warehouse.Status = req.Status
	}
	warehouse.IsDefault = req.IsDefault

	// Save to repository
	if err := s.warehouseRepo.UpdateWarehouse(ctx, warehouse); err != nil {
		return fmt.Errorf("failed to update warehouse: %w", err)
	}

	// Update cache
	if err := s.cache.SetWarehouse(ctx, warehouse); err != nil {
		// Log error but continue
		fmt.Printf("failed to update warehouse cache: %v", err)
	}

	return nil
}

// DeleteWarehouse deletes a warehouse
func (s *WarehouseServiceImpl) DeleteWarehouse(ctx context.Context, id string) error {
	// Delete from repository
	if err := s.warehouseRepo.DeleteWarehouse(ctx, id); err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}

	// Delete from cache
	if err := s.cache.DeleteWarehouse(ctx, id); err != nil {
		// Log error but continue
		fmt.Printf("failed to delete warehouse from cache: %v", err)
	}

	return nil
}

// ListWarehouses lists all warehouses
func (s *WarehouseServiceImpl) ListWarehouses(ctx context.Context) ([]dto.WarehouseDTO, error) {
	warehouses, err := s.warehouseRepo.ListWarehouses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list warehouses: %w", err)
	}

	// Map to DTOs
	dtos := make([]dto.WarehouseDTO, len(warehouses))
	for i, warehouse := range warehouses {
		dtos[i] = *s.mapToWarehouseDTO(warehouse)
	}

	return dtos, nil
}

// SetDefaultWarehouse sets a warehouse as the default
func (s *WarehouseServiceImpl) SetDefaultWarehouse(ctx context.Context, req dto.SetDefaultWarehouseRequest) error {
	// Set as default in repository
	if err := s.warehouseRepo.SetDefaultWarehouse(ctx, req.ID); err != nil {
		return fmt.Errorf("failed to set default warehouse: %w", err)
	}

	// Clear cache for all warehouses (can be optimized in the future)
	warehouses, err := s.warehouseRepo.ListWarehouses(ctx)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to list warehouses for cache update: %v", err)
		return nil
	}

	// Update cache for all warehouses
	for _, warehouse := range warehouses {
		if err := s.cache.SetWarehouse(ctx, warehouse); err != nil {
			// Log error but continue
			fmt.Printf("failed to update warehouse cache: %v", err)
		}
	}

	return nil
}

// EnableWarehouse enables a warehouse
func (s *WarehouseServiceImpl) EnableWarehouse(ctx context.Context, id string) error {
	// Enable in repository
	if err := s.warehouseRepo.EnableWarehouse(ctx, id); err != nil {
		return fmt.Errorf("failed to enable warehouse: %w", err)
	}

	// Update cache
	if err := s.cache.DeleteWarehouse(ctx, id); err != nil {
		// Log error but continue
		fmt.Printf("failed to update warehouse cache: %v", err)
	}

	return nil
}

// DisableWarehouse disables a warehouse
func (s *WarehouseServiceImpl) DisableWarehouse(ctx context.Context, id string) error {
	// Disable in repository
	if err := s.warehouseRepo.DisableWarehouse(ctx, id); err != nil {
		return fmt.Errorf("failed to disable warehouse: %w", err)
	}

	// Update cache
	if err := s.cache.DeleteWarehouse(ctx, id); err != nil {
		// Log error but continue
		fmt.Printf("failed to update warehouse cache: %v", err)
	}

	return nil
}

// mapToWarehouseDTO maps a model.Warehouse to dto.WarehouseDTO
func (s *WarehouseServiceImpl) mapToWarehouseDTO(warehouse *model.Warehouse) *dto.WarehouseDTO {
	return &dto.WarehouseDTO{
		ID:           warehouse.ID,
		Name:         warehouse.Name,
		Code:         warehouse.Code,
		Address:      warehouse.Address,
		ContactName:  warehouse.ContactName,
		ContactPhone: warehouse.ContactPhone,
		Status:       warehouse.Status,
		IsDefault:    warehouse.IsDefault,
	}
}
