package dto

import (
	"time"

	"finance-assistant/internal/domain/entity"
	"github.com/google/uuid"
)

// UserRequest representa os dados enviados para criar/atualizar um usuário
type UserRequest struct {
	Name  string `json:"name" binding:"required" example:"João Silva"`                    // Nome completo do usuário
	Email string `json:"email" binding:"required,email" example:"joao.silva@example.com"` // Email do usuário
	Phone string `json:"phone,omitempty" example:"(11) 98765-4321"`                       // Telefone do usuário (opcional)
}

// UserResponse representa os dados retornados pela API
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"` // ID externo do usuário
	Name      string    `json:"name" example:"João Silva"`                         // Nome do usuário
	Email     string    `json:"email" example:"joao.silva@example.com"`            // Email do usuário
	Phone     string    `json:"phone,omitempty" example:"(11) 98765-4321"`         // Telefone do usuário
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`         // Data de criação
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`         // Data de atualização
}

// UpdateUserRequest representa os dados enviados para atualizar um usuário
type UpdateUserRequest struct {
	Name  string `json:"name" example:"João Silva Atualizado"`                            // Nome atualizado do usuário (opcional)
	Email string `json:"email" binding:"omitempty,email" example:"joao.novo@example.com"` // Email atualizado (opcional)
	Phone string `json:"phone" example:"(11) 99999-8888"`                                 // Telefone atualizado (opcional)
}

// UserListResponse representa a resposta de uma listagem paginada
type UserListResponse struct {
	Users []UserResponse `json:"users"`              // Lista de usuários
	Total int            `json:"total" example:"42"` // Número total de usuários
	Page  int            `json:"page" example:"1"`   // Página atual
	Limit int            `json:"limit" example:"10"` // Limite de itens por página
}

// FromEntity converte uma entidade User para UserResponse
func FromEntity(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ExternalID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToEntity converte UserRequest para uma entidade User
func (r *UserRequest) ToEntity() (*entity.User, error) {
	return entity.NewUser(r.Name, r.Email, r.Phone)
}
