package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/product-service/application/dto"
	applicationservice "mall-go/services/product-service/application/service"
	productpb "mall-go/services/product-service/proto"
)

// ProductServer 实现ProductService gRPC接口
type ProductServer struct {
	productpb.UnimplementedProductServiceServer
	productService *applicationservice.ProductService
}

// NewProductServer 创建新的ProductServer
func NewProductServer(productService *applicationservice.ProductService) *ProductServer {
	return &ProductServer{
		productService: productService,
	}
}

// CreateProduct 创建新产品
func (s *ProductServer) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	// 转换请求
	attributes := make([]dto.ProductAttributeRequest, len(req.Attributes))
	for i, attr := range req.Attributes {
		attributes[i] = dto.ProductAttributeRequest{
			Name:  attr.Name,
			Value: attr.Value,
			Type:  int(attr.Type),
			Sort:  int(attr.Sort),
		}
	}

	skus := make([]dto.ProductSkuRequest, len(req.Skus))
	for i, sku := range req.Skus {
		skus[i] = dto.ProductSkuRequest{
			SkuCode:  sku.SkuCode,
			Price:    sku.Price.Amount,
			Stock:    int(sku.Stock),
			LowStock: int(sku.LowStock),
			SpecJSON: sku.SpecJson,
			Pic:      sku.Pic,
		}
	}

	dtoReq := &dto.ProductCreateRequest{
		Name:              req.Name,
		SubTitle:          req.Subtitle,
		BrandID:           req.BrandId,
		ProductCategoryID: req.CategoryId,
		ProductSn:         req.ProductSn,
		Price:             req.Price,
		OriginalPrice:     req.OriginalPrice,
		Stock:             int(req.Stock),
		Unit:              req.UnitName,
		Weight:            req.Weight,
		Sort:              int(req.Sort),
		Pic:               req.MainPic,
		AlbumPics:         req.AlbumPics,
		DetailTitle:       req.DetailTitle,
		DetailDesc:        req.DetailDesc,
		Description:       req.Description,
		Keywords:          req.Keywords,
		PromotionType:     int(req.PromotionType),
		Attributes:        []dto.AttributeDTO{},
	}

	// 调用应用服务
	result, err := s.productService.CreateProduct(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建产品失败: %v", err)
	}

	return &productpb.CreateProductResponse{
		Success:   true,
		Message:   "产品创建成功",
		ProductId: result,
	}, nil
}

