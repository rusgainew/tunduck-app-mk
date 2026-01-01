# ✅ Приоритет 3 завершен: RBAC через auth-service

## 🎯 Выполненные задачи

### 1. Обновлен Proto - добавлена роль ✅

**Файл:** `api/proto/auth.proto`

**Изменения:**

```protobuf
message User {
  string id = 1;
  string email = 2;
  string first_name = 3;
  string last_name = 4;
  int64 created_at = 5;
  int64 updated_at = 6;
  string status = 7;
  string role = 8;  // NEW: "admin", "user", "organization_owner"
}
```

### 2. Обновлен AuthProxyService ✅

**Файл:** `go-api/internal/services/service_impl/auth_proxy_service_impl.go`

**Изменения:**

- ✅ Register response включает `Role` из gRPC
- ✅ Login response включает `Role` из gRPC
- ✅ ValidateToken возвращает `Role` из gRPC

```go
// Конвертация gRPC → UserInfo
userInfo := &models.UserInfo{
    ID:       grpcUser.Id,
    Email:    grpcUser.Email,
    FullName: grpcUser.FirstName + " " + grpcUser.LastName,
    Role:     grpcUser.Role, // NEW: RBAC роль
}
```

### 3. Обновлен JWT Middleware ✅

**Файл:** `go-api/pkg/middleware/jwt_auth_middleware.go`

**Изменения:**

- ✅ Сохраняет `role` в контексте после валидации через auth-service
- ✅ Логирует роль при успешной аутентификации
- ✅ Добавлена функция `GetRoleFromContext()`

```go
// В JWTAuthMiddleware
c.Locals("user", user)
c.Locals("user_id", user.ID)
c.Locals("email", user.Email)
c.Locals("role", user.Role) // NEW
```

### 4. Обновлен RBAC Middleware ✅

**Файл:** `go-api/pkg/middleware/rbac_middleware.go`

#### До (старая реализация):

```go
func RBACMiddleware(logger *logrus.Logger, requiredRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Получаем JWT токен из контекста
        user := c.Locals("user")
        claims, ok := user.(*jwt.Token)
        mapClaims := claims.Claims.(jwt.MapClaims)
        role, ok := mapClaims["role"].(string) // Парсинг JWT локально

        // Проверка роли
        if role == requiredRole {
            return c.Next()
        }
    }
}
```

#### После (новая реализация):

```go
func RBACMiddleware(logger *logrus.Logger, requiredRoles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Получаем роль напрямую из контекста (проверена auth-service)
        role := c.Locals("role")
        if role == nil {
            return apperror.New(apperror.ErrUnauthorized, "Authentication required")
        }

        roleStr, ok := role.(string)
        if !ok {
            return apperror.New(apperror.ErrForbidden, "Invalid role format")
        }

        // Проверка роли
        hasRole := false
        for _, requiredRole := range requiredRoles {
            if roleStr == requiredRole {
                hasRole = true
                break
            }
        }

        if !hasRole {
            return apperror.New(apperror.ErrForbidden, "Insufficient permissions")
        }

        return c.Next()
    }
}
```

**Ключевые изменения:**

- ❌ Удалена зависимость от `github.com/golang-jwt/jwt/v5`
- ❌ Не парсит JWT токен локально
- ✅ Использует role из контекста (установлен после проверки auth-service)
- ✅ Единый источник истины для ролей - auth-service

### 5. Обновлен AdminOrSelfMiddleware ✅

#### До:

```go
func AdminOrSelfMiddleware(logger *logrus.Logger, userIDParam string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user")
        claims, ok := user.(*jwt.Token)
        mapClaims := claims.Claims.(jwt.MapClaims)
        role, _ := mapClaims["role"].(string)
        userID, _ := mapClaims["sub"].(string)

        if role == "admin" || userID == targetUserID {
            return c.Next()
        }
    }
}
```

#### После:

```go
func AdminOrSelfMiddleware(logger *logrus.Logger, userIDParam string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Извлекаем данные из контекста (установлены JWTAuthMiddleware)
        userID := c.Locals("user_id")
        role := c.Locals("role")

        if userID == nil || role == nil {
            return apperror.New(apperror.ErrUnauthorized, "Authentication required")
        }

        userIDStr, ok := userID.(string)
        roleStr, ok := role.(string)

        // Админ имеет доступ ко всем профилям
        if roleStr == "admin" {
            return c.Next()
        }

        // Пользователь может редактировать только свой профиль
        targetUserID := c.Params(userIDParam)
        if userIDStr == targetUserID {
            return c.Next()
        }

        return apperror.New(apperror.ErrForbidden, "You can only edit your own profile")
    }
}
```

