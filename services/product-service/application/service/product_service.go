package service

import (
	"context"

	"mall-go/services/product-service/application/dto"
	"mall-go/services/product-service/domain/event"
	"mall-go/services/product-service/domain/model"
	"mall-go/services/product-service/domain/repository"
	domainService "mall-go/services/product-service/domain/service"
)

// EventPublisher 事件发布接口
type EventPublisher interface {
	Publish(ctx context.Context, event event.Event) error
}

// ProductService 商品应用服务
type ProductService struct {
	productRepository  repository.ProductRepository
	brandRepository    repository.BrandRepository
	categoryRepository repository.CategoryRepository
	domainService      *domainService.ProductDomainService
	eventPublisher     EventPublisher
}

// NewProductService 创建商品应用服务
func NewProductService(
	productRepository repository.ProductRepository,
	brandRepository repository.BrandRepository,
	categoryRepository repository.CategoryRepository,
	domainService *domainService.ProductDomainService,
	eventPublisher EventPublisher,
) *ProductService {
	return &ProductService{
		productRepository:  productRepository,
		brandRepository:    brandRepository,
		categoryRepository: categoryRepository,
		domainService:      domainService,
		eventPublisher:     eventPublisher,
	}
}

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(ctx context.Context, req *dto.ProductCreateRequest) (string, error) {
	// 创建商品领域模型
	product, err := model.NewProduct(req.Name, req.Price, req.Description, req.BrandID, req.ProductCategoryID, req.Stock)
	if err != nil {
		return "", err
	}

	// 补充其他信息
	product.ProductSn = req.ProductSn
	product.SubTitle = req.SubTitle
	product.Unit = req.Unit
	product.Weight = req.Weight
	product.Sort = req.Sort
	product.OriginalPrice = model.Money{
		Amount:   req.OriginalPrice,
		Currency: "CNY",
	}
	product.Pic = req.Pic
	product.DetailTitle = req.DetailTitle
	product.DetailDesc = req.DetailDesc
	product.DetailHTML = req.DetailHTML
	product.DetailMobileHTML = req.DetailMobileHTML
	product.Keywords = req.Keywords
	product.Note = req.Note
	product.AlbumPics = req.AlbumPics
	product.ServiceIDs = req.ServiceIDs

	if !req.PromotionStartTime.IsZero() && !req.PromotionEndTime.IsZero() {
		product.PromotionStartTime = req.PromotionStartTime
		product.PromotionEndTime = req.PromotionEndTime
		product.PromotionPerLimit = req.PromotionPerLimit
		product.PromotionType = req.PromotionType
	}

	// 添加属性
	for _, attr := range req.Attributes {
		product.AddAttribute(attr.Name, attr.Value)
	}

	// 使用领域服务创建商品（处理跨聚合根的业务逻辑）
	if err := s.domainService.CreateProduct(ctx, product); err != nil {
		return "", err
	}

	// 发布领域事件
	createEvent := event.NewProductCreatedEvent(product.ID, product.Name, product.Price.Amount)
	if err := s.eventPublisher.Publish(ctx, createEvent); err != nil {
		// 通常只记录日志，不影响主流程
		// log.Printf("Failed to publish event: %v", err)
	}

	return product.ID, nil
}

