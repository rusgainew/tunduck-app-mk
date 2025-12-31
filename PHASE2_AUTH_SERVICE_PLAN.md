# Phase 2: Auth-Service Implementation Plan

## 📋 Общая информация

**Duration:** 2-3 недели  
**Team:** 1-2 разработчика  
**Priority:** Высокий (foundation service)  
**Dependencies:** Phase 1 ✅

---

## 🎯 Цели Phase 2

1. ✅ Мигрировать логику аутентификации из monolith
2. ✅ Применить DDD паттерны для auth domain
3. ✅ Реализовать gRPC сервис для других микросервисов
4. ✅ Добавить RabbitMQ event publishing
5. ✅ Написать comprehensive тесты
6. ✅ Развернуть первый микросервис в Docker

---

## 📁 Структура Auth-Service

```
auth-service/
├── cmd/
│   └── main.go                          # Entry point, инициализация контейнера
│
├── internal/
│   ├── domain/                          # 🔴 DOMAIN LAYER (чистая бизнес-логика)
│   │   ├── user.go                      # User aggregate root
│   │   ├── user_repository.go           # Repository interface
│   │   ├── credential.go                # Value object (email, password)
│   │   ├── token.go                     # Value object (JWT token)
│   │   ├── errors.go                    # Domain errors
│   │   └── events.go                    # Domain events
│   │
│   ├── application/                     # 🟡 APPLICATION LAYER (use cases)
│   │   ├── services/
│   │   │   ├── register_user_service.go
│   │   │   ├── login_user_service.go
│   │   │   ├── validate_token_service.go
│   │   │   └── logout_user_service.go
│   │   │
│   │   ├── dto/                         # Data Transfer Objects
│   │   │   ├── register_user_dto.go
│   │   │   ├── login_user_dto.go
│   │   │   └── token_dto.go
│   │   │
│   │   └── events/
│   │       ├── user_registered_event.go
│   │       ├── user_logged_in_event.go
│   │       └── event_publisher.go
│   │
│   ├── infrastructure/                  # 🟢 INFRASTRUCTURE LAYER (external)
│   │   ├── persistence/
│   │   │   ├── user_repository.go       # Interface implementation
│   │   │   ├── postgresql/
│   │   │   │   ├── user_postgres_repo.go
│   │   │   │   └── queries.go
│   │   │   └── redis/
│   │   │       └── token_blacklist.go   # Token revocation storage
│   │   │
│   │   ├── rabbitmq/
│   │   │   ├── event_publisher.go       # RabbitMQ publisher
│   │   │   └── connection.go
│   │   │
│   │   ├── http/
│   │   │   └── clients/
│   │   │       └── identity_verifier.go # Optional: email verification
│   │   │
│   │   └── config/
│   │       └── config.go                # Environment configuration
│   │
│   ├── interfaces/                      # 🔵 INTERFACES LAYER (API contracts)
│   │   ├── http/
│   │   │   ├── routes.go                # HTTP routes definition
│   │   │   ├── middleware/
│   │   │   │   ├── auth_middleware.go   # JWT validation middleware
│   │   │   │   └── error_handler.go     # Error response formatting
│   │   │   └── handlers/
│   │   │       ├── register_handler.go  # POST /auth/register
│   │   │       ├── login_handler.go     # POST /auth/login
│   │   │       ├── me_handler.go        # GET /auth/me
│   │   │       ├── logout_handler.go    # POST /auth/logout
│   │   │       └── refresh_handler.go   # POST /auth/refresh
│   │   │
│   │   ├── grpc/
│   │   │   ├── handlers/
│   │   │   │   └── auth_grpc_handler.go # gRPC server implementation
│   │   │   └── interceptors/
│   │   │       └── logging_interceptor.go
│   │   │
│   │   └── openapi.yaml                 # REST API documentation
│   │
│   └── container.go                     # Dependency Injection (wire or fx)
│
├── api/
│   ├── proto/                           # 📌 Proto files (from centralized location)
│   │   ├── auth_service.proto           # gRPC service definition
│   │   ├── auth.proto                   # Message definitions
│   │   ├── common.proto                 # Shared messages
│   │   └── Makefile                     # Compilation
│   │
│   └── openapi.yaml                     # REST API specification
│
├── migrations/
│   ├── 001_create_users_table.sql
│   ├── 002_create_user_sessions_table.sql
│   └── migrations.go                    # Migration runner
│
├── tests/
│   ├── domain/
│   │   ├── user_test.go
│   │   └── credential_test.go
│   ├── application/
│   │   ├── register_user_service_test.go
│   │   └── login_user_service_test.go
│   ├── integration/
│   │   ├── http_handlers_test.go        # HTTP endpoint tests
│   │   ├── grpc_handlers_test.go        # gRPC service tests
│   │   └── database_test.go             # Repository tests
│   └── fixtures/
│       └── test_data.go
│
├── docker/
│   ├── Dockerfile                       # Multi-stage build
│   └── docker-compose.yml               # Local development
│
├── docs/
│   ├── API.md                           # REST API documentation
│   ├── GRPC.md                          # gRPC documentation
│   ├── ARCHITECTURE.md                  # Service architecture
│   └── EVENTS.md                        # RabbitMQ events guide
│
├── go.mod                               # Dependencies
├── go.sum
├── Makefile                             # Build and test commands
└── README.md                            # Service overview
```

