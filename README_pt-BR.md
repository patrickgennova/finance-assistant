# Finance Assistant API

### Descrição

O Finance Assistant é uma API RESTful desenvolvida em Go para servir como um assistente financeiro pessoal. Este projeto foi criado como parte do meu percurso de aprendizado de Go, focando em boas práticas de programação, arquitetura limpa (Clean Architecture) e princípios de Domain-Driven Design (DDD).

A API permite que usuários façam upload de documentos financeiros (extratos bancários, notas fiscais e recibos), processa esses documentos e, posteriormente, extrai transações para oferecer insights financeiros.

### Funcionalidades Principais

- **CRUD completo de usuários**: Gerenciamento de usuários com dados básicos como nome, email e telefone.
- **Processamento de documentos financeiros**: Upload, armazenamento e processamento de documentos.
- **Integração com Kafka**: Sistema de mensageria para processamento assíncrono de documentos.
- **Análise de documentos com IA**: Integração planejada para extrair e categorizar transações de documentos.

### Tecnologias e Padrões

- **Go**: Linguagem principal de desenvolvimento
- **PostgreSQL**: Banco de dados relacional
- **Kafka**: Sistema de mensageria para processamento assíncrono
- **Gin**: Framework web para a API
- **SQLX**: Biblioteca para acesso ao banco de dados
- **Swagger/OpenAPI**: Documentação da API
- **Clean Architecture**: Separação em camadas (Domain, Repository, Service, Handler)
- **Domain-Driven Design (DDD)**: Estruturação do código baseada no domínio do problema

### Estrutura do Projeto

```
finance-assistant/
├── cmd/                # Pontos de entrada da aplicação
│   └── api/            # Aplicação API
│       └── main.go     # Ponto de entrada principal
├── config/             # Configurações da aplicação
├── internal/           # Código interno da aplicação
│   ├── domain/         # Camada de domínio (entidades e regras de negócio)
│   │   ├── entity/     # Entidades de domínio
│   │   ├── repository/ # Interfaces de repositórios
│   │   └── service/    # Serviços de domínio
│   ├── infrastructure/ # Implementações de infraestrutura
│   │   ├── database/   # Configuração de banco de dados
│   │   ├── kafka/      # Configuração de Kafka
│   │   └── repository/ # Implementações concretas de repositórios
│   └── interface/      # Interfaces externas
│       ├── api/        # API HTTP
│       │   ├── dto/    # Objetos de transferência de dados
│       │   ├── handler/# Handlers HTTP
│       │   └── middleware/ # Middlewares
│       └── http/       # Configuração HTTP
├── migrations/         # Migrações de banco de dados
├── pkg/                # Pacotes compartilhados
└── docs/               # Documentação adicional
```

### Como Executar

1. Clone o repositório:
```bash
git clone https://github.com/seuusuario/finance-assistant.git
cd finance-assistant
```

2. Configure as variáveis de ambiente:
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

3. Inicie os serviços com Docker:
```bash
docker-compose up -d
```

4. Execute a aplicação:
```bash
make run
```

5. Acesse a documentação Swagger:
```
http://localhost:8080/swagger/index.html
```

### Desenvolvimento

Este projeto foi desenvolvido como parte do meu aprendizado em Go, com foco em:

1. **Arquitetura Limpa**: Separação do código em camadas com dependências controladas
2. **Domain-Driven Design**: Estruturação do código baseada no domínio do problema
3. **REST APIs**: Implementação de uma API seguindo práticas RESTful
4. **Concorrência**: Uso de goroutines e canais para operações assíncronas
5. **Integração com bancos de dados**: Uso do PostgreSQL com a biblioteca SQLX
6. **Mensageria**: Integração com Kafka para processamento assíncrono

### Próximos Passos

- Implementação completa do módulo de processamento de documentos
- Integração com IA para extrair e categorizar transações
- Adição de autenticação e autorização
- Implementação de dashboard para visualização de insights financeiros
- Testes unitários e de integração

### Contribuições

Este é um projeto pessoal de aprendizado, mas contribuições são bem-vindas através de issues e pull requests.

### Licença

Este projeto está licenciado sob a licença MIT.