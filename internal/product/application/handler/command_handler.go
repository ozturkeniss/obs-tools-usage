package handler

import (
	"obs-tools-usage/internal/product/application/command"
	"obs-tools-usage/internal/product/application/usecase"
	"obs-tools-usage/internal/product/domain/entity"
)

// CommandHandler handles all commands
type CommandHandler struct {
	productUseCase *usecase.ProductUseCase
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(productUseCase *usecase.ProductUseCase) *CommandHandler {
	return &CommandHandler{
		productUseCase: productUseCase,
	}
}

// HandleCreateProduct handles CreateProductCommand
func (h *CommandHandler) HandleCreateProduct(cmd command.CreateProductCommand) (*entity.Product, error) {
	return h.productUseCase.CreateProduct(cmd.ToDTO())
}

// HandleUpdateProduct handles UpdateProductCommand
func (h *CommandHandler) HandleUpdateProduct(cmd command.UpdateProductCommand) (*entity.Product, error) {
	return h.productUseCase.UpdateProduct(cmd.ID, cmd.ToDTO())
}

// HandleDeleteProduct handles DeleteProductCommand
func (h *CommandHandler) HandleDeleteProduct(cmd command.DeleteProductCommand) error {
	return h.productUseCase.DeleteProduct(cmd.ID)
}
