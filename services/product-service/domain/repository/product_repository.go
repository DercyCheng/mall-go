package repository

import (
	"context"

	"mall-go/services/product-service/domain/model"
)

// ProductRepository 定义了商品聚合根的仓储接口
type ProductRepository interface {
	Save(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id string) (*model.Product, error)
	FindByCategory(ctx context.Context, categoryID string, page, size int) ([]*model.Product, int64, error)
	FindByBrand(ctx context.Context, brandID string, page, size int) ([]*model.Product, int64, error)
	Search(ctx context.Context, query string, filters map[string]interface{}, page, size int) ([]*model.Product, int64, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status model.ProductStatus) error
	FindAll(ctx context.Context, page, size int) ([]*model.Product, int64, error)
}