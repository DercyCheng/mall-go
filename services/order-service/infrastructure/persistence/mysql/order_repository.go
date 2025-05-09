package mysql

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"mall-go/services/order-service/domain/model"
	"mall-go/services/order-service/domain/repository"
	"mall-go/services/order-service/infrastructure/util"
)

// OrderEntity 订单数据库实体
type OrderEntity struct {
	ID                    string    `gorm:"primaryKey;type:varchar(36)"`
	MemberID              string    `gorm:"column:member_id;type:varchar(36)"`
	OrderSn               string    `gorm:"column:order_sn;type:varchar(64);uniqueIndex"`
	MemberUsername        string    `gorm:"column:member_username;type:varchar(64)"`
	TotalAmount           float64   `gorm:"column:total_amount;type:decimal(10,2)"`
	PayAmount             float64   `gorm:"column:pay_amount;type:decimal(10,2)"`
	FreightAmount         float64   `gorm:"column:freight_amount;type:decimal(10,2)"`
	PromotionAmount       float64   `gorm:"column:promotion_amount;type:decimal(10,2)"`
	IntegrationAmount     float64   `gorm:"column:integration_amount;type:decimal(10,2)"`
	CouponAmount          float64   `gorm:"column:coupon_amount;type:decimal(10,2)"`
	DiscountAmount        float64   `gorm:"column:discount_amount;type:decimal(10,2)"`
	PayType               string    `gorm:"column:pay_type;type:varchar(20)"`
	SourceType            int       `gorm:"column:source_type;type:int"`
	Status                string    `gorm:"column:status;type:varchar(20);index"`
	OrderType             int       `gorm:"column:order_type;type:int"`
	DeliveryCompany       string    `gorm:"column:delivery_company;type:varchar(64)"`
	DeliverySn            string    `gorm:"column:delivery_sn;type:varchar(64)"`
	AutoConfirmDay        int       `gorm:"column:auto_confirm_day;type:int"`
	Integration           int       `gorm:"column:integration;type:int"`
	Growth                int       `gorm:"column:growth;type:int"`
	PromotionInfo         string    `gorm:"column:promotion_info;type:varchar(500)"`
	BillType              int       `gorm:"column:bill_type;type:int"`
	BillHeader            string    `gorm:"column:bill_header;type:varchar(200)"`
	BillContent           string    `gorm:"column:bill_content;type:varchar(200)"`
	BillReceiverPhone     string    `gorm:"column:bill_receiver_phone;type:varchar(32)"`
	BillReceiverEmail     string    `gorm:"column:bill_receiver_email;type:varchar(200)"`
	ReceiverName          string    `gorm:"column:receiver_name;type:varchar(100)"`
	ReceiverPhone         string    `gorm:"column:receiver_phone;type:varchar(32)"`
	ReceiverPostCode      string    `gorm:"column:receiver_post_code;type:varchar(32)"`
	ReceiverProvince      string    `gorm:"column:receiver_province;type:varchar(32)"`
	ReceiverCity          string    `gorm:"column:receiver_city;type:varchar(32)"`
	ReceiverDistrict      string    `gorm:"column:receiver_district;type:varchar(32)"`
	ReceiverDetailAddress string    `gorm:"column:receiver_detail_address;type:varchar(200)"`
	Note                  string    `gorm:"column:note;type:varchar(500)"`
	ConfirmStatus         int       `gorm:"column:confirm_status;type:int"`
	DeleteStatus          int       `gorm:"column:delete_status;type:int;default:0"` // 0->未删除，1->已删除
	UseIntegration        int       `gorm:"column:use_integration;type:int"`
	PaymentTime           time.Time `gorm:"column:payment_time;type:datetime"`
	DeliveryTime          time.Time `gorm:"column:delivery_time;type:datetime"`
	ReceiveTime           time.Time `gorm:"column:receive_time;type:datetime"`
	CommentTime           time.Time `gorm:"column:comment_time;type:datetime"`
	ModifyTime            time.Time `gorm:"column:modify_time;type:datetime"`
	CreatedAt             time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt             time.Time `gorm:"column:updated_at;type:datetime"`
}

// TableName 返回表名
func (OrderEntity) TableName() string {
	return "oms_order"
}

