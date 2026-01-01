# 🎯 PROJECT MASTER GUIDE - Полное руководство по проекту

**Дата:** 1 января 2026  
**Статус:** ✅ Phase 1 Complete - Ready for Phase 2  
**Repository:** https://github.com/rusgainew/tunduck-app-mk

---

## ⚡ QUICKSTART (10 минут)

### Минимум для старта

```bash
# 1. Клонировать и перейти
git clone https://github.com/rusgainew/tunduck-app-mk.git
cd tunduck-app-mk

# 2. Установить инструменты
brew install protobuf          # macOS
# или: sudo apt-get install protobuf-compiler  # Linux

# 3. Запустить инфраструктуру
docker-compose -f docker-compose.microservices.yml up -d postgres redis rabbitmq

# 4. Компилировать proto файлы
cd api/proto && make proto

# 5. Проверить здоровье сервисов
curl http://localhost:5432/  # PostgreSQL
redis-cli ping              # Redis
# RabbitMQ Management: http://localhost:15672 (guest/guest)
```

---

## 📊 СОСТОЯНИЕ ПРОЕКТА

### Текущая ситуация

- **Текущая архитектура:** Monolith (go-api)
- **Размер:** 11,282 строк кода
- **Стек:** Go 1.25, Fiber, PostgreSQL, Redis, JWT, Prometheus, Docker
- **Стартегия:** Миграция на микросервисы + DDD

### Целевая архитектура

```
┌──────────────────────────────────────────┐
│         API GATEWAY (Port 3000)          │
│  Routes, Rate Limit, Auth, Circuit Break│
└──────┬──────────┬──────────┬─────────────┘
       ↓          ↓          ↓
┌─────────────┐ ┌──────────┐ ┌──────────┐
│AUTH-SERVICE │ │COMPANY-  │ │DOCUMENT- │
│(Port 3001)  │ │SERVICE   │ │SERVICE   │
│             │ │(Port 3002)│ │(Port 3003)│
└────┬────────┘ └──────┬───┘ └──────┬───┘
     ↓                ↓             ↓
┌──────────┐  ┌──────────┐  ┌──────────┐
│ auth_db  │  │company_db│  │ doc_db   │
│PostgreSQL│  │PostgreSQL│  │PostgreSQL│
└──────────┘  └──────────┘  └──────────┘
```

---

## 🏗️ АРХИТЕКТУРА

### Domain-Driven Design (DDD)

Каждый микросервис состоит из **4 слоёв**:

#### 1️⃣ **Domain Layer** (`internal/domain/*`)

- Чистая бизнес-логика (без фреймворков)
- Entities, Value Objects, Aggregates
- Domain Events
- Пример: `user.go`, `credential.go`, `token.go`

```go
// domain/user.go
type User struct {
    id       UUID
    email    Email           // Value Object
    password Password        // Value Object (хешированный)
    roles    []Role
}

func (u *User) Register(email Email, pwd Password) error {
    if email.IsInvalid() {
        return ErrInvalidEmail
    }
    // ...
}
```

#### 2️⃣ **Application Layer** (`internal/application/*`)

- Use Cases (бизнес-процессы)
- DTOs (Data Transfer Objects)
- Commands & Queries
- Пример: `services/register_user_service.go`, `dto/register_user_dto.go`

```go
// application/services/register_user_service.go
type RegisterUserService struct {
    userRepo UserRepository
    cache    Cache
}

func (s *RegisterUserService) Execute(dto RegisterUserDTO) (*User, error) {
    // Валидация, создание, сохранение
}
```

#### 3️⃣ **Infrastructure Layer** (`internal/infrastructure/*`)

- Database (Repository pattern)
- HTTP clients для других сервисов
- gRPC clients
- Cache, Config, External APIs

```
infrastructure/
├── persistence/
│   ├── user_repository.go        (Interface)
│   └── postgres/
│       └── user_postgres_repo.go (Implementation)
├── http/
│   └── handlers/
│       └── company_client.go     (HTTP client)
└── config/
    ├── database.go
    └── cache.go
```

#### 4️⃣ **Interfaces Layer** (`internal/interfaces/*`)

- HTTP handlers (REST API)
- gRPC handlers
- Routes
- Middleware (auth, logging, errors)

