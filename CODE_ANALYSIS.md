# Детальный код-анализ текущей архитектуры

## 📈 Метрики кода

### По слоям

| Слой        | Файлы      | Строк       | Комментарий                           |
| ----------- | ---------- | ----------- | ------------------------------------- |
| Controllers | 5          | ~1,328      | HTTP handlers для всех доменов        |
| Services    | 5 + 5 impl | ~2,000      | Бизнес-логика + интеграционные тесты  |
| Models      | 4          | ~?          | Domain entities (User, Org, Document) |
| Repository  | 8          | ~?          | Data access layer                     |
| **ИТОГО**   | **27**     | **~11,282** | **Monolith**                          |

### Топ 5 файлов (по размеру)

1. `services/service_impl/services_comprehensive_test.go` - 640 строк (тесты)
2. `services/service_impl/user_service_impl.go` - 454 строк
3. `services/service_impl/esf_document_service_impl.go` - 362 строк
4. `controllers/auth_controller_test.go` - 354 строк (тесты)
5. `controllers/esf_document_controller.go` - 321 строк

---

## 🏗️ Текущая архитектура (Monolith)

### Controllers (HTTP Layer)

```
✓ auth_controller.go (280 LOC)
  - POST /register - регистрация
  - POST /login - вход
  - GET /me - профиль
  - POST /logout - выход

✓ user_controller.go (207 LOC)
  - CRUD пользователя
  - Управление правами

✓ role_controller.go (285 LOC)
  - Управление ролями
  - Назначение прав

✓ esf_organization_controller.go (235 LOC)
  - Управление компаниями
  - Структура организации

✓ esf_document_controller.go (321 LOC)
  - Управление документами
  - Документооборот
```

### Services (Business Logic)

```
Interfaces:
  - UserService (21 LOC)
  - RoleService (29 LOC)
  - ESFDocumentService (25 LOC)
  - ESFOrganizationService (25 LOC)

Implementations:
  - UserServiceImpl (454 LOC) ← LARGEST
  - RoleServiceImpl (241 LOC)
  - ESFDocumentServiceImpl (362 LOC)
  - ESFOrganizationServiceImpl (284 LOC)
  - OrganizationDBServiceImpl (159 LOC)

Tests:
  - services_comprehensive_test.go (640 LOC)
  - user_service_integration_test.go (182 LOC)
  - esf_organization_service_test.go (21 LOC)
```

### Models (Entities)

```
Вероятно:
  - user_model.go - User aggregate
  - esf_organization_model.go - Organization aggregate
  - esf_document_model.go - Document aggregate
  - esf_entries_model.go - Document entries
```

### Repository (Data Access)

```
Interfaces & Implementations:
  - user_repository.go / user_repository_postgres.go
  - esf_document_repository.go / esf_document_postgres.go
  - esf_organization.go / esf_organization_postgres.go
```

---

## 🔗 Зависимости между компонентами

```
Controller Layer
       ↓
   ↓---┴---↓---↓
User   Role   Document   Organization
Service Layer
       ↓
   Repository Layer (PostgreSQL)
       ↓
   Shared Layer (Cache, Logger, Auth)
```

### Проблемы текущей архитектуры

1. **Tight Coupling** - Controllers напрямую зависят от Services
2. **No Bounded Contexts** - все вместе в одном модуле
3. **Mixed Concerns** - бизнес-логика, HTTP, database в одном месте
4. **Hard to Scale** - невозможно развертывать отдельно
5. **Testing Complexity** - нужны мок-объекты для всего
6. **No Event System** - жесткие синхронные вызовы

---

## 🎯 Доменные области для разделения

### 1. AUTH & REGISTRATION Domain

**Файлы для миграции:**

- `controllers/auth_controller.go`
- `services/user_service.py` (частично)
- `services/service_impl/user_service_impl.go` (частично)
- `models/user_model.go` (User entity)
- `repository/user_repository*` (User persistence)

**Entities:**

- User (aggregate root)
- Credential (value object)
- Token (value object)

**Use Cases:**

- Register(email, password)
- Login(email, password)
- RefreshToken()
- Logout()
- ChangePassword()

---

### 2. COMPANY/ORGANIZATION Domain

**Файлы для миграции:**

- `controllers/esf_organization_controller.go`
- `services/esf_organization_service*`
- `models/esf_organization_model.go`
- `repository/esf_organization*`

