package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mall-go/services/payment-service/application/dto"
	"mall-go/services/payment-service/application/service"
	"mall-go/services/payment-service/proto/paymentpb"
)

// PaymentHandler handles payment-related gRPC requests
type PaymentHandler struct {
	paymentService service.PaymentService
	paymentpb.UnimplementedPaymentServiceServer
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePayment creates a new payment
func (h *PaymentHandler) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.PaymentResponse, error) {
	// Validate request
	if req.OrderId == "" || req.UserId == "" || req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order ID, user ID and positive amount are required")
	}

	// Map request to DTO
	createReq := dto.CreatePaymentRequest{
		OrderID:         req.OrderId,
		UserID:          req.UserId,
		Amount:          req.Amount,
		Currency:        req.Currency,
		PaymentMethod:   req.PaymentMethod,
		PaymentProvider: req.PaymentProvider,
		ClientIP:        req.ClientIp,
		ExpireMinutes:   int(req.ExpireMinutes),
	}

	// Set defaults if needed
	if createReq.Currency == "" {
		createReq.Currency = "CNY"
	}

	// Call service
	payment, err := h.paymentService.CreatePayment(ctx, createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	// Map response
	return &paymentpb.PaymentResponse{
		Success: true,
		Message: "Payment created successfully",
		Payment: mapToPaymentProto(payment),
	}, nil
}

// GetPayment retrieves a payment by ID
func (h *PaymentHandler) GetPayment(ctx context.Context, req *paymentpb.GetPaymentRequest) (*paymentpb.PaymentResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "payment ID is required")
	}

	// Call service
	payment, err := h.paymentService.GetPaymentByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payment: %v", err)
	}

	// Map response
	return &paymentpb.PaymentResponse{
		Success: true,
		Message: "Payment retrieved successfully",
		Payment: mapToPaymentProto(payment),
	}, nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (h *PaymentHandler) GetPaymentByOrderID(ctx context.Context, req *paymentpb.GetPaymentByOrderIDRequest) (*paymentpb.PaymentResponse, error) {
	// Validate request
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	// Call service
	payment, err := h.paymentService.GetPaymentByOrderID(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payment by order ID: %v", err)
	}

	// Map response
	return &paymentpb.PaymentResponse{
		Success: true,
		Message: "Payment retrieved successfully",
		Payment: mapToPaymentProto(payment),
	}, nil
}

// ListPayments lists payments with pagination
func (h *PaymentHandler) ListPayments(ctx context.Context, req *paymentpb.ListPaymentsRequest) (*paymentpb.PaymentsResponse, error) {
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
	payments, total, err := h.paymentService.ListPaymentsByUserID(ctx, req.UserId, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list payments: %v", err)
	}

	// Map response
	paymentsList := make([]*paymentpb.Payment, len(payments))
	for i, payment := range payments {
		paymentsList[i] = mapToPaymentProto(payment)
	}

	return &paymentpb.PaymentsResponse{
		Success:   true,
		Message:   "Payments retrieved successfully",
		Payments:  paymentsList,
		Total:     int32(total),
		Page:      int32(page),
		PageSize:  int32(pageSize),
	}, nil
}

// UpdatePaymentStatus updates a payment's status
func (h *PaymentHandler) UpdatePaymentStatus(ctx context.Context, req *paymentpb.UpdatePaymentStatusRequest) (*paymentpb.PaymentResponse, error) {
	// Validate request
	if req.Id == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "payment ID and status are required")
	}

	// Map request to DTO
	updateReq := dto.UpdatePaymentStatusRequest{
		ID:            req.Id,
		Status:        req.Status,
		TransactionID: req.TransactionId,
		ErrorCode:     req.ErrorCode,
		ErrorMsg:      req.ErrorMsg,
	}

	// Call service
	payment, err := h.paymentService.UpdatePaymentStatus(ctx, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update payment status: %v", err)
	}

	// Map response
	return &paymentpb.PaymentResponse{
		Success: true,
		Message: "Payment status updated successfully",
		Payment: mapToPaymentProto(payment),
	}, nil
}

