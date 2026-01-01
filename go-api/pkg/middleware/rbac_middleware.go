package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/response"
	"github.com/sirupsen/logrus"
)

// RBACMiddleware создает middleware для проверки ролей (Role-Based Access Control)
// Требует предварительной валидации через JWTAuthMiddleware (проверка через auth-service)
func RBACMiddleware(logger *logrus.Logger, requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")

		// Получаем роль из контекста (установлена JWTAuthMiddleware после проверки через auth-service)
		role := c.Locals("role")
		if role == nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
			}).Warn("RBAC check failed: role not found in context")

			appErr := apperror.New(apperror.ErrUnauthorized, "Authentication required - no role in context")
			return response.Error(c, appErr)
		}

		// Типируем роль
		roleStr, ok := role.(string)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("RBAC check failed: invalid role type")

			appErr := apperror.New(apperror.ErrForbidden, "Invalid role format in context")
			return response.Error(c, appErr)
		}

		// Проверяем, есть ли требуемая роль
		hasRole := false
		for _, requiredRole := range requiredRoles {
			if roleStr == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			email := c.Locals("email")
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"email":      email,
				"user_role":  roleStr,
				"required":   requiredRoles,
			}).Warn("RBAC check failed: insufficient permissions")

			appErr := apperror.New(apperror.ErrForbidden, "Insufficient permissions to access this resource")
			return response.Error(c, appErr)
		}

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       c.Path(),
			"user_role":  roleStr,
		}).Debug("RBAC check passed")

		return c.Next()
	}
}

// AdminOnlyMiddleware создает middleware для проверки, что пользователь администратор
func AdminOnlyMiddleware(logger *logrus.Logger) fiber.Handler {
	return RBACMiddleware(logger, "admin")
}

// AdminOrSelfMiddleware создает middleware для проверки:
// - Пользователь является администратором, ИЛИ
// - Пользователь редактирует свой собственный профиль
// userIDParam - название параметра маршрута с ID пользователя (например, "id")
func AdminOrSelfMiddleware(logger *logrus.Logger, userIDParam string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")

		// Получаем user_id и role из контекста (установлены JWTAuthMiddleware)
		userID := c.Locals("user_id")
		role := c.Locals("role")

		if userID == nil || role == nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("AdminOrSelf check failed: user_id or role not found in context")

			appErr := apperror.New(apperror.ErrUnauthorized, "Authentication required")
			return response.Error(c, appErr)
		}

		// Типируем userID и role
		userIDStr, ok := userID.(string)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("AdminOrSelf check failed: invalid user_id type")

			appErr := apperror.New(apperror.ErrForbidden, "Invalid user_id format in context")
			return response.Error(c, appErr)
		}

		roleStr, ok := role.(string)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("AdminOrSelf check failed: invalid role type")

			appErr := apperror.New(apperror.ErrForbidden, "Invalid role format in context")
			return response.Error(c, appErr)
		}

		// Если администратор - разрешаем доступ
		if roleStr == "admin" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
				"user_role":  roleStr,
			}).Debug("AdminOrSelf check passed: user is admin")

			return c.Next()
		}

		// Иначе проверяем, редактирует ли пользователь свой профиль
		targetUserID := c.Params(userIDParam)
		if userIDStr == targetUserID {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
				"user_id":    userIDStr,
			}).Debug("AdminOrSelf check passed: user editing own profile")

			return c.Next()
		}

		// Доступ запрещён
		email := c.Locals("email")
		logger.WithFields(logrus.Fields{
			"request_id":     requestID,
			"method":         c.Method(),
			"path":           c.Path(),
			"email":          email,
			"user_role":      roleStr,
			"user_id":        userIDStr,
			"target_user_id": targetUserID,
		}).Warn("AdminOrSelf check failed: insufficient permissions")

		appErr := apperror.New(apperror.ErrForbidden, "You can only edit your own profile")
		return response.Error(c, appErr)
	}
}
