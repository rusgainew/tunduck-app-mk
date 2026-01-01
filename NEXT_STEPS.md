# Анализ кода и план продолжения разработки

## 📊 Анализ текущего кода

### Auth-Service ✅ (Завершен - Phase 2)

**Структура:** DDD 4-layer architecture

- ✅ Domain Layer: User aggregate, value objects, repository interfaces
- ✅ Application Layer: 6 use case services (Register, Login, Token, GetUser, Logout, RefreshToken)
- ✅ Infrastructure Layer: PostgreSQL, Redis, RabbitMQ, JWT manager
- ✅ Interfaces Layer: HTTP (:8001), gRPC (:9001)

**Метрики качества:**

- 19 файлов, ~1500 строк кода
- Unit тесты для domain и application слоев
- Чистая архитектура без внешних зависимостей в domain
- Bootstrap в main.go с graceful shutdown

**Готовность:** 100% - Production ready

---

### Go-API 🔄 (Требует рефакторинга)

**Текущая структура:** Монолитная с элементами Clean Architecture

**Проблемы:**

1. **Дублирование логики auth** - имеет собственные Register/Login сервисы
2. **Смешанная ответственность** - управляет и пользователями, и ESF документами
3. **Нет gRPC клиента** - не может общаться с auth-service
4. **Entity дублирование** - User существует в 2 местах
5. **Неполное разделение слоев** - сервисы вызывают GORM напрямую

**Что работает хорошо:**

- ✅ Middleware stack (Recovery, CORS, Logging, Metrics)
- ✅ Dependency Injection через контейнер
- ✅ Redis caching
- ✅ Health checks
- ✅ Prometheus metrics
- ✅ Rate limiting

**Компоненты для сохранения:**

```
go-api/pkg/
├── middleware/      # ✅ Оставить: хорошие middleware
├── cache/           # ✅ Оставить: Redis cache manager
├── response/        # ✅ Оставить: unified response format
├── apperror/        # ✅ Оставить: error handling
├── health/          # ✅ Оставить: health checker
├── metrics/         # ✅ Оставить: Prometheus metrics
└── container/       # ✅ Оставить: DI container
```

**Компоненты для рефакторинга:**

```
go-api/
├── internal/services/user_service.go         # ❌ Удалить Register/Login
├── internal/controllers/auth_controller.go   # 🔄 Переделать в proxy
├── pkg/entity/user_entity.go                 # 🔄 Упростить до UserReference
└── pkg/auth/auth.go                          # 🔄 Использовать только для middleware
```

---

## 🎯 Выполненные задачи

### ✅ Фаза 1: gRPC Client (Завершена)

1. **Создан gRPC клиент для auth-service**

   - Файл: `go-api/internal/clients/grpc/auth_client.go`
   - Методы: Register, Login, ValidateToken, RefreshToken, Logout, GetUser
   - Connection pooling с keepalive
   - Proper error handling

2. **Создан план рефакторинга**
   - Файл: `REFACTORING_GO_API_PLAN.md`
   - 5 фаз рефакторинга
   - Архитектурная диаграмма
   - Чеклист задач

---

## 🚀 Следующие шаги

### Приоритет 1: Завершение интеграции Auth-Service

#### Задача 1.1: Auth Proxy Service

Создать сервис-прокси для переадресации auth запросов на auth-service

```bash
# Создать файлы:
go-api/internal/services/auth_proxy_service.go       # Interface + Implementation
go-api/internal/services/service_impl/auth_proxy_service_impl.go
```

**Что делает:**

- Получает HTTP запросы от клиента
- Конвертирует в gRPC запросы
- Вызывает auth-service через gRPC клиент
- Конвертирует gRPC ответы обратно в HTTP

#### Задача 1.2: Обновить Auth Controller

Переделать `auth_controller.go` для использования AuthProxyService

```go
// ❌ БЫЛО:
userService.Register()  // вызов локального сервиса

// ✅ СТАНЕТ:
authProxyService.Register()  // вызов auth-service через gRPC
```

#### Задача 1.3: Обновить Auth Middleware

Middleware должен проверять токены через auth-service

```go
// В jwt_auth_middleware.go:
// Вместо локальной проверки JWT, вызывать:
user, err := authClient.ValidateToken(ctx, token)
```

---

### Приоритет 2: Разделение User Entity

#### Задача 2.1: Создать UserReference

Легковесная версия User для ESF документов

```go
// go-api/internal/domain/entity/user_reference.go
type UserReference struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
}
```

#### Задача 2.2: Удалить auth поля из User

```go
// Удалить из go-api/pkg/entity/user_entity.go:
// - Password
// - IsActive
// - Все auth-related методы
```

---

### Приоритет 3: Настройка Docker Compose

#### Задача 3.1: Создать docker-compose.microservices.yml

```yaml
services:
  auth-service:
    ports: ["8001:8001", "9001:9001"]

  go-api:
    environment:
      - AUTH_SERVICE_GRPC_URL=auth-service:9001
    depends_on:
      - auth-service
```

#### Задача 3.2: Обновить .env файлы

```bash
# go-api/.env
AUTH_SERVICE_GRPC_URL=localhost:9001  # для development
AUTH_SERVICE_HTTP_URL=localhost:8001  # для fallback

# auth-service/.env
JWT_SECRET=shared-secret-key
```

---

## 📝 Детальный план выполнения

### День 1: Auth Proxy Service

**Шаг 1:** Создать интерфейс

```bash
touch go-api/internal/services/auth_proxy_service.go
```

