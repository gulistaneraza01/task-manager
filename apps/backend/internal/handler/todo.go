package handler

import (
	"net/http"

	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/todo"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/service"
	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	Handler
	todoService *service.TodoService
}

func NewTodoHandler(s *server.Server, todoSer *service.TodoService) *TodoHandler {
	return &TodoHandler{
		Handler:     NewHandler(s),
		todoService: todoSer,
	}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.CreateTodoPayload) (*todo.Todo, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.CreateTodo(c, userID, payload)
		},
		http.StatusCreated,
		&todo.CreateTodoPayload{},
	)(c)
}

func (h *TodoHandler) GetTodo(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.GetTodoQuery) (*todo.PopulatedTodo, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.GetTodo(c, userID, payload.ID)
		},
		http.StatusOK,
		&todo.GetTodoQuery{},
	)(c)
}

func (h *TodoHandler) GetAllTodos(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *todo.GetAllTodosQuery) (*model.PaginatedResponse[todo.PopulatedTodo], error) {
			userID := middleware.GetUserID(c)
			return h.todoService.GetAllTodos(c, userID, query)
		},
		http.StatusOK,
		&todo.GetAllTodosQuery{},
	)(c)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.UpdateTodoPayload) (*todo.Todo, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.UpdateTodo(c, userID, payload.ID, payload)
		},
		http.StatusOK,
		&todo.UpdateTodoPayload{},
	)(c)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *todo.DeleteTodoPayload) error {
			userID := middleware.GetUserID(c)
			return h.todoService.DeleteTodo(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&todo.DeleteTodoPayload{},
	)(c)
}

func (h *TodoHandler) GetTodoStat(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *todo.GetTodoStatePayload) (*todo.TodoStats, error) {
			userID := middleware.GetUserID(c)
			return h.todoService.GetTodoStat(c, userID)
		},
		http.StatusOK,
		&todo.GetTodoStatePayload{},
	)(c)
}
