package shortner

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Routes(router fiber.Router, service ShortUrlService) {
	apirouter := router.Group("/api/v1/urls")
	router.Get("/:url", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		url, err := service.GetUrl(ctx, c.Params("url"))
		if err != nil {
			return response(nil, http.StatusBadRequest, err, c)
		}
		return c.Redirect(url, http.StatusFound)

	})
	apirouter.Post("/", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		var shortUrl ShortUrl
		if err := c.BodyParser(&shortUrl); err != nil {
			return response(nil, http.StatusBadRequest, err, c)
		}
		if shortUrlRes, err := service.CreateUrl(ctx, shortUrl); err != nil {
			return response(nil, http.StatusBadRequest, err, c)
		} else {
			return response(shortUrlRes, http.StatusOK, err, c)
		}
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
