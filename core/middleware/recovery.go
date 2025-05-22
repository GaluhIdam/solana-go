package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case *fiber.Error:
					_ = c.Status(e.Code).JSON(fiber.Map{
						"status":  e.Code,
						"message": e.Message,
					})
				case fiber.Map:
					_ = c.Status(fiber.StatusBadRequest).JSON(e)
				default:
					log.Println("[Panic Error]:", r)
					_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status":  fiber.StatusInternalServerError,
						"message": "Internal Server Error",
					})
				}
			}
		}()
		return c.Next()
	}
}
