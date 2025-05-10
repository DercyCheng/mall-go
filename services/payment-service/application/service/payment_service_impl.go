package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mall-go/services/payment-service/application/dto"
	"mall-go/services/payment-service/domain/model"
	"mall-go/services/payment-service/domain/repository"
)

// PaymentServiceImpl implements the PaymentService interface
type PaymentServiceImpl struct {
	paymentRepo    repository.PaymentRepository
	paymentCache   repository.PaymentCache
	refundRepo     repository.PaymentRefundRepository
	paymentProviders map[string]PaymentProvider
}

// NewPaymentServiceImpl creates a new PaymentServiceImpl
func NewPaymentServiceImpl(
	paymentRepo repository.PaymentRepository,
	paymentCache repository.PaymentCache,
	refundRepo repository.PaymentRefundRepository,
	paymentProviders map[string]PaymentProvider,
) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		paymentRepo:      paymentRepo,
		paymentCache:     paymentCache,
		refundRepo:       refundRepo,
		paymentProviders: paymentProviders,
	}
}

// CreatePayment creates a new payment
func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, request dto.CreatePaymentRequest) (*dto.PaymentDTO, error) {
	// Validate payment provider
	provider, exists := s.paymentProviders[request.PaymentProvider]
	if !exists {
		return nil, fmt.Errorf("unsupported payment provider: %s", request.PaymentProvider)
	}

	// Check if payment for this order already exists
	existing, err := s.paymentRepo.GetByOrderID(ctx, request.OrderID)
	if err == nil && existing != nil {
		// Payment already exists, check if it's still valid
		if !existing.IsExpired() && existing.Status == model.PaymentStatusPending {
			// Return existing payment
			paymentDTO := mapPaymentToDTO(existing)
			
			// Generate payment URL if not already set
			if paymentDTO.PaymentURL == "" {
				paymentURL, err := provider.GeneratePaymentURL(ctx, existing)
				if err == nil {
					paymentDTO.PaymentURL = paymentURL
				}
			}
			
			return paymentDTO, nil
		}
	}

	// Generate payment ID
	paymentID := uuid.New().String()

	// Set expire time
	var expireTime time.Time
	if request.ExpireMinutes > 0 {
		expireTime = time.Now().Add(time.Duration(request.ExpireMinutes) * time.Minute)
	} else {
		// Default to 30 minutes
		expireTime = time.Now().Add(30 * time.Minute)
	}

	// Create payment model
	payment := model.NewPayment(
		paymentID,
		request.OrderID,
		request.UserID,
		request.Amount,
		request.Currency,
		request.PaymentMethod,
		request.PaymentProvider,
		expireTime,
	)
	payment.ClientIP = request.ClientIP

	// Save payment to repository
	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Cache the payment
	_ = s.paymentCache.SetPayment(ctx, payment, 3600) // Cache for 1 hour
	_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)

	// Generate payment URL
	paymentURL, err := provider.GeneratePaymentURL(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment URL: %w", err)
	}

	// Map to DTO
	paymentDTO := mapPaymentToDTO(payment)
	paymentDTO.PaymentURL = paymentURL

	return paymentDTO, nil
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentServiceImpl) GetPaymentByID(ctx context.Context, id string) (*dto.PaymentDTO, error) {
	// Try cache first
	payment, err := s.paymentCache.GetPayment(ctx, id)
	if err == nil {
		return mapPaymentToDTO(payment), nil
	}

	// Fall back to database
	payment, err = s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("payment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Cache the result
	_ = s.paymentCache.SetPayment(ctx, payment, 3600) // Cache for 1 hour

	return mapPaymentToDTO(payment), nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (s *PaymentServiceImpl) GetPaymentByOrderID(ctx context.Context, orderID string) (*dto.PaymentDTO, error) {
	// Try cache first
	payment, err := s.paymentCache.GetPaymentByOrderID(ctx, orderID)
	if err == nil {
		return mapPaymentToDTO(payment), nil
	}

	// Fall back to database
	payment, err = s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("payment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Cache the result
	_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600) // Cache for 1 hour
	_ = s.paymentCache.SetPayment(ctx, payment, 3600)

	return mapPaymentToDTO(payment), nil
}

// GetPaymentByTransactionID retrieves a payment by transaction ID
func (s *PaymentServiceImpl) GetPaymentByTransactionID(ctx context.Context, transactionID string) (*dto.PaymentDTO, error) {
	payment, err := s.paymentRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("payment not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Cache the result
	_ = s.paymentCache.SetPayment(ctx, payment, 3600) // Cache for 1 hour
	_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)

	return mapPaymentToDTO(payment), nil
}

// ListPaymentsByUserID lists all payments for a user
func (s *PaymentServiceImpl) ListPaymentsByUserID(ctx context.Context, userID string, page, pageSize int) ([]*dto.PaymentDTO, int, error) {
	payments, total, err := s.paymentRepo.ListByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}

	// Map to DTOs
	paymentDTOs := make([]*dto.PaymentDTO, len(payments))
	for i, payment := range payments {
		paymentDTOs[i] = mapPaymentToDTO(payment)
	}

	return paymentDTOs, total, nil
}

// UpdatePaymentStatus updates a payment's status
func (s *PaymentServiceImpl) UpdatePaymentStatus(ctx context.Context, request dto.UpdatePaymentStatusRequest) (*dto.PaymentDTO, error) {
	// Get payment
	payment, err := s.paymentRepo.GetByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Update payment based on status
	switch request.Status {
	case model.PaymentStatusCompleted:
		err = payment.SetCompleted(request.TransactionID)
	case model.PaymentStatusFailed:
		err = payment.SetFailed(request.ErrorCode, request.ErrorMsg)
	case model.PaymentStatusProcessing:
		err = payment.SetProcessing()
	case model.PaymentStatusRefunded:
		err = payment.SetRefunded()
	case model.PaymentStatusCanceled:
		err = payment.SetCanceled()
	default:
		return nil, fmt.Errorf("invalid payment status: %s", request.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	// Save payment
	err = s.paymentRepo.Update(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Update cache
	_ = s.paymentCache.SetPayment(ctx, payment, 3600)
	_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)

	return mapPaymentToDTO(payment), nil
}

// ProcessPaymentCallback processes a payment callback from a payment provider
func (s *PaymentServiceImpl) ProcessPaymentCallback(ctx context.Context, request dto.PaymentCallbackRequest) (*dto.PaymentCallbackResponse, error) {
	// Get the payment provider
	provider, exists := s.paymentProviders[request.PaymentProvider]
	if !exists {
		return nil, fmt.Errorf("unsupported payment provider: %s", request.PaymentProvider)
	}

	// Verify callback
	paymentID, orderID, status, transactionID, err := provider.VerifyCallback(ctx, request.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to verify callback: %w", err)
	}

	// Get payment by ID or order ID
	var payment *model.Payment
	if paymentID != "" {
		payment, err = s.paymentRepo.GetByID(ctx, paymentID)
	} else if orderID != "" {
		payment, err = s.paymentRepo.GetByOrderID(ctx, orderID)
	} else {
		return nil, errors.New("callback missing payment ID and order ID")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Record callback time
	payment.RecordCallback()

	// Update payment status based on callback status
	if status == "success" {
		err = payment.SetCompleted(transactionID)
		if err != nil {
			return nil, fmt.Errorf("failed to update payment status: %w", err)
		}
	} else if status == "failed" {
		err = payment.SetFailed("callback_failed", "Payment callback indicates failure")
		if err != nil {
			return nil, fmt.Errorf("failed to update payment status: %w", err)
		}
	}

	// Save payment
	err = s.paymentRepo.Update(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// Update cache
	_ = s.paymentCache.SetPayment(ctx, payment, 3600)
	_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)

	return &dto.PaymentCallbackResponse{
		Success: true,
		Message: "Callback processed successfully",
	}, nil
}

// QueryPaymentStatus queries a payment's status from the payment provider
func (s *PaymentServiceImpl) QueryPaymentStatus(ctx context.Context, request dto.QueryPaymentRequest) (*dto.PaymentDTO, error) {
	// Get payment from repository
	var payment *model.Payment
	var err error

	if request.ID != "" {
		payment, err = s.paymentRepo.GetByID(ctx, request.ID)
	} else if request.OrderID != "" {
		payment, err = s.paymentRepo.GetByOrderID(ctx, request.OrderID)
	} else {
		return nil, errors.New("payment ID or order ID is required")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Get the payment provider
	provider, exists := s.paymentProviders[payment.PaymentProvider]
	if !exists {
		return nil, fmt.Errorf("unsupported payment provider: %s", payment.PaymentProvider)
	}

	// Query payment status from provider
	status, transactionID, err := provider.QueryPaymentStatus(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment status: %w", err)
	}

	// Update payment status if changed
	if payment.Status != status {
		if status == model.PaymentStatusCompleted {
			err = payment.SetCompleted(transactionID)
		} else if status == model.PaymentStatusFailed {
			err = payment.SetFailed("query_failed", "Payment query indicates failure")
		} else if status == model.PaymentStatusProcessing {
			err = payment.SetProcessing()
		}

		if err != nil {
			return nil, fmt.Errorf("failed to update payment status: %w", err)
		}

		// Save payment
		err = s.paymentRepo.Update(ctx, payment)
		if err != nil {
			return nil, fmt.Errorf("failed to save payment: %w", err)
		}

		// Update cache
		_ = s.paymentCache.SetPayment(ctx, payment, 3600)
		_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)
	}

	return mapPaymentToDTO(payment), nil
}

// ListPaymentMethods lists all available payment methods
func (s *PaymentServiceImpl) ListPaymentMethods(ctx context.Context) ([]dto.PaymentMethodInfo, error) {
	// For now, return a hardcoded list of payment methods
	methods := []dto.PaymentMethodInfo{
		{
			Code:        "alipay",
			Name:        "Alipay",
			Icon:        "alipay.png",
			Description: "Pay with Alipay",
			IsEnabled:   true,
		},
		{
			Code:        "wxpay",
			Name:        "WeChat Pay",
			Icon:        "wxpay.png",
			Description: "Pay with WeChat Pay",
			IsEnabled:   true,
		},
		{
			Code:        "creditcard",
			Name:        "Credit Card",
			Icon:        "creditcard.png",
			Description: "Pay with Credit Card",
			IsEnabled:   true,
		},
	}

	return methods, nil
}

// Helper function to map a payment domain model to a DTO
func mapPaymentToDTO(payment *model.Payment) *dto.PaymentDTO {
	return &dto.PaymentDTO{
		ID:              payment.ID,
		OrderID:         payment.OrderID,
		TransactionID:   payment.TransactionID,
		UserID:          payment.UserID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		PaymentMethod:   payment.PaymentMethod,
		PaymentProvider: payment.PaymentProvider,
		Status:          payment.Status,
		PaymentTime:     payment.PaymentTime,
		ExpireTime:      payment.ExpireTime,
		CallbackTime:    payment.CallbackTime,
		SuccessTime:     payment.SuccessTime,
		ErrorMsg:        payment.ErrorMsg,
		ErrorCode:       payment.ErrorCode,
		ClientIP:        payment.ClientIP,
		PaymentData:     payment.PaymentData,
		CreatedAt:       payment.CreatedAt,
		UpdatedAt:       payment.UpdatedAt,
	}
}
