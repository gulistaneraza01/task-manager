package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/errs"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/category"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	server *server.Server
}

func NewCategoryRepository(server *server.Server) *CategoryRepository {
	return &CategoryRepository{server: server}
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(ctx context.Context, userID string, payload *category.CreateCategoryPayload) (*category.Category, error) {
	stmt := `
	    INSERT INTO categories (name, user_id, color)
	    VALUES (@name, @user_id, @color)
	    RETURNING *
	`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":    payload.Name,
		"user_id": userID,
		"color":   payload.Color,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	categoryItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[category.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to collect category: %w", err)
	}

	return &categoryItem, nil

}

// GetCategory gets a category by id
func (r *CategoryRepository) GetCategory(ctx context.Context, categoryID uuid.UUID, userID string) (*category.Category, error) {
	stmt := `SELECT * FROM categories WHERE id = @category_id AND user_id = @user_id`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"category_id": categoryID,
		"user_id":     userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	categoryItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[category.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to collect category: %w", err)
	}

	return &categoryItem, nil
}

// GetCategories gets all categories for a user
func (r *CategoryRepository) GetCategories(ctx context.Context, userID string, payload *category.GetAllCategoryQuery) (*model.PaginatedResponse[category.Category], error) {
	stmt := `SELECT * FROM categories`
	args := pgx.NamedArgs{"user_id": userID}

	conditions := []string{"user_id = @user_id"}

	if payload.Search != nil {
		conditions = append(conditions, "name ILIKE '%' || @search || '%'")
		args["search"] = *payload.Search
	}

	orderBy := "name ASC"
	if payload.Sort != nil {
		var sortCol string
		switch *payload.Sort {
		case "name":
			sortCol = "name"
		case "color":
			sortCol = "color"
		case "createdAt":
			sortCol = "created_at"
		case "updatedAt":
			sortCol = "updated_at"
		}
		if sortCol != "" {
			order := "ASC"
			if payload.Order != nil && strings.EqualFold(*payload.Order, "desc") {
				order = "DESC"
			}
			orderBy = sortCol + " " + order
		}
	}

	stmt += " WHERE " + strings.Join(conditions, " AND ")

	countstmt := "SELECT COUNT(*) FROM categories WHERE " + strings.Join(conditions, " AND ")

	var total int
	err := r.server.DB.Pool.QueryRow(ctx, countstmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total categories user_id: %s %w", userID, err)
	}

	stmt += " ORDER BY " + orderBy + " LIMIT @limit OFFSET @offset"
	args["limit"] = *payload.Limit
	args["offset"] = (*payload.Page - 1) * *payload.Limit

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[category.Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[category.Category]{
				Data:       []category.Category{},
				Page:       *payload.Page,
				Limit:      *payload.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:todo_categories for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[category.Category]{
		Data:       categories,
		Page:       *payload.Page,
		Limit:      *payload.Limit,
		Total:      total,
		TotalPages: (total + *payload.Limit - 1) / *payload.Limit,
	}, nil
}

// UpdateCategory updates a category by ID and user ID
func (r *CategoryRepository) UpdateCategory(ctx context.Context, userID string, categoryID uuid.UUID, payload *category.UpdateCategoryPayload) (*category.Category, error) {
	stmt := `UPDATE categories SET`
	args := pgx.NamedArgs{
		"user_id": userID,
		"id":      categoryID,
	}

	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}

	if payload.Color != nil {
		setClauses = append(setClauses, "color = @color")
		args["color"] = *payload.Color
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	stmt += " " + strings.Join(setClauses, ", ") + " WHERE id = @id AND user_id = @user_id RETURNING *"

	row, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	category, err := pgx.CollectOneRow(row, pgx.RowToStructByName[category.Category])
	if err != nil {
		return nil, fmt.Errorf("failed to collect updated category: %w", err)
	}

	return &category, nil
}

// DeleteCategory deletes a category by ID and user ID
func (r *CategoryRepository) DeleteCategory(ctx context.Context, userID string, categoryID uuid.UUID) error {
	stmt := `DELETE FROM categories WHERE id = @id AND user_id = @user_id`

	row, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"id":      categoryID,
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if row.RowsAffected() == 0 {
		code := "CATEGORY_NOT_FOUND"
		return errs.NewNotFoundError("category not found", false, &code)
	}

	return nil
}
