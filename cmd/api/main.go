package main

import (
	"context"
	"finance-assistant/internal/infrastructure/kafka"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"finance-assistant/config"
	_ "finance-assistant/docs"
	"finance-assistant/internal/domain/service"
	"finance-assistant/internal/infrastructure/database"
	repo "finance-assistant/internal/infrastructure/repository"
	"finance-assistant/internal/interface/api/handler"
	"finance-assistant/internal/interface/http"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title           Finance Assistant API
// @version         1.0
// @description     API para gerenciamento de finanças pessoais com análise de documentos
// @termsOfService  http://swagger.io/terms/

// @contact.name   Seu Nome
// @contact.url    http://seusite.com
// @contact.email  seu.email@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	// Carregar a configuração
	cfg := config.LoadConfig()

	// Conectar ao banco de dados
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Executar migrações
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Erro ao inicializar driver de migração: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"finance_assistant", driver)
	if err != nil {
		log.Fatalf("Erro ao configurar migrações: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erro ao executar migrações: %v", err)
	}

	log.Println("Migrações aplicadas com sucesso")

	// Inicializar produtor Kafka
	kafkaProducer, err := kafka.NewProducer(cfg)
	if err != nil {
		log.Printf("Aviso: Falha ao conectar ao Kafka: %v", err)
		log.Println("Continuando sem suporte ao Kafka - os documentos serão salvos mas não processados")
		kafkaProducer = nil
	} else {
		// Verificar conexão com Kafka
		if err := kafkaProducer.CheckKafkaConnection(); err != nil {
			log.Printf("Aviso: Verificação de conexão Kafka falhou: %v", err)
			log.Println("Continuando sem suporte ao Kafka - os documentos serão salvos mas não processados")
			kafkaProducer.Close()
			kafkaProducer = nil
		} else {
			defer kafkaProducer.Close()
		}
	}

	// Inicializar repositórios
	userRepo := repo.NewPostgresUserRepository(db)
	documentRepo := repo.NewPostgresDocumentRepository(db)

	// Inicializar serviços
	userService := service.NewUserService(userRepo)
	documentService := service.NewDocumentService(documentRepo, userRepo, kafkaProducer)

	// Inicializar handlers
	userHandler := handler.NewUserHandler(userService)
	documentHandler := handler.NewDocumentHandler(documentService)
	systemHandler := handler.NewSystemHandler(kafkaProducer)

	// Configurar o router
	router := inhttp.SetupRouter(userHandler, documentHandler, systemHandler)

	// Iniciar servidor HTTP
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Iniciar o servidor em uma goroutine
	go func() {
		log.Printf("Servidor iniciado na porta %d\n", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Configurar graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar servidor: %v", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