**Entities:**

- Organization (aggregate root)
- Employee
- Department

**Use Cases:**

- CreateOrganization()
- UpdateOrganization()
- ManageEmployees()

---

### 3. DOCUMENT Domain

**Файлы для миграции:**

- `controllers/esf_document_controller.go`
- `services/esf_document_service*`
- `models/esf_document_model.go`
- `models/esf_entries_model.go`
- `repository/esf_document*`

**Entities:**

- Document (aggregate root)
- DocumentEntry (nested entity)

**Use Cases:**

- CreateDocument()
- SendDocument()
- ReceiveDocument()
- ApproveDocument()

---

### 4. RBAC (Role-Based Access Control)

**Файлы для миграции:**

- `controllers/role_controller.go`
- `services/role_service*`

**Note:** Может быть shared service или часть каждого сервиса

---

## 📦 Shared/Cross-Cutting Concerns

### Текущие pkg utilities:

```
pkg/
  ├── apperror/       - Error handling
  ├── auth/           - JWT, Auth
  ├── cache/          - Redis caching
  ├── container/      - DI container
  ├── entity/         - Base entities
  ├── health/         - Health checks
  ├── logger/         - Logging (Logrus)
  ├── metrics/        - Prometheus
  ├── middleware/     - HTTP middleware (JWT, etc)
  ├── pagination/     - Pagination helpers
  ├── ratelimit/      - Rate limiting
  ├── rbac/           - RBAC
  ├── response/       - HTTP response wrappers
  └── transaction/    - Transaction management
```

### Данные должны остаться Shared:

✓ logger - общий для всех сервисов
✓ cache - Redis может быть shared
✓ metrics - Prometheus registry
✓ middleware - базовые middleware
✓ health - health checks
✓ errors - error types

### Данные, которые должны быть частью каждого сервиса:

✓ domain entities
✓ repositories
✓ services (use cases)
✓ handlers (controllers)
✓ dto (transfer objects)

---

## 🗄️ Стратегия Database

### Вариант 1: Database Per Service (рекомендуется для масштабирования)

```
PostgreSQL Cluster
  ├── auth_db (для auth-service)
  ├── company_db (для company-service)
  ├── document_db (для document-service)
  └── shared_db (для cross-service data)
```

### Вариант 2: Shared Database (проще для начала)

```
PostgreSQL
  ├── users (auth-service читает, пишет)
  ├── organizations
  ├── documents
  ├── roles
  └── (все в одной базе, но логически разделено)
```

**Рекомендация:** Начать с Вариант 2, потом перейти на Вариант 1

---

## 🔄 Миграция - примеры

### BEFORE: User Service (Monolith)

```go
// internal/services/service_impl/user_service_impl.go
type UserServiceImpl struct {
    repo   repository.UserRepository
    log    *logrus.Logger
    cache  cache.CacheManager
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
    // Проверяем кэш
    if cached, _ := s.cache.Get(...); cached != nil {
        return cached.(*models.User), nil
    }

    // Идем в БД
    user, err := s.repo.GetByID(ctx, id)
    if err != nil { ... }

    // Кэшируем
    s.cache.Set(...)

    return user, nil
}
```

### AFTER: Auth Service (DDD Microservice)

```
auth-service/
├── cmd/main.go
├── internal/
│   ├── domain/
│   │   ├── user.go (User aggregate)
│   │   ├── credential.go (Value Object)
│   │   └── errors.go
│   ├── application/
│   │   ├── services/
│   │   │   └── get_user_service.go
│   │   └── dto/
│   │       └── user_dto.go
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   ├── user_repository.go (interface)
│   │   │   └── postgres/user_postgres_repo.go
│   │   └── config/
│   ├── interfaces/
│   │   └── http/
│   │       └── handlers/
│   │           └── get_user_handler.go
│   └── container.go (DI)
└── go.mod
```

---

## 🚀 Next Steps

1. **Утвердить Decision:**

   - [ ] Database per service или Shared?
   - [ ] REST между сервисами или gRPC?
   - [ ] Event Bus (Redis Pub/Sub)?

2. **Создать шаблон микросервиса (scaffold)**

3. **Мигрировать AUTH-SERVICE первым** (самый простой)

4. **Мигрировать остальные** по очереди