// OrderItemEntity 订单项数据库实体
type OrderItemEntity struct {
	ID                string    `gorm:"primaryKey;type:varchar(36)"`
	OrderId           string    `gorm:"column:order_id;type:varchar(36);index"`
	OrderSn           string    `gorm:"column:order_sn;type:varchar(64)"`
	ProductID         string    `gorm:"column:product_id;type:varchar(36)"`
	ProductSn         string    `gorm:"column:product_sn;type:varchar(64)"`
	ProductName       string    `gorm:"column:product_name;type:varchar(200)"`
	ProductPic        string    `gorm:"column:product_pic;type:varchar(500)"`
	ProductPrice      float64   `gorm:"column:product_price;type:decimal(10,2)"`
	ProductQuantity   int       `gorm:"column:product_quantity;type:int"`
	ProductAttr       string    `gorm:"column:product_attr;type:varchar(500)"`
	CouponAmount      float64   `gorm:"column:coupon_amount;type:decimal(10,2)"`
	PromotionAmount   float64   `gorm:"column:promotion_amount;type:decimal(10,2)"`
	RealAmount        float64   `gorm:"column:real_amount;type:decimal(10,2)"`
	GiftIntegration   int       `gorm:"column:gift_integration;type:int"`
	GiftGrowth        int       `gorm:"column:gift_growth;type:int"`
	ProductCategoryId string    `gorm:"column:product_category_id;type:varchar(36)"`
	CreatedAt         time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt         time.Time `gorm:"column:updated_at;type:datetime"`
}

// TableName 返回表名
func (OrderItemEntity) TableName() string {
	return "oms_order_item"
}

