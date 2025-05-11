package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"finance-assistant/internal/domain/entity"
	"finance-assistant/internal/domain/service"
	"finance-assistant/internal/interface/api/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	documentService *service.DocumentService
}

func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
	}
}

// Create godoc
// @Summary      Criar documento
// @Description  Envia um novo documento para um usuário
// @Tags         documents
// @Accept       multipart/form-data
// @Produce      json
// @Param        id              path      string   true  "ID do usuário"
// @Param        document_type   formData  string   true  "Tipo de documento (ex: bank_statement, invoice, receipt)"
// @Param        categories      formData  []string false "Categorias do documento (opcional)"
// @Param        file            formData  file     true  "Arquivo do documento (PDF, DOCX, XLS, PNG, JPEG)"
// @Success      201             {object}  dto.DocumentResponse
// @Failure      400             {object}  map[string]interface{}
// @Failure      404             {object}  map[string]interface{}
// @Failure      500             {object}  map[string]interface{}
// @Router       /users/{id}/documents [post]
func (h *DocumentHandler) Create(c *gin.Context) {
	// Obter ID do usuário a partir do parâmetro da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	// Carregar dados do form-data
	var req dto.DocumentUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Erro ao processar formulário",
			"details": err.Error(),
		})
		return
	}

	// Verificar se o tipo de documento foi fornecido
	if req.DocumentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de documento é obrigatório"})
		return
	}

	// Obter arquivo enviado
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não encontrado ou inválido"})
		return
	}
	defer file.Close()

	// Validar tamanho do arquivo (exemplo: limite de 10MB)
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo muito grande, tamanho máximo permitido é 10MB"})
		return
	}

	// Obter nome do arquivo e extensão
	filename := fileHeader.Filename
	fileExt := strings.ToLower(filepath.Ext(filename))

	// Validar tipo de arquivo
	allowedExtensions := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}

	if !allowedExtensions[fileExt] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Tipo de arquivo não suportado",
			"details": "Tipos permitidos: PDF, DOC, DOCX, XLS, XLSX, PNG, JPG, JPEG",
		})
		return
	}

	// Determinar o content type com base na extensão do arquivo
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	case ".doc", ".docx":
		contentType = "application/msword"
	case ".xls", ".xlsx":
		contentType = "application/vnd.ms-excel"
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	default:
		contentType = "application/octet-stream"
	}

	// Ler o conteúdo do arquivo
	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler arquivo"})
		return
	}

	// Verificar se o arquivo está vazio
	if len(fileContent) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo vazio"})
		return
	}

	// Converter o conteúdo para base64
	base64Content := base64.StdEncoding.EncodeToString(fileContent)

	// Log para debug
	log.Printf("Processando upload de arquivo: %s, tamanho: %d bytes, tipo: %s",
		filename, len(fileContent), contentType)

	// Criar documento
	document, err := h.documentService.CreateDocument(
		c.Request.Context(),
		userID,
		req.DocumentType,
		filename,
		contentType,
		base64Content,
		req.Categories,
	)
	if err != nil {
		var status int
		var message string

		switch err {
		case service.ErrUserNotFoundForDocument:
			status = http.StatusNotFound
			message = "Usuário não encontrado"
		case entity.ErrInvalidDocumentType:
			status = http.StatusBadRequest
			message = "Tipo de documento inválido"
		case entity.ErrInvalidDocumentContent:
			status = http.StatusBadRequest
			message = "Conteúdo do documento inválido"
		case entity.ErrInvalidDocumentFilename:
			status = http.StatusBadRequest
			message = "Nome de arquivo inválido"
		default:
			status = http.StatusInternalServerError
			message = fmt.Sprintf("Erro ao criar documento: %v", err)
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusCreated, dto.DocumentFromEntity(document))
}

