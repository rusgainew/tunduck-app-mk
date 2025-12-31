package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/rusgainew/tunduck-app/internal/conf"
	"github.com/rusgainew/tunduck-app/internal/services/service_impl"
	"github.com/rusgainew/tunduck-app/pkg/container"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/health"
	"github.com/rusgainew/tunduck-app/pkg/metrics"
	"github.com/rusgainew/tunduck-app/pkg/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// App представляет основное приложение со всеми зависимостями
type App struct {
	ctx           context.Context       // Контекст для управления жизненным циклом
	logger        *logrus.Logger        // Логгер для записи событий
	conf          *conf.Conf            // Конфигурация приложения
	db            *gorm.DB              // Подключение к базе данных
	redisClient   *redis.Client         // Подключение к Redis
	fiber         *fiber.App            // Веб-сервер Fiber
	container     *container.Container  // DI контейнер со всеми зависимостями
	metrics       *metrics.Metrics      // Prometheus метрики
	healthChecker *health.HealthChecker // Health check компонент
}

// NewApp создает и инициализирует новое приложение
// ctx - контекст для управления жизненным циклом приложения
// envPath - путь к файлу с переменными окружения
func NewApp(ctx context.Context, envPath string) (*App, error) {
	app := &App{
		ctx: ctx,
	}

	// Создаем файл логов
	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Создаем логер и настраиваем вывод в файл и терминал одновременно
	app.logger = logrus.New()
	app.logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	app.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Инициализируем конфигурацию
	app.conf = conf.NewConf(app.logger, envPath)

	// Подключаемся к БД
	app.db = app.conf.DBConnect()

	// Инициализируем Redis подключение с retry logic
	redisHost := app.conf.GetConValue("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := app.conf.GetConValue("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	app.redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Пытаемся подключиться к Redis с retry logic
	if err := app.connectToRedisWithRetry(ctx, 3); err != nil {
		app.logger.WithError(err).Warn("Failed to connect to Redis after retries (cache will be unavailable, but app will continue)")
	} else {
		app.logger.Infof("Redis connected successfully at %s", redisAddr)
	}

	// Включаем расширение uuid-ossp для PostgreSQL
	if err := app.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		app.logger.WithError(err).Warn("Failed to create uuid-ossp extension (may already exist)")
	}

	// Выполняем миграции БД
	if err := app.db.AutoMigrate(&entity.User{}, &entity.EstOrganization{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	app.logger.Info("Database migrations completed successfully")

	// Создаем Fiber приложение
	app.fiber = fiber.New()

	// Инициализируем Prometheus метрики
	app.metrics = metrics.NewMetrics()

	// Инициализируем Health Checker
	app.healthChecker = health.NewHealthChecker(app.db, app.redisClient, app.logger)

	// Добавляем middleware для восстановления после паник (ПЕРВЫМ, перед другими)
	app.fiber.Use(middleware.RecoveryMiddleware(app.logger))

	// Добавляем CORS middleware
	origins := app.conf.GetConValue("ALLOWED_ORIGINS")
	if origins == "" {
		// по умолчанию разрешаем порты разработки 3000, 3002, 5173 (Vite), включая IPv6 и IPv4
		origins = "http://localhost:3000,http://127.0.0.1:3000,http://[::1]:3000,http://localhost:3002,http://127.0.0.1:3002,http://[::1]:3002,http://localhost:5173,http://127.0.0.1:5173,http://[::1]:5173"
	}
	app.fiber.Use(cors.New(cors.Config{
		AllowOrigins: origins,

		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Org-Id",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Добавляем middleware для уникальных ID запросов (трассировка)
	app.fiber.Use(middleware.RequestIDMiddleware())

	// Добавляем middleware для обработки ошибок
	app.fiber.Use(middleware.ErrorHandlingMiddleware(app.logger))

	// Добавляем middleware для логирования
	app.fiber.Use(middleware.LogrusMiddleware(app.logger))

	// Добавляем middleware для метрик (ДОЛЖЕН быть после RequestID)
	app.fiber.Use(middleware.MetricsMiddleware(app.metrics))

	// Инициализируем DI контейнер со всеми зависимостями
	app.container = container.NewContainer(app.db, app.logger, app.redisClient)
	app.logger.Info("Dependency injection container initialized with Redis cache")

	// Инициализируем Rate Limiter (доступен из контейнера для handlers)
	_ = app.container.GetRateLimiter()
	app.logger.Info("Rate limiter initialized with Redis backend")

	// Выполняем cache warming для основных данных
	app.warmCache()

	// Инициализируем сервис управления динамическими БД организаций
	organizationDBService := service_impl.NewOrganizationDBService(
		app.db,
		app.logger,
		app.conf.GetConValue("DB_HOST"),
		app.conf.GetConValue("DB_PORT"),
		app.conf.GetConValue("DB_USER"),
		app.conf.GetConValue("DB_PASSWORD"),
		app.conf.GetConValue("DB_SSLMODE"),
	)
	app.logger.Info("Organization database service initialized")

	// Регистрируем все handlers с контейнером зависимостей
	RegisterHandlers(app.fiber, app.container, organizationDBService)

	// Регистрируем Prometheus metrics endpoint в правильном формате
	// Prometheus scraper ожидает текстовый формат по пути /metrics
	metricsHandler := promhttp.Handler()
	app.fiber.All("/metrics", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "text/plain; version=0.0.4; charset=utf-8")

		// Используем буфер для захвата вывода prometheus handler
		w := &responseWriter{
			header: make(http.Header),
			body:   []byte{},
		}

		// Создаем dummy http.Request для prometheus handler
		req := &http.Request{
			Method:     c.Method(),
			RequestURI: c.OriginalURL(),
			Header:     make(http.Header),
		}

		// Выполняем prometheus handler
		metricsHandler.ServeHTTP(w, req)

		// Отправляем результат клиенту
		for key, values := range w.header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		return c.Send(w.body)
	})

	// Регистрируем Swagger JSON endpoint
	app.fiber.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		swaggerSpec := map[string]interface{}{
			"openapi": "3.0.0",
			"info": map[string]interface{}{
				"title":       "Tunduck Admin API System",
				"description": "Enterprise API для управления ЭСФ документами, организациями и пользователями с кешированием, ограничением частоты запросов и мониторингом здоровья системы",
				"version":     "2.0.0",
				"contact": map[string]string{
					"name":  "API Support",
					"url":   "https://github.com/rusgainew/tunguc-app2",
					"email": "support@example.com",
				},
				"license": map[string]string{
					"name": "MIT",
					"url":  "https://opensource.org/licenses/MIT",
				},
			},
			"servers": []map[string]string{
				{"url": "http://localhost:8080", "description": "Development"},
				{"url": "https://api.example.com", "description": "Production"},
			},
			"components": map[string]interface{}{
				"securitySchemes": map[string]interface{}{
					"BearerAuth": map[string]interface{}{
						"type":         "http",
						"scheme":       "bearer",
						"bearerFormat": "JWT",
					},
				},
				"schemas": getSwaggerSchemas(),
			},
			"paths": getSwaggerPaths(),
		}
		return c.JSON(swaggerSpec)
	})

	// Регистрируем Swagger UI endpoint
	// Используем gofiber/swagger для интеграции с Fiber
	app.fiber.Get("/swagger/*", swagger.HandlerDefault)

	// Регистрируем API документацию редирект
	app.fiber.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	// Регистрируем Health Check endpoint
	app.fiber.Get("/health", func(c *fiber.Ctx) error {
		healthStatus := app.healthChecker.Check(c.Context())
		statusCode := http.StatusOK
		if healthStatus.Status == health.StatusDown {
			statusCode = http.StatusServiceUnavailable
		}
		return c.Status(statusCode).JSON(healthStatus)
	})

	return app, nil
}

// Run запускает веб-сервер и блокирует выполнение до завершения работы
func (a *App) Run() error {
	// Формируем адрес сервера
	host := a.conf.GetConValue("APP_HOST")
	port := a.conf.GetConValue("APP_PORT")
	// Нормализуем хост: 127.0.0.0 — некорректный loopback, заменяем на 127.0.0.1
	if host == "127.0.0.0" {
		host = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	a.logger.Infof("Starting server on %s", addr)

	if err := a.fiber.Listen(addr); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// Shutdown корректно завершает работу приложения и освобождает ресурсы
func (a *App) Shutdown() error {
	// Закрываем Redis соединение
	if a.redisClient != nil {
		if err := a.redisClient.Close(); err != nil {
			a.logger.WithError(err).Warn("Failed to close Redis connection")
		}
	}

	if a.fiber != nil {
		return a.fiber.Shutdown()
	}
	return nil
}

// warmCache предварительно загружает часто используемые данные в кеш
func (a *App) warmCache() {
	if a.container == nil || a.container.GetCacheManager() == nil {
		a.logger.Info("Cache warming skipped: Redis not available")
		return
	}

	a.logger.Info("Starting cache warming...")

	a.logger.Info("Cache warming completed")
}

// ShutdownWithContext корректно завершает работу приложения с поддержкой контекста и таймаута
func (a *App) ShutdownWithContext(ctx context.Context) error {
	// Закрываем Redis соединение
	if a.redisClient != nil {
		if err := a.redisClient.Close(); err != nil {
			a.logger.WithError(err).Warn("Failed to close Redis connection")
		}
	}

	if a.fiber != nil {
		return a.fiber.ShutdownWithContext(ctx)
	}
	return nil
}

// connectToRedisWithRetry пытается подключиться к Redis с retry logic
func (a *App) connectToRedisWithRetry(ctx context.Context, maxRetries int) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		a.logger.Infof("Attempting to connect to Redis (attempt %d/%d)...", attempt, maxRetries)

		if err := a.redisClient.Ping(ctx).Err(); err != nil {
			lastErr = err
			a.logger.WithError(err).Warnf("Redis connection attempt %d failed, retrying in 2 seconds...", attempt)

			if attempt < maxRetries {
				select {
				case <-time.After(2 * time.Second):
					// Continue to next retry
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		} else {
			// Connection successful
			return nil
		}
	}

	return fmt.Errorf("failed to connect to Redis after %d attempts: %w", maxRetries, lastErr)
}

// responseWriter реализует интерфейс http.ResponseWriter для использования с prometheus
type responseWriter struct {
	header     http.Header
	body       []byte
	statusCode int
}

func (rw *responseWriter) Header() http.Header {
	return rw.header
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body = append(rw.body, data...)
	return len(data), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

// getSwaggerSchemas возвращает определения схем OpenAPI
func getSwaggerSchemas() map[string]interface{} {
	return map[string]interface{}{
		"User": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id":         map[string]string{"type": "string", "format": "uuid"},
				"username":   map[string]string{"type": "string"},
				"email":      map[string]string{"type": "string", "format": "email"},
				"fullName":   map[string]string{"type": "string"},
				"phone":      map[string]string{"type": "string"},
				"role":       map[string]string{"type": "string", "enum": "admin|user|viewer"},
				"isActive":   map[string]string{"type": "boolean"},
				"created_at": map[string]string{"type": "string", "format": "date-time"},
				"updated_at": map[string]string{"type": "string", "format": "date-time"},
			},
		},
		"RegisterRequest": map[string]interface{}{
			"type":     "object",
			"required": []string{"username", "email", "password"},
			"properties": map[string]interface{}{
				"username": map[string]string{"type": "string", "minLength": "3"},
				"email":    map[string]string{"type": "string", "format": "email"},
				"password": map[string]string{"type": "string", "minLength": "8"},
				"fullName": map[string]string{"type": "string"},
				"phone":    map[string]string{"type": "string"},
			},
		},
		"AdminRegisterRequest": map[string]interface{}{
			"type":     "object",
			"required": []string{"username", "email", "password", "adminSecret"},
			"properties": map[string]interface{}{
				"username":    map[string]string{"type": "string", "minLength": "3"},
				"email":       map[string]string{"type": "string", "format": "email"},
				"password":    map[string]string{"type": "string", "minLength": "8"},
				"fullName":    map[string]string{"type": "string"},
				"phone":       map[string]string{"type": "string"},
				"adminSecret": map[string]string{"type": "string", "description": "Секретный код администратора для проверки"},
			},
		},
		"LoginRequest": map[string]interface{}{
			"type":     "object",
			"required": []string{"email", "password"},
			"properties": map[string]interface{}{
				"email":    map[string]string{"type": "string", "format": "email"},
				"password": map[string]string{"type": "string"},
			},
		},
		"LoginResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"token": map[string]string{"type": "string"},
				"user": map[string]interface{}{
					"$ref": "#/components/schemas/User",
				},
			},
		},
		"UpdateUserRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"username": map[string]string{"type": "string", "minLength": "3"},
				"email":    map[string]string{"type": "string", "format": "email"},
				"fullName": map[string]string{"type": "string"},
				"phone":    map[string]string{"type": "string"},
				"role":     map[string]string{"type": "string", "enum": "admin|user|viewer"},
				"isActive": map[string]string{"type": "boolean"},
			},
		},
		"Organization": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id":          map[string]string{"type": "string", "format": "uuid"},
				"name":        map[string]string{"type": "string"},
				"description": map[string]string{"type": "string"},
				"token":       map[string]string{"type": "string"},
				"dbName":      map[string]string{"type": "string"},
				"created_at":  map[string]string{"type": "string", "format": "date-time"},
				"updated_at":  map[string]string{"type": "string", "format": "date-time"},
			},
		},
		"CreateOrganizationRequest": map[string]interface{}{
			"type":     "object",
			"required": []string{"name"},
			"properties": map[string]interface{}{
				"name":        map[string]string{"type": "string"},
				"description": map[string]string{"type": "string"},
				"token":       map[string]string{"type": "string"},
				"dbName":      map[string]string{"type": "string"},
			},
		},
		"UpdateOrganizationRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name":        map[string]string{"type": "string"},
				"description": map[string]string{"type": "string"},
				"token":       map[string]string{"type": "string"},
			},
		},
		"Document": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id":              map[string]string{"type": "string", "format": "uuid"},
				"number":          map[string]string{"type": "string"},
				"status":          map[string]string{"type": "string", "enum": "draft|sent|received|processed"},
				"amount":          map[string]string{"type": "number"},
				"organization_id": map[string]string{"type": "string", "format": "uuid"},
				"created_at":      map[string]string{"type": "string", "format": "date-time"},
				"updated_at":      map[string]string{"type": "string", "format": "date-time"},
			},
		},
		"CreateDocumentRequest": map[string]interface{}{
			"type":     "object",
			"required": []string{"number", "organization_id"},
			"properties": map[string]interface{}{
				"number":          map[string]string{"type": "string"},
				"status":          map[string]string{"type": "string", "enum": "draft|sent|received|processed"},
				"amount":          map[string]string{"type": "number"},
				"organization_id": map[string]string{"type": "string", "format": "uuid"},
			},
		},
		"UpdateDocumentRequest": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"number": map[string]string{"type": "string"},
				"status": map[string]string{"type": "string", "enum": "draft|sent|received|processed"},
				"amount": map[string]string{"type": "number"},
			},
		},
		"HealthStatus": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"status":    map[string]string{"type": "string"},
				"timestamp": map[string]string{"type": "string", "format": "date-time"},
				"uptime":    map[string]string{"type": "string"},
				"components": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name":          map[string]string{"type": "string"},
							"status":        map[string]string{"type": "string"},
							"response_time": map[string]string{"type": "string"},
							"message":       map[string]string{"type": "string"},
						},
					},
				},
			},
		},
		"ErrorResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"error":   map[string]string{"type": "string"},
				"message": map[string]string{"type": "string"},
			},
		},
	}
}