// OrderRepository 订单仓储MySQL实现
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// Save 保存订单
func (r *orderRepository) Save(ctx context.Context, order *model.Order) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 转换为数据库实体
		orderEntity := mapOrderToEntity(order)

		// 保存订单
		if err := tx.Create(orderEntity).Error; err != nil {
			return err
		}

		// 保存订单项
		for i := range order.OrderItems {
			// 生成订单项ID
			if order.OrderItems[i].ID == "" {
				order.OrderItems[i].ID = util.GenerateUUID()
			}

			// 创建订单项实体
			orderItemEntity := mapOrderItemToEntity(&order.OrderItems[i], order.ID, order.OrderSn)

			// 保存订单项
			if err := tx.Create(orderItemEntity).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FindByID 根据ID查找订单
func (r *orderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	var orderEntity OrderEntity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&orderEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// 查询订单项
	var orderItemEntities []OrderItemEntity
	err = r.db.WithContext(ctx).Where("order_id = ?", id).Find(&orderItemEntities).Error
	if err != nil {
		return nil, err
	}

	// 转换为领域模型
	return mapEntityToOrder(&orderEntity, orderItemEntities), nil
}

// FindByOrderSn 根据订单编号查找订单
func (r *orderRepository) FindByOrderSn(ctx context.Context, orderSn string) (*model.Order, error) {
	var orderEntity OrderEntity
	err := r.db.WithContext(ctx).Where("order_sn = ?", orderSn).First(&orderEntity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// 查询订单项
	var orderItemEntities []OrderItemEntity
	err = r.db.WithContext(ctx).Where("order_sn = ?", orderSn).Find(&orderItemEntities).Error
	if err != nil {
		return nil, err
	}

	// 转换为领域模型
	return mapEntityToOrder(&orderEntity, orderItemEntities), nil
}

// FindByMemberID 查询会员的订单列表
func (r *orderRepository) FindByMemberID(ctx context.Context, memberID string, page, size int) ([]*model.Order, int64, error) {
	var orderEntities []OrderEntity
	var total int64

	// 查询总数
	err := r.db.WithContext(ctx).Model(&OrderEntity{}).
		Where("member_id = ? AND delete_status = 0", memberID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = r.db.WithContext(ctx).
		Where("member_id = ? AND delete_status = 0", memberID).
		Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&orderEntities).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询订单项
	orders := make([]*model.Order, len(orderEntities))
	for i, orderEntity := range orderEntities {
		var orderItemEntities []OrderItemEntity
		err = r.db.WithContext(ctx).Where("order_id = ?", orderEntity.ID).Find(&orderItemEntities).Error
		if err != nil {
			return nil, 0, err
		}

		// 转换为领域模型
		orders[i] = mapEntityToOrder(&orderEntity, orderItemEntities)
	}

	return orders, total, nil
}

// Update 更新订单
func (r *orderRepository) Update(ctx context.Context, order *model.Order) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 转换为数据库实体
		orderEntity := mapOrderToEntity(order)

		// 更新订单
		if err := tx.Save(orderEntity).Error; err != nil {
			return err
		}

		// 更新订单项（先删除后插入）
		if err := tx.Where("order_id = ?", order.ID).Delete(&OrderItemEntity{}).Error; err != nil {
			return err
		}

		// 重新插入订单项
		for i := range order.OrderItems {
			// 生成订单项ID
			if order.OrderItems[i].ID == "" {
				order.OrderItems[i].ID = util.GenerateUUID()
			}

			// 创建订单项实体
			orderItemEntity := mapOrderItemToEntity(&order.OrderItems[i], order.ID, order.OrderSn)

			// 保存订单项
			if err := tx.Create(orderItemEntity).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete 删除订单(逻辑删除)
func (r *orderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&OrderEntity{}).
		Where("id = ?", id).
		Update("delete_status", 1).
		Error
}

// List 分页查询订单列表
func (r *orderRepository) List(ctx context.Context, query repository.OrderQuery) ([]*model.Order, int64, error) {
	var orderEntities []OrderEntity
	var total int64

	// 构建查询条件
	db := r.db.WithContext(ctx).Model(&OrderEntity{}).Where("delete_status = 0")

	if query.OrderSn != "" {
		db = db.Where("order_sn = ?", query.OrderSn)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.MemberUsername != "" {
		db = db.Where("member_username LIKE ?", "%"+query.MemberUsername+"%")
	}

	if query.ReceiverName != "" {
		db = db.Where("receiver_name LIKE ?", "%"+query.ReceiverName+"%")
	}

	if query.ReceiverPhone != "" {
		db = db.Where("receiver_phone LIKE ?", "%"+query.ReceiverPhone+"%")
	}

	if query.CreateTimeBegin != "" && query.CreateTimeEnd != "" {
		db = db.Where("created_at BETWEEN ? AND ?", query.CreateTimeBegin, query.CreateTimeEnd)
	} else if query.CreateTimeBegin != "" {
		db = db.Where("created_at >= ?", query.CreateTimeBegin)
	} else if query.CreateTimeEnd != "" {
		db = db.Where("created_at <= ?", query.CreateTimeEnd)
	}

	if query.SourceType != -1 {
		db = db.Where("source_type = ?", query.SourceType)
	}

	if query.OrderType != -1 {
		db = db.Where("order_type = ?", query.OrderType)
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = db.Order("created_at DESC").
		Offset((query.Page - 1) * query.Size).
		Limit(query.Size).
		Find(&orderEntities).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询订单项
	orders := make([]*model.Order, len(orderEntities))
	for i, orderEntity := range orderEntities {
		var orderItemEntities []OrderItemEntity
		err = r.db.WithContext(ctx).Where("order_id = ?", orderEntity.ID).Find(&orderItemEntities).Error
		if err != nil {
			return nil, 0, err
		}

		// 转换为领域模型
		orders[i] = mapEntityToOrder(&orderEntity, orderItemEntities)
	}

	return orders, total, nil
}

// UpdateStatus 更新订单状态
func (r *orderRepository) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&OrderEntity{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// UpdateNote 更新订单备注
func (r *orderRepository) UpdateNote(ctx context.Context, id string, note string) error {
	return r.db.WithContext(ctx).
		Model(&OrderEntity{}).
		Where("id = ?", id).
		Update("note", note).
		Error
}

// UpdateReceiverInfo 更新收货人信息
func (r *orderRepository) UpdateReceiverInfo(ctx context.Context, id string, receiverInfo map[string]string) error {
	// 构建更新字段
	updates := make(map[string]interface{})
	for key, value := range receiverInfo {
		if key == "receiverName" || key == "receiverPhone" || key == "receiverPostCode" ||
			key == "receiverProvince" || key == "receiverCity" || key == "receiverDistrict" ||
			key == "receiverDetailAddress" {
			updates[key] = value
		}
	}

	if len(updates) > 0 {
		return r.db.WithContext(ctx).
			Model(&OrderEntity{}).
			Where("id = ?", id).
			Updates(updates).
			Error
	}

	return nil
}

// mapOrderToEntity 将订单领域模型转换为数据库实体
func mapOrderToEntity(order *model.Order) *OrderEntity {
	if order == nil {
		return nil
	}

	// 处理零值时间
	paymentTime := order.PaymentTime
	deliveryTime := order.DeliveryTime
	receiveTime := order.ReceiveTime
	commentTime := order.CommentTime
	modifyTime := order.ModifyTime

	return &OrderEntity{
		ID:                    order.ID,
		MemberID:              order.MemberID,
		OrderSn:               order.OrderSn,
		MemberUsername:        order.MemberUsername,
		TotalAmount:           order.TotalAmount.Amount,
		PayAmount:             order.PayAmount.Amount,
		FreightAmount:         order.FreightAmount.Amount,
		PromotionAmount:       order.PromotionAmount.Amount,
		IntegrationAmount:     order.IntegrationAmount.Amount,
		CouponAmount:          order.CouponAmount.Amount,
		DiscountAmount:        order.DiscountAmount.Amount,
		PayType:               string(order.PayType),
		SourceType:            order.SourceType,
		Status:                string(order.Status),
		OrderType:             order.OrderType,
		DeliveryCompany:       order.DeliveryCompany,
		DeliverySn:            order.DeliverySn,
		AutoConfirmDay:        order.AutoConfirmDay,
		Integration:           order.Integration,
		Growth:                order.Growth,
		PromotionInfo:         order.PromotionInfo,
		BillType:              order.BillType,
		BillHeader:            order.BillHeader,
		BillContent:           order.BillContent,
		BillReceiverPhone:     order.BillReceiverPhone,
		BillReceiverEmail:     order.BillReceiverEmail,
		ReceiverName:          order.ReceiverName,
		ReceiverPhone:         order.ReceiverPhone,
		ReceiverPostCode:      order.ReceiverPostCode,
		ReceiverProvince:      order.ReceiverProvince,
		ReceiverCity:          order.ReceiverCity,
		ReceiverDistrict:      order.ReceiverDistrict,
		ReceiverDetailAddress: order.ReceiverDetailAddress,
		Note:                  order.Note,
		ConfirmStatus:         order.ConfirmStatus,
		DeleteStatus:          order.DeleteStatus,
		UseIntegration:        order.UseIntegration,
		PaymentTime:           paymentTime,
		DeliveryTime:          deliveryTime,
		ReceiveTime:           receiveTime,
		CommentTime:           commentTime,
		ModifyTime:            modifyTime,
		CreatedAt:             order.CreatedAt,
		UpdatedAt:             order.UpdatedAt,
	}
}

// mapOrderItemToEntity 将订单项领域模型转换为数据库实体
func mapOrderItemToEntity(orderItem *model.OrderItem, orderID string, orderSn string) *OrderItemEntity {
	if orderItem == nil {
		return nil
	}

	return &OrderItemEntity{
		ID:                orderItem.ID,
		OrderId:           orderID,
		OrderSn:           orderSn,
		ProductID:         orderItem.ProductID,
		ProductSn:         orderItem.ProductSn,
		ProductName:       orderItem.ProductName,
		ProductPic:        orderItem.ProductPic,
		ProductPrice:      orderItem.ProductPrice.Amount,
		ProductQuantity:   orderItem.ProductQuantity,
		ProductAttr:       orderItem.ProductAttr,
		CouponAmount:      orderItem.CouponAmount,
		PromotionAmount:   orderItem.PromotionAmount,
		RealAmount:        orderItem.RealAmount,
		GiftIntegration:   orderItem.GiftIntegration,
		GiftGrowth:        orderItem.GiftGrowth,
		ProductCategoryId: orderItem.ProductCategoryId,
		CreatedAt:         orderItem.CreatedAt,
		UpdatedAt:         time.Now(),
	}
}

// mapEntityToOrder 将数据库实体转换为订单领域模型
func mapEntityToOrder(orderEntity *OrderEntity, orderItemEntities []OrderItemEntity) *model.Order {
	if orderEntity == nil {
		return nil
	}

	// 转换订单项
	orderItems := make([]model.OrderItem, len(orderItemEntities))
	for i, entity := range orderItemEntities {
		orderItems[i] = model.OrderItem{
			ID:                entity.ID,
			ProductID:         entity.ProductID,
			ProductSn:         entity.ProductSn,
			ProductName:       entity.ProductName,
			ProductPic:        entity.ProductPic,
			ProductPrice:      model.Money{Amount: entity.ProductPrice, Currency: "CNY"},
			ProductQuantity:   entity.ProductQuantity,
			ProductAttr:       entity.ProductAttr,
			CouponAmount:      entity.CouponAmount,
			PromotionAmount:   entity.PromotionAmount,
			RealAmount:        entity.RealAmount,
			GiftIntegration:   entity.GiftIntegration,
			GiftGrowth:        entity.GiftGrowth,
			ProductCategoryId: entity.ProductCategoryId,
			CreatedAt:         entity.CreatedAt,
		}
	}

	return &model.Order{
		ID:                    orderEntity.ID,
		MemberID:              orderEntity.MemberID,
		OrderSn:               orderEntity.OrderSn,
		MemberUsername:        orderEntity.MemberUsername,
		TotalAmount:           model.Money{Amount: orderEntity.TotalAmount, Currency: "CNY"},
		PayAmount:             model.Money{Amount: orderEntity.PayAmount, Currency: "CNY"},
		FreightAmount:         model.Money{Amount: orderEntity.FreightAmount, Currency: "CNY"},
		PromotionAmount:       model.Money{Amount: orderEntity.PromotionAmount, Currency: "CNY"},
		IntegrationAmount:     model.Money{Amount: orderEntity.IntegrationAmount, Currency: "CNY"},
		CouponAmount:          model.Money{Amount: orderEntity.CouponAmount, Currency: "CNY"},
		DiscountAmount:        model.Money{Amount: orderEntity.DiscountAmount, Currency: "CNY"},
		PayType:               model.PaymentType(orderEntity.PayType),
		SourceType:            orderEntity.SourceType,
		Status:                model.OrderStatus(orderEntity.Status),
		OrderType:             orderEntity.OrderType,
		DeliveryCompany:       orderEntity.DeliveryCompany,
		DeliverySn:            orderEntity.DeliverySn,
		AutoConfirmDay:        orderEntity.AutoConfirmDay,
		Integration:           orderEntity.Integration,
		Growth:                orderEntity.Growth,
		PromotionInfo:         orderEntity.PromotionInfo,
		BillType:              orderEntity.BillType,
		BillHeader:            orderEntity.BillHeader,
		BillContent:           orderEntity.BillContent,
		BillReceiverPhone:     orderEntity.BillReceiverPhone,
		BillReceiverEmail:     orderEntity.BillReceiverEmail,
		ReceiverName:          orderEntity.ReceiverName,
		ReceiverPhone:         orderEntity.ReceiverPhone,
		ReceiverPostCode:      orderEntity.ReceiverPostCode,
		ReceiverProvince:      orderEntity.ReceiverProvince,
		ReceiverCity:          orderEntity.ReceiverCity,
		ReceiverDistrict:      orderEntity.ReceiverDistrict,
		ReceiverDetailAddress: orderEntity.ReceiverDetailAddress,
		Note:                  orderEntity.Note,
		ConfirmStatus:         orderEntity.ConfirmStatus,
		DeleteStatus:          orderEntity.DeleteStatus,
		UseIntegration:        orderEntity.UseIntegration,
		PaymentTime:           orderEntity.PaymentTime,
		DeliveryTime:          orderEntity.DeliveryTime,
		ReceiveTime:           orderEntity.ReceiveTime,
		CommentTime:           orderEntity.CommentTime,
		ModifyTime:            orderEntity.ModifyTime,
		OrderItems:            orderItems,
		CreatedAt:             orderEntity.CreatedAt,
		UpdatedAt:             orderEntity.UpdatedAt,
	}
}
