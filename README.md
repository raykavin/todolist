# Todo List Application

A scalable Todo List application built with Go, following Domain-Driven Design (DDD) principles and clean architecture patterns.

## 🚀 Features

- **User Management**: Complete authentication and authorization system with JWT tokens
- **Todo Management**: Create, read, update, delete, and complete todos
- **Person Management**: Associate todos with persons
- **Statistics**: Track todo completion rates and daily statistics
- **Tags**: Organize todos with tags
- **Priority System**: Set priorities for todos
- **OIDC Support**: OpenID Connect authentication integration
- **Multi-Database Support**: PostgreSQL, MySQL, and SQLite
- **API Documentation**: Swagger/OpenAPI documentation
- **Comprehensive Testing**: Unit tests

## 🏗️ Architecture

This application follows clean architecture principles with the following layers:

- **Domain Layer**: Core business logic and entities
- **Use Case Layer**: Application-specific business rules
- **Infrastructure Layer**: External frameworks, databases, and tools
- **Delivery Layer**: HTTP handlers and presenters

### Project Structure

# Project Structure Documentation

```
├── cmd/                                      # Application entry points (main packages)
│   ├── api/                                  # HTTP API server application
│   │   └── main.go                           # API server entry point - initializes DI container and starts HTTP server
│   ├── migration/                            # Database migration application
│   │   └── main.go                           # Migration runner - executes database schema migrations
│   └── worker/                               # Background worker application
│       └── main.go                           # Worker entry point - runs background jobs and async tasks
│
├── config/                                   # Configuration files
│   ├── config.dev.yml                        # Development environment configuration
│   ├── config.prod.yml                       # Production environment configuration
│   ├── config.test.yml                       # Test environment configuration
│   └── config.yml                            # Base/default configuration (inherited by others)
│
├── deploy/                                   # Deployment configurations and scripts
│   ├── docker/                               # Docker-related files
│   │   ├── Dockerfile                        # Production Docker image definition
│   │   ├── Dockerfile.dev                    # Development Docker image with hot reload
│   │   └── healthcheck.sh                    # Docker health check script
│   ├── kubernetes/                           # Kubernetes manifests
│   │   ├── configmap.yaml                    # K8s ConfigMap for app configuration
│   │   ├── deployment.yaml                   # K8s Deployment definition
│   │   └── service.yaml                      # K8s Service to expose the app
│   └── terraform/                            # Infrastructure as Code
│       ├── main.tf                           # Main Terraform configuration
│       └── variables.tf                      # Terraform variable definitions
│
├── docker-compose.dev.yml                    # Development environment orchestration
├── docker-compose.test.yml                   # Test environment orchestration
├── docker-compose.yml                        # Production environment orchestration
│
├── docs/                                     # Documentation and diagrams
│   ├── diagram.xml                           # Architecture diagrams (XML format)
│   ├── docs.go                               # Go file for Swagger documentation generation
│   ├── swagger.json                          # Generated OpenAPI specification (JSON)
│   ├── swagger.yaml                          # Generated OpenAPI specification (YAML)
│   ├── Todo_List_App_Mermaid_Chart.mmd       # Mermaid diagram source
│   └── Todo_List_App_Mermaid_Chart.png       # Generated diagram image
│
├── go.mod                                    # Go module definition and dependencies
├── go.sum                                    # Go module checksums for dependencies
│
├── internal/                                 # Private application code (not importable by external packages)
│   ├── adapter/                              # Interface adapters (Hexagonal Architecture)
│   │   ├── delivery/                         # Delivery mechanisms (how users interact with the app)
│   │   │   └── http/                         # HTTP delivery adapter
│   │   │       ├── gin.go                    # Gin framework setup and router configuration
│   │   │       ├── handler/                  # HTTP request handlers
│   │   │       │   ├── auth.go               # Authentication endpoints (login, logout, refresh)
│   │   │       │   ├── health.go             # Health check and readiness endpoints
│   │   │       │   ├── helper.go             # Common handler utilities and helpers
│   │   │       │   ├── person.go             # Person CRUD endpoints
│   │   │       │   └── todo.go               # Todo CRUD and business operation endpoints
│   │   │       ├── interfaces.go             # Handler interfaces definitions
│   │   │       ├── middleware/               # HTTP middleware components
│   │   │       │   ├── auth.go               # JWT authentication middleware
│   │   │       │   ├── cors.go               # CORS configuration middleware
│   │   │       │   ├── log.go                # Request/response logging middleware
│   │   │       │   ├── oidc.go               # OpenID Connect middleware
│   │   │       │   ├── rate_limit.go         # API rate limiting middleware
│   │   │       │   └── recovery.go           # Panic recovery middleware
│   │   │       ├── presenter/                # Response formatters and presenters
│   │   │       │   ├── login.go              # Login response formatting
│   │   │       │   ├── person.go             # Person entity response formatting
│   │   │       │   ├── todo.go               # Todo entity response formatting
│   │   │       │   └── user.go               # User entity response formatting
│   │   │       └── wrappers.go               # HTTP handler wrappers and utilities
│   │   ├── notification/                     # Notification adapters
│   │   │   ├── email.go                      # Email notification implementation
│   │   │   ├── interfaces.go                 # Notification interfaces
│   │   │   ├── push.go                       # Push notification implementation
│   │   │   └── sms.go                        # SMS notification implementation
│   │   └── storage/                          # File storage adapters
│   │       ├── interfaces.go                 # Storage interfaces
│   │       ├── local.go                      # Local filesystem storage
│   │       └── s3.go                         # AWS S3 storage implementation
│   │
│   ├── config/                               # Configuration structures and loaders
│   │   ├── application.go                    # Application-wide configuration
│   │   ├── cache.go                          # Cache configuration (Redis, memory)
│   │   ├── config.go                         # Main configuration aggregator
│   │   ├── database.go                       # Database connection configuration
│   │   ├── interfaces.go                     # Configuration interfaces
│   │   ├── jwt.go                            # JWT token configuration
│   │   ├── oidc.go                           # OpenID Connect configuration
│   │   └── web.go                            # Web server configuration
│   │
│   ├── di/                                   # Dependency injection container setup
│   │   ├── application_services.go           # Application service wiring
│   │   ├── core.go                           # Core DI container setup
│   │   ├── databases.go                      # Database connection injection
│   │   ├── domain_services.go                # Domain service wiring
│   │   ├── handlers.go                       # HTTP handler injection
│   │   ├── http_server.go                    # HTTP server setup
│   │   ├── logger.go                         # Logger injection
│   │   ├── repositories.go                   # Repository injection
│   │   └── usecases.go                       # Use case injection
│   │
│   ├── domain/                               # Domain layer (business logic and rules)
│   │   ├── person/                           # Person aggregate
│   │   │   ├── entity/                       # Person entities
│   │   │   │   └── person.go                 # Person entity definition
│   │   │   ├── repository/                   # Person repository interface
│   │   │   │   └── person.go                 # Person repository contract
│   │   │   └── valueobject/                  # Person value objects
│   │   │       ├── email.go                  # Email value object
│   │   │       ├── email_test.go             # Email value object tests
│   │   │       ├── tax_id.go                 # Tax ID value object
│   │   │       └── tax_id_test.go            # Tax ID value object tests
│   │   ├── shared/                           # Shared domain components
│   │   │   ├── entity.go                     # Base entity with common fields
│   │   │   ├── repository.go                 # Base repository interface
│   │   │   └── valueobject/                  # Shared value objects
│   │   │       ├── date.go                   # Date value object
│   │   │       ├── date_test.go              # Date value object tests
│   │   │       └── priority.go               # Priority enumeration
│   │   ├── todo/                             # Todo aggregate
│   │   │   ├── entity/                       # Todo entities
│   │   │   │   ├── todo.go                   # Todo entity definition
│   │   │   │   └── todo_test.go              # Todo entity tests
│   │   │   ├── repository/                   # Todo repository interface
│   │   │   │   └── repository.go             # Todo repository contract
│   │   │   ├── service/                      # Todo domain services
│   │   │   │   ├── todo.go                   # Todo business logic service
│   │   │   │   └── todo_test.go              # Todo service tests
│   │   │   └── valueobject/                  # Todo value objects
│   │   │       ├── tag_count.go              # Tag count value object
│   │   │       ├── todo_description.go       # Todo description value object
│   │   │       ├── todo_filter_criteria.go   # Filter criteria for queries
│   │   │       ├── todo_statistics.go        # Todo statistics value object
│   │   │       ├── todo_status.go            # Todo status enumeration
│   │   │       └── todo_title.go             # Todo title value object
│   │   └── user/                             # User aggregate
│   │       ├── entity/                       # User entities
│   │       │   └── user.go                   # User entity definition
│   │       ├── repository/                   # User repository interface
│   │       │   └── user.go                   # User repository contract
│   │       └── valueobject/                  # User value objects
│   │           ├── password.go               # Password value object with hashing
│   │           ├── password_test.go          # Password value object tests
│   │           ├── user_role.go              # User role enumeration
│   │           └── user_status.go            # User status enumeration
│   │
│   ├── dto/                                  # Data Transfer Objects
│   │   ├── auth.go                           # Authentication DTOs (login, token)
│   │   ├── http_response.go                  # Standard HTTP response wrapper
│   │   ├── person.go                         # Person request/response DTOs
│   │   ├── todo.go                           # Todo request/response DTOs
│   │   └── user.go                           # User request/response DTOs
│   │
│   ├── infrastructure/                       # Infrastructure layer implementations
│   │   ├── auth/                             # Authentication infrastructure
│   │   │   ├── jwt_token_adapter.go          # JWT token generation and validation
│   │   │   ├── oauth2.go                     # OAuth2 client implementation
│   │   │   └── oidc.go                       # OpenID Connect implementation
│   │   ├── cache/                            # Cache implementations
│   │   │   ├── distributed.go                # Distributed cache interface
│   │   │   ├── memory.go                     # In-memory cache implementation
│   │   │   └── redis.go                      # Redis cache implementation
│   │   ├── database/                         # Database infrastructure
│   │   │   ├── helper.go                     # Database helper functions
│   │   │   ├── mapper/                       # Entity-Model mappers
│   │   │   │   ├── person.go                 # Person entity <-> model mapping
│   │   │   │   ├── todo.go                   # Todo entity <-> model mapping
│   │   │   │   └── user.go                   # User entity <-> model mapping
│   │   │   ├── migrate_default.go            # Default migration runner
│   │   │   ├── model/                        # Database models (GORM)
│   │   │   │   ├── audit.go                  # Audit log model
│   │   │   │   ├── login_attempt.go          # Login attempt tracking model
│   │   │   │   ├── person.go                 # Person table model
│   │   │   │   ├── tag.go                    # Tag table model
│   │   │   │   ├── todo_daily_statistics.go  # Daily statistics model
│   │   │   │   ├── todo.go                   # Todo table model
│   │   │   │   ├── todo_tag.go               # Todo-tag relation model
│   │   │   │   ├── todo_view.go              # Todo materialized view model
│   │   │   │   ├── user.go                   # User table model
│   │   │   │   └── user_statistics_view.go   # User statistics view model
│   │   │   ├── repository/                   # Repository implementations
│   │   │   │   ├── person.go                 # Person repository implementation
│   │   │   │   ├── person_query.go           # Person complex queries
│   │   │   │   ├── todo.go                   # Todo repository implementation
│   │   │   │   ├── todo_query.go             # Todo complex queries
│   │   │   │   ├── user.go                   # User repository implementation
│   │   │   │   └── user_query.go             # User complex queries
│   │   │   └── seed_default.go               # Database seeding for development
│   │   ├── email/                            # Email infrastructure
│   │   │   └── email.go                      # Email service implementation (SMTP)
│   │   ├── logger/                           # Logging infrastructure
│   │   │   └── smart.go                      # Smart logger with context
│   │   ├── messaging/                        # Message queue infrastructure
│   │   │   ├── kafka.go                      # Apache Kafka implementation
│   │   │   ├── rabbitmq.go                   # RabbitMQ implementation
│   │   │   └── sqs.go                        # AWS SQS implementation
│   │   ├── push/                             # Push notification infrastructure
│   │   │   └── firebase.go                   # Firebase Cloud Messaging
│   │   └── storage/                          # File storage infrastructure
│   │       ├── gcs.go                        # Google Cloud Storage
│   │       ├── local.go                      # Local file storage
│   │       └── s3.go                         # AWS S3 storage
│   │
│   ├── service/                              # Application services
│   │   ├── token_service.go                  # JWT token management service
│   │   └── user_security.go                  # User security service (auth, permissions)
│   │
│   └── usecase/                              # Use cases (application business rules)
│       ├── person/                           # Person use cases
│       │   ├── create_person.go              # Create person use case
│       │   ├── get_person.go                 # Get person details use case
│       │   └── update_person.go              # Update person use case
│       ├── todo/                             # Todo use cases
│       │   ├── complete_todo.go              # Mark todo as complete use case
│       │   ├── create_todo.go                # Create new todo use case
│       │   ├── delete_todo.go                # Delete todo use case
│       │   ├── get_statistics.go             # Get todo statistics use case
│       │   ├── get_todo.go                   # Get todo details use case
│       │   ├── list_todo.go                  # List todos with filters use case
│       │   └── update_todo.go                # Update todo use case
│       └── user/                             # User use cases
│           ├── change_password.go            # Change user password use case
│           ├── create_user.go                # Create new user use case
│           └── login.go                      # User login use case
│
├── Makefile                                  # Build automation and common tasks
│
├── pkg/                                      # Public packages (can be imported by external projects)
│   ├── auth/                                 # Authentication utilities
│   │   ├── jwt_token.go                      # JWT token utilities
│   │   └── oidc_client.go                    # OIDC client utilities
│   ├── config/                               # Configuration utilities
│   │   └── viper.go                          # Viper configuration loader
│   ├── database/                             # Database utilities
│   │   ├── custom_log.go                     # Custom GORM logger
│   │   ├── gorm.go                           # GORM setup and configuration
│   │   ├── gorm_test.go                      # GORM utilities tests
│   │   ├── mysql.go                          # MySQL specific configuration
│   │   ├── postgres.go                       # PostgreSQL specific configuration
│   │   └── sqlite.go                         # SQLite specific configuration
│   ├── errors/                               # Error handling utilities
│   │   ├── errors.go                         # Custom error types
│   │   └── http.go                           # HTTP error responses
│   ├── log/                                  # Logging utilities
│   │   ├── interfaces.go                     # Logger interfaces
│   │   └── smart.go                          # Smart logger implementation
│   ├── terminal/                             # Terminal utilities
│   │   └── banner.go                         # Application banner display
│   ├── validator/                            # Validation utilities
│   │   ├── messages.go                       # Validation error messages
│   │   ├── rules.go                          # Custom validation rules
│   │   └── validator.go                      # Validator setup
│   └── web/                                  # Web utilities
│       ├── client.go                         # HTTP client utilities
│       └── gin.go                            # Gin framework utilities
│
├── scripts/                                  # Utility scripts
│   ├── build.sh                              # Build application binaries
│   ├── deploy.sh                             # Deployment automation script
│   ├── generate-mocks.sh                     # Generate test mocks
│   ├── setup.sh                              # Initial project setup
│   ├── swagger-doc.sh                        # Generate Swagger documentation
│   └── test.sh                               # Run test suites
│
└── test/                                     # Test suites
    ├── e2e/                                  # End-to-end tests
    │   ├── api_test.go                       # General API E2E tests
    │   ├── auth_flow_test.go                 # Authentication flow E2E tests
    │   ├── person_flow_test.go               # Person CRUD flow E2E tests
    │   ├── todo_flow_test.go                 # Todo CRUD flow E2E tests
    │   └── user_flow_test.go                 # User management flow E2E tests
    └── integration/                          # Integration tests
        ├── person_test.go                    # Person integration tests
        ├── todo_test.go                      # Todo integration tests
        └── user_test.go                      # User integration tests
```

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL/MySQL/SQLite (depending on your choice)
- Redis (optional, for caching)
- Make

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone https://gitlab.mba.corp/templates/todo_example
   cd todo_example
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up configuration**
   ```bash
   cp config/config.yml config/config.dev.yml
   # Edit config.dev.yml with your settings
   ```

