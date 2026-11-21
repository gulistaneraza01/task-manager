package v1

import (
	"github.com/gulistaneraza01/go-task-manager-echo/internal/handler"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerCategoryRoutes(router *echo.Group, handler *handler.CategoryHandler, auth *middleware.AuthMiddleware) {
	//Category opertion

	category := router.Group("/categories")
	category.Use(auth.RequireAuth)

	//collection operation
	category.POST("", handler.CreateCategory)
	category.GET("", handler.GetCategories)

	//item operation
	dynamicCategory := category.Group("/:id")
	dynamicCategory.GET("", handler.GetCategory)
	dynamicCategory.PATCH("", handler.UpdateCategory)
	dynamicCategory.DELETE("", handler.DeleteCategory)

}
