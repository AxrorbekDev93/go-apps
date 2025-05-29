package handlers

import (
	"database/sql"
	"fmt"
	"go-api/db"
	"go-api/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
	DepoID   int    `json:"depo_id"`
	TabelNum string `json:"tabel_num"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

func RegisterUser(c *fiber.Ctx) error {
	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверные данные"})
	}

	// Хешируем пароль
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	// Сохраняем только username, password и роль. Остальные поля — null или default
	_, err := db.DB.Exec(`
		INSERT INTO users (username, password, role)
		VALUES (?, ?, ?)`,
		input.Username, hash, input.Role,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка регистрации"})
	}

	return c.JSON(fiber.Map{"message": "Регистрация прошла успешно"})
}

func Login(c *fiber.Ctx) error {
	var input AuthInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат"})
	}

	var id int
	var hashedPwd, role string
	var depoID sql.NullInt64
	var isActive sql.NullBool

	err := db.DB.QueryRow(`
		SELECT id, password, role, depo_id, is_active
		FROM users
		WHERE username = ?`, input.Username).
		Scan(&id, &hashedPwd, &role, &depoID, &isActive)

	if err == sql.ErrNoRows {
		return c.Status(401).JSON(fiber.Map{"error": "Пользователь не найден"})
	}
	if err != nil {
		fmt.Println("❌ Scan error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при чтении данных"})
	}

	if !isActive.Valid || !isActive.Bool {
		return c.Status(401).JSON(fiber.Map{"error": "Пользователь не активен"})
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(input.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Неверный пароль"})
	}

	depoIDValue := 0
	if depoID.Valid {
		depoIDValue = int(depoID.Int64)
	}

	token, _ := utils.GenerateToken(id, role, depoIDValue)

	return c.JSON(fiber.Map{"token": token, "role": role})
}
