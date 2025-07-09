# Projeto Golang

## Visão Geral

Este template implementa uma arquitetura limpa e modular para aplicações Go, seguindo os princípios SOLID e as melhores práticas da linguagem. Utiliza:

* **Uber FX** para injeção de dependências
* **Gin** como framework HTTP
* **GORM** ou drivers nativos para persistência
* **Swagger** para documentação automática da API
* **Viper** para leitura de arquivos de configuração e/ou variáveis de ambiente
* **Clean Architecture** para separação de responsabilidades

## Princípios e Conceitos

### Clean Architecture

A arquitetura limpa é uma abordagem para estruturar aplicações de software, visando facilitar a manutenção, a testabilidade e a independência de tecnologias.

* **Independência de frameworks**: A lógica de negócios não depende de bibliotecas externas
* **Testabilidade**: Cada camada pode ser testada independentemente
* **Independência de UI**: A lógica pode funcionar sem interface web
* **Independência de banco de dados**: Regras de negócio não conhecem detalhes de persistência

### SOLID

SOLID é um conjunto de princípios de design orientado a objetos:

* **S**ingle Responsibility: Cada módulo tem uma única responsabilidade
* **O**pen/Closed: Aberto para extensão, fechado para modificação
* **L**iskov Substitution: Interfaces bem definidas e substituíveis
* **I**nterface Segregation: Interfaces pequenas e específicas
* **D**ependency Inversion: Dependências apontam para abstrações

## Estrutura do Projeto

Este é apenas um modelo esperado da estrutura do projeto.

