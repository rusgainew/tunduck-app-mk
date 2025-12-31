#!/bin/bash
set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функция для вывода
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Проверка аргументов
if [ $# -lt 2 ]; then
    print_error "Использование: ./generate-service.sh <service-name> <module-path>"
    echo "Пример: ./generate-service.sh auth-service github.com/rusgainew/tunduck-app"
    exit 1
fi

SERVICE_NAME=$1
MODULE_PATH=$2
SERVICE_DIR="./$SERVICE_NAME"

print_info "Создание микросервиса: $SERVICE_NAME"
print_info "Module path: $MODULE_PATH/$SERVICE_NAME"

# Проверяем, что папка не существует
if [ -d "$SERVICE_DIR" ]; then
    print_error "Папка $SERVICE_DIR уже существует!"
    exit 1
fi

# Создаем структуру папок
print_info "Создание структуры папок..."
mkdir -p "$SERVICE_DIR"/{cmd,internal/{domain,application/{services,dto,commands},infrastructure/{persistence/postgres,config,http/{client,middleware}},interfaces/http/{handlers,routes}},migrations,api,tests,config}

print_success "Структура папок создана"

# Создаем main.go
cat > "$SERVICE_DIR/cmd/main.go" << 'EOF'
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize application
	app, err := NewApp(ctx, ".env")
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	errChan := make(chan error, 1)
	go func() {
		errChan <- app.Run()
	}()

	// Wait for shutdown signal or server error
	select {
	case <-sigChan:
		log.Println("Shutdown signal received, starting graceful shutdown...")
		cancel()
	case err := <-errChan:
		if err != nil {
			log.Printf("Server error: %v", err)
		}
		cancel()
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
EOF

print_success "cmd/main.go создан"

# Создаем app.go
cat > "$SERVICE_DIR/cmd/app.go" << 'EOF'
package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type App struct {
	app    *fiber.App
	log    *logrus.Logger
	config Config
}

type Config struct {
	Port     string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
	RedisURL string
}

func NewApp(ctx context.Context, envFile string) (*App, error) {
	// Load configuration (from env file or environment variables)
	cfg := loadConfig(envFile)

	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// Create Fiber app
	fiberApp := fiber.New()

	// Initialize container (DI)
	// container := container.New(db, cache, log)

	// Register middleware
	// fiberApp.Use(middleware.LoggingMiddleware(log))
	// fiberApp.Use(middleware.ErrorHandler)

	// Register routes
	// handlers.RegisterRoutes(fiberApp, container)

	return &App{
		app:    fiberApp,
		log:    log,
		config: cfg,
	}, nil
}

func (a *App) Run() error {
	return a.app.Listen(":" + a.config.Port)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.app.ShutdownWithContext(ctx)
}

func loadConfig(envFile string) Config {
	return Config{
		Port:     getEnv("PORT", "3001"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "service_db"),
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
EOF

print_success "cmd/app.go создан"

# Создаем Domain layer examples
cat > "$SERVICE_DIR/internal/domain/entity.go" << 'EOF'
package domain

import "github.com/google/uuid"

// TODO: Define your domain entities here
// Example: User, Organization, Document, etc.

// Entity is a base type for domain entities
type Entity struct {
	ID        uuid.UUID
	CreatedAt int64
	UpdatedAt int64
}

// NewEntity creates a new entity with a unique ID
func NewEntity() Entity {
	return Entity{
		ID:        uuid.New(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}
EOF

print_success "internal/domain/entity.go создан"

# Создаем Application layer example
cat > "$SERVICE_DIR/internal/application/services/example_service.go" << 'EOF'
package services

import (
	"context"

	// "github.com/your-module/internal/domain"
)

// ExampleService demonstrates the structure of a use case
type ExampleService struct {
	// Inject dependencies here
	// repo      domain.Repository
	// logger    Logger
	// cache     Cache
}

// NewExampleService creates a new example service
func NewExampleService(
	// dependencies...
) *ExampleService {
	return &ExampleService{
		// Initialize dependencies
	}
}

// Execute demonstrates the basic structure of a use case
func (s *ExampleService) Execute(ctx context.Context) error {
	// TODO: Implement your business logic here
	return nil
}
EOF

print_success "internal/application/services/example_service.go создан"

# Создаем Infrastructure layer example
cat > "$SERVICE_DIR/internal/infrastructure/config/config.go" << 'EOF'
package config

import (
	"os"
)

type Config struct {
	ServiceName string
	Port        string
	Environment string

	Database DatabaseConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	URL string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "service"),
		Port:        getEnv("PORT", "3001"),
		Environment: getEnv("ENVIRONMENT", "development"),

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Database: getEnv("DB_NAME", "service_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},

		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
EOF

print_success "internal/infrastructure/config/config.go создан"

# Создаем Interfaces layer example
cat > "$SERVICE_DIR/internal/interfaces/http/handlers/example_handler.go" << 'EOF'
package handlers

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/your-module/internal/application/services"
)

// ExampleHandler demonstrates HTTP handler structure
type ExampleHandler struct {
	// service *services.ExampleService
}

// NewExampleHandler creates a new example handler
func NewExampleHandler(
	// service *services.ExampleService,
) *ExampleHandler {
	return &ExampleHandler{
		// service: service,
	}
}

// GetExample handles GET /example
// @Summary Get example data
// @Tags example
// @Success 200
// @Router /example [get]
func (h *ExampleHandler) GetExample(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello from handler",
	})
}
EOF

print_success "internal/interfaces/http/handlers/example_handler.go создан"

# Создаем Routes
cat > "$SERVICE_DIR/internal/interfaces/http/routes.go" << 'EOF'
package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/your-module/internal/interfaces/http/handlers"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(app *fiber.App, exampleHandler *handlers.ExampleHandler) {
	api := app.Group("/api")

	// Example routes
	api.Get("/example", exampleHandler.GetExample)

	// TODO: Register your own routes here
}
EOF

print_success "internal/interfaces/http/routes.go создан"

# Создаем go.mod
cat > "$SERVICE_DIR/go.mod" << EOF
module $MODULE_PATH/$SERVICE_NAME

go 1.25

require (
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/google/uuid v1.6.0
	github.com/sirupsen/logrus v1.9.3
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)
EOF

print_success "go.mod создан"

# Создаем Dockerfile
cat > "$SERVICE_DIR/Dockerfile" << 'EOF'
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/service .
COPY .env.example .env

EXPOSE 3001

CMD ["./service"]
EOF

print_success "Dockerfile создан"

# Создаем docker-compose.yml
cat > "$SERVICE_DIR/docker-compose.yml" << 'EOF'
version: '3.8'

services:
  service:
    build: .
    ports:
      - "3001:3001"
    environment:
      - PORT=3001
      - ENVIRONMENT=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=service_db
      - REDIS_URL=redis://redis:6379
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: service_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
EOF

print_success "docker-compose.yml создан"

# Создаем .env.example
cat > "$SERVICE_DIR/.env.example" << 'EOF'
SERVICE_NAME=$SERVICE_NAME
PORT=3001
ENVIRONMENT=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=service_db
DB_SSLMODE=disable

REDIS_URL=redis://localhost:6379

JWT_SECRET=your_jwt_secret_key
EOF

print_success ".env.example создан"

# Создаем Makefile
cat > "$SERVICE_DIR/Makefile" << 'EOF'
.PHONY: help build run test docker-build docker-up docker-down clean

help:
	@echo "Available commands:"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make clean         - Clean build artifacts"

build:
	go build -o service ./cmd

run:
	go run ./cmd/main.go

test:
	go test -v ./...

docker-build:
	docker build -t service .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

clean:
	rm -f service
	go clean
EOF

print_success "Makefile создан"

# Создаем README.md
cat > "$SERVICE_DIR/README.md" << EOF
# $SERVICE_NAME Service

Описание сервиса.

## Quick Start

### Prerequisites
- Go 1.25+
- PostgreSQL 16+
- Redis 7+
- Docker & Docker Compose

### Local Development

1. Clone the repository
2. Copy .env.example to .env
3. Run: \`make docker-up\`
4. Run: \`make run\`

### Docker

\`\`\`bash
make docker-build
make docker-up
\`\`\`

### Testing

\`\`\`bash
make test
\`\`\`

## Project Structure

- \`cmd/\` - Application entry points
- \`internal/\` - Private application code
  - \`domain/\` - Domain layer (business logic)
  - \`application/\` - Application layer (use cases)
  - \`infrastructure/\` - Infrastructure layer (DB, HTTP clients)
  - \`interfaces/\` - Interface layer (HTTP handlers)
- \`migrations/\` - Database migrations
- \`api/\` - OpenAPI/Swagger documentation
- \`tests/\` - Integration tests

## API Documentation

See \`api/openapi.yaml\` or run the server and visit \`http://localhost:3001/swagger\`

## License

MIT
EOF

print_success "README.md создан"

# Создаем OpenAPI template
mkdir -p "$SERVICE_DIR/api"
cat > "$SERVICE_DIR/api/openapi.yaml" << 'EOF'
openapi: 3.0.0
info:
  title: Service API
  version: 1.0.0
  description: Service API documentation

servers:
  - url: http://localhost:3001
    description: Local server

paths:
  /api/example:
    get:
      summary: Get example data
      tags:
        - example
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
EOF

print_success "api/openapi.yaml создан"

# Создаем migration template
cat > "$SERVICE_DIR/migrations/001_init.sql" << 'EOF'
-- TODO: Add your database migrations here
-- Example:
-- CREATE TABLE users (
--     id UUID PRIMARY KEY,
--     email VARCHAR(255) NOT NULL UNIQUE,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
EOF

print_success "migrations/001_init.sql создан"

# Создаем example test
cat > "$SERVICE_DIR/tests/example_test.go" << 'EOF'
package tests

import (
	"testing"
)

func TestExample(t *testing.T) {
	// TODO: Write your tests here
	t.Log("Example test")
}
EOF

print_success "tests/example_test.go создан"

echo ""
print_success "✓ Микросервис '$SERVICE_NAME' успешно создан!"
echo ""
print_info "Следующие шаги:"
echo "  1. cd $SERVICE_DIR"
echo "  2. go mod tidy"
echo "  3. Отредактировать internal/domain/ (ваши entities)"
echo "  4. Отредактировать internal/application/services/"
echo "  5. Отредактировать internal/infrastructure/persistence/"
echo "  6. Отредактировать internal/interfaces/http/handlers/"
echo ""
print_info "Для запуска:"
echo "  make docker-up"
echo "  make run"
echo ""
print_info "Для тестирования:"
echo "  make test"
echo ""
