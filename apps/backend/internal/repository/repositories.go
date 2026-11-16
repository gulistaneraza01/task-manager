package repository

import "github.com/gulistaneraza01/go-task-manager-echo/internal/server"

type Repositories struct {
	Todo    *TodoRepository
	Comment *CommentRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Todo:    NewTodoRepository(s),
		Comment: NewCommentRepository(s),
	}
}