```
interfaces/
├── http/
│   ├── handlers/
│   │   ├── register_handler.go
│   │   └── login_handler.go
│   └── routes.go
└── grpc/
    ├── handlers/
    │   └── auth_grpc_handler.go
    └── client/
        └── company_client.go
```

### Структура микросервиса

```
auth-service/
├── cmd/auth-service/main.go      # Entry point
├── internal/
│   ├── domain/                   # Business logic (entities, values)
│   │   ├── user.go
│   │   ├── credential.go
│   │   ├── token.go
│   │   ├── errors.go
│   │   └── events.go
│   │
│   ├── application/              # Use cases, DTOs
│   │   ├── services/
│   │   │   ├── register_user_service.go
│   │   │   ├── login_user_service.go
│   │   │   └── *_test.go
│   │   └── dto/
│   │       ├── register_user_dto.go
│   │       └── user_response_dto.go
│   │
│   ├── infrastructure/           # Database, HTTP, Cache
│   │   ├── persistence/
│   │   │   ├── user_repository.go (interface)
│   │   │   └── postgres/
│   │   │       ├── user_postgres_repo.go
│   │   │       └── migration.go
│   │   ├── http/
│   │   │   └── company_client.go (gRPC/REST client)
│   │   └── config/
│   │       ├── database.go
│   │       └── cache.go
│   │
│   ├── interfaces/               # API contracts
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── register_handler.go
│   │   │   │   ├── login_handler.go
│   │   │   │   └── routes.go
│   │   │   └── middleware/
│   │   │       └── jwt_middleware.go
│   │   └── grpc/
│   │       ├── handlers/
│   │       │   └── auth_grpc_handler.go
│   │       └── client/
│   │           └── company_client.go
│   │
│   └── container.go              # Dependency Injection
│
├── api/proto/
│   ├── auth_service.proto        # gRPC service definition
│   ├── auth.proto                # Message definitions
│   └── Makefile
│
├── migrations/
│   └── 001_create_users_table.sql
│
├── go.mod
├── Dockerfile
└── README.md
```

---

## 🔌 SERVICE COMMUNICATION

### 1. Синхронные вызовы → **gRPC**

```
AuthService ──gRPC──> CompanyService
↓
Используем protobuf для типизированных быстрых вызовов
```

**Proto пример:**

```protobuf
service CompanyService {
    rpc GetOrganization(GetOrganizationRequest) returns (Organization);
    rpc CreateOrganization(CreateOrganizationRequest) returns (Organization);
}

message GetOrganizationRequest {
    string organization_id = 1;
}

message Organization {
    string id = 1;
    string name = 2;
    repeated Employee employees = 3;
}
```

### 2. Асинхронные события → **RabbitMQ**

```
AuthService ──Event──> RabbitMQ ──> CompanyService (subscribes)

Events:
- UserRegistered { user_id, email, created_at }
- UserRoleChanged { user_id, roles }
- DocumentCreated { document_id, owner_id }
```

### 3. Кэширование → **Redis**

```
CompanyService ──Query──> Redis Cache
    ↓ (miss)
    └──> PostgreSQL ──> Cache ──> Client
```

---

## 📋 РАЗДЕЛ ПО РАЗРАБОТКЕ

### Доменные области (Bounded Contexts)

#### 🔐 **AUTH-SERVICE**

**Файлы из монолита:** `controllers/auth_controller.go`, `services/user_service_impl.go`, `models/user_model.go`, `repository/user_repository*`

**Entities:**

- `User` (aggregate root)
- `Credential` (email, password)
- `Token` (JWT)
- `Role` (user roles)

**Use Cases (RPC методы):**

- `Register(email, password) → AuthResponse`
- `Login(email, password) → AuthResponse`
- `ValidateToken(token) → User`
- `RefreshToken(token) → Token`
- `Logout() → Empty`
- `GetUser(user_id) → User`

**Dependencies:** PostgreSQL

---

#### 🏢 **COMPANY-SERVICE**

**Файлы из монолита:** `controllers/esf_organization_controller.go`, `services/esf_organization_service_impl.go`, `models/esf_organization_model.go`, `repository/esf_organization_postgres.go`

**Entities:**

- `Organization` (aggregate root)
- `Employee`
- `Department`
- `OrganizationRole`

**Use Cases (RPC методы):**

