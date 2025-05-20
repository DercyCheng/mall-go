package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mall-go/services/inventory-service/application/dto"
	"mall-go/services/inventory-service/application/service"
)

// InventoryHandler 处理库存相关HTTP请求
type InventoryHandler struct {
	inventoryService service.InventoryService
}

// NewInventoryHandler 创建库存处理器
func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

// CreateInventory 创建库存
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req dto.InventoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建库存
	resp, err := h.inventoryService.CreateInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateInventory 更新库存信息
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	var req dto.InventoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保ID参数一致
	if req.ID == "" {
		req.ID = c.Param("id")
	}

	// 更新库存
	err := h.inventoryService.UpdateInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory updated successfully"})
}

// AdjustInventory 调整库存数量
func (h *InventoryHandler) AdjustInventory(c *gin.Context) {
	var req dto.InventoryAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人ID（实际项目中应该从认证中间件获取）
	operatorID := c.GetString("userID")
	if operatorID != "" {
		req.OperatorID = operatorID
	}

	// 调整库存
	err := h.inventoryService.AdjustInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory adjusted successfully"})
}

// InboundInventory 入库操作
func (h *InventoryHandler) InboundInventory(c *gin.Context) {
	var req dto.InventoryInboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人ID
	operatorID := c.GetString("userID")
	if operatorID != "" {
		req.OperatorID = operatorID
	}

	// 执行入库
	err := h.inventoryService.InboundInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory inbound successfully"})
}

// OutboundInventory 出库操作
func (h *InventoryHandler) OutboundInventory(c *gin.Context) {
	var req dto.InventoryOutboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人ID
	operatorID := c.GetString("userID")
	if operatorID != "" {
		req.OperatorID = operatorID
	}

	// 执行出库
	err := h.inventoryService.OutboundInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory outbound successfully"})
}

// ReserveInventory 预留库存
func (h *InventoryHandler) ReserveInventory(c *gin.Context) {
	var req dto.InventoryReserveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人ID
	operatorID := c.GetString("userID")
	if operatorID != "" {
		req.OperatorID = operatorID
	}

	// 预留库存
	err := h.inventoryService.ReserveInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory reserved successfully"})
}

// ReleaseInventory 释放预留库存
func (h *InventoryHandler) ReleaseInventory(c *gin.Context) {
	var req dto.InventoryReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人ID
	operatorID := c.GetString("userID")
	if operatorID != "" {
		req.OperatorID = operatorID
	}

	// 释放库存
	err := h.inventoryService.ReleaseInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory released successfully"})
}

// LockInventory 锁定库存
func (h *InventoryHandler) LockInventory(c *gin.Context) {
	var req dto.InventoryLockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人ID
	operatorID := c.GetString("userID")
	if operatorID != "" {
		req.OperatorID = operatorID
	}

	// 锁定库存
	err := h.inventoryService.LockInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory locked successfully"})
}

// UnlockInventory 解锁库存
func (h *InventoryHandler) UnlockInventory(c *gin.Context) {
	var req dto.InventoryUnlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录操作人ID
	operatorID := c.GetString("userID")
	if operatorID != "" {
		req.OperatorID = operatorID
	}

	// 解锁库存
	err := h.inventoryService.UnlockInventory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory unlocked successfully"})
}

// GetInventory 获取库存详情
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	id := c.Param("id")

	// 获取库存
	inventory, err := h.inventoryService.GetInventory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// GetInventoryByProductID 根据商品ID获取库存
func (h *InventoryHandler) GetInventoryByProductID(c *gin.Context) {
	productID := c.Param("productId")

	// 获取库存
	inventory, err := h.inventoryService.GetInventoryByProductID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// GetInventoryBySKU 根据SKU获取库存
func (h *InventoryHandler) GetInventoryBySKU(c *gin.Context) {
	sku := c.Param("sku")

	// 获取库存
	inventory, err := h.inventoryService.GetInventoryBySKU(c.Request.Context(), sku)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

// GetLowStockInventories 获取低库存商品列表
func (h *InventoryHandler) GetLowStockInventories(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 获取低库存列表
	inventories, err := h.inventoryService.GetLowStockInventories(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventories)
}

// SearchInventories 搜索库存
func (h *InventoryHandler) SearchInventories(c *gin.Context) {
	// 构建查询请求
	req := dto.InventoryQueryRequest{
		Keyword:     c.Query("keyword"),
		Status:      c.Query("status"),
		WarehouseID: c.Query("warehouseId"),
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	req.Page = page
	req.Size = size

	// 搜索库存
	inventories, err := h.inventoryService.SearchInventories(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventories)
}

// GetInventoryOperations 获取库存操作记录
func (h *InventoryHandler) GetInventoryOperations(c *gin.Context) {
	productID := c.Param("productId")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 获取操作记录
	operations, err := h.inventoryService.GetInventoryOperations(c.Request.Context(), productID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, operations)
}

// GetInventoryOperationsByOrderID 根据订单ID获取库存操作记录
func (h *InventoryHandler) GetInventoryOperationsByOrderID(c *gin.Context) {
	orderID := c.Param("orderId")

	// 获取操作记录
	operations, err := h.inventoryService.GetInventoryOperationsByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, operations)
}