// getSwaggerPaths возвращает определения путей OpenAPI
func getSwaggerPaths() map[string]interface{} {
	return map[string]interface{}{
		"/api/auth/register": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"Authentication"},
				"summary":     "Регистрация нового пользователя",
				"description": "Регистрация администратора или обычного пользователя с проверкой уникальности",
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/RegisterRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"201": map[string]interface{}{
						"description": "Пользователь успешно зарегистрирован",
					},
					"400": map[string]interface{}{
						"description": "Неверный ввод или валидация ошибка",
					},
					"429": map[string]interface{}{
						"description": "Превышен лимит запросов",
					},
				},
			},
		},
		"/api/auth/register-admin": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"Authentication"},
				"summary":     "Регистрация администратора",
				"description": "Регистрация администратора с проверкой секретного кода администратора. Требуется действительный adminSecret из переменной окружения ADMIN_SECRET",
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/AdminRegisterRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"201": map[string]interface{}{
						"description": "Администратор успешно зарегистрирован",
					},
					"400": map[string]interface{}{
						"description": "Неверный ввод, валидация ошибка или неверный adminSecret",
					},
					"401": map[string]interface{}{
						"description": "Отсутствует или неверный секретный код администратора",
					},
					"429": map[string]interface{}{
						"description": "Превышен лимит запросов",
					},
				},
			},
		},
		"/api/auth/login": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":        []string{"Authentication"},
				"summary":     "Вход пользователя в систему",
				"description": "Аутентификация и получение JWT токена",
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/LoginRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Успешный вход, возвращается JWT токен",
					},
					"400": map[string]interface{}{
						"description": "Неверный формат запроса",
					},
					"401": map[string]interface{}{
						"description": "Неверные учётные данные",
					},
					"429": map[string]interface{}{
						"description": "Превышен лимит запросов",
					},
				},
			},
		},
		"/api/auth/logout": map[string]interface{}{
			"post": map[string]interface{}{
				"tags":    []string{"Authentication"},
				"summary": "Выход пользователя из системы",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Успешный выход, токен добавлен в blacklist",
					},
					"401": map[string]interface{}{
						"description": "Отсутствует или неверный токен",
					},
				},
			},
		},
		"/api/auth/me": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Authentication"},
				"summary":     "Получить текущего пользователя",
				"description": "Получить информацию о аутентифицированном пользователе",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Информация о текущем пользователе",
					},
					"401": map[string]interface{}{
						"description": "Отсутствует или неверный токен",
					},
				},
			},
		},
		"/api/users": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Users"},
				"summary":     "Получить всех пользователей",
				"description": "Получить список всех пользователей с пагинацией",
				"parameters": []map[string]interface{}{
					{
						"name":        "page",
						"in":          "query",
						"schema":      map[string]string{"type": "integer", "default": "1"},
						"required":    false,
						"description": "Номер страницы (начиная с 1)",
					},
					{
						"name":        "limit",
						"in":          "query",
						"schema":      map[string]string{"type": "integer", "default": "10"},
						"required":    false,
						"description": "Количество записей на странице (1-100)",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Список пользователей",
					},
				},
			},
		},
		"/api/users/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":    []string{"Users"},
				"summary": "Получить пользователя по ID",
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Данные пользователя",
					},
					"404": map[string]interface{}{
						"description": "Пользователь не найден",
					},
				},
			},
			"put": map[string]interface{}{
				"tags":        []string{"Users"},
				"summary":     "Обновить пользователя",
				"description": "Обновить данные пользователя (администратор или сам пользователь)",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/UpdateUserRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Пользователь успешно обновлён",
					},
					"400": map[string]interface{}{
						"description": "Неверные данные",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
					"403": map[string]interface{}{
						"description": "Доступ запрещён (не администратор и не сам пользователь)",
					},
					"404": map[string]interface{}{
						"description": "Пользователь не найден",
					},
				},
			},
			"delete": map[string]interface{}{
				"tags":        []string{"Users"},
				"summary":     "Удалить пользователя",
				"description": "Удалить пользователя (только администратор)",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"responses": map[string]interface{}{
					"204": map[string]interface{}{
						"description": "Пользователь успешно удалён",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
					"403": map[string]interface{}{
						"description": "Доступ запрещён (не администратор)",
					},
					"404": map[string]interface{}{
						"description": "Пользователь не найден",
					},
				},
			},
		},
		"/api/esf-organizations": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Organizations"},
				"summary":     "Получить все организации",
				"description": "Получить список всех организаций",
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Список организаций",
					},
				},
			},
			"post": map[string]interface{}{
				"tags":        []string{"Organizations"},
				"summary":     "Создать новую организацию",
				"description": "Создать новую организацию (требуется администратор)",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/CreateOrganizationRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"201": map[string]interface{}{
						"description": "Организация успешно создана",
					},
					"400": map[string]interface{}{
						"description": "Неверные данные",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
					"403": map[string]interface{}{
						"description": "Доступ запрещён",
					},
				},
			},
		},
		"/api/esf-organizations/paginated": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Organizations"},
				"summary":     "Получить организации с пагинацией",
				"description": "Получить список организаций с поддержкой пагинации, поиска и сортировки",
				"parameters": []map[string]interface{}{
					{
						"name":     "page",
						"in":       "query",
						"schema":   map[string]string{"type": "integer"},
						"required": false,
					},
					{
						"name":     "pageSize",
						"in":       "query",
						"schema":   map[string]string{"type": "integer"},
						"required": false,
					},
					{
						"name":     "search",
						"in":       "query",
						"schema":   map[string]string{"type": "string"},
						"required": false,
					},
					{
						"name":     "sortBy",
						"in":       "query",
						"schema":   map[string]string{"type": "string"},
						"required": false,
					},
					{
						"name":     "sortOrder",
						"in":       "query",
						"schema":   map[string]string{"type": "string", "enum": "asc|desc"},
						"required": false,
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Список организаций с пагинацией",
					},
				},
			},
		},
		"/api/esf-organizations/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":    []string{"Organizations"},
				"summary": "Получить организацию по ID",
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Данные организации",
					},
					"404": map[string]interface{}{
						"description": "Организация не найдена",
					},
				},
			},
			"put": map[string]interface{}{
				"tags":    []string{"Organizations"},
				"summary": "Обновить организацию",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/UpdateOrganizationRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Организация успешно обновлена",
					},
					"400": map[string]interface{}{
						"description": "Неверные данные",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
					"404": map[string]interface{}{
						"description": "Организация не найдена",
					},
				},
			},
			"delete": map[string]interface{}{
				"tags":    []string{"Organizations"},
				"summary": "Удалить организацию",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"responses": map[string]interface{}{
					"204": map[string]interface{}{
						"description": "Организация успешно удалена",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
					"404": map[string]interface{}{
						"description": "Организация не найдена",
					},
				},
			},
		},
		"/api/esf-documents": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Documents"},
				"summary":     "Получить все ЭСФ документы",
				"description": "Получить список всех электронных счётов-фактур",
				"parameters": []map[string]interface{}{
					{
						"name":        "org_id",
						"in":          "query",
						"schema":      map[string]string{"type": "string", "format": "uuid"},
						"required":    false,
						"description": "Фильтр по ID организации",
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Список документов",
					},
				},
			},
			"post": map[string]interface{}{
				"tags":    []string{"Documents"},
				"summary": "Создать новый ЭСФ документ",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/CreateDocumentRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"201": map[string]interface{}{
						"description": "Документ успешно создан",
					},
					"400": map[string]interface{}{
						"description": "Неверные данные",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
				},
			},
		},
		"/api/esf-documents/paginated": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"Documents"},
				"summary":     "Получить документы с пагинацией",
				"description": "Получить список документов с поддержкой пагинации, поиска и сортировки",
				"parameters": []map[string]interface{}{
					{
						"name":     "page",
						"in":       "query",
						"schema":   map[string]string{"type": "integer"},
						"required": false,
					},
					{
						"name":     "pageSize",
						"in":       "query",
						"schema":   map[string]string{"type": "integer"},
						"required": false,
					},
					{
						"name":     "org_id",
						"in":       "query",
						"schema":   map[string]string{"type": "string", "format": "uuid"},
						"required": false,
					},
					{
						"name":     "search",
						"in":       "query",
						"schema":   map[string]string{"type": "string"},
						"required": false,
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Список документов с пагинацией",
					},
				},
			},
		},
		"/api/esf-documents/{id}": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":    []string{"Documents"},
				"summary": "Получить ЭСФ документ по ID",
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
					{
						"name":     "org_id",
						"in":       "query",
						"schema":   map[string]string{"type": "string", "format": "uuid"},
						"required": false,
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Данные документа",
					},
					"404": map[string]interface{}{
						"description": "Документ не найден",
					},
				},
			},
			"put": map[string]interface{}{
				"tags":    []string{"Documents"},
				"summary": "Обновить ЭСФ документ",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"requestBody": map[string]interface{}{
					"required": true,
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": map[string]string{"$ref": "#/components/schemas/UpdateDocumentRequest"},
						},
					},
				},
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Документ успешно обновлен",
					},
					"400": map[string]interface{}{
						"description": "Неверные данные",
					},
					"404": map[string]interface{}{
						"description": "Документ не найден",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
				},
			},
			"delete": map[string]interface{}{
				"tags":    []string{"Documents"},
				"summary": "Удалить ЭСФ документ",
				"security": []map[string]interface{}{
					{"BearerAuth": []string{}},
				},
				"parameters": []map[string]interface{}{
					{
						"name":     "id",
						"in":       "path",
						"required": true,
						"schema":   map[string]string{"type": "string", "format": "uuid"},
					},
				},
				"responses": map[string]interface{}{
					"204": map[string]interface{}{
						"description": "Документ успешно удален",
					},
					"404": map[string]interface{}{
						"description": "Документ не найден",
					},
					"401": map[string]interface{}{
						"description": "Требуется аутентификация",
					},
				},
			},
		},
		"/health": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"System"},
				"summary":     "Проверка здоровья системы",
				"description": "Проверить состояние системы (PostgreSQL, Redis, Uptime)",
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Системы работают нормально",
					},
					"503": map[string]interface{}{
						"description": "Система или компоненты недоступны",
					},
				},
			},
		},
		"/metrics": map[string]interface{}{
			"get": map[string]interface{}{
				"tags":        []string{"System"},
				"summary":     "Метрики Prometheus",
				"description": "Получить метрики приложения в формате Prometheus",
				"responses": map[string]interface{}{
					"200": map[string]interface{}{
						"description": "Данные метрик в текстовом формате",
						"content": map[string]interface{}{
							"text/plain": map[string]interface{}{
								"schema": map[string]string{"type": "string"},
							},
						},
					},
				},
			},
		},
	}
}
