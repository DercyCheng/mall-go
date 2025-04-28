package repository

import (
	"context"

	"mall-go/services/product-service/domain/model"
)

// BrandRepository 品牌仓储接口
type BrandRepository interface {
	FindByID(ctx context.Context, id string) (*model.Brand, error)
	FindAll(ctx context.Context, page, size int) ([]*model.Brand, int64, error)
	Save(ctx context.Context, brand *model.Brand) error
	Update(ctx context.Context, brand *model.Brand) error
	Delete(ctx context.Context, id string) error
}