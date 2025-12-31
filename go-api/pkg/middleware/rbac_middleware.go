package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/response"
	"github.com/sirupsen/logrus"
)

// RBACMiddleware создает middleware для проверки ролей (Role-Based Access Control)
// Требует предварительной валидации JWT токена
func RBACMiddleware(logger *logrus.Logger, requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")

		// Получаем JWT токен из контекста (установлен JWTAuthMiddleware)
		user := c.Locals("user")
		if user == nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
			}).Warn("RBAC check failed: user not found in context")

			appErr := apperror.New(apperror.ErrUnauthorized, "Authentication required")
			return response.Error(c, appErr)
		}

		// Типируем токен
		claims, ok := user.(*jwt.Token)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("RBAC check failed: invalid token type")

			appErr := apperror.New(apperror.ErrForbidden, "Invalid token format")
			return response.Error(c, appErr)
		}

		// Получаем роль из claims
		mapClaims, ok := claims.Claims.(jwt.MapClaims)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("RBAC check failed: invalid claims type")

			appErr := apperror.New(apperror.ErrForbidden, "Invalid token claims")
			return response.Error(c, appErr)
		}

		role, ok := mapClaims["role"].(string)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("RBAC check failed: role not found in claims")

			appErr := apperror.New(apperror.ErrForbidden, "Role not found in token")
			return response.Error(c, appErr)
		}

		// Проверяем, есть ли требуемая роль
		hasRole := false
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			username, _ := mapClaims["username"].(string)
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Method(),
				"path":       c.Path(),
				"username":   username,
				"user_role":  role,
				"required":   requiredRoles,
			}).Warn("RBAC check failed: insufficient permissions")

			appErr := apperror.New(apperror.ErrForbidden, "Insufficient permissions to access this resource")
			return response.Error(c, appErr)
		}

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"path":       c.Path(),
			"user_role":  role,
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

		// Получаем JWT токен из контекста
		user := c.Locals("user")
		if user == nil {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("AdminOrSelf check failed: user not found in context")

			appErr := apperror.New(apperror.ErrUnauthorized, "Authentication required")
			return response.Error(c, appErr)
		}

		// Типируем токен
		claims, ok := user.(*jwt.Token)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("AdminOrSelf check failed: invalid token type")

			appErr := apperror.New(apperror.ErrForbidden, "Invalid token format")
			return response.Error(c, appErr)
		}

		// Получаем роль и ID пользователя из claims
		mapClaims, ok := claims.Claims.(jwt.MapClaims)
		if !ok {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
			}).Warn("AdminOrSelf check failed: invalid claims type")

			appErr := apperror.New(apperror.ErrForbidden, "Invalid token claims")
			return response.Error(c, appErr)
		}

		role, _ := mapClaims["role"].(string)
		userID, _ := mapClaims["sub"].(string) // "sub" - стандартное поле для user ID в JWT

		// Если администратор - разрешаем доступ
		if role == "admin" {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
				"user_role":  role,
			}).Debug("AdminOrSelf check passed: user is admin")

			return c.Next()
		}

		// Иначе проверяем, редактирует ли пользователь свой профиль
		targetUserID := c.Params(userIDParam)
		if userID == targetUserID {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"path":       c.Path(),
				"user_id":    userID,
			}).Debug("AdminOrSelf check passed: user editing own profile")

			return c.Next()
		}

		// Доступ запрещён
		username, _ := mapClaims["username"].(string)
		logger.WithFields(logrus.Fields{
			"request_id":     requestID,
			"method":         c.Method(),
			"path":           c.Path(),
			"username":       username,
			"user_role":      role,
			"user_id":        userID,
			"target_user_id": targetUserID,
		}).Warn("AdminOrSelf check failed: insufficient permissions")

		appErr := apperror.New(apperror.ErrForbidden, "You can only edit your own profile")
		return response.Error(c, appErr)
	}
}
