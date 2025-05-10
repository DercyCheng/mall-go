package repository

import (
	"context"

	"mall-go/services/payment-service/domain/model"
)

// PaymentRepository defines the interface for payment persistence operations
type PaymentRepository interface {
	// Create creates a new payment record
	Create(ctx context.Context, payment *model.Payment) error
	
	// GetByID retrieves a payment by ID
	GetByID(ctx context.Context, id string) (*model.Payment, error)
	
	// GetByOrderID retrieves a payment by order ID
	GetByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
	
	// GetByTransactionID retrieves a payment by transaction ID
	GetByTransactionID(ctx context.Context, transactionID string) (*model.Payment, error)
	
	// Update updates an existing payment
	Update(ctx context.Context, payment *model.Payment) error
	
	// ListByUserID retrieves all payments for a user with pagination
	ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.Payment, int, error)
	
	// ListByStatus retrieves all payments with a specific status with pagination
	ListByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.Payment, int, error)
	
	// ListAll retrieves all payments with pagination
	ListAll(ctx context.Context, page, pageSize int) ([]*model.Payment, int, error)
}

// PaymentRefundRepository defines the interface for payment refund persistence operations
type PaymentRefundRepository interface {
	// Create creates a new payment refund record
	Create(ctx context.Context, refund *model.PaymentRefund) error
	
	// GetByID retrieves a refund by ID
	GetByID(ctx context.Context, id string) (*model.PaymentRefund, error)
	
	// GetByPaymentID retrieves refunds by payment ID
	GetByPaymentID(ctx context.Context, paymentID string) ([]*model.PaymentRefund, error)
	
	// GetByOrderID retrieves refunds by order ID
	GetByOrderID(ctx context.Context, orderID string) ([]*model.PaymentRefund, error)
	
	// Update updates an existing refund
	Update(ctx context.Context, refund *model.PaymentRefund) error
	
	// ListByUserID retrieves all refunds for a user with pagination
	ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*model.PaymentRefund, int, error)
	
	// ListByStatus retrieves all refunds with a specific status with pagination
	ListByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.PaymentRefund, int, error)
	
	// ListAll retrieves all refunds with pagination
	ListAll(ctx context.Context, page, pageSize int) ([]*model.PaymentRefund, int, error)
}

// PaymentCache defines the interface for payment caching operations
type PaymentCache interface {
	// SetPayment caches a payment
	SetPayment(ctx context.Context, payment *model.Payment, ttl int) error
	
	// GetPayment retrieves a cached payment by ID
	GetPayment(ctx context.Context, id string) (*model.Payment, error)
	
	// DeletePayment removes a cached payment
	DeletePayment(ctx context.Context, id string) error
	
	// SetPaymentByOrderID caches a payment by order ID
	SetPaymentByOrderID(ctx context.Context, payment *model.Payment, ttl int) error
	
	// GetPaymentByOrderID retrieves a cached payment by order ID
	GetPaymentByOrderID(ctx context.Context, orderID string) (*model.Payment, error)
	
	// DeletePaymentByOrderID removes a cached payment by order ID
	DeletePaymentByOrderID(ctx context.Context, orderID string) error
}