// GetByID godoc
// @Summary      Buscar documento por ID
// @Description  Retorna um documento pelo seu ID
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id        path    string  true   "ID do documento"
// @Param        detailed  query   bool    false  "Se verdadeiro, inclui o conteúdo do arquivo na resposta (default: false)"
// @Success      200       {object}  dto.DocumentResponse
// @Success      200       {object}  dto.DocumentDetailResponse "Quando detailed=true"
// @Failure      400       {object}  map[string]interface{}
// @Failure      404       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]interface{}
// @Router       /documents/{id} [get]
func (h *DocumentHandler) GetByID(c *gin.Context) {
	// Obter ID do documento a partir do parâmetro da URL
	documentIDStr := c.Param("id")
	documentID, err := uuid.Parse(documentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
		return
	}

	// Verificar se o cliente quer os detalhes completos, incluindo o conteúdo
	detailed := c.DefaultQuery("detailed", "false") == "true"

	// Buscar documento
	document, err := h.documentService.GetDocumentByExternalID(c.Request.Context(), documentID)
	if err != nil {
		if err == service.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar a resposta adequada com base no parâmetro 'detailed'
	if detailed {
		c.JSON(http.StatusOK, dto.DocumentDetailFromEntity(document))
	} else {
		c.JSON(http.StatusOK, dto.DocumentFromEntity(document))
	}
}

// GetByUserID godoc
// @Summary      Listar documentos de um usuário
// @Description  Retorna uma lista paginada de documentos de um usuário específico
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id     path      string  true   "ID do usuário"
// @Param        page   query     int     false  "Página atual (padrão: 1)"
// @Param        limit  query     int     false  "Limite de itens por página (padrão: 10)"
// @Success      200    {object}  dto.DocumentListResponse
// @Failure      400    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /users/{id}/documents [get]
func (h *DocumentHandler) GetByUserID(c *gin.Context) {
	// Obter ID do usuário a partir do parâmetro da URL
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	// Parâmetros de paginação
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Buscar documentos do usuário
	documents, total, err := h.documentService.GetDocumentsByUserExternalID(c.Request.Context(), userID, page, limit)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converter entidades para DTOs
	response := make([]dto.DocumentResponse, len(documents))
	for i, doc := range documents {
		response[i] = dto.DocumentFromEntity(doc)
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": response,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// UpdateStatus godoc
// @Summary      Atualizar status do documento
// @Description  Atualiza o status de processamento de um documento
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id      path      string                        true  "ID do documento"
// @Param        status  body      dto.DocumentStatusUpdateRequest  true  "Novo status"
// @Success      200     {object}  dto.DocumentResponse
// @Failure      400     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /documents/{id}/status [put]
func (h *DocumentHandler) UpdateStatus(c *gin.Context) {
	// Obter ID do documento a partir do parâmetro da URL
	documentIDStr := c.Param("id")
	documentID, err := uuid.Parse(documentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
		return
	}

	var req dto.DocumentStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dados inválido"})
		return
	}

	// Validar status
	var status entity.DocumentStatus
	switch req.Status {
	case "pending":
		status = entity.DocumentStatusPending
	case "processing":
		status = entity.DocumentStatusProcessing
	case "processed":
		status = entity.DocumentStatusProcessed
	case "failed":
		status = entity.DocumentStatusFailed
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido"})
		return
	}

	// Atualizar status
	document, err := h.documentService.UpdateDocumentStatus(c.Request.Context(), documentID, status)
	if err != nil {
		if err == service.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.DocumentFromEntity(document))
}

// Delete godoc
// @Summary      Excluir documento
// @Description  Remove um documento do sistema
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do documento"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /documents/{id} [delete]
func (h *DocumentHandler) Delete(c *gin.Context) {
	// Obter ID do documento a partir do parâmetro da URL
	documentIDStr := c.Param("id")
	documentID, err := uuid.Parse(documentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
		return
	}

	// Excluir documento
	if err := h.documentService.DeleteDocument(c.Request.Context(), documentID); err != nil {
		if err == service.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary      Listar documentos
// @Description  Retorna uma lista paginada de todos os documentos
// @Tags         documents
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Página atual (padrão: 1)"
// @Param        limit  query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200    {object}  dto.DocumentListResponse
// @Failure      500    {object}  map[string]interface{}
// @Router       /documents [get]
func (h *DocumentHandler) List(c *gin.Context) {
	// Parâmetros de paginação
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Listar documentos
	documents, total, err := h.documentService.ListDocuments(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converter entidades para DTOs
	response := make([]dto.DocumentResponse, len(documents))
	for i, doc := range documents {
		response[i] = dto.DocumentFromEntity(doc)
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": response,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// DownloadDocument godoc
// @Summary      Download do documento
// @Description  Retorna o arquivo do documento para download
// @Tags         documents
// @Accept       json
// @Produce      octet-stream
// @Param        id   path      string  true  "ID do documento"
// @Success      200  {file}    binary
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /documents/{id}/download [get]
func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	// Obter ID do documento
	documentIDStr := c.Param("id")
	documentID, err := uuid.Parse(documentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de documento inválido"})
		return
	}

	// Buscar documento
	document, err := h.documentService.GetDocumentByExternalID(c.Request.Context(), documentID)
	if err != nil {
		if err == service.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Documento não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converter o conteúdo base64 para binário
	fileContent, err := base64.StdEncoding.DecodeString(document.FileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar conteúdo do arquivo"})
		return
	}

	// Configurar cabeçalhos para download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", document.Filename))
	c.Header("Content-Type", document.ContentType)
	c.Header("Content-Length", strconv.Itoa(len(fileContent)))

	// Enviar o arquivo
	c.Data(http.StatusOK, document.ContentType, fileContent)
}
