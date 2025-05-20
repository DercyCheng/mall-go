package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"mall-go/services/order-service/domain/model"
	"mall-go/services/order-service/domain/repository"
)

// OrderRepository 订单仓储实现
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Save 保存订单
func (r *OrderRepository) Save(ctx context.Context, order *model.Order) error {
	// 开始事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 序列化订单项
	orderItemsJSON, err := json.Marshal(order.OrderItems)
	if err != nil {
		return err
	}

	// 序列化地址
	shippingAddrJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return err
	}
	billingAddrJSON, err := json.Marshal(order.BillingAddress)
	if err != nil {
		return err
	}

	// 插入订单
	query := `
		INSERT INTO orders (
			id, user_id, order_sn, status, 
			total_amount, total_amount_currency, 
			pay_amount, pay_amount_currency,
			freight_amount, freight_amount_currency,
			discount_amount, discount_amount_currency,
			coupon_amount, coupon_amount_currency,
			point_amount, point_amount_currency,
			payment_method, payment_time,
			delivery_method, delivery_company, delivery_tracking_no, delivery_time,
			receive_time, comment_time,
			shipping_address, billing_address,
			note, order_items,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		order.ID, order.UserID, order.OrderSN, order.Status,
		order.TotalAmount.Amount, order.TotalAmount.Currency,
		order.PayAmount.Amount, order.PayAmount.Currency,
		order.FreightAmount.Amount, order.FreightAmount.Currency,
		order.DiscountAmount.Amount, order.DiscountAmount.Currency,
		order.CouponAmount.Amount, order.CouponAmount.Currency,
		order.PointAmount.Amount, order.PointAmount.Currency,
		order.PaymentMethod, order.PaymentTime,
		order.DeliveryMethod, order.DeliveryCompany, order.DeliveryTrackingNo, order.DeliveryTime,
		order.ReceiveTime, order.CommentTime,
		shippingAddrJSON, billingAddrJSON,
		order.Note, orderItemsJSON,
		order.CreatedAt, order.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// FindByID 根据ID查询订单
func (r *OrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	query := `
		SELECT 
			id, user_id, order_sn, status, 
			total_amount, total_amount_currency, 
			pay_amount, pay_amount_currency,
			freight_amount, freight_amount_currency,
			discount_amount, discount_amount_currency,
			coupon_amount, coupon_amount_currency,
			point_amount, point_amount_currency,
			payment_method, payment_time,
			delivery_method, delivery_company, delivery_tracking_no, delivery_time,
			receive_time, comment_time,
			shipping_address, billing_address,
			note, order_items,
			created_at, updated_at
		FROM orders
		WHERE id = ? AND deleted_at IS NULL
	`

	var order model.Order
	var totalAmountCurrency, payAmountCurrency, freightAmountCurrency, discountAmountCurrency, couponAmountCurrency, pointAmountCurrency string
	var shippingAddrJSON, billingAddrJSON, orderItemsJSON []byte
	var paymentMethod, deliveryMethod, status sql.NullString
	var paymentTime, deliveryTime, receiveTime, commentTime sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&order.ID, &order.UserID, &order.OrderSN, &status,
		&order.TotalAmount.Amount, &totalAmountCurrency,
		&order.PayAmount.Amount, &payAmountCurrency,
		&order.FreightAmount.Amount, &freightAmountCurrency,
		&order.DiscountAmount.Amount, &discountAmountCurrency,
		&order.CouponAmount.Amount, &couponAmountCurrency,
		&order.PointAmount.Amount, &pointAmountCurrency,
		&paymentMethod, &paymentTime,
		&deliveryMethod, &order.DeliveryCompany, &order.DeliveryTrackingNo, &deliveryTime,
		&receiveTime, &commentTime,
		&shippingAddrJSON, &billingAddrJSON,
		&order.Note, &orderItemsJSON,
		&order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, err
	}

	// 设置可空字段
	if status.Valid {
		order.Status = model.OrderStatus(status.String)
	}
	if paymentMethod.Valid {
		order.PaymentMethod = model.PaymentMethod(paymentMethod.String)
	}
	if deliveryMethod.Valid {
		order.DeliveryMethod = model.DeliveryMethod(deliveryMethod.String)
	}
	if paymentTime.Valid {
		order.PaymentTime = paymentTime.Time
	}
	if deliveryTime.Valid {
		order.DeliveryTime = deliveryTime.Time
	}
	if receiveTime.Valid {
		order.ReceiveTime = receiveTime.Time
	}
	if commentTime.Valid {
		order.CommentTime = commentTime.Time
	}

	// 设置货币类型
	order.TotalAmount.Currency = totalAmountCurrency
	order.PayAmount.Currency = payAmountCurrency
	order.FreightAmount.Currency = freightAmountCurrency
	order.DiscountAmount.Currency = discountAmountCurrency
	order.CouponAmount.Currency = couponAmountCurrency
	order.PointAmount.Currency = pointAmountCurrency

	// 反序列化地址
	if err = json.Unmarshal(shippingAddrJSON, &order.ShippingAddress); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(billingAddrJSON, &order.BillingAddress); err != nil {
		return nil, err
	}

	// 反序列化订单项
	if err = json.Unmarshal(orderItemsJSON, &order.OrderItems); err != nil {
		return nil, err
	}

	return &order, nil
}

// FindByOrderSN 根据订单编号查询订单
func (r *OrderRepository) FindByOrderSN(ctx context.Context, orderSN string) (*model.Order, error) {
	query := `
		SELECT id FROM orders
		WHERE order_sn = ? AND deleted_at IS NULL
	`

	var id string
	err := r.db.QueryRowContext(ctx, query, orderSN).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found with SN %s: %w", orderSN, err)
		}
		return nil, err
	}

	return r.FindByID(ctx, id)
}

// FindByUserID 根据用户ID查询订单列表
func (r *OrderRepository) FindByUserID(ctx context.Context, userID string, page, size int) ([]*model.Order, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM orders
		WHERE user_id = ? AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Order{}, 0, nil
	}

	// 分页查询订单ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM orders
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集订单ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询订单详情
	orders := make([]*model.Order, 0, len(ids))
	for _, id := range ids {
		order, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// FindByStatus 根据订单状态查询订单列表
func (r *OrderRepository) FindByStatus(ctx context.Context, status model.OrderStatus, page, size int) ([]*model.Order, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM orders
		WHERE status = ? AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery, status).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Order{}, 0, nil
	}

	// 分页查询订单ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM orders
		WHERE status = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, status, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集订单ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询订单详情
	orders := make([]*model.Order, 0, len(ids))
	for _, id := range ids {
		order, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// FindAll 分页查询所有订单
func (r *OrderRepository) FindAll(ctx context.Context, page, size int) ([]*model.Order, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM orders
		WHERE deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Order{}, 0, nil
	}

	// 分页查询订单ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM orders
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集订单ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询订单详情
	orders := make([]*model.Order, 0, len(ids))
	for _, id := range ids {
		order, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// Search 搜索订单
func (r *OrderRepository) Search(ctx context.Context, keyword string, page, size int) ([]*model.Order, int64, error) {
	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM orders
		WHERE (order_sn LIKE ? OR user_id LIKE ? OR note LIKE ?)
		AND deleted_at IS NULL
	`
	likeKeyword := "%" + keyword + "%"
	err := r.db.QueryRowContext(ctx, countQuery, likeKeyword, likeKeyword, likeKeyword).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Order{}, 0, nil
	}

	// 分页查询订单ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM orders
		WHERE (order_sn LIKE ? OR user_id LIKE ? OR note LIKE ?)
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, likeKeyword, likeKeyword, likeKeyword, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集订单ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询订单详情
	orders := make([]*model.Order, 0, len(ids))
	for _, id := range ids {
		order, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// Update 更新订单
func (r *OrderRepository) Update(ctx context.Context, order *model.Order) error {
	// 开始事务
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 序列化订单项
	orderItemsJSON, err := json.Marshal(order.OrderItems)
	if err != nil {
		return err
	}

	// 序列化地址
	shippingAddrJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return err
	}
	billingAddrJSON, err := json.Marshal(order.BillingAddress)
	if err != nil {
		return err
	}

	// 更新订单
	query := `
		UPDATE orders SET
			user_id = ?, order_sn = ?, status = ?, 
			total_amount = ?, total_amount_currency = ?, 
			pay_amount = ?, pay_amount_currency = ?,
			freight_amount = ?, freight_amount_currency = ?,
			discount_amount = ?, discount_amount_currency = ?,
			coupon_amount = ?, coupon_amount_currency = ?,
			point_amount = ?, point_amount_currency = ?,
			payment_method = ?, payment_time = ?,
			delivery_method = ?, delivery_company = ?, delivery_tracking_no = ?, delivery_time = ?,
			receive_time = ?, comment_time = ?,
			shipping_address = ?, billing_address = ?,
			note = ?, order_items = ?,
			updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		order.UserID, order.OrderSN, order.Status,
		order.TotalAmount.Amount, order.TotalAmount.Currency,
		order.PayAmount.Amount, order.PayAmount.Currency,
		order.FreightAmount.Amount, order.FreightAmount.Currency,
		order.DiscountAmount.Amount, order.DiscountAmount.Currency,
		order.CouponAmount.Amount, order.CouponAmount.Currency,
		order.PointAmount.Amount, order.PointAmount.Currency,
		order.PaymentMethod, order.PaymentTime,
		order.DeliveryMethod, order.DeliveryCompany, order.DeliveryTrackingNo, order.DeliveryTime,
		order.ReceiveTime, order.CommentTime,
		shippingAddrJSON, billingAddrJSON,
		order.Note, orderItemsJSON,
		time.Now(),
		order.ID,
	)

	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// Delete 删除订单（软删除）
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE orders SET
			deleted_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

// CountByStatus 统计各状态订单数量
func (r *OrderRepository) CountByStatus(ctx context.Context) (map[model.OrderStatus]int64, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM orders
		WHERE deleted_at IS NULL
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[model.OrderStatus]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[model.OrderStatus(status)] = count
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// FindByDateRange 根据日期范围查询订单
func (r *OrderRepository) FindByDateRange(ctx context.Context, startDate, endDate string, page, size int) ([]*model.Order, int64, error) {
	// 解析日期
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid start date format: %w", err)
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid end date format: %w", err)
	}
	// 设置为结束日期的最后一刻
	endTime = endTime.Add(24*time.Hour - time.Second)

	// 计算总数
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM orders
		WHERE created_at BETWEEN ? AND ?
		AND deleted_at IS NULL
	`
	err = r.db.QueryRowContext(ctx, countQuery, startTime, endTime).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Order{}, 0, nil
	}

	// 分页查询订单ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM orders
		WHERE created_at BETWEEN ? AND ?
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, startTime, endTime, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集订单ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询订单详情
	orders := make([]*model.Order, 0, len(ids))
	for _, id := range ids {
		order, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// FindByProductID 根据商品ID查询包含该商品的订单
func (r *OrderRepository) FindByProductID(ctx context.Context, productID string, page, size int) ([]*model.Order, int64, error) {
	// 计算总数 - 这里需要搜索JSON字段，根据具体数据库可能有不同实现
	var total int64
	countQuery := `
		SELECT COUNT(*) FROM orders
		WHERE JSON_CONTAINS(order_items, JSON_OBJECT('productID', ?))
		AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, countQuery, productID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回空列表
	if total == 0 {
		return []*model.Order{}, 0, nil
	}

	// 分页查询订单ID
	offset := (page - 1) * size
	query := `
		SELECT id FROM orders
		WHERE JSON_CONTAINS(order_items, JSON_OBJECT('productID', ?))
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, productID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 收集订单ID
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// 批量查询订单详情
	orders := make([]*model.Order, 0, len(ids))
	for _, id := range ids {
		order, err := r.FindByID(ctx, id)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	return orders, total, nil
}

// Ensure OrderRepository implements repository.OrderRepository
var _ repository.OrderRepository = (*OrderRepository)(nil)
