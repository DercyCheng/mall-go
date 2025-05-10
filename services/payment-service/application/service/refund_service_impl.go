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

// RefundServiceImpl implements the RefundService interface
type RefundServiceImpl struct {
	refundRepo       repository.PaymentRefundRepository
	paymentRepo      repository.PaymentRepository
	paymentCache     repository.PaymentCache
	paymentProviders map[string]PaymentProvider
}

// NewRefundServiceImpl creates a new RefundServiceImpl
func NewRefundServiceImpl(
	refundRepo repository.PaymentRefundRepository,
	paymentRepo repository.PaymentRepository,
	paymentCache repository.PaymentCache,
	paymentProviders map[string]PaymentProvider,
) *RefundServiceImpl {
	return &RefundServiceImpl{
		refundRepo:       refundRepo,
		paymentRepo:      paymentRepo,
		paymentCache:     paymentCache,
		paymentProviders: paymentProviders,
	}
}

// CreateRefund creates a new refund request
func (s *RefundServiceImpl) CreateRefund(ctx context.Context, request dto.CreateRefundRequest) (*dto.PaymentRefundDTO, error) {
	// Get payment
	payment, err := s.paymentRepo.GetByID(ctx, request.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Validate payment status
	if payment.Status != model.PaymentStatusCompleted {
		return nil, errors.New("can only refund completed payments")
	}

	// Validate refund amount
	if request.RefundAmount > payment.Amount {
		return nil, errors.New("refund amount cannot exceed payment amount")
	}

	// Check if there are existing refunds
	existingRefunds, err := s.refundRepo.GetByPaymentID(ctx, request.PaymentID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing refunds: %w", err)
	}

	// Calculate total refunded amount
	var totalRefunded float64
	for _, r := range existingRefunds {
		if r.Status == model.RefundStatusCompleted || r.Status == model.RefundStatusProcessing {
			totalRefunded += r.RefundAmount
		}
	}

	// Validate total refund amount
	if totalRefunded+request.RefundAmount > payment.Amount {
		return nil, errors.New("total refund amount exceeds payment amount")
	}

	// Generate refund ID
	refundID := uuid.New().String()

	// Create refund model
	refund := model.NewPaymentRefund(
		refundID,
		payment.ID,
		payment.OrderID,
		payment.UserID,
		request.RefundAmount,
		request.RefundReason,
		payment.PaymentProvider,
		request.OperatorID,
		request.OperatorName,
	)

	// Get payment provider
	provider, exists := s.paymentProviders[payment.PaymentProvider]
	if !exists {
		return nil, fmt.Errorf("unsupported payment provider: %s", payment.PaymentProvider)
	}

	// Initiate refund with provider
	transactionID, err := provider.Refund(ctx, payment, refund)
	if err != nil {
		refund.SetFailed("provider_error", err.Error())
	} else {
		// Set to processing status
		refund.SetProcessing()
		refund.TransactionID = transactionID
	}

	// Save refund
	err = s.refundRepo.Create(ctx, refund)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	// If refund was successful and completed the whole payment amount, update payment status
	if refund.Status == model.RefundStatusCompleted && request.RefundAmount == payment.Amount {
		payment.SetRefunded()
		err = s.paymentRepo.Update(ctx, payment)
		if err != nil {
			// Log the error but continue
			fmt.Printf("failed to update payment status: %v", err)
		} else {
			// Update payment cache
			_ = s.paymentCache.SetPayment(ctx, payment, 3600)
			_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)
		}
	}

	return mapRefundToDTO(refund), nil
}

// GetRefundByID retrieves a refund by ID
func (s *RefundServiceImpl) GetRefundByID(ctx context.Context, id string) (*dto.PaymentRefundDTO, error) {
	refund, err := s.refundRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("refund not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	return mapRefundToDTO(refund), nil
}

// GetRefundsByPaymentID retrieves all refunds for a payment
func (s *RefundServiceImpl) GetRefundsByPaymentID(ctx context.Context, paymentID string) ([]*dto.PaymentRefundDTO, error) {
	refunds, err := s.refundRepo.GetByPaymentID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return []*dto.PaymentRefundDTO{}, nil
		}
		return nil, fmt.Errorf("failed to get refunds: %w", err)
	}

	// Map to DTOs
	refundDTOs := make([]*dto.PaymentRefundDTO, len(refunds))
	for i, refund := range refunds {
		refundDTOs[i] = mapRefundToDTO(refund)
	}

	return refundDTOs, nil
}

