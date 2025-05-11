package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"finance-assistant/internal/domain/entity"
	"finance-assistant/internal/domain/repository"
	"finance-assistant/internal/infrastructure/kafka"
	"github.com/google/uuid"
)

var (
	ErrDocumentNotFound        = errors.New("Documento não encontrado")
	ErrUserNotFoundForDocument = errors.New("Usuário não encontrado para este documento")
)

type DocumentService struct {
	repo          repository.DocumentRepository
	userRepo      repository.UserRepository
	kafkaProducer *kafka.Producer
}

func NewDocumentService(
	repo repository.DocumentRepository,
	userRepo repository.UserRepository,
	kafkaProducer *kafka.Producer,
) *DocumentService {
	return &DocumentService{
		repo:          repo,
		userRepo:      userRepo,
		kafkaProducer: kafkaProducer,
	}
}

// CreateDocument cria um novo documento e o envia para processamento
func (s *DocumentService) CreateDocument(
	ctx context.Context,
	userExternalID uuid.UUID,
	documentType,
	filename,
	contentType,
	fileContent string,
	categories []string,
) (*entity.Document, error) {
	// Buscar usuário pelo externalID
	user, err := s.userRepo.FindByExternalID(ctx, userExternalID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFoundForDocument
	}

	// Criar novo documento
	document, err := entity.NewDocument(
		user.ID,
		documentType,
		filename,
		contentType,
		fileContent,
		categories,
	)
	if err != nil {
		return nil, err
	}

	// Primeiro salva como pendente
	document.Status = entity.DocumentStatusPending

	// Salvar no banco de dados
	if err := s.repo.Create(ctx, document); err != nil {
		return nil, fmt.Errorf("erro ao salvar documento: %w", err)
	}

	// Log para debug
	log.Printf("Documento %s criado com sucesso. Enviando para processamento...", document.ExternalID)

	// Tentar enviar para o Kafka
	if err := s.kafkaProducer.SendDocument(document); err != nil {
		// Se falhar no envio, atualiza status para falha
		log.Printf("Erro ao enviar documento %s para Kafka: %v", document.ExternalID, err)
		_ = s.repo.UpdateStatus(ctx, document.ID, entity.DocumentStatusFailed)
		return nil, fmt.Errorf("erro ao enviar documento para processamento: %w", err)
	}

	// Atualizar status para processando
	log.Printf("Documento %s enviado para Kafka com sucesso. Atualizando status...", document.ExternalID)
	document.UpdateStatus(entity.DocumentStatusProcessing)
	if err := s.repo.UpdateStatus(ctx, document.ID, entity.DocumentStatusProcessing); err != nil {
		log.Printf("Aviso: Não foi possível atualizar o status do documento %s: %v", document.ExternalID, err)
		// Não retornamos erro aqui, pois o documento já foi enviado para o Kafka
	}

	return document, nil
}

// GetDocumentByExternalID obtém um documento pelo seu ID externo
func (s *DocumentService) GetDocumentByExternalID(ctx context.Context, externalID uuid.UUID) (*entity.Document, error) {
	document, err := s.repo.FindByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, ErrDocumentNotFound
	}
	return document, nil
}

// GetDocumentsByUserExternalID lista documentos de um usuário específico
func (s *DocumentService) GetDocumentsByUserExternalID(ctx context.Context, userExternalID uuid.UUID, page, perPage int) ([]*entity.Document, int, error) {
	// Buscar usuário pelo externalID
	user, err := s.userRepo.FindByExternalID(ctx, userExternalID)
	if err != nil {
		return nil, 0, err
	}
	if user == nil {
		return nil, 0, ErrUserNotFound
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	// Contar total de documentos do usuário
	total, err := s.repo.CountByUserID(ctx, user.ID)
	if err != nil {
		return nil, 0, err
	}

	// Buscar documentos do usuário
	documents, err := s.repo.FindByUserID(ctx, user.ID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	// Se não houver documentos, retornar slice vazio
	if documents == nil {
		return []*entity.Document{}, total, nil
	}

	return documents, total, nil
}

// UpdateDocumentStatus atualiza o status de um documento
func (s *DocumentService) UpdateDocumentStatus(ctx context.Context, externalID uuid.UUID, status entity.DocumentStatus) (*entity.Document, error) {
	document, err := s.repo.FindByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, ErrDocumentNotFound
	}

	document.UpdateStatus(status)
	if err := s.repo.Update(ctx, document); err != nil {
		return nil, err
	}

	return document, nil
}

// DeleteDocument exclui um documento
func (s *DocumentService) DeleteDocument(ctx context.Context, externalID uuid.UUID) error {
	document, err := s.repo.FindByExternalID(ctx, externalID)
	if err != nil {
		return err
	}
	if document == nil {
		return ErrDocumentNotFound
	}

	return s.repo.Delete(ctx, document.ID)
}

// ListDocuments lista todos os documentos com paginação
func (s *DocumentService) ListDocuments(ctx context.Context, page, perPage int) ([]*entity.Document, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	documents, err := s.repo.List(ctx, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	// Se não houver documentos, retornar slice vazio
	if documents == nil {
		documents = []*entity.Document{}
	}

	// Contagem total seria ideal, mas simplificamos aqui
	return documents, len(documents), nil
}
