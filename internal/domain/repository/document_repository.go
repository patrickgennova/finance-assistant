package repository

import (
	"context"

	"finance-assistant/internal/domain/entity"
	"github.com/google/uuid"
)

type DocumentRepository interface {
	Create(ctx context.Context, document *entity.Document) error
	FindByID(ctx context.Context, id int64) (*entity.Document, error)
	FindByExternalID(ctx context.Context, externalID uuid.UUID) (*entity.Document, error)
	FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.Document, error)
	Update(ctx context.Context, document *entity.Document) error
	UpdateStatus(ctx context.Context, id int64, status entity.DocumentStatus) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*entity.Document, error)
	CountByUserID(ctx context.Context, userID int64) (int, error)
}