---

## 🔄 DDD Pattern Implementation

### Domain Layer: User Aggregate Root

```go
// internal/domain/user.go

package domain

import "time"

// User aggregate root
type User struct {
    id        UserID
    email     Email
    password  Password
    firstName string
    lastName  string
    status    UserStatus
    createdAt time.Time
    updatedAt time.Time
}

// Value Objects
type UserID string
type Email string
type Password string

// Domain Events
type UserRegistered struct {
    UserID    UserID
    Email     Email
    Timestamp time.Time
}

type UserLoggedIn struct {
    UserID    UserID
    Timestamp time.Time
}

// Repository Interface (defined in domain, implemented in infrastructure)
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email Email) (*User, error)
    FindByID(ctx context.Context, id UserID) (*User, error)
}
```

### Application Layer: Use Cases

```go
// internal/application/services/register_user_service.go

package services

type RegisterUserService struct {
    repo      domain.UserRepository
    publisher domain.EventPublisher
}

func (s *RegisterUserService) Execute(ctx context.Context, cmd RegisterUserCommand) (*UserDTO, error) {
    // Validate input
    // Check if email already exists
    // Hash password
    // Create user aggregate
    // Save to database
    // Publish UserRegistered event
    // Return DTO
}
```

### Infrastructure Layer: Implementations

```go
// internal/infrastructure/persistence/postgresql/user_postgres_repo.go

package postgresql

type UserPostgresRepository struct {
    db *pgxpool.Pool
}

func (r *UserPostgresRepository) Save(ctx context.Context, user *domain.User) error {
    // Insert into users table
}

func (r *UserPostgresRepository) FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
    // Query by email
}
```

### Interfaces Layer: HTTP Handlers

```go
// internal/interfaces/http/handlers/register_handler.go

package handlers

type RegisterHandler struct {
    service services.RegisterUserService
}

func (h *RegisterHandler) Handle(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{Message: "Invalid request"})
    }

    // Validate request
    // Call service
    // Return response
}
```

---

## 📋 Implementation Checklist

### Week 1: Setup & Domain Layer

- [ ] **Setup Project Structure**

  - [ ] Create `auth-service/` directory
  - [ ] Initialize `go.mod`
  - [ ] Setup project layout

- [ ] **Domain Layer**

  - [ ] Implement `User` aggregate
  - [ ] Implement `Credential` value object
  - [ ] Implement `Token` value object
  - [ ] Define domain errors
  - [ ] Define domain events (UserRegistered, UserLoggedIn)
  - [ ] Define `UserRepository` interface
  - [ ] Write domain layer unit tests

- [ ] **Proto Files**
  - [ ] Copy proto files from `api/proto/`
  - [ ] Run `make proto` to generate Go code
  - [ ] Verify generated `*.pb.go` files

### Week 2: Application & Infrastructure