4. **Generate documentation**
   ```bash
   make swagger
   ```

## 🚀 Running the Application

### Using Docker Compose

```bash
# Development environment
docker-compose -f docker-compose.dev.yml up

# Production environment
docker-compose up
```

### Running Locally

```bash
# Run the API server
go run cmd/api/main.go
```

### Using Make

```bash
# Build the application
make build

# Run tests
make test

# Generate mocks
make mocks

# Deploy
make deploy
```

## 🔧 Configuration

The application uses YAML configuration files located in the `config/` directory:

- `config.yml`: Base configuration
- `config.dev.yml`: Development environment
- `config.test.yml`: Test environment
- `config.prod.yml`: Production environment

## 📚 API Documentation

API documentation is available via Swagger UI:

1. Start the application
2. Navigate to `http://localhost:8080/swagger/index.html`

### Main Endpoints

#### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/logout` - User logout
- `PUT /api/v1/auth/change-password` - Change password

#### Todos
- `GET /api/v1/todos` - List todos with filters
- `POST /api/v1/todos` - Create new todo
- `GET /api/v1/todos/:id` - Get todo details
- `PUT /api/v1/todos/:id` - Update todo
- `DELETE /api/v1/todos/:id` - Delete todo
- `POST /api/v1/todos/:id/complete` - Mark todo as complete
- `GET /api/v1/todos/statistics` - Get todo statistics

#### People
- `GET /api/v1/people/:id` - Get person details
- `POST /api/v1/people` - Create new person
- `PUT /api/v1/people/:id` - Update person

## 🧪 Testing

The application includes comprehensive test coverage:

```bash
# Run all tests
make test

# Run unit tests
go test ./internal/...

# Run integration tests
go test ./test/integration/...

# Run E2E tests
go test ./test/e2e/...

# Run with coverage
go test -cover ./...
```

## 🚢 Deployment

### Docker

```bash
# Build Docker image
docker build -f deploy/docker/Dockerfile -t todo-app:latest .

# Run container
docker run -p 8080:8080 todo-app:latest
```

## 🔒 Security Features

- JWT-based authentication
- OIDC/OAuth2 support
- Password hashing with bcrypt
- Rate limiting middleware
- CORS configuration
- SQL injection protection via ORM
- Input validation

## 📊 Monitoring

The application includes:
- Health check endpoint: `GET /health`
- Structured logging with different log levels
- Metrics collection ready
- Request/Response logging middleware
---

**Note**: This is a sample application demonstrating clean architecture principles in Go. Adapt security settings and configurations according to your production requirements.