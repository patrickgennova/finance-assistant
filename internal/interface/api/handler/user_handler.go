package handler

import (
	"finance-assistant/internal/domain/service"
	"finance-assistant/internal/interface/api/dto"
	"finance-assistant/internal/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Create godoc
// @Summary      Criar usuário
// @Description  Cria um novo usuário no sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.UserRequest  true  "Dados do usuário"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Verificar se é um erro de validação
		if _, ok := err.(validator.ValidationErrors); ok {
			// Reutilizar a função Validate para traduzir as mensagens
			valid, errs := validation.Validate(req)
			if !valid {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Erro de validação",
					"details": errs.Errors,
				})
				return
			}
		}

		// Se não for erro de validação
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de dados inválido"})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Name, req.Email, req.Phone)
	if err != nil {
		if err == service.ErrEmailAlreadyUsed {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.FromEntity(user))
}

// GetByID godoc
// @Summary      Buscar usuário por ID
// @Description  Retorna um usuário pelo seu ID externo
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do usuário"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	externalID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	user, err := h.userService.GetUserByExternalID(c.Request.Context(), externalID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.FromEntity(user))
}

// Update godoc
// @Summary      Atualizar usuário
// @Description  Atualiza os dados de um usuário existente
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "ID do usuário"
// @Param        user  body      dto.UpdateUserRequest  true  "Dados para atualização"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	externalID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), externalID, req.Name, req.Email, req.Phone)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		if err == service.ErrEmailAlreadyUsed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email já está em uso"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.FromEntity(user))
}

// Delete godoc
// @Summary      Excluir usuário
// @Description  Remove um usuário do sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do usuário"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	externalID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), externalID); err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary      Listar usuários
// @Description  Retorna uma lista paginada de usuários
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Página atual (padrão: 1)"
// @Param        limit  query     int  false  "Limite de itens por página (padrão: 10)"
// @Success      200    {object}  dto.UserListResponse
// @Failure      500    {object}  map[string]interface{}
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
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

	users, err := h.userService.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Converter entidades para DTOs
	response := make([]dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = dto.FromEntity(user)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": response,
		"page":  page,
		"limit": limit,
	})
}
