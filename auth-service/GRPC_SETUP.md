# gRPC Setup Guide for Auth-Service

## Overview

Auth-Service implements both **HTTP REST** and **gRPC** interfaces for maximum compatibility and performance.

## Quick Start

### Prerequisites

1. **Install protoc** (Protocol Buffer Compiler)

   ```bash
   # macOS
   brew install protobuf

   # Ubuntu/Debian
   sudo apt-get install protobuf-compiler

   # Verify installation
   protoc --version
   ```

2. **Install Go protoc plugins**
   ```bash
   go install github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc@latest
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```

### Compile Proto Files

```bash
cd api/proto
make proto
```

This generates:

- `gen/proto/auth/auth_service.pb.go` - Message types
- `gen/proto/auth/auth_service_grpc.pb.go` - Service interface and client

### Run Auth-Service

```bash
cd auth-service
go run cmd/auth-service/main.go
```

Services will be available on:

- **HTTP**: `http://localhost:8001`
- **gRPC**: `localhost:9001`

## gRPC Service Definition

The `AuthService` provides:

```proto
service AuthService {
  rpc Register(RegisterRequest) returns (AuthResponse);
  rpc Login(LoginRequest) returns (AuthResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (User);
  rpc GetUser(GetUserRequest) returns (User);
  rpc Logout(LogoutRequest) returns (api.common.Empty);
  rpc RefreshToken(RefreshTokenRequest) returns (Token);
}
```

## Testing gRPC with grpcurl

### Install grpcurl

```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### List services

```bash
grpcurl -plaintext localhost:9001 list
grpcurl -plaintext localhost:9001 list api.auth.AuthService
```

### Register user

```bash
grpcurl -plaintext -d '{"email":"user@example.com","password":"SecurePassword123","first_name":"John","last_name":"Doe"}' \
  localhost:9001 api.auth.AuthService/Register
```

### Login

```bash
grpcurl -plaintext -d '{"email":"user@example.com","password":"SecurePassword123"}' \
  localhost:9001 api.auth.AuthService/Login
```

### Validate token

```bash
grpcurl -plaintext -d '{"access_token":"<your_token_here>"}' \
  localhost:9001 api.auth.AuthService/ValidateToken
```

### Get user

```bash
grpcurl -plaintext -d '{"user_id":"<user_id>"}' \
  localhost:9001 api.auth.AuthService/GetUser
```

### Logout

```bash
grpcurl -plaintext -d '{"access_token":"<your_token_here>"}' \
  localhost:9001 api.auth.AuthService/Logout
```

### Refresh token

```bash
grpcurl -plaintext -d '{"refresh_token":"<your_refresh_token>"}' \
  localhost:9001 api.auth.AuthService/RefreshToken
```

## Calling Auth-Service from Other Services

### Company-Service (Go example)

```go
package main

import (
	"context"
	pb "github.com/rusgainew/tunduck-app-mk/gen/proto/auth"
	"google.golang.org/grpc"
)

func main() {
	// Connect to auth-service
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Validate token
	resp, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		AccessToken: "eyJhbGciOiJIUzI1NiIs...",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("User: %v", resp)
}
```

### API Gateway (reverse proxy)

Map gRPC calls to HTTP endpoints via gateway:

```go
// In api-gateway service
import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

func main() {
	// Create gRPC client
	conn, _ := grpc.Dial("auth-service:9001", grpc.WithInsecure())
	gwmux := runtime.NewServeMux()

	// Register handlers
	pb.RegisterAuthServiceHandlerClient(context.Background(), gwmux, conn)

	// HTTP on 8080
	http.ListenAndServe(":8080", gwmux)
}
```

## Architecture

### HTTP vs gRPC

| Aspect    | HTTP REST     | gRPC                  |
| --------- | ------------- | --------------------- |
| Speed     | Slower (JSON) | 10x faster (Protobuf) |
| Payload   | Larger (JSON) | Smaller (binary)      |
| Browser   | ✅ Native     | ❌ Need gateway       |
| Bandwidth | More          | Less                  |
| Mobile    | Good          | Excellent             |
| Protocol  | HTTP/1.1      | HTTP/2                |

**Use HTTP for**: Web browsers, public APIs, legacy systems
**Use gRPC for**: Service-to-service, real-time, high-throughput

## Implementation Details

### Current Status

✅ **Completed**:

- AuthServiceServer with all 6 RPC methods
- Business logic integration
- Error handling with gRPC status codes
- Token validation and blacklisting

⏳ **Requires protoc**:

- Proto file compilation (generate .pb.go files)
- Protobuf message types

### How It Works

1. **Proto Definition** (`api/proto/auth_service.proto`)

   - Defines service interface and messages
   - Language-agnostic contract

2. **Code Generation** (`make proto`)

   - Generates Go interfaces and types
   - Creates gRPC server/client stubs

3. **Implementation** (`AuthServiceServer`)

   - Implements generated interface methods
   - Contains business logic

4. **Registration** (`RegisterAuthService`)

   - Registers service with gRPC server
   - Routes incoming RPC calls to handlers

5. **Server Start** (`main.go`)
   - Listens on port 9001
   - Accepts gRPC connections

## Running Tests

```bash
# Unit tests
go test ./internal/application/service -v

# All tests
go test ./... -v

# With coverage
go test ./... -cover
```

## Next Steps

1. ✅ Implement gRPC handlers (DONE)
2. ⏳ Compile proto files (requires protoc installation)
3. ⏳ Integration with other services
4. ⏳ Add gRPC middleware (auth, logging)
5. ⏳ Deploy with Docker

## Troubleshooting

### `proto.RegisterAuthServiceServer not found`

- Solution: Run `cd api/proto && make proto` to generate code

### gRPC connection refused

- Check: `auth-service` is running on port 9001
- Check: No firewall blocking port 9001

### Type not found errors

- Solution: Run proto compilation step again

## Resources

- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [gRPC-Gateway](https://grpc-ecosystem.github.io/grpc-gateway/) (HTTP gateway)

## Additional Configuration

### TLS/SSL for Production

```go
creds, _ := credentials.NewServerTLSFromFile("server.crt", "server.key")
server := grpc.NewServer(grpc.Creds(creds))
```

### Interceptors (middleware)

```go
server := grpc.NewServer(
	grpc.UnaryInterceptor(authInterceptor),
	grpc.UnaryInterceptor(loggingInterceptor),
)
```

---

**Status**: ✅ Implementation complete, ⏳ Requires protoc for compilation
