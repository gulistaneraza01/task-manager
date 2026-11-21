package v1

import (
	"github.com/gulistaneraza01/go-task-manager-echo/internal/handler"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {
	//Task route
	registerTodoRoutes(router, handlers.Todo, handlers.Comment, middleware.Auth)

	//Category route
	registerCategoryRoutes(router, handlers.Category, middleware.Auth)

	//Comment route
	registerCommentRoutes(router, handlers.Comment, middleware.Auth)

}
