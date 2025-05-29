package handlers

import (
	"go-api/db"

	"github.com/gofiber/fiber/v2"
)

type Depo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetDepos(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, name FROM depos")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка получения депо"})
	}
	defer rows.Close()

	var depos []Depo
	for rows.Next() {
		var d Depo
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка сканирования"})
		}
		depos = append(depos, d)
	}

	return c.JSON(depos)
}

func CreateDepo(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Только супер-админ может добавлять депо"})
	}

	var depo Depo
	if err := c.BodyParser(&depo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверные данные"})
	}

	_, err := db.DB.Exec("INSERT INTO depos (name) VALUES (?)", depo.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка добавления депо"})
	}

	return c.JSON(fiber.Map{"message": "Депо добавлено успешно"})
}
