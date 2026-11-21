package service

import (
	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/comment"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/repository"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/labstack/echo/v4"
)

type CommentService struct {
	server      *server.Server
	commentRepo *repository.CommentRepository
	todoRepo    *repository.TodoRepository
}

func NewCommentService(server *server.Server, commentRepo *repository.CommentRepository, todoRepo *repository.TodoRepository) *CommentService {
	return &CommentService{
		server:      server,
		commentRepo: commentRepo,
		todoRepo:    todoRepo,
	}
}

func (s *CommentService) CreateComment(ctx echo.Context, userID string, payload *comment.CreateCommentPayload) (*comment.Comment, error) {

	logger := middleware.GetLogger(ctx)

	//check todoexist

	_, err := s.todoRepo.GetTodoByID(ctx.Request().Context(), userID, payload.TodoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	commentItem, err := s.commentRepo.AddComment(ctx.Request().Context(), userID, payload.TodoID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create comment")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_added").
		Str("comment_id", commentItem.ID.String()).
		Str("todo_id", payload.TodoID.String()).
		Msg("Comment added successfully")

	return commentItem, nil
}

func (s *CommentService) GetComment(ctx echo.Context, userID string, commentID uuid.UUID) (*comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	comment, err := s.commentRepo.GetCommentByID(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment")
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetCommentsByTodoID(ctx echo.Context, userID string, todoID uuid.UUID) ([]comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	// Validate todo exists and belongs to user
	_, err := s.todoRepo.CheckTodoExists(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("todo validation failed")
		return nil, err
	}

	comments, err := s.commentRepo.GetCommentsByTodoID(ctx.Request().Context(), userID, todoID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch comments by todo ID")
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) UpdateComment(ctx echo.Context, userID string, commentID uuid.UUID, payload *comment.UpdateCommentPayload) (*comment.Comment, error) {
	logger := middleware.GetLogger(ctx)

	//validate comment exits
	_, err := s.commentRepo.GetCommentByID(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("comment validation failed")
		return nil, err
	}

	comment, err := s.commentRepo.UpdateComment(ctx.Request().Context(), userID, commentID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update comment")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_updated").
		Str("comment_id", comment.ID.String()).
		Msg("Comment updated successfully")

	return comment, nil
}

func (s *CommentService) DeleteComment(ctx echo.Context, userID string, commentID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	// Validate comment exists
	_, err := s.commentRepo.GetCommentByID(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("comment validation failed")
		return err
	}

	err = s.commentRepo.DeleteComment(ctx.Request().Context(), userID, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete comment")
		return err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "comment_deleted").
		Str("comment_id", commentID.String()).
		Msg("Comment deleted successfully")

	return nil
}
