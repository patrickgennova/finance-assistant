package inhttp

import (
	"finance-assistant/internal/interface/api/dto"
	"finance-assistant/internal/interface/api/handler"
	"finance-assistant/internal/interface/api/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "finance-assistant/docs"
)

func SetupRouter(userHandler *handler.UserHandler, documentHandler *handler.DocumentHandler, systemHandler *handler.SystemHandler) *gin.Engine {
	router := gin.Default()

	// Configurar tamanho máximo de upload (10MB)
	router.MaxMultipartMemory = 10 << 20 // 10 MB

	// Middleware global para CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Endpoint para Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	// @Summary      Verificar saúde da API
	// @Description  Retorna o status de saúde da API
	// @Tags         system
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  dto.HealthResponse
	// @Router       /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, dto.HealthResponse{Status: "ok"})
	})

	// Endpoints de sistema
	system := router.Group("/system")
	{
		system.GET("/kafka", systemHandler.KafkaStatus)
	}

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Usuários
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			// Documentos por usuário
			users.POST("/:id/documents", middleware.ProcessArrayFields(), documentHandler.Create)
			users.GET("/:id/documents", documentHandler.GetByUserID)
		}

		// Documentos
		documents := v1.Group("/documents")
		{
			documents.GET("", documentHandler.List)
			documents.GET("/:id", documentHandler.GetByID)
			documents.GET("/:id/download", documentHandler.DownloadDocument)
			documents.PUT("/:id/status", documentHandler.UpdateStatus)
			documents.DELETE("/:id", documentHandler.Delete)
		}
	}

	return router
}
