package comment

import (
	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
)

// Comment represents a comment in the database
type Comment struct {
	model.Base
	TodoID   uuid.UUID  `json:"todo_id" db:"todo_id"`
	UserID   uuid.UUID  `json:"user_id" db:"user_id"`
	Content  string     `json:"content" db:"content"`
	ParentID *uuid.UUID `json:"parent_id" db:"parent_id"`
}
