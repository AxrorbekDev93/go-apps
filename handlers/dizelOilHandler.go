// 📁 handlers/dieselOilHandler.go
package handlers

import (
	"go-api/db"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Структура данных
type DieselOilInput struct {
	AnalysisDate   string  `json:"analysis_date"`
	RepairType     string  `json:"repair_type"`
	Locomotive     string  `json:"locomotive"`
	Section        string  `json:"section"`
	FlashPoint     float64 `json:"flash_point"`
	Viscosity      float64 `json:"viscosity"`
	Contamination  float64 `json:"contamination"`
	WaterContent   float64 `json:"water_content"`
	Comment        string  `json:"comment"`
	EmployeeNumber string  `json:"employee_number"`
	LastOilDate    string  `json:"last_oil_date"`
}

// 🔽 Логика расчёта заключения
func getConclusion(input DieselOilInput) string {
	if input.FlashPoint > 170 &&
		input.Viscosity >= 11.5 && input.Viscosity <= 16.5 &&
		input.Contamination < 1300 &&
		input.WaterContent < 0.06 {
		return "Яроқли"
	}
	return "Яроқсиз"
}

// ✅ GET /diesel-oil
func GetDieselOil(c *fiber.Ctx) error {
	depoID := c.Locals("depo_id").(int)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	offset := (page - 1) * limit

	query := `SELECT id, analysis_date, repair_type, locomotive, section, flash_point,
		viscosity, contamination, water_content, comment, employee_number,
		last_oil_date, conclusion FROM dizel_oil_teplovoz WHERE depo_id = ?
		ORDER BY analysis_date DESC LIMIT ? OFFSET ?`

	rows, err := db.DB.Query(query, depoID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB query error"})
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var (
			id                                                                                              int
			analysisDate, repairType, locomotive, section, comment, employeeNumber, lastOilDate, conclusion string
			flashPoint, viscosity, contamination, waterContent                                              float64
		)
		err = rows.Scan(&id, &analysisDate, &repairType, &locomotive, &section,
			&flashPoint, &viscosity, &contamination, &waterContent, &comment,
			&employeeNumber, &lastOilDate, &conclusion)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Scan error"})
		}
		results = append(results, map[string]interface{}{
			"id": id, "analysis_date": analysisDate, "repair_type": repairType,
			"locomotive": locomotive, "section": section, "flash_point": flashPoint,
			"viscosity": viscosity, "contamination": contamination, "water_content": waterContent,
			"comment": comment, "employee_number": employeeNumber,
			"last_oil_date": lastOilDate, "conclusion": conclusion,
		})
	}

	var total int
	db.DB.QueryRow("SELECT COUNT(*) FROM dizel_oil_teplovoz WHERE depo_id = ?", depoID).Scan(&total)
	return c.JSON(fiber.Map{"rows": results, "total": total})
}

// ✅ POST /diesel-oil
func AddDieselOil(c *fiber.Ctx) error {
	var input DieselOilInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	depoID := c.Locals("depo_id").(int)
	conclusion := getConclusion(input)

	query := `INSERT INTO dizel_oil_teplovoz (
		analysis_date, repair_type, locomotive, section,
		flash_point, viscosity, contamination, water_content,
		comment, employee_number, last_oil_date, conclusion, depo_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.DB.Exec(query,
		input.AnalysisDate, input.RepairType, input.Locomotive, input.Section,
		input.FlashPoint, input.Viscosity, input.Contamination, input.WaterContent,
		input.Comment, input.EmployeeNumber, input.LastOilDate, conclusion, depoID,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB insert error"})
	}
	return c.JSON(fiber.Map{"message": "Анализ добавлен"})
}

// ✅ DELETE /diesel-oil/:id
func DeleteDieselOil(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := db.DB.Exec("DELETE FROM dizel_oil_teplovoz WHERE id = ?", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка удаления"})
	}
	return c.JSON(fiber.Map{"message": "Удалено"})
}
