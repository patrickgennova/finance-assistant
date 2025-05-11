package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"finance-assistant/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (external_id, name, email, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.ExternalID,
		user.Name,
		user.Email,
		user.Phone,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT id, external_id, name, email, phone, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by ID: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByExternalID(ctx context.Context, externalID uuid.UUID) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT id, external_id, name, email, phone, created_at, updated_at
		FROM users
		WHERE external_id = $1
	`

	err := r.db.GetContext(ctx, &user, query, externalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by external ID: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	query := `
		SELECT id, external_id, name, email, phone, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, phone = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Phone,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %d", user.ID)
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %d", id)
	}

	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var users []*entity.User

	query := `
		SELECT id, external_id, name, email, phone, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	if err := r.db.SelectContext(ctx, &users, query, limit, offset); err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}