**Ключевые изменения:**

- ❌ Не парсит JWT claims
- ✅ Использует `user_id` и `role` из контекста
- ✅ Более простая логика
- ✅ Меньше точек отказа

---

## 📊 Архитектура RBAC

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT REQUEST                           │
│                                                               │
│  PUT /api/users/123                                          │
│  Authorization: Bearer eyJhbGc...                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                GO-API Middleware Chain                       │
│                                                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  1. JWTAuthMiddleware                                   │ │
│  │     - Извлекает token                                   │ │
│  │     - Валидирует через authProxyService.ValidateToken() │ │
│  │     - Получает UserInfo с role                          │ │
│  │     - Сохраняет в контекст:                             │ │
│  │       * c.Locals("user", UserInfo)                      │ │
│  │       * c.Locals("user_id", "user_123")                 │ │
│  │       * c.Locals("email", "user@example.com")           │ │
│  │       * c.Locals("role", "admin")                       │ │
│  └────────────────────────────────────────────────────────┘ │
│                          │                                    │
│                          ▼                                    │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  2. AdminOrSelfMiddleware                               │ │
│  │     - Читает role из c.Locals("role")                   │ │
│  │     - Читает user_id из c.Locals("user_id")             │ │
│  │     - Если role == "admin" → Allow                      │ │
│  │     - Если user_id == :id param → Allow                 │ │
│  │     - Иначе → 403 Forbidden                             │ │
│  └────────────────────────────────────────────────────────┘ │
│                          │                                    │
│                          ▼                                    │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  3. UserController.updateUser                           │ │
│  │     - Обновляет профиль пользователя                    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                         │
                         │ gRPC (при JWTAuthMiddleware)
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              AUTH-SERVICE (Port 9001 gRPC)                   │
│                                                               │
│  ValidateToken RPC                                           │
│  - Проверяет JWT signature                                   │
│  - Проверяет expiration                                      │
│  - Загружает User из PostgreSQL                              │
│  - Возвращает User с role                                    │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔄 Поток RBAC проверки

### Сценарий 1: Admin редактирует любой профиль

```http
PUT /api/users/456 HTTP/1.1
Authorization: Bearer <admin_token>
```

**Шаг 1: JWTAuthMiddleware**

```go
token → auth-service.ValidateToken(token)
      ← User{id: "123", email: "admin@example.com", role: "admin"}

c.Locals("user_id", "123")
c.Locals("role", "admin")
c.Next() // Переход к следующему middleware
```

**Шаг 2: AdminOrSelfMiddleware**

```go
role := c.Locals("role") // "admin"
userID := c.Locals("user_id") // "123"
targetID := c.Params("id") // "456"

if role == "admin" {
    return c.Next() // ✅ Admin может редактировать любой профиль
}
```

**Результат:** ✅ 200 OK - профиль обновлен

---

### Сценарий 2: User редактирует свой профиль

```http
PUT /api/users/123 HTTP/1.1
Authorization: Bearer <user_token>
```

**Шаг 1: JWTAuthMiddleware**

```go
token → auth-service.ValidateToken(token)
      ← User{id: "123", email: "user@example.com", role: "user"}

c.Locals("user_id", "123")
c.Locals("role", "user")
c.Next()
```

**Шаг 2: AdminOrSelfMiddleware**

```go
role := c.Locals("role") // "user"
userID := c.Locals("user_id") // "123"
targetID := c.Params("id") // "123"

if role == "admin" {
    // false - пропускаем
}

if userID == targetID { // "123" == "123"
    return c.Next() // ✅ User может редактировать свой профиль
}
```

**Результат:** ✅ 200 OK - профиль обновлен

---

### Сценарий 3: User пытается редактировать чужой профиль

```http
PUT /api/users/456 HTTP/1.1
Authorization: Bearer <user_token>
```

**Шаг 1: JWTAuthMiddleware**

```go
token → auth-service.ValidateToken(token)
      ← User{id: "123", email: "user@example.com", role: "user"}

c.Locals("user_id", "123")
c.Locals("role", "user")
c.Next()
```

**Шаг 2: AdminOrSelfMiddleware**

