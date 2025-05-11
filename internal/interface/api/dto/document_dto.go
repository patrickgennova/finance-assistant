package dto

import (
	"finance-assistant/internal/domain/entity"
	"time"

	"github.com/google/uuid"
)

// DocumentUploadRequest representa os dados enviados para criar um documento
// @Description Dados do formulário para upload de um novo documento
type DocumentUploadRequest struct {
	DocumentType string   `form:"document_type" binding:"required" example:"bank_statement"` // Tipo de documento
	Categories   []string `form:"categories" example:"banco,mensal"`                         // Categorias para classificação (opcional)
	// O arquivo é enviado via multipart/form-data com o campo "file"
}

// DocumentResponse representa os dados retornados pela API
// @Description Informações de um documento armazenado
type DocumentResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"` // ID externo do documento
	DocumentType string    `json:"document_type" example:"bank_statement"`            // Tipo de documento
	Filename     string    `json:"filename" example:"extrato_janeiro.pdf"`            // Nome do arquivo
	ContentType  string    `json:"content_type" example:"application/pdf"`            // Tipo MIME do arquivo
	FileSize     int       `json:"file_size" example:"125000"`                        // Tamanho aproximado do arquivo em bytes
	Categories   []string  `json:"categories" example:"[\"banco\",\"mensal\"]"`       // Categorias do documento
	Status       string    `json:"status" example:"processing"`                       // Status de processamento (pending, processing, processed, failed)
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`         // Data de criação
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`         // Data de última atualização
}

// DocumentDetailResponse representa os dados detalhados do documento, incluindo o conteúdo
// @Description Informações detalhadas de um documento, incluindo seu conteúdo
type DocumentDetailResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`          // ID externo do documento
	DocumentType string    `json:"document_type" example:"bank_statement"`                     // Tipo de documento
	Filename     string    `json:"filename" example:"extrato_janeiro.pdf"`                     // Nome do arquivo
	ContentType  string    `json:"content_type" example:"application/pdf"`                     // Tipo MIME do arquivo
	FileSize     int       `json:"file_size" example:"125000"`                                 // Tamanho aproximado do arquivo em bytes
	FileContent  string    `json:"file_content" example:"JVBERi0xLjUKJYCBgoMKMSAwIG9iago8..."` // Conteúdo do arquivo em Base64
	Categories   []string  `json:"categories" example:"[\"banco\",\"mensal\"]"`                // Categorias do documento
	Status       string    `json:"status" example:"processing"`                                // Status de processamento (pending, processing, processed, failed)
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`                  // Data de criação
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`                  // Data de última atualização
}

// DocumentListResponse representa a resposta de uma listagem paginada de documentos
// @Description Lista paginada de documentos
type DocumentListResponse struct {
	Documents []DocumentResponse `json:"documents"`          // Lista de documentos
	Total     int                `json:"total" example:"42"` // Número total de documentos
	Page      int                `json:"page" example:"1"`   // Página atual
	Limit     int                `json:"limit" example:"10"` // Limite de itens por página
}

// DocumentStatusUpdateRequest representa a requisição para atualizar o status de um documento
// @Description Requisição para mudar o status de um documento
type DocumentStatusUpdateRequest struct {
	Status string `json:"status" binding:"required" example:"processed" enums:"pending,processing,processed,failed"` // Novo status do documento
}

// DocumentFromEntity converte uma entidade Document para DocumentResponse
func DocumentFromEntity(document *entity.Document) DocumentResponse {
	// Calcular tamanho do arquivo a partir do conteúdo base64
	fileSize := 0
	if document.FileContent != "" {
		// Estimar tamanho do arquivo a partir do base64
		// A string base64 é aproximadamente 33% maior que o binário original
		fileSize = len(document.FileContent) * 3 / 4
	}

	return DocumentResponse{
		ID:           document.ExternalID,
		DocumentType: document.DocumentType,
		Filename:     document.Filename,
		ContentType:  document.ContentType,
		FileSize:     fileSize,
		Categories:   document.Categories,
		Status:       string(document.Status),
		CreatedAt:    document.CreatedAt,
		UpdatedAt:    document.UpdatedAt,
	}
}

// DocumentDetailFromEntity converte uma entidade Document para DocumentDetailResponse
func DocumentDetailFromEntity(document *entity.Document) DocumentDetailResponse {
	// Calcular tamanho do arquivo
	fileSize := 0
	if document.FileContent != "" {
		fileSize = len(document.FileContent) * 3 / 4
	}

	return DocumentDetailResponse{
		ID:           document.ExternalID,
		DocumentType: document.DocumentType,
		Filename:     document.Filename,
		ContentType:  document.ContentType,
		FileSize:     fileSize,
		FileContent:  document.FileContent,
		Categories:   document.Categories,
		Status:       string(document.Status),
		CreatedAt:    document.CreatedAt,
		UpdatedAt:    document.UpdatedAt,
	}
}
