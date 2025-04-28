package repository

import (
	"context"

	"mall-go/services/product-service/domain/model"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	FindByID(ctx context.Context, id string) (*model.Category, error)
	FindByParentID(ctx context.Context, parentID string) ([]*model.Category, error)
	FindAll(ctx context.Context, page, size int) ([]*model.Category, int64, error)
	Save(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
}