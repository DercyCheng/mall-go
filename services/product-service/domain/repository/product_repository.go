package repository

import (
	"context"

	"mall-go/services/product-service/domain/model"
)

// ProductRepository 定义了商品聚合根的仓储接口
type ProductRepository interface {
	// 保存商品
	Save(ctx context.Context, product *model.Product) error
	// 根据ID查找商品
	FindByID(ctx context.Context, id string) (*model.Product, error)
	// 根据分类查找商品
	FindByCategory(ctx context.Context, categoryID string, page, size int) ([]*model.Product, int64, error)
	// 根据品牌查找商品
	FindByBrand(ctx context.Context, brandID string, page, size int) ([]*model.Product, int64, error)
	// 搜索商品
	Search(ctx context.Context, query string, filters map[string]interface{}, page, size int) ([]*model.Product, int64, error)
	// 更新商品
	Update(ctx context.Context, product *model.Product) error
	// 删除商品
	Delete(ctx context.Context, id string) error
	// 分页查询商品列表
	List(ctx context.Context, pageNum, pageSize int, name, productSn string, publishStatus, verifyStatus int, brandId, productCategoryId string) ([]*model.Product, int64, error)
	// 批量更新上架状态
	UpdatePublishStatus(ctx context.Context, ids []string, publishStatus int) error
	// 批量更新推荐状态
	UpdateRecommendStatus(ctx context.Context, ids []string, recommendStatus int) error
	// 批量更新新品状态
	UpdateNewStatus(ctx context.Context, ids []string, newStatus int) error
	// 批量删除
	DeleteBatch(ctx context.Context, ids []string) error
}
