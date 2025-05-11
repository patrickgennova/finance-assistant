package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidUserName  = errors.New("nome de usuário inválido")
	ErrInvalidUserEmail = errors.New("email de usuário inválido")
)

type User struct {
	ID         int64     `db:"id" json:"id"`
	ExternalID uuid.UUID `db:"external_id" json:"external_id"`
	Name       string    `db:"name" json:"name"`
	Email      string    `db:"email" json:"email"`
	Phone      string    `db:"phone" json:"phone"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// NewUser cria um novo usuário com validações
func NewUser(name, email, phone string) (*User, error) {
	if name == "" {
		return nil, ErrInvalidUserName
	}
	if email == "" {
		return nil, ErrInvalidUserEmail
	}

	return &User{
		ExternalID: uuid.New(),
		Name:       name,
		Email:      email,
		Phone:      phone,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

// Validate valida os dados do usuário
func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidUserName
	}
	if u.Email == "" {
		return ErrInvalidUserEmail
	}
	return nil
}

// Update atualiza os dados do usuário
func (u *User) Update(name, email, phone string) error {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	u.Phone = phone
	u.UpdatedAt = time.Now()
	return u.Validate()
}
