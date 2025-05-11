package handler

import (
	"net/http"

	"finance-assistant/internal/infrastructure/kafka"
	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	kafkaProducer *kafka.Producer
}

func NewSystemHandler(kafkaProducer *kafka.Producer) *SystemHandler {
	return &SystemHandler{
		kafkaProducer: kafkaProducer,
	}
}

func (h *SystemHandler) KafkaStatus(c *gin.Context) {
	status := "ok"
	message := "Conexão com Kafka estabelecida"

	if h.kafkaProducer == nil {
		status = "unavailable"
		message = "Produtor Kafka não inicializado"
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  status,
			"message": message,
		})
		return
	}

	err := h.kafkaProducer.CheckKafkaConnection()
	if err != nil {
		status = "error"
		message = err.Error()
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  status,
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": message,
	})
}
