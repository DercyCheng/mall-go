package model

import (
	"errors"
	"time"
)

// Brand 品牌实体
type Brand struct {
	ID                  string
	Name                string
	Logo                string
	Description         string
	FirstLetter         string
	Sort                int
	FactoryStatus       int    // 是否为品牌制造商：0->不是；1->是
	ShowStatus          int    // 是否显示：0->不显示；1->显示
	ProductCount        int    // 产品数量
	ProductCommentCount int    // 产品评价数量
	BigPic              string // 专区大图
	BrandStory          string // 品牌故事
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// NewBrand 创建新品牌的工厂方法
func NewBrand(name string, logo string, description string) (*Brand, error) {
	if name == "" {
		return nil, errors.New("brand name cannot be empty")
	}

	return &Brand{
		Name:          name,
		Logo:          logo,
		Description:   description,
		ShowStatus:    1,
		FactoryStatus: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// Update 更新品牌信息
func (b *Brand) Update(name, logo, description string) {
	if name != "" {
		b.Name = name
	}

	if logo != "" {
		b.Logo = logo
	}

	if description != "" {
		b.Description = description
	}

	b.UpdatedAt = time.Now()
}

// Show 显示品牌
func (b *Brand) Show() {
	b.ShowStatus = 1
	b.UpdatedAt = time.Now()
}

// Hide 隐藏品牌
func (b *Brand) Hide() {
	b.ShowStatus = 0
	b.UpdatedAt = time.Now()
}

// SetFactoryStatus 设置品牌制造商状态
func (b *Brand) SetFactoryStatus(isFactory bool) {
	if isFactory {
		b.FactoryStatus = 1
	} else {
		b.FactoryStatus = 0
	}
	b.UpdatedAt = time.Now()
}
