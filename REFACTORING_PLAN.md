# План рефакторинга на микросервисы DDD

## 📊 Анализ текущего кода

### Метрики

- **Общий объём:** 11,282 строк кода
- **Стек:** Go 1.25, Fiber, PostgreSQL, Redis, JWT
- **Тестовое покрытие:** Есть unit и integration тесты
- **CI/CD:** Docker, docker-compose

### Текущая архитектура (Monolith)

```
go-api/
├── cmd/api/              # Entry point
├── internal/
│   ├── controllers/      # HTTP handlers (auth, user, org, doc)
│   ├── services/         # Business logic
│   ├── models/           # Domain entities
│   ├── repository/       # Data access (PostgreSQL)
│   └── conf/             # Configuration
└── pkg/                  # Shared utilities (cache, auth, logger, etc)
```

### Текущие доменные области

1. **Auth & Registration** - login, register, JWT
2. **User Management** - RBAC, user profiles
3. **Organization (Company)** - компании, структура
4. **ESF Documents** - документооборот
5. **Infrastructure** - caching, logging, metrics, health

---

## 🎯 Целевая архитектура (Микросервисы + DDD)

### Bounded Contexts (DDD)

```
1. AUTH-SERVICE (Microservice)
   └── domain: User, Credential, Token

2. COMPANY-SERVICE (Microservice)
   └── domain: Organization, Role, Permission

3. DOCUMENT-SERVICE (Microservice)
   └── domain: Document, DocumentEntry

4. API-GATEWAY (Entry point)
   └── Маршрутизация между сервисами

5. SHARED-SERVICES (Shared libs)
   └── Logger, Cache, Metrics, Health, Transaction
```

### Структура каждого микросервиса (DDD)

```
auth-service/
├── cmd/                          # Entry point
│   └── main.go
├── internal/
│   ├── domain/                   # Core бизнес-логика (entities, value objects)
│   │   ├── user.go              # User aggregate root
│   │   ├── credential.go
│   │   └── token.go
│   ├── application/              # Use cases, DTOs
│   │   ├── services/
│   │   │   └── user_service.go
│   │   └── dto/
│   ├── infrastructure/           # DB, HTTP clients, external services
│   │   ├── persistence/
│   │   │   ├── user_repository.go
│   │   │   └── postgres/
│   │   ├── http/
│   │   │   └── handlers/
│   │   ├── grpc/
│   │   │   └── client/           # gRPC clients to other services
│   │   └── config/
│   └── interfaces/              # API contracts
│       ├── http/
│       │   └── handlers/
│       └── grpc/
│           └── handlers/         # gRPC handlers
├── api/
│   ├── proto/                    # 📁 ВСЕ PROTO ФАЙЛЫ ЗДЕСЬ!
│   │   ├── auth_service.proto   # Service definition
│   │   ├── auth.proto           # Message definitions
│   │   └── Makefile             # protoc compilation
│   └── openapi.yaml             # REST API documentation
├── migrations/
│   └── 001_create_users_table.sql
└── go.mod
```

---

## 📋 Пошаговый план реализации

### Phase 1: Подготовка

- [x] Анализ текущего кода
- [x] **Создать структуру `api/proto/` для всех сервисов** ✅
- [x] **Создать proto файлы:** ✅
  - [x] `auth_service.proto` - gRPC service definition
  - [x] `auth.proto` - User, Token, Credential messages
  - [x] `company_service.proto` - Organization management service
  - [x] `company.proto` - Organization, Employee, Department messages
  - [x] `document_service.proto` - Document workflow service
  - [x] `document.proto` - Document, DocumentEntry, DocumentWorkflow messages
  - [x] `common.proto` - Shared messages (Empty, Error, Pagination)
- [x] **Создать Makefile для protoc compilation** ✅

### Phase 2: Рефакторинг Auth-Service

