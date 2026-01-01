package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	// Infrastructure
	redisInfra "github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/cache/redis"
	rabbitmqInfra "github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/event/rabbitmq"
	jwtInfra "github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/jwt"
	postgresInfra "github.com/rusgainew/tunduck-app-mk/auth-service/internal/infrastructure/persistence/postgres"

	// Application
	authService "github.com/rusgainew/tunduck-app-mk/auth-service/internal/application/service"

	// Interfaces
	authgrpc "github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/grpc"
	httprouter "github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/http"
	"github.com/rusgainew/tunduck-app-mk/auth-service/internal/interfaces/http/server"
)

func main() {
	// Load configuration from environment
	config := &Config{
		PostgreSQL: postgresInfra.DefaultConfig(),
		Redis:      redisInfra.DefaultConfig(),
		RabbitMQ:   rabbitmqInfra.DefaultConfig(),
		Server: ServerConfig{
			Port:     getEnv("AUTH_SERVICE_PORT", "8001"),
			GrpcPort: getEnv("AUTH_SERVICE_GRPC_PORT", "9001"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessTTL:  3600,   // 1 hour
			RefreshTTL: 604800, // 7 days
		},
	}

	// Initialize database
	db, err := postgresInfra.NewConnection(config.PostgreSQL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
	log.Println("✓ Connected to PostgreSQL")

	// Initialize Redis
	redisClient, err := redisInfra.NewConnection(config.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Test Redis connection
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}
	log.Printf("✓ Connected to Redis: %s", pong)

	// Initialize RabbitMQ
	amqpConn, err := rabbitmqInfra.NewConnection(config.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	// Create RabbitMQ channel
	ch, err := rabbitmqInfra.NewChannel(amqpConn)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ channel: %v", err)
	}
	defer ch.Close()

	// Declare exchange
	if err := rabbitmqInfra.DeclareExchange(ch, "tunduck.auth"); err != nil {
		log.Fatalf("Failed to declare RabbitMQ exchange: %v", err)
	}
	log.Println("✓ Connected to RabbitMQ")

	// Initialize repositories
	userRepository := postgresInfra.NewUserRepositoryPostgres(db)

	// Initialize infrastructure implementations
	eventPublisher := rabbitmqInfra.NewEventPublisherRabbitMQ(ch)
	tokenBlacklist := redisInfra.NewTokenBlacklistRedis(redisClient)

	// Initialize JWT manager
	jwtManager := jwtInfra.NewJWTManager(config.JWT.Secret, int64(config.JWT.AccessTTL))

	// Token service built on JWT manager
	tokenService := authService.NewTokenService(jwtManager)

	// Initialize services
	registerService := authService.NewRegisterUserService(
		userRepository,
		eventPublisher,
	)

	loginService := authService.NewLoginUserService(userRepository, tokenService, eventPublisher)

	getUserService := authService.NewGetUserService(userRepository)

	logoutService := authService.NewLogoutUserService(
		tokenBlacklist,
		eventPublisher,
	)

	refreshTokenService := authService.NewRefreshTokenService(userRepository, tokenService, tokenBlacklist)

	validateTokenService := authService.NewValidateTokenService(tokenService, tokenBlacklist, userRepository)

	// Initialize gRPC handler
	grpcHandler := authgrpc.NewAuthServiceServer(
		registerService,
		loginService,
		tokenService,
		userRepository,
		tokenBlacklist,
	)

	// Initialize HTTP router
	httpRouter := httprouter.NewRouter(
		registerService,
		loginService,
		validateTokenService,
		getUserService,
		logoutService,
		refreshTokenService,
	)

	// Initialize HTTP server
	httpServer := server.NewHTTPServer(
		config.Server.Port,
		httpRouter,
	)

	// Initialize gRPC server
	grpcAddr := ":" + config.Server.GrpcPort
	grpcServer := authgrpc.NewGRPCServer(grpcAddr, grpcHandler)

	// Start HTTP server in a goroutine
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Start gRPC server in a goroutine
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("\nShutting down gracefully...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 30*1000000000) // 30 seconds
	defer cancel()

	grpcServer.Stop()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("✓ Server stopped")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Config holds application configuration
type Config struct {
	PostgreSQL postgresInfra.Config
	Redis      redisInfra.Config
	RabbitMQ   rabbitmqInfra.Config
	Server     ServerConfig
	JWT        JWTConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port     string
	GrpcPort string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	AccessTTL  int
	RefreshTTL int
}
