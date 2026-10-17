package handler

import (
	"github.com/gulistaneraza01/go-task-manager-echo/internal/server"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
	}
}
