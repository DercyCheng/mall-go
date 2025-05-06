package dto

import "time"

// ProductResponse 商品响应DTO
type ProductResponse struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	SubTitle            string         `json:"subTitle"`
	ProductSn           string         `json:"productSn"`
	Description         string         `json:"description"`
	Price               float64        `json:"price"`
	OriginalPrice       float64        `json:"originalPrice"`
	Stock               int            `json:"stock"`
	LowStock            int            `json:"lowStock"`
	Unit                string         `json:"unit"`
	Weight              float64        `json:"weight"`
	Sort                int            `json:"sort"`
	Sale                int            `json:"sale"`
	BrandID             string         `json:"brandId"`
	BrandName           string         `json:"brandName"`
	ProductCategoryID   string         `json:"productCategoryId"`
	ProductCategoryName string         `json:"productCategoryName"`
	Pic                 string         `json:"pic"`
	AlbumPics           []string       `json:"albumPics"`
	DetailTitle         string         `json:"detailTitle"`
	DetailDesc          string         `json:"detailDesc"`
	DetailHTML          string         `json:"detailHtml"`
	DetailMobileHTML    string         `json:"detailMobileHtml"`
	PromotionStartTime  time.Time      `json:"promotionStartTime"`
	PromotionEndTime    time.Time      `json:"promotionEndTime"`
	PromotionPerLimit   int            `json:"promotionPerLimit"`
	PromotionType       int            `json:"promotionType"`
	Keywords            string         `json:"keywords"`
	Note                string         `json:"note"`
	ServiceIDs          []string       `json:"serviceIds"`
	PublishStatus       int            `json:"publishStatus"`
	NewStatus           int            `json:"newStatus"`
	RecommendStatus     int            `json:"recommendStatus"`
	VerifyStatus        int            `json:"verifyStatus"`
	PreviewStatus       int            `json:"previewStatus"`
	DeleteStatus        int            `json:"deleteStatus"`
	Attributes          []AttributeDTO `json:"attributes"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
}

// ProductBriefResponse 商品简要响应DTO
type ProductBriefResponse struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	SubTitle            string    `json:"subTitle"`
	ProductSn           string    `json:"productSn"`
	Price               float64   `json:"price"`
	OriginalPrice       float64   `json:"originalPrice"`
	Pic                 string    `json:"pic"`
	Sale                int       `json:"sale"`
	BrandName           string    `json:"brandName"`
	ProductCategoryName string    `json:"productCategoryName"`
	PublishStatus       int       `json:"publishStatus"`
	NewStatus           int       `json:"newStatus"`
	RecommendStatus     int       `json:"recommendStatus"`
	CreatedAt           time.Time `json:"createdAt"`
}

// ProductListResponse 商品列表响应DTO
type ProductListResponse struct {
	List  []ProductBriefResponse `json:"list"`
	Total int64                  `json:"total"`
}

// BrandResponse 品牌响应DTO
type BrandResponse struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Logo                string    `json:"logo"`
	Description         string    `json:"description"`
	FirstLetter         string    `json:"firstLetter"`
	Sort                int       `json:"sort"`
	FactoryStatus       int       `json:"factoryStatus"`
	ShowStatus          int       `json:"showStatus"`
	ProductCount        int       `json:"productCount"`
	ProductCommentCount int       `json:"productCommentCount"`
	BigPic              string    `json:"bigPic"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// CategoryResponse 分类响应DTO
type CategoryResponse struct {
	ID           string             `json:"id"`
	ParentID     string             `json:"parentId"`
	Name         string             `json:"name"`
	Level        int                `json:"level"`
	ProductCount int                `json:"productCount"`
	ProductUnit  string             `json:"productUnit"`
	NavStatus    int                `json:"navStatus"`
	ShowStatus   int                `json:"showStatus"`
	Sort         int                `json:"sort"`
	Icon         string             `json:"icon"`
	Keywords     string             `json:"keywords"`
	Description  string             `json:"description"`
	Children     []CategoryResponse `json:"children,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`
}
