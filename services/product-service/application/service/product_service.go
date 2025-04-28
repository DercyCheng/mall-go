package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"mall-go/services/product-service/application/dto"
	"mall-go/services/product-service/domain/model"
	"mall-go/services/product-service/domain/repository"
)

// ProductService 商品应用服务
type ProductService struct {
	productRepo repository.ProductRepository
	brandRepo   repository.BrandRepository
	categoryRepo repository.CategoryRepository
}

// NewProductService 创建商品应用服务实例
func NewProductService(
	productRepo repository.ProductRepository,
	brandRepo repository.BrandRepository,
	categoryRepo repository.CategoryRepository,
) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		brandRepo:   brandRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateProduct 创建新商品
func (s *ProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// 获取品牌信息
	brand, err := s.brandRepo.FindByID(ctx, req.BrandID)
	if err != nil {
		return nil, err
	}

	// 获取分类信息
	category, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}

	// 创建商品领域模型
	product := &model.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price: model.Money{
			Amount:   req.Price,
			Currency: "CNY", // 默认使用人民币
		},
		Status: model.ProductStatusDraft,
		Inventory: model.Inventory{
			AvailableQuantity: req.StockQuantity,
			ReservedQuantity:  0,
			LowStockThreshold: req.LowStockThreshold,
		},
		Brand:       *brand,
		Category:    *category,
		Attributes:  req.Attributes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 处理促销信息
	if req.PromotionInfo != nil {
		promotion := &model.Promotion{
			ID:           uuid.New().String(),
			Type:         model.PromotionType(req.PromotionInfo.Type),
			Discount:     req.PromotionInfo.Discount,
			StartTime:    req.PromotionInfo.StartTime,
			EndTime:      req.PromotionInfo.EndTime,
			Requirements: req.PromotionInfo.Requirements,
		}
		product.Promotion = promotion
	}

	// 保存商品
	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return s.toProductResponse(product), nil
}

// GetProduct 获取商品详情
func (s *ProductService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// UpdateProduct 更新商品信息
func (s *ProductService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// 获取原商品信息
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新商品信息
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price.Amount = req.Price
	}
	if req.StockQuantity >= 0 {
		product.Inventory.AvailableQuantity = req.StockQuantity
	}
	if req.LowStockThreshold >= 0 {
		product.Inventory.LowStockThreshold = req.LowStockThreshold
	}
	if len(req.Attributes) > 0 {
		product.Attributes = req.Attributes
	}

	// 更新品牌信息
	if req.BrandID != "" && req.BrandID != product.Brand.ID {
		brand, err := s.brandRepo.FindByID(ctx, req.BrandID)
		if err != nil {
			return nil, err
		}
		product.Brand = *brand
	}

	// 更新分类信息
	if req.CategoryID != "" && req.CategoryID != product.Category.ID {
		category, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
		if err != nil {
			return nil, err
		}
		product.Category = *category
	}

	// 更新促销信息
	if req.PromotionInfo != nil {
		promotion := &model.Promotion{
			ID:           uuid.New().String(),
			Type:         model.PromotionType(req.PromotionInfo.Type),
			Discount:     req.PromotionInfo.Discount,
			StartTime:    req.PromotionInfo.StartTime,
			EndTime:      req.PromotionInfo.EndTime,
			Requirements: req.PromotionInfo.Requirements,
		}

		if err := product.ApplyPromotion(promotion); err != nil {
			return nil, err
		}
	}

	// 更新时间
	product.UpdatedAt = time.Now()

	// 保存更新
	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	return s.toProductResponse(product), nil
}

// DeleteProduct 删除商品
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.productRepo.Delete(ctx, id)
}

// UpdateProductStatus 更新商品状态
func (s *ProductService) UpdateProductStatus(ctx context.Context, id string, req dto.UpdateProductStatusRequest) error {
	status := model.ProductStatus(req.Status)
	return s.productRepo.UpdateStatus(ctx, id, status)
}

// PublishProduct 发布商品
func (s *ProductService) PublishProduct(ctx context.Context, id string) error {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := product.Publish(); err != nil {
		return err
	}

	return s.productRepo.Update(ctx, product)
}

// SearchProducts 搜索商品
func (s *ProductService) SearchProducts(ctx context.Context, req dto.SearchProductRequest) (*dto.ProductListResponse, error) {
	products, total, err := s.productRepo.Search(ctx, req.Query, req.Filters, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, *s.toProductResponse(product))
	}

	return &dto.ProductListResponse{
		List:     productResponses,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ListProducts 获取商品列表
func (s *ProductService) ListProducts(ctx context.Context, page, pageSize int) (*dto.ProductListResponse, error) {
	products, total, err := s.productRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var productResponses []dto.ProductResponse
	for _, product := range products {
		productResponses = append(productResponses, *s.toProductResponse(product))
	}

	return &dto.ProductListResponse{
		List:     productResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// 辅助方法: 将领域模型转换为响应DTO
func (s *ProductService) toProductResponse(product *model.Product) *dto.ProductResponse {
	resp := &dto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price.Amount,
		Status:      string(product.Status),
		Brand: dto.BrandResponse{
			ID:   product.Brand.ID,
			Name: product.Brand.Name,
			Logo: product.Brand.Logo,
		},
		Category: dto.CategoryResponse{
			ID:       product.Category.ID,
			Name:     product.Category.Name,
			ParentID: product.Category.ParentID,
			Level:    product.Category.Level,
		},
		Inventory: dto.InventoryResponse{
			AvailableQuantity: product.Inventory.AvailableQuantity,
			LowStockThreshold: product.Inventory.LowStockThreshold,
		},
		Attributes: product.Attributes,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}

	// 处理可选的促销信息
	if product.Promotion != nil {
		resp.Promotion = &dto.PromotionResponse{
			ID:           product.Promotion.ID,
			Type:         string(product.Promotion.Type),
			Discount:     product.Promotion.Discount,
			StartTime:    product.Promotion.StartTime,
			EndTime:      product.Promotion.EndTime,
			Requirements: product.Promotion.Requirements,
		}
	}

	return resp
}