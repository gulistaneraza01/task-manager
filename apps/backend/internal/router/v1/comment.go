package v1

import (
	"github.com/gulistaneraza01/go-task-manager-echo/internal/handler"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerCommentRoutes(router *echo.Group, handler *handler.CommentHandler, auth *middleware.AuthMiddleware) {

	// Comment routes
	comments := router.Group("/comments")
	comments.Use(auth.RequireAuth)

	// dynamic routes
	dynamicCommentRoute := comments.Group("/:id")
	dynamicCommentRoute.GET("", handler.GetComment)
	dynamicCommentRoute.PATCH("", handler.UpdateComment)
	dynamicCommentRoute.DELETE("", handler.DeleteComment)

}
