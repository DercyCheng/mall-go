package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/response"
	"mall-go/services/product-service/application/dto"
	"mall-go/services/product-service/application/service"
)

// ProductHandler 商品API处理器
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler 创建商品API处理器
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// Create 创建商品
// @Summary 创建商品
// @Description 创建新商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "商品信息"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	productResp, err := h.productService.CreateProduct(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "创建商品失败: "+err.Error())
		return
	}

	response.Success(c, productResp)
}

// Get 获取商品详情
// @Summary 获取商品详情
// @Description 根据商品ID获取商品详细信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "商品ID不能为空")
		return
	}

	productResp, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "获取商品失败: "+err.Error())
		return
	}

	response.Success(c, productResp)
}

// Update 更新商品
// @Summary 更新商品
// @Description 根据商品ID更新商品信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Param product body dto.UpdateProductRequest true "商品信息"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "商品ID不能为空")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	productResp, err := h.productService.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新商品失败: "+err.Error())
		return
	}

	response.Success(c, productResp)
}

// Delete 删除商品
// @Summary 删除商品
// @Description 根据商品ID删除商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "商品ID不能为空")
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除商品失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateStatus 更新商品状态
// @Summary 更新商品状态
// @Description 根据商品ID更新商品状态
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Param status body dto.UpdateProductStatusRequest true "商品状态"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id}/status [patch]
func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "商品ID不能为空")
		return
	}

	var req dto.UpdateProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	if err := h.productService.UpdateProductStatus(c.Request.Context(), id, req); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新商品状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Publish 发布商品
// @Summary 发布商品
// @Description 根据商品ID发布商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/{id}/publish [post]
func (h *ProductHandler) Publish(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "商品ID不能为空")
		return
	}

	if err := h.productService.PublishProduct(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "发布商品失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Search 搜索商品
// @Summary 搜索商品
// @Description 根据查询条件搜索商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param query query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.ProductListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products/search [get]
func (h *ProductHandler) Search(c *gin.Context) {
	var req dto.SearchProductRequest
	req.Query = c.Query("query")
	req.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	req.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 校验分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	result, err := h.productService.SearchProducts(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "搜索商品失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// List 获取商品列表
// @Summary 获取商品列表
// @Description 分页获取商品列表
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.ProductListResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 校验分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	result, err := h.productService.ListProducts(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取商品列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}