package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDocumentType     = errors.New("Tipo de documento inválido")
	ErrInvalidDocumentContent  = errors.New("Conteúdo do documento inválido")
	ErrInvalidDocumentUserID   = errors.New("ID de usuário inválido")
	ErrInvalidDocumentFilename = errors.New("Nome de arquivo inválido")
)

type DocumentStatus string

const (
	DocumentStatusPending    DocumentStatus = "pending"
	DocumentStatusProcessing DocumentStatus = "processing"
	DocumentStatusProcessed  DocumentStatus = "processed"
	DocumentStatusFailed     DocumentStatus = "failed"
)

type Document struct {
	ID           int64          `db:"id" json:"id"`
	ExternalID   uuid.UUID      `db:"external_id" json:"external_id"`
	UserID       int64          `db:"user_id" json:"user_id"`
	DocumentType string         `db:"document_type" json:"document_type"`
	Filename     string         `db:"filename" json:"filename"`
	ContentType  string         `db:"content_type" json:"content_type"`
	FileContent  string         `db:"file_content" json:"file_content"`
	Categories   []string       `db:"categories" json:"categories"`
	Status       DocumentStatus `db:"status" json:"status"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

// NewDocument cria um novo documento
func NewDocument(userID int64, documentType, filename, contentType, fileContent string, categories []string) (*Document, error) {
	if userID <= 0 {
		return nil, ErrInvalidDocumentUserID
	}
	if documentType == "" {
		return nil, ErrInvalidDocumentType
	}
	if filename == "" {
		return nil, ErrInvalidDocumentFilename
	}
	if fileContent == "" {
		return nil, ErrInvalidDocumentContent
	}

	if categories == nil {
		categories = []string{}
	}

	now := time.Now()
	return &Document{
		ExternalID:   uuid.New(),
		UserID:       userID,
		DocumentType: documentType,
		Filename:     filename,
		ContentType:  contentType,
		FileContent:  fileContent,
		Categories:   categories,
		Status:       DocumentStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Validate valida os dados do documento
func (d *Document) Validate() error {
	if d.UserID <= 0 {
		return ErrInvalidDocumentUserID
	}
	if d.DocumentType == "" {
		return ErrInvalidDocumentType
	}
	if d.Filename == "" {
		return ErrInvalidDocumentFilename
	}
	if d.FileContent == "" {
		return ErrInvalidDocumentContent
	}
	return nil
}

// UpdateStatus atualiza o status do documento
func (d *Document) UpdateStatus(status DocumentStatus) {
	d.Status = status
	d.UpdatedAt = time.Now()
}

// UpdateCategories atualiza as categorias do documento
func (d *Document) UpdateCategories(categories []string) {
	d.Categories = categories
	d.UpdatedAt = time.Now()
}