- [ ] **Application Layer**

  - [ ] Implement `RegisterUserService`
  - [ ] Implement `LoginUserService`
  - [ ] Implement `ValidateTokenService`
  - [ ] Implement `LogoutUserService`
  - [ ] Create DTOs for each service
  - [ ] Write application layer unit tests

- [ ] **Infrastructure Layer**

  - [ ] Implement `UserPostgresRepository`
  - [ ] Create database migrations
  - [ ] Implement `RabbitMQEventPublisher`
  - [ ] Setup configuration management
  - [ ] Write repository tests

- [ ] **Database**
  - [ ] Create `users` table schema
  - [ ] Create `user_sessions` table (for logout)
  - [ ] Add indexes for email and ID

### Week 3: Interfaces & Testing

- [ ] **HTTP Handlers**

  - [ ] Implement `RegisterHandler`
  - [ ] Implement `LoginHandler`
  - [ ] Implement `MeHandler`
  - [ ] Implement `LogoutHandler`
  - [ ] Implement `RefreshHandler`
  - [ ] Setup HTTP routes

- [ ] **gRPC Handlers**

  - [ ] Implement `AuthServiceServer`
  - [ ] Implement `ValidateToken` RPC
  - [ ] Implement `GetUser` RPC
  - [ ] Setup gRPC server

- [ ] **Testing**

  - [ ] Write HTTP endpoint tests
  - [ ] Write gRPC service tests
  - [ ] Write integration tests
  - [ ] Achieve 80%+ code coverage

- [ ] **Documentation**

  - [ ] Write API documentation
  - [ ] Write gRPC documentation
  - [ ] Write architecture guide
  - [ ] Write events guide

- [ ] **Deployment**
  - [ ] Create Dockerfile
  - [ ] Create docker-compose.yml
  - [ ] Setup CI/CD pipeline
  - [ ] Deploy to development environment

---

## 🔧 Key Technologies & Tools

| Component     | Technology       | Purpose                     |
| ------------- | ---------------- | --------------------------- |
| Framework     | Fiber (Go)       | HTTP server                 |
| gRPC          | Protocol Buffers | Inter-service communication |
| Database      | PostgreSQL       | User persistence            |
| Caching       | Redis            | Token blacklist             |
| Message Queue | RabbitMQ         | Event publishing            |
| Testing       | Testify, GoMock  | Unit & integration tests    |
| DI Container  | wire or uber/fx  | Dependency injection        |
| Logging       | Logrus           | Structured logging          |
| Migration     | golang-migrate   | Database versioning         |

---

## 🚀 Development Workflow

### Local Development

```bash
# 1. Clone proto files
cp -r ../api/proto .

# 2. Compile proto files
cd api/proto && make proto

# 3. Run migrations
make db-migrate

# 4. Start service in watch mode
make dev

# 5. Run tests
make test

# 6. Build Docker image
make docker-build
```

### Testing Strategy

```bash
# Unit tests (domain & services)
make test-unit

# Integration tests
make test-integration

# All tests with coverage
make test-coverage

# End-to-end tests
make test-e2e
```

---

## 📊 Success Criteria

✅ Proto files compiled successfully  
✅ All domain entities implemented with tests  
✅ All application services implemented  
✅ All repository implementations working  
✅ HTTP handlers returning correct responses  
✅ gRPC server accepting requests from clients  
✅ RabbitMQ events being published  
✅ 80%+ code coverage  
✅ All integration tests passing  
✅ Docker image building successfully  
✅ Service running in docker-compose

---

## 📝 Next Steps (Phase 3)

After Auth-Service is complete:

1. Extract Company-Service (esf_organization)
2. Integrate with Auth-Service via gRPC
3. Add RabbitMQ event consumption
4. Repeat for Document-Service
5. Build API-Gateway

---

## 🔗 References

- [DDD in Go](https://threedots.tech/post/ddd-lite-in-go/)
- [Fiber Documentation](https://docs.gofiber.io)
- [gRPC Best Practices](https://grpc.io/docs/guides/performance-best-practices/)
- [RabbitMQ Go Client](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
- [PostgreSQL in Go](https://pkg.go.dev/github.com/jackc/pgx/v5)