```
├── cmd/
│   ├── api/
│   │   └── main.go                      # Entry point da API REST
│   ├── worker/
│   │   └── main.go                      # Entry point para workers/jobs
│   └── migration/
│       └── main.go                      # Entry point para migrações
├── config/
│   ├── config.yml                       # Configuração padrão
│   ├── config.dev.yml                   # Configuração de desenvolvimento
│   ├── config.prod.yml                  # Configuração de produção
│   └── config.test.yml                  # Configuração para testes
├── internal/                            # Código privado da aplicação
│   ├── adapter/                         # Adaptadores para serviços externos
│   │   ├── cache/
│   │   │   ├── interfaces.go            # Interface do adaptador de cache
│   │   │   ├── redis.go                 # Implementação Redis
│   │   │   └── memory.go                # Implementação em memória
│   │   ├── payment/
│   │   │   ├── interfaces.go            # Interface do gateway de pagamento
│   │   │   ├── stripe.go                # Adaptador para Stripe
│   │   │   └── paypal.go                # Adaptador para PayPal
│   │   ├── notification/
│   │   │   ├── interfaces.go            # Interfaces de notificação
│   │   │   ├── email.go                 # Adaptador para envio de emails
│   │   │   ├── sms.go                   # Adaptador para envio de SMS
│   │   │   └── push.go                  # Adaptador para push notifications
│   │   └── storage/
│   │       ├── interfaces.go            # Interface de storage
│   │       ├── s3.go                    # Adaptador S3
│   │       └── local.go                 # Adaptador local
│   ├── config/                          # Estruturas de configuração
│   │   ├── application.go               # Config gerais da aplicação
│   │   ├── config.go                    # Config principal
│   │   ├── database.go                  # Config de banco de dados
│   │   ├── interfaces.go                # Interfaces de configuração
│   │   ├── oidc.go                      # Config de OIDC/OAuth2
│   │   ├── redis.go                     # Config do Redis
│   │   ├── services.go                  # Config de serviços externos
│   │   └── web.go                       # Config do servidor web
│   ├── dto/                             # Data Transfer Objects
│   │   ├── request/
│   │   │   ├── user.go                  # DTOs de requisição de usuário
│   │   │   ├── product.go               # DTOs de requisição de produto
│   │   │   ├── order.go                 # DTOs de requisição de pedido
│   │   │   ├── auth.go                  # DTOs de requisição de autenticação
│   │   │   └── pagination.go            # DTO de paginação para requests
│   │   └── response/
│   │       ├── user.go                  # DTOs de resposta de usuário
│   │       ├── product.go               # DTOs de resposta de produto
│   │       ├── order.go                 # DTOs de resposta de pedido
│   │       ├── auth.go                  # DTOs de resposta de autenticação
│   │       ├── error.go                 # DTO de resposta de erro
│   │       └── pagination.go            # DTO genérico de paginação
│   ├── entity/                          # Entidades de domínio
│   │   ├── user.go                      # Entidade usuário
│   │   ├── product.go                   # Entidade produto
│   │   ├── order.go                     # Entidade pedido
│   │   ├── payment.go                   # Entidade pagamento
│   │   ├── category.go                  # Entidade categoria
│   │   └── value_objects.go             # Value objects (Money, Address, etc)
│   ├── fx/                              # Módulos de injeção de dependência
│   │   └── module/
│   │       ├── core.go                  # Módulo core (config, logger, etc)
│   │       ├── database.go              # Módulo de banco de dados
│   │       ├── handler.go               # Módulo de handlers HTTP
│   │       ├── repository.go            # Módulo de repositories
│   │       ├── usecase.go               # Módulo de use cases
│   │       ├── service.go               # Módulo de services
│   │       ├── adapter.go               # Módulo de adaptadores
│   │       ├── infrastructure.go        # Módulo de infraestrutura
│   │       ├── middleware.go            # Módulo de middlewares
│   │       └── router.go                # Módulo de roteamento
│   ├── http/
│   │   ├── handler/                     # Handlers HTTP
│   │   │   ├── user.go                  # Handler de usuários
│   │   │   ├── product.go               # Handler de produtos
│   │   │   ├── order.go                 # Handler de pedidos
│   │   │   ├── auth.go                  # Handler de autenticação
│   │   │   ├── payment.go               # Handler de pagamentos
│   │   │   └── health.go                # Handler de health check
│   │   ├── middleware/                  # Middlewares customizados
│   │   │   ├── auth.go                  # Middleware de autenticação
│   │   │   ├── cors.go                  # Middleware de CORS
│   │   │   ├── oidc.go                  # Middleware de OIDC
│   │   │   ├── rate_limit.go            # Middleware de rate limiting
│   │   │   ├── logger.go                # Middleware de logging
│   │   │   ├── recovery.go              # Middleware de recovery
│   │   │   └── validation.go            # Middleware de validação
│   │   └── presenter/                   # Presenters para formatar respostas
│   │       ├── user.go                  # Presenter de usuário
│   │       ├── product.go               # Presenter de produto
│   │       ├── order.go                 # Presenter de pedido
│   │       └── payment.go               # Presenter de pagamento
│   ├── infrastructure/                  # Camada de infraestrutura
│   │   ├── auth/
│   │   │   ├── oidc.go                  # Implementação OIDC
│   │   │   ├── oauth2.go                # Implementação OAuth2
│   │   │   ├── jwt.go                   # Implementação JWT
│   │   │   └── provider.go              # Interface de autenticação
│   │   ├── cache/
│   │   │   ├── redis.go                 # Implementação Redis
│   │   │   ├── memory.go                # Cache em memória
│   │   │   └── distributed.go           # Cache distribuído
│   │   ├── database/
│   │   │   ├── connection.go            # Gerenciamento de conexões
│   │   │   ├── transaction.go           # Gerenciamento de transações
│   │   │   ├── migrations.go            # Sistema de migrações
│   │   │   └── seeder.go                # Seeders para desenvolvimento
│   │   ├── model/                       # Modelos GORM
│   │   │   ├── user.go                  # Modelo de usuário
│   │   │   ├── product.go               # Modelo de produto
│   │   │   ├── order.go                 # Modelo de pedido
│   │   │   ├── payment.go               # Modelo de pagamento
│   │   │   ├── category.go              # Modelo de categoria
│   │   │   └── base.go                  # Modelo base com campos comuns
│   │   ├── messaging/
│   │   │   ├── interfaces.go            # Interfaces de messaging
│   │   │   ├── rabbitmq.go              # Implementação RabbitMQ
│   │   │   ├── kafka.go                 # Implementação Kafka
│   │   │   └── sqs.go                   # Implementação AWS SQS
│   │   └── storage/
│   │       ├── s3.go                    # Implementação S3
│   │       ├── gcs.go                   # Implementação Google Cloud Storage
│   │       └── local.go                 # Storage local
│   ├── repository/                      # Implementações de repositórios
│   │   ├── user_repository.go           # Repository de usuários
│   │   ├── product_repository.go        # Repository de produtos
│   │   ├── order_repository.go          # Repository de pedidos
│   │   ├── payment_repository.go        # Repository de pagamentos
│   │   ├── category_repository.go       # Repository de categorias
│   │   └── interfaces.go                # Interfaces dos repositories (se não estiverem em usecase)
│   ├── service/                         # Serviços de domínio
│   │   ├── auth/
│   │   │   └── auth_service.go          # Serviço de autenticação
│   │   ├── email/
│   │   │   ├── email_service.go         # Serviço de email
│   │   │   └── templates.go             # Templates de email
│   │   ├── notification/
│   │   │   └── notification_service.go  # Serviço de notificação
│   │   ├── payment/
│   │   │   └── payment_service.go       # Serviço de pagamento
│   │   ├── storage/
│   │   │   └── storage_service.go       # Serviço de armazenamento
│   │   └── interfaces.go                # Interfaces dos serviços
│   └── usecase/                         # Casos de uso (regras de negócio)
│       ├── user/
│       │   ├── create_user.go           # Use case criar usuário
│       │   ├── get_user.go              # Use case buscar usuário
│       │   ├── update_user.go           # Use case atualizar usuário
│       │   ├── delete_user.go           # Use case deletar usuário
│       │   ├── list_users.go            # Use case listar usuários
│       │   └── change_password.go       # Use case alterar senha
│       ├── product/
│       │   ├── create_product.go        # Use case criar produto
│       │   ├── update_product.go        # Use case atualizar produto
│       │   ├── delete_product.go        # Use case deletar produto
│       │   ├── get_product.go           # Use case buscar produto
│       │   ├── list_products.go         # Use case listar produtos
│       │   └── update_stock.go          # Use case atualizar estoque
│       ├── order/
│       │   ├── create_order.go          # Use case criar pedido
│       │   ├── cancel_order.go          # Use case cancelar pedido
│       │   ├── list_orders.go           # Use case listar pedidos
│       │   ├── get_order.go             # Use case buscar pedido
│       │   └── process_payment.go       # Use case processar pagamento
│       └── interfaces.go                # Interfaces dos use cases e repositories
├── pkg/                                 # Pacotes reutilizáveis
│   ├── auth/
│   │   ├── jwt.go                       # Utilitários JWT
│   │   ├── oidc_client.go               # Cliente OIDC genérico
│   │   ├── password.go                  # Hash e validação de senha
│   │   └── token.go                     # Gerenciamento de tokens
│   ├── database/
│   │   ├── gorm.go                      # Configuração GORM
│   │   ├── postgres.go                  # Driver PostgreSQL
│   │   ├── mysql.go                     # Driver MySQL
│   │   └── sqlite.go                    # Driver SQLite (para testes)
│   ├── errors/
│   │   ├── errors.go                    # Erros customizados
│   │   ├── http.go                      # Erros HTTP
│   │   └── domain.go                    # Erros de domínio
│   ├── http/
│   │   ├── client.go                    # Cliente HTTP customizado
│   │   ├── gin/
│   │   │   ├── adapter.go               # Adaptador para Gin
│   │   │   ├── gin.go                   # Configuração Gin
│   │   │   ├── helpers.go               # Helpers HTTP
│   │   │   └── validators.go            # Validadores customizados
│   │   └── fiber/
│   │       ├── adapter.go               # Adaptador para Fiber
│   │       └── fiber.go                 # Configuração Fiber
│   ├── log/
│   │   ├── log.go                       # Interface de logging
│   │   ├── zerolog/
│   │   │   └── zerolog.go               # Implementação com zerolog
│   │   └── zap/
│   │       └── zap.go                   # Implementação com zap
│   ├── validator/
│   │   ├── validator.go                 # Validações customizadas
│   │   ├── rules.go                     # Regras de validação
│   │   └── messages.go                  # Mensagens de erro customizadas
│   └── utils/
│       ├── env.go                       # Utilitários de environment
│       ├── response.go                  # Helpers de resposta
│       ├── converter.go                 # Conversores gerais
│       ├── crypto.go                    # Utilitários de criptografia
│       └── string.go                    # Manipulação de strings
├── migrations/                          # Migrações de banco de dados
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_products_table.up.sql
│   ├── 000002_create_products_table.down.sql
│   ├── 000003_create_orders_table.up.sql
│   ├── 000003_create_orders_table.down.sql
│   ├── 000004_create_payments_table.up.sql
│   └── 000004_create_payments_table.down.sql
├── seeders/                            # População do banco de dados
│   ├── 000001_users.sql
│   ├── 000002_products.sql
│   ├── 000003_orders.sql
│   └── 000004_payments.sql
├── scripts/                             # Scripts auxiliares
│   ├── setup.sh                         # Script de setup inicial
│   ├── test.sh                          # Script de testes
│   ├── build.sh                         # Script de build
│   ├── deploy.sh                        # Script de deploy
│   └── generate-mocks.sh                # Script para gerar mocks
├── test/                                # Testes de integração/e2e
│   ├── integration/
│   │   ├── user_test.go                 # Testes de integração de usuário
│   │   ├── product_test.go              # Testes de integração de produto
│   │   ├── order_test.go                # Testes de integração de pedido
│   │   └── payment_test.go              # Testes de integração de pagamento
│   ├── e2e/
│   │   ├── api_test.go                  # Testes end-to-end da API
│   │   ├── auth_flow_test.go            # Teste do fluxo de autenticação
│   │   └── order_flow_test.go           # Teste do fluxo de pedido
│   └── fixtures/
│       ├── users.json                   # Fixtures de usuários
│       ├── products.json                # Fixtures de produtos
│       └── orders.json                  # Fixtures de pedidos
├── docs/                                # Documentação
│   ├── swagger.json                     # Documentação Swagger gerada
│   ├── swagger.yaml                     # Documentação Swagger fonte
│   ├── api.md                           # Documentação manual da API
│   ├── architecture.md                  # Documentação da arquitetura
│   └── deployment.md                    # Guia de deployment
├── deployments/                         # Configurações de deployment
│   ├── docker/
│   │   ├── Dockerfile                   # Dockerfile da aplicação
│   │   └── Dockerfile.dev               # Dockerfile para desenvolvimento
│   ├── kubernetes/
│   │   ├── deployment.yaml              # Deployment do Kubernetes
│   │   ├── service.yaml                 # Service do Kubernetes
│   │   └── configmap.yaml               # ConfigMap do Kubernetes
│   └── terraform/
│       ├── main.tf                      # Configuração principal Terraform
│       └── variables.tf                 # Variáveis Terraform
├── .gitignore                           # Arquivos ignorados pelo Git
├── .dockerignore                        # Arquivos ignorados pelo Docker Build
├── .env.example                         # Exemplo de variáveis de ambiente
├── docker-compose.yml                   # Configuração Docker Compose
├── docker-compose.dev.yml               # Override para desenvolvimento
├── docker-compose.test.yml              # Override para testes
├── Makefile                             # Comandos make
├── README.md                            # Documentação do projeto
├── go.mod                               # Dependências Go
└── go.sum                               # Lock de dependências
```