```go
package services

type AuthProxyService interface {
    Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
    Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
    ValidateToken(token string) (*models.UserInfo, error)
    RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error)
    Logout(ctx context.Context, token string) error
}
```

**Шаг 2:** Создать имплементацию

```bash
touch go-api/internal/services/service_impl/auth_proxy_service_impl.go
```

```go
package service_impl

import (
    grpcClient "github.com/rusgainew/tunduck-app-mk/go-api/internal/clients/grpc"
)

type authProxyService struct {
    authClient *grpcClient.AuthClient
    logger     *logrus.Logger
}

func NewAuthProxyService(authClient *grpcClient.AuthClient, logger *logrus.Logger) services.AuthProxyService {
    return &authProxyService{
        authClient: authClient,
        logger:     logger,
    }
}

func (s *authProxyService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
    // 1. Валидация входных данных
    // 2. Вызов auth-service через gRPC
    grpcResp, err := s.authClient.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName)
    // 3. Конвертация gRPC response → HTTP response
    // 4. Возврат результата
}
```

**Шаг 3:** Обновить Container

```go
// go-api/pkg/container/container.go
func (c *Container) GetAuthProxyService() services.AuthProxyService {
    if c.authProxyService == nil {
        authClient, _ := grpc.NewAuthClient(os.Getenv("AUTH_SERVICE_GRPC_URL"))
        c.authProxyService = service_impl.NewAuthProxyService(authClient, c.logger)
    }
    return c.authProxyService
}
```

**Шаг 4:** Обновить Auth Controller

```go
// go-api/internal/controllers/auth_controller.go
func (ctrl *AuthController) Register(c *fiber.Ctx) error {
    var req models.RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return err
    }

    // ВМЕСТО: resp, err := ctrl.userService.Register(...)
    // НОВОЕ:
    resp, err := ctrl.authProxyService.Register(c.Context(), &req)

    return c.JSON(resp)
}
```

---

### День 2: Middleware + Entity

**Задачи:**

1. Обновить `jwt_auth_middleware.go` для вызова auth-service
2. Создать `UserReference` entity
3. Удалить auth логику из `user_entity.go`
4. Обновить ESF controllers для использования UserReference

---

### День 3: Docker + Тестирование

**Задачи:**

1. Создать `docker-compose.microservices.yml`
2. Обновить `.env.example` файлы
3. Написать integration тесты
4. Протестировать весь flow: Register → Login → ValidateToken
5. Проверить health checks обоих сервисов

---

## 🧪 Тестирование

### Юнит-тесты для AuthProxyService

```go
// auth_proxy_service_test.go
func TestRegister(t *testing.T) {
    // Mock gRPC client
    mockClient := &MockAuthClient{}
    service := NewAuthProxyService(mockClient, logger)

    // Test success case
    // Test validation errors
    // Test gRPC errors
}
```

### Интеграционные тесты

```bash
# 1. Запустить auth-service
cd auth-service && go run cmd/auth-service/main.go

# 2. Запустить go-api
cd go-api && go run cmd/api/main.go

# 3. Тест регистрации
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!","firstName":"John","lastName":"Doe"}'

# 4. Тест логина
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test123!"}'

# 5. Тест protected endpoint
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer <token>"
```

---

## 📊 Метрики прогресса

### Текущее состояние (1 января 2026)

- ✅ Auth-Service: 100% (полностью завершен)
- 🔄 Go-API Refactoring: 15% (создан gRPC клиент)
- 🔲 Document-Service: 0% (в планах)
- 🔲 Company-Service: 0% (в планах)

### Цели на ближайшие 3 дня

- ✅ День 1: AuthProxyService + Controller update → 40%
- ✅ День 2: Middleware + Entity refactoring → 70%
- ✅ День 3: Docker + Integration tests → 100%

---

## 🎯 Финальная цель

**Микросервисная архитектура:**

```
┌─────────────────────────────────────────────────────────────┐
│                      FRONTEND (Next.js)                      │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP REST
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   GO-API (API Gateway)                       │
│  Routing, Auth Proxy, Rate Limiting, Caching                │
└──┬───────────────┬───────────────┬───────────────────────┬──┘
   │ gRPC          │ gRPC          │ gRPC                  │
   ▼               ▼               ▼                       ▼
┌──────────┐  ┌──────────┐  ┌───────────┐  ┌──────────────────┐
│  Auth    │  │ Document │  │  Company  │  │  Notification    │
│ Service  │  │ Service  │  │  Service  │  │   Service        │
└──────────┘  └──────────┘  └───────────┘  └──────────────────┘
```

**Преимущества:**

- ✅ Независимое развертывание сервисов
- ✅ Горизонтальное масштабирование
- ✅ Изоляция отказов
- ✅ Технологическое разнообразие
- ✅ Команды работают независимо

---

## 📚 Полезные команды

```bash
# Запуск auth-service
cd auth-service && go run cmd/auth-service/main.go

# Запуск go-api
cd go-api && go run cmd/api/main.go

# Генерация proto файлов
cd api/proto && make proto

# Запуск всех сервисов через Docker
docker-compose -f docker-compose.microservices.yml up

# Просмотр логов auth-service
docker-compose logs -f auth-service

# Тесты auth-service
cd auth-service && go test ./...

# Тесты go-api
cd go-api && go test ./...

# Проверка gRPC endpoints
grpcurl -plaintext localhost:9001 list
grpcurl -plaintext localhost:9001 api.auth.AuthService/ValidateToken
```

---

**Следующий шаг:** Создать AuthProxyService для проксирования auth запросов на auth-service через gRPC.

Готов продолжить? 🚀
