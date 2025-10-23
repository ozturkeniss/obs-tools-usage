package handler

import (
	"obs-tools-usage/internal/payment/application/command"
	"obs-tools-usage/internal/payment/application/dto"
	"obs-tools-usage/internal/payment/application/usecase"
)

// CommandHandler handles all commands
type CommandHandler struct {
	paymentUseCase *usecase.PaymentUseCase
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(paymentUseCase *usecase.PaymentUseCase) *CommandHandler {
	return &CommandHandler{
		paymentUseCase: paymentUseCase,
	}
}

// HandleCreatePayment handles CreatePaymentCommand
func (h *CommandHandler) HandleCreatePayment(cmd command.CreatePaymentCommand) (*dto.PaymentResponse, error) {
	return h.paymentUseCase.CreatePayment(
		cmd.UserID,
		cmd.BasketID,
		cmd.Method,
		cmd.Provider,
		cmd.Currency,
		cmd.Description,
		cmd.Metadata,
	)
}

// HandleUpdatePayment handles UpdatePaymentCommand
func (h *CommandHandler) HandleUpdatePayment(cmd command.UpdatePaymentCommand) (*dto.PaymentResponse, error) {
	return h.paymentUseCase.UpdatePayment(
		cmd.PaymentID,
		cmd.Status,
		cmd.Metadata,
	)
}

// HandleProcessPayment handles ProcessPaymentCommand
func (h *CommandHandler) HandleProcessPayment(cmd command.ProcessPaymentCommand) (*dto.PaymentResponse, error) {
	return h.paymentUseCase.ProcessPayment(
		cmd.PaymentID,
		cmd.ProviderID,
	)
}

// HandleRefundPayment handles RefundPaymentCommand
func (h *CommandHandler) HandleRefundPayment(cmd command.RefundPaymentCommand) (*dto.PaymentResponse, error) {
	return h.paymentUseCase.RefundPayment(
		cmd.PaymentID,
		cmd.Amount,
		cmd.Reason,
	)
}