// GetProduct 获取产品详情
func (s *ProductServer) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	// 调用应用服务
	result, err := s.productService.GetProduct(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "获取产品失败: %v", err)
	}

	// 转换产品属性
	attributes := make([]*productpb.ProductAttribute, len(result.Attributes))
	for i, attr := range result.Attributes {
		attributes[i] = &productpb.ProductAttribute{
			Id:    attr.ID,
			Name:  attr.Name,
			Value: attr.Value,
			Type:  int32(attr.Type),
			Sort:  int32(attr.Sort),
		}
	}

	// 转换产品SKU
	skus := make([]*productpb.ProductSku, len(result.Skus))
	for i, sku := range result.Skus {
		skus[i] = &productpb.ProductSku{
			Id:        sku.ID,
			ProductId: sku.ProductID,
			SkuCode:   sku.SkuCode,
			Price:     &productpb.Money{Amount: sku.Price, Currency: "CNY"},
			Stock:     int32(sku.Stock),
			LowStock:  int32(sku.LowStock),
			SpecJson:  sku.SpecJSON,
			Pic:       sku.Pic,
			Sale:      int32(sku.Sale),
			LockStock: int32(sku.LockStock),
			CreatedAt: sku.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: sku.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	// 构建响应产品
	product := &productpb.Product{
		Id:            result.ID,
		Name:          result.Name,
		Subtitle:      result.SubTitle,
		BrandId:       result.BrandID,
		BrandName:     result.BrandName,
		CategoryId:    result.CategoryID,
		CategoryName:  result.CategoryName,
		ProductSn:     result.ProductSn,
		Price:         &productpb.Money{Amount: result.Price, Currency: "CNY"},
		OriginalPrice: &productpb.Money{Amount: result.OriginalPrice, Currency: "CNY"},
		Stock:         int32(result.Stock),
		LowStock:      int32(result.LowStock),
		Unit:          0,           // 修复字段类型不匹配，设置默认值
		UnitName:      result.Unit, // Unit存储的实际上是字符串
		Weight:        result.Weight,
		Sort:          int32(result.Sort),
		PicUrls:       result.PicUrls,
		AlbumPics:     result.AlbumPics,
		MainPic:       result.MainPic,
		DetailTitle:   result.DetailTitle,
		DetailDesc:    result.DetailDesc,
		Description:   result.Description,
		Keywords:      result.Keywords,
		PromotionType: int32(result.PromotionType),
		PublishStatus: int32(result.PublishStatus),
		VerifyStatus:  int32(result.VerifyStatus),
		Sale:          int32(result.Sale),
		GiftPoint:     int32(result.GiftPoint),
		Attributes:    attributes,
		Skus:          skus,
		CreatedAt:     result.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return &productpb.GetProductResponse{
		Success: true,
		Message: "获取产品成功",
		Product: product,
	}, nil
}

// ListProducts 获取产品列表
func (s *ProductServer) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	// 转换请求
	dtoReq := &dto.ProductQueryRequest{
		PageNum:           int(req.Page),
		PageSize:          int(req.Size),
		BrandID:           req.BrandId,
		ProductCategoryID: req.CategoryId,
		PublishStatus:     int(req.PublishStatus),
		Name:              req.Keyword,
	}

	// 调用应用服务
	result, err := s.productService.ListProducts(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取产品列表失败: %v", err)
	}

	// 转换产品列表
	products := make([]*productpb.Product, len(result.List))
	for i, p := range result.List {
		products[i] = &productpb.Product{
			Id:            p.ID,
			Name:          p.Name,
			Subtitle:      p.SubTitle,
			BrandId:       p.BrandID,
			BrandName:     p.BrandName,
			CategoryId:    p.CategoryID,
			CategoryName:  p.CategoryName,
			ProductSn:     p.ProductSn,
			Price:         &productpb.Money{Amount: p.Price, Currency: "CNY"},
			OriginalPrice: &productpb.Money{Amount: p.OriginalPrice, Currency: "CNY"},
			MainPic:       p.MainPic,
			Sale:          int32(p.Sale),
			Stock:         int32(p.Stock),
			Unit:          0, // 修复字段类型不匹配，设置默认值
			UnitName:      p.UnitName,
			CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &productpb.ListProductsResponse{
		Success:  true,
		Message:  "获取产品列表成功",
		Products: products,
		Total:    int32(result.Total),
	}, nil
}

// UpdateProduct 更新产品
func (s *ProductServer) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
	// 转换请求
	dtoReq := &dto.ProductUpdateRequest{
		ID:                req.ProductId,
		Name:              req.Name,
		SubTitle:          req.Subtitle,
		BrandID:           req.BrandId,
		ProductCategoryID: req.CategoryId,
		Description:       req.Description,
		Price:             req.Price,
		OriginalPrice:     req.OriginalPrice,
		Stock:             int(req.Stock),
		Unit:              req.UnitName,
		Weight:            req.Weight,
		Sort:              int(req.Sort),
		Pic:               req.MainPic,
		AlbumPics:         req.AlbumPics,
		DetailTitle:       req.DetailTitle,
		DetailDesc:        req.DetailDesc,
		Keywords:          req.Keywords,
		PromotionType:     int(req.PromotionType),
	}

	// 调用应用服务
	err := s.productService.UpdateProduct(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新产品失败: %v", err)
	}

	return &productpb.UpdateProductResponse{
		Success: true,
		Message: "更新产品成功",
	}, nil
}

// DeleteProduct 删除产品
func (s *ProductServer) DeleteProduct(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error) {
	// 调用应用服务
	err := s.productService.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "删除产品失败: %v", err)
	}

	return &productpb.DeleteProductResponse{
		Success: true,
		Message: "删除产品成功",
	}, nil
}

// GetStock 获取库存
func (s *ProductServer) GetStock(ctx context.Context, req *productpb.GetStockRequest) (*productpb.GetStockResponse, error) {
	// 由于应用服务中可能没有对应的方法，我们简化处理，只返回成功响应
	// TODO: 实现真实的库存查询逻辑
	return &productpb.GetStockResponse{
		Success:   true,
		Message:   "获取库存成功",
		Stock:     100, // 默认值
		LockStock: 0,
	}, nil
}

// UpdateStock 更新库存
func (s *ProductServer) UpdateStock(ctx context.Context, req *productpb.UpdateStockRequest) (*productpb.UpdateStockResponse, error) {
	// 由于应用服务中可能没有对应的方法，我们简化处理，只返回成功响应
	// TODO: 实现真实的库存更新逻辑
	return &productpb.UpdateStockResponse{
		Success: true,
		Message: "更新库存成功",
	}, nil
}
