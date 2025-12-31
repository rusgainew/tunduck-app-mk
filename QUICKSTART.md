# 🚀 Quick Start Guide для Microservices Architecture

## 📋 Предварительные требования

- Docker & Docker Compose (v20.10+)
- Go 1.25+ (для локальной разработки)
- Make (для выполнения команд)
- protoc compiler (для work с proto файлами)

---

## 🔧 Setup & Installation

### 1️⃣ Установка инструментов

```bash
# Установить protobuf compiler
# macOS
brew install protobuf

# Linux (Ubuntu/Debian)
sudo apt-get install protobuf-compiler

# Установить Go plugins для protoc
cd api/proto
make proto-install-tools
```

### 2️⃣ Скачать проект

```bash
git clone https://github.com/rusgainew/tunduck-app-mk.git
cd tunduck-app-mk
```

### 3️⃣ Компилировать proto файлы

```bash
cd api/proto
make proto
# Сгенерирует *.pb.go и *_grpc.pb.go файлы
```

---

## 🐳 Запуск Docker Compose

### Phase 1: Базовая инфраструктура (Database, Cache, Message Queue)

```bash
# Запустить PostgreSQL, Redis, RabbitMQ
docker-compose -f docker-compose.microservices.yml up -d postgres redis rabbitmq

# Проверить статус
docker-compose -f docker-compose.microservices.yml ps

# Просмотреть логи
docker-compose -f docker-compose.microservices.yml logs -f postgres
docker-compose -f docker-compose.microservices.yml logs -f redis
docker-compose -f docker-compose.microservices.yml logs -f rabbitmq
```

### Phase 2: Запуск Auth-Service (когда будет готов)

```bash
# Раскомментировать в docker-compose.microservices.yml секцию auth-service
# Затем запустить:

docker-compose -f docker-compose.microservices.yml up -d auth-service

# Проверить:
curl http://localhost:8001/health
```

### Запустить все сервисы сразу

```bash
docker-compose -f docker-compose.microservices.yml up -d
```

---

## 🧪 Тестирование

### PostgreSQL

```bash
# Подключиться к БД
docker exec -it tunduck-postgres psql -U tunduck_user -d tunduck_db

# Проверить таблицы
\dt

# Выйти
\q
```

### Redis

```bash
# Проверить Redis
docker exec -it tunduck-redis redis-cli ping
# Output: PONG

# Просмотреть ключи
docker exec -it tunduck-redis redis-cli KEYS "*"
```

### RabbitMQ

```bash
# Веб-интерфейс Management
# http://localhost:15672
# Логин: tunduck_user
# Пароль: tunduck_password_dev

# Проверить статус через CLI
docker exec -it tunduck-rabbitmq rabbitmqctl status

# Просмотреть exchanges
docker exec -it tunduck-rabbitmq rabbitmqctl list_exchanges

# Просмотреть queues
docker exec -it tunduck-rabbitmq rabbitmqctl list_queues

# Просмотреть bindings
docker exec -it tunduck-rabbitmq rabbitmqctl list_bindings
```

---

## 📁 Структура проекта после setup

```
tunduck-app-mk/
├── api/
│   └── proto/                      # 📌 Proto files (скомпилированы в *.pb.go)
│       ├── *.proto
│       ├── *.pb.go               # Generated
│       └── *_grpc.pb.go          # Generated
├── auth-service/                   # Phase 2 (under development)
├── company-service/                # Phase 3 (planned)
├── document-service/               # Phase 4 (planned)
├── api-gateway/                    # Phase 5 (planned)
├── docker-compose.microservices.yml
├── config/
│   └── rabbitmq.conf
├── scripts/
│   └── init-rabbitmq.sh
└── ...
```

---

## 🔄 Development Workflow

### Локальная разработка Auth-Service

```bash
# 1. Перейти в директорию сервиса
cd auth-service

# 2. Скопировать и компилировать proto файлы
cp -r ../api/proto .
cd api/proto && make proto && cd ../..

# 3. Установить зависимости
go mod download

# 4. Запустить миграции БД
go run cmd/main.go migrate

# 5. Запустить сервис в режиме разработки
go run cmd/main.go serve

# 6. В другом терминале - запустить тесты
go test ./...

# 7. С покрытием
go test -cover ./...
```

### Быстрая перекомпиляция proto

```bash
cd api/proto
make proto-clean
make proto
```

---

## 🐛 Debugging & Logs

### Просмотреть логи конкретного сервиса

```bash
docker-compose -f docker-compose.microservices.yml logs -f auth-service
docker-compose -f docker-compose.microservices.yml logs -f company-service
docker-compose -f docker-compose.microservices.yml logs -f rabbitmq --tail=100
```

