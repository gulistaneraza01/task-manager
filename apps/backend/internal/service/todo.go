package service

import (
	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/errs"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/todo"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/repository"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/labstack/echo/v4"
)

type TodoService struct {
	server       *server.Server
	todoRepo     *repository.TodoRepository
	categoryRepo *repository.CategoryRepository
}

func NewTodoService(server *server.Server, todoRepo *repository.TodoRepository,
	categoryRepo *repository.CategoryRepository,
) *TodoService {
	return &TodoService{
		server:       server,
		todoRepo:     todoRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *TodoService) CreateTodo(ctx echo.Context, userID uuid.UUID, payload *todo.CreateTodoPayload) (*todo.Todo, error) {
	logger := middleware.GetLogger(ctx)

	// Validate parent todo exists and belongs to user (if provided)
	if payload.ParentTodoID != nil {
		parentTodo, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, *payload.ParentTodoID)
		if err != nil {
			logger.Error().Err(err).Msg("parent todo validation failed")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}

	}

	// Validate category exists and belongs to user (if provided)
	if payload.CategoryID != nil {
		_, err := s.categoryRepo.GetCategory(ctx.Request().Context(), userID, *payload.CategoryID)
		if err != nil {
			logger.Error().Err(err).Msg("category validation failed")
			return nil, err
		}
	}

	todoItem, err := s.todoRepo.CreateTodo(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create todo")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "todo_created").
		Str("todo_id", todoItem.ID.String()).
		Str("title", todoItem.Title).
		Str("category_id", func() string {
			if todoItem.CategoryID != nil {
				return todoItem.CategoryID.String()
			}
			return ""
		}()).
		Str("priority", string(todoItem.Priority)).
		Msg("Todo created successfully")

	return todoItem, nil
}

func (s *TodoService) GetTodo(ctx echo.Context, userID uuid.UUID, todoID uuid.UUID) (*todo.PopulatedTodo, error) {
	logger := middleware.GetLogger(ctx)

	todoItem, err := s.todoRepo.GetTodoByID(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get todo")
		return nil, err
	}

	return todoItem, nil
}

func (s *TodoService) GetAllTodos(ctx echo.Context, userID uuid.UUID, query *todo.GetAllTodosQuery) (*model.PaginatedResponse[todo.PopulatedTodo], error) {
	logger := middleware.GetLogger(ctx)

	todos, err := s.todoRepo.GetAllTodos(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get todos")
		return nil, err
	}

	return todos, nil
}

func (s *TodoService) UpdateTodo(ctx echo.Context, userID uuid.UUID, todoID uuid.UUID, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {
	logger := middleware.GetLogger(ctx)

	//validate parant
	if payload.ParentTodoID != nil {
		parentTodo, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, *payload.ParentTodoID)
		if err != nil {
			logger.Error().Err(err).Msg("parent todo validation failed")
			return nil, err
		}

		if parentTodo.ID == payload.ID {
			err := errs.NewBadRequestError("Todo cannot be its own parent", false, nil, nil, nil)
			logger.Warn().Msg("todo cannot be its own parent")
			return nil, err
		}

		if !parentTodo.CanHaveChildren() {
			err := errs.NewBadRequestError("Parent todo cannot have children (subtasks can't have subtasks)", false, nil, nil, nil)
			logger.Warn().Msg("parent todo cannot have children")
			return nil, err
		}

		logger.Debug().Msg("parent todo validation passed")
	}

	// Validate category exists and belongs to user (if provided)
	if payload.CategoryID != nil {
		_, err := s.categoryRepo.GetCategory(ctx.Request().Context(), userID, *payload.CategoryID)
		if err != nil {
			logger.Error().Err(err).Msg("category validation failed")
			return nil, err
		}

		logger.Debug().Msg("category validation passed")
	}

	updatedTodo, err := s.todoRepo.UpdateTodo(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update todo")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "todo_updated").
		Str("todo_id", updatedTodo.ID.String()).
		Str("title", updatedTodo.Title).
		Str("category_id", func() string {
			if updatedTodo.CategoryID != nil {
				return updatedTodo.CategoryID.String()
			}
			return ""
		}()).
		Str("priority", string(updatedTodo.Priority)).
		Str("status", string(updatedTodo.Status)).
		Msg("Todo updated successfully")

	return updatedTodo, nil
}

func (s *TodoService) GetTodoStat(ctx echo.Context, userID uuid.UUID) (*todo.TodoStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.todoRepo.GetTodoStat(ctx.Request().Context(), userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch todo statistics")
		return nil, err
	}

	return stats, nil
}

func (s *TodoService) DeleteTodo(ctx echo.Context, userID uuid.UUID, todoID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	// Validate todo exists
	_, err := s.todoRepo.GetTodoByID(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return err
	}

	err = s.todoRepo.DeleteTodo(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete todo")
		return err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "todo_deleted").
		Str("todo_id", todoID.String()).
		Msg("Todo deleted successfully")

	return nil
}
