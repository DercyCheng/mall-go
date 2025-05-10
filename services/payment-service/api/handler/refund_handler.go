package handler

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/payment-service/application/dto"
	"mall-go/services/payment-service/application/service"
	"mall-go/services/payment-service/proto/paymentpb"
)

// RefundHandler handles refund-related gRPC requests
type RefundHandler struct {
	refundService service.RefundService
	paymentpb.UnimplementedRefundServiceServer
}

// NewRefundHandler creates a new RefundHandler
func NewRefundHandler(refundService service.RefundService) *RefundHandler {
	return &RefundHandler{
		refundService: refundService,
	}
}

// CreateRefund creates a new refund
func (h *RefundHandler) CreateRefund(ctx context.Context, req *paymentpb.CreateRefundRequest) (*paymentpb.RefundResponse, error) {
	// Validate request
	if req.PaymentId == "" || req.OrderId == "" || req.UserId == "" || req.RefundAmount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "payment ID, order ID, user ID, and positive refund amount are required")
	}

	// Map request to DTO
	createReq := dto.CreateRefundRequest{
		PaymentID:      req.PaymentId,
		OrderID:        req.OrderId,
		UserID:         req.UserId,
		RefundAmount:   req.RefundAmount,
		RefundReason:   req.RefundReason,
		PaymentProvider: req.PaymentProvider,
		OperatorID:     req.OperatorId,
		OperatorName:   req.OperatorName,
	}

	// Call service
	refund, err := h.refundService.CreateRefund(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refund: %v", err)
	}

	// Map response
	return &paymentpb.RefundResponse{
		Success: true,
		Message: "Refund created successfully",
		Refund:  mapToRefundProto(refund),
	}, nil
}

// GetRefund retrieves a refund by ID
func (h *RefundHandler) GetRefund(ctx context.Context, req *paymentpb.GetRefundRequest) (*paymentpb.RefundResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "refund ID is required")
	}

	// Call service
	refund, err := h.refundService.GetRefundByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get refund: %v", err)
	}

	// Map response
	return &paymentpb.RefundResponse{
		Success: true,
		Message: "Refund retrieved successfully",
		Refund:  mapToRefundProto(refund),
	}, nil
}

// GetRefundsByPaymentID retrieves refunds by payment ID
func (h *RefundHandler) GetRefundsByPaymentID(ctx context.Context, req *paymentpb.GetRefundsByPaymentIDRequest) (*paymentpb.RefundsResponse, error) {
	// Validate request
	if req.PaymentId == "" {
		return nil, status.Error(codes.InvalidArgument, "payment ID is required")
	}

	// Call service
	refunds, err := h.refundService.GetRefundsByPaymentID(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get refunds by payment ID: %v", err)
	}

	// Map response
	refundsList := make([]*paymentpb.Refund, len(refunds))
	for i, refund := range refunds {
		refundsList[i] = mapToRefundProto(refund)
	}

	return &paymentpb.RefundsResponse{
		Success: true,
		Message: "Refunds retrieved successfully",
		Refunds: refundsList,
		Total:   int32(len(refunds)),
	}, nil
}

// GetRefundsByOrderIDRequest retrieves refunds by order ID
func (h *RefundHandler) GetRefundsByOrderIDRequest(ctx context.Context, req *paymentpb.GetRefundsByOrderIDRequest) (*paymentpb.RefundsResponse, error) {
	// Validate request
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	// Call service
	refunds, err := h.refundService.GetRefundsByOrderID(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get refunds by order ID: %v", err)
	}

	// Map response
	refundsList := make([]*paymentpb.Refund, len(refunds))
	for i, refund := range refunds {
		refundsList[i] = mapToRefundProto(refund)
	}

	return &paymentpb.RefundsResponse{
		Success: true,
		Message: "Refunds retrieved successfully",
		Refunds: refundsList,
		Total:   int32(len(refunds)),
	}, nil
}

