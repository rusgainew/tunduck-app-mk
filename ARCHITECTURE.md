# 🏗️ Архитектура микросервисов на основе DDD

## Текущее состояние vs Целевое

### ДО: Monolith (11,282 LOC)

```
┌─────────────────────────────────────────┐
│         go-api (Single Module)          │
├─────────────────────────────────────────┤
│  Controllers (1,328 LOC)                │
│  ├── auth_controller                    │
│  ├── user_controller                    │
│  ├── role_controller                    │
│  ├── esf_organization_controller        │
│  └── esf_document_controller            │
├─────────────────────────────────────────┤
│  Services (2,000 LOC)                   │
│  ├── user_service                       │
│  ├── role_service                       │
│  ├── esf_organization_service           │
│  └── esf_document_service               │
├─────────────────────────────────────────┤
│  Models, Repository, Config             │
└─────────────────────────────────────────┘
        ↓
┌─────────────────────────────────────────┐
│     PostgreSQL (single database)        │
└─────────────────────────────────────────┘
```

### ПОСЛЕ: Microservices + DDD

```
┌──────────────────────────────────────────────────────────────┐
│                    API GATEWAY (Port 3000)                   │
│  Routes, Rate Limiting, Circuit Breaker, Auth Validation    │
└──────────────────────────┬─────────────────────────────────┘
        ↓          ↓           ↓              ↓
┌────────────┐ ┌────────────┐ ┌───────────┐ ┌─────────────┐
│   AUTH     │ │  COMPANY   │ │ DOCUMENT  │ │   (Future)  │
│  SERVICE   │ │  SERVICE   │ │ SERVICE   │ │  Services   │
│ (Port 3001)│ │(Port 3002) │ │(Port 3003)│ │             │
└─────┬──────┘ └─────┬──────┘ └─────┬─────┘ └─────────────┘
      ↓              ↓              ↓
┌──────────┐    ┌──────────┐   ┌──────────┐
│ auth_db  │    │company_db│   │doc_db    │
│PostgreSQL│    │PostgreSQL│   │PostgreSQL│
└──────────┘    └──────────┘   └──────────┘

              ┌─────────────────┐
              │  Shared Service │
              ├─────────────────┤
              │ Redis (Cache)   │
              │ Prometheus      │
              │ Logger (ELK?)   │
              │ Event Bus       │
              └─────────────────┘
```

---

## 📦 Структура каждого Microservice

### Папка микросервиса (например, auth-service)

