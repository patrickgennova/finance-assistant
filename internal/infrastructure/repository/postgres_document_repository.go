package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"finance-assistant/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresDocumentRepository struct {
	db *sqlx.DB
}

func NewPostgresDocumentRepository(db *sqlx.DB) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{
		db: db,
	}
}

func (r *PostgresDocumentRepository) Create(ctx context.Context, document *entity.Document) error {
	query := `
		INSERT INTO documents (
			external_id, user_id, document_type, filename, content_type, 
			file_content, categories, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	// Converter categories para JSON
	categoriesJSON, err := json.Marshal(document.Categories)
	if err != nil {
		return fmt.Errorf("error marshaling categories: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		document.ExternalID,
		document.UserID,
		document.DocumentType,
		document.Filename,
		document.ContentType,
		document.FileContent,
		categoriesJSON,
		document.Status,
		document.CreatedAt,
		document.UpdatedAt,
	).Scan(&document.ID)

	if err != nil {
		return fmt.Errorf("error creating document: %w", err)
	}

	return nil
}

func (r *PostgresDocumentRepository) FindByID(ctx context.Context, id int64) (*entity.Document, error) {
	query := `
		SELECT
			id, external_id, user_id, document_type, filename, content_type, 
			file_content, categories, status, created_at, updated_at
		FROM documents
		WHERE id = $1
	`

	var documentDB struct {
		ID           int64                 `db:"id"`
		ExternalID   uuid.UUID             `db:"external_id"`
		UserID       int64                 `db:"user_id"`
		DocumentType string                `db:"document_type"`
		Filename     string                `db:"filename"`
		ContentType  string                `db:"content_type"`
		FileContent  string                `db:"file_content"`
		Categories   []byte                `db:"categories"`
		Status       entity.DocumentStatus `db:"status"`
		CreatedAt    sql.NullTime          `db:"created_at"`
		UpdatedAt    sql.NullTime          `db:"updated_at"`
	}

	err := r.db.GetContext(ctx, &documentDB, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding document by ID: %w", err)
	}

	// Converter categories de JSON para []string
	var categories []string
	if err := json.Unmarshal(documentDB.Categories, &categories); err != nil {
		return nil, fmt.Errorf("error unmarshaling categories: %w", err)
	}

	// Converter para entity.Document
	document := &entity.Document{
		ID:           documentDB.ID,
		ExternalID:   documentDB.ExternalID,
		UserID:       documentDB.UserID,
		DocumentType: documentDB.DocumentType,
		Filename:     documentDB.Filename,
		ContentType:  documentDB.ContentType,
		FileContent:  documentDB.FileContent,
		Categories:   categories,
		Status:       documentDB.Status,
		CreatedAt:    documentDB.CreatedAt.Time,
		UpdatedAt:    documentDB.UpdatedAt.Time,
	}

	return document, nil
}

func (r *PostgresDocumentRepository) FindByExternalID(ctx context.Context, externalID uuid.UUID) (*entity.Document, error) {
	query := `
		SELECT
			id, external_id, user_id, document_type, filename, content_type, 
			file_content, categories, status, created_at, updated_at
		FROM documents
		WHERE external_id = $1
	`

	var documentDB struct {
		ID           int64                 `db:"id"`
		ExternalID   uuid.UUID             `db:"external_id"`
		UserID       int64                 `db:"user_id"`
		DocumentType string                `db:"document_type"`
		Filename     string                `db:"filename"`
		ContentType  string                `db:"content_type"`
		FileContent  string                `db:"file_content"`
		Categories   []byte                `db:"categories"`
		Status       entity.DocumentStatus `db:"status"`
		CreatedAt    sql.NullTime          `db:"created_at"`
		UpdatedAt    sql.NullTime          `db:"updated_at"`
	}

	err := r.db.GetContext(ctx, &documentDB, query, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding document by external ID: %w", err)
	}

	// Converter categories de JSON para []string
	var categories []string
	if err := json.Unmarshal(documentDB.Categories, &categories); err != nil {
		return nil, fmt.Errorf("error unmarshaling categories: %w", err)
	}

	// Converter para entity.Document
	document := &entity.Document{
		ID:           documentDB.ID,
		ExternalID:   documentDB.ExternalID,
		UserID:       documentDB.UserID,
		DocumentType: documentDB.DocumentType,
		Filename:     documentDB.Filename,
		ContentType:  documentDB.ContentType,
		FileContent:  documentDB.FileContent,
		Categories:   categories,
		Status:       documentDB.Status,
		CreatedAt:    documentDB.CreatedAt.Time,
		UpdatedAt:    documentDB.UpdatedAt.Time,
	}

	return document, nil
}

func (r *PostgresDocumentRepository) FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.Document, error) {
	query := `
		SELECT
			id, external_id, user_id, document_type, filename, content_type, 
			file_content, categories, status, created_at, updated_at
		FROM documents
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error finding documents by user ID: %w", err)
	}
	defer rows.Close()

	var documents []*entity.Document
	for rows.Next() {
		var documentDB struct {
			ID           int64                 `db:"id"`
			ExternalID   uuid.UUID             `db:"external_id"`
			UserID       int64                 `db:"user_id"`
			DocumentType string                `db:"document_type"`
			Filename     string                `db:"filename"`
			ContentType  string                `db:"content_type"`
			FileContent  string                `db:"file_content"`
			Categories   []byte                `db:"categories"`
			Status       entity.DocumentStatus `db:"status"`
			CreatedAt    sql.NullTime          `db:"created_at"`
			UpdatedAt    sql.NullTime          `db:"updated_at"`
		}

		err := rows.Scan(
			&documentDB.ID,
			&documentDB.ExternalID,
			&documentDB.UserID,
			&documentDB.DocumentType,
			&documentDB.Filename,
			&documentDB.ContentType,
			&documentDB.FileContent,
			&documentDB.Categories,
			&documentDB.Status,
			&documentDB.CreatedAt,
			&documentDB.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning document row: %w", err)
		}

		// Converter categories de JSON para []string
		var categories []string
		if err := json.Unmarshal(documentDB.Categories, &categories); err != nil {
			return nil, fmt.Errorf("error unmarshaling categories: %w", err)
		}

		document := &entity.Document{
			ID:           documentDB.ID,
			ExternalID:   documentDB.ExternalID,
			UserID:       documentDB.UserID,
			DocumentType: documentDB.DocumentType,
			Filename:     documentDB.Filename,
			ContentType:  documentDB.ContentType,
			FileContent:  documentDB.FileContent,
			Categories:   categories,
			Status:       documentDB.Status,
			CreatedAt:    documentDB.CreatedAt.Time,
			UpdatedAt:    documentDB.UpdatedAt.Time,
		}

		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating document rows: %w", err)
	}

	return documents, nil
}