```go
role := c.Locals("role") // "user"
userID := c.Locals("user_id") // "123"
targetID := c.Params("id") // "456"

if role == "admin" {
    // false - пропускаем
}

if userID == targetID { // "123" != "456"
    // false - пропускаем
}

// Доступ запрещен
return apperror.New(apperror.ErrForbidden, "You can only edit your own profile")
```

**Результат:** ❌ 403 Forbidden

---

## 📈 Преимущества новой RBAC архитектуры

### 1. Единый источник истины ✅

- Роли управляются только в auth-service
- go-api не парсит JWT локально
- Нет дублирования логики RBAC

### 2. Безопасность ✅

- Роль проверяется auth-service
- Невозможно подделать роль в go-api
- JWT_SECRET известен только auth-service

### 3. Производительность ✅

- Роль извлекается один раз (при ValidateToken)
- Не нужно парсить JWT в каждом middleware
- Меньше нагрузки на CPU

### 4. Простота кода ✅

- Удалена зависимость от `jwt/v5` в RBAC middleware
- Меньше кода
- Меньше точек отказа

### 5. Масштабируемость ✅

- Можно добавить новые роли в auth-service
- go-api автоматически получит новые роли
- Единая точка управления RBAC

---

## 🧪 Тестирование

### 1. Проверка роли Admin

```bash
# Регистрация админа (через auth-service напрямую)
curl -X POST http://localhost:8001/api/auth/register-admin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!",
    "firstName": "Admin",
    "lastName": "User",
    "adminSecret": "secret_key"
  }'

# Login через go-api
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!"
  }'

# Сохранить токен
ADMIN_TOKEN="<token_from_response>"

# Проверка доступа к чужому профилю (должно работать)
curl -X PUT http://localhost:8080/api/users/any_user_id \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Updated Name"
  }'
```

**Ожидается:** ✅ 200 OK - админ может редактировать любой профиль

### 2. Проверка роли User

```bash
# Регистрация обычного пользователя
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "User123!",
    "fullName": "Regular User",
    "confirmPassword": "User123!"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "User123!"
  }'

# Извлечь user_id и token из ответа
USER_ID="<user_id_from_response>"
USER_TOKEN="<token_from_response>"

# Попытка редактировать свой профиль (должно работать)
curl -X PUT http://localhost:8080/api/users/$USER_ID \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Updated Regular User"
  }'
```

**Ожидается:** ✅ 200 OK - пользователь может редактировать свой профиль

```bash
# Попытка редактировать чужой профиль (должно быть запрещено)
curl -X PUT http://localhost:8080/api/users/another_user_id \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "fullName": "Hacked Name"
  }'
```

**Ожидается:** ❌ 403 Forbidden - пользователь не может редактировать чужой профиль

### 3. Проверка AdminOnlyMiddleware

```bash
# Попытка удалить пользователя обычным user (должно быть запрещено)
curl -X DELETE http://localhost:8080/api/users/$USER_ID \
  -H "Authorization: Bearer $USER_TOKEN"
```

**Ожидается:** ❌ 403 Forbidden - требуется роль admin

```bash
# Удаление пользователя админом (должно работать)
curl -X DELETE http://localhost:8080/api/users/$USER_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Ожидается:** ✅ 200 OK - админ может удалять пользователей

---

## 📝 Изменённые файлы

```
tunduck-app-mk/
├── api/
│   └── proto/
│       └── auth.proto                                # UPDATED: добавлено поле role
├── go-api/
│   ├── internal/
│   │   ├── models/
│   │   │   └── user_model.go                        # EXISTS: Role уже был
│   │   └── services/
│   │       └── service_impl/
│   │           └── auth_proxy_service_impl.go       # UPDATED: добавлен Role в конвертацию
│   └── pkg/
│       └── middleware/
│           ├── jwt_auth_middleware.go               # UPDATED: добавлены role в контекст и GetRoleFromContext()
│           └── rbac_middleware.go                   # REFACTORED: убрана зависимость от JWT, использует контекст
```

**Всего:** 1 proto обновление, 3 файла обновлены

---

## ✅ Готово к использованию!

RBAC теперь полностью интегрирован с auth-service. Все проверки ролей используют данные из auth-service через gRPC.

**Следующий шаг:** Создать Docker Compose для оркестрации микросервисов.

---

_Дата завершения: 1 января 2026_
