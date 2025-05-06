package service

import (
	"context"
	"errors"

	"mall-go/services/product-service/domain/model"
	"mall-go/services/product-service/domain/repository"
)

// ProductDomainService 商品领域服务
type ProductDomainService struct {
	productRepository  repository.ProductRepository
	brandRepository    repository.BrandRepository
	categoryRepository repository.CategoryRepository
}

// NewProductDomainService 创建商品领域服务
func NewProductDomainService(
	productRepository repository.ProductRepository,
	brandRepository repository.BrandRepository,
	categoryRepository repository.CategoryRepository,
) *ProductDomainService {
	return &ProductDomainService{
		productRepository:  productRepository,
		brandRepository:    brandRepository,
		categoryRepository: categoryRepository,
	}
}

// CreateProduct 创建商品（包含跨聚合根的业务规则）
func (s *ProductDomainService) CreateProduct(ctx context.Context, product *model.Product) error {
	// 验证品牌是否存在
	brand, err := s.brandRepository.FindByID(ctx, product.Brand.ID)
	if err != nil {
		return err
	}
	if brand == nil {
		return errors.New("brand not found")
	}

	// 验证分类是否存在
	category, err := s.categoryRepository.FindByID(ctx, product.Category.ID)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	// 补充品牌和分类的名称
	product.Brand.Name = brand.Name
	product.Category.Name = category.Name

	// 保存商品
	return s.productRepository.Save(ctx, product)
}

// TransferCategory 将商品转移到其他分类
func (s *ProductDomainService) TransferCategory(ctx context.Context, productID string, newCategoryID string) error {
	// 检查产品是否存在
	product, err := s.productRepository.FindByID(ctx, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	// 检查新分类是否存在
	newCategory, err := s.categoryRepository.FindByID(ctx, newCategoryID)
	if err != nil {
		return err
	}
	if newCategory == nil {
		return errors.New("new category not found")
	}

	// 更新商品分类
	product.Category = *newCategory
	return s.productRepository.Update(ctx, product)
}

// ApplyPromotionToCategory 对某个分类的所有商品应用促销
func (s *ProductDomainService) ApplyPromotionToCategory(ctx context.Context, categoryID string, promotion *model.Promotion) error {
	// 分批处理以避免一次性加载过多数据
	pageSize := 100
	pageNum := 1

	for {
		products, total, err := s.productRepository.FindByCategory(ctx, categoryID, pageNum, pageSize)
		if err != nil {
			return err
		}

		for _, product := range products {
			if err := product.ApplyPromotion(promotion); err != nil {
				continue // 如果某个商品应用促销失败，继续处理其他商品
			}
			if err := s.productRepository.Update(ctx, product); err != nil {
				return err
			}
		}

		// 如果已经处理完所有数据，退出循环
		if int64((pageNum-1)*pageSize+len(products)) >= total {
			break
		}

		pageNum++
	}

	return nil
}

// CheckAndUpdateProductStatus 检查并更新商品状态（如库存状态）
func (s *ProductDomainService) CheckAndUpdateProductStatus(ctx context.Context, productID string) error {
	product, err := s.productRepository.FindByID(ctx, productID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	// 检查库存是否低于阈值
	if product.Inventory.AvailableQuantity <= product.Inventory.LowStockThreshold {
		// 这里可以触发领域事件通知库存不足
		// 或者直接更新商品状态
	}

	return s.productRepository.Update(ctx, product)
}
