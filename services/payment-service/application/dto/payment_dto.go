package dto

import (
	"time"
)

// CreatePaymentRequest represents a request to create a new payment
type CreatePaymentRequest struct {
	OrderID         string  `json:"order_id" validate:"required"`
	UserID          string  `json:"user_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Currency        string  `json:"currency" validate:"required,len=3"`
	PaymentMethod   string  `json:"payment_method" validate:"required"`
	PaymentProvider string  `json:"payment_provider" validate:"required"`
	ClientIP        string  `json:"client_ip"`
	ExpireMinutes   int     `json:"expire_minutes"`
}

// PaymentDTO represents a payment in the system
type PaymentDTO struct {
	ID              string    `json:"id"`
	OrderID         string    `json:"order_id"`
	TransactionID   string    `json:"transaction_id"`
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentMethod   string    `json:"payment_method"`
	PaymentProvider string    `json:"payment_provider"`
	Status          string    `json:"status"`
	PaymentTime     time.Time `json:"payment_time,omitempty"`
	ExpireTime      time.Time `json:"expire_time,omitempty"`
	CallbackTime    time.Time `json:"callback_time,omitempty"`
	SuccessTime     time.Time `json:"success_time,omitempty"`
	ErrorMsg        string    `json:"error_msg,omitempty"`
	ErrorCode       string    `json:"error_code,omitempty"`
	ClientIP        string    `json:"client_ip,omitempty"`
	PaymentURL      string    `json:"payment_url,omitempty"`  // URL for redirection to payment provider
	PaymentData     string    `json:"payment_data,omitempty"` // serialized payment data
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PaymentResponse represents a standard payment response
type PaymentResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *PaymentDTO `json:"data,omitempty"`
}

// PaymentsResponse represents a paginated list of payments response
type PaymentsResponse struct {
	Success  bool          `json:"success"`
	Message  string        `json:"message"`
	Data     []*PaymentDTO `json:"data"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// PaymentCallbackRequest represents a callback request from a payment provider
type PaymentCallbackRequest struct {
	PaymentProvider string                 `json:"payment_provider"`
	Parameters      map[string]interface{} `json:"parameters"`
}

// PaymentCallbackResponse represents a response to a payment callback
type PaymentCallbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdatePaymentStatusRequest represents a request to update a payment's status
type UpdatePaymentStatusRequest struct {
	ID            string `json:"id" validate:"required"`
	Status        string `json:"status" validate:"required"`
	TransactionID string `json:"transaction_id"`
	ErrorCode     string `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
}

// QueryPaymentRequest represents a request to query a payment's status
type QueryPaymentRequest struct {
	ID        string `json:"id,omitempty"`
	OrderID   string `json:"order_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	PaymentID string `json:"payment_id,omitempty"`
}

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	PaymentID      string  `json:"payment_id" validate:"required"`
	OrderID        string  `json:"order_id" validate:"required"`
	UserID         string  `json:"user_id" validate:"required"`
	RefundAmount   float64 `json:"refund_amount" validate:"required,gt=0"`
	RefundReason   string  `json:"refund_reason" validate:"required"`
	PaymentProvider string `json:"payment_provider" validate:"required"`
	OperatorID     string  `json:"operator_id"`
	OperatorName   string  `json:"operator_name"`
}

// PaymentRefundDTO represents a payment refund in the system
type PaymentRefundDTO struct {
	ID              string    `json:"id"`
	PaymentID       string    `json:"payment_id"`
	OrderID         string    `json:"order_id"`
	UserID          string    `json:"user_id"`
	RefundAmount    float64   `json:"refund_amount"`
	RefundReason    string    `json:"refund_reason"`
	Status          string    `json:"status"`
	TransactionID   string    `json:"transaction_id,omitempty"`
	OperatorID      string    `json:"operator_id,omitempty"`
	OperatorName    string    `json:"operator_name,omitempty"`
	RefundTime      time.Time `json:"refund_time"`
	SuccessTime     time.Time `json:"success_time,omitempty"`
	ErrorMsg        string    `json:"error_msg,omitempty"`
	ErrorCode       string    `json:"error_code,omitempty"`
	PaymentProvider string    `json:"payment_provider"`
	RefundData      string    `json:"refund_data,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RefundResponse represents a standard refund response
type RefundResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    *PaymentRefundDTO `json:"data,omitempty"`
}

// RefundsResponse represents a paginated list of refunds response
type RefundsResponse struct {
	Success  bool                `json:"success"`
	Message  string              `json:"message"`
	Data     []*PaymentRefundDTO `json:"data"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

// UpdateRefundStatusRequest represents a request to update a refund's status
type UpdateRefundStatusRequest struct {
	ID            string `json:"id" validate:"required"`
	Status        string `json:"status" validate:"required"`
	TransactionID string `json:"transaction_id"`
	ErrorCode     string `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
}

// RefundCallbackRequest represents a callback request from a payment provider for refund
type RefundCallbackRequest struct {
	PaymentProvider string                 `json:"payment_provider"`
	Parameters      map[string]interface{} `json:"parameters"`
}

// PaymentMethodInfo represents information about a payment method
type PaymentMethodInfo struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	IsEnabled   bool   `json:"is_enabled"`
}

// PaymentMethodsResponse represents a response with available payment methods
type PaymentMethodsResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    []PaymentMethodInfo `json:"data"`
}
