# Clean Architecture Go Template

Template de API REST em Go seguindo os princípios da Clean Architecture, SOLID e as melhores práticas de desenvolvimento.

## 📋 Índice

- [Características](#-características)
- [Arquitetura](#-arquitetura)
- [Tecnologias](#-tecnologias)
- [Pré-requisitos](#-pré-requisitos)
- [Instalação](#-instalação)
- [Configuração](#-configuração)
- [Executando a Aplicação](#-executando-a-aplicação)
- [Desenvolvimento](#-desenvolvimento)
- [Testes](#-testes)
- [Documentação da API](#-documentação-da-api)
- [Deployment](#-deployment)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Padrões e Convenções](#-padrões-e-convenções)
- [Contribuindo](#-contribuindo)
- [Licença](#-licença)

## ✨ Características

- 🏗️ **Clean Architecture** - Separação clara de responsabilidades em camadas
- 🔌 **Injeção de Dependências** - Usando Uber FX para IoC container
- 📚 **RESTful API** - Seguindo as melhores práticas REST
- 📝 **Swagger Documentation** - Documentação automática da API
- 🔐 **JWT Authentication** - Autenticação segura com tokens
- 🚀 **Alta Performance** - Otimizado para baixa latência
- 🐳 **Docker Ready** - Containerização completa
- ☸️ **Kubernetes Ready** - Manifestos para deploy em K8s
- 🧪 **Testes Completos** - Unitários, integração e E2E
- 📊 **Observabilidade** - Logs estruturados, métricas e traces
- 🔄 **CI/CD** - Pipelines automatizados com GitHub Actions
- 🛡️ **Segurança** - Middlewares de segurança e rate limiting

## 🏛️ Arquitetura

Este projeto segue os princípios da **Clean Architecture** proposta por Robert C. Martin (Uncle Bob):

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Presentation Layer                         │
│   ┌─────────────┐  ┌──────────────┐  ┌─────────────┐  ┌──────────────┐  │
│   │   Handlers  │  │  Middlewares │  │  Presenters │  │     DTOs     │  │
│   └─────────────┘  └──────────────┘  └─────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
                                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              Application Layer                          │
│           ┌─────────────┐  ┌──────────────┐  ┌─────────────┐            │
│           │  Use Cases  │  │   Services   │  │ Interfaces  │            │
│           └─────────────┘  └──────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────────────┘
                                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                               Domain Layer                              │
│           ┌─────────────┐  ┌──────────────┐  ┌─────────────┐            │
│           │  Entities   │  │Value Objects │  │Domain Rules │            │
│           └─────────────┘  └──────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────────────┘
                                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                               Infrastructure Layer                      │
│    ┌─────────────┐  ┌──────────────┐  ┌─────────────┐  ┌──────────────┐ │
│    │Repositories │  │   Adapters   │  │  Database   │  │External APIs │ │
│    └─────────────┘  └──────────────┘  └─────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
```

### Princípios SOLID

- **S**ingle Responsibility: Cada módulo/classe tem uma única responsabilidade
- **O**pen/Closed: Código aberto para extensão, fechado para modificação
- **L**iskov Substitution: Interfaces bem definidas e substituíveis
- **I**nterface Segregation: Interfaces pequenas e específicas
- **D**ependency Inversion: Dependências apontam sempre para abstrações

## 🛠️ Tecnologias

### Core
- **[Go 1.21+](https://golang.org/)** - Linguagem de programação
- **[Gin](https://gin-gonic.com/)** - Framework HTTP
- **[Uber FX](https://uber-go.github.io/fx/)** - Framework de injeção de dependências

### Banco de Dados
- **[PostgreSQL](https://www.postgresql.org/)** - Banco de dados principal
- **[GORM](https://gorm.io/)** - ORM para Go
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Ferramenta de migração

### Cache e Mensageria
- **[Redis](https://redis.io/)** - Cache e rate limiting
- **[RabbitMQ](https://www.rabbitmq.com/)** - Message broker (opcional)
- **[Apache Kafka](https://kafka.apache.org/)** - Streaming de eventos (opcional)

### Observabilidade
- **[Zerolog](https://github.com/rs/zerolog)** - Logging estruturado
- **[Prometheus](https://prometheus.io/)** - Métricas
- **[Jaeger](https://www.jaegertracing.io/)** - Distributed tracing

### Desenvolvimento
- **[Swagger](https://swagger.io/)** - Documentação da API
- **[Air](https://github.com/cosmtrek/air)** - Hot reload
- **[Mockery](https://github.com/vektra/mockery)** - Geração de mocks
- **[golangci-lint](https://golangci-lint.run/)** - Linter

## 📋 Pré-requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- PostgreSQL 14+ (ou usar via Docker)
- Redis 6+ (ou usar via Docker)
- Make (opcional, mas recomendado)

## 🚀 Instalação

### 1. Clone o repositório

```bash
git clone https://github.com/raykavin/go-clean-template
cd seu-projeto
```

### 2. Copie as variáveis de ambiente

```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

### 3. Instale as dependências

```bash
go mod download
go mod tidy
```

### 4. Configure o banco de dados

```bash
# Via Docker Compose
docker-compose up -d postgres redis

# Ou configure manualmente e ajuste o .env
```

### 5. Execute as migrações

```bash
make migrate-up

# Ou manualmente
go run cmd/migration/main.go up
```

## ⚙️ Configuração

### Variáveis de Ambiente

```env
# Application
APP_NAME=clean-arch-api
APP_ENV=development
APP_PORT=8080
APP_VERSION=1.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cleanarch_db
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-key
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# External Services
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-password

# AWS (opcional)
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
S3_BUCKET=
```

### Arquivo de Configuração

O projeto também suporta configuração via arquivo YAML:

```yaml
# config/config.yml
app:
  name: "Clean Architecture API"
  version: "1.0.0"
  environment: "development"
  
database:
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "postgres"
    dbname: "cleanarch_db"
```

## 🏃‍♂️ Executando a Aplicação

### Desenvolvimento

```bash
# Com hot reload (recomendado)
make run-dev

# Ou usando air diretamente
air

# Ou executando diretamente
go run cmd/api/main.go
```

### Produção

```bash
# Build da aplicação
make build

# Executar binário
./bin/api

# Ou via Docker
docker-compose up api
```

### Docker Compose

```bash
# Subir todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f api

# Parar serviços
docker-compose down
```

## 💻 Desenvolvimento

### Criando um Novo Módulo

1. **Defina a Entidade** (`internal/entity/`)

```go
// internal/entity/book.go
type Book struct {
    ID        uint
    Title     string
    Author    string
    ISBN      string
    CreatedAt time.Time
}
```

2. **Crie o Use Case** (`internal/usecase/book/`)

```go
// internal/usecase/book/create_book.go
type CreateBook struct {
    bookRepo BookRepository
}

func (uc *CreateBook) Execute(ctx context.Context, input CreateBookInput) (*entity.Book, error) {
    // Implementar lógica
}
```

3. **Implemente o Repository** (`internal/repository/`)

```go
// internal/repository/book_repository.go
type bookRepository struct {
    db *gorm.DB
}

func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
    // Implementar persistência
}
```

4. **Adicione o Handler** (`internal/http/handler/`)

```go
// internal/http/handler/book.go
func (h *BookHandler) CreateBook(c *gin.Context) {
    // Implementar endpoint
}
```

5. **Configure a Injeção de Dependências** (`internal/fx/module/`)

### Comandos Make Úteis

```bash
make help         # Ver todos os comandos disponíveis
make run          # Executar aplicação
make test         # Executar testes
make lint         # Executar linter
make swagger      # Gerar documentação Swagger
make migrate-up   # Aplicar migrações
make migrate-down # Reverter migrações
make build        # Build da aplicação
make docker-build # Build da imagem Docker
```

## 🧪 Testes

### Executando Testes

```bash
# Todos os testes
make test

# Com cobertura
make test-coverage

# Testes unitários
go test ./internal/usecase/...

# Testes de integração
go test ./test/integration/...

# Testes E2E
go test ./test/e2e/...
```

### Gerando Mocks

```bash
# Gerar todos os mocks
make generate-mocks

# Ou manualmente
mockery --name=UserRepository --dir=internal/usecase --output=mocks
```

### Exemplo de Teste

```go
func TestCreateUser_Execute(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.UserRepository)
    useCase := user.NewCreateUser(mockRepo)
    
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
    
    // Act
    result, err := useCase.Execute(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

## 📚 Documentação da API

### Swagger

A documentação Swagger está disponível em:

```
http://localhost:8080/swagger/index.html
```

Para regenerar a documentação:

```bash
make swagger
```

### Exemplos de Requisições

#### Autenticação

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"senha123"}'

# Response
{
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "user@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Criar Usuário

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {token}" \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"senha123"}'
```

## 🚀 Deployment

### Docker

```bash
# Build da imagem
docker build -t meu-app:latest .

# Executar container
docker run -p 8080:8080 --env-file .env meu-app:latest
```

### Kubernetes

```bash
# Aplicar configurações
kubectl apply -f deployments/kubernetes/

# Verificar deployment
kubectl get pods -n meu-namespace

# Ver logs
kubectl logs -f deployment/api -n meu-namespace
```

### CI/CD com GitHub Actions

O projeto inclui workflows para:

- **CI**: Testes, linting e build em cada push
- **CD**: Deploy automático em merges para main
- **Security**: Scan de vulnerabilidades

## 📏 Padrões e Convenções

### Nomenclatura

- **Arquivos**: `snake_case` (ex: `user_repository.go`)
- **Interfaces**: `PascalCase` (ex: `UserRepository`)
- **Structs**: `PascalCase` (ex: `CreateUserUseCase`)
- **Funções/Métodos**: `PascalCase` para públicos, `camelCase` para privados

### Commits

[Conventional Commits](https://www.conventionalcommits.org/):

```
feat(user): add user registration endpoint
fix(auth): correct JWT expiration time
docs(readme): update installation instructions
chore(deps): update gin to v1.9.1
```

### Branches

- `main` - Produção
- `develop` - Desenvolvimento
- `feature/*` - Novas funcionalidades
- `fix/*` - Correções
- `hotfix/*` - Correções urgentes

### Checklist do PR

- [ ] Testes passando
- [ ] Código seguindo padrões (lint)
- [ ] Documentação atualizada
- [ ] Commits seguindo convenção

---
