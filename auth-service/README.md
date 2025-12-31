# Auth Service - Phase 2

Микросервис аутентификации для Tunduck App.

## 📁 Структура

```
auth-service/
├── cmd/
│   └── auth-service/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/                  # Domain Layer (DDD)
│   │   ├── entity/
│   │   │   └── user.go          # User aggregate, Credential, Token value objects
│   │   └── repository/
│   │       └── interfaces.go    # Repository & Event Publisher interfaces
│   │
│   ├── application/             # Application Layer (DDD)
│   │   ├── dto/
│   │   │   └── auth_dto.go      # DTOs: Request, Response mappers
│   │   └── service/
│   │       └── auth_service.go  # RegisterUserService, LoginUserService
│   │
│   ├── infrastructure/          # Infrastructure Layer (DDD)
│   │   ├── config/
│   │   │   └── config.go        # Environment loading
│   │   ├── persistence/
│   │   │   └── postgres/
│   │   │       └── user_repository.go  # PostgreSQL implementation
│   │   ├── event/
│   │   │   └── rabbitmq/
│   │   │       └── event_publisher.go  # RabbitMQ implementation
│   │   └── cache/
│   │       └── redis/
│   │           └── token_blacklist.go  # Redis implementation
│   │
│   └── interfaces/              # Interfaces Layer (DDD)
│       └── http/
│           ├── handler/
│           │   └── auth_handler.go    # HTTP endpoints
│           └── server/
│               └── server.go          # HTTP server
│
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## 🎯 DDD Слои

### Domain Layer (Чистая бизнес-логика)

- **User** - Aggregate root (ID, Email, Name, Password, Status, CreatedAt, UpdatedAt, LastLogin)
- **Credential** - Value object (Email, Password for verification)
- **Token** - Value object (AccessToken, RefreshToken, ExpiresIn, IssuedAt)
- **Role** - Value object (ID, Name)
- **Permission** - Value object (ID, Name, Action)

### Application Layer (Use cases)

- **RegisterUserService** - Регистрация новых пользователей
- **LoginUserService** - Вход и создание токенов
- **TokenService** - Управление JWT токенами
- **DTOs** - Request/Response контракты

### Infrastructure Layer (External dependencies)

- **PostgreSQL Repository** - Персистентность User aggregate
- **RabbitMQ Publisher** - События (user_registered, user_logged_in)
- **Redis Blacklist** - Токены в черном списке (logout)

### Interfaces Layer (API endpoints)

- **HTTP Handlers**:
  - POST /auth/register - Регистрация
  - POST /auth/login - Вход
  - GET /auth/me - Получить профиль
  - POST /auth/logout - Выход
  - POST /auth/refresh - Обновить токен
- **gRPC Handlers** - TODO (использовать api/proto/auth_service.proto)

## 🚀 Быстрый старт

### Требования

- Go 1.25+
- PostgreSQL 16
- Redis 7
- RabbitMQ 3.13

### 1. Запустить инфраструктуру

```bash
cd ..
docker-compose -f docker-compose.microservices.yml up -d
```

### 2. Инициализировать переменные окружения

```bash
cp .env.example .env
```

### 3. Запустить сервис

```bash
cd auth-service
go run cmd/auth-service/main.go
```

Сервис запустится на:

- HTTP: http://localhost:8001
- gRPC: localhost:9001

### 4. Тестировать endpoints

#### Регистрация

```bash
curl -X POST http://localhost:8001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "Test User",
    "password": "SecurePassword123"
  }'
```

#### Вход

```bash
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

#### Получить профиль

```bash
curl -X GET http://localhost:8001/auth/me \
  -H "Authorization: Bearer <access_token>"
```

#### Health Check

```bash
curl http://localhost:8001/health
```

## 🔐 Аутентификация

### JWT токены

- **Access Token**: 1 час (3600 сек) - для каждого запроса
- **Refresh Token**: 7 дней - для получения нового access token

### Password Security

- Хеширование: bcrypt (cost = 12)
- Никогда не возвращаем пароль в API

## 📊 Последовательность событий (RabbitMQ)

Сервис публикует события на exchange `tunduck.auth`:

```
user.registered → user_registered queue
user.logged_in  → user_logged_in queue
user.logged_out → user_logged_out queue
```

Другие сервисы могут подписаться и получать события.

## 🧪 Тестирование

### Unit Tests

```bash
go test ./internal/application/service/...
go test ./internal/domain/entity/...
```

### Integration Tests

```bash
go test ./internal/infrastructure/...
```

### Coverage

```bash
go test -cover ./...
```

## 📝 Гайдлайны разработки

### 1. Ни когда не смешивайте слои

❌ Domain должен быть независим от Infrastructure
❌ Handlers не должны содержать бизнес-логику

### 2. Используйте DTOs

✓ Всегда преобразуйте Entities в DTOs перед возвратом
✓ Принимайте RequestDTOs в handlers

### 3. Обработка ошибок

✓ Заверните external errors: `fmt.Errorf("context: %w", err)`
✓ Используйте domain errors для бизнес-ошибок

### 4. Контекст

✓ Всегда передавайте context.Context первым параметром
✓ Уважайте deadline и cancellation

### 5. Dependency Injection

✓ Используйте factories для создания сервисов
✓ Инжектируйте interfaces, не конкретные типы

## 📚 Дополнительно

- [PHASE2_AUTH_SERVICE_PLAN.md](../../PHASE2_AUTH_SERVICE_PLAN.md) - Детальный план
- [api/proto/README.md](../../api/proto/README.md) - Работа с proto файлами
- [QUICKSTART.md](../../QUICKSTART.md) - Общий quickstart

## 📋 TODO

- [ ] Реализовать JWT (sign, verify)
- [ ] Реализовать gRPC handlers
- [ ] Добавить логирование (zap/logrus)
- [ ] Добавить middleware (CORS, logging, etc)
- [ ] Покрыть тестами (>80% coverage)
- [ ] Добавить Dockerfile
- [ ] Интеграция с API Gateway
- [ ] OpenAPI/Swagger documentation

## 👨‍💼 Статус

- ✅ Структура проекта
- ✅ DDD слои созданы
- ✅ Domain entities
- ✅ Application services (basic)
- ✅ Infrastructure layer
- ✅ HTTP handlers (basic)
- 🔄 JWT implementation
- 🔄 gRPC handlers
- ⏳ Testing & CI/CD