- `GetOrganization(org_id) → Organization`
- `CreateOrganization(name, owner_id) → Organization`
- `UpdateOrganization(org_id, updates) → Organization`
- `DeleteOrganization(org_id) → Empty`
- `ListOrganizations(page) → List<Organization>`
- `GetOrganizationMembers(org_id) → List<Employee>`
- `AddMember(org_id, user_id, role) → Employee`
- `RemoveMember(org_id, user_id) → Empty`

**Dependencies:** PostgreSQL, Auth-Service (gRPC client)

**Events subscribed:**

- `UserRegistered` (для автоматического добавления в организацию)

---

#### 📄 **DOCUMENT-SERVICE**

**Файлы из монолита:** `controllers/esf_document_controller.go`, `services/esf_document_service_impl.go`, `models/esf_document_model.go`, `repository/esf_document_postgres.go`

**Entities:**

- `Document` (aggregate root)
- `DocumentEntry`
- `DocumentWorkflow`
- `DocumentTemplate`

**Use Cases (RPC методы):**

- `GetDocument(doc_id) → Document`
- `CreateDocument(owner_id, company_id, template_id) → Document`
- `UpdateDocument(doc_id, changes) → Document`
- `SendDocument(doc_id, recipient_ids) → Document`
- `ApproveDocument(doc_id, reviewer_id) → Document`
- `RejectDocument(doc_id, reason) → Document`
- `ArchiveDocument(doc_id) → Document`
- `ListDocuments(filter, page) → List<Document>`
- `AddDocumentEntry(doc_id, data) → DocumentEntry`
- `UpdateDocumentEntry(entry_id, data) → DocumentEntry`
- `RemoveDocumentEntry(entry_id) → Empty`

**Dependencies:** PostgreSQL, Auth-Service (gRPC), Company-Service (gRPC)

**Events subscribed:**

- `UserRegistered`
- `OrganizationCreated`

---

### Proto Files Location

**Все proto файлы находятся в:** `api/proto/`

```
api/proto/
├── auth_service.proto          # ✅ AuthService definition (6 methods)
├── auth.proto                  # ✅ User, Token, Credential, Role messages
├── company_service.proto       # ✅ CompanyService definition (8 methods)
├── company.proto               # ✅ Organization, Employee, Department messages
├── document_service.proto      # ✅ DocumentService definition (11 methods)
├── document.proto              # ✅ Document, DocumentEntry, Workflow messages
├── common.proto                # ✅ Empty, Error, PageInfo shared messages
├── Makefile                    # ✅ Proto compilation automation
└── README.md                   # Proto documentation
```

**Total:** 25 RPC методов, 3 service definitions

**Компиляция:**

```bash
cd api/proto
make proto              # Генерирует *.pb.go и *_grpc.pb.go
make proto-install-tools # Устанавливает protoc plugins
```

---

## 📅 ПЛАН РАЗРАБОТКИ (6 фаз)

### ✅ Phase 1: Proto Files Centralization

**Статус:** COMPLETED (Jan 1, 2026)

- ✅ Создана структура `api/proto/`
- ✅ Все proto файлы созданы (8 files, 25 RPC methods)
- ✅ Makefile для компиляции
- ✅ Документация обновлена

### 📅 Phase 2: Auth-Service Implementation

**Статус:** READY FOR DEVELOPMENT
**Время:** 3 недели
**Tasks:**

1. Создать модуль `auth-service/`
2. Implement DDD слои (domain, application, infrastructure, interfaces)
3. REST endpoints (register, login, logout, get-user)
4. gRPC handlers для других сервисов
5. Unit тесты (domain layer)
6. Integration тесты
7. Docker image + docker-compose

**Milestone:** Auth-Service работает со всеми endpoints

### 📅 Phase 3: Company-Service Implementation

**Время:** 2.5 недели
**Tasks:**

1. Создать модуль `company-service/`
2. Implement DDD слои
3. REST endpoints для организаций
4. gRPC handlers
5. Добавить RabbitMQ subscriptions (слушать UserRegisteredEvent)
6. gRPC client для Auth-Service
7. Тесты

### 📅 Phase 4: Document-Service Implementation

**Время:** 2.5 недели
**Аналогично Company-Service**

### 📅 Phase 5: API Gateway & Integration

**Время:** 2 недели
**Tasks:**

1. Создать API Gateway (Kong или custom)
2. gRPC load balancing
3. Circuit breaker patterns
4. Rate limiting
5. Полная интеграция всех сервисов

### 📅 Phase 6: DevOps & Deployment

**Время:** 2 недели
**Tasks:**

