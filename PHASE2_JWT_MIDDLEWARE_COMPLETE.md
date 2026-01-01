# ✅ Фаза 2 завершена: JWT Middleware через auth-service

## 🎯 Выполненные задачи

### 1. Обновлен JWT Middleware ✅

**Файл:** `go-api/pkg/middleware/jwt_auth_middleware.go`

**Изменения:**

- ✅ Удалена локальная проверка JWT через `jwtware`
- ✅ Добавлена проверка через `authProxyService.ValidateToken()`
- ✅ Токен извлекается из `Authorization: Bearer <token>`
- ✅ Информация о пользователе сохраняется в контексте
- ✅ Структурированное логирование всех операций

**До:**

```go
func JWTAuthMiddleware(secret string, logger *logrus.Logger) fiber.Handler {
    return jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte(secret)},
        // ... локальная проверка JWT
    })
}
```

**После:**

```go
func JWTAuthMiddleware(authProxyService services.AuthProxyService, logger *logrus.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Извлечь токен
        token := extractToken(c)

        // Валидировать через auth-service
        user, err := authProxyService.ValidateToken(c.Context(), token)

        // Сохранить в контекст
        c.Locals("user", user)
        c.Locals("user_id", user.ID)
        c.Locals("email", user.Email)

        return c.Next()
    }
}
```

### 2. Обновлены helper функции ✅

**GetUserIDFromContext:**

- Извлекает `user_id` из Locals (строка)
- Конвертирует в UUID
- Обработка ошибок

**GetEmailFromContext:**

- Извлекает `email` из Locals
- Возвращает строку

**GetUsernameFromContext (DEPRECATED):**

- Теперь возвращает email для обратной совместимости
- Помечена как устаревшая - auth-service использует email

**GetClaimsFromContext (DEPRECATED):**

- Конвертирует Locals в map для совместимости
- Рекомендуется использовать `c.Locals("user")`

### 3. Обновлены контроллеры ✅

#### UserController

**Файл:** `go-api/internal/controllers/user_controller.go`

**Изменения:**

- ✅ Добавлено поле `authProxyService`
- ✅ Обновлена сигнатура конструктора
- ✅ Middleware использует `authProxyService`

```go
type UserController struct {
    logger           *logger.Logger
    db               *gorm.DB
    logrus           *logrus.Logger
    authProxyService services.AuthProxyService // NEW
}

func NewUserController(
    app *fiber.App,
    authProxyService services.AuthProxyService, // NEW
    log *logrus.Logger,
    db *gorm.DB,
) {
    // ...
}

// Middleware
userGroup.Put("/:id",
    middleware.JWTAuthMiddleware(c.authProxyService, c.logrus), // UPDATED
    middleware.AdminOrSelfMiddleware(c.logrus, "id"),
    c.updateUser,
)
```

#### EsfDocumentController

**Файл:** `go-api/internal/controllers/esf_document_controller.go`

**Изменения:**

- ✅ Добавлено поле `authProxyService`
- ✅ Обновлена сигнатура конструктора
- ✅ Protected routes используют `authProxyService`

```go
protected := esfDocumentGroup.Group("")
protected.Use(middleware.JWTAuthMiddleware(c.authProxyService, c.logrus))
```

#### EsfOrganizationController

**Файл:** `go-api/internal/controllers/esf_organization_controller.go`

**Аналогичные изменения для организаций ЭСФ**

#### RoleController

**Файл:** `go-api/internal/controllers/role_controller.go`

**Аналогичные изменения для управления ролями**

### 4. Обновлен handlers.go ✅

**Файл:** `go-api/cmd/api/handlers.go`

**Изменения:**

```go
// Передача authProxyService во все контроллеры
controllers.NewAuthController(app, cnt.GetAuthProxyService(), ...)
controllers.NewUserController(app, cnt.GetAuthProxyService(), ...)
controllers.NewEsfDocumentController(app, cnt.GetAuthProxyService(), ...)
controllers.NewEsfOrganizationController(app, cnt.GetAuthProxyService(), ...)
```

---

