package v1

import (
	"github.com/gulistaneraza01/go-task-manager-echo/internal/handler"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerTodoRoutes(r *echo.Group, h *handler.TodoHandler, ch *handler.CommentHandler, auth *middleware.AuthMiddleware) {
	//task opertaion

	todos := r.Group("/todos")
	todos.Use(auth.RequireAuth)

	//collection opertaion
	todos.GET("", h.GetAllTodos)
	todos.POST("", h.CreateTodo)
	todos.POST("/stat", h.GetTodoStat)

	//item opertaion
	dynamicTodos := todos.Group("/:id")
	dynamicTodos.GET("", h.GetTodo)
	dynamicTodos.PATCH("", h.UpdateTodo)
	dynamicTodos.DELETE("", h.DeleteTodo)

	//task comment
	todoComments := dynamicTodos.Group("/comments")
	todoComments.POST("", ch.CreateComment)
	todoComments.GET("", ch.GetCommentsByTodoID)

}
