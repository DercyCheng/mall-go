package service

import "errors"

// 应用服务层错误定义
var (
	ErrProductNotFound    = errors.New("product not found")
	ErrBrandNotFound      = errors.New("brand not found")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrDuplicateProductSn = errors.New("duplicate product sn")
	ErrInsufficientStock  = errors.New("insufficient stock")
)
