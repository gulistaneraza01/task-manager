package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/errs"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/todo"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/jackc/pgx/v5"
)

type TodoRepository struct {
	server *server.Server
}

func NewTodoRepository(server *server.Server) *TodoRepository {
	return &TodoRepository{server: server}
}

// CreateTodo creates a new task item
func (r *TodoRepository) CreateTodo(ctx context.Context, userId uuid.UUID, payload *todo.CreateTodoPayload) (*todo.Todo, error) {

	stmt := `INSERT INTO todos
	(user_id, title, description, priority, due_date, category_id, parent_todo_id, metadata)
	VALUES (@user_id, @title, @description, @priority, @due_date, @category_id, @parent_todo_id, @metadata)
	RETURNING *`

	priority := todo.PriorityMedium

	if payload.Priority != nil {
		priority = *payload.Priority
	}

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{"user_id": userId, "title": payload.Title, "description": payload.Description, "priority": priority, "due_date": payload.DueDate, "category_id": payload.CategoryID, "parent_todo_id": payload.ParentTodoID, "metadata": payload.Metadata})

	if err != nil {
		return nil, fmt.Errorf("failed to create todo query user_id: %s title: %s %w", userId, payload.Title, err)
	}

	todoItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[todo.Todo])
	if err != nil {
		return nil, fmt.Errorf("f ailed to collect todo items user_id: %s title: %s %w", userId, payload.Title, err)
	}

	return &todoItem, nil
}

// GetTodoByID retrieves a task item by its ID
func (r *TodoRepository) GetTodoByID(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) (*todo.PopulatedTodo, error) {

	stmt := `SELECT t.*, CASE WHEN c.id IS NOT NULL THEN to_jsonb(camel(c)) ELSE NULL END AS category,
		COALESCE(jsonb_agg(to_jsonb(camel(child)) ORDER BY child.sort_order ASC, child.created_at ASC) FILTER (WHERE child.id IS NOT NULL), '[]'::JSONB) AS children,
		COALESCE(jsonb_agg(to_jsonb(camel(com)) ORDER BY com.created_at ASC) FILTER (com.id IS NOT NULL), '[]'::JSONB) AS comments,
		COALESCE(jsonb_agg(to_jsonb(camel(ds)) ORDER BY ds.created_at DESC) FILTER (ds.id IS NOT NULL), '[]'::JSONB) AS attachments
	FROM
		todos t
		LEFT JOIN categories c ON c.id = t.category_id AND c.user_id = @user_id
		LEFT JOIN todos child ON child.parent_todo_id = t.id AND child.user_id = @user_id
		LEFT JOIN comments com ON com.todo_id = t.id AND com.user_id = @user_id
		LEFT JOIN document_stores ds ON ds.todo_id = t.id AND ds.user_id = @user_id
	WHERE
		t.id = @id AND t.user_id = @user_id
	GROUP BY
		t.id,
		c.id,
		ds.id`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{"id": todoID, "user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get todo by id user_id: %s todo_id: %s %w", userID, todoID, err)
	}

	todoItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[todo.PopulatedTodo])
	if err != nil {
		return nil, fmt.Errorf("failed to collect todo item user_id: %s todo_id: %s %w", userID, todoID, err)
	}

	return &todoItem, nil
}

// CheckTodoExists checks if a task exists for a user
func (r *TodoRepository) CheckTodoExists(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) (*todo.Todo, error) {

	stmt := `SELECT * FROM todos WHERE id = @id AND user_id = @user_id`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{"id": todoID, "user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to check if todo exists user_id: %s todo_id: %s %w", userID, todoID, err)
	}

	todoItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[todo.Todo])
	if err != nil {
		return nil, fmt.Errorf("failed to collect todo user_id: %s todo_id: %s %w", userID, todoID, err)
	}

	return &todoItem, nil
}

