package dto

import (
	"time"

	"mall-go/services/product-service/domain/model"
)

// 请求DTO

// CreateProductRequest 创建商品请求DTO
type CreateProductRequest struct {
	Name             string                 `json:"name" binding:"required"`
	Description      string                 `json:"description"`
	Price            float64                `json:"price" binding:"required,gt=0"`
	BrandID          string                 `json:"brandId" binding:"required"`
	CategoryID       string                 `json:"categoryId" binding:"required"`
	StockQuantity    int                    `json:"stockQuantity" binding:"required,gte=0"`
	LowStockThreshold int                   `json:"lowStockThreshold"`
	Attributes       []model.Attribute      `json:"attributes"`
	PromotionInfo    *PromotionInfo         `json:"promotionInfo"`
}

// PromotionInfo 促销信息DTO
type PromotionInfo struct {
	Type         string                 `json:"type" binding:"required,oneof=percentage fixed_amount bundle"`
	Discount     float64                `json:"discount" binding:"required,gt=0"`
	StartTime    time.Time              `json:"startTime" binding:"required"`
	EndTime      time.Time              `json:"endTime" binding:"required"`
	Requirements map[string]interface{} `json:"requirements"`
}

// UpdateProductRequest 更新商品请求DTO
type UpdateProductRequest struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Price            float64                `json:"price" binding:"omitempty,gt=0"`
	BrandID          string                 `json:"brandId"`
	CategoryID       string                 `json:"categoryId"`
	StockQuantity    int                    `json:"stockQuantity" binding:"omitempty,gte=0"`
	LowStockThreshold int                   `json:"lowStockThreshold" binding:"omitempty,gte=0"`
	Attributes       []model.Attribute      `json:"attributes"`
	PromotionInfo    *PromotionInfo         `json:"promotionInfo"`
}

// UpdateProductStatusRequest 更新商品状态请求DTO
type UpdateProductStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft active inactive deleted"`
}

// SearchProductRequest 搜索商品请求DTO
type SearchProductRequest struct {
	Query    string                 `json:"query" form:"query"`
	Filters  map[string]interface{} `json:"filters" form:"filters"`
	Page     int                    `json:"page" form:"page,default=1" binding:"min=1"`
	PageSize int                    `json:"pageSize" form:"pageSize,default=10" binding:"min=1,max=100"`
}

// 响应DTO

// ProductResponse 商品响应DTO
type ProductResponse struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Price           float64            `json:"price"`
	Status          string             `json:"status"`
	Brand           BrandResponse      `json:"brand"`
	Category        CategoryResponse   `json:"category"`
	Inventory       InventoryResponse  `json:"inventory"`
	Attributes      []model.Attribute  `json:"attributes"`
	Promotion       *PromotionResponse `json:"promotion,omitempty"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

// BrandResponse 品牌响应DTO
type BrandResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// CategoryResponse 分类响应DTO
type CategoryResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parentId"`
	Level    int    `json:"level"`
}

// InventoryResponse 库存响应DTO
type InventoryResponse struct {
	AvailableQuantity int `json:"availableQuantity"`
	LowStockThreshold int `json:"lowStockThreshold"`
}

// PromotionResponse 促销响应DTO
type PromotionResponse struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Discount     float64                `json:"discount"`
	StartTime    time.Time              `json:"startTime"`
	EndTime      time.Time              `json:"endTime"`
	Requirements map[string]interface{} `json:"requirements"`
}

// ProductListResponse 商品列表响应DTO
type ProductListResponse struct {
	List     []ProductResponse `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}