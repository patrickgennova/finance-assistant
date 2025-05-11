package validation

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// ValidationError contém informações sobre um erro de validação
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors é uma coleção de erros de validação
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Validate valida uma struct usando o validator do Gin
func Validate(v interface{}) (bool, ValidationErrors) {
	validate := binding.Validator.Engine().(*validator.Validate)

	// Traduzir mensagens de erro
	translate := func(fe validator.FieldError) string {
		field := fe.Field()
		tag := fe.Tag()
		param := fe.Param()

		switch tag {
		case "required":
			return fmt.Sprintf("O campo %s é obrigatório", field)
		case "email":
			return fmt.Sprintf("O campo %s deve ser um email válido", field)
		case "min":
			return fmt.Sprintf("O campo %s deve ter pelo menos %s caracteres", field, param)
		case "max":
			return fmt.Sprintf("O campo %s deve ter no máximo %s caracteres", field, param)
		default:
			return fmt.Sprintf("O campo %s é inválido", field)
		}
	}

	err := validate.Struct(v)
	if err == nil {
		return true, ValidationErrors{}
	}

	var validationErrs ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		validationErrs.Errors = append(validationErrs.Errors, ValidationError{
			Field:   strings.ToLower(err.Field()),
			Message: translate(err),
		})
	}

	return false, validationErrs
}
