package repository

import "github.com/gulistaneraza01/go-task-manager-echo/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}