// GetRefundsByOrderID retrieves all refunds for an order
func (s *RefundServiceImpl) GetRefundsByOrderID(ctx context.Context, orderID string) ([]*dto.PaymentRefundDTO, error) {
	refunds, err := s.refundRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return []*dto.PaymentRefundDTO{}, nil
		}
		return nil, fmt.Errorf("failed to get refunds: %w", err)
	}

	// Map to DTOs
	refundDTOs := make([]*dto.PaymentRefundDTO, len(refunds))
	for i, refund := range refunds {
		refundDTOs[i] = mapRefundToDTO(refund)
	}

	return refundDTOs, nil
}

// ListRefundsByUserID lists all refunds for a user
func (s *RefundServiceImpl) ListRefundsByUserID(ctx context.Context, userID string, page, pageSize int) ([]*dto.PaymentRefundDTO, int, error) {
	refunds, total, err := s.refundRepo.ListByUserID(ctx, userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list refunds: %w", err)
	}

	// Map to DTOs
	refundDTOs := make([]*dto.PaymentRefundDTO, len(refunds))
	for i, refund := range refunds {
		refundDTOs[i] = mapRefundToDTO(refund)
	}

	return refundDTOs, total, nil
}

// UpdateRefundStatus updates a refund's status
func (s *RefundServiceImpl) UpdateRefundStatus(ctx context.Context, request dto.UpdateRefundStatusRequest) (*dto.PaymentRefundDTO, error) {
	// Get refund
	refund, err := s.refundRepo.GetByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	// Update refund based on status
	switch request.Status {
	case model.RefundStatusCompleted:
		err = refund.SetCompleted(request.TransactionID)
	case model.RefundStatusFailed:
		err = refund.SetFailed(request.ErrorCode, request.ErrorMsg)
	case model.RefundStatusProcessing:
		err = refund.SetProcessing()
	default:
		return nil, fmt.Errorf("invalid refund status: %s", request.Status)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update refund status: %w", err)
	}

	// Save refund
	err = s.refundRepo.Update(ctx, refund)
	if err != nil {
		return nil, fmt.Errorf("failed to save refund: %w", err)
	}

	// If refund was completed, update payment status if needed
	if refund.Status == model.RefundStatusCompleted {
		// Get payment
		payment, err := s.paymentRepo.GetByID(ctx, refund.PaymentID)
		if err == nil {
			// Check if total refunded equals payment amount
			refunds, err := s.refundRepo.GetByPaymentID(ctx, refund.PaymentID)
			if err == nil {
				var totalRefunded float64
				for _, r := range refunds {
					if r.Status == model.RefundStatusCompleted {
						totalRefunded += r.RefundAmount
					}
				}

				if totalRefunded == payment.Amount {
					// Update payment status
					payment.SetRefunded()
					err = s.paymentRepo.Update(ctx, payment)
					if err == nil {
						// Update payment cache
						_ = s.paymentCache.SetPayment(ctx, payment, 3600)
						_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)
					}
				}
			}
		}
	}

	return mapRefundToDTO(refund), nil
}

// ProcessRefundCallback processes a refund callback from a payment provider
func (s *RefundServiceImpl) ProcessRefundCallback(ctx context.Context, request dto.RefundCallbackRequest) (*dto.PaymentCallbackResponse, error) {
	// Get the payment provider
	provider, exists := s.paymentProviders[request.PaymentProvider]
	if !exists {
		return nil, fmt.Errorf("unsupported payment provider: %s", request.PaymentProvider)
	}

	// Verify callback
	refundID, status, transactionID, err := provider.VerifyRefundCallback(ctx, request.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to verify callback: %w", err)
	}

	// Get refund
	refund, err := s.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	// Update refund status based on callback status
	if status == "success" {
		err = refund.SetCompleted(transactionID)
		if err != nil {
			return nil, fmt.Errorf("failed to update refund status: %w", err)
		}
	} else if status == "failed" {
		err = refund.SetFailed("callback_failed", "Refund callback indicates failure")
		if err != nil {
			return nil, fmt.Errorf("failed to update refund status: %w", err)
		}
	}

	// Save refund
	err = s.refundRepo.Update(ctx, refund)
	if err != nil {
		return nil, fmt.Errorf("failed to save refund: %w", err)
	}

	// If refund was completed, update payment status if needed
	if refund.Status == model.RefundStatusCompleted {
		// Get payment
		payment, err := s.paymentRepo.GetByID(ctx, refund.PaymentID)
		if err == nil {
			// Check if total refunded equals payment amount
			refunds, err := s.refundRepo.GetByPaymentID(ctx, refund.PaymentID)
			if err == nil {
				var totalRefunded float64
				for _, r := range refunds {
					if r.Status == model.RefundStatusCompleted {
						totalRefunded += r.RefundAmount
					}
				}

				if totalRefunded == payment.Amount {
					// Update payment status
					payment.SetRefunded()
					err = s.paymentRepo.Update(ctx, payment)
					if err == nil {
						// Update payment cache
						_ = s.paymentCache.SetPayment(ctx, payment, 3600)
						_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)
					}
				}
			}
		}
	}

	return &dto.PaymentCallbackResponse{
		Success: true,
		Message: "Refund callback processed successfully",
	}, nil
}

