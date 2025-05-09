package dto

import "time"

// ProductCreateRequest 创建商品的请求DTO
type ProductCreateRequest struct {
	Name               string         `json:"name" binding:"required"`
	SubTitle           string         `json:"subTitle"`
	ProductSn          string         `json:"productSn" binding:"required"`
	Description        string         `json:"description"`
	Price              float64        `json:"price" binding:"required"`
	OriginalPrice      float64        `json:"originalPrice"`
	Stock              int            `json:"stock" binding:"required"`
	Unit               string         `json:"unit"`
	Weight             float64        `json:"weight"`
	Sort               int            `json:"sort"`
	BrandID            string         `json:"brandId" binding:"required"`
	ProductCategoryID  string         `json:"productCategoryId" binding:"required"`
	Pic                string         `json:"pic"`
	AlbumPics          []string       `json:"albumPics"`
	DetailTitle        string         `json:"detailTitle"`
	DetailDesc         string         `json:"detailDesc"`
	DetailHTML         string         `json:"detailHtml"`
	DetailMobileHTML   string         `json:"detailMobileHtml"`
	PromotionStartTime time.Time      `json:"promotionStartTime"`
	PromotionEndTime   time.Time      `json:"promotionEndTime"`
	PromotionPerLimit  int            `json:"promotionPerLimit"`
	PromotionType      int            `json:"promotionType"`
	Keywords           string         `json:"keywords"`
	Note               string         `json:"note"`
	ServiceIDs         []string       `json:"serviceIds"`
	Attributes         []AttributeDTO `json:"attributes"`
}

// ProductUpdateRequest 更新商品的请求DTO
type ProductUpdateRequest struct {
	ID                 string         `json:"id" binding:"required"`
	Name               string         `json:"name"`
	SubTitle           string         `json:"subTitle"`
	Description        string         `json:"description"`
	Price              float64        `json:"price"`
	OriginalPrice      float64        `json:"originalPrice"`
	Stock              int            `json:"stock"`
	Unit               string         `json:"unit"`
	Weight             float64        `json:"weight"`
	Sort               int            `json:"sort"`
	BrandID            string         `json:"brandId"`
	ProductCategoryID  string         `json:"productCategoryId"`
	Pic                string         `json:"pic"`
	AlbumPics          []string       `json:"albumPics"`
	DetailTitle        string         `json:"detailTitle"`
	DetailDesc         string         `json:"detailDesc"`
	DetailHTML         string         `json:"detailHtml"`
	DetailMobileHTML   string         `json:"detailMobileHtml"`
	PromotionStartTime time.Time      `json:"promotionStartTime"`
	PromotionEndTime   time.Time      `json:"promotionEndTime"`
	PromotionPerLimit  int            `json:"promotionPerLimit"`
	PromotionType      int            `json:"promotionType"`
	Keywords           string         `json:"keywords"`
	Note               string         `json:"note"`
	ServiceIDs         []string       `json:"serviceIds"`
	Attributes         []AttributeDTO `json:"attributes"`
}

// ProductQueryRequest 查询商品的请求DTO
type ProductQueryRequest struct {
	PageNum           int     `form:"pageNum" binding:"required,min=1"`
	PageSize          int     `form:"pageSize" binding:"required,min=1,max=100"`
	Name              string  `form:"name"`
	ProductSn         string  `form:"productSn"`
	PublishStatus     int     `form:"publishStatus" binding:"omitempty,oneof=-1 0 1"`
	VerifyStatus      int     `form:"verifyStatus" binding:"omitempty,oneof=-1 0 1"`
	BrandID           string  `form:"brandId"`
	ProductCategoryID string  `form:"productCategoryId"`
	CategoryID        string  `form:"categoryId"` // 添加以修复错误
	Keyword           string  `form:"keyword"`    // 添加以修复错误
	MinPrice          float64 `form:"minPrice"`   // 添加以修复错误
	MaxPrice          float64 `form:"maxPrice"`   // 添加以修复错误
	Page              int     `form:"page"`       // 添加以修复错误
	Size              int     `form:"size"`       // 添加以修复错误
	SortBy            string  `form:"sortBy"`     // 添加以修复错误
	SortDesc          bool    `form:"sortDesc"`   // 添加以修复错误
}