### Подключиться в контейнер

```bash
# Интерактивный bash в контейнере
docker-compose -f docker-compose.microservices.yml exec auth-service bash

# Выполнить команду
docker-compose -f docker-compose.microservices.yml exec postgres psql -U tunduck_user -d tunduck_db
```

### Очистить все данные (Warning! Удаляет volumes)

```bash
# Остановить контейнеры
docker-compose -f docker-compose.microservices.yml down

# Удалить volumes (данные БД)
docker-compose -f docker-compose.microservices.yml down -v

# Запустить заново
docker-compose -f docker-compose.microservices.yml up -d
```

---

## 📊 Health Checks

### HTTP Endpoints

```bash
# Auth Service health
curl http://localhost:8001/health

# Company Service health
curl http://localhost:8002/health

# Document Service health
curl http://localhost:8003/health

# API Gateway health
curl http://localhost:8000/health
```

### gRPC Health Check

```bash
# Используем grpcurl (install: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest)

# Auth Service gRPC
grpcurl -plaintext localhost:9001 list
grpcurl -plaintext localhost:9001 api.auth.AuthService/ValidateToken

# Company Service gRPC
grpcurl -plaintext localhost:9002 api.company.CompanyService/GetOrganization
```

---

## 🔐 Environment Variables

### Для локальной разработки (используется в docker-compose)

```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=tunduck_user
DB_PASSWORD=tunduck_password_dev
DB_NAME=tunduck_db

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_URL=amqp://tunduck_user:tunduck_password_dev@rabbitmq:5672/tunduck

# Service Configuration
SERVICE_PORT=8001          # HTTP port
GRPC_PORT=9001            # gRPC port
ENV=development           # Environment

# gRPC Service Addresses
AUTH_SERVICE_GRPC=auth-service:9001
COMPANY_SERVICE_GRPC=company-service:9002
DOCUMENT_SERVICE_GRPC=document-service:9003
```

### Для production (переопределить в deployment)

```bash
# Same as above, но с production values
DB_PASSWORD=<strong-password>
REDIS_PASSWORD=<strong-password>
RABBITMQ_URL=amqp://<user>:<password>@rabbitmq.prod:5672/tunduck
ENV=production
```

---

## 📚 Документация

| Документ                                                   | Описание                               |
| ---------------------------------------------------------- | -------------------------------------- |
| [ARCHITECTURE.md](ARCHITECTURE.md)                         | Общая архитектура проекта              |
| [REFACTORING_PLAN.md](REFACTORING_PLAN.md)                 | План рефакторинга на микросервисы      |
| [PHASE2_AUTH_SERVICE_PLAN.md](PHASE2_AUTH_SERVICE_PLAN.md) | Детальный план разработки Auth-Service |
| [api/proto/README.md](api/proto/README.md)                 | Руководство по proto файлам            |
| [RECOMMENDATIONS.md](RECOMMENDATIONS.md)                   | Рекомендации и best practices          |

---

## 🆘 Troubleshooting

### RabbitMQ Management UI не доступен

```bash
# Проверить, запущен ли rabbitmq
docker-compose -f docker-compose.microservices.yml ps rabbitmq

# Проверить логи
docker-compose -f docker-compose.microservices.yml logs rabbitmq

# Перезапустить
docker-compose -f docker-compose.microservices.yml restart rabbitmq
```

### PostgreSQL connection refused

```bash
# Убедиться, что postgres запущен
docker-compose -f docker-compose.microservices.yml up -d postgres

# Дождаться healthcheck-а (может занять 30 сек)
docker-compose -f docker-compose.microservices.yml logs postgres
```

### Proto files не компилируются

```bash
# Проверить версию protoc
protoc --version

# Проверить установку plugins
which protoc-gen-go
which protoc-gen-go-grpc

# Переустановить plugins
cd api/proto && make proto-install-tools

# Очистить и перекомпилировать
make proto-clean && make proto
```

---

## 🚀 Next Steps

1. ✅ Phase 1: Proto files setup (DONE)
2. 📅 Phase 2: Разработка Auth-Service (следующий этап)
3. 📅 Phase 3: Разработка Company-Service
4. 📅 Phase 4: Разработка Document-Service
5. 📅 Phase 5: Разработка API-Gateway
6. 📅 Phase 6: Database per Service migration
7. 📅 Phase 7: DevOps & Production deployment

---

## 📞 Support

Для вопросов и проблем:

1. Проверить логи: `docker-compose logs -f <service>`
2. Проверить документацию: [RECOMMENDATIONS.md](RECOMMENDATIONS.md)
3. Проверить health endpoints
4. Проверить environment variables

---

## 📜 License

Этот проект является частью tunduck-app-mk.
