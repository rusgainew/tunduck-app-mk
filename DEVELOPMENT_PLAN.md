# 📈 ПЛАН РАЗРАБОТКИ ПРИЛОЖЕНИЯ - TUNDUCK

**Дата создания:** 1 января 2026  
**Версия:** 1.0  
**Статус:** Ready for Implementation

---

## 📊 АНАЛИЗ ТЕКУЩЕГО СОСТОЯНИЯ

### Что уже сделано ✅

#### Phase 1: Proto Files Centralization (ЗАВЕРШЕНА)

- ✅ Централизованная структура `api/proto/` для всех сервисов
- ✅ 3 Service definitions (auth, company, document) - 25 RPC методов
- ✅ 4 Message definition файла (auth, company, document, common)
- ✅ Makefile для компиляции proto
- ✅ Полная документация

#### Текущая инфраструктура

- ✅ Docker & Docker Compose настроены
- ✅ PostgreSQL для БД
- ✅ Redis для кеша
- ✅ RabbitMQ для событий
- ✅ CI/CD pipeline

#### Документация

- ✅ PROJECT_MASTER_GUIDE.md - полный гайд (23 KB)
- ✅ Детальные планы для каждой фазы
- ✅ DDD паттерны задокументированы
- ✅ Примеры кода для каждого слоя

### Что нужно сделать 🔨

#### Phase 2: Auth-Service (2-3 недели) - NEXT

- 🔲 Создать модуль auth-service/
- 🔲 Implement 4 DDD слоя
- 🔲 REST endpoints (register, login, logout)
- 🔲 gRPC handlers
- 🔲 Unit & Integration тесты
- 🔲 Docker image

#### Phase 3: Company-Service (2.5 недели)

- 🔲 Создать модуль company-service/
- 🔲 Implement DDD слои
- 🔲 REST endpoints
- 🔲 gRPC + RabbitMQ интеграция
- 🔲 Тесты

#### Phase 4: Document-Service (2.5 недели)

- 🔲 Создать модуль document-service/
- 🔲 Implement DDD слои
- 🔲 REST endpoints
- 🔲 Интеграция с Auth & Company
- 🔲 Тесты

#### Phase 5: API Gateway & Integration (2 недели)

- 🔲 Создать API Gateway
- 🔲 gRPC load balancing
- 🔲 Circuit breaker patterns
- 🔲 Rate limiting

#### Phase 6: DevOps & Deployment (2 недели)

- 🔲 Kubernetes manifests
- 🔲 Helm charts
- 🔲 Monitoring (Prometheus, Grafana)
- 🔲 Logging (ELK)
- 🔲 Distributed tracing (Jaeger)

---

## 📅 ДЕТАЛЬНЫЙ ПЛАН РАЗРАБОТКИ

### 🎯 PHASE 2: AUTH-SERVICE (Недели 1-3)

#### Неделя 1: Setup & Domain Layer

**День 1-2: Проект структура**

```bash
mkdir auth-service
cd auth-service

# Структура
mkdir -p cmd internal/{domain,application,infrastructure,interfaces}
mkdir -p internal/infrastructure/{persistence,rabbitmq,http,config}
mkdir -p internal/application/{services,dto,events}
mkdir -p internal/interfaces/{http,grpc}
mkdir -p migrations tests/{domain,application,integration}
mkdir -p docs docker
```

**День 3-5: Domain Layer (чистая бизнес-логика)**

Файлы для создания:

```
internal/domain/
├── user.go                    # User aggregate root
│   ├── User struct (id, email, password, roles, status)
│   ├── Methods: Register(), Login(), ChangePassword(), Logout()
│   └── Value Objects: UserID, Email, Password, Token
│
├── repository.go              # Repository interface
│   └── UserRepository interface (Save, FindByEmail, FindByID, Delete)
│
├── errors.go                  # Domain errors
│   ├── ErrInvalidEmail
│   ├── ErrPasswordTooShort
│   ├── ErrUserNotFound
│   └── ErrInvalidCredentials
│
└── events.go                  # Domain events
    ├── UserRegistered event
    ├── UserLoggedIn event
    └── UserLoggedOut event
```

