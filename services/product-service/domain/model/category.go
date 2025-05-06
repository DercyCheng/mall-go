package model

import (
	"errors"
	"time"
)

// Category 商品分类实体
type Category struct {
	ID           string
	ParentID     string
	Name         string
	Level        int    // 分类级别：0->1级；1->2级；2->3级
	ProductCount int    // 商品数量
	ProductUnit  string // 商品单位
	NavStatus    int    // 是否显示在导航栏：0->不显示；1->显示
	ShowStatus   int    // 显示状态：0->不显示；1->显示
	Sort         int    // 排序
	Icon         string // 图标
	Keywords     string // 关键字
	Description  string // 描述
	Children     []*Category // 子分类
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewCategory 创建新分类的工厂方法
func NewCategory(name string, parentID string, level int) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	return &Category{
		Name:        name,
		ParentID:    parentID,
		Level:       level,
		NavStatus:   0,
		ShowStatus:  1,
		ProductUnit: "件",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// Update 更新分类信息
func (c *Category) Update(name, icon, keywords, description string) {
	if name != "" {
		c.Name = name
	}

	if icon != "" {
		c.Icon = icon
	}

	if keywords != "" {
		c.Keywords = keywords
	}

	if description != "" {
		c.Description = description
	}

	c.UpdatedAt = time.Now()
}

// Show 显示分类
func (c *Category) Show() {
	c.ShowStatus = 1
	c.UpdatedAt = time.Now()
}

// Hide 隐藏分类
func (c *Category) Hide() {
	c.ShowStatus = 0
	c.UpdatedAt = time.Now()
}

// ShowInNav 在导航栏显示
func (c *Category) ShowInNav() {
	c.NavStatus = 1
	c.UpdatedAt = time.Now()
}

// HideInNav 在导航栏隐藏
func (c *Category) HideInNav() {
	c.NavStatus = 0
	c.UpdatedAt = time.Now()
}

// IsRoot 是否为一级分类
func (c *Category) IsRoot() bool {
	return c.Level == 0
}

// IsLeaf 是否为叶子分类
func (c *Category) IsLeaf(hasChildren bool) bool {
	return !hasChildren
}
