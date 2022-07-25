package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Routes(router fiber.Router, service AuthService) {

	router.Post("/login", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		type LoginInput struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return response(nil, http.StatusBadRequest, err, c)
		}
		if claims, err := service.NormalLogin(ctx, input.Email, input.Password); err != nil {
			return response(nil, http.StatusUnauthorized, err, c)
		} else {
			if token, err := service.GenerateToken(ctx, claims); err != nil {
				return response(nil, http.StatusInternalServerError, err, c)
			} else {
				return response(map[string]string{"token": token}, http.StatusOK, nil, c)
			}
		}
	})
	router.Post("/register", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		var input RegisterInput
		if err := c.BodyParser(&input); err != nil {
			return response(nil, http.StatusBadRequest, err, c)
		}
		if claims, err := service.RegisterUser(ctx, input); err != nil {
			return response(nil, http.StatusInternalServerError, err, c)
		} else {
			if token, err := service.GenerateToken(ctx, claims); err != nil {
				return response(nil, http.StatusInternalServerError, err, c)
			} else {
				return response(map[string]string{"token": token}, http.StatusOK, nil, c)
			}
		}
	})
	router.Post("/facebook", func(c *fiber.Ctx) error {
		// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		// defer cancel()
		return nil
	})
}

func response(data interface{}, httpStatus int, err error, c *fiber.Ctx) error {
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(map[string]string{
			"error": err.Error(),
		})
	} else {
		if data != nil {
			return c.Status(httpStatus).JSON(data)
		} else {
			c.Status(httpStatus)
			return nil
		}
	}
}