## Camadas da Arquitetura

### 1. **Entities** (`internal/entity/`)

Objetos de negócio centrais, sem dependências externas. Contêm as regras de negócio fundamentais.

```go
// internal/entity/user.go
package entity

import (
    "errors"
    "regexp"
    "time"
)

type User struct {
    ID        uint
    Name      string
    Email     string
    Password  string 
    Role      UserRole
    Active    bool
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UserRole string

const (
    RoleAdmin    UserRole  = "admin"
    RoleUser     UserRole  = "user"
    RoleModerator UserRole = "moderator"
)

func NewUser(name, email, hashedPassword string) (*User, error) {
    user := &User{
        Name:      name,
        Email:     email,
        Password:  hashedPassword,
        Role:      RoleUser,
        Active:    true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := user.Validate(); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (u *User) Validate() error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    
    if !u.IsValidEmail() {
        return errors.New("invalid email format")
    }
    
    if u.Password == "" {
        return errors.New("password is required")
    }
    
    return nil
}

func (u *User) IsValidEmail() bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(u.Email)
}
```

### 2. **Services** (`internal/service/`)

Services encapsulam lógica de negócio complexa que pode envolver múltiplos use cases ou operações externas.

```go
// internal/service/interfaces.go
package service

import (
    "context"
    "meu-projeto/internal/entity"
)

// EmailService define operações de envio de email
type EmailService interface {
    SendWelcomeEmail(ctx context.Context, user *entity.User) error
    SendPasswordResetEmail(ctx context.Context, user *entity.User, token string) error
    SendOrderConfirmation(ctx context.Context, order *entity.Order) error
}

// NotificationService define operações de notificação
type NotificationService interface {
    NotifyUserCreated(ctx context.Context, user *entity.User) error
    NotifyOrderCreated(ctx context.Context, order *entity.Order) error
    NotifyPaymentProcessed(ctx context.Context, payment *entity.Payment) error
}

// PaymentService define operações de pagamento
type PaymentService interface {
    ProcessPayment(ctx context.Context, order *entity.Order, paymentMethod string) (*entity.Payment, error)
    RefundPayment(ctx context.Context, paymentID uint, amount float64) error
    GetPaymentStatus(ctx context.Context, paymentID uint) (string, error)
}

// AuthService define operações de autenticação
type AuthService interface {
    GenerateToken(ctx context.Context, user *entity.User) (string, error)
    ValidateToken(ctx context.Context, token string) (*entity.User, error)
    RefreshToken(ctx context.Context, refreshToken string) (string, error)
    GeneratePasswordResetToken(ctx context.Context, email string) (string, error)
}
```

### 3. **Use Cases** (`internal/usecase/`)

Podem definir regras de negócio complexas da aplicação.

Define interfaces para repositories e implementa a lógica de negócio e/ou implementam um caso de uso especifico para uma ação.

#### Interfaces

```go
// internal/usecase/interfaces.go
package usecase

import (
    "context"
    "time"
    "meu-projeto/internal/entity"
)

// Repository interfaces
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id uint) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)
    Search(ctx context.Context, query string) ([]*entity.User, error)
}

type ProductRepository interface {
    Create(ctx context.Context, product *entity.Product) error
    GetByID(ctx context.Context, id uint) (*entity.Product, error)
    GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
    Update(ctx context.Context, product *entity.Product) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, filters ProductFilters, offset, limit int) ([]*entity.Product, int64, error)
    UpdateStock(ctx context.Context, id uint, quantity int) error
    GetByCategory(ctx context.Context, categoryID uint) ([]*entity.Product, error)
}

// Filtros
type ProductFilters struct {
    Name       string
    CategoryID *uint
    MinPrice   *float64
    MaxPrice   *float64
    InStock    *bool
}
```

