package model

import (
	"errors"
	"time"
)

// Payment represents a payment transaction in the system
type Payment struct {
	ID              string    `json:"id"`
	OrderID         string    `json:"order_id"`
	TransactionID   string    `json:"transaction_id"`  // ID assigned by payment provider
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentMethod   string    `json:"payment_method"`  // alipay, wxpay, creditcard, etc.
	PaymentProvider string    `json:"payment_provider"` // provider name (alipay, wxpay, stripe, etc.)
	Status          string    `json:"status"`           // pending, processing, completed, failed, refunded
	PaymentTime     time.Time `json:"payment_time"`     // when payment was made
	ExpireTime      time.Time `json:"expire_time"`      // when payment will expire
	CallbackTime    time.Time `json:"callback_time"`    // when provider sent callback
	SuccessTime     time.Time `json:"success_time"`     // when payment was confirmed
	ErrorMsg        string    `json:"error_msg"`        // error message if failed
	ErrorCode       string    `json:"error_code"`       // error code if failed
	ClientIP        string    `json:"client_ip"`        // client IP address
	PaymentData     string    `json:"payment_data"`     // serialized payment data (JSON)
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PaymentStatusPending represents a pending payment status
const PaymentStatusPending = "pending"

// PaymentStatusProcessing represents a processing payment status
const PaymentStatusProcessing = "processing"

// PaymentStatusCompleted represents a completed payment status
const PaymentStatusCompleted = "completed"

// PaymentStatusFailed represents a failed payment status
const PaymentStatusFailed = "failed"

// PaymentStatusRefunded represents a refunded payment status
const PaymentStatusRefunded = "refunded"

// PaymentStatusCanceled represents a canceled payment status
const PaymentStatusCanceled = "canceled"

// NewPayment creates a new payment
func NewPayment(
	id, orderID, userID string,
	amount float64,
	currency, paymentMethod, paymentProvider string,
	expireTime time.Time,
) *Payment {
	now := time.Now()
	return &Payment{
		ID:              id,
		OrderID:         orderID,
		UserID:          userID,
		Amount:          amount,
		Currency:        currency,
		PaymentMethod:   paymentMethod,
		PaymentProvider: paymentProvider,
		Status:          PaymentStatusPending,
		ExpireTime:      expireTime,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// IsExpired checks if the payment is expired
func (p *Payment) IsExpired() bool {
	return !p.ExpireTime.IsZero() && time.Now().After(p.ExpireTime)
}

// SetCompleted marks the payment as completed
func (p *Payment) SetCompleted(transactionID string) error {
	if p.Status == PaymentStatusCompleted {
		return errors.New("payment already completed")
	}
	
	p.Status = PaymentStatusCompleted
	p.TransactionID = transactionID
	p.SuccessTime = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// SetFailed marks the payment as failed
func (p *Payment) SetFailed(errorCode, errorMsg string) error {
	if p.Status == PaymentStatusCompleted {
		return errors.New("cannot fail a completed payment")
	}
	
	p.Status = PaymentStatusFailed
	p.ErrorCode = errorCode
	p.ErrorMsg = errorMsg
	p.UpdatedAt = time.Now()
	return nil
}

// SetProcessing marks the payment as processing
func (p *Payment) SetProcessing() error {
	if p.Status != PaymentStatusPending {
		return errors.New("only pending payments can be set to processing")
	}
	
	p.Status = PaymentStatusProcessing
	p.UpdatedAt = time.Now()
	return nil
}

// SetRefunded marks the payment as refunded
func (p *Payment) SetRefunded() error {
	if p.Status != PaymentStatusCompleted {
		return errors.New("only completed payments can be refunded")
	}
	
	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	return nil
}

// SetCanceled marks the payment as canceled
func (p *Payment) SetCanceled() error {
	if p.Status == PaymentStatusCompleted || p.Status == PaymentStatusRefunded {
		return errors.New("completed or refunded payments cannot be canceled")
	}
	
	p.Status = PaymentStatusCanceled
	p.UpdatedAt = time.Now()
	return nil
}

// RecordCallback records when a callback was received
func (p *Payment) RecordCallback() {
	p.CallbackTime = time.Now()
	p.UpdatedAt = time.Now()
}

// UpdatePaymentData updates the payment data
func (p *Payment) UpdatePaymentData(data string) {
	p.PaymentData = data
	p.UpdatedAt = time.Now()
}

// PaymentRefund represents a refund transaction
type PaymentRefund struct {
	ID              string    `json:"id"`
	PaymentID       string    `json:"payment_id"`
	OrderID         string    `json:"order_id"`
	UserID          string    `json:"user_id"`
	RefundAmount    float64   `json:"refund_amount"`
	RefundReason    string    `json:"refund_reason"`
	Status          string    `json:"status"` // pending, processing, completed, failed
	TransactionID   string    `json:"transaction_id"`
	OperatorID      string    `json:"operator_id"` // admin user who processed the refund
	OperatorName    string    `json:"operator_name"`
	RefundTime      time.Time `json:"refund_time"`      // when refund was requested
	SuccessTime     time.Time `json:"success_time"`     // when refund was confirmed
	ErrorMsg        string    `json:"error_msg"`        // error message if failed
	ErrorCode       string    `json:"error_code"`       // error code if failed
	PaymentProvider string    `json:"payment_provider"` // provider name (alipay, wxpay, stripe, etc.)
	RefundData      string    `json:"refund_data"`      // serialized refund data (JSON)
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RefundStatusPending represents a pending refund status
const RefundStatusPending = "pending"

// RefundStatusProcessing represents a processing refund status
const RefundStatusProcessing = "processing"

// RefundStatusCompleted represents a completed refund status
const RefundStatusCompleted = "completed"

// RefundStatusFailed represents a failed refund status
const RefundStatusFailed = "failed"

// NewPaymentRefund creates a new payment refund
func NewPaymentRefund(
	id, paymentID, orderID, userID string,
	refundAmount float64,
	refundReason, paymentProvider string,
	operatorID, operatorName string,
) *PaymentRefund {
	now := time.Now()
	return &PaymentRefund{
		ID:              id,
		PaymentID:       paymentID,
		OrderID:         orderID,
		UserID:          userID,
		RefundAmount:    refundAmount,
		RefundReason:    refundReason,
		Status:          RefundStatusPending,
		PaymentProvider: paymentProvider,
		OperatorID:      operatorID,
		OperatorName:    operatorName,
		RefundTime:      now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// SetCompleted marks the refund as completed
func (r *PaymentRefund) SetCompleted(transactionID string) error {
	if r.Status == RefundStatusCompleted {
		return errors.New("refund already completed")
	}
	
	r.Status = RefundStatusCompleted
	r.TransactionID = transactionID
	r.SuccessTime = time.Now()
	r.UpdatedAt = time.Now()
	return nil
}

// SetFailed marks the refund as failed
func (r *PaymentRefund) SetFailed(errorCode, errorMsg string) error {
	if r.Status == RefundStatusCompleted {
		return errors.New("cannot fail a completed refund")
	}
	
	r.Status = RefundStatusFailed
	r.ErrorCode = errorCode
	r.ErrorMsg = errorMsg
	r.UpdatedAt = time.Now()
	return nil
}

// SetProcessing marks the refund as processing
func (r *PaymentRefund) SetProcessing() error {
	if r.Status != RefundStatusPending {
		return errors.New("only pending refunds can be set to processing")
	}
	
	r.Status = RefundStatusProcessing
	r.UpdatedAt = time.Now()
	return nil
}

// UpdateRefundData updates the refund data
func (r *PaymentRefund) UpdateRefundData(data string) {
	r.RefundData = data
	r.UpdatedAt = time.Now()
}
