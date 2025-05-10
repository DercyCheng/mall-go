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

// PaymentRefundRepository implements repository.PaymentRefundRepository interface using MySQL
type PaymentRefundRepository struct {
	db *sqlx.DB
}

// NewPaymentRefundRepository creates a new PaymentRefundRepository instance
func NewPaymentRefundRepository(db *sqlx.DB) *PaymentRefundRepository {
	return &PaymentRefundRepository{
		db: db,
	}
}

// ensureTableExists creates the payment_refunds table if it doesn't exist
func (r *PaymentRefundRepository) ensureTableExists() error {
	schema := `
	CREATE TABLE IF NOT EXISTS payment_refunds (
		id VARCHAR(36) PRIMARY KEY,
		payment_id VARCHAR(36) NOT NULL,
		order_id VARCHAR(36) NOT NULL,
		user_id VARCHAR(36) NOT NULL,
		refund_amount DECIMAL(10,2) NOT NULL,
		refund_reason VARCHAR(255) NOT NULL,
		status VARCHAR(20) NOT NULL,
		transaction_id VARCHAR(128) DEFAULT NULL,
		operator_id VARCHAR(36) DEFAULT NULL,
		operator_name VARCHAR(50) DEFAULT NULL,
		refund_time TIMESTAMP NOT NULL,
		success_time TIMESTAMP NULL,
		error_msg VARCHAR(255) DEFAULT NULL,
		error_code VARCHAR(50) DEFAULT NULL,
		payment_provider VARCHAR(20) NOT NULL,
		refund_data TEXT DEFAULT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_payment_id (payment_id),
		INDEX idx_order_id (order_id),
		INDEX idx_user_id (user_id),
		INDEX idx_status (status)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := r.db.Exec(schema)
	return err
}

// Create creates a new payment refund record
func (r *PaymentRefundRepository) Create(ctx context.Context, refund *model.PaymentRefund) error {
	if err := r.ensureTableExists(); err != nil {
		return fmt.Errorf("ensure table exists error: %w", err)
	}
	
	// If no ID is provided, generate one
	if refund.ID == "" {
		refund.ID = uuid.New().String()
	}
	
	query := `
	INSERT INTO payment_refunds (
		id, payment_id, order_id, user_id, refund_amount, refund_reason, 
		status, transaction_id, operator_id, operator_name, refund_time, 
		success_time, error_msg, error_code, payment_provider, refund_data,
		created_at, updated_at
	) VALUES (
		:id, :payment_id, :order_id, :user_id, :refund_amount, :refund_reason, 
		:status, :transaction_id, :operator_id, :operator_name, :refund_time, 
		:success_time, :error_msg, :error_code, :payment_provider, :refund_data,
		:created_at, :updated_at
	)`
	
	params := map[string]interface{}{
		"id":               refund.ID,
		"payment_id":       refund.PaymentID,
		"order_id":         refund.OrderID,
		"user_id":          refund.UserID,
		"refund_amount":    refund.RefundAmount,
		"refund_reason":    refund.RefundReason,
		"status":           refund.Status,
		"transaction_id":   refund.TransactionID,
		"operator_id":      refund.OperatorID,
		"operator_name":    refund.OperatorName,
		"refund_time":      refund.RefundTime,
		"success_time":     sqlTimeOrNull(refund.SuccessTime),
		"error_msg":        refund.ErrorMsg,
		"error_code":       refund.ErrorCode,
		"payment_provider": refund.PaymentProvider,
		"refund_data":      refund.RefundData,
		"created_at":       refund.CreatedAt,
		"updated_at":       refund.UpdatedAt,
	}
	
	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("create payment refund error: %w", err)
	}
	
	return nil
}

// GetByID retrieves a refund by ID
func (r *PaymentRefundRepository) GetByID(ctx context.Context, id string) (*model.PaymentRefund, error) {
	var refund model.PaymentRefund
	
	query := `
	SELECT 
		id, payment_id, order_id, user_id, refund_amount, refund_reason, 
		status, transaction_id, operator_id, operator_name, refund_time, 
		success_time, error_msg, error_code, payment_provider, refund_data,
		created_at, updated_at
	FROM payment_refunds 
	WHERE id = ?
	`
	
	err := r.db.GetContext(ctx, &refund, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("payment refund not found: %w", repository.ErrNotFound)
		}
		return nil, fmt.Errorf("get payment refund error: %w", err)
	}
	
	return &refund, nil
}

// GetByPaymentID retrieves refunds by payment ID
func (r *PaymentRefundRepository) GetByPaymentID(ctx context.Context, paymentID string) ([]*model.PaymentRefund, error) {
	var refunds []*model.PaymentRefund
	
	query := `
	SELECT 
		id, payment_id, order_id, user_id, refund_amount, refund_reason, 
		status, transaction_id, operator_id, operator_name, refund_time, 
		success_time, error_msg, error_code, payment_provider, refund_data,
		created_at, updated_at
	FROM payment_refunds 
	WHERE payment_id = ?
	ORDER BY created_at DESC
	`
	
	err := r.db.SelectContext(ctx, &refunds, query, paymentID)
	if err != nil {
		return nil, fmt.Errorf("get payment refunds error: %w", err)
	}
	
	return refunds, nil
}

// GetByOrderID retrieves refunds by order ID
func (r *PaymentRefundRepository) GetByOrderID(ctx context.Context, orderID string) ([]*model.PaymentRefund, error) {
	var refunds []*model.PaymentRefund
	
	query := `
	SELECT 
		id, payment_id, order_id, user_id, refund_amount, refund_reason, 
		status, transaction_id, operator_id, operator_name, refund_time, 
		success_time, error_msg, error_code, payment_provider, refund_data,
		created_at, updated_at
	FROM payment_refunds 
	WHERE order_id = ?
	ORDER BY created_at DESC
	`
	
	err := r.db.SelectContext(ctx, &refunds, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get payment refunds error: %w", err)
	}
	
	return refunds, nil
}

// Update updates an existing refund
func (r *PaymentRefundRepository) Update(ctx context.Context, refund *model.PaymentRefund) error {
	query := `
	UPDATE payment_refunds 
	SET 
		status = :status,
		transaction_id = :transaction_id,
		success_time = :success_time,
		error_msg = :error_msg,
		error_code = :error_code,
		refund_data = :refund_data,
		updated_at = :updated_at
	WHERE id = :id
	`
	
	params := map[string]interface{}{
		"id":             refund.ID,
		"status":         refund.Status,
		"transaction_id": refund.TransactionID,
		"success_time":   sqlTimeOrNull(refund.SuccessTime),
		"error_msg":      refund.ErrorMsg,
		"error_code":     refund.ErrorCode,
		"refund_data":    refund.RefundData,
		"updated_at":     time.Now(),
	}
	
	result, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("update payment refund error: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected error: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("payment refund not found: %w", repository.ErrNotFound)
	}
	
	return nil
}

// ListByUserID retrieves all refunds for a user with pagination
func (r *PaymentRefundRepository) ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.PaymentRefund, int, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `SELECT COUNT(*) FROM payment_refunds WHERE user_id = ?`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, userID); err != nil {
		return nil, 0, fmt.Errorf("count payment refunds error: %w", err)
	}
	
	query := `
	SELECT 
		id, payment_id, order_id, user_id, refund_amount, refund_reason, 
		status, transaction_id, operator_id, operator_name, refund_time, 
		success_time, error_msg, error_code, payment_provider, refund_data,
		created_at, updated_at
	FROM payment_refunds 
	WHERE user_id = ?
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`
	
	var refunds []*model.PaymentRefund
	err := r.db.SelectContext(ctx, &refunds, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payment refunds error: %w", err)
	}
	
	return refunds, total, nil
}

// ListByStatus retrieves all refunds with a specific status with pagination
func (r *PaymentRefundRepository) ListByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.PaymentRefund, int, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `SELECT COUNT(*) FROM payment_refunds WHERE status = ?`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, status); err != nil {
		return nil, 0, fmt.Errorf("count payment refunds error: %w", err)
	}
	
	query := `
	SELECT 
		id, payment_id, order_id, user_id, refund_amount, refund_reason, 
		status, transaction_id, operator_id, operator_name, refund_time, 
		success_time, error_msg, error_code, payment_provider, refund_data,
		created_at, updated_at
	FROM payment_refunds 
	WHERE status = ?
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`
	
	var refunds []*model.PaymentRefund
	err := r.db.SelectContext(ctx, &refunds, query, status, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payment refunds error: %w", err)
	}
	
	return refunds, total, nil
}

// ListAll retrieves all refunds with pagination
func (r *PaymentRefundRepository) ListAll(ctx context.Context, page, pageSize int) ([]*model.PaymentRefund, int, error) {
	offset := (page - 1) * pageSize
	
	countQuery := `SELECT COUNT(*) FROM payment_refunds`
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("count payment refunds error: %w", err)
	}
	
	query := `
	SELECT 
		id, payment_id, order_id, user_id, refund_amount, refund_reason, 
		status, transaction_id, operator_id, operator_name, refund_time, 
		success_time, error_msg, error_code, payment_provider, refund_data,
		created_at, updated_at
	FROM payment_refunds 
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`
	
	var refunds []*model.PaymentRefund
	err := r.db.SelectContext(ctx, &refunds, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payment refunds error: %w", err)
	}
	
	return refunds, total, nil
}
