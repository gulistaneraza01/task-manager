package category

import (
	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
)

// Category represents a category in the database
type Category struct {
	model.Base
	Name   string    `json:"name" db:"name"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Color  string    `json:"color" db:"color"`
}
