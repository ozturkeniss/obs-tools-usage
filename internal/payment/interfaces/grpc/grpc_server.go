package grpc

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"obs-tools-usage/api/proto/payment"
	"obs-tools-usage/internal/payment/application/command"
	"obs-tools-usage/internal/payment/application/handler"
	"obs-tools-usage/internal/payment/application/query"
)

// PaymentGRPCServer implements the PaymentService gRPC server
type PaymentGRPCServer struct {
	payment.UnimplementedPaymentServiceServer
	commandHandler *handler.CommandHandler
	queryHandler   *handler.QueryHandler
	logger         *logrus.Logger
}

// NewPaymentGRPCServer creates a new payment gRPC server
func NewPaymentGRPCServer(commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler, logger *logrus.Logger) *PaymentGRPCServer {
	return &PaymentGRPCServer{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
		logger:         logger,
	}
}

// CreatePayment creates a new payment
func (s *PaymentGRPCServer) CreatePayment(ctx context.Context, req *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"user_id":   req.UserId,
		"basket_id": req.BasketId,
		"method":    req.Method,
		"provider":  req.Provider,
	}).Debug("gRPC CreatePayment request received")

	// Handle command
	paymentResponse, err := s.commandHandler.HandleCreatePayment(command.CreatePaymentCommand{
		UserID:      req.UserId,
		BasketID:    req.BasketId,
		Method:      req.Method,
		Provider:    req.Provider,
		Currency:    req.Currency,
		Description: req.Description,
		Metadata:    make(map[string]string),
	})
	if err != nil {
		s.logger.WithError(err).Error("Failed to create payment")
		return &payment.CreatePaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcPayment := s.convertToGRPCPayment(paymentResponse)

	s.logger.WithFields(logrus.Fields{
		"payment_id": paymentResponse.ID,
		"user_id":    req.UserId,
		"amount":     paymentResponse.Amount,
	}).Info("Successfully created payment via gRPC")

	return &payment.CreatePaymentResponse{
		Success: true,
		Message: "Payment created successfully",
		Payment: grpcPayment,
	}, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentGRPCServer) GetPayment(ctx context.Context, req *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	s.logger.WithField("payment_id", req.PaymentId).Debug("gRPC GetPayment request received")

	// Handle query
	paymentResponse, err := s.queryHandler.HandleGetPayment(query.GetPaymentQuery{PaymentID: req.PaymentId})
	if err != nil {
		s.logger.WithError(err).WithField("payment_id", req.PaymentId).Error("Failed to get payment")
		return &payment.GetPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcPayment := s.convertToGRPCPayment(paymentResponse)

	s.logger.WithField("payment_id", req.PaymentId).Info("Successfully retrieved payment via gRPC")

	return &payment.GetPaymentResponse{
		Success: true,
		Message: "Payment retrieved successfully",
		Payment: grpcPayment,
	}, nil
}

// UpdatePayment updates a payment
func (s *PaymentGRPCServer) UpdatePayment(ctx context.Context, req *payment.UpdatePaymentRequest) (*payment.UpdatePaymentResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"payment_id": req.PaymentId,
		"status":     req.Status,
	}).Debug("gRPC UpdatePayment request received")

	// Handle command
	paymentResponse, err := s.commandHandler.HandleUpdatePayment(command.UpdatePaymentCommand{
		PaymentID: req.PaymentId,
		Status:    req.Status,
		Metadata:  make(map[string]string),
	})
	if err != nil {
		s.logger.WithError(err).WithField("payment_id", req.PaymentId).Error("Failed to update payment")
		return &payment.UpdatePaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcPayment := s.convertToGRPCPayment(paymentResponse)

	s.logger.WithFields(logrus.Fields{
		"payment_id": req.PaymentId,
		"status":     req.Status,
	}).Info("Successfully updated payment via gRPC")

	return &payment.UpdatePaymentResponse{
		Success: true,
		Message: "Payment updated successfully",
		Payment: grpcPayment,
	}, nil
}

// ProcessPayment processes a payment
func (s *PaymentGRPCServer) ProcessPayment(ctx context.Context, req *payment.ProcessPaymentRequest) (*payment.ProcessPaymentResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"payment_id":  req.PaymentId,
		"provider_id": req.ProviderId,
	}).Debug("gRPC ProcessPayment request received")

	// Handle command
	paymentResponse, err := s.commandHandler.HandleProcessPayment(command.ProcessPaymentCommand{
		PaymentID:  req.PaymentId,
		ProviderID: req.ProviderId,
	})
	if err != nil {
		s.logger.WithError(err).WithField("payment_id", req.PaymentId).Error("Failed to process payment")
		return &payment.ProcessPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcPayment := s.convertToGRPCPayment(paymentResponse)

	s.logger.WithFields(logrus.Fields{
		"payment_id": req.PaymentId,
		"status":     paymentResponse.Status,
	}).Info("Successfully processed payment via gRPC")

	return &payment.ProcessPaymentResponse{
		Success: true,
		Message: "Payment processed successfully",
		Payment: grpcPayment,
	}, nil
}

