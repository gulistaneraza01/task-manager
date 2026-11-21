package category

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ------------------------------------------------------------------------------

// CreateCategoryPayload represents the payload for creating a new category
type CreateCategoryPayload struct {
	Name  string `json:"name" validate:"required min=1 max=255"`
	Color string `json:"color" validate:"required hexcolor"`
}

func (c *CreateCategoryPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(c)
}

//--------------------------------------------------------------------------------

// GetAllCategoryQuery represents the query for getting all categories
type GetAllCategoryQuery struct {
	Page   *int    `query:"page" validate:"min=1"`
	Limit  *int    `query:"limit" validate:"min=1,max=100"`
	Sort   *string `query:"sort" validate:"omitempty oneof=name color createdAt updatedAt"`
	Order  *string `query:"order" validate:"omitempty oneof=asc desc"`
	Search *string `query:"search" validate:"omitempty min=1 max=255"`
}

func (q *GetAllCategoryQuery) Validate() error {
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
		defaultSort := "name"
		q.Sort = &defaultSort
	}

	if q.Order == nil {
		defaultOrder := "asc"
		q.Order = &defaultOrder
	}
	return nil
}

// ------------------------------------------------------------------------------

// GetCategoryQuery represents the query for getting a category
type GetCategoryQuery struct {
	ID uuid.UUID `param:"id" validate:"required uuid"`
}

func (q *GetCategoryQuery) Validate() error {
	validation := validator.New()
	return validation.Struct(q)
}

// ------------------------------------------------------------------------------

// UpdateCategoryPayload represents the payload for updating a category
type UpdateCategoryPayload struct {
	ID    uuid.UUID `param:"id" validate:"required uuid"`
	Name  *string   `json:"name" validate:"omitempty min=1 max=255"`
	Color *string   `json:"color" validate:"omitempty hexcolor"`
}

func (c *UpdateCategoryPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(c)
}

// ------------------------------------------------------------------------------

// DeleteCategoryPayload represents the payload for deleting a category
type DeleteCategoryPayload struct {
	ID uuid.UUID `param:"id" validate:"required uuid"`
}

func (c *DeleteCategoryPayload) Validate() error {
	validation := validator.New()
	return validation.Struct(c)
}
