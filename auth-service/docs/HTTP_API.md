# HTTP API Documentation

## Endpoints

### 1. Health Check

**Endpoint:** `GET /health`

**Authentication:** Not required

**Response (200 OK):**

```json
{
  "status": "ok",
  "service": "auth-service"
}
```

---

### 2. Register

**Endpoint:** `POST /api/v1/auth/register`

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "SecurePass123"
}
```

**Validation Rules:**

- Email: valid email format
- Name: minimum 2 characters
- Password: minimum 8 characters, must contain uppercase, lowercase, and numbers

**Response (201 Created):**

```json
{
  "id": "user_1234567890",
  "email": "user@example.com",
  "name": "John Doe",
  "status": "active",
  "created_at": "2024-01-01T12:00:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid email/password format
- `409 Conflict` - User with this email already exists
- `500 Internal Server Error` - Internal error

---

### 3. Login

**Endpoint:** `POST /api/v1/auth/login`

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

**Response (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Invalid email or password
- `403 Forbidden` - User account is inactive or blocked
- `500 Internal Server Error` - Internal error

---

### 4. Get Profile (Me)

**Endpoint:** `GET /api/v1/auth/me`

**Authentication:** Required (Bearer Token)

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (200 OK):**

```json
{
  "id": "user_1234567890",
  "email": "user@example.com",
  "name": "John Doe",
  "status": "active",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "last_login": "2024-01-01T15:30:00Z"
}
```

**Error Responses:**

- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - User not found
- `500 Internal Server Error` - Internal error

---

### 5. Logout

**Endpoint:** `POST /api/v1/auth/logout`

**Authentication:** Required (Bearer Token)

**Headers:**

```
Authorization: Bearer <access_token>
```

**Response (200 OK):**

```json
{
  "message": "Successfully logged out"
}
```

**Error Responses:**

- `401 Unauthorized` - Missing or invalid token
- `500 Internal Server Error` - Internal error

---

### 6. Refresh Token

**Endpoint:** `POST /api/v1/auth/refresh`

**Authentication:** Not required

**Request Body:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Invalid or revoked refresh token
- `403 Forbidden` - User account is inactive
- `500 Internal Server Error` - Internal error

---

## Error Response Format

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message"
}
```

### Error Codes

| Code                | HTTP Status | Description                         |
| ------------------- | ----------- | ----------------------------------- |
| UNAUTHORIZED        | 401         | Missing or invalid token            |
| USER_NOT_FOUND      | 404         | User not found                      |
| USER_EXISTS         | 409         | User with this email already exists |
| USER_BLOCKED        | 403         | User account is blocked             |
| USER_INACTIVE       | 403         | User account is not active          |
| INVALID_CREDENTIALS | 401         | Invalid email or password           |
| INVALID_EMAIL       | 400         | Invalid email format                |
| INVALID_PASSWORD    | 400         | Invalid password format             |
| PASSWORD_TOO_SHORT  | 400         | Password too short (min 8 chars)    |
| PASSWORD_WEAK       | 400         | Password too weak                   |
| TOKEN_EXPIRED       | 401         | Token has expired                   |
| TOKEN_INVALID       | 401         | Token is invalid                    |
| TOKEN_REVOKED       | 401         | Token has been revoked              |
| INTERNAL_ERROR      | 500         | Internal server error               |

---

## Middleware

### Authentication Middleware

- Validates JWT token from Authorization header
- Extracts userID and token from JWT claims
- Adds userID and token to request context
- Skips validation for: `/health`, `/api/v1/auth/register`, `/api/v1/auth/login`

### CORS Middleware

- Allows all origins
- Allows GET, POST, PUT, DELETE, OPTIONS, PATCH methods
- Handles preflight requests

### Logging Middleware

- Logs all requests and responses
- Includes method, path, status code, and duration

### Recovery Middleware

- Catches panics and returns 500 error
- Prevents server crash

---

## Examples

### Register a new user

```bash
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "password": "SecurePass123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123"
  }'
```

### Get profile

```bash
curl -X GET http://localhost:8001/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

### Logout

```bash
curl -X POST http://localhost:8001/api/v1/auth/logout \
  -H "Authorization: Bearer <access_token>"
```

### Refresh token

```bash
curl -X POST http://localhost:8001/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token>"
  }'
```
