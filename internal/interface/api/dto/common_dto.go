package dto

// ErrorResponse modelo padrão para respostas de erro
// @Description Resposta padrão para erros na API
type ErrorResponse struct {
	Error   string        `json:"error" example:"Erro de validação"` // Mensagem de erro
	Details []ErrorDetail `json:"details,omitempty"`                 // Detalhes do erro (opcional)
}

// ErrorDetail detalhes adicionais para erros
// @Description Detalhes adicionais sobre um erro
type ErrorDetail struct {
	Field   string `json:"field" example:"email"`                         // Campo com erro
	Message string `json:"message" example:"O campo email é obrigatório"` // Mensagem descritiva do erro
}

// HealthResponse resposta do endpoint de health check
// @Description Resposta do endpoint de verificação de saúde da API
type HealthResponse struct {
	Status string `json:"status" example:"ok"` // Status da API
}
