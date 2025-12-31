package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/cache/redis"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/config"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/event/rabbitmq"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/jwt"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/http/server"
)

func main() {
	// 1. Загрузить конфигурацию
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	fmt.Printf("Starting Auth Service\n")
	fmt.Printf("Environment: %s\n", cfg.Environment)
	fmt.Printf("HTTP Port: %d\n", cfg.HttpPort)
	fmt.Printf("gRPC Port: %d\n", cfg.GrpcPort)

	// 2. Подключиться к PostgreSQL
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✓ Connected to PostgreSQL")

	// Инициализировать schema
	if err := postgres.InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("✓ Database schema initialized")

	// 3. Подключиться к Redis
	redisClient := redis.NewClient(redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("✓ Connected to Redis")

	// 4. Подключиться к RabbitMQ
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ channel: %v", err)
	}
	defer ch.Close()
	fmt.Println("✓ Connected to RabbitMQ")

	// 5. Инициализировать repositories и services
	userRepo := postgres.NewUserRepositoryPostgres(db)
	eventPublisher := rabbitmq.NewEventPublisherRabbitMQ(ch)
	tokenBlacklist := redis.NewTokenBlacklistRedis(redisClient)

	// JWT Manager and Token service
	jwtManager := jwt.NewJWTManager(cfg.JwtSecret, cfg.JwtExpires)
	tokenService := service.NewTokenService(jwtManager)

	// Register and Login services
	registerService := service.NewRegisterUserService(userRepo, eventPublisher)
	loginService := service.NewLoginUserService(userRepo, tokenService, eventPublisher)

	// 6. Запустить HTTP и gRPC серверы параллельно
	httpServer := server.NewHTTPServer(
		cfg,
		registerService,
		loginService,
		tokenService,
		userRepo,
	)

	// gRPC server
	grpcAddr := fmt.Sprintf(":%d", cfg.GrpcPort)
	authServiceServer := grpc.NewAuthServiceServer(
		registerService,
		loginService,
		tokenService,
		userRepo,
		tokenBlacklist,
	)
	grpcServer := grpc.NewGRPCServer(grpcAddr, authServiceServer)

	// Start HTTP server in goroutine
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start gRPC server (blocking)
	if err := grpcServer.Start(); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