1. Kubernetes manifests
2. Helm charts
3. CI/CD pipeline
4. Monitoring (Prometheus, Grafana)
5. Logging (ELK)
6. Distributed tracing (Jaeger)

---

## 🛠️ ИНСТРУМЕНТЫ & ТЕХНОЛОГИИ

| Компонент            | Выбор            | Причина                            |
| -------------------- | ---------------- | ---------------------------------- |
| **Framework**        | Fiber            | Быстрый, минимальный overhead      |
| **Language**         | Go               | Производительность, concurrency    |
| **DB**               | PostgreSQL       | ACID, надежность, масштабируемость |
| **Cache**            | Redis            | Быстрый, поддерживает Pub/Sub      |
| **Message Queue**    | RabbitMQ         | Надежность, persistence, DLQ       |
| **gRPC**             | Protocol Buffers | Типизирован, быстрый               |
| **Containerization** | Docker           | Стандартный, простой               |
| **Orchestration**    | Kubernetes       | Production-ready                   |
| **Monitoring**       | Prometheus       | Metrics collection                 |
| **Logging**          | ELK Stack        | Centralized logging                |
| **Tracing**          | Jaeger           | Distributed tracing                |

---

## 🧪 TESTING STRATEGY

### 1. Unit Tests

- **Слой:** Domain Layer
- **Что тестировать:** Entities, Value Objects, Domain Logic
- **Как:** Без mock-объектов, чистые функции

```go
// domain/user_test.go
func TestRegisterUser(t *testing.T) {
    user := NewUser("john@example.com", "password123")
    assert.NoError(t, user.Register())
}
```

### 2. Integration Tests

- **Слой:** Application Layer
- **Что тестировать:** Use Cases + Repository + Database
- **Как:** testcontainers для PostgreSQL

```go
// application/services/register_user_service_test.go
func TestRegisterUserService_Integration(t *testing.T) {
    // Setup: testcontainers PostgreSQL
    // Execute: service.Register()
    // Verify: User created in DB
}
```

### 3. Contract Tests

- **Слой:** Inter-service Communication
- **Что тестировать:** gRPC contracts между сервисами
- **Как:** Pact or Protocol Buffer schema validation

```go
// Test that Company-Service can call Auth-Service
// Test that expected request/response matches proto
```

---

## 📊 DATABASE STRATEGY

### Текущий подход: Shared Database

```
┌──────────────────────────────────┐
│    PostgreSQL (Single Database)  │
├──────────────────────────────────┤
│ users_table                      │
│ organizations_table              │
│ documents_table                  │
│ ...                              │
└──────────────────────────────────┘
```

**Плюсы:** Простая интеграция, ACID транзакции  
**Минусы:** Сложнее масштабировать потом

### Будущий подход: Database Per Service

```
├─ auth_db      (PostgreSQL)
├─ company_db   (PostgreSQL)
└─ document_db  (PostgreSQL)
```

**Требует:** Event Sourcing, Saga pattern для распределенных транзакций

---

## 📊 DEPLOYMENT ARCHITECTURE

### Local Development

```bash
docker-compose -f docker-compose.microservices.yml up -d
# Запускает: PostgreSQL, Redis, RabbitMQ, Auth-Service, Company-Service, Document-Service
```

### Docker Compose (Current)

```yaml
services:
  postgres:
    image: postgres:15
    ports: 5432

  redis:
    image: redis:7
    ports: 6379

  rabbitmq:
    image: rabbitmq:3.11-management
    ports: 5672, 15672

  auth-service:
    build: ./auth-service
    ports: 3001

  company-service:
    build: ./company-service
    ports: 3002

  document-service:
    build: ./document-service
    ports: 3003
```

### Kubernetes (Future)

```
Namespace: tunduck
├─ auth-service (Deployment, Service, ConfigMap, Secret)
├─ company-service (Deployment, Service, ConfigMap, Secret)
├─ document-service (Deployment, Service, ConfigMap, Secret)
├─ api-gateway (Deployment, Service, Ingress)
├─ postgres (StatefulSet, PersistentVolume)
├─ redis (StatefulSet)
└─ rabbitmq (StatefulSet)
```

---

## 🚀 GETTING STARTED FOR DEVELOPERS

### Для разработки Auth-Service

