package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mall-go/services/order-service/application/dto"
	"mall-go/services/order-service/application/service"
)

// OrderHandler 处理订单相关HTTP请求
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从token中获取用户ID（实际项目中应该从认证中间件获取）
	userID := c.GetString("userID")
	if userID != "" {
		req.UserID = userID
	}

	// 创建订单
	resp, err := h.orderService.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// PayOrder 支付订单
func (h *OrderHandler) PayOrder(c *gin.Context) {
	var req dto.OrderPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 支付订单
	err := h.orderService.PayOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order paid successfully"})
}

// ShipOrder 发货订单
func (h *OrderHandler) ShipOrder(c *gin.Context) {
	var req dto.OrderShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 发货订单
	err := h.orderService.ShipOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order shipped successfully"})
}

// UpdateOrderStatus 更新订单状态
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req dto.OrderStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新订单状态
	err := h.orderService.UpdateOrderStatus(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	reason := c.Query("reason")

	// 取消订单
	err := h.orderService.CancelOrder(c.Request.Context(), orderID, reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}

// GetOrder 获取订单详情
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	// 获取订单
	order, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderByOrderSN 根据订单编号获取订单
func (h *OrderHandler) GetOrderByOrderSN(c *gin.Context) {
	orderSN := c.Param("orderSN")

	// 获取订单
	order, err := h.orderService.GetOrderByOrderSN(c.Request.Context(), orderSN)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetUserOrders 获取用户订单列表
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.Param("userID")

	// 如果没有提供用户ID，则尝试从认证信息中获取
	if userID == "" {
		userID = c.GetString("userID")
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 获取用户订单
	orders, err := h.orderService.GetUserOrders(c.Request.Context(), userID, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// SearchOrders 搜索订单
func (h *OrderHandler) SearchOrders(c *gin.Context) {
	// 构建查询请求
	req := dto.OrderQueryRequest{
		Keyword:   c.Query("keyword"),
		Status:    c.Query("status"),
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
		ProductID: c.Query("productId"),
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	req.Page = page
	req.Size = size

	// 搜索订单
	orders, err := h.orderService.SearchOrders(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrderStatistics 获取订单统计数据
func (h *OrderHandler) GetOrderStatistics(c *gin.Context) {
	// 获取订单统计
	stats, err := h.orderService.GetOrderStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DeleteOrder 删除订单
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	// 删除订单
	err := h.orderService.DeleteOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
