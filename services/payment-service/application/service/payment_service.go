package service

import (
	"context"

	"mall-go/services/payment-service/application/dto"
)

// PaymentService defines the interface for payment business operations
type PaymentService interface {
	// CreatePayment creates a new payment
	CreatePayment(ctx context.Context, request dto.CreatePaymentRequest) (*dto.PaymentDTO, error)
	
	// GetPaymentByID retrieves a payment by ID
	GetPaymentByID(ctx context.Context, id string) (*dto.PaymentDTO, error)
	
	// GetPaymentByOrderID retrieves a payment by order ID
	GetPaymentByOrderID(ctx context.Context, orderID string) (*dto.PaymentDTO, error)
	
	// GetPaymentByTransactionID retrieves a payment by transaction ID
	GetPaymentByTransactionID(ctx context.Context, transactionID string) (*dto.PaymentDTO, error)
	
	// ListPaymentsByUserID lists all payments for a user
	ListPaymentsByUserID(ctx context.Context, userID string, page, pageSize int) ([]*dto.PaymentDTO, int, error)
	
	// UpdatePaymentStatus updates a payment's status
	UpdatePaymentStatus(ctx context.Context, request dto.UpdatePaymentStatusRequest) (*dto.PaymentDTO, error)
	
	// ProcessPaymentCallback processes a payment callback from a payment provider
	ProcessPaymentCallback(ctx context.Context, request dto.PaymentCallbackRequest) (*dto.PaymentCallbackResponse, error)
	
	// QueryPaymentStatus queries a payment's status from the payment provider
	QueryPaymentStatus(ctx context.Context, request dto.QueryPaymentRequest) (*dto.PaymentDTO, error)
	
	// ListPaymentMethods lists all available payment methods
	ListPaymentMethods(ctx context.Context) ([]dto.PaymentMethodInfo, error)
}

// RefundService defines the interface for refund business operations
type RefundService interface {
	// CreateRefund creates a new refund request
	CreateRefund(ctx context.Context, request dto.CreateRefundRequest) (*dto.PaymentRefundDTO, error)
	
	// GetRefundByID retrieves a refund by ID
	GetRefundByID(ctx context.Context, id string) (*dto.PaymentRefundDTO, error)
	
	// GetRefundsByPaymentID retrieves all refunds for a payment
	GetRefundsByPaymentID(ctx context.Context, paymentID string) ([]*dto.PaymentRefundDTO, error)
	
	// GetRefundsByOrderID retrieves all refunds for an order
	GetRefundsByOrderID(ctx context.Context, orderID string) ([]*dto.PaymentRefundDTO, error)
	
	// ListRefundsByUserID lists all refunds for a user
	ListRefundsByUserID(ctx context.Context, userID string, page, pageSize int) ([]*dto.PaymentRefundDTO, int, error)
	
	// UpdateRefundStatus updates a refund's status
	UpdateRefundStatus(ctx context.Context, request dto.UpdateRefundStatusRequest) (*dto.PaymentRefundDTO, error)
	
	// ProcessRefundCallback processes a refund callback from a payment provider
	ProcessRefundCallback(ctx context.Context, request dto.RefundCallbackRequest) (*dto.PaymentCallbackResponse, error)
	
	// QueryRefundStatus queries a refund's status from the payment provider
	QueryRefundStatus(ctx context.Context, refundID string) (*dto.PaymentRefundDTO, error)
}
