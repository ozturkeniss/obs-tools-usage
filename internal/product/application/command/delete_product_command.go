package command

// DeleteProductCommand represents a command to delete a product
type DeleteProductCommand struct {
	ID int `json:"id" binding:"required"`
}
