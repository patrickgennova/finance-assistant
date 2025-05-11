package service

import (
	"context"
	"errors"

	"finance-assistant/internal/domain/entity"
	"finance-assistant/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound     = errors.New("usuário não encontrado")
	ErrEmailAlreadyUsed = errors.New("email já está em uso")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, phone string) (*entity.User, error) {
	// Verifica se já existe um usuário com este email
	existingUser, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}

	user, err := entity.NewUser(name, email, phone)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByExternalID(ctx context.Context, externalID uuid.UUID) (*entity.User, error) {
	user, err := s.repo.FindByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, externalID uuid.UUID, name, email, phone string) (*entity.User, error) {
	user, err := s.repo.FindByExternalID(ctx, externalID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Se o email for alterado, verificar se já está em uso
	if email != "" && email != user.Email {
		existingUser, err := s.repo.FindByEmail(ctx, email)
		if err == nil && existingUser != nil {
			return nil, ErrEmailAlreadyUsed
		}
	}

	if err := user.Update(name, email, phone); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, externalID uuid.UUID) error {
	user, err := s.repo.FindByExternalID(ctx, externalID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(ctx, user.ID)
}

func (s *UserService) ListUsers(ctx context.Context, page, perPage int) ([]*entity.User, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	return s.repo.List(ctx, perPage, offset)
}