**Код для User aggregate:**

```go
package domain

type User struct {
    ID        string
    Email     Email
    Password  Password
    FirstName string
    LastName  string
    Roles     []string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (u *User) Register(email Email, password Password) error {
    if !email.IsValid() {
        return ErrInvalidEmail
    }
    if !password.IsStrong() {
        return ErrWeakPassword
    }
    u.Email = email
    u.Password = password
    u.Status = "active"
    return nil
}

func (u *User) VerifyPassword(pwd Password) bool {
    return u.Password.Match(pwd)
}
```

**День 6-7: Unit Tests для Domain Layer**

- TestUserRegister
- TestPasswordValidation
- TestEmailValidation

---

#### Неделя 2: Application & Infrastructure Layers

**День 8-10: Application Layer (Use Cases)**

Файлы:

```
internal/application/services/
├── register_user_service.go
│   ├── Execute(ctx, RegisterUserDTO) (*User, error)
│   ├── Валидация email
│   ├── Хеширование пароля
│   ├── Сохранение в БД
│   └── Публикация UserRegistered event
│
├── login_user_service.go
│   ├── FindUser by email
│   ├── Verify password
│   ├── Generate JWT token
│   └── Publish UserLoggedIn event
│
├── validate_token_service.go
│   └── Validate JWT token
│
└── logout_user_service.go
    └── Revoke token in Redis

internal/application/dto/
├── register_user_dto.go       # RegisterUserRequest
├── login_user_dto.go          # LoginRequest
├── user_response_dto.go       # UserResponse
└── token_response_dto.go      # TokenResponse
```

**День 11-12: Infrastructure Layer**

Файлы:

```
internal/infrastructure/persistence/
├── user_repository.go         # Interface implementation
└── postgres/
    ├── user_postgres_repo.go
    │   ├── Save(user) error
    │   ├── FindByEmail(email) (*User, error)
    │   ├── FindByID(id) (*User, error)
    │   └── Delete(id) error
    └── queries.go             # SQL queries

internal/infrastructure/rabbitmq/
├── event_publisher.go         # Publish events
│   ├── PublishUserRegistered()
│   ├── PublishUserLoggedIn()
│   └── PublishUserLoggedOut()
└── connection.go              # RabbitMQ connection

internal/infrastructure/config/
└── config.go                  # Load env vars
```

**День 13-14: Integration Tests**

- TestRepository
- TestEventPublishing
- TestDatabaseConnection

---

#### Неделя 3: Interfaces & Testing

**День 15-17: Interfaces Layer (API)**

REST Handlers:

```
internal/interfaces/http/handlers/
├── register_handler.go        # POST /api/v1/auth/register
│   ├── Parse request JSON
│   ├── Call RegisterUserService
│   └── Return UserResponse
│
├── login_handler.go           # POST /api/v1/auth/login
│   ├── Validate credentials
│   ├── Call LoginUserService
│   └── Return TokenResponse
│
├── me_handler.go              # GET /api/v1/auth/me
│   ├── Extract token from header
│   ├── Validate token
│   └── Return current user
│
├── logout_handler.go          # POST /api/v1/auth/logout
│   ├── Extract token
│   ├── Call LogoutUserService
│   └── Return success
│
└── refresh_handler.go         # POST /api/v1/auth/refresh
    ├── Validate old token
    └── Generate new token

internal/interfaces/http/
├── routes.go                  # Route definition
│   ├── /api/v1/auth/register
│   ├── /api/v1/auth/login
│   ├── /api/v1/auth/me
│   ├── /api/v1/auth/logout
│   └── /api/v1/auth/refresh
│
├── middleware/
│   ├── auth_middleware.go     # JWT validation
│   ├── error_handler.go       # Error responses
│   └── logging_middleware.go  # Request logging
└── openapi.yaml               # Swagger specification
```