// GetAllTodos retrieves all tasks for a user
func (r *TodoRepository) GetAllTodos(ctx context.Context, userID uuid.UUID, query *todo.GetAllTodosQuery) (*model.PaginatedResponse[todo.PopulatedTodo], error) {
	stmt := `
	SELECT
		t.*,
		CASE WHEN c.id IS NOT NULL THEN to_jsonb(camel (c)) ELSE NULL END AS category,
		COALESCE(jsonb_agg(	to_jsonb(camel (child))	ORDER BY child.sort_order ASC, child.created_at ASC) FILTER (	WHERE	child.id IS NOT NULL),'[]'::JSONB) AS children,
		COALESCE(jsonb_agg(	to_jsonb(camel (com))	ORDER BY com.created_at ASC) FILTER (	WHERE	com.id IS NOT NULL),'[]'::JSONB) AS comments,
		COALESCE(jsonb_agg(to_jsonb(camel (ds)) ORDER BY ds.created_at DESC) FILTER (WHERE ds.id IS NOT NULL), '[]'::JSONB) AS attachments
	FROM
		todos t
		LEFT JOIN categories c ON c.id=t.category_id AND c.user_id=@user_id
		LEFT JOIN todos child ON child.parent_todo_id=t.id AND child.user_id=@user_id
		LEFT JOIN comments com ON com.todo_id=t.id AND com.user_id=@user_id
		LEFT JOIN document_stores ds ON ds.todo_id=t.id AND ds.user_id=@user_id
`
	args := pgx.NamedArgs{"user_id": userID}

	conditions := []string{"t.user_id = @user_id"}

	if query.Status != nil {
		conditions = append(conditions, "t.status = @status")
		args["status"] = *query.Status
	}

	if query.Priority != nil {
		conditions = append(conditions, "t.priority = @priority")
		args["priority"] = *query.Priority
	}

	if query.CategoryID != nil {
		conditions = append(conditions, "t.category_id = @category_id")
		args["category_id"] = *query.CategoryID
	}

	if query.ParentTodoID != nil {
		conditions = append(conditions, "t.parent_todo_id = @parent_todo_id")
		args["parent_todo_id"] = *query.ParentTodoID
	} else {
		conditions = append(conditions, "t.parent_todo_id IS NULL")
	}

	if query.DueFrom != nil {
		conditions = append(conditions, "t.due_date >= @due_from")
		args["due_from"] = *query.DueFrom
	}

	if query.DueTo != nil {
		conditions = append(conditions, "t.due_date <= @due_to")
		args["due_to"] = *query.DueTo
	}

	if query.Overdue != nil {
		conditions = append(conditions, "t.due_date < NOW() AND t.status != 'completed'")
	}

	if query.Completed != nil {
		if *query.Completed {
			conditions = append(conditions, "t.status = 'completed'")
		} else {
			conditions = append(conditions, "t.status != 'completed'")
		}
	}

	if query.Search != nil {
		conditions = append(conditions, "(t.title ILIKE @search OR COALESCE(t.description, '') ILIKE @search)")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	countstmt := "SELECT COUNT(*) FROM todos t WHERE t.user_id = @user_id"
	if len(conditions) > 0 {
		countstmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.server.DB.Pool.QueryRow(ctx, countstmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total todos user_id: %s %w", userID, err)
	}

	orderBy := "t.created_at DESC"
	if query.Sort != nil {
		var sortCol string
		switch *query.Sort {
		case "title":
			sortCol = "t.title"
		case "dueDate":
			sortCol = "t.due_date"
		case "priority":
			sortCol = "t.priority"
		case "createdAt":
			sortCol = "t.created_at"
		case "updatedAt":
			sortCol = "t.updated_at"
		}
		if sortCol != "" {
			order := "ASC"
			if query.Order != nil && strings.EqualFold(*query.Order, "desc") {
				order = "DESC"
			}
			orderBy = sortCol + " " + order
		}
	}

	stmt += " GROUP BY t.id,c.id ORDER BY " + orderBy

	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * *query.Limit

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos user_id: %s %w", userID, err)
	}

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[todo.PopulatedTodo])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[todo.PopulatedTodo]{
				Data:       []todo.PopulatedTodo{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect todos user_id: %s %w", userID, err)
	}

	return &model.PaginatedResponse[todo.PopulatedTodo]{
		Data:       todos,
		Total:      total,
		Page:       *query.Page,
		Limit:      *query.Limit,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

// UpdateTodo updates a todo by ID and user ID
func (r *TodoRepository) UpdateTodo(ctx context.Context, userID uuid.UUID, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {
	stmt := `UPDATE todos SET`
	args := pgx.NamedArgs{
		"user_id": userID,
		"todo_id": payload.ID,
	}

	setClauses := []string{}

	if payload.Title != nil {
		setClauses = append(setClauses, "title = @title")
		args["title"] = *payload.Title
	}

	if payload.Description != nil {
		setClauses = append(setClauses, "description = @description")
		args["description"] = *payload.Description
	}

	if payload.Priority != nil {
		setClauses = append(setClauses, "priority = @priority")
		args["priority"] = *payload.Priority
	}

	if payload.DueDate != nil {
		setClauses = append(setClauses, "due_date = @due_date")
		args["due_date"] = *payload.DueDate
	}

	if payload.CategoryID != nil {
		setClauses = append(setClauses, "category_id = @category_id")
		args["category_id"] = *payload.CategoryID
	}

	if payload.ParentTodoID != nil {
		setClauses = append(setClauses, "parent_todo_id = @parent_todo_id")
		args["parent_todo_id"] = *payload.ParentTodoID
	}

	if payload.Metadata != nil {
		setClauses = append(setClauses, "metadata = @metadata")
		args["metadata"] = *payload.Metadata
	}

	if payload.Status != nil {
		setClauses = append(setClauses, "status = @status")
		args["status"] = *payload.Status

		//auto-set completed_at when status is completed
		if *payload.Status == todo.StatusCompleted {
			setClauses = append(setClauses, "completed_at = @completed_at")
			args["completed_at"] = time.Now()
		} else {
			setClauses = append(setClauses, "completed_at = NULL")
		}
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("no field to update", false, nil, nil, nil)
	}

	stmt += " " + strings.Join(setClauses, ", ")
	stmt += " WHERE id = @todo_id AND user_id = @user_id RETURNING *"

	row, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo user_id: %s %w", userID, err)
	}

	updatedTodo, err := pgx.CollectOneRow(row, pgx.RowToStructByName[todo.Todo])
	if err != nil {
		return nil, fmt.Errorf("failed to collect updated todo: %w", err)
	}

	return &updatedTodo, nil
}

// DeleteTodo deletes a todo by ID and user ID
func (r *TodoRepository) DeleteTodo(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) error {
	stmt := `DELETE FROM todos WHERE id = @todo_id AND user_id = @user_id`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{"user_id": userID, "todo_id": todoID})
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	if result.RowsAffected() == 0 {
		code := "TODO_NOT_FOUND"
		return errs.NewNotFoundError("todo not found", false, &code)
	}

	return nil
}

// GetTodoStat retrieves task statistics for a user
func (r *TodoRepository) GetTodoStat(ctx context.Context, userID uuid.UUID) (*todo.TodoStats, error) {
	stmt := `
		SELECT
			COUNT(*) AS total,
			COUNT(CASE WHEN status='draft' THEN 1 END) AS draft,
			COUNT(CASE WHEN status='active' THEN 1 END) AS active,
			COUNT(CASE WHEN status='completed' THEN 1 END) AS completed,
			COUNT(CASE WHEN status='archived' THEN 1 END) AS archived,
			COUNT(CASE WHEN due_date<NOW() AND status!='completed' THEN 1 END) AS overdue
		FROM
			todos
		WHERE
			user_id=@user_id
	`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get todo stats user_id: %s %w", userID, err)
	}

	stats, err := pgx.CollectOneRow(row, pgx.RowToStructByName[todo.TodoStats])
	if err != nil {
		return nil, fmt.Errorf("failed to collect todo stats user_id: %s %w", userID, err)
	}

	return &stats, nil
}

// IsTodoExists checks if a todo exists for a user
func (r *TodoRepository) IsTodoExists(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) (bool, error) {
	stmt := `SELECT EXISTS (
		SELECT 1 FROM todos
		WHERE id = @id AND user_id = @user_id
	)`

	var exists bool
	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"id":      todoID,
		"user_id": userID,
	}).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check todo existence: %w", err)
	}

	return exists, nil
}
