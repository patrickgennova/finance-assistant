package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"finance-assistant/config"
	"finance-assistant/internal/domain/entity"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Producer representa um produtor Kafka
type Producer struct {
	producer *kafka.Producer
	topic    string
}

// DocumentMessage representa a mensagem que será enviada para o Kafka
type DocumentMessage struct {
	ID           string    `json:"id"`
	ExternalID   string    `json:"external_id"`
	UserID       string    `json:"user_id"`
	DocumentType string    `json:"document_type"`
	Filename     string    `json:"filename"`
	ContentType  string    `json:"content_type"`
	FileContent  string    `json:"file_content"`
	Categories   []string  `json:"categories"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewProducer cria um novo produtor Kafka
func NewProducer(cfg *config.Config) (*Producer, error) {
	// Configurações adicionais para o produtor
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.KafkaBrokers[0],
		"client.id":                "finance-assistant",
		"acks":                     "all",
		"delivery.timeout.ms":      "30000",  // 30 segundos timeout
		"request.timeout.ms":       "15000",  // 15 segundo timeout
		"message.max.bytes":        16777216, // 16MB
		"socket.keepalive.enable":  "true",
		"socket.max.fails":         "3",
		"reconnect.backoff.ms":     "100",
		"reconnect.backoff.max.ms": "10000",
		"retry.backoff.ms":         "100",
		"compression.type":         "gzip",
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar produtor Kafka: %w", err)
	}

	// Configurar handler para eventos de entrega
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Falha ao entregar mensagem ao Kafka: %v", ev.TopicPartition.Error)
				} else {
					log.Printf("Mensagem enviada com sucesso para Kafka: %v", ev.TopicPartition)
				}
			case kafka.Error:
				log.Printf("Erro Kafka: %v", ev)
			}
		}
	}()

	log.Printf("Produtor Kafka conectado com sucesso ao broker: %s", cfg.KafkaBrokers[0])
	log.Printf("Tópico configurado: %s", cfg.KafkaTopic)

	return &Producer{
		producer: producer,
		topic:    cfg.KafkaTopic,
	}, nil
}

// SendDocument envia um documento para o Kafka
func (p *Producer) SendDocument(document *entity.Document) error {
	log.Printf("Enviando documento %s para processamento...", document.ExternalID)

	// Converter documento para mensagem
	message := DocumentMessage{
		ID:           fmt.Sprintf("%d", document.ID),
		ExternalID:   document.ExternalID.String(),
		UserID:       fmt.Sprintf("%d", document.UserID),
		DocumentType: document.DocumentType,
		Filename:     document.Filename,
		ContentType:  document.ContentType,
		FileContent:  document.FileContent,
		Categories:   document.Categories,
		CreatedAt:    document.CreatedAt,
		UpdatedAt:    document.UpdatedAt,
	}

	// Converter mensagem para JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	// Criar tópico se não existir (opcional)
	// Isso requer permissões administrativas no Kafka
	/*
		adminClient, err := kafka.NewAdminClientFromProducer(p.producer)
		if err == nil {
			topicResults, err := adminClient.CreateTopics(
				[]kafka.TopicSpecification{{
					Topic:             p.topic,
					NumPartitions:     1,
					ReplicationFactor: 1,
				}},
				kafka.SetAdminOperationTimeout(time.Second*10),
			)
			if err != nil {
				log.Printf("Aviso: Não foi possível criar tópico: %v", err)
			}
			for _, result := range topicResults {
				if result.Error.Code() != kafka.ErrNoError &&
				   result.Error.Code() != kafka.ErrTopicAlreadyExists {
					log.Printf("Aviso: Erro ao criar tópico %s: %v", result.Topic, result.Error.String())
				}
			}
			adminClient.Close()
		}
	*/

	// Enviar mensagem com timeout
	deliveryChan := make(chan kafka.Event, 1)
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: messageJSON,
		Key:   []byte(document.ExternalID.String()),
		Headers: []kafka.Header{
			{
				Key:   "content_type",
				Value: []byte("application/json"),
			},
			{
				Key:   "source",
				Value: []byte("finance-assistant"),
			},
		},
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("erro ao produzir mensagem: %w", err)
	}

	// Aguardar confirmação de entrega com timeout
	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("erro ao entregar mensagem: %w", m.TopicPartition.Error)
		}
		log.Printf("Documento %s enviado com sucesso para o tópico %s [%d] @ %v",
			document.ExternalID, *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout ao aguardar confirmação de entrega")
	}

	return nil
}

// CheckKafkaConnection verifica se a conexão com o Kafka está funcionando
func (p *Producer) CheckKafkaConnection() error {
	// Metadados para verificar a conexão
	metadata, err := p.producer.GetMetadata(nil, true, 5000)
	if err != nil {
		return fmt.Errorf("erro ao obter metadados do Kafka: %w", err)
	}

	log.Printf("Conexão com Kafka verificada. Brokers disponíveis: %d", len(metadata.Brokers))
	return nil
}

// Close fecha o produtor
func (p *Producer) Close() {
	// Liberar mensagens pendentes antes de fechar
	remain := p.producer.Flush(5000) // 5 segundos timeout
	if remain > 0 {
		log.Printf("Aviso: %d mensagens ainda não foram entregues ao fechar o produtor", remain)
	}
	p.producer.Close()
	log.Println("Produtor Kafka fechado")
}
