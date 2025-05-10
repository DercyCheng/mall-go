package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/inventory-service/application/dto"
	"mall-go/services/inventory-service/application/service"
	"mall-go/services/inventory-service/proto/inventorypb"
)

// WarehouseHandler handles warehouse-related gRPC requests
type WarehouseHandler struct {
	warehouseService service.WarehouseService
	inventorypb.UnimplementedWarehouseServiceServer
}

// NewWarehouseHandler creates a new warehouse handler
func NewWarehouseHandler(warehouseService service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseService: warehouseService,
	}
}

// CreateWarehouse creates a new warehouse
func (h *WarehouseHandler) CreateWarehouse(ctx context.Context, req *inventorypb.CreateWarehouseRequest) (*inventorypb.WarehouseResponse, error) {
	// Validate request
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "warehouse name is required")
	}

	// Map request to DTO
	createReq := dto.CreateWarehouseRequest{
		Name:        req.Name,
		Code:        req.Code,
		Address:     req.Address,
		ContactName: req.ContactName,
		ContactPhone: req.ContactPhone,
		Type:        req.Type,
		Status:      req.Status,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	// Call service
	warehouse, err := h.warehouseService.CreateWarehouse(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create warehouse: %v", err)
	}

	// Map response
	return &inventorypb.WarehouseResponse{
		Success:   true,
		Message:   "Warehouse created successfully",
		Warehouse: mapToWarehouseProto(warehouse),
	}, nil
}

// GetWarehouse retrieves a warehouse by ID
func (h *WarehouseHandler) GetWarehouse(ctx context.Context, req *inventorypb.GetWarehouseRequest) (*inventorypb.WarehouseResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "warehouse ID is required")
	}

	// Call service
	warehouse, err := h.warehouseService.GetWarehouse(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get warehouse: %v", err)
	}

	// Map response
	return &inventorypb.WarehouseResponse{
		Success:   true,
		Message:   "Warehouse retrieved successfully",
		Warehouse: mapToWarehouseProto(warehouse),
	}, nil
}

// GetWarehouseByCode retrieves a warehouse by code
func (h *WarehouseHandler) GetWarehouseByCode(ctx context.Context, req *inventorypb.GetWarehouseByCodeRequest) (*inventorypb.WarehouseResponse, error) {
	// Validate request
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "warehouse code is required")
	}

	// Call service
	warehouse, err := h.warehouseService.GetWarehouseByCode(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get warehouse by code: %v", err)
	}

	// Map response
	return &inventorypb.WarehouseResponse{
		Success:   true,
		Message:   "Warehouse retrieved successfully",
		Warehouse: mapToWarehouseProto(warehouse),
	}, nil
}

// ListWarehouses lists all warehouses with pagination
func (h *WarehouseHandler) ListWarehouses(ctx context.Context, req *inventorypb.ListWarehousesRequest) (*inventorypb.WarehousesResponse, error) {
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Call service
	warehouses, total, err := h.warehouseService.ListWarehouses(ctx, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list warehouses: %v", err)
	}

	// Map response
	items := make([]*inventorypb.Warehouse, len(warehouses))
	for i, warehouse := range warehouses {
		items[i] = mapToWarehouseProto(&warehouse)
	}

	return &inventorypb.WarehousesResponse{
		Success:  true,
		Message:  "Warehouses retrieved successfully",
		Items:    items,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// UpdateWarehouse updates a warehouse
func (h *WarehouseHandler) UpdateWarehouse(ctx context.Context, req *inventorypb.UpdateWarehouseRequest) (*inventorypb.WarehouseResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "warehouse ID is required")
	}

	// Map request to DTO
	updateReq := dto.UpdateWarehouseRequest{
		ID:          req.Id,
		Name:        req.Name,
		Address:     req.Address,
		ContactName: req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:      req.Status,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	// Call service
	warehouse, err := h.warehouseService.UpdateWarehouse(ctx, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update warehouse: %v", err)
	}

	// Map response
	return &inventorypb.WarehouseResponse{
		Success:   true,
		Message:   "Warehouse updated successfully",
		Warehouse: mapToWarehouseProto(warehouse),
	}, nil
}

// SetWarehouseStatus sets the warehouse status
func (h *WarehouseHandler) SetWarehouseStatus(ctx context.Context, req *inventorypb.SetWarehouseStatusRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "warehouse ID and status are required")
	}

	// Map request to DTO
	setReq := dto.SetWarehouseStatusRequest{
		ID:     req.Id,
		Status: req.Status,
	}

	// Call service
	if err := h.warehouseService.SetWarehouseStatus(ctx, setReq); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set warehouse status: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Warehouse status set successfully",
	}, nil
}

// DeleteWarehouse deletes a warehouse
func (h *WarehouseHandler) DeleteWarehouse(ctx context.Context, req *inventorypb.DeleteWarehouseRequest) (*inventorypb.GenericResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "warehouse ID is required")
	}

	// Call service
	if err := h.warehouseService.DeleteWarehouse(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete warehouse: %v", err)
	}

	return &inventorypb.GenericResponse{
		Success: true,
		Message: "Warehouse deleted successfully",
	}, nil
}

// mapToWarehouseProto maps a dto.WarehouseDTO to inventorypb.Warehouse
func mapToWarehouseProto(warehouse *dto.WarehouseDTO) *inventorypb.Warehouse {
	return &inventorypb.Warehouse{
		Id:          warehouse.ID,
		Name:        warehouse.Name,
		Code:        warehouse.Code,
		Address:     warehouse.Address,
		ContactName: warehouse.ContactName,
		ContactPhone: warehouse.ContactPhone,
		Type:        warehouse.Type,
		Status:      warehouse.Status,
		Latitude:    warehouse.Latitude,
		Longitude:   warehouse.Longitude,
		CreatedAt:   warehouse.CreatedAt,
		UpdatedAt:   warehouse.UpdatedAt,
	}
}