#### Implementação de Use Case

```go
// internal/usecase/user/create_user.go
package user

import (
    "context"
    "errors"
    "meu-projeto/internal/entity"
    "meu-projeto/internal/usecase"
    "meu-projeto/pkg/auth"
    "meu-projeto/pkg/log"
)

type CreateUserInput struct {
    Name     string
    Email    string
    Password string
    Role     string
}

type CreateUser struct {
    userRepo usecase.UserRepository
    log      log.Logger
}

func NewCreateUser(userRepo usecase.UserRepository, log log.Logger) *CreateUser {
    return &CreateUser{
        userRepo: userRepo,
        log:      log,
    }
}

func (uc *CreateUser) Execute(ctx context.Context, input CreateUserInput) (*entity.User, error) {
    // Validação de entrada
    if err := uc.validate(input); err != nil {
        uc.log.WithError(err).Error("validation failed")
        return nil, err
    }
    
    // Verificar se email já existe
    existing, _ := uc.userRepo.GetByEmail(ctx, input.Email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }
    
    // Hash da senha
    hashedPassword, err := auth.HashPassword(input.Password)
    if err != nil {
        uc.log.WithError(err).Error("failed to hash password")
        return nil, errors.New("failed to process password")
    }
    
    // Criar entidade
    user, err := entity.NewUser(input.Name, input.Email, hashedPassword)
    if err != nil {
        return nil, err
    }
    
    // Definir role se especificado
    if input.Role != "" {
        if err := user.ChangeRole(entity.UserRole(input.Role)); err != nil {
            return nil, err
        }
    }
    
    // Persistir
    if err := uc.userRepo.Create(ctx, user); err != nil {
        uc.log.WithError(err).Error("failed to create user")
        return nil, errors.New("failed to create user")
    }
    
    uc.log.WithField("user_id", user.ID).Info("user created successfully")
    return user, nil
}

func (uc *CreateUser) validate(input CreateUserInput) error {
    if input.Name == "" {
        return errors.New("name is required")
    }
    if input.Email == "" {
        return errors.New("email is required")
    }
    if len(input.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}
```


3. \

### 3. **Repositories** (`internal/repository/`)

Implementações concretas de persistência.

```go
// internal/repository/user_repository.go
package repository

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "meu-projeto/internal/entity"
    "meu-projeto/internal/infrastructure/model"
    "meu-projeto/internal/usecase"
)

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) usecase.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    userModel := r.toModel(user)
    
    if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
        return err
    }
    
    user.ID = userModel.ID
    return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
    var userModel model.User
    
    if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    
    return r.toEntity(&userModel), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    var userModel model.User
    
    if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    
    return r.toEntity(&userModel), nil
}

// Conversores
func (r *userRepository) toModel(user *entity.User) model.User {
    return model.User{
        BaseModel: model.BaseModel{
            ID:        user.ID,
            CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
        },
        Name:     user.Name,
        Email:    user.Email,
        Password: user.Password,
        Role:     string(user.Role),
        Active:   user.Active,
    }
}

func (r *userRepository) toEntity(userModel *model.User) *entity.User {
    return &entity.User{
        ID:        userModel.ID,
        Name:      userModel.Name,
        Email:     userModel.Email,
        Password:  userModel.Password,
        Role:      entity.UserRole(userModel.Role),
        Active:    userModel.Active,
        CreatedAt: userModel.CreatedAt,
        UpdatedAt: userModel.UpdatedAt,
    }
}
```

### 4. **Handlers** (`internal/http/handler/`)

Controllers HTTP que recebem requisições e delegam para use cases.

```go
// internal/http/handler/user.go
package handler

import (
    "net/http"
    "strconv"
    "meu-projeto/internal/dto/request"
    "meu-projeto/internal/dto/response"
    "meu-projeto/internal/usecase/user"
    "meu-projeto/pkg/log"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    createUser *user.CreateUser
    getUser    *user.GetUser
    listUsers  *user.ListUsers
    log        log.Logger
}

func NewUserHandler(
    createUser *user.CreateUser,
    getUser *user.GetUser,
    listUsers *user.ListUsers,
    log log.Logger,
) *UserHandler {
    return &UserHandler{
        createUser: createUser,
        getUser:    getUser,
        listUsers:  listUsers,
        log:        log,
    }
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body request.CreateUserRequest true "Create user"
// @Success 201 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req request.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.ErrorResponse{
            Error: "Invalid request body",
        })
        return
    }
    
    input := user.CreateUserInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        Role:     req.Role,
    }
    
    newUser, err := h.createUser.Execute(c.Request.Context(), input)
    if err != nil {
        h.log.WithError(err).Error("failed to create user")
        c.JSON(http.StatusInternalServerError, response.ErrorResponse{
            Error: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, response.ToUserResponse(newUser))
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get details of user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.UserResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.ErrorResponse{
            Error: "Invalid user ID",
        })
        return
    }
    
    user, err := h.getUser.Execute(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, response.ErrorResponse{
            Error: "User not found",
        })
        return
    }
    
    c.JSON(http.StatusOK, response.ToUserResponse(user))
}

// RegisterRoutes registra as rotas do handler
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
    users := router.Group("/users")
    {
        users.POST("", h.CreateUser)
        users.GET("/:id", h.GetUser)
        users.GET("", h.ListUsers)
    }
}
```

## Configuração Inicial

### 1. Clone o template

```bash
git clone https://github.com/templates/golang.git meu-projeto
cd meu-projeto
```

### 2. Instale as dependências

```bash
go mod download
```

### 3. Configure o ambiente

```bash
cp config/config.yml config/config.local.yml
# Edite config/config.local.yml com suas configurações
```

### 4. Execute as migrações

```bash
make migrate-up
```

## Guia de Desenvolvimento

### Criando um Novo Módulo


