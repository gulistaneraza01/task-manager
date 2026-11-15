package todo

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//---------------------------------------------------------------------------------------

// CreateTodoPayload represents the payload for creating a new task
type CreateTodoPayload struct {
	Title        string     `json:"title" validate:"required min=1 max=255"`
	Description  *string    `json:"description" validate:"max=1000"`
	Priority     *Priority  `json:"priority" validate:"omitempty oneof=low medium high"`
	DueDate      *time.Time `json:"dueDate"`
	CategoryID   *uuid.UUID `json:"categoryId" validate:"omitempty uuid"`
	ParentTodoID *uuid.UUID `json:"parentTodoId" validate:"omitempty uuid"`
	Metadata     *MetaData  `json:"metadata"`
}

func (p *CreateTodoPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}

//---------------------------------------------------------------------------------------

// GetAllTodosQuery represents the query parameters for getting all tasks
type GetAllTodosQuery struct {
	Page         *int       `query:"page" validate:"omitempty min=1"`
	Limit        *int       `query:"limit" validate:"omitempty min=1 max=100"`
	Sort         *string    `query:"sort" validate:"omitempty oneof=title dueDate priority createdAt updatedAt"`
	Order        *string    `query:"order" validate:"omitempty oneof=asc desc"`
	Search       *string    `query:"search" validate:"omitempty min=1 max=255"`
	Status       *Status    `query:"status" validate:"omitempty oneof=draft active completed archived"`
	Priority     *Priority  `query:"priority" validate:"omitempty oneof=low medium high"`
	CategoryID   *uuid.UUID `query:"categoryId" validate:"omitempty uuid"`
	ParentTodoID *uuid.UUID `query:"parentTodoId" validate:"omitempty uuid"`
	DueFrom      *time.Time `query:"dueFrom"`
	DueTo        *time.Time `query:"dueTo"`
	Overdue      *bool      `query:"overdue"`
	Completed    *bool      `query:"completed"`
}

func (q *GetAllTodosQuery) Validate() error {
	validation := validator.New()

	if err := validation.Struct(q); err != nil {
		return err
	}

	if q.Page == nil {
		defaultPage := 1
		q.Page = &defaultPage
	}
	if q.Limit == nil {
		defaultLimit := 20
		q.Limit = &defaultLimit
	}

	if q.Sort == nil {
		defaultSort := "createdAt"
		q.Sort = &defaultSort
	}

	if q.Order == nil {
		defaultOrder := "desc"
		q.Order = &defaultOrder
	}

	return nil
}

//---------------------------------------------------------------------------------------

// GetTodoQuery represents the query parameters for getting a single task
type GetTodoQuery struct {
	ID uuid.UUID `param:"id" validate:"required uuid"`
}

func (q *GetTodoQuery) Validate() error {
	validation := validator.New()
	return validation.Struct(q)
}

//---------------------------------------------------------------------------------------

// GetSearchTodoQuery represents the query parameters for searching tasks
type GetSearchTodoQuery struct {
	Query string `query:"query" validate:"required min=1 max=255"`
}

func (q *GetSearchTodoQuery) Validate() error {
	validation := validator.New()
	return validation.Struct(q)
}

//---------------------------------------------------------------------------------------

// GetTodoStatePayload represents the payload for getting the state of a task
type GetTodoStatePayload struct {
}

func (p *GetTodoStatePayload) Validate() error {
	return nil
}

//---------------------------------------------------------------------------------------

// UpdateTodoPayload represents the payload for (patch) updating a task
type UpdateTodoPayload struct {
	ID           uuid.UUID  `param:"id" validate:"required uuid"`
	Title        *string    `json:"title" validate:"omitempty min=1 max=255"`
	Description  *string    `json:"description" validate:"omitempty max=1000"`
	Priority     *Priority  `json:"priority" validate:"omitempty oneof=low medium high"`
	DueDate      *time.Time `json:"dueDate"`
	CategoryID   *uuid.UUID `json:"categoryId" validate:"omitempty uuid"`
	ParentTodoID *uuid.UUID `json:"parentTodoId" validate:"omitempty uuid"`
	Metadata     *MetaData  `json:"metadata"`
	Status       *Status    `json:"status" validate:"omitempty oneof=draft active completed archived"`
}

func (p *UpdateTodoPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}

//---------------------------------------------------------------------------------------

// DeleteTodoPayload represents the payload for deleting a task
type DeleteTodoPayload struct {
	ID uuid.UUID `param:"id" validate:"required uuid"`
}

func (p *DeleteTodoPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}
