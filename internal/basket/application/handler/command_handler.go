package handler

import (
	"obs-tools-usage/internal/basket/application/command"
	"obs-tools-usage/internal/basket/application/dto"
	"obs-tools-usage/internal/basket/application/usecase"
)

// CommandHandler handles all commands
type CommandHandler struct {
	basketUseCase *usecase.BasketUseCase
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(basketUseCase *usecase.BasketUseCase) *CommandHandler {
	return &CommandHandler{
		basketUseCase: basketUseCase,
	}
}

// HandleCreateBasket handles CreateBasketCommand
func (h *CommandHandler) HandleCreateBasket(cmd command.CreateBasketCommand) (*dto.BasketResponse, error) {
	return h.basketUseCase.CreateBasket(cmd.UserID)
}

// HandleAddItem handles AddItemCommand
func (h *CommandHandler) HandleAddItem(cmd command.AddItemCommand) (*dto.BasketResponse, error) {
	return h.basketUseCase.AddItem(cmd.UserID, cmd.ProductID, cmd.Quantity)
}

// HandleUpdateItem handles UpdateItemCommand
func (h *CommandHandler) HandleUpdateItem(cmd command.UpdateItemCommand) (*dto.BasketResponse, error) {
	return h.basketUseCase.UpdateItem(cmd.UserID, cmd.ProductID, cmd.Quantity)
}

// HandleRemoveItem handles RemoveItemCommand
func (h *CommandHandler) HandleRemoveItem(cmd command.RemoveItemCommand) (*dto.BasketResponse, error) {
	return h.basketUseCase.RemoveItem(cmd.UserID, cmd.ProductID)
}

// HandleClearBasket handles ClearBasketCommand
func (h *CommandHandler) HandleClearBasket(cmd command.ClearBasketCommand) (*dto.BasketResponse, error) {
	return h.basketUseCase.ClearBasket(cmd.UserID)
}

// HandleDeleteBasket handles DeleteBasketCommand
func (h *CommandHandler) HandleDeleteBasket(cmd command.ClearBasketCommand) error {
	return h.basketUseCase.DeleteBasket(cmd.UserID)
}