gRPC Handlers:

```
internal/interfaces/grpc/
├── auth_grpc_handler.go       # AuthService implementation
│   ├── Register() gRPC
│   ├── Login() gRPC
│   ├── ValidateToken() gRPC
│   ├── GetUser() gRPC
│   ├── Logout() gRPC
│   └── RefreshToken() gRPC
└── interceptors/
    └── logging_interceptor.go
```

**День 18-19: HTTP & gRPC Tests**

- TestRegisterEndpoint
- TestLoginEndpoint
- TestAuthMiddleware
- TestgRPCHandlers

**День 20-21: Docker & Documentation**

- Dockerfile (multi-stage build)
- docker-compose.yml (local dev)
- README.md
- API.md (REST documentation)
- GRPC.md (gRPC documentation)

---

### 📋 PHASE 2 DELIVERABLES

#### Code Structure

```
auth-service/
├── ✅ Domain Layer (100% domain logic)
├── ✅ Application Layer (6 use cases)
├── ✅ Infrastructure Layer (PostgreSQL + RabbitMQ)
├── ✅ Interfaces Layer (REST + gRPC)
├── ✅ DI Container (wire or manual)
└── ✅ All tests (80%+ coverage)
```

#### RPC Methods Implemented

1. ✅ Register(RegisterRequest) → AuthResponse
2. ✅ Login(LoginRequest) → AuthResponse
3. ✅ ValidateToken(ValidateTokenRequest) → User
4. ✅ GetUser(GetUserRequest) → User
5. ✅ Logout(LogoutRequest) → Empty
6. ✅ RefreshToken(RefreshTokenRequest) → Token

#### Tests

- 🧪 Domain layer: 15+ unit tests
- 🧪 Application layer: 20+ unit tests
- 🧪 Infrastructure: 10+ integration tests
- 🧪 HTTP Handlers: 15+ integration tests
- 🧪 gRPC Handlers: 10+ tests
- **Total:** 70+ tests with 80% code coverage

#### Documentation

- 📖 README.md
- 📖 API.md (REST documentation)
- 📖 GRPC.md (gRPC documentation)
- 📖 ARCHITECTURE.md (service design)
- 📖 EVENTS.md (RabbitMQ events guide)

---

### 🎯 PHASE 3: COMPANY-SERVICE (Недели 4-6)

**Аналогично Phase 2, но с:**

- Organization aggregate root
- Employee, Department entities
- gRPC client для Auth-Service
- RabbitMQ subscriber (слушает UserRegisteredEvent)
- 8 RPC методов CompanyService

**Timeline:**

- Неделя 4: Setup + Domain + Application
- Неделя 5: Infrastructure + Interfaces
- Неделя 6: Tests + Docker + Documentation

---

### 🎯 PHASE 4: DOCUMENT-SERVICE (Недели 7-9)

**Аналогично Phase 3, но с:**

- Document aggregate root
- DocumentEntry, DocumentWorkflow entities
- gRPC clients для Auth & Company Services
- RabbitMQ subscriber (слушает все события)
- 11 RPC методов DocumentService
- Complex workflow logic

**Timeline:**

- Неделя 7: Setup + Domain + Application
- Неделя 8: Infrastructure + Interfaces + Integration
- Неделя 9: Tests + Docker + Documentation

---

### 🎯 PHASE 5: API GATEWAY (Недели 10-11)

**Задачи:**

1. Создать API Gateway (Kong или custom Fiber)
2. Route mapping:
   - `/api/auth/*` → auth-service:3001
   - `/api/company/*` → company-service:3002
   - `/api/document/*` → document-service:3003
3. Rate limiting (per IP, per user)
4. Circuit breaker для gRPC вызовов
5. gRPC load balancing
6. Request/response logging
7. Distributed tracing (Jaeger)

