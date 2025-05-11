# Finance Assistant API

### Description

Finance Assistant is a RESTful API developed in Go to serve as a personal financial assistant. This project was created as part of my Go learning journey, focusing on good programming practices, Clean Architecture, and Domain-Driven Design (DDD) principles.

The API allows users to upload financial documents (bank statements, invoices, and receipts), processes these documents, and later extracts transactions to provide financial insights.

### Main Features

- **Complete user CRUD**: User management with basic data such as name, email, and phone.
- **Financial document processing**: Upload, storage, and processing of documents.
- **Kafka integration**: Messaging system for asynchronous document processing.
- **Document analysis with AI**: Planned integration to extract and categorize transactions from documents.

### Technologies and Patterns

- **Go**: Main development language
- **PostgreSQL**: Relational database
- **Kafka**: Messaging system for asynchronous processing
- **Gin**: Web framework for the API
- **SQLX**: Library for database access
- **Swagger/OpenAPI**: API documentation
- **Clean Architecture**: Layered separation (Domain, Repository, Service, Handler)
- **Domain-Driven Design (DDD)**: Code structure based on the problem domain

### Project Structure

```
finance-assistant/
├── cmd/                # Application entry points
│   └── api/            # API application
│       └── main.go     # Main entry point
├── config/             # Application configurations
├── internal/           # Internal application code
│   ├── domain/         # Domain layer (entities and business rules)
│   │   ├── entity/     # Domain entities
│   │   ├── repository/ # Repository interfaces
│   │   └── service/    # Domain services
│   ├── infrastructure/ # Infrastructure implementations
│   │   ├── database/   # Database configuration
│   │   ├── kafka/      # Kafka configuration
│   │   └── repository/ # Concrete repository implementations
│   └── interface/      # External interfaces
│       ├── api/        # HTTP API
│       │   ├── dto/    # Data transfer objects
│       │   ├── handler/# HTTP handlers
│       │   └── middleware/ # Middlewares
│       └── http/       # HTTP configuration
├── migrations/         # Database migrations
├── pkg/                # Shared packages
└── docs/               # Additional documentation
```

### How to Run

1. Clone the repository:
```bash
git clone https://github.com/yourusername/finance-assistant.git
cd finance-assistant
```

2. Configure environment variables:
```bash
cp .env.example .env
# Edit the .env file with your settings
```

3. Start services with Docker:
```bash
make docker-up
```

4. Run the application:
```bash
make run
```

5. Access Swagger documentation:
```
http://localhost:8080/swagger/index.html
```

### Development

This project was developed as part of my learning journey in Go, focusing on:

1. **Clean Architecture**: Separation of code into layers with controlled dependencies
2. **Domain-Driven Design**: Structuring code based on the problem domain
3. **REST APIs**: Implementing an API following RESTful practices
4. **Concurrency**: Using goroutines and channels for asynchronous operations
5. **Database integration**: Using PostgreSQL with the SQLX library
6. **Messaging**: Kafka integration for asynchronous processing

### Next Steps

- Complete implementation of the document processing module
- AI integration to extract and categorize transactions
- Adding authentication and authorization
- Implementation of a dashboard for visualizing financial insights
- Unit and integration tests

### Contributions

This is a personal learning project, but contributions are welcome through issues and pull requests.

### License

This project is licensed under the MIT License.