package comment

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//------------------------------------------------------------------------------
// Create comment for task

type CreateCommentPayload struct {
	TodoID   uuid.UUID  `json:"todoId" validate:"required,uuid"`
	Content  string     `json:"content" validate:"required,min=1,max=1000"`
	ParentID *uuid.UUID `json:"parentId" validate:"uuid"`
}

func (c *CreateCommentPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(c)
}

// ------------------------------------------------------------------------------

// GetAllCommentsForTaskPayload represents the payload for getting all comments for a task
type GetAllCommentsForTaskPayload struct {
	TodoID uuid.UUID `param:"todoId" validate:"required,uuid"`
}

func (p *GetAllCommentsForTaskPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}

// ------------------------------------------------------------------------------

// GetCommentPayload represents the payload for getting a comment
type GetCommentPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetCommentPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}

// ------------------------------------------------------------------------------

// UpdateCommentPayload represents the payload for updating a comment
type UpdateCommentPayload struct {
	ID      uuid.UUID `param:"id" validate:"required,uuid"`
	Content string    `json:"content" validate:"required,min=1,max=1000"`
}

func (p *UpdateCommentPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}

// ------------------------------------------------------------------------------

// DeleteCommentPayload represents the payload for deleting a comment
type DeleteCommentPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteCommentPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(p)
}
