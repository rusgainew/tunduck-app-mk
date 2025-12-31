# Microservice Scaffold Generator

Скрипт для автоматического создания новых микросервисов с DDD структурой.

## 🚀 Использование

```bash
./generate-service.sh auth-service github.com/rusgainew/tunduck-app
./generate-service.sh company-service github.com/rusgainew/tunduck-app
./generate-service.sh document-service github.com/rusgainew/tunduck-app
```

## 📁 Что создается

Каждый микросервис имеет структуру:

```
{service-name}/
├── cmd/
│   ├── main.go
│   └── app.go
├── internal/
│   ├── domain/          # Domain Layer
│   ├── application/     # Application Layer
│   ├── infrastructure/  # Infrastructure Layer
│   └── interfaces/      # Interface Layer
├── migrations/
├── api/
├── tests/
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── Makefile
└── README.md
```

## 🔄 Следующие шаги

1. Отредактировать `go.mod` (module path)
2. Реализовать Domain entities в `internal/domain/`
3. Реализовать Use Cases в `internal/application/services/`
4. Реализовать Repository в `internal/infrastructure/persistence/`
5. Реализовать HTTP handlers в `internal/interfaces/http/handlers/`
