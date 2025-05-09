package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/response"
	"mall-go/services/order-service/application/dto"
	service "mall-go/services/order-service/application/service"
)

// OrderHandler 订单控制器
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler 创建订单控制器
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Create godoc
// @Summary 创建订单
// @Description 创建新订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param order body dto.OrderCreateRequest true "订单信息"
// @Success 200 {object} response.Response{data=dto.CreateOrderResponse}
// @Router /api/v1/orders [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 创建订单
	result, err := h.orderService.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建订单失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// Get godoc
// @Summary 获取订单详情
// @Description 根据ID获取订单详情
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response{data=dto.OrderResponse}
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取订单详情失败: "+err.Error())
		return
	}

	response.Success(c, order)
}

// GetByOrderSn godoc
// @Summary 根据订单编号获取订单
// @Description 根据订单编号获取订单详情
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orderSn path string true "订单编号"
// @Success 200 {object} response.Response{data=dto.OrderResponse}
// @Router /api/v1/orders/sn/{orderSn} [get]
func (h *OrderHandler) GetByOrderSn(c *gin.Context) {
	orderSn := c.Param("orderSn")
	if orderSn == "" {
		response.Fail(c, http.StatusBadRequest, "订单编号不能为空")
		return
	}

	order, err := h.orderService.GetOrderByOrderSn(c.Request.Context(), orderSn)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取订单详情失败: "+err.Error())
		return
	}

	response.Success(c, order)
}

// List godoc
// @Summary 获取订单列表
// @Description 分页查询订单列表
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param orderSn query string false "订单编号"
// @Param status query string false "订单状态"
// @Param memberUsername query string false "会员用户名"
// @Param receiverName query string false "收货人姓名"
// @Param receiverPhone query string false "收货人电话"
// @Param createTimeBegin query string false "创建时间开始"
// @Param createTimeEnd query string false "创建时间结束"
// @Param sourceType query int false "订单来源：0->PC订单；1->app订单"
// @Param orderType query int false "订单类型：0->正常订单；1->秒杀订单"
// @Param page query int true "页码"
// @Param size query int true "每页数量"
// @Success 200 {object} response.Response{data=dto.OrderListResponse}
// @Router /api/v1/orders [get]
func (h *OrderHandler) List(c *gin.Context) {
	var query dto.OrderQueryRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 调用服务查询订单列表
	result, err := h.orderService.ListOrders(c.Request.Context(), &query)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取订单列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// ListMemberOrders godoc
// @Summary 获取会员订单列表
// @Description 分页查询指定会员的订单列表
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param memberId path string true "会员ID"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} response.Response{data=dto.OrderListResponse}
// @Router /api/v1/members/{memberId}/orders [get]
func (h *OrderHandler) ListMemberOrders(c *gin.Context) {
	memberID := c.Param("memberId")
	if memberID == "" {
		response.Fail(c, http.StatusBadRequest, "会员ID不能为空")
		return
	}

	// 获取分页参数，默认为第1页，每页10条
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	// 调用服务查询会员订单列表
	result, err := h.orderService.ListMemberOrders(c.Request.Context(), memberID, page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取会员订单列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// Pay godoc
// @Summary 订单支付
// @Description 支付订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param payRequest body dto.OrderPayRequest true "支付信息"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/pay [post]
func (h *OrderHandler) Pay(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	var req dto.OrderPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 支付订单
	err := h.orderService.PayOrder(c.Request.Context(), id, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "支付订单失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Ship godoc
// @Summary 订单发货
// @Description 订单发货
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param deliveryRequest body dto.OrderDeliveryRequest true "发货信息"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/ship [post]
func (h *OrderHandler) Ship(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	var req dto.OrderDeliveryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 发货
	err := h.orderService.ShipOrder(c.Request.Context(), id, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "订单发货失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Receive godoc
// @Summary 确认收货
// @Description 确认订单收货
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/receive [post]
func (h *OrderHandler) Receive(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	// 确认收货
	err := h.orderService.ReceiveOrder(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "确认收货失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Cancel godoc
// @Summary 取消订单
// @Description 取消订单
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param cancelRequest body dto.OrderCancelRequest true "取消原因"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/cancel [post]
func (h *OrderHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	var req dto.OrderCancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 取消订单
	err := h.orderService.CancelOrder(c.Request.Context(), id, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "取消订单失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete godoc
// @Summary 删除订单
// @Description 删除订单(逻辑删除)
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	// 删除订单
	err := h.orderService.DeleteOrder(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除订单失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateNote godoc
// @Summary 更新订单备注
// @Description 更新订单备注
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param noteRequest body dto.OrderUpdateNoteRequest true "订单备注"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/note [put]
func (h *OrderHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	var req dto.OrderUpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 更新订单备注
	err := h.orderService.UpdateOrderNote(c.Request.Context(), id, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新订单备注失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateReceiverInfo godoc
// @Summary 更新收货人信息
// @Description 更新订单收货人信息
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param receiverRequest body dto.OrderUpdateReceiverInfoRequest true "收货人信息"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/receiver [put]
func (h *OrderHandler) UpdateReceiverInfo(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	var req dto.OrderUpdateReceiverInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 更新收货人信息
	err := h.orderService.UpdateOrderReceiverInfo(c.Request.Context(), id, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新收货人信息失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ApplyRefund godoc
// @Summary 申请退款
// @Description 订单申请退款
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Param refundRequest body dto.OrderRefundRequest true "退款信息"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/refund/apply [post]
func (h *OrderHandler) ApplyRefund(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	var req dto.OrderRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 申请退款
	err := h.orderService.ApplyRefund(c.Request.Context(), id, &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "申请退款失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// ConfirmRefund godoc
// @Summary 确认退款
// @Description 确认订单退款
// @Tags 订单管理
// @Accept json
// @Produce json
// @Param id path string true "订单ID"
// @Success 200 {object} response.Response
// @Router /api/v1/orders/{id}/refund/confirm [post]
func (h *OrderHandler) ConfirmRefund(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "订单ID不能为空")
		return
	}

	// 确认退款
	err := h.orderService.ConfirmRefund(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "确认退款失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}
