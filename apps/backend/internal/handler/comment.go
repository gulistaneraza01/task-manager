package handler

import (
	"net/http"

	"github.com/gulistaneraza01/go-task-manager-echo/internal/middleware"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/comment"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/service"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	Handler
	commentService *service.CommentService
}

func NewCommentHandler(s *server.Server, commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		Handler:        NewHandler(s),
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *comment.CreateCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.CreateComment(c, userID, req)
		},
		http.StatusCreated,
		&comment.CreateCommentPayload{},
	)(c)
}

func (h *CommentHandler) GetComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *comment.GetCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.GetComment(c, userID, req.ID)
		},
		http.StatusOK,
		&comment.GetCommentPayload{},
	)(c)
}

func (h *CommentHandler) GetCommentsByTodoID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, req *comment.GetAllCommentsForTaskPayload) ([]comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.GetCommentsByTodoID(c, userID, req.TodoID)
		},
		http.StatusOK,
		&comment.GetAllCommentsForTaskPayload{},
	)(c)
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *comment.UpdateCommentPayload) (*comment.Comment, error) {
			userID := middleware.GetUserID(c)
			return h.commentService.UpdateComment(c, userID, payload.ID, payload)
		},
		http.StatusOK,
		&comment.UpdateCommentPayload{},
	)(c)
}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *comment.DeleteCommentPayload) error {
			userID := middleware.GetUserID(c)
			return h.commentService.DeleteComment(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&comment.DeleteCommentPayload{},
	)(c)
}
