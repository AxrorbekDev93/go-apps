package middleware

import (
	"go-api/utils"

	"github.com/gofiber/fiber/v2"
)

func Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Нет токена"})
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Неверный токен"})
		}

		c.Locals("user_id", int(claims["user_id"].(float64)))
		c.Locals("role", claims["role"])
		c.Locals("depo_id", int(claims["depo_id"].(float64)))

		return c.Next()
	}
}
