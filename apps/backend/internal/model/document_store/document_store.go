package document_store

import (
	"github.com/google/uuid"
	"github.com/gulistaneraza01/go-task-manager-echo/internal/model"
)

type DocumentStore struct {
	model.Base
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Bucket   string    `json:"bucket" db:"bucket"`
	Key      string    `json:"key" db:"key"`
	URL      string    `json:"url" db:"url"`
	FileName string    `json:"file_name" db:"file_name"`
	FileSize int64     `json:"file_size" db:"file_size"`
	FileType string    `json:"file_type" db:"file_type"`
	Metadata MetaData  `json:"metadata" db:"metadata"`
}

type MetaData struct {
	Error *string `json:"error" db:"error"`
	// Add other metadata fields as needed
}