1. **Defina a entidade** em `internal/entity/`
2. **Crie os DTOs** em `internal/dto/`
3. **Defina as interfaces** em `internal/usecase/interfaces.go`
4. **Implemente o use case** em `internal/usecase/`
5. **Crie o modelo GORM** em `internal/infrastructure/model/`
6. **Implemente o repository** em `internal/repository/`
7. **Adicione o handler** em `internal/http/handler/`
8. **Configure a injeção de dependências** em `internal/fx/module/`

### Exemplo de implementação de entidade

Entidades são objetos de domínio puros que contêm regras de negócio fundamentais:

```go
// internal/entity/product.go
package entity

import (
    "errors"
    "time"
)

type Product struct {
    ID          uint
    Name        string
    Description string
    SKU         string
    Price       Money
    Stock       int
    CategoryID  uint
    Active      bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Money struct {
    Amount   float64
    Currency string
}

func NewProduct(name, description, sku string, price float64, stock int) (*Product, error) {
    product := &Product{
        Name:        name,
        Description: description,
        SKU:         sku,
        Price:       Money{Amount: price, Currency: "BRL"},
        Stock:       stock,
        Active:      true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    if err := product.Validate(); err != nil {
        return nil, err
    }
    
    return product, nil
}

func (p *Product) Validate() error {
    if p.Name == "" {
        return errors.New("product name is required")
    }
    
    if p.SKU == "" {
        return errors.New("SKU is required")
    }
    
    if p.Price.Amount < 0 {
        return errors.New("price cannot be negative")
    }
    
    if p.Stock < 0 {
        return errors.New("stock cannot be negative")
    }
    
    return nil
}

func (p *Product) UpdateStock(quantity int) error {
    if p.Stock+quantity < 0 {
        return errors.New("insufficient stock")
    }
    p.Stock += quantity
    p.UpdatedAt = time.Now()
    return nil
}

func (p *Product) IsAvailable() bool {
    return p.Active && p.Stock > 0
}
```

### Exemplo de implementação de Use Cases

Use cases contêm a lógica de negócio da aplicação:

```go
// internal/usecase/product/create_product.go
package product

import (
    "context"
    "errors"
    "meu-projeto/internal/entity"
    "meu-projeto/internal/usecase"
    "meu-projeto/pkg/log"
)

type CreateProductInput struct {
    Name        string
    Description string
    SKU         string
    Price       float64
    Stock       int
    CategoryID  uint
}

type CreateProduct struct {
    productRepo usecase.ProductRepository
    log         log.Logger
}

func NewCreateProduct(productRepo usecase.ProductRepository, log log.Logger) *CreateProduct {
    return &CreateProduct{
        productRepo: productRepo,
        log:         log,
    }
}

func (uc *CreateProduct) Execute(ctx context.Context, input CreateProductInput) (*entity.Product, error) {
    // Verificar se SKU já existe
    existing, _ := uc.productRepo.GetBySKU(ctx, input.SKU)
    if existing != nil {
        return nil, errors.New("SKU already exists")
    }
    
    // Criar produto
    product, err := entity.NewProduct(
        input.Name,
        input.Description,
        input.SKU,
        input.Price,
        input.Stock,
    )
    if err != nil {
        return nil, err
    }
    
    product.CategoryID = input.CategoryID
    
    // Persistir
    if err := uc.productRepo.Create(ctx, product); err != nil {
        uc.log.WithError(err).Error("failed to create product")
        return nil, errors.New("failed to create product")
    }
    
    uc.log.WithField("product_id", product.ID).Info("product created successfully")
    return product, nil
}
```

### Exemplo de implementação de serviço

Serviçes implementam lógica de negócio complexa que pode envolver múltiplos casos de uso:

```go
// internal/service/auth/auth_service.go
package auth

import (
    "context"
    "errors"
    "time"
    "meu-projeto/internal/entity"
    "meu-projeto/internal/service"
    "meu-projeto/pkg/log"
    "github.com/golang-jwt/jwt/v5"
)

type authService struct {
    jwtSecret          string
    jwtExpiration      time.Duration
    refreshExpiration  time.Duration
    log                log.Logger
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func NewAuthService(
    jwtSecret string,
    jwtExpiration time.Duration,
    refreshExpiration time.Duration,
    log log.Logger,
) service.AuthService {
    return &authService{
        jwtSecret:         jwtSecret,
        jwtExpiration:     jwtExpiration,
        refreshExpiration: refreshExpiration,
        log:               log,
    }
}

func (s *authService) GenerateToken(ctx context.Context, user *entity.User) (string, error) {
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   string(user.Role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        s.log.WithError(err).Error("failed to generate token")
        return "", err
    }
    
    s.log.WithField("user_id", user.ID).Info("token generated successfully")
    return tokenString, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*entity.User, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })
    
    if err != nil {
        s.log.WithError(err).Error("failed to parse token")
        return nil, errors.New("invalid token")
    }
    
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    // Aqui você normalmente buscaria o usuário completo do banco
    // Este é um exemplo simplificado
    user := &entity.User{
        ID:    claims.UserID,
        Email: claims.Email,
        Role:  entity.UserRole(claims.Role),
    }
    
    return user, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
    // Validar refresh token
    user, err := s.ValidateToken(ctx, refreshToken)
    if err != nil {
        return "", err
    }
    
    // Gerar novo access token
    return s.GenerateToken(ctx, user)
}

func (s *authService) GeneratePasswordResetToken(ctx context.Context, email string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "type":  "password_reset",
        "exp":   time.Now().Add(1 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}
```

### Exemplo de implementação de repositorio

Repositories fazem a ponte entre entities e modelos de banco de dados:

```go
// internal/repository/product_repository.go
package repository

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "meu-projeto/internal/entity"
    "meu-projeto/internal/infrastructure/model"
    "meu-projeto/internal/usecase"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) usecase.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
    productModel := r.toModel(product)
    
    if err := r.db.WithContext(ctx).Create(&productModel).Error; err != nil {
        return err
    }
    
    product.ID = productModel.ID
    return nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
    var productModel model.Product
    
    if err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&productModel).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("product not found")
        }
        return nil, err
    }
    
    return r.toEntity(&productModel), nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var product model.Product
        
        // Lock para evitar condições de corrida
        if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, id).Error; err != nil {
            return err
        }
        
        newStock := product.Stock + quantity
        if newStock < 0 {
            return errors.New("insufficient stock")
        }
        
        return tx.Model(&product).Update("stock", newStock).Error
    })
}

// Conversores
func (r *productRepository) toModel(product *entity.Product) model.Product {
    return model.Product{
        BaseModel: model.BaseModel{
            ID:        product.ID,
            CreatedAt: product.CreatedAt,
            UpdatedAt: product.UpdatedAt,
        },
        Name:        product.Name,
        Description: product.Description,
        SKU:         product.SKU,
        Price:       product.Price.Amount,
        Stock:       product.Stock,
        CategoryID:  product.CategoryID,
        Active:      product.Active,
    }
}

func (r *productRepository) toEntity(productModel *model.Product) *entity.Product {
    return &entity.Product{
        ID:          productModel.ID,
        Name:        productModel.Name,
        Description: productModel.Description,
        SKU:         productModel.SKU,
        Price:       entity.Money{Amount: productModel.Price, Currency: "BRL"},
        Stock:       productModel.Stock,
        CategoryID:  productModel.CategoryID,
        Active:      productModel.Active,
        CreatedAt:   productModel.CreatedAt,
        UpdatedAt:   productModel.UpdatedAt,
    }
}
```

### Exemplo de implementação manipulador HTTP

Handlers processam requisições HTTP e delegam para use cases:

```go
// internal/http/handler/product.go
package handler

import (
    "net/http"
    "strconv"
    "meu-projeto/internal/dto/request"
    "meu-projeto/internal/dto/response"
    "meu-projeto/internal/usecase/product"
    "meu-projeto/pkg/log"
    "github.com/gin-gonic/gin"
)

type ProductHandler struct {
    createProduct *product.CreateProduct
    listProducts  *product.ListProducts
    log           log.Logger
}

func NewProductHandler(
    createProduct *product.CreateProduct,
    listProducts *product.ListProducts,
    log log.Logger,
) *ProductHandler {
    return &ProductHandler{
        createProduct: createProduct,
        listProducts:  listProducts,
        log:           log,
    }
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Param product body request.CreateProductRequest true "Create product"
// @Success 201 {object} response.ProductResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
    var req request.CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.ErrorResponse{
            Error: "Invalid request body",
        })
        return
    }
    
    input := product.CreateProductInput{
        Name:        req.Name,
        Description: req.Description,
        SKU:         req.SKU,
        Price:       req.Price,
        Stock:       req.Stock,
        CategoryID:  req.CategoryID,
    }
    
    newProduct, err := h.createProduct.Execute(c.Request.Context(), input)
    if err != nil {
        h.log.WithError(err).Error("failed to create product")
        c.JSON(http.StatusInternalServerError, response.ErrorResponse{
            Error: err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, response.ToProductResponse(newProduct))
}

func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup) {
    products := router.Group("/products")
    {
        products.POST("", h.CreateProduct)
        products.GET("", h.ListProducts)
        products.GET("/:id", h.GetProduct)
        products.PUT("/:id", h.UpdateProduct)
        products.DELETE("/:id", h.DeleteProduct)
    }
}
```

### Configurando Injeção de Dependências

Configure todos os módulos usando Uber FX:

```go
// internal/fx/module/repository.go
package module

import (
    "context"
    "gorm.io/gorm"
    "meu-projeto/internal/repository"
    "meu-projeto/internal/usecase"
    "meu-projeto/pkg/log"
    "go.uber.org/fx"
)

type RepositoryContainer struct {
    fx.Out
    UserRepository    usecase.UserRepository
    ProductRepository usecase.ProductRepository
    OrderRepository   usecase.OrderRepository
}

func NewRepositories(db *gorm.DB) (RepositoryContainer, error) {
    return RepositoryContainer{
        UserRepository:    repository.NewUserRepository(db),
        ProductRepository: repository.NewProductRepository(db),
        OrderRepository:   repository.NewOrderRepository(db),
    }, nil
}

func Repository() fx.Option {
    return fx.Module("repository",
        fx.Provide(NewRepositories),
        fx.Invoke(func(lc fx.Lifecycle, log log.Logger) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    log.Info("Repositories initialized")
                    return nil
                },
            })
        }),
    )
}
```

```go
// internal/fx/module/usecase.go
package module

import (
    "context"
    "meu-projeto/internal/usecase"
    "meu-projeto/internal/usecase/user"
    "meu-projeto/internal/usecase/product"
    "meu-projeto/internal/usecase/order"
    "meu-projeto/pkg/log"
    "go.uber.org/fx"
)

type UseCaseContainer struct {
    fx.Out
    // User use cases
    CreateUser *user.CreateUser
    GetUser    *user.GetUser
    ListUsers  *user.ListUsers
    
    // Product use cases
    CreateProduct *product.CreateProduct
    ListProducts  *product.ListProducts
    
    // Order use cases
    CreateOrder *order.CreateOrder
}

type UseCaseParams struct {
    fx.In
    UserRepo    usecase.UserRepository
    ProductRepo usecase.ProductRepository
    OrderRepo   usecase.OrderRepository
    Log         log.Logger
}

func NewUseCases(params UseCaseParams) (UseCaseContainer, error) {
    return UseCaseContainer{
        // User use cases
        CreateUser: user.NewCreateUser(params.UserRepo, params.Log),
        GetUser:    user.NewGetUser(params.UserRepo, params.Log),
        ListUsers:  user.NewListUsers(params.UserRepo, params.Log),
        
        // Product use cases
        CreateProduct: product.NewCreateProduct(params.ProductRepo, params.Log),
        ListProducts:  product.NewListProducts(params.ProductRepo, params.Log),
        
        // Order use cases
        CreateOrder: order.NewCreateOrder(params.OrderRepo, params.ProductRepo, params.Log),
    }, nil
}

func UseCase() fx.Option {
    return fx.Module("usecase",
        fx.Provide(NewUseCases),
    )
}
```

```go
// internal/fx/module/handler.go
package module

import (
    "context"
    "meu-projeto/internal/http/handler"
    "meu-projeto/internal/usecase/user"
    "meu-projeto/internal/usecase/product"
    "meu-projeto/pkg/log"
    "github.com/gin-gonic/gin"
    "go.uber.org/fx"
)

type HandlerContainer struct {
    fx.Out
    UserHandler    *handler.UserHandler
    ProductHandler *handler.ProductHandler
}

type HandlerParams struct {
    fx.In
    // User use cases
    CreateUser *user.CreateUser
    GetUser    *user.GetUser
    ListUsers  *user.ListUsers
    
    // Product use cases
    CreateProduct *product.CreateProduct
    ListProducts  *product.ListProducts
    
    Log log.Logger
}

func NewHandlers(params HandlerParams) (HandlerContainer, error) {
    return HandlerContainer{
        UserHandler: handler.NewUserHandler(
            params.CreateUser,
            params.GetUser,
            params.ListUsers,
            params.Log,
        ),
        ProductHandler: handler.NewProductHandler(
            params.CreateProduct,
            params.ListProducts,
            params.Log,
        ),
    }, nil
}

func Handler() fx.Option {
    return fx.Module("handler",
        fx.Provide(NewHandlers),
        fx.Invoke(func(
            router *gin.Engine,
            userHandler *handler.UserHandler,
            productHandler *handler.ProductHandler,
        ) {
            api := router.Group("/api/v1")
            
            userHandler.RegisterRoutes(api)
            productHandler.RegisterRoutes(api)
        }),
    )
}
```

## Banco de Dados

### Usando GORM

#### Modelo Base

```go
// internal/infrastructure/model/base.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
```

#### Modelos de Domínio

```go
// internal/infrastructure/model/user.go
package model

import (
    "gorm.io/gorm"
)

type User struct {
    BaseModel
    Name     string    `gorm:"size:255;not null" json:"name"`
    Email    string    `gorm:"uniqueIndex;not null" json:"email"`
    Password string    `gorm:"not null" json:"-"`
    Role     string    `gorm:"size:50;default:'user'" json:"role"`
    Active   bool      `gorm:"default:true" json:"active"`
    
    // Relacionamentos
    Orders   []Order   `gorm:"foreignKey:UserID" json:"orders,omitempty"`
    Profile  *Profile  `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
    // Validações e hooks
    if u.Email == "" {
        return gorm.ErrInvalidData
    }
    return nil
}
```

#### Configuração GORM

```go
// internal/infrastructure/database/connection.go
package database

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

type DatabaseConfig struct {
    Host        string
    Port        int
    User        string
    Password    string
    DBName      string
    SSLMode     string
    TablePrefix string
    AutoMigrate bool
}

func NewGormDB(config DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
        config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   config.TablePrefix,
            SingularTable: true,
        },
    })
    
    if err != nil {
        return nil, err
    }
    
    // Auto migrate
    if config.AutoMigrate {
        err = db.AutoMigrate(
            &model.User{},
            &model.Product{},
            &model.Order{},
            &model.OrderItem{},
            &model.Payment{},
        )
    }
    
    return db, err
}
```

### Usando golang-migrate


1. **Instale a ferramenta**:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```


2. **Crie uma migração**:

```bash
migrate create -ext sql -dir migrations -seq create_users_table
```


3. **Escreva a migração**:

```sql
-- migrations/000001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- migrations/000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```


4. **Execute as migrações**:

```bash
migrate -path migrations -database "postgresql://user:pass@localhost/dbname?sslmode=disable" up
```

## Documentação da API

### Configurando Swagger


1. **Instale swag**:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```


2. **Adicione anotações no main.go**:

```go
// @title           Clean Architecture API
// @version         1.0
// @description     API Server with Clean Architecture

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
    // ...
}
```


3. **Gere a documentação**:

```bash
swag init -g cmd/api/main.go --output docs
```

## Arquivo de Configuração e Variáveis de Ambiente

Use o pacote viper para ler arquivos de configuração e/ou variáveis de ambiente.

### Estrutura do Arquivo `config.yml`

```yaml
app:
  name: "Clean Architecture API"                                 # Nome da aplicação
  description: "API com arquitetura limpa"                       # Descrição da aplicação
  version: "1.0.0"                                               # Versão da aplicação
  
  web:
    port: 8080                                                   # Porta onde a API escutará
    host: "0.0.0.0"                                              # Host da aplicação (0.0.0.0 para aceitar de qualquer IP)
    ssl_enabled: false                                           # Habilita ou não HTTPS (true/false)
    ssl_cert: "/etc/ssl/certs/cert.pem"                          # Caminho do certificado SSL
    ssl_key: "/etc/ssl/private/key.pem"                          # Caminho da chave privada SSL
    
    cors:                                                        # Configurações de CORS (Cross-Origin Resource Sharing)
      allow_origins: ["*"]                                       # Lista de domínios permitidos (ex.: ["https://meusite.com"])
      allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"] # Métodos HTTP permitidos
      allow_headers: ["Content-Type", "Authorization"]           # Headers permitidos na requisição
      allow_credentials: true                                    # Permite envio de cookies e credenciais (true/false)
      
  logger:
    level: "info"                                                # Nível de log: debug, info, warn, error
    format: "json"                                               # Formato do log: json ou console
    output: "stdout"                                             # Destino dos logs: stdout (terminal) ou file (arquivo)
    file_path: "./logs/app.log"                                  # Caminho do arquivo de log (usado se output for "file")

databases:
  postgres:
    dsn: |
      host=localhost              # Endereço do host do banco de dados
      port=5432                   # Porta padrão do PostgreSQL
      user=postgres               # Usuário para autenticação
      password=postgres           # Senha do usuário
      dbname=cleanarch_db         # Nome do banco de dados
      sslmode=disable             # Configuração de SSL: disable, require, verify-ca, verify-full
      application_name=cleanarch  # Nome da aplicação (aparece nas estatísticas do PostgreSQL)
      TimeZone=America/Sao_Paulo  # Timezone da conexão (ajuda na consistência de datas)
      search_path=public          # Schema padrão utilizado nas consultas SQL
    max_open_conns: 100           # Máximo de conexões abertas simultaneamente no pool
    max_idle_conns: 10            # Máximo de conexões ociosas no pool
    conn_max_lifetime: "1h"       # Tempo máximo que uma conexão pode ficar aberta (ex.: "1h", "30m") 
    auto_migrate: true            # Se verdadeiro, executa migrations automaticamente ao iniciar
    
cache:
  redis:
    host: "localhost"                                          # Host do Redis
    port: 6379                                                 # Porta do Redis
    password: ""                                               # Senha do Redis (se necessário)
    db: 0                                                      # Database index (Redis permite múltiplos DBs, padrão é 0)
    pool_size: 10                                              # Tamanho do pool de conexões com o Redis

