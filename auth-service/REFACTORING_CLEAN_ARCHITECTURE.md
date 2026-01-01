# Рефакторинг auth-service: Чистая архитектура и SRP

## Проведенные улучшения

### 1. Разделение ответственностей (SRP)

#### Validator (Валидатор)

- **Файл**: `internal/interfaces/grpc/validator/request_validator.go`
- **Ответственность**: Валидация входящих gRPC запросов
- **Методы**:
  - `ValidateRegisterRequest`
  - `ValidateLoginRequest`
  - `ValidateTokenRequest`
  - `ValidateGetUserRequest`
  - `ValidateLogoutRequest`
  - `ValidateRefreshTokenRequest`

#### Adapter (Адаптер)

- **Файл**: `internal/interfaces/grpc/adapter/request_adapter.go`
- **Ответственность**: Преобразование proto-запросов в DTO приложения
- **Методы**:
  - `ToRegisterDTO` - конвертирует RegisterRequest в application DTO
  - `ToLoginDTO` - конвертирует LoginRequest в application DTO

#### Mapper (Маппер)

- **Файлы**:
  - `internal/interfaces/grpc/mapper/user_mapper.go`
  - `internal/interfaces/grpc/mapper/token_mapper.go`
  - `internal/interfaces/grpc/mapper/auth_response_mapper.go`
- **Ответственность**: Преобразование доменных моделей в proto-сообщения
- **Методы**:
  - `ToProtoUser` - entity.User → authpb.User
  - `ToProtoUserWithNames` - entity.User с разделением имен
  - `ToProtoTokenFromLogin` - создание proto Token
  - `ToAuthResponseWithToken` - создание AuthResponse

#### Handler (Обработчик)

- **Файлы**:
  - `internal/interfaces/grpc/handler/register_handler.go`
  - `internal/interfaces/grpc/handler/login_handler.go`
  - `internal/interfaces/grpc/handler/token_handler.go`
  - `internal/interfaces/grpc/handler/user_handler.go`
- **Ответственность**: Координация обработки запросов
- **Каждый handler отвечает за конкретную операцию**:
  - `RegisterHandler` - регистрация
  - `LoginHandler` - вход
  - `TokenHandler` - операции с токенами (validate, refresh)
  - `UserHandler` - операции с пользователями (get, logout)

### 2. Чистая архитектура

#### Слой интерфейсов (interfaces)

```
interfaces/grpc/
├── adapter/          # Адаптация внешних запросов
├── validator/        # Валидация запросов
├── mapper/           # Преобразование моделей
├── handler/          # Бизнес-логика обработки
└── auth_service_server.go  # Точка входа (делегирование)
```

#### AuthServiceServer

Теперь выполняет только **одну ответственность** - делегирование запросов:

```go
func (s *AuthServiceServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
    return s.registerHandler.Handle(ctx, req)
}
```

### 3. Преимущества рефакторинга

#### До рефакторинга

- ❌ AuthServiceServer содержал 200+ строк логики
- ❌ Смешивание валидации, маппинга, бизнес-логики
- ❌ Сложность тестирования
- ❌ Нарушение SRP

#### После рефакторинга

- ✅ Каждый компонент имеет одну ответственность
- ✅ Легко тестировать отдельные части
- ✅ Легко заменить/расширить любой компонент
- ✅ Следование Clean Architecture
- ✅ Улучшенная читаемость кода

### 4. Пример использования

```go
// Создание зависимостей
validator := validator.NewRequestValidator()
adapter := adapter.NewRequestAdapter()
userMapper := mapper.NewUserMapper()
tokenMapper := mapper.NewTokenMapper()
responseMapper := mapper.NewAuthResponseMapper(userMapper, tokenMapper)

// Создание обработчика
registerHandler := handler.NewRegisterHandler(
    validator,
    adapter,
    registerService,
    responseMapper,
    tokenMapper,
)

// Использование
response, err := registerHandler.Handle(ctx, request)
```

### 5. Тестируемость

Каждый компонент можно тестировать изолированно:

```go
// Тест валидатора
func TestValidateRegisterRequest(t *testing.T) {
    v := validator.NewRequestValidator()
    err := v.ValidateRegisterRequest(&authpb.RegisterRequest{
        Email: "test@example.com",
        Password: "password123",
    })
    assert.NoError(t, err)
}

// Тест маппера
func TestToProtoUser(t *testing.T) {
    m := mapper.NewUserMapper()
    user := &entity.User{ID: "1", Email: "test@test.com"}
    proto := m.ToProtoUser(user)
    assert.Equal(t, "1", proto.Id)
}
```

## Итог

Рефакторинг следует принципам:

- ✅ **Single Responsibility Principle** - каждый класс делает одно
- ✅ **Clean Architecture** - разделение по слоям
- ✅ **Dependency Inversion** - зависимости через интерфейсы
- ✅ **Testability** - легко тестируемый код
- ✅ **Maintainability** - легко поддерживаемый код
