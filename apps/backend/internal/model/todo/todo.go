package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
)

// Status represents the status of a task
type Status string

const (
	StatusDraft     Status = "draft"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

// Priority represents the priority of a task
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Task represents a task in the database
type Task struct {
	model.Base
	Title       string     `json:"title" db:"title"`
	Description *string    `json:"description" db:"description"`
	Status      Status     `json:"status" db:"status"`
	Priority    Priority   `json:"priority" db:"priority"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	Metadata    *MetaData  `json:"metadata" db:"metadata"`
	SortOrder   int        `json:"sort_order" db:"sort_order"`
	CategoryID  *uuid.UUID `json:"category_id" db:"category_id"`
	ParentID    *uuid.UUID `json:"parent_id" db:"parent_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
}

// MetaData represents the metadata of a task
type MetaData struct {
	Tag        []string   `json:"tag"`
	Reminder   *time.Time `json:"reminder"`
	Color      *string    `json:"color"`
	Difficulty int        `json:"difficulty"`
}