// QueryRefundStatus queries a refund's status from the payment provider
func (s *RefundServiceImpl) QueryRefundStatus(ctx context.Context, refundID string) (*dto.PaymentRefundDTO, error) {
	// Get refund
	refund, err := s.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	// Get payment provider
	provider, exists := s.paymentProviders[refund.PaymentProvider]
	if !exists {
		return nil, fmt.Errorf("unsupported payment provider: %s", refund.PaymentProvider)
	}

	// Query refund status from provider
	status, transactionID, err := provider.QueryRefundStatus(ctx, refund)
	if err != nil {
		return nil, fmt.Errorf("failed to query refund status: %w", err)
	}

	// Update refund status if changed
	if refund.Status != status {
		if status == model.RefundStatusCompleted {
			err = refund.SetCompleted(transactionID)
		} else if status == model.RefundStatusFailed {
			err = refund.SetFailed("query_failed", "Refund query indicates failure")
		} else if status == model.RefundStatusProcessing {
			err = refund.SetProcessing()
		}

		if err != nil {
			return nil, fmt.Errorf("failed to update refund status: %w", err)
		}

		// Save refund
		err = s.refundRepo.Update(ctx, refund)
		if err != nil {
			return nil, fmt.Errorf("failed to save refund: %w", err)
		}

		// If refund was completed, update payment status if needed
		if refund.Status == model.RefundStatusCompleted {
			// Get payment
			payment, err := s.paymentRepo.GetByID(ctx, refund.PaymentID)
			if err == nil {
				// Check if total refunded equals payment amount
				refunds, err := s.refundRepo.GetByPaymentID(ctx, refund.PaymentID)
				if err == nil {
					var totalRefunded float64
					for _, r := range refunds {
						if r.Status == model.RefundStatusCompleted {
							totalRefunded += r.RefundAmount
						}
					}

					if totalRefunded == payment.Amount {
						// Update payment status
						payment.SetRefunded()
						err = s.paymentRepo.Update(ctx, payment)
						if err == nil {
							// Update payment cache
							_ = s.paymentCache.SetPayment(ctx, payment, 3600)
							_ = s.paymentCache.SetPaymentByOrderID(ctx, payment, 3600)
						}
					}
				}
			}
		}
	}

	return mapRefundToDTO(refund), nil
}

// Helper function to map a refund domain model to a DTO
func mapRefundToDTO(refund *model.PaymentRefund) *dto.PaymentRefundDTO {
	return &dto.PaymentRefundDTO{
		ID:              refund.ID,
		PaymentID:       refund.PaymentID,
		OrderID:         refund.OrderID,
		UserID:          refund.UserID,
		RefundAmount:    refund.RefundAmount,
		RefundReason:    refund.RefundReason,
		Status:          refund.Status,
		TransactionID:   refund.TransactionID,
		OperatorID:      refund.OperatorID,
		OperatorName:    refund.OperatorName,
		RefundTime:      refund.RefundTime,
		SuccessTime:     refund.SuccessTime,
		ErrorMsg:        refund.ErrorMsg,
		ErrorCode:       refund.ErrorCode,
		PaymentProvider: refund.PaymentProvider,
		RefundData:      refund.RefundData,
		CreatedAt:       refund.CreatedAt,
		UpdatedAt:       refund.UpdatedAt,
	}
}