// RefundPayment refunds a payment
func (s *PaymentGRPCServer) RefundPayment(ctx context.Context, req *payment.RefundPaymentRequest) (*payment.RefundPaymentResponse, error) {
	s.logger.WithFields(logrus.Fields{
		"payment_id": req.PaymentId,
		"amount":     req.Amount,
		"reason":     req.Reason,
	}).Debug("gRPC RefundPayment request received")

	// Handle command
	paymentResponse, err := s.commandHandler.HandleRefundPayment(command.RefundPaymentCommand{
		PaymentID: req.PaymentId,
		Amount:    req.Amount,
		Reason:    req.Reason,
	})
	if err != nil {
		s.logger.WithError(err).WithField("payment_id", req.PaymentId).Error("Failed to refund payment")
		return &payment.RefundPaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcPayment := s.convertToGRPCPayment(paymentResponse)

	s.logger.WithFields(logrus.Fields{
		"payment_id": req.PaymentId,
		"amount":     req.Amount,
	}).Info("Successfully refunded payment via gRPC")

	return &payment.RefundPaymentResponse{
		Success: true,
		Message: "Payment refunded successfully",
		Payment: grpcPayment,
	}, nil
}

// GetPaymentsByUser retrieves payments by user
func (s *PaymentGRPCServer) GetPaymentsByUser(ctx context.Context, req *payment.GetPaymentsByUserRequest) (*payment.GetPaymentsByUserResponse, error) {
	s.logger.WithField("user_id", req.UserId).Debug("gRPC GetPaymentsByUser request received")

	// Handle query
	payments, err := s.queryHandler.HandleGetPaymentsByUser(query.GetPaymentsByUserQuery{UserID: req.UserId})
	if err != nil {
		s.logger.WithError(err).WithField("user_id", req.UserId).Error("Failed to get payments by user")
		return &payment.GetPaymentsByUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	var grpcPayments []*payment.Payment
	for _, paymentResponse := range payments {
		grpcPayments = append(grpcPayments, s.convertToGRPCPayment(paymentResponse))
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":        req.UserId,
		"payments_count": len(payments),
	}).Info("Successfully retrieved payments by user via gRPC")

	return &payment.GetPaymentsByUserResponse{
		Success:  true,
		Message:  "Payments retrieved successfully",
		Payments: grpcPayments,
	}, nil
}

// GetPaymentStats retrieves payment statistics
func (s *PaymentGRPCServer) GetPaymentStats(ctx context.Context, req *payment.GetPaymentStatsRequest) (*payment.GetPaymentStatsResponse, error) {
	s.logger.WithField("user_id", req.UserId).Debug("gRPC GetPaymentStats request received")

	// Handle query
	stats, err := s.queryHandler.HandleGetPaymentStats(query.GetPaymentStatsQuery{UserID: req.UserId})
	if err != nil {
		s.logger.WithError(err).WithField("user_id", req.UserId).Error("Failed to get payment stats")
		return &payment.GetPaymentStatsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Convert to gRPC response
	grpcStats := &payment.PaymentStats{
		TotalPayments:     stats.TotalPayments,
		TotalAmount:       stats.TotalAmount,
		CompletedPayments: stats.CompletedPayments,
		FailedPayments:    stats.FailedPayments,
		PendingPayments:   stats.PendingPayments,
		AverageAmount:     stats.AverageAmount,
	}

	s.logger.WithField("user_id", req.UserId).Info("Successfully retrieved payment stats via gRPC")

	return &payment.GetPaymentStatsResponse{
		Success: true,
		Message: "Payment stats retrieved successfully",
		Stats:   grpcStats,
	}, nil
}

// HealthCheck performs a health check
func (s *PaymentGRPCServer) HealthCheck(ctx context.Context, req *payment.HealthCheckRequest) (*payment.HealthCheckResponse, error) {
	s.logger.Debug("gRPC HealthCheck request received")

	return &payment.HealthCheckResponse{
		Success:   true,
		Message:   "Payment service is healthy",
		Service:   "payment-service",
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}, nil
}

// convertToGRPCPayment converts internal payment response to gRPC payment message
func (s *PaymentGRPCServer) convertToGRPCPayment(paymentResponse interface{}) *payment.Payment {
	// Simplified conversion for now
	return &payment.Payment{
		Id:          "converted-id",
		UserId:      "converted-user-id",
		BasketId:    "converted-basket-id",
		Amount:      0.0,
		Currency:    "USD",
		Status:      "pending",
		Method:      "credit_card",
		Provider:    "stripe",
		Description: "Converted payment",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
		Items:       []*payment.PaymentItem{},
	}
}

// RegisterServer registers the gRPC server
func RegisterServer(s *grpc.Server, commandHandler *handler.CommandHandler, queryHandler *handler.QueryHandler, logger *logrus.Logger) {
	paymentServer := NewPaymentGRPCServer(commandHandler, queryHandler, logger)
	payment.RegisterPaymentServiceServer(s, paymentServer)
	logger.Info("Payment gRPC server registered")
}