// UpdateProduct 更新商品
func (s *ProductService) UpdateProduct(ctx context.Context, req *dto.ProductUpdateRequest) error {
	// 获取现有商品
	product, err := s.productRepository.FindByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	// 更新基础信息
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.SubTitle != "" {
		product.SubTitle = req.SubTitle
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		oldPrice := product.Price.Amount
		product.Price.Amount = req.Price

		// 发布价格变更事件
		priceEvent := event.NewProductPriceChangedEvent(product.ID, oldPrice, req.Price)
		_ = s.eventPublisher.Publish(ctx, priceEvent)
	}
	if req.OriginalPrice > 0 {
		product.OriginalPrice.Amount = req.OriginalPrice
	}
	if req.Stock > 0 {
		oldStock := product.Inventory.AvailableQuantity
		product.Inventory.AvailableQuantity = req.Stock

		// 发布库存变更事件
		stockEvent := event.NewProductStockChangedEvent(product.ID, oldStock, req.Stock)
		_ = s.eventPublisher.Publish(ctx, stockEvent)
	}

	// 更新其他信息
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.Weight > 0 {
		product.Weight = req.Weight
	}
	if req.Sort > 0 {
		product.Sort = req.Sort
	}
	if req.Pic != "" {
		product.Pic = req.Pic
	}
	if len(req.AlbumPics) > 0 {
		product.AlbumPics = req.AlbumPics
	}
	if req.DetailTitle != "" {
		product.DetailTitle = req.DetailTitle
	}
	if req.DetailDesc != "" {
		product.DetailDesc = req.DetailDesc
	}
	if req.DetailHTML != "" {
		product.DetailHTML = req.DetailHTML
	}
	if req.DetailMobileHTML != "" {
		product.DetailMobileHTML = req.DetailMobileHTML
	}
	if req.Keywords != "" {
		product.Keywords = req.Keywords
	}
	if req.Note != "" {
		product.Note = req.Note
	}
	if len(req.ServiceIDs) > 0 {
		product.ServiceIDs = req.ServiceIDs
	}

	// 更新促销信息
	if !req.PromotionStartTime.IsZero() && !req.PromotionEndTime.IsZero() {
		product.PromotionStartTime = req.PromotionStartTime
		product.PromotionEndTime = req.PromotionEndTime
		product.PromotionPerLimit = req.PromotionPerLimit
		product.PromotionType = req.PromotionType
	}

	// 处理品牌和分类更新
	if req.BrandID != "" && req.BrandID != product.Brand.ID {
		brand, err := s.brandRepository.FindByID(ctx, req.BrandID)
		if err != nil {
			return err
		}
		if brand != nil {
			product.Brand = *brand
		}
	}

	if req.ProductCategoryID != "" && req.ProductCategoryID != product.Category.ID {
		category, err := s.categoryRepository.FindByID(ctx, req.ProductCategoryID)
		if err != nil {
			return err
		}
		if category != nil {
			product.Category = *category
		}
	}

	// 如果有新的属性，添加或更新
	if len(req.Attributes) > 0 {
		// 简单实现：直接替换属性
		// 实际项目中，可能需要更复杂的合并逻辑
		product.Attributes = make([]model.Attribute, len(req.Attributes))
		for i, attr := range req.Attributes {
			product.Attributes[i] = model.Attribute{
				Name:  attr.Name,
				Value: attr.Value,
			}
		}
	}

	// 保存更新
	if err := s.productRepository.Update(ctx, product); err != nil {
		return err
	}

	// 发布更新事件
	updateEvent := event.NewProductUpdatedEvent(product.ID, product.Name, product.Price.Amount)
	_ = s.eventPublisher.Publish(ctx, updateEvent)

	return nil
}

// GetProduct 获取商品详情
func (s *ProductService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := s.productRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// 转换为DTO
	return convertProductToDTO(product), nil
}

// ListProducts 分页查询商品列表
func (s *ProductService) ListProducts(ctx context.Context, req *dto.ProductQueryRequest) (*dto.ProductListResponse, error) {
	products, total, err := s.productRepository.List(
		ctx,
		req.PageNum,
		req.PageSize,
		req.Name,
		req.ProductSn,
		req.PublishStatus,
		req.VerifyStatus,
		req.BrandID,
		req.ProductCategoryID,
	)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	result := &dto.ProductListResponse{
		Total: total,
		List:  make([]dto.ProductBriefResponse, len(products)),
	}

	for i, product := range products {
		result.List[i] = convertProductToBriefDTO(product)
	}

	return result, nil
}

// DeleteProduct 删除商品
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	// 检查商品是否存在
	product, err := s.productRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrProductNotFound
	}

	// 删除商品
	if err := s.productRepository.Delete(ctx, id); err != nil {
		return err
	}

	// 发布删除事件
	deleteEvent := event.NewProductDeletedEvent(id)
	_ = s.eventPublisher.Publish(ctx, deleteEvent)

	return nil
}

// BatchDeleteProducts 批量删除商品
func (s *ProductService) BatchDeleteProducts(ctx context.Context, ids []string) error {
	if err := s.productRepository.DeleteBatch(ctx, ids); err != nil {
		return err
	}

	// 发布批量删除事件
	for _, id := range ids {
		deleteEvent := event.NewProductDeletedEvent(id)
		_ = s.eventPublisher.Publish(ctx, deleteEvent)
	}

	return nil
}

