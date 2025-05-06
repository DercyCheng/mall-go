package model

import (
	"errors"
	"time"
)

// Promotion 促销实体
type Promotion struct {
	ID           string
	Type         PromotionType
	Name         string
	Description  string
	Discount     float64 // 折扣值或满减金额
	StartTime    time.Time
	EndTime      time.Time
	Status       int                    // 0->未开始；1->进行中；2->已结束
	Requirements map[string]interface{} // 促销要求参数
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewPromotion 创建促销的工厂方法
func NewPromotion(name string, promotionType PromotionType, discount float64, startTime, endTime time.Time) (*Promotion, error) {
	if name == "" {
		return nil, errors.New("promotion name cannot be empty")
	}

	if discount <= 0 {
		return nil, errors.New("discount must be greater than zero")
	}

	if startTime.After(endTime) {
		return nil, errors.New("start time must be before end time")
	}

	status := 0
	now := time.Now()
	if now.After(startTime) && now.Before(endTime) {
		status = 1
	} else if now.After(endTime) {
		status = 2
	}

	return &Promotion{
		Name:         name,
		Type:         promotionType,
		Discount:     discount,
		StartTime:    startTime,
		EndTime:      endTime,
		Status:       status,
		Requirements: make(map[string]interface{}),
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsActive 检查促销是否活动中
func (p *Promotion) IsActive() bool {
	now := time.Now()
	return now.After(p.StartTime) && now.Before(p.EndTime)
}

// Update 更新促销信息
func (p *Promotion) Update(name string, discount float64, startTime, endTime time.Time) error {
	if name != "" {
		p.Name = name
	}

	if discount > 0 {
		p.Discount = discount
	}

	if !startTime.IsZero() && !endTime.IsZero() {
		if startTime.After(endTime) {
			return errors.New("start time must be before end time")
		}
		p.StartTime = startTime
		p.EndTime = endTime
	}

	status := 0
	now := time.Now()
	if now.After(p.StartTime) && now.Before(p.EndTime) {
		status = 1
	} else if now.After(p.EndTime) {
		status = 2
	}

	p.Status = status
	p.UpdatedAt = now
	return nil
}

// AddRequirement 添加促销要求
func (p *Promotion) AddRequirement(key string, value interface{}) {
	p.Requirements[key] = value
	p.UpdatedAt = time.Now()
}

// RemoveRequirement 移除促销要求
func (p *Promotion) RemoveRequirement(key string) {
	delete(p.Requirements, key)
	p.UpdatedAt = time.Now()
}

// CalculateDiscountedPrice 计算折扣后价格
func (p *Promotion) CalculateDiscountedPrice(originalPrice float64) float64 {
	if !p.IsActive() {
		return originalPrice
	}

	switch p.Type {
	case PromotionTypePercentage:
		return originalPrice * (100 - p.Discount) / 100
	case PromotionTypeFixedAmount:
		if originalPrice > p.Requirements["minAmount"].(float64) {
			return originalPrice - p.Discount
		}
	}

	return originalPrice
}
