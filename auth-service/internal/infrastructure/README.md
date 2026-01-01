# Infrastructure Layer

## Компоненты

### 1. PostgreSQL Repository

**Файл:** `persistence/postgres/user_repository.go`

**Методы:**

- `CreateUser` - создание пользователя
- `GetUserByID` - получение по ID
- `GetUserByEmail` - получение по email
- `UpdateUser` - обновление данных
- `DeleteUser` - удаление
- `UserExists` - проверка существования

**Connection:** `connection.go`

- Connection pool: 25 max open, 5 idle
- Max lifetime: 5 minutes

### 2. RabbitMQ Event Publisher

**Файл:** `event/rabbitmq/event_publisher.go`

**Метод:**

- `Publish(ctx, event DomainEvent)` - универсальная публикация

**Exchange:** `tunduck.auth` (topic)
**Routing Keys:**

- `user.registered`
- `user.logged_in`
- `user.logged_out`
- `user.blocked`
- `user.password_changed`

### 3. Redis Token Blacklist

**Файл:** `cache/redis/token_blacklist.go`

**Методы:**

- `AddToBlacklist(ctx, token)` - добавить в blacklist (TTL 24h)
- `IsBlacklisted(ctx, token)` - проверить

**Key Format:** `blacklist:{token}`

### 4. JWT Manager

**Файл:** `jwt/jwt_manager.go`

**Методы:**

- `GenerateAccessToken(userID, email)`
- `GenerateRefreshToken(userID)`
- `ValidateAccessToken(token)`
- `ValidateRefreshToken(token)`

## Миграции

**Файл:** `migrations/001_create_users_table.sql`

```sql
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NULL
);
```

**Индексы:**

- `idx_users_email` - на email
- `idx_users_status` - на status
- `idx_users_created_at` - на created_at

## Конфигурация

### PostgreSQL

```go
Config{
    Host: "localhost",
    Port: 5432,
    User: "postgres",
    Password: "postgres",
    DBName: "auth_service",
    SSLMode: "disable",
    MaxOpenConns: 25,
    MaxIdleConns: 5,
    ConnMaxLifetime: 5 * time.Minute,
}
```

### RabbitMQ

```go
Config{
    Host: "localhost",
    Port: 5672,
    User: "guest",
    Password: "guest",
    VHost: "/",
}
```

### Redis

```go
Config{
    Host: "localhost",
    Port: 6379,
    Password: "",
    DB: 0,
}
```

## Использование

```go
// PostgreSQL
db, _ := postgres.NewConnection(postgres.DefaultConfig())
userRepo := postgres.NewUserRepositoryPostgres(db)

// RabbitMQ
conn, _ := rabbitmq.NewConnection(rabbitmq.DefaultConfig())
ch, _ := rabbitmq.NewChannel(conn)
rabbitmq.DeclareExchange(ch, "tunduck.auth")
eventPublisher := rabbitmq.NewEventPublisherRabbitMQ(ch)

// Redis
redisClient, _ := redis.NewConnection(redis.DefaultConfig())
tokenBlacklist := redis.NewTokenBlacklistRedis(redisClient)
```