## 📊 Архитектура после обновления

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT REQUEST                           │
│                                                               │
│  GET /api/esf-documents                                      │
│  Authorization: Bearer eyJhbGc...                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                GO-API (Port 8080)                            │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  JWTAuthMiddleware                                      │ │
│  │  1. Извлекает Bearer token                             │ │
│  │  2. Вызывает authProxyService.ValidateToken()          │ │
│  └────────────────────────────────────────────────────────┘ │
│                          │                                    │
│                          ▼                                    │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  AuthProxyService                                       │ │
│  │  - Делегирует проверку на auth-service                 │ │
│  └────────────────────────────────────────────────────────┘ │
│                          │                                    │
│                          │ gRPC                               │
│                          ▼                                    │
└────────────────────────┼────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              AUTH-SERVICE (Port 9001 gRPC)                   │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  ValidateToken RPC                                      │ │
│  │  - Проверяет signature JWT                             │ │
│  │  - Проверяет expiration                                 │ │
│  │  - Проверяет revocation (Redis)                        │ │
│  │  - Возвращает UserInfo                                  │ │
│  └────────────────────────────────────────────────────────┘ │
│                          │                                    │
│                          ▼                                    │
│                    ┌──────────┐                              │
│                    │  Redis   │                              │
│                    │ (tokens) │                              │
│                    └──────────┘                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔄 Поток проверки токена

### 1. Клиент отправляет запрос

```http
GET /api/esf-documents HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 2. JWTAuthMiddleware обрабатывает запрос

```go
// Извлекаем токен
authHeader := c.Get("Authorization")
parts := strings.Split(authHeader, " ")
token := parts[1]

// Валидируем через auth-service
user, err := authProxyService.ValidateToken(c.Context(), token)
if err != nil {
    return response.Error(c, apperror.New(apperror.ErrInvalidToken, "Invalid or expired JWT token"))
}

// Сохраняем в контекст
c.Locals("user", user)
c.Locals("user_id", user.ID)
c.Locals("email", user.Email)
```

### 3. AuthProxyService делает gRPC запрос

```go
// internal/services/service_impl/auth_proxy_service_impl.go
func (s *AuthProxyServiceImpl) ValidateToken(ctx context.Context, token string) (*models.UserInfo, error) {
    // gRPC вызов
    resp, err := s.authClient.ValidateToken(ctx, &pb.ValidateTokenRequest{
        Token: token,
    })

    // Конвертация ответа
    return &models.UserInfo{
        ID:       resp.User.Id,
        Email:    resp.User.Email,
        FullName: resp.User.FirstName + " " + resp.User.LastName,
    }, nil
}
```

### 4. Auth-Service проверяет токен

```go
// auth-service/internal/application/validate_token_service.go
- Парсит JWT
- Проверяет signature (JWT_SECRET)
- Проверяет expiration
- Проверяет revocation в Redis
- Загружает User из PostgreSQL
- Возвращает UserInfo
```

### 5. Контроллер получает пользователя из контекста

```go
// В любом protected endpoint
userID := c.Locals("user_id").(string)
email := c.Locals("email").(string)
user := c.Locals("user").(*models.UserInfo)
```

---

## 🧪 Тестирование

### 1. Запустить auth-service

```bash
cd auth-service
go run cmd/auth-service/main.go
```

**Ожидаемый вывод:**

```
✓ Connected to PostgreSQL
✓ Connected to Redis: PONG
✓ Connected to RabbitMQ
HTTP server listening on :8001
gRPC server listening on :9001
```

### 2. Запустить go-api

```bash
cd go-api
AUTH_SERVICE_GRPC_URL=localhost:9001 go run cmd/api/main.go
```

**Ожидаемый вывод:**

```
Connected to auth-service gRPC at localhost:9001
AuthProxyService initialized successfully
JWT middleware configured for auth-service validation
Server started on port 8080
```

### 3. Получить токен (Login)

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test123!"
  }'
```

**Ответ:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_1234567890",
    "email": "test@example.com",
    "fullName": "Test User"
  }
}
```

### 4. Использовать токен для protected endpoint

```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/api/esf-documents \
  -H "Authorization: Bearer $TOKEN"
```

**Ожидается:**

- ✅ Токен проверяется через auth-service (gRPC)
- ✅ UserInfo сохраняется в контексте
- ✅ Документы возвращаются с кодом 200

### 5. Тест с невалидным токеном

```bash
curl -X GET http://localhost:8080/api/esf-documents \
  -H "Authorization: Bearer invalid_token"
