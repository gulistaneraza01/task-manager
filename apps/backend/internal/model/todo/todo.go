package todo

import (
	"time"

	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/category"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/comment"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/document_store"
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

// Todo represents a task in the database
type Todo struct {
	model.Base
	Title       string     `json:"title" db:"title"`
	Description *string    `json:"description" db:"description"`
	Status      Status     `json:"status" db:"status"`
	Priority    Priority   `json:"priority" db:"priority"`
	DueDate     *time.Time `json:"dueDate" db:"due_date"`
	CompletedAt *time.Time `json:"completedAt" db:"completed_at"`
	Metadata    *MetaData  `json:"metadata" db:"metadata"`
	SortOrder   int        `json:"sortOrder" db:"sort_order"`
	CategoryID  *uuid.UUID `json:"categoryId" db:"category_id"`
	ParentID    *uuid.UUID `json:"parentTodoId" db:"parent_todo_id"`
	UserID      uuid.UUID  `json:"userId" db:"user_id"`
}

// MetaData represents the metadata of a task
type MetaData struct {
	Tag        []string   `json:"tag"`
	Reminder   *time.Time `json:"reminder"`
	Color      *string    `json:"color"`
	Difficulty int        `json:"difficulty"`
}

// PopulatedTodo represents a task with its related data
type PopulatedTodo struct {
	Todo
	Category    *category.Category             `json:"category" db:"category"`
	Children    []Todo                         `json:"children" db:"children"`
	Comments    []comment.Comment              `json:"comments" db:"comments"`
	Attachments []document_store.DocumentStore `json:"attachments" db:"attachments"`
}

type TodoStats struct {
	Total     int `json:"total"`
	Draft     int `json:"draft"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
	Archived  int `json:"archived"`
	Overdue   int `json:"overdue"`
}
