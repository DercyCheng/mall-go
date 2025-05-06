package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mall-go/pkg/response"
	"mall-go/services/product-service/application/dto"
	"mall-go/services/product-service/application/service"
)

// ProductHandler 商品控制器
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler 创建商品控制器
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// List godoc
// @Summary 获取商品列表
// @Description 根据条件分页查询商品列表
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param pageNum query int true "页码"
// @Param pageSize query int true "每页数量"
// @Param name query string false "商品名称"
// @Param productSn query string false "商品货号"
// @Param publishStatus query int false "上架状态：0->下架；1->上架"
// @Param verifyStatus query int false "审核状态：0->未审核；1->审核通过"
// @Param brandId query string false "品牌ID"
// @Param productCategoryId query string false "商品分类ID"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]dto.ProductResponse}}
// @Router /api/v1/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var query dto.ProductQueryRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	products, err := h.productService.ListProducts(c.Request.Context(), &query)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取商品列表失败: "+err.Error())
		return
	}

	response.PageSuccess(c, products.List, products.Total, query.PageNum, query.PageSize)
}

// Get godoc
// @Summary 获取商品详情
// @Description 根据ID获取商品详情
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Success 200 {object} response.Response{data=dto.ProductResponse}
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "商品ID不能为空")
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取商品详情失败: "+err.Error())
		return
	}

	if product == nil {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}

	response.Success(c, product)
}

// Create godoc
// @Summary 创建商品
// @Description 创建新商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param product body dto.ProductCreateRequest true "商品信息"
// @Success 200 {object} response.Response{data=string}
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	id, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建商品失败: "+err.Error())
		return
	}

	response.Success(c, id)
}

// Update godoc
// @Summary 更新商品
// @Description 更新商品信息
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Param product body dto.ProductUpdateRequest true "商品信息"
// @Success 200 {object} response.Response
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "商品ID不能为空")
		return
	}

	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}
	req.ID = id

	err := h.productService.UpdateProduct(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新商品失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete godoc
// @Summary 删除商品
// @Description 根据ID删除商品
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param id path string true "商品ID"
// @Success 200 {object} response.Response
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "商品ID不能为空")
		return
	}

	err := h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除商品失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdatePublishStatus godoc
// @Summary 批量更新上架状态
// @Description 批量更新商品上架状态
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param request body dto.ProductStatusBatchRequest true "批量更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/products/publish-status [put]
func (h *ProductHandler) UpdatePublishStatus(c *gin.Context) {
	var req dto.ProductStatusBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := h.productService.UpdatePublishStatusDirect(c.Request.Context(), req.IDs, req.Status)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新上架状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateNewStatus godoc
// @Summary 批量更新新品状态
// @Description 批量更新商品新品状态
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param request body dto.ProductStatusBatchRequest true "批量更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/products/new-status [put]
func (h *ProductHandler) UpdateNewStatus(c *gin.Context) {
	var req dto.ProductStatusBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := h.productService.UpdateNewStatusDirect(c.Request.Context(), req.IDs, req.Status)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新新品状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateRecommendStatus godoc
// @Summary 批量更新推荐状态
// @Description 批量更新商品推荐状态
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param request body dto.ProductStatusBatchRequest true "批量更新请求"
// @Success 200 {object} response.Response
// @Router /api/v1/products/recommend-status [put]
func (h *ProductHandler) UpdateRecommendStatus(c *gin.Context) {
	var req dto.ProductStatusBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := h.productService.UpdateRecommendStatusDirect(c.Request.Context(), req.IDs, req.Status)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新推荐状态失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// TransferCategory godoc
// @Summary 转移商品分类
// @Description 将商品转移到其他分类下
// @Tags 商品管理
// @Accept json
// @Produce json
// @Param request body dto.CategoryTransferRequest true "分类转移请求"
// @Success 200 {object} response.Response
// @Router /api/v1/products/category/transfer [post]
func (h *ProductHandler) TransferCategory(c *gin.Context) {
	var req dto.CategoryTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	for _, productID := range req.ProductIDs {
		err := h.productService.TransferCategoryDirect(c.Request.Context(), productID, req.CategoryID)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "转移商品分类失败: "+err.Error())
			return
		}
	}

	response.Success(c, nil)
}