auth:
  jwt:
    secret: "your-secret-key"                                  # Chave secreta para assinar os tokens JWT
    expiration: "24h"                                          # Tempo de expiração do access token (ex.: 24h)
    refresh_expiration: "168h"                                 # Expiração do refresh token (ex.: 168h = 7 dias)
    
  oidc:
    enabled: false                                             # Habilita ou não autenticação via OIDC
    client_id: "your-client-id"                                # Client ID registrado no provedor OIDC
    client_secret: "your-client-secret"                        # Client Secret do OIDC
    issuer: "https://your-oidc-provider.com"                   # URL do issuer OIDC (ex.: Keycloak, Auth0)
    redirect_url: "http://localhost:8080/auth/callback"        # URL de callback após autenticação

services:
  email:
    provider: "smtp"                                           # Provedor de email: smtp, sendgrid, ses
    smtp:
      host: "smtp.gmail.com"                                   # Host SMTP
      port: 587                                                # Porta SMTP (geralmente 587 para TLS)
      username: "your-email@gmail.com"                         # Usuário SMTP
      password: "your-password"                                # Senha SMTP
      from: "noreply@yourapp.com"                              # Email remetente

  storage:
    driver: "local"                                            # Tipo de storage: local ou s3
    local:
      path: "./uploads"                                        # Caminho local onde os arquivos serão armazenados
    s3:
      bucket: "your-bucket"                                    # Nome do bucket S3
      region: "us-east-1"                                      # Região do bucket (ex.: us-east-1)
      access_key: "your-access-key"                            # Chave de acesso AWS
      secret_key: "your-secret-key"                            # Chave secreta AWS
```

### Lendo Configurações com Viper

```go
// internal/config/config.go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Database DatabaseConfig `mapstructure:"databases"`
    Cache    CacheConfig    `mapstructure:"cache"`
    Auth     AuthConfig     `mapstructure:"auth"`
    Services ServicesConfig `mapstructure:"services"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.AddConfigPath(".")
    
    // Permitir variáveis de ambiente
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")
    
    // Bind de variáveis específicas
    viper.BindEnv("app.web.port", "APP_PORT")
    viper.BindEnv("databases.postgres.password", "DB_PASSWORD")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## Execução e Testes

### Executando localmente

```bash
# Desenvolvimento
make run

# Com hot reload
air

# Docker
docker-compose up

# Build
make build
```

### Makefile

```makefile
.PHONY: help run build test clean

help:
	@echo "Available commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v -cover ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

fmt:
	go fmt ./...

migrate-up:
	migrate -path migrations -database "${DATABASE_URL}" up

migrate-down:
	migrate -path migrations -database "${DATABASE_URL}" down

swagger:
	swag init -g cmd/api/main.go --output docs

docker-build:
	docker build -t myapp:latest .

docker-run:
	docker-compose up -d
```

### Executando testes

```bash
# Todos os testes
make test

# Com cobertura
make test-coverage

# Testes específicos
go test ./internal/usecase/...

# Testes de integração
go test ./test/integration/...
```

### Exemplo de Teste

```go
// internal/usecase/user/create_user_test.go
package user_test

import (
    "context"
    "testing"
    "meu-projeto/internal/entity"
    "meu-projeto/internal/usecase/user"
    "meu-projeto/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreateUser_Execute(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    mockLog := new(mocks.Logger)
    
    createUser := user.NewCreateUser(mockRepo, mockLog)
    
    t.Run("success", func(t *testing.T) {
        input := user.CreateUserInput{
            Name:     "John Doe",
            Email:    "john@example.com",
            Password: "password123",
        }
        
        mockRepo.On("GetByEmail", mock.Anything, input.Email).Return(nil, errors.New("not found"))
        mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
        mockLog.On("WithField", mock.Anything, mock.Anything).Return(mockLog)
        mockLog.On("Info", mock.Anything).Return()
        
        result, err := createUser.Execute(context.Background(), input)
        
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Equal(t, input.Name, result.Name)
        assert.Equal(t, input.Email, result.Email)
        
        mockRepo.AssertExpectations(t)
    })
    
    t.Run("email already exists", func(t *testing.T) {
        input := user.CreateUserInput{
            Name:     "John Doe",
            Email:    "existing@example.com",
            Password: "password123",
        }
        
        existingUser := &entity.User{Email: input.Email}
        mockRepo.On("GetByEmail", mock.Anything, input.Email).Return(existingUser, nil)
        
        result, err := createUser.Execute(context.Background(), input)
        
        assert.Error(t, err)
        assert.Nil(t, result)
        assert.Equal(t, "email already registered", err.Error())
    })
}
```

## 📏 Padrões e Convenções

### Nomenclatura

* **Arquivos**: snake_case (ex: `user_repository.go`)
* **Interfaces**: PascalCase com sufixo do tipo (ex: `UserRepository`)
* **Structs**: PascalCase (ex: `CreateUserUseCase`)
* **Funções/Métodos**: PascalCase para públicos, camelCase para privados
* **Pacotes**: lowercase sem underscore (ex: `usecase`, não `use_case`)

### Estrutura de Commits

```
tipo(escopo): descrição

[corpo opcional]

[rodapé opcional]
```

Tipos: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Exemplos:

```
feat(user): add user registration endpoint
fix(auth): correct JWT token expiration
docs(readme): update installation instructions
```

### Formatação e Linting

```bash
# Formatar código
make fmt

# Verificar linting
make lint

# Configuração do golangci-lint (.golangci.yml)
linters:
  enable:
    - gofmt
    - golint
    - govet
    - ineffassign
    - misspell
    - unconvert
    - prealloc
    - nakedret
```

### Versionamento

Seguir [Semantic Versioning](https://semver.org/):

* MAJOR: mudanças incompatíveis na API
* MINOR: novas funcionalidades compatíveis
* PATCH: correções de bugs compatíveis

### Estrutura de Erros

```go
// pkg/errors/errors.go
package errors

import "fmt"

type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Erros comuns
var (
    ErrNotFound     = &AppError{Code: "NOT_FOUND", Message: "Resource not found"}
    ErrUnauthorized = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
    ErrBadRequest   = &AppError{Code: "BAD_REQUEST", Message: "Invalid request"}
)
```


---

**Observações**: Use este template como base e adapte conforme as necessidades específicas do seu projeto. Mantenha a arquitetura limpa e as responsabilidades bem separadas!