// ListRefunds lists refunds with pagination
func (h *RefundHandler) ListRefunds(ctx context.Context, req *paymentpb.ListRefundsRequest) (*paymentpb.RefundsResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Call service
	refunds, total, err := h.refundService.ListRefundsByUserID(ctx, req.UserId, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list refunds: %v", err)
	}

	// Map response
	refundsList := make([]*paymentpb.Refund, len(refunds))
	for i, refund := range refunds {
		refundsList[i] = mapToRefundProto(refund)
	}

	return &paymentpb.RefundsResponse{
		Success:  true,
		Message:  "Refunds retrieved successfully",
		Refunds:  refundsList,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// UpdateRefundStatus updates a refund's status
func (h *RefundHandler) UpdateRefundStatus(ctx context.Context, req *paymentpb.UpdateRefundStatusRequest) (*paymentpb.RefundResponse, error) {
	// Validate request
	if req.Id == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "refund ID and status are required")
	}

	// Map request to DTO
	updateReq := dto.UpdateRefundStatusRequest{
		ID:            req.Id,
		Status:        req.Status,
		TransactionID: req.TransactionId,
		ErrorCode:     req.ErrorCode,
		ErrorMsg:      req.ErrorMsg,
	}

	// Call service
	refund, err := h.refundService.UpdateRefundStatus(ctx, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update refund status: %v", err)
	}

	// Map response
	return &paymentpb.RefundResponse{
		Success: true,
		Message: "Refund status updated successfully",
		Refund:  mapToRefundProto(refund),
	}, nil
}

// QueryRefundStatus queries a refund's status from the payment provider
func (h *RefundHandler) QueryRefundStatus(ctx context.Context, req *paymentpb.QueryRefundStatusRequest) (*paymentpb.RefundResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "refund ID is required")
	}

	// Call service
	refund, err := h.refundService.QueryRefundStatus(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query refund status: %v", err)
	}

	// Map response
	return &paymentpb.RefundResponse{
		Success: true,
		Message: "Refund status queried successfully",
		Refund:  mapToRefundProto(refund),
	}, nil
}

// ProcessRefundCallback processes a refund callback from a payment provider
func (h *RefundHandler) ProcessRefundCallback(ctx context.Context, req *paymentpb.RefundCallbackRequest) (*paymentpb.PaymentCallbackResponse, error) {
	// Validate request
	if req.PaymentProvider == "" {
		return nil, status.Error(codes.InvalidArgument, "payment provider is required")
	}

	// Convert parameters map
	parameters := make(map[string]interface{})
	if len(req.Parameters) > 0 {
		for k, v := range req.Parameters {
			parameters[k] = v
		}
	} else if req.RawData != "" {
		// Try to parse raw data as JSON
		if err := json.Unmarshal([]byte(req.RawData), &parameters); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid raw data format: %v", err)
		}
	}

	// Map request to DTO
	callbackReq := dto.RefundCallbackRequest{
		PaymentProvider: req.PaymentProvider,
		Parameters:      parameters,
	}

	// Call service
	resp, err := h.refundService.ProcessRefundCallback(ctx, callbackReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process refund callback: %v", err)
	}

	// Map response
	return &paymentpb.PaymentCallbackResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// Helper function to map a DTO refund to a proto refund
func mapToRefundProto(refund *dto.PaymentRefundDTO) *paymentpb.Refund {
	return &paymentpb.Refund{
		Id:              refund.ID,
		PaymentId:       refund.PaymentID,
		OrderId:         refund.OrderID,
		UserId:          refund.UserID,
		RefundAmount:    refund.RefundAmount,
		RefundReason:    refund.RefundReason,
		Status:          refund.Status,
		TransactionId:   refund.TransactionID,
		OperatorId:      refund.OperatorID,
		OperatorName:    refund.OperatorName,
		RefundTime:      formatTime(refund.RefundTime),
		SuccessTime:     formatTime(refund.SuccessTime),
		ErrorMsg:        refund.ErrorMsg,
		ErrorCode:       refund.ErrorCode,
		PaymentProvider: refund.PaymentProvider,
		RefundData:      refund.RefundData,
		CreatedAt:       formatTime(refund.CreatedAt),
		UpdatedAt:       formatTime(refund.UpdatedAt),
	}
}