```
auth-service/
│
├── cmd/
│   ├── main.go                 # Entry point
│   └── app.go                  # Application setup
│
├── internal/                   # Private code
│   │
│   ├── domain/                 # Domain Layer (Pure Business Logic)
│   │   ├── user.go            # User Aggregate Root
│   │   ├── user_test.go       # Domain tests
│   │   ├── credential.go      # Value Object (Email, Password)
│   │   ├── credential_test.go
│   │   ├── token.go           # Value Object (JWT Token)
│   │   ├── errors.go          # Domain errors (InvalidCredential, etc)
│   │   └── events.go          # Domain events (UserRegistered, etc)
│   │
│   ├── application/           # Application Layer (Use Cases)
│   │   ├── services/
│   │   │   ├── register_user_service.go      # Use case
│   │   │   ├── login_user_service.go         # Use case
│   │   │   ├── get_current_user_service.go   # Use case
│   │   │   ├── logout_user_service.go        # Use case
│   │   │   └── *_service_test.go             # Use case tests
│   │   │
│   │   ├── dto/
│   │   │   ├── register_user_dto.go
│   │   │   ├── login_user_dto.go
│   │   │   ├── user_response_dto.go
│   │   │   └── token_response_dto.go
│   │   │
│   │   └── commands/
│   │       ├── register_user_command.go
│   │       └── login_user_command.go
│   │
│   ├── infrastructure/        # Infrastructure Layer
│   │   │
│   │   ├── persistence/
│   │   │   ├── user_repository.go        # Interface
│   │   │   ├── postgres/
│   │   │   │   ├── user_postgres_repo.go
│   │   │   │   └── migration.go
│   │   │   └── postgres_test.go
│   │   │
│   │   ├── config/
│   │   │   ├── config.go     # Load from env
│   │   │   └── database.go   # DB connection
│   │   │
│   │   ├── http/
│   │   │   ├── client/       # HTTP clients to other services
│   │   │   │   └── company_service_client.go
│   │   │   └── middleware/
│   │   │       ├── jwt_middleware.go
│   │   │       ├── error_handler.go
│   │   │       └── logging_middleware.go
│   │   │
│   │   └── cache/
│   │       └── redis_cache.go
│   │
│   ├── interfaces/           # Interface Layer (API Contracts)
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── register_handler.go
│   │   │   │   ├── login_handler.go
│   │   │   │   ├── get_user_handler.go
│   │   │   │   └── logout_handler.go
│   │   │   └── routes.go
│   │   │
│   │   └── grpc/
│   │       ├── handlers/
│   │       │   └── auth_grpc_handler.go
│   │       └── client/              # gRPC clients to other services
│   │           └── company_client.go
│   │
│   ├── container.go          # DI Container (Dependency Injection)
│   └── errors.go             # Common errors
│
├── migrations/               # Database migrations
│   └── 001_create_users_table.sql
│
├── api/
│   ├── proto/                # gRPC Protocol Buffer definitions
│   │   ├── auth_service.proto
│   │   ├── auth.proto        # Message definitions
│   │   └── Makefile          # protoc compilation
│   │
│   └── openapi.yaml          # OpenAPI/Swagger для REST API
│
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml        # For local development
├── .env.example
├── Makefile
├── README.md
│
├── tests/                    # Integration tests
│   ├── integration_test.go
│   └── fixtures/
│
└── config/                   # Configuration templates
    └── logger.yaml
```

---

## 🎯 Domain-Driven Design (DDD) слои

### 1️⃣ Domain Layer (Ядро бизнес-логики)

**Файлы:** `internal/domain/*`
**Особенности:**

- ❌ Нет зависимостей от фреймворков
- ❌ Нет зависимостей от БД
- ✅ Только чистая бизнес-логика
- ✅ Entities, Value Objects, Aggregates, Domain Events
- ✅ Unit-тестируемо

**Пример:**

```go
// domain/user.go
package domain

import "errors"

type User struct {
    id         UUID
    email      Email        // Value Object
    password   Password     // Value Object (хэшированный)
    createdAt  time.Time
}

func (u *User) IsPasswordValid(plainPassword string) bool {
    return u.password.IsValid(plainPassword)
}

func (u *User) ChangePassword(newPassword Password) error {
    if !newPassword.IsStrong() {
        return errors.New("password is too weak")
    }
    u.password = newPassword
    return nil
}

// Value Object
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: email}, nil
}
```

### 2️⃣ Application Layer (Использование ЯДР)

**Файлы:** `internal/application/services/*`, `internal/application/dto/*`
**Особенности:**

- ✅ Orchestrates domain objects
- ✅ Use Cases / Services
- ✅ DTOs (Data Transfer Objects)
- ✅ Commands / Queries
- ❌ Не содержит бизнес-логику

**Пример:**

```go
// application/services/register_user_service.go
package services

import "context"

type RegisterUserService struct {
    userRepo      UserRepository  // Interface
    eventPublisher EventPublisher  // Interface
    logger        Logger          // Interface
}

func (s *RegisterUserService) Execute(ctx context.Context, cmd RegisterUserCommand) (*UserDTO, error) {
    // Проверяем, что юзер не существует
    exists, err := s.userRepo.ExistsByEmail(ctx, cmd.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, domain.NewUserAlreadyExistsError(cmd.Email)
    }

    // Создаем Domain Entity
    user := domain.NewUser(cmd.Email, cmd.Password)

    // Сохраняем в БД
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }

    // Публикуем Domain Event
    s.eventPublisher.Publish(domain.UserRegisteredEvent{
        UserID: user.ID(),
        Email:  cmd.Email,
    })

    return toUserDTO(user), nil
}
```

### 3️⃣ Infrastructure Layer (Технические детали)

**Файлы:** `internal/infrastructure/*`
**Особенности:**

