package grpc

import (
	"google.golang.org/grpc"
)

// RegisterAuthService - Register the AuthServiceServer with the gRPC server
// NOTE: This is a simplified implementation. In a real setup:
// 1. Run: cd api/proto && make proto
// 2. This generates: gen/proto/auth/auth_service_grpc.pb.go
// 3. Implement: type AuthServiceServer interface { ... } methods
// 4. Use: pb.RegisterAuthServiceServer(grpcServer, handler)
//
// For now, we provide the business logic that would be called by the generated code.
func RegisterAuthService(server *grpc.Server, handler *AuthServiceServer) {
	// In production, this would be:
	// pb.RegisterAuthServiceServer(server, handler)
	// Where pb is the generated proto package

	// The handler implements:
	// - Register(ctx, req) -> (*AuthResponse, error)
	// - Login(ctx, req) -> (*AuthResponse, error)
	// - ValidateToken(ctx, req) -> (*User, error)
	// - GetUser(ctx, req) -> (*User, error)
	// - Logout(ctx, req) -> (*Empty, error)
	// - RefreshToken(ctx, req) -> (*Token, error)
	//
	// To use this in production:
	// 1. Install protoc: https://grpc.io/docs/protoc-installation/
	// 2. cd api/proto && make proto
	// 3. The generated code will handle registration automatically
}
