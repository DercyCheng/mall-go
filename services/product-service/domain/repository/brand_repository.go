package repository

import (
	"context"

	"mall-go/services/product-service/domain/model"
)

// BrandRepository 定义了品牌实体的仓储接口
type BrandRepository interface {
	// 保存品牌
	Save(ctx context.Context, brand *model.Brand) error
	// 根据ID查找品牌
	FindByID(ctx context.Context, id string) (*model.Brand, error)
	// 根据名称查找品牌
	FindByName(ctx context.Context, name string) (*model.Brand, error)
	// 获取所有品牌
	FindAll(ctx context.Context) ([]*model.Brand, error)
	// 分页查询品牌列表
	List(ctx context.Context, pageNum, pageSize int, name string) ([]*model.Brand, int64, error)
	// 更新品牌
	Update(ctx context.Context, brand *model.Brand) error
	// 删除品牌
	Delete(ctx context.Context, id string) error
	// 批量更新显示状态
	UpdateShowStatus(ctx context.Context, ids []string, showStatus int) error
	// 批量更新厂商状态
	UpdateFactoryStatus(ctx context.Context, ids []string, factoryStatus int) error
}