func (r *PostgresDocumentRepository) Update(ctx context.Context, document *entity.Document) error {
	query := `
		UPDATE documents
		SET document_type = $1, filename = $2, content_type = $3, 
			file_content = $4, categories = $5, status = $6, updated_at = $7
		WHERE id = $8
	`

	// Converter categories para JSON
	categoriesJSON, err := json.Marshal(document.Categories)
	if err != nil {
		return fmt.Errorf("error marshaling categories: %w", err)
	}

	result, err := r.db.ExecContext(
		ctx,
		query,
		document.DocumentType,
		document.Filename,
		document.ContentType,
		document.FileContent,
		categoriesJSON,
		document.Status,
		document.UpdatedAt,
		document.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no document found with ID: %d", document.ID)
	}

	return nil
}

func (r *PostgresDocumentRepository) UpdateStatus(ctx context.Context, id int64, status entity.DocumentStatus) error {
	query := `
		UPDATE documents
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("error updating document status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no document found with ID: %d", id)
	}

	return nil
}

func (r *PostgresDocumentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM documents WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no document found with ID: %d", id)
	}

	return nil
}

func (r *PostgresDocumentRepository) List(ctx context.Context, limit, offset int) ([]*entity.Document, error) {
	query := `
		SELECT
			id, external_id, user_id, document_type, filename, content_type, 
			file_content, categories, status, created_at, updated_at
		FROM documents
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing documents: %w", err)
	}
	defer rows.Close()

	var documents []*entity.Document
	for rows.Next() {
		var documentDB struct {
			ID           int64                 `db:"id"`
			ExternalID   uuid.UUID             `db:"external_id"`
			UserID       int64                 `db:"user_id"`
			DocumentType string                `db:"document_type"`
			Filename     string                `db:"filename"`
			ContentType  string                `db:"content_type"`
			FileContent  string                `db:"file_content"`
			Categories   []byte                `db:"categories"`
			Status       entity.DocumentStatus `db:"status"`
			CreatedAt    sql.NullTime          `db:"created_at"`
			UpdatedAt    sql.NullTime          `db:"updated_at"`
		}

		err := rows.Scan(
			&documentDB.ID,
			&documentDB.ExternalID,
			&documentDB.UserID,
			&documentDB.DocumentType,
			&documentDB.Filename,
			&documentDB.ContentType,
			&documentDB.FileContent,
			&documentDB.Categories,
			&documentDB.Status,
			&documentDB.CreatedAt,
			&documentDB.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning document row: %w", err)
		}

		// Converter categories de JSON para []string
		var categories []string
		if err := json.Unmarshal(documentDB.Categories, &categories); err != nil {
			return nil, fmt.Errorf("error unmarshaling categories: %w", err)
		}

		document := &entity.Document{
			ID:           documentDB.ID,
			ExternalID:   documentDB.ExternalID,
			UserID:       documentDB.UserID,
			DocumentType: documentDB.DocumentType,
			Filename:     documentDB.Filename,
			ContentType:  documentDB.ContentType,
			FileContent:  documentDB.FileContent,
			Categories:   categories,
			Status:       documentDB.Status,
			CreatedAt:    documentDB.CreatedAt.Time,
			UpdatedAt:    documentDB.UpdatedAt.Time,
		}

		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating document rows: %w", err)
	}

	return documents, nil
}

func (r *PostgresDocumentRepository) CountByUserID(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM documents WHERE user_id = $1`

	var count int
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("error counting documents by user ID: %w", err)
	}

	return count, nil
}