- ✅ Database access (PostgreSQL)
- ✅ External API clients
- ✅ Caching (Redis)
- ✅ Configuration
- ✅ Logging, Metrics

**Пример:**

```go
// infrastructure/persistence/postgres/user_postgres_repo.go
package postgres

import (
    "context"
    "github.com/rusgainew/tunduck/internal/domain"
    "gorm.io/gorm"
)

type UserPostgresRepository struct {
    db *gorm.DB
}

func (r *UserPostgresRepository) Save(ctx context.Context, user *domain.User) error {
    entity := mapUserToDB(user)
    return r.db.WithContext(ctx).Create(entity).Error
}

func (r *UserPostgresRepository) GetByID(ctx context.Context, id domain.UUID) (*domain.User, error) {
    var entity UserEntity
    if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return mapUserFromDB(&entity), nil
}
```

### 4️⃣ Interfaces Layer (API контракты)

**Файлы:** `internal/interfaces/http/handlers/*`
**Особенности:**

- ✅ HTTP handlers / REST endpoints
- ✅ gRPC handlers (опционально)
- ✅ Валидация input
- ✅ Преобразование request → command
- ✅ Преобразование response ← DTO

**Пример:**

```go
// interfaces/http/handlers/register_handler.go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/rusgainew/tunduck/internal/application/services"
)

type RegisterHandler struct {
    service *services.RegisterUserService
    logger  Logger
}

// POST /auth/register
func (h *RegisterHandler) Handle(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{Message: "Invalid request"})
    }

    // Валидация
    if err := h.validate(req); err != nil {
        return c.Status(400).JSON(ErrorResponse{Message: err.Error()})
    }

    // Выполняем Use Case
    dto, err := h.service.Execute(c.Context(), toCommand(req))
    if err != nil {
        return c.Status(500).JSON(ErrorResponse{Message: err.Error()})
    }

    return c.Status(201).JSON(toResponse(dto))
}
```

---

## 🔄 Communication Between Services

### Вариант 1: REST + Event Bus (Рекомендуется для начала)

```
Service A ──HTTP GET/POST──> Service B
Service A ──Event──> Redis Pub/Sub ──> Service B
Service A ──Event──> RabbitMQ ──> Service B
```

**Преимущества:**

- Простая интеграция
- Легко дебажить (HTTP)
- Async communication с Event Bus

**Недостатки:**

- Медленнее чем gRPC
- Больше данных

### Вариант 2: gRPC (Для будущего)

```
Service A ──gRPC──> Service B (протоколбуферы)
```

**Преимущества:**

- Быстрее (бинарный протокол)
- Type-safe (protobuf)
- Bi-directional streaming

**Недостатки:**

- Сложнее интегрировать с фронтом
- Требует proto файлов

**Рекомендация:** Начать с REST/HTTP, потом добавить gRPC для critical path

---

## 🗄️ Database Strategy

### Вариант 1: Shared Database (СЕЙЧАС - простой старт)

```
PostgreSQL
├── users (auth-service owns)
├── organizations (company-service owns)
├── documents (document-service owns)
├── roles (shared)
└── audit_logs (shared)

Pros:
  ✅ Простая интеграция
  ✅ ACID транзакции между сервисами
  ✅ Легче мигрировать

Cons:
  ❌ Сложно масштабировать
  ❌ Tight coupling на уровне БД
```

### Вариант 2: Database Per Service (БУДУЩЕЕ - масштабирование)

```
PostgreSQL Cluster:
├── auth_db (auth-service)
├── company_db (company-service)
├── document_db (document-service)
└── shared_db (shared tables)

Communication:
  - REST API для данных других сервисов
  - Event Bus для синхронизации
  - API Gateway для join-ов

Pros:
  ✅ Отличное масштабирование
  ✅ Автономность сервисов
  ✅ Разные DB для разных сервисов

Cons:
  ❌ Сложнее с транзакциями (Saga pattern)
  ❌ Требует Event Sourcing
  ❌ Дороже в операционном плане
```

**Рекомендация:** Начать с Shared Database, потом per-service

---

## 🔌 Service Interfaces

