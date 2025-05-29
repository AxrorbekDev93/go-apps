package handlers

import (
	"database/sql"
	"go-api/db"

	"github.com/gofiber/fiber/v2"
)

func GetTabelByNumber(c *fiber.Ctx) error {
	tabelNum := c.Params("tabel_num")

	var fullName, position, phone string
	err := db.DB.QueryRow(`
		SELECT full_name, position, phone
		FROM tabels
		WHERE tabel_num = ?
	`, tabelNum).Scan(&fullName, &position, &phone)

	if err == sql.ErrNoRows {
		return c.Status(404).JSON(fiber.Map{"error": "Сотрудник не найден"})
	}
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка поиска"})
	}

	return c.JSON(fiber.Map{
		"full_name": fullName,
		"position":  position,
		"phone":     phone,
	})
}