**Technologies:**

- Kong / Envoy / Custom Fiber
- gRPC load balancer
- Prometheus для метрик

---

### 🎯 PHASE 6: DEVOPS & MONITORING (Недели 12-13)

**Kubernetes:**

- Service manifests для каждого микросервиса
- Ingress для API Gateway
- ConfigMaps для конфигурации
- Secrets для sensitive data
- StatefulSets для БД/Cache/MQ

**Monitoring:**

- Prometheus (metrics collection)
- Grafana (dashboards)
- ELK Stack (logging)
- Jaeger (distributed tracing)

**CI/CD:**

- GitHub Actions workflows
- Unit test automation
- Integration test automation
- Docker image building & pushing
- Deployment automation

---

## 🛠️ ТЕХНИЧЕСКИЙ СТЕК

| Компонент     | Выбор              | Версия | Причина                            |
| ------------- | ------------------ | ------ | ---------------------------------- |
| Language      | Go                 | 1.25+  | Производительность, concurrency    |
| Framework     | Fiber              | Latest | Быстрый, минимальный overhead      |
| DB            | PostgreSQL         | 15+    | ACID, надежность, масштабируемость |
| Cache         | Redis              | 7+     | Быстрый, поддерживает Pub/Sub      |
| Message Queue | RabbitMQ           | 3.11+  | Надежность, persistence, DLQ       |
| gRPC          | Protocol Buffers 3 | Latest | Типизирован, быстрый               |
| Testing       | testify            | Latest | Assertions, mocking, test suites   |
| DI            | wire / fx          | Latest | Compile-time or runtime DI         |
| Logging       | logrus / zap       | Latest | Structured logging                 |
| Metrics       | Prometheus         | Latest | Metrics collection                 |
| Tracing       | Jaeger             | Latest | Distributed tracing                |
| Container     | Docker             | 20.10+ | Стандартный, простой               |
| Orchestration | Kubernetes         | 1.25+  | Production-ready                   |
| API Docs      | Swagger/OpenAPI    | 3.0    | Интерактивная документация         |

---

## 📊 РИСКИ И МИTIGATIONS

| Риск                      | Вероятность | Влияние | Mitigation                              |
| ------------------------- | ----------- | ------- | --------------------------------------- |
| Сложность DDD             | Средняя     | Высокое | Обучение команды, паттерны из примеров  |
| Несовместимость gRPC      | Низкая      | Высокое | Четкое определение proto контрактов     |
| Performance bottleneck    | Низкая      | Высокое | Load testing каждой фазы                |
| Database schema conflicts | Средняя     | Среднее | Использование миграций, версионирование |
| Service discovery issues  | Низкая      | Высокое | Kubernetes service mesh (Istio)         |
| Monitoring complexity     | Средняя     | Среднее | Готовые dashboards в Grafana            |

---

## 📈 МЕТРИКИ УСПЕХА

### Phase 2 (Auth-Service)

- ✅ 6 RPC методов работают корректно
- ✅ 70+ тестов с 80% coverage
- ✅ API документирована в Swagger
- ✅ Docker image собирается и запускается
- ✅ Integration тесты с реальной БД проходят

### Phase 3 (Company-Service)

- ✅ 8 RPC методов работают корректно
- ✅ Интеграция с Auth-Service через gRPC
- ✅ RabbitMQ events публикуются и получаются
- ✅ 70+ тестов
- ✅ Performance: <100ms для основных операций

### Phase 4 (Document-Service)

- ✅ 11 RPC методов работают корректно
- ✅ Интеграция с Auth & Company Services
- ✅ Complex workflow engine работает
- ✅ 70+ тестов
- ✅ Все сервисы масштабируются горизонтально

### Phase 5 (API Gateway)

- ✅ Все requests маршрутизируются корректно
- ✅ Rate limiting работает
- ✅ Circuit breaker срабатывает на ошибки
- ✅ Requests логируются для трейсинга

