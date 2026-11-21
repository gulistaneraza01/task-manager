package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/errs"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model/comment"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/jackc/pgx/v5"
)

type CommentRepository struct {
	server *server.Server
}

func NewCommentRepository(server *server.Server) *CommentRepository {
	return &CommentRepository{server: server}
}

// AddComment to Task
func (r *CommentRepository) AddComment(ctx context.Context, userID string, todoID uuid.UUID, payload *comment.CreateCommentPayload) (*comment.Comment, error) {

	stmt := `INSERT INTO comments
		(todo_id, user_id, content, parent_comment_id)
		VALUES (@todo_id, @user_id, @content, @parent_comment_id)
		RETURNING *`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"todo_id":           todoID,
		"user_id":           userID,
		"content":           payload.Content,
		"parent_comment_id": payload.ParentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	commentItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect comment: %w", err)
	}

	return &commentItem, nil
}

// GetCommentsByTodoID for user todo
func (r *CommentRepository) GetCommentsByTodoID(ctx context.Context, userID string, todoID uuid.UUID) ([]comment.Comment, error) {
	stmt := `SELECT * FROM comments
	WHERE todo_id = @todo_id AND user_id = @user_id
	ORDER BY created_at ASC`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"todo_id": todoID,
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	comments, err := pgx.CollectRows(rows, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect comments: %w", err)
	}

	return comments, nil
}

// GetCommentByID gets a single comment detail
func (r *CommentRepository) GetCommentByID(ctx context.Context, userID string, commentID uuid.UUID) (*comment.Comment, error) {
	stmt := `SELECT * FROM comments
	WHERE id = @comment_id AND user_id = @user_id`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"comment_id": commentID,
		"user_id":    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}

	commentItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect comment: %w", err)
	}

	return &commentItem, nil
}

// UpdateComment updates a comment for a todo
func (r *CommentRepository) UpdateComment(ctx context.Context, userID string, commentID uuid.UUID, payload *comment.UpdateCommentPayload) (*comment.Comment, error) {
	stmt := `UPDATE comments
	SET content = @content
	WHERE id = @comment_id AND user_id = @user_id
	RETURNING *`

	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"content":    payload.Content,
		"comment_id": commentID,
		"user_id":    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	updatedComment, err := pgx.CollectOneRow(row, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to collect updated comment: %w", err)
	}

	return &updatedComment, nil
}

// DeleteComment deletes the single comment
func (r *CommentRepository) DeleteComment(ctx context.Context, userID string, commentID uuid.UUID) error {

	stmt := `
		DELETE FROM comments
		WHERE id = @comment_id AND user_id = @user_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"comment_id": commentID,
		"user_id":    userID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		code := "COMMENT_NOT_FOUND"
		return errs.NewNotFoundError("comment not found", false, &code)
	}

	return nil
}

// CheckcommentExists
func (r *CommentRepository) CheckCommentExist(ctx context.Context, userID string, commentID uuid.UUID) (*comment.Comment, error) {

	stmt := `
		Select * FROM comments
		WHERE id = @comment_id AND user_id = @user_id
		RETURNING *
	`
	row, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{"comment_id": commentID, "user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to check comment existence: %w", err)
	}

	commentData, err := pgx.CollectOneRow(row, pgx.RowToStructByName[comment.Comment])
	if err != nil {
		return nil, fmt.Errorf("failed to get comment data: %w", err)
	}

	return &commentData, nil
}

// IsCommentExists checks if a comment exists for a user
func (r *CommentRepository) IsCommentExists(ctx context.Context, userID string, commentID uuid.UUID) (bool, error) {
	stmt := `SELECT EXISTS (
		SELECT 1 FROM comments
		WHERE id = @comment_id AND user_id = @user_id
	)`

	var exists bool
	err := r.server.DB.Pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"comment_id": commentID,
		"user_id":    userID,
	}).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check comment existence: %w", err)
	}

	return exists, nil
}