// UpdatePublishStatus 更新上架状态
func (s *ProductService) UpdatePublishStatus(ctx context.Context, req *dto.ProductStatusBatchRequest) error {
	err := s.productRepository.UpdatePublishStatus(ctx, req.IDs, req.Status)
	if err != nil {
		return err
	}

	// 发布状态变更事件
	if req.Status == 1 { // 上架
		for _, id := range req.IDs {
			product, _ := s.productRepository.FindByID(ctx, id)
			if product != nil {
				publishEvent := event.NewProductPublishedEvent(id, product.Name)
				_ = s.eventPublisher.Publish(ctx, publishEvent)
			}
		}
	}

	return nil
}

// UpdateNewStatus 更新新品状态
func (s *ProductService) UpdateNewStatus(ctx context.Context, req *dto.ProductStatusBatchRequest) error {
	return s.productRepository.UpdateNewStatus(ctx, req.IDs, req.Status)
}

// UpdateRecommendStatus 更新推荐状态
func (s *ProductService) UpdateRecommendStatus(ctx context.Context, req *dto.ProductStatusBatchRequest) error {
	return s.productRepository.UpdateRecommendStatus(ctx, req.IDs, req.Status)
}

// TransferCategory 转移商品分类
func (s *ProductService) TransferCategory(ctx context.Context, req *dto.CategoryTransferRequest) error {
	for _, productID := range req.ProductIDs {
		if err := s.domainService.TransferCategory(ctx, productID, req.CategoryID); err != nil {
			return err
		}
	}
	return nil
}

// 工具函数：商品模型转DTO
func convertProductToDTO(product *model.Product) *dto.ProductResponse {
	response := &dto.ProductResponse{
		ID:                  product.ID,
		Name:                product.Name,
		SubTitle:            product.SubTitle,
		ProductSn:           product.ProductSn,
		Description:         product.Description,
		Price:               product.Price.Amount,
		OriginalPrice:       product.OriginalPrice.Amount,
		Stock:               product.Inventory.AvailableQuantity,
		LowStock:            product.Inventory.LowStockThreshold,
		Unit:                product.Unit,
		Weight:              product.Weight,
		Sort:                product.Sort,
		Sale:                product.Sale,
		BrandID:             product.Brand.ID,
		BrandName:           product.Brand.Name,
		ProductCategoryID:   product.Category.ID,
		ProductCategoryName: product.Category.Name,
		Pic:                 product.Pic,
		AlbumPics:           product.AlbumPics,
		DetailTitle:         product.DetailTitle,
		DetailDesc:          product.DetailDesc,
		DetailHTML:          product.DetailHTML,
		DetailMobileHTML:    product.DetailMobileHTML,
		PromotionStartTime:  product.PromotionStartTime,
		PromotionEndTime:    product.PromotionEndTime,
		PromotionPerLimit:   product.PromotionPerLimit,
		PromotionType:       product.PromotionType,
		Keywords:            product.Keywords,
		Note:                product.Note,
		ServiceIDs:          product.ServiceIDs,
		PublishStatus:       product.PublishStatus,
		NewStatus:           product.NewStatus,
		RecommendStatus:     product.RecommendStatus,
		VerifyStatus:        product.VerifyStatus,
		PreviewStatus:       product.PreviewStatus,
		DeleteStatus:        product.DeleteStatus,
		CreatedAt:           product.CreatedAt,
		UpdatedAt:           product.UpdatedAt,
	}

	// 转换属性
	response.Attributes = make([]dto.AttributeDTO, len(product.Attributes))
	for i, attr := range product.Attributes {
		response.Attributes[i] = dto.AttributeDTO{
			Name:  attr.Name,
			Value: attr.Value,
		}
	}

	return response
}

// 工具函数：商品模型转简要DTO
func convertProductToBriefDTO(product *model.Product) dto.ProductBriefResponse {
	return dto.ProductBriefResponse{
		ID:                  product.ID,
		Name:                product.Name,
		SubTitle:            product.SubTitle,
		ProductSn:           product.ProductSn,
		Price:               product.Price.Amount,
		OriginalPrice:       product.OriginalPrice.Amount,
		Pic:                 product.Pic,
		Sale:                product.Sale,
		BrandName:           product.Brand.Name,
		ProductCategoryName: product.Category.Name,
		PublishStatus:       product.PublishStatus,
		NewStatus:           product.NewStatus,
		RecommendStatus:     product.RecommendStatus,
		CreatedAt:           product.CreatedAt,
	}
}
