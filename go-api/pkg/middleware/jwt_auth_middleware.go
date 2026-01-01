package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/response"
	"github.com/sirupsen/logrus"
)

// JWTAuthMiddleware создает улучшенный middleware для проверки JWT токенов
// через auth-service (gRPC). Это обеспечивает единый источник истины для токенов
func JWTAuthMiddleware(authProxyService services.AuthProxyService, logger *logrus.Logger) fiber.Handler {
	if authProxyService == nil {
		logger.Fatal("AuthProxyService is required for JWT middleware")
		return func(c *fiber.Ctx) error {
			return response.Error(c, apperror.New(apperror.ErrInternal, "JWT middleware configuration error"))
		}
	}

	return func(c *fiber.Ctx) error {
		// Извлекаем токен из заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			requestID := c.Get("X-Request-ID")
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("Authorization header is missing")

			appErr := apperror.New(apperror.ErrUnauthorized, "Authorization header is required")
			return response.Error(c, appErr)
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			requestID := c.Get("X-Request-ID")
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("Invalid Authorization header format")

			appErr := apperror.New(apperror.ErrInvalidToken, "Invalid authorization header format")
			return response.Error(c, appErr)
		}

		token := parts[1]

		// Валидируем токен через auth-service
		user, err := authProxyService.ValidateToken(c.Context(), token)
		if err != nil {
			requestID := c.Get("X-Request-ID")
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"error":      err.Error(),
			}).Warn("JWT validation failed via auth-service")

			// Если это уже apperror, возвращаем как есть
			if appErr, ok := err.(*apperror.AppError); ok {
				return response.Error(c, appErr)
			}

			// Иначе создаем новую ошибку
			appErr := apperror.New(apperror.ErrInvalidToken, "Invalid or expired JWT token")
			return response.Error(c, appErr)
		}

		// Сохраняем информацию о пользователе в контексте
		c.Locals("user", user)
		c.Locals("user_id", user.ID)
		c.Locals("email", user.Email)
		c.Locals("role", user.Role)

		requestID := c.Get("X-Request-ID")
		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"path":       c.Path(),
		}).Debug("JWT validated successfully via auth-service")

		return c.Next()
	}
}

// JWTOptionalMiddleware создает опциональный JWT middleware
// Если токен присутствует - валидирует его через auth-service, если нет - продолжает выполнение
func JWTOptionalMiddleware(authProxyService services.AuthProxyService, logger *logrus.Logger) fiber.Handler {
	if authProxyService == nil {
		logger.Fatal("AuthProxyService is required for JWT middleware")
		return func(c *fiber.Ctx) error {
			return response.Error(c, apperror.New(apperror.ErrInternal, "JWT middleware configuration error"))
		}
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		requestID := c.Get("X-Request-ID")

		// Если заголовок Authorization отсутствует, продолжаем без токена
		if authHeader == "" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Debug("No Authorization header provided, proceeding as anonymous")
			return c.Next()
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("Invalid Authorization header format in optional auth")

			appErr := apperror.New(apperror.ErrInvalidToken, "Invalid authorization header format")
			return response.Error(c, appErr)
		}

		token := parts[1]

		// Валидируем токен через auth-service
		user, err := authProxyService.ValidateToken(c.Context(), token)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"error":      err.Error(),
			}).Warn("JWT validation failed for optional auth")

			// Если это уже apperror, возвращаем как есть
			if appErr, ok := err.(*apperror.AppError); ok {
				return response.Error(c, appErr)
			}

			// Иначе создаем новую ошибку
			appErr := apperror.New(apperror.ErrInvalidToken, "Invalid or expired JWT token")
			return response.Error(c, appErr)
		}

		// Сохраняем информацию о пользователе в контексте
		c.Locals("user", user)
		c.Locals("user_id", user.ID)
		c.Locals("email", user.Email)
		c.Locals("role", user.Role)

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"path":       c.Path(),
		}).Debug("JWT validated successfully via auth-service (optional)")

		return c.Next()
	}
}

// GetUserIDFromContext извлекает ID пользователя из контекста
// Теперь использует models.UserInfo вместо JWT claims
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userID := c.Locals("user_id")
	if userID == nil {
		return uuid.UUID{}, apperror.New(apperror.ErrUnauthorized, "No user_id in context")
	}

	// После проверки через auth-service, user_id уже строка
	userIDStr, ok := userID.(string)
	if !ok {
		return uuid.UUID{}, apperror.New(apperror.ErrInvalidToken, fmt.Sprintf("Invalid user_id format in context, got type: %T", userID))
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, apperror.New(apperror.ErrInvalidToken, fmt.Sprintf("Invalid user_id UUID format: %v", err))
	}

	return parsedID, nil
}

// GetEmailFromContext извлекает email из контекста
func GetEmailFromContext(c *fiber.Ctx) (string, error) {
	email := c.Locals("email")
	if email == nil {
		return "", apperror.New(apperror.ErrUnauthorized, "No email in context")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", apperror.New(apperror.ErrInvalidToken, "Invalid email format in context")
	}

	return emailStr, nil
}

// GetRoleFromContext извлекает role из контекста
func GetRoleFromContext(c *fiber.Ctx) (string, error) {
	role := c.Locals("role")
	if role == nil {
		return "", apperror.New(apperror.ErrUnauthorized, "No role in context")
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", apperror.New(apperror.ErrInvalidToken, "Invalid role format in context")
	}

	return roleStr, nil
}

// GetUsernameFromContext извлекает username из контекста
// DEPRECATED: Используйте GetEmailFromContext - auth-service использует email для аутентификации
func GetUsernameFromContext(c *fiber.Ctx) (string, error) {
	// Возвращаем email для обратной совместимости
	return GetEmailFromContext(c)
}

// GetClaimsFromContext извлекает информацию о пользователе из контекста
// DEPRECATED: Используйте c.Locals("user") для получения models.UserInfo
func GetClaimsFromContext(c *fiber.Ctx) (map[string]interface{}, error) {
	user := c.Locals("user")
	if user == nil {
		return nil, apperror.New(apperror.ErrUnauthorized, "No user in context")
	}

	// Конвертируем в map для обратной совместимости
	result := make(map[string]interface{})
	result["user_id"] = c.Locals("user_id")
	result["email"] = c.Locals("email")

	return result, nil
}
