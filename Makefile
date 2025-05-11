# Makefile
.PHONY: build run test migrate-up migrate-down docker-up docker-down

# Variáveis
APP_NAME=finance-assistant
DOCKER_COMPOSE=docker-compose

# Build da aplicação
build:
	go build -o bin/$(APP_NAME) ./cmd/api

# Executar a aplicação
run:
	go run ./cmd/api/main.go

# Executar testes
test:
	go test -v ./...

# Iniciar infraestrutura com Docker
docker-up:
	$(DOCKER_COMPOSE) up -d

# Parar infraestrutura com Docker
docker-down:
	$(DOCKER_COMPOSE) down

# Executar migrações (up)
migrate-up:
	migrate -path migrations -database "postgresql://finance_assistant:finance_assistant@localhost:5432/finance_assistant?sslmode=disable" up

# Reverter migrações (down)
migrate-down:
	migrate -path migrations -database "postgresql://finance_assistant:finance_assistant@localhost:5432/finance_assistant?sslmode=disable" down

# Criar nova migração
migrate-create:
	@read -p "Nome da migração: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Limpar binários
clean:
	rm -rf bin/

# Instalar dependências
deps:
	go mod download

# Verificar código
lint:
	go vet ./...
	golangci-lint run

# Gerar documentação Swagger
swagger:
	$(HOME)/go/bin/swag init -g cmd/api/main.go