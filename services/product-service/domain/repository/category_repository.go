package repository

import (
	"context"

	"mall-go/services/product-service/domain/model"
)

// CategoryRepository 定义了分类实体的仓储接口
type CategoryRepository interface {
	// 保存分类
	Save(ctx context.Context, category *model.Category) error
	// 根据ID查找分类
	FindByID(ctx context.Context, id string) (*model.Category, error)
	// 根据父ID查找子分类
	FindByParentID(ctx context.Context, parentID string) ([]*model.Category, error)
	// 获取所有分类
	FindAll(ctx context.Context) ([]*model.Category, error)
	// 查询分类树形结构
	FindTree(ctx context.Context) ([]*model.Category, error)
	// 更新分类
	Update(ctx context.Context, category *model.Category) error
	// 删除分类
	Delete(ctx context.Context, id string) error
	// 批量更新导航栏显示状态
	UpdateNavStatus(ctx context.Context, ids []string, navStatus int) error
	// 批量更新显示状态
	UpdateShowStatus(ctx context.Context, ids []string, showStatus int) error
	// 查询分类及其子分类
	FindWithChildren(ctx context.Context, id string) (*model.Category, []*model.Category, error)
}
