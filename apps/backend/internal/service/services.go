package service

import (
	"github.com/gulistaneraza01/go-task-manager-echo/internal/lib/job"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/repository"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
)

type Services struct {
	Auth     *AuthService
	Job      *job.JobService
	Todo     *TodoService
	Comment  *CommentService
	Category *CategoryService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:      s.Job,
		Auth:     authService,
		Todo:     NewTodoService(s, repos.Todo, repos.Category),
		Comment:  NewCommentService(s, repos.Comment, repos.Todo),
		Category: NewCategoryService(s, repos.Category),
	}, nil
}
