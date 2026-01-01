package controllers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/rusgainew/tunduck-app/internal/services"
	"github.com/rusgainew/tunduck-app/pkg/apperror"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/logger"
	"github.com/rusgainew/tunduck-app/pkg/middleware"
	"github.com/rusgainew/tunduck-app/pkg/rbac"
)

type UserController struct {
	logger           *logger.Logger
	db               *gorm.DB
	logrus           *logrus.Logger
	authProxyService services.AuthProxyService
}

func NewUserController(app *fiber.App, authProxyService services.AuthProxyService, log *logrus.Logger, db *gorm.DB) {
	controller := &UserController{
		logger:           logger.New(log),
		db:               db,
		logrus:           log,
		authProxyService: authProxyService,
	}

	controller.logger.Info(context.Background(), "UserController initialized", logrus.Fields{})
	controller.registerRoutes(app)
}

func (c *UserController) registerRoutes(app *fiber.App) {
	userGroup := app.Group("/api/users")

	// Публичные GET routes (без JWT)
	userGroup.Get("/", c.getAllUsers)
	userGroup.Get("/:id", c.getUserByID)

	// Защищённые routes для администраторов
	// Для обновления: администратор или редактирование собственного профиля
	userGroup.Put("/:id", middleware.JWTAuthMiddleware(c.authProxyService, c.logrus), middleware.AdminOrSelfMiddleware(c.logrus, "id"), c.updateUser)

	// Для удаления: только администраторы
	userGroup.Delete("/:id", middleware.JWTAuthMiddleware(c.authProxyService, c.logrus), middleware.AdminOnlyMiddleware(c.logrus), c.deleteUser)
}

// getAllUsers возвращает всех пользователей с пагинацией
func (c *UserController) getAllUsers(ctx *fiber.Ctx) error {
	c.logger.Info(ctx.Context(), "Fetching all users", logrus.Fields{})

	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []entity.User
	var total int64

	// Получаем общее количество пользователей
	if err := c.db.Model(&entity.User{}).Count(&total).Error; err != nil {
		c.logger.Error(ctx.Context(), "Failed to count users", err, logrus.Fields{})
		appErr := apperror.New(apperror.ErrInternal, "failed to count users").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Получаем пользователей с пагинацией
	if err := c.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.logger.Error(ctx.Context(), "Failed to fetch users", err, logrus.Fields{})
		appErr := apperror.New(apperror.ErrInternal, "failed to fetch users").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// getUserByID возвращает пользователя по ID
func (c *UserController) getUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Fetching user by ID", logrus.Fields{"id": id})

	var user entity.User
	if err := c.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			appErr := apperror.New(apperror.ErrNotFound, "user not found")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}
		c.logger.Error(ctx.Context(), "Failed to fetch user", err, logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInternal, "failed to fetch user").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

// updateUser обновляет пользователя (только администратор или пользователь редактирует себя)
func (c *UserController) updateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Updating user", logrus.Fields{"id": id})

	var updateData struct {
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		IsActive bool   `json:"isActive"`
	}

	if err := ctx.BodyParser(&updateData); err != nil {
		c.logger.Error(ctx.Context(), "Failed to parse request body", err, logrus.Fields{})
		appErr := apperror.New(apperror.ErrInvalidRequest, "invalid request body")
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Проверяем, существует ли пользователь
	var user entity.User
	if err := c.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			appErr := apperror.New(apperror.ErrNotFound, "user not found")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}
		c.logger.Error(ctx.Context(), "Failed to fetch user", err, logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInternal, "failed to fetch user").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Если email изменилась, проверяем её уникальность
	if updateData.Email != "" && updateData.Email != user.Email {
		var count int64
		if err := c.db.Model(&entity.User{}).Where("email = ?", updateData.Email).Count(&count).Error; err != nil {
			c.logger.Error(ctx.Context(), "Failed to check email uniqueness", err, logrus.Fields{"email": updateData.Email})
			appErr := apperror.New(apperror.ErrInternal, "failed to verify email")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}
		if count > 0 {
			appErr := apperror.New(apperror.ErrEmailExists, "email already in use")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}
		user.Email = updateData.Email
	}

	// Обновляем прочие поля
	if updateData.FullName != "" {
		user.FullName = updateData.FullName
	}
	if updateData.Phone != "" {
		user.Phone = updateData.Phone
	}
	if updateData.Role != "" {
		user.Role = rbac.Role(updateData.Role)
	}
	user.IsActive = updateData.IsActive

	// Сохраняем обновления
	if err := c.db.Save(&user).Error; err != nil {
		c.logger.Error(ctx.Context(), "Failed to update user", err, logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInternal, "failed to update user").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "User updated successfully", logrus.Fields{"id": id})
	return ctx.Status(http.StatusOK).JSON(user)
}

// deleteUser удаляет пользователя (только администратор)
func (c *UserController) deleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c.logger.Info(ctx.Context(), "Deleting user", logrus.Fields{"id": id})

	// Проверяем, существует ли пользователь
	var user entity.User
	if err := c.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			appErr := apperror.New(apperror.ErrNotFound, "user not found")
			return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
		}
		c.logger.Error(ctx.Context(), "Failed to fetch user", err, logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInternal, "failed to fetch user").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	// Удаляем пользователя (GORM soft delete)
	if err := c.db.Delete(&user).Error; err != nil {
		c.logger.Error(ctx.Context(), "Failed to delete user", err, logrus.Fields{"id": id})
		appErr := apperror.New(apperror.ErrInternal, "failed to delete user").WithError(err)
		return ctx.Status(appErr.HTTPStatus).JSON(appErr.ToResponse())
	}

	c.logger.Info(ctx.Context(), "User deleted successfully", logrus.Fields{"id": id})
	return ctx.SendStatus(http.StatusNoContent)
}
