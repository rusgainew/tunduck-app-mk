package grpc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/rusgainew/tunduck-app-mk/gen/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// AuthClient - gRPC клиент для общения с auth-service
type AuthClient struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

// NewAuthClient создает новый gRPC клиент для auth-service
func NewAuthClient(target string) (*AuthClient, error) {
	// Настройки connection pooling и keepalive
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	// Опции подключения
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TODO: использовать TLS в production
		grpc.WithKeepaliveParams(kacp),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024), // 10MB
			grpc.MaxCallSendMsgSize(10*1024*1024), // 10MB
		),
	}

	// Подключение к auth-service
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth-service at %s: %w", target, err)
	}

	// Создание клиента
	client := pb.NewAuthServiceClient(conn)

	return &AuthClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close закрывает соединение с auth-service
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Register вызывает метод регистрации пользователя
func (c *AuthClient) Register(ctx context.Context, email, password, firstName, lastName string) (*pb.AuthResponse, error) {
	req := &pb.RegisterRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	resp, err := c.client.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return resp, nil
}

// Login вызывает метод входа пользователя
func (c *AuthClient) Login(ctx context.Context, email, password string) (*pb.AuthResponse, error) {
	req := &pb.LoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.client.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to login user: %w", err)
	}

	return resp, nil
}

// ValidateToken проверяет валидность токена и возвращает информацию о пользователе
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*pb.User, error) {
	req := &pb.ValidateTokenRequest{
		AccessToken: token,
	}

	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return resp, nil
}

// RefreshToken обновляет токен доступа
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (*pb.Token, error) {
	req := &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := c.client.RefreshToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return resp, nil
}

// Logout выходит из системы (добавляет токен в blacklist)
func (c *AuthClient) Logout(ctx context.Context, token string) error {
	req := &pb.LogoutRequest{
		AccessToken: token,
	}

	_, err := c.client.Logout(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

// GetUser получает информацию о пользователе по ID
func (c *AuthClient) GetUser(ctx context.Context, userID string) (*pb.User, error) {
	req := &pb.GetUserRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp, nil
}