### Auth Service API

```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "StrongPassword123!",
  "firstName": "John",
  "lastName": "Doe"
}

Response: 201 Created
{
  "id": "uuid",
  "email": "user@example.com",
  "token": "jwt_token",
  "expiresAt": "2025-01-01T00:00:00Z"
}
```

```http
POST /api/auth/login
{
  "email": "user@example.com",
  "password": "StrongPassword123!"
}

Response: 200 OK
{
  "token": "jwt_token",
  "user": { ... }
}
```

### Internal gRPC Interface (для других сервисов)

```proto
service AuthService {
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

---

## 🧪 Testing Strategy

### Unit Tests (Domain Layer)

```go
// internal/domain/user_test.go
func TestUserRegistration(t *testing.T) {
    email, _ := domain.NewEmail("test@example.com")
    password, _ := domain.NewPassword("StrongPass123!")

    user := domain.NewUser(email, password)

    assert.NotNil(t, user.ID())
    assert.Equal(t, email, user.Email())
}
```

### Integration Tests (Full Stack)

```go
// tests/integration_test.go
func TestRegisterUserFlow(t *testing.T) {
    // Setup DB
    db := setupTestDB(t)
    repo := postgres.NewUserRepository(db)
    service := services.NewRegisterUserService(repo)
    handler := handlers.NewRegisterHandler(service)

    // Test request
    req := testutil.NewRegisterRequest("test@example.com", "Pass123!")
    res := handler.Handle(req)

    assert.Equal(t, 201, res.StatusCode)
}
```

### Contract Tests (Service Communication)

```go
// tests/contract_test.go
func TestAuthServiceContractWithCompanyService(t *testing.T) {
    // Company Service вызывает Auth Service API
    user, err := authClient.GetUser(ctx, userID)
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

---

## 🚀 Deployment & DevOps

### Docker для каждого сервиса

```dockerfile
# auth-service/Dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o auth-service ./cmd

FROM alpine:latest
COPY --from=builder /app/auth-service .
EXPOSE 3001
CMD ["./auth-service"]
```

### Docker Compose для local development

```yaml
# docker-compose.yml
version: "3.8"

services:
  auth-service:
    build: ./auth-service
    ports:
      - "3001:3001"
    environment:
      - DB_HOST=postgres
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis

  company-service:
    build: ./company-service
    ports:
      - "3002:3002"
    depends_on:
      - postgres
      - redis

  document-service:
    build: ./document-service
    ports:
      - "3003:3003"
    depends_on:
      - postgres
      - redis

  api-gateway:
    build: ./api-gateway
    ports:
      - "3000:3000"
    depends_on:
      - auth-service
      - company-service
      - document-service

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

---

## 📊 Жизненный цикл разработки

```
Sprint 1: Подготовка
├── Утвердить DDD архитектуру ✅ (текущее)
├── Создать шаблон микросервиса
└── Настроить build/test pipeline

Sprint 2: Auth Service
├── Мигрировать auth логику
├── Применить DDD паттерны
├── Написать тесты
└── Развернуть первый микросервис

Sprint 3: Company Service
├── Мигрировать organization логику
├── Применить DDD паттерны
└── Интеграция с Auth Service

Sprint 4: Document Service
└── Аналогично Company Service

Sprint 5: API Gateway
├── Создать gateway для маршрутизации
├── Service discovery
├── Circuit breaker patterns
└── Rate limiting

Sprint 6+: Расширение
├── Add gRPC если нужно
├── Добавить Event Sourcing
├── Масштабирование
└── Database per service миграция
```

---

## ✅ Чек-лист перед началом разработки

- [ ] Утвердить Communication Strategy (REST vs gRPC)
- [ ] Утвердить Database Strategy (Shared vs Per-service)
- [ ] Создать шаблон проекта (scaffold)
- [ ] Настроить CI/CD для микросервисов
- [ ] Подготовить Docker registry
- [ ] Настроить monitoring & logging
- [ ] Определить naming conventions
- [ ] Создать shared library как separate module
- [ ] Написать guidelines для разработчиков
- [ ] Подготовить документацию по API (OpenAPI)
