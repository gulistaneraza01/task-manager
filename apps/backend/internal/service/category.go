package service

import (
	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/category"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/repository"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/labstack/echo/v4"
)

type CategoryService struct {
	server       *server.Server
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(server *server.Server, categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		server:       server,
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) CreateCategory(ctx echo.Context, userID uuid.UUID, payload *category.CreateCategoryPayload) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	CategoryItem, err := s.categoryRepo.CreateCategory(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create category")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_created").
		Str("category_id", CategoryItem.ID.String()).
		Str("name", CategoryItem.Name).
		Str("color", CategoryItem.Color).
		Msg("Category created successfully")

	return CategoryItem, nil
}

func (s *CategoryService) GetCategory(ctx echo.Context, categoryID uuid.UUID, userID uuid.UUID) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	categoryItem, err := s.categoryRepo.GetCategory(ctx.Request().Context(), categoryID, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get category")
		return nil, err
	}

	return categoryItem, nil
}

func (s *CategoryService) GetCategories(ctx echo.Context, userID uuid.UUID, query *category.GetAllCategoryQuery) (*model.PaginatedResponse[category.Category], error) {
	logger := middleware.GetLogger(ctx)

	categories, err := s.categoryRepo.GetCategories(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get categories")
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) UpdateCategory(ctx echo.Context, categoryID uuid.UUID, userID uuid.UUID, payload *category.UpdateCategoryPayload) (*category.Category, error) {
	logger := middleware.GetLogger(ctx)

	_, err := s.categoryRepo.GetCategory(ctx.Request().Context(), categoryID, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get category")
		return nil, err
	}

	categoryItem, err := s.categoryRepo.UpdateCategory(ctx.Request().Context(), userID, categoryID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update category")
		return nil, err
	}

	//Bussiness event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_updated").
		Str("category_id", categoryItem.ID.String()).
		Str("name", categoryItem.Name).
		Str("color", categoryItem.Color).
		Msg("Category updated successfully")

	return categoryItem, nil
}

func (s *CategoryService) DeleteCategory(ctx echo.Context, categoryID uuid.UUID, userID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	_, err := s.categoryRepo.GetCategory(ctx.Request().Context(), categoryID, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get category")
		return err
	}

	err = s.categoryRepo.DeleteCategory(ctx.Request().Context(), userID, categoryID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete category")
		return err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "category_deleted").
		Str("category_id", categoryID.String()).
		Msg("Category deleted successfully")

	return nil
}