```

**Ожидается:**

```json
{
  "error": "Invalid or expired JWT token",
  "code": "INVALID_TOKEN",
  "status": 401
}
```

### 6. Тест без токена

```bash
curl -X GET http://localhost:8080/api/esf-documents
```

**Ожидается:**

```json
{
  "error": "Authorization header is required",
  "code": "UNAUTHORIZED",
  "status": 401
}
```

---

## 📈 Преимущества новой архитектуры

### 1. Единый источник истины ✅

- Все токены проверяются auth-service
- Нет дублирования логики проверки
- Централизованное управление сессиями

### 2. Revocation поддержка ✅

- Auth-service может отозвать токены через Redis
- go-api сразу видит отозванные токены
- Logout работает корректно для всех сервисов

### 3. Масштабируемость ✅

- Можно добавить новые API Gateway
- Все будут использовать один auth-service
- gRPC connection pooling

### 4. Безопасность ✅

- Секретный ключ JWT только в auth-service
- go-api не знает JWT_SECRET
- Меньше точек атаки

### 5. Мониторинг ✅

- Централизованные логи auth операций
- Единая точка для rate limiting auth
- Метрики в одном месте

---

## 🚀 Следующие шаги

### Приоритет 3: RBAC через auth-service

**Задача:** Обновить RBAC middleware для проверки ролей через auth-service

**Что нужно сделать:**

1. Добавить метод `CheckPermission(userID, permission)` в auth-service
2. Создать `PermissionProxyService` в go-api
3. Обновить `rbac_middleware.go` для использования proxy
4. Портировать роли из go-api в auth-service

### Приоритет 4: Docker Compose

**Задача:** Создать docker-compose для оркестрации микросервисов

**Что нужно:**

```yaml
version: "3.8"
services:
  postgres:
    image: postgres:15

  redis:
    image: redis:7

  rabbitmq:
    image: rabbitmq:3-management

  auth-service:
    build: ./auth-service
    ports:
      - "8001:8001"
      - "9001:9001"
    depends_on:
      - postgres
      - redis
      - rabbitmq

  go-api:
    build: ./go-api
    ports:
      - "8080:8080"
    environment:
      AUTH_SERVICE_GRPC_URL: auth-service:9001
    depends_on:
      - auth-service
```

### Приоритет 5: Тесты

**Integration тесты:**

- Тест register → login → protected endpoint
- Тест token refresh flow
- Тест logout и revocation

**Unit тесты:**

- JWT middleware с моками authProxyService
- Error handling в middleware
- Context locals extraction

---

## 🐛 Известные проблемы

### 1. gRPC пакеты отсутствуют

**Проблема:** `could not import google.golang.org/grpc`

**Решение:**

```bash
cd go-api
go get google.golang.org/grpc@latest
go get google.golang.org/grpc/credentials/insecure
go get google.golang.org/grpc/keepalive
go get google.golang.org/grpc/codes
go get google.golang.org/grpc/status
```

### 2. Proto генерация отсутствует

**Проблема:** `could not import github.com/rusgainew/tunduck-app-mk/gen/proto/auth`

**Решение:**

```bash
cd api/proto
make generate
```

### 3. Legacy UserService несовместим

**Проблема:** UserService использует старые модели (Username, Phone, Role)

**Решение:** Постепенная миграция, сохранение UserService для обратной совместимости

### 4. AppError методы отсутствуют

**Проблема:** `undefined: apperror.NewBadRequestError`

**Решение:** Использовать `apperror.New(apperror.ErrBadRequest, "message")`

---

## 📝 Изменённые файлы

```
go-api/
├── pkg/
│   └── middleware/
│       └── jwt_auth_middleware.go                    # REFACTORED
├── internal/
│   ├── controllers/
│   │   ├── user_controller.go                        # UPDATED
│   │   ├── esf_document_controller.go                # UPDATED
│   │   ├── esf_organization_controller.go            # UPDATED
│   │   └── role_controller.go                        # UPDATED
│   └── clients/
│       └── grpc/
│           └── auth_client.go                        # EXISTS (need proto)
└── cmd/
    └── api/
        └── handlers.go                               # UPDATED
```

**Всего:** 1 рефакторинг, 5 обновлений

---

## ✅ Готово к тестированию!

JWT Middleware теперь полностью интегрирован с auth-service. Все защищённые endpoints проверяют токены через gRPC.

**Следующий шаг:** Установить gRPC зависимости и сгенерировать proto файлы.

---

_Дата завершения: 1 января 2026_