```bash
# 1. Клонировать репо и перейти
git clone https://github.com/rusgainew/tunduck-app-mk.git
cd tunduck-app-mk

# 2. Запустить инфраструктуру
docker-compose -f docker-compose.microservices.yml up -d postgres redis rabbitmq

# 3. Перейти в auth-service
cd auth-service

# 4. Компилировать proto (если изменились)
cd ../api/proto && make proto && cd ../auth-service

# 5. Запустить сервис локально
go run cmd/auth-service/main.go

# 6. Тестировать
go test ./...

# 7. Тестировать API
curl -X POST http://localhost:3001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Для разработки новых Proto методов

```bash
# 1. Редактировать proto файл
vim api/proto/auth_service.proto

# 2. Скомпилировать
cd api/proto && make proto

# 3. Реализовать handler в сервисе
vim auth-service/internal/interfaces/grpc/handlers/auth_grpc_handler.go

# 4. Написать тесты
vim auth-service/internal/interfaces/grpc/handlers/auth_grpc_handler_test.go

# 5. Скомпилировать и запустить
cd auth-service && go test ./...
```

---

## ❓ FAQ & TROUBLESHOOTING

### Q: Как добавить новый RPC метод?

**A:**

1. Добавить в `api/proto/service_name.proto`
2. Запустить `cd api/proto && make proto`
3. Реализовать в `internal/interfaces/grpc/handlers/`

### Q: Как вызвать другой сервис из gRPC?

**A:**

```go
// internal/infrastructure/http/company_client.go
type CompanyClient struct {
    conn *grpc.ClientConn
}

func (c *CompanyClient) GetOrganization(ctx context.Context, orgID string) (*Organization, error) {
    req := &GetOrganizationRequest{OrganizationId: orgID}
    return c.client.GetOrganization(ctx, req)
}
```

### Q: Как подписаться на RabbitMQ событие?

**A:**

```go
// internal/infrastructure/events/rabbit_subscriber.go
func (s *RabbitSubscriber) SubscribeUserRegistered() error {
    ch, err := s.conn.Channel()
    q, err := ch.QueueDeclare("user.registered", ...)
    forever := make(chan amqp.Delivery)
    ch.Consume(q.Name, "", true, false, false, false, forever)

    for d := range forever {
        var event UserRegisteredEvent
        json.Unmarshal(d.Body, &event)
        // Handle event
    }
}
```

### Q: Как подключиться к PostgreSQL локально?

**A:**

```bash
# Из docker-compose
docker-compose -f docker-compose.microservices.yml exec postgres psql -U postgres -d tunduck
```

### Q: Как очистить docker контейнеры?

**A:**

```bash
docker-compose -f docker-compose.microservices.yml down -v  # -v удаляет volumes
docker system prune -a  # Удаляет неиспользуемые образы
```

---

## 📚 ДОПОЛНИТЕЛЬНЫЕ ДОКУМЕНТЫ

Для специализированных тем смотрите отдельные документы:

| Документ                           | Для чего                      |
| ---------------------------------- | ----------------------------- |
| `api/proto/README.md`              | Детали работы с proto файлами |
| `auth-service/README.md`           | Специфика Auth-Service        |
| `auth-service/GRPC_SETUP.md`       | Настройка gRPC                |
| `go-api/docs/API_DOCUMENTATION.md` | Текущий API (monolith)        |
| `PHASE2_AUTH_SERVICE_PLAN.md`      | Детальный план Phase 2        |
| `PHASE1_COMPLETION_REPORT.md`      | Статус Phase 1                |

---

## ✅ ЧЕК-ЛИСТ ДЛЯ НОВИЧКОВ

Когда начинаете разработку:

- [ ] Прочитал этот документ (PROJECT_MASTER_GUIDE.md)
- [ ] Запустил docker-compose локально
- [ ] Скомпилировал proto файлы
- [ ] Запустил auth-service локально
- [ ] Написал простой unit test
- [ ] Прочитал api/proto/README.md
- [ ] Понял DDD слои
- [ ] Готов к разработке Phase 2

---

## 📞 КОНТАКТЫ

**Repository:** https://github.com/rusgainew/tunduck-app-mk  
**Current Branch:** `dev`  
**Default Branch:** `main`  
**Issues:** https://github.com/rusgainew/tunduck-app-mk/issues

**Team Members:** [Добавить контакты]

---

**Last Updated:** 1 января 2026  
**Next Review:** 15 января 2026 (или после Phase 2 начала)
