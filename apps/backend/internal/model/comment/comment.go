package comment

import (
	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
)

// Comment represents a comment in the database
type Comment struct {
	model.Base
	TodoID   uuid.UUID  `json:"todoId" db:"todo_id"`
	UserID   uuid.UUID  `json:"userId" db:"user_id"`
	Content  string     `json:"content" db:"content"`
	ParentID *uuid.UUID `json:"parentId" db:"parent_id"`
}