// ProductStatusBatchRequest 批量更新商品状态的请求DTO
type ProductStatusBatchRequest struct {
	IDs    []string `json:"ids" binding:"required,min=1"`
	Status int      `json:"status" binding:"required,oneof=0 1"`
}

// AttributeDTO 商品属性DTO
type AttributeDTO struct {
	ID    string `json:"id"` // 添加ID字段
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
	Type  int    `json:"type"` // 添加Type字段
	Sort  int    `json:"sort"` // 添加Sort字段
}

// CategoryTransferRequest 商品分类转移的请求DTO
type CategoryTransferRequest struct {
	ProductIDs []string `json:"productIds" binding:"required,min=1"`
	CategoryID string   `json:"categoryId" binding:"required"`
}

// 添加缺失的DTO类型

// ProductAttributeRequest 产品属性请求DTO
type ProductAttributeRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  int    `json:"type"`
	Sort  int    `json:"sort"`
}

// ProductSkuRequest 产品SKU请求DTO
type ProductSkuRequest struct {
	SkuCode  string  `json:"skuCode"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	LowStock int     `json:"lowStock"`
	SpecJSON string  `json:"specJson"`
	Pic      string  `json:"pic"`
}

// CreateProductRequest 创建产品请求DTO
type CreateProductRequest struct {
	Name          string                    `json:"name"`
	Subtitle      string                    `json:"subtitle"`
	BrandID       string                    `json:"brandId"`
	CategoryID    string                    `json:"categoryId"`
	ProductSn     string                    `json:"productSn"`
	Price         float64                   `json:"price"`
	OriginalPrice float64                   `json:"originalPrice"`
	Stock         int                       `json:"stock"`
	LowStock      int                       `json:"lowStock"`
	Unit          int                       `json:"unit"`
	UnitName      string                    `json:"unitName"`
	Weight        float64                   `json:"weight"`
	Sort          int                       `json:"sort"`
	PicUrls       []string                  `json:"picUrls"`
	AlbumPics     []string                  `json:"albumPics"`
	MainPic       string                    `json:"mainPic"`
	DetailTitle   string                    `json:"detailTitle"`
	DetailDesc    string                    `json:"detailDesc"`
	Description   string                    `json:"description"`
	Keywords      string                    `json:"keywords"`
	PromotionType int                       `json:"promotionType"`
	PublishStatus int                       `json:"publishStatus"`
	Attributes    []ProductAttributeRequest `json:"attributes"`
	Skus          []ProductSkuRequest       `json:"skus"`
}

// UpdateProductRequest 更新产品请求DTO
type UpdateProductRequest struct {
	Name          string                    `json:"name"`
	Subtitle      string                    `json:"subtitle"`
	BrandID       string                    `json:"brandId"`
	CategoryID    string                    `json:"categoryId"`
	ProductSn     string                    `json:"productSn"`
	Price         float64                   `json:"price"`
	OriginalPrice float64                   `json:"originalPrice"`
	Stock         int                       `json:"stock"`
	LowStock      int                       `json:"lowStock"`
	Unit          int                       `json:"unit"`
	UnitName      string                    `json:"unitName"`
	Weight        float64                   `json:"weight"`
	Sort          int                       `json:"sort"`
	PicUrls       []string                  `json:"picUrls"`
	AlbumPics     []string                  `json:"albumPics"`
	MainPic       string                    `json:"mainPic"`
	DetailTitle   string                    `json:"detailTitle"`
	DetailDesc    string                    `json:"detailDesc"`
	Description   string                    `json:"description"`
	Keywords      string                    `json:"keywords"`
	PromotionType int                       `json:"promotionType"`
	PublishStatus int                       `json:"publishStatus"`
	Attributes    []ProductAttributeRequest `json:"attributes"`
	Skus          []ProductSkuRequest       `json:"skus"`
}
