package middleware

import (
	"github.com/gin-gonic/gin"
)

// ProcessArrayFields processa campos de formulário que devem ser tratados como arrays
func ProcessArrayFields() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Processar o corpo da requisição
		err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
		if err == nil {
			// Verificar campos múltiplos como 'categories'
			if values, exists := c.Request.PostForm["categories[]"]; exists {
				c.Request.Form["categories"] = values
			}
		}
		c.Next()
	}
}
