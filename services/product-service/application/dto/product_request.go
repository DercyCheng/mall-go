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
	PageNum           int    `form:"pageNum" binding:"required,min=1"`
	PageSize          int    `form:"pageSize" binding:"required,min=1,max=100"`
	Name              string `form:"name"`
	ProductSn         string `form:"productSn"`
	PublishStatus     int    `form:"publishStatus" binding:"omitempty,oneof=-1 0 1"`
	VerifyStatus      int    `form:"verifyStatus" binding:"omitempty,oneof=-1 0 1"`
	BrandID           string `form:"brandId"`
	ProductCategoryID string `form:"productCategoryId"`
}

// ProductStatusBatchRequest 批量更新商品状态的请求DTO
type ProductStatusBatchRequest struct {
	IDs    []string `json:"ids" binding:"required,min=1"`
	Status int      `json:"status" binding:"required,oneof=0 1"`
}

// AttributeDTO 商品属性DTO
type AttributeDTO struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// CategoryTransferRequest 商品分类转移的请求DTO
type CategoryTransferRequest struct {
	ProductIDs []string `json:"productIds" binding:"required,min=1"`
	CategoryID string   `json:"categoryId" binding:"required"`
}
