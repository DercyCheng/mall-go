package service

import (
	"context"
	
	"mall-go/services/payment-service/domain/model"
)

// PaymentProvider defines the interface for external payment providers
type PaymentProvider interface {
	// GeneratePaymentURL generates a payment URL for the given payment
	GeneratePaymentURL(ctx context.Context, payment *model.Payment) (string, error)
	
	// VerifyCallback verifies a payment callback from the provider
	// Returns payment ID, order ID, status (success or failed), transaction ID, and error
	VerifyCallback(ctx context.Context, parameters map[string]interface{}) (string, string, string, string, error)
	
	// QueryPaymentStatus queries the payment status from the provider
	// Returns status, transaction ID, and error
	QueryPaymentStatus(ctx context.Context, payment *model.Payment) (string, string, error)
	
	// Refund initiates a refund for a payment
	Refund(ctx context.Context, payment *model.Payment, refund *model.PaymentRefund) (string, error)
	
	// VerifyRefundCallback verifies a refund callback from the provider
	// Returns refund ID, status (success or failed), transaction ID, and error
	VerifyRefundCallback(ctx context.Context, parameters map[string]interface{}) (string, string, string, error)
	
	// QueryRefundStatus queries the refund status from the provider
	// Returns status, transaction ID, and error
	QueryRefundStatus(ctx context.Context, refund *model.PaymentRefund) (string, string, error)
}
