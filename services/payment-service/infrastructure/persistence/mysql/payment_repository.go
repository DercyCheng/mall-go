package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"

	"mall-go/services/payment-service/domain/model"
	"mall-go/services/payment-service/domain/repository"
)

// PaymentRepository implements repository.PaymentRepository interface using MySQL
type PaymentRepository struct {
	db *sqlx.DB
}

// NewPaymentRepository creates a new PaymentRepository instance
func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

// ensureTableExists creates the payments table if it doesn't exist
func (r *PaymentRepository) ensureTableExists() error {
	schema := `
	CREATE TABLE IF NOT EXISTS payments (
		id VARCHAR(36) PRIMARY KEY,
		order_id VARCHAR(36) NOT NULL,
		transaction_id VARCHAR(128) DEFAULT NULL,
		user_id VARCHAR(36) NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		currency VARCHAR(3) NOT NULL DEFAULT 'CNY',
		payment_method VARCHAR(20) NOT NULL,
		payment_provider VARCHAR(20) NOT NULL,
		status VARCHAR(20) NOT NULL,
		payment_time TIMESTAMP NULL,
		expire_time TIMESTAMP NULL,
		callback_time TIMESTAMP NULL,
		success_time TIMESTAMP NULL,
		error_msg VARCHAR(255) DEFAULT NULL,
		error_code VARCHAR(50) DEFAULT NULL,
		client_ip VARCHAR(50) DEFAULT NULL,
		payment_data TEXT DEFAULT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_order_id (order_id),
		INDEX idx_user_id (user_id),
		INDEX idx_status (status),
		INDEX idx_transaction_id (transaction_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := r.db.Exec(schema)
	return err
}

// Create creates a new payment record
func (r *PaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	if err := r.ensureTableExists(); err != nil {
		return fmt.Errorf("ensure table exists error: %w", err)
	}
	
	// If no ID is provided, generate one
	if payment.ID == "" {
		payment.ID = uuid.New().String()
	}
	
	query := `
	INSERT INTO payments (
		id, order_id, transaction_id, user_id, amount, currency, payment_method, 
		payment_provider, status, payment_time, expire_time, callback_time, 
		success_time, error_msg, error_code, client_ip, payment_data, 
		created_at, updated_at
	) VALUES (
		:id, :order_id, :transaction_id, :user_id, :amount, :currency, :payment_method, 
		:payment_provider, :status, :payment_time, :expire_time, :callback_time, 
		:success_time, :error_msg, :error_code, :client_ip, :payment_data, 
		:created_at, :updated_at
	)`
	
	params := map[string]interface{}{
		"id":               payment.ID,
		"order_id":         payment.OrderID,
		"transaction_id":   payment.TransactionID,
		"user_id":          payment.UserID,
		"amount":           payment.Amount,
		"currency":         payment.Currency,
		"payment_method":   payment.PaymentMethod,
		"payment_provider": payment.PaymentProvider,
		"status":           payment.Status,
		"payment_time":     sqlTimeOrNull(payment.PaymentTime),
		"expire_time":      sqlTimeOrNull(payment.ExpireTime),
		"callback_time":    sqlTimeOrNull(payment.CallbackTime),
		"success_time":     sqlTimeOrNull(payment.SuccessTime),
		"error_msg":        payment.ErrorMsg,
		"error_code":       payment.ErrorCode,
		"client_ip":        payment.ClientIP,
		"payment_data":     payment.PaymentData,
		"created_at":       payment.CreatedAt,
		"updated_at":       payment.UpdatedAt,
	}
	
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("create payment error: %w", err)
	}
	
	return nil
}

// GetByID retrieves a payment by ID
func (r *PaymentRepository) GetByID(ctx context.Context, id string) (*model.Payment, error) {
	var payment model.Payment
	
	query := `
	SELECT 
		id, order_id, transaction_id, user_id, amount, currency, payment_method, 
		payment_provider, status, payment_time, expire_time, callback_time, 
		success_time, error_msg, error_code, client_ip, payment_data, 
		created_at, updated_at
	FROM payments 
	WHERE id = ?
	`
	
	err := r.db.GetContext(ctx, &payment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found: %w", repository.ErrNotFound)
		}
		return nil, fmt.Errorf("get payment error: %w", err)
	}
	
	return &payment, nil
}

// GetByOrderID retrieves a payment by order ID
func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error) {
	var payment model.Payment
	
	query := `
	SELECT 
		id, order_id, transaction_id, user_id, amount, currency, payment_method, 
		payment_provider, status, payment_time, expire_time, callback_time, 
		success_time, error_msg, error_code, client_ip, payment_data, 
		created_at, updated_at
	FROM payments 
	WHERE order_id = ?
	ORDER BY created_at DESC
	LIMIT 1
	`
	
	err := r.db.GetContext(ctx, &payment, query, orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found: %w", repository.ErrNotFound)
		}
		return nil, fmt.Errorf("get payment error: %w", err)
	}
	
	return &payment, nil
}

// GetByTransactionID retrieves a payment by transaction ID
func (r *PaymentRepository) GetByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error) {
	var payment model.Payment
	
	query := `
	SELECT 
		id, order_id, transaction_id, user_id, amount, currency, payment_method, 
		payment_provider, status, payment_time, expire_time, callback_time, 
		success_time, error_msg, error_code, client_ip, payment_data, 
		created_at, updated_at
	FROM payments 
	WHERE transaction_id = ?
	`
	
	err := r.db.GetContext(ctx, &payment, query, transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment not found: %w", repository.ErrNotFound)
		}
		return nil, fmt.Errorf("get payment error: %w", err)
	}
	
	return &payment, nil
}

// Update updates an existing payment
func (r *PaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	query := `
	UPDATE payments 
	SET 
		transaction_id = :transaction_id,
		status = :status,
		payment_time = :payment_time,
		callback_time = :callback_time,
		success_time = :success_time,
		error_msg = :error_msg,
		error_code = :error_code,
		payment_data = :payment_data,
		updated_at = :updated_at
	WHERE id = :id
	`
	
	params := map[string]interface{}{
		"id":             payment.ID,
		"transaction_id": payment.TransactionID,
		"status":         payment.Status,
		"payment_time":   sqlTimeOrNull(payment.PaymentTime),
		"callback_time":  sqlTimeOrNull(payment.CallbackTime),
		"success_time":   sqlTimeOrNull(payment.SuccessTime),
		"error_msg":      payment.ErrorMsg,
		"error_code":     payment.ErrorCode,
		"payment_data":   payment.PaymentData,
		"updated_at":     time.Now(),
	}
	
	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("update payment error: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected error: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("payment not found: %w", repository.ErrNotFound)
	}
	
	return nil
}

// ListByUserID retrieves all payments for a user with pagination
func (r *PaymentRepository) ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Payment, int, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `SELECT COUNT(*) FROM payments WHERE user_id = ?`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, userID); err != nil {
		return nil, 0, fmt.Errorf("count payments error: %w", err)
	}
	
	query := `
	SELECT 
		id, order_id, transaction_id, user_id, amount, currency, payment_method, 
		payment_provider, status, payment_time, expire_time, callback_time, 
		success_time, error_msg, error_code, client_ip, payment_data, 
		created_at, updated_at
	FROM payments 
	WHERE user_id = ?
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`
	
	var payments []*model.Payment
	err := r.db.SelectContext(ctx, &payments, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments error: %w", err)
	}
	
	return payments, total, nil
}

// ListByStatus retrieves all payments with a specific status with pagination
func (r *PaymentRepository) ListByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.Payment, int, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `SELECT COUNT(*) FROM payments WHERE status = ?`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, status); err != nil {
		return nil, 0, fmt.Errorf("count payments error: %w", err)
	}
	
	query := `
	SELECT 
		id, order_id, transaction_id, user_id, amount, currency, payment_method, 
		payment_provider, status, payment_time, expire_time, callback_time, 
		success_time, error_msg, error_code, client_ip, payment_data, 
		created_at, updated_at
	FROM payments 
	WHERE status = ?
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`
	
	var payments []*model.Payment
	err := r.db.SelectContext(ctx, &payments, query, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments error: %w", err)
	}
	
	return payments, total, nil
}

// ListAll retrieves all payments with pagination
func (r *PaymentRepository) ListAll(ctx context.Context, page, pageSize int) ([]*model.Payment, int, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `SELECT COUNT(*) FROM payments`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("count payments error: %w", err)
	}
	
	query := `
	SELECT 
		id, order_id, transaction_id, user_id, amount, currency, payment_method, 
		payment_provider, status, payment_time, expire_time, callback_time, 
		success_time, error_msg, error_code, client_ip, payment_data, 
		created_at, updated_at
	FROM payments 
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`
	
	var payments []*model.Payment
	err := r.db.SelectContext(ctx, &payments, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments error: %w", err)
	}
	
	return payments, total, nil
}

// Helper function to handle null time values
func sqlTimeOrNull(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}