1. Создать отдельный модуль `auth-service`
2. Применить DDD паттерны
3. Создать HTTP handlers (REST)
4. Создать gRPC handlers для других сервисов
5. Написать unit тесты для domain layer
6. Написать integration тесты

### Phase 3: Company-Service

1. Аналогично Auth Service
2. Добавить RabbitMQ event subscription (слушать UserRegisteredEvent)
3. gRPC клиент для вызова Auth Service

### Phase 4: Document-Service

1. Аналогично Company Service
2. RabbitMQ event subscription
3. gRPC клиенты для Company и Auth Service

### Phase 5: API-Gateway

1. Создать gateway для HTTP routing
2. gRPC load balancing
3. Circuit breaker patterns для gRPC
4. Rate limiting

### Phase 6: Event System & Integration

1. Полная интеграция RabbitMQ
2. Dead Letter Queue handling
3. Event Sourcing (опционально)
4. Distributed tracing (Jaeger)

### Phase 7: DevOps/Deployment

1. docker-compose для всех сервисов + RabbitMQ
2. gRPC и HTTP communication настроены
3. Database per service strategy (если нужно)

---

## 🔄 Примеры миграции

### ДО (Monolith)

```go
// controllers/auth_controller.go
func (c *AuthController) register(ctx *fiber.Ctx) error {
    var req RegisterRequest
    if err := ctx.BodyParser(&req); err != nil { ... }
    user, err := c.service.Register(ctx.Context(), req)
    if err != nil { ... }
    return ctx.JSON(user)
}
```

### ПОСЛЕ (DDD)

```go
// domain/user.go
type User struct {
    id        UUID
    email     Email // Value Object
    password  Password // Value Object
}

func (u *User) Register(email Email, password Password) error {
    if err := u.validateEmail(email); err != nil {
        return DomainError{...}
    }
    u.email = email
    u.password = password.hash()
    return nil
}

// application/services/register_user.go
type RegisterUserService struct {
    repo UserRepository
    evt  EventPublisher
}

func (s *RegisterUserService) Execute(cmd RegisterUserCommand) (*UserDTO, error) {
    user := domain.NewUser(cmd.Email, cmd.Password)
    if err := s.repo.Save(user); err != nil { ... }
    s.evt.Publish(UserRegisteredEvent{...})
    return toDTO(user), nil
}

// interfaces/handlers/auth_handler.go
func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
    var req RegisterRequest
    dto, err := h.service.Execute(toCommand(req))
    return ctx.JSON(toResponse(dto))
}
```

---

## 🔌 Service Communication

### Вариант 1: REST + Event Bus

```
Auth-Service ──HTTP──> Company-Service
Auth-Service ──Event──> Redis Pub/Sub ──> Company-Service
```

### Вариант 2: gRPC

```
Auth-Service ──gRPC──> Company-Service
```

**Рекомендация:** Начать с REST, позже перейти на gRPC

---

## 📦 Shared Library

```go
// shared/pkg/logger
// shared/pkg/cache (Redis)
// shared/pkg/metrics (Prometheus)
// shared/pkg/middleware (Auth, Logging, etc)
// shared/pkg/event (Event Bus)
// shared/pkg/health (Health checks)
```

---

## ✅ Чек-лист для каждого сервиса

- [ ] Создана папка структура
- [ ] Определены Domain entities
- [ ] Написаны интерфейсы Repository
- [ ] Реализовано PostgreSQL repository
- [ ] Написаны Application services (Use cases)
- [ ] Созданы HTTP handlers
- [ ] Написаны unit тесты (domain layer)
- [ ] Написаны integration тесты
- [ ] Документирован API (Swagger/OpenAPI)
- [ ] Настроена CI/CD для сервиса

---

## 🚀 Следующие шаги

1. **Утвердить архитектуру** - обсудить варианты communication
2. **Создать шаблон (scaffold)** для нового микросервиса
3. **Начать с Auth-Service** - самый простой для миграции
4. **Постепенно расширять** - Company, Documents, и т.д.