### Phase 6 (DevOps)

- ✅ Kubernetes manifests готовы
- ✅ Prometheus собирает метрики
- ✅ Grafana показывает dashboards
- ✅ ELK Stack собирает логи
- ✅ Jaeger трассирует распределенные вызовы

---

## 🎯 CHECKPOINTS

### End of Phase 2 (День 21)

- [ ] Auth-Service полностью реализована
- [ ] Все REST endpoints работают
- [ ] Все gRPC методы реализованы
- [ ] 70+ тестов проходят
- [ ] Docker image готов
- [ ] Документация полная

### End of Phase 3 (День 42)

- [ ] Company-Service реализована
- [ ] Integration с Auth-Service работает
- [ ] RabbitMQ events отправляются/получаются
- [ ] Все 16 микросервис RPC методов готовы (Auth 6 + Company 8 + 2 новых)
- [ ] 140+ тестов проходят

### End of Phase 4 (День 63)

- [ ] Document-Service реализована
- [ ] Все 25 RPC методов готовы
- [ ] Workflow engine работает
- [ ] 200+ тестов проходят
- [ ] Все 3 сервиса интегрированы

### End of Phase 5 (День 77)

- [ ] API Gateway маршрутизирует все запросы
- [ ] Rate limiting работает
- [ ] Circuit breaker работает
- [ ] Performance <100ms

### End of Phase 6 (День 91)

- [ ] Kubernetes deployment готов
- [ ] Мониторинг полностью настроен
- [ ] Логирование централизовано
- [ ] Трейсинг работает

---

## 🚀 БЫСТРЫЙ СТАРТ

### Сегодня (День 1)

1. Прочитать PROJECT_MASTER_GUIDE.md
2. Обсудить с командой архитектурные решения
3. Выделить ресурсы для Phase 2

### Завтра (День 2)

1. Создать auth-service структуру
2. Начать реализацию Domain Layer
3. Написать первые unit тесты

### Неделя 1

1. Завершить Domain + Application Layer
2. Начать Infrastructure
3. Написать unit тесты

### Неделя 2

1. Завершить Infrastructure + Interfaces
2. Написать integration тесты
3. Создать HTTP handlers

### Неделя 3

1. Написать gRPC handlers
2. Создать Docker image
3. Написать документацию

---

## 📚 ДОПОЛНИТЕЛЬНЫЕ РЕСУРСЫ

**В проекте:**

- `PROJECT_MASTER_GUIDE.md` - полный гайд архитектуры
- `PHASE2_AUTH_SERVICE_PLAN.md` - детальный план Phase 2
- `api/proto/README.md` - работа с proto файлами
- `generate-service.sh` - scaffold для новых сервисов

**Внешние ресурсы:**

- Domain-Driven Design (Eric Evans)
- Building Microservices with Go (Nic Jackson)
- gRPC Best Practices (официальная документация)

**Инструменты:**

- `wire` для dependency injection
- `testify` для testing
- `sqlc` для type-safe SQL
- `protobuf` для gRPC

---

## 📞 ВОПРОСЫ

**Q: Сколько в сумме времени займет разработка?**  
A: 13 недель (3+ месяца) при команде из 2-3 разработчиков

**Q: Можно ли начать с другого микросервиса?**  
A: Нет, Auth-Service - foundation для остальных

**Q: Какой процент кода переиспользуется из monolith?**  
A: ~30-40% (бизнес-логика), остальное нужно рефакторить под DDD

**Q: Нужны ли гарантии обратной совместимости?**  
A: Да, при каждом изменении proto файла нужна миграция

**Q: Как обновлять proto файлы?**  
A: Редактировать `api/proto/*.proto` → `make proto` → обновить сервисы

---

**Статус:** ✅ Ready for Implementation  
**Дата обновления:** 1 января 2026  
**Версия:** 1.0