// QueryPaymentStatus queries a payment's status from the payment provider
func (h *PaymentHandler) QueryPaymentStatus(ctx context.Context, req *paymentpb.QueryPaymentStatusRequest) (*paymentpb.PaymentResponse, error) {
	// Validate request
	if req.Id == "" && req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "either payment ID or order ID is required")
	}

	// Map request to DTO
	queryReq := dto.QueryPaymentRequest{
		ID:      req.Id,
		OrderID: req.OrderId,
	}

	// Call service
	payment, err := h.paymentService.QueryPaymentStatus(ctx, queryReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query payment status: %v", err)
	}

	// Map response
	return &paymentpb.PaymentResponse{
		Success: true,
		Message: "Payment status queried successfully",
		Payment: mapToPaymentProto(payment),
	}, nil
}

// ProcessPaymentCallback processes a payment callback from a payment provider
func (h *PaymentHandler) ProcessPaymentCallback(ctx context.Context, req *paymentpb.PaymentCallbackRequest) (*paymentpb.PaymentCallbackResponse, error) {
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
	callbackReq := dto.PaymentCallbackRequest{
		PaymentProvider: req.PaymentProvider,
		Parameters:      parameters,
	}

	// Call service
	resp, err := h.paymentService.ProcessPaymentCallback(ctx, callbackReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment callback: %v", err)
	}

	// Map response
	return &paymentpb.PaymentCallbackResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// ListPaymentMethods lists all available payment methods
func (h *PaymentHandler) ListPaymentMethods(ctx context.Context, req *paymentpb.ListPaymentMethodsRequest) (*paymentpb.PaymentMethodsResponse, error) {
	// Call service
	methods, err := h.paymentService.ListPaymentMethods(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list payment methods: %v", err)
	}

	// Map response
	methodsList := make([]*paymentpb.PaymentMethod, len(methods))
	for i, method := range methods {
		methodsList[i] = &paymentpb.PaymentMethod{
			Code:        method.Code,
			Name:        method.Name,
			Icon:        method.Icon,
			Description: method.Description,
			IsEnabled:   method.IsEnabled,
		}
	}

	return &paymentpb.PaymentMethodsResponse{
		Success: true,
		Message: "Payment methods retrieved successfully",
		Methods: methodsList,
	}, nil
}

// Helper function to map a DTO payment to a proto payment
func mapToPaymentProto(payment *dto.PaymentDTO) *paymentpb.Payment {
	return &paymentpb.Payment{
		Id:              payment.ID,
		OrderId:         payment.OrderID,
		TransactionId:   payment.TransactionID,
		UserId:          payment.UserID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		PaymentMethod:   payment.PaymentMethod,
		PaymentProvider: payment.PaymentProvider,
		Status:          payment.Status,
		PaymentTime:     formatTime(payment.PaymentTime),
		ExpireTime:      formatTime(payment.ExpireTime),
		CallbackTime:    formatTime(payment.CallbackTime),
		SuccessTime:     formatTime(payment.SuccessTime),
		ErrorMsg:        payment.ErrorMsg,
		ErrorCode:       payment.ErrorCode,
		ClientIp:        payment.ClientIP,
		PaymentUrl:      payment.PaymentURL,
		PaymentData:     payment.PaymentData,
		CreatedAt:       formatTime(payment.CreatedAt),
		UpdatedAt:       formatTime(payment.UpdatedAt),
	}
}

// formatTime formats a time.Time value as string
func formatTime(t interface{}) string {
	if t == nil {
		return ""
	}
	
	switch v := t.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		// Try to use String() method
		if str, ok := t.(interface{ String() string }); ok {
			return str.String()
		}
		return ""
	}
}
