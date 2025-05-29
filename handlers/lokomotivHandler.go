package handlers

import (
	"go-api/db"

	"github.com/gofiber/fiber/v2"
)

type Locomotive struct {
	ID     int    `json:"id"`
	Model  string `json:"model"`
	Number string `json:"number"`
	Depo   string `json:"depo"`
}

// GET /locomotives
func GetLocomotives(c *fiber.Ctx) error {
	depoID, ok := c.Locals("depo_id").(int)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Не найден depo_id"})
	}

	rows, err := db.DB.Query(`
		SELECT l.id, l.model, l.number, d.name 
		FROM locomotives l
		LEFT JOIN depos d ON l.depo_id = d.id
		WHERE l.depo_id = ?
		ORDER BY l.id DESC
	`, depoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при получении локомотивов"})
	}
	defer rows.Close()

	var locos []Locomotive
	for rows.Next() {
		var l Locomotive
		err := rows.Scan(&l.ID, &l.Model, &l.Number, &l.Depo)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка чтения данных"})
		}
		locos = append(locos, l)
	}

	return c.JSON(locos)
}

// POST /locomotives
func AddLocomotive(c *fiber.Ctx) error {
	depoID, ok := c.Locals("depo_id").(int)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Не найден depo_id"})
	}

	var input struct {
		Model  string `json:"model"`
		Number string `json:"number"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный JSON"})
	}

	if input.Model == "" || input.Number == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Заполните все поля"})
	}

	_, err := db.DB.Exec(`INSERT INTO locomotives (model, number, depo_id) VALUES (?, ?, ?)`,
		input.Model, input.Number, depoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при добавлении локомотива"})
	}

	return c.JSON(fiber.Map{"message": "Локомотив добавлен"})
}

// DELETE /locomotives/:id
func DeleteLocomotive(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.DB.Exec("DELETE FROM locomotives WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления"})
	}
	return c.JSON(fiber.Map{"message": "Локомотив удалён"})
}
