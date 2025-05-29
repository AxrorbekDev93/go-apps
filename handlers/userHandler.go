package handlers

import (
	"database/sql"
	"fmt"
	"go-api/db"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Position string `json:"position"`
	DepoID   int    `json:"depo_id"`
	TabelNum string `json:"tabel_num"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
	DepoName string `json:"depo_name"`
}

func GetUsers(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Роль не найдена"})
	}

	depoIDRaw := c.Locals("depo_id")
	var depoID int
	if depoIDInt, ok := depoIDRaw.(int); ok {
		depoID = depoIDInt
	} else {
		depoID = 0
	}

	var rows *sql.Rows
	var err error

	if role == "superadmin" {
		rows, err = db.DB.Query(`
			SELECT u.id, u.username, u.full_name, u.position, u.depo_id, u.tabel_num, u.phone, u.role, u.is_active, d.name
			FROM users u
			LEFT JOIN depos d ON u.depo_id = d.id`)
	} else {
		rows, err = db.DB.Query(`
			SELECT u.id, u.username, u.full_name, u.position, u.depo_id, u.tabel_num, u.phone, u.role, u.is_active, d.name
			FROM users u
			LEFT JOIN depos d ON u.depo_id = d.id
			WHERE u.depo_id = ?`, depoID)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка запроса пользователей"})
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		var fullName, position, tabelNum, phone, depoName sql.NullString
		var depoID sql.NullInt64
		var isActive sql.NullBool

		err := rows.Scan(&u.ID, &u.Username, &fullName, &position, &depoID, &tabelNum, &phone, &u.Role, &isActive, &depoName)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка сканирования данных"})
		}

		u.FullName = nullToStr(fullName)
		u.Position = nullToStr(position)
		u.TabelNum = nullToStr(tabelNum)
		u.Phone = nullToStr(phone)
		u.DepoName = nullToStr(depoName)
		u.DepoID = int(depoID.Int64)
		u.IsActive = isActive.Valid && isActive.Bool

		users = append(users, u)
	}

	return c.JSON(users)
}

func nullToStr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func UpdateUserBySuperAdmin(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Нет прав"})
	}

	id := c.Params("id")
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"` // необязателен
		FullName string `json:"full_name"`
		Position string `json:"position"`
		DepoID   int    `json:"depo_id"`
		TabelNum string `json:"tabel_num"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		IsActive bool   `json:"is_active"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверные данные"})
	}

	// Обновление с хешированием пароля при необходимости
	query := `
		UPDATE users
		SET username = ?, full_name = ?, position = ?, depo_id = ?, tabel_num = ?, phone = ?, role = ?, is_active = ?
		WHERE id = ?`

	_, err := db.DB.Exec(query,
		input.Username,
		input.FullName,
		input.Position,
		input.DepoID,
		input.TabelNum,
		input.Phone,
		input.Role,
		input.IsActive,
		id,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления пользователя"})
	}

	// Отдельно обновим пароль, если он передан
	if input.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		_, err := db.DB.Exec(`UPDATE users SET password = ? WHERE id = ?`, hash, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления пароля"})
		}
	}

	return c.JSON(fiber.Map{"message": "Пользователь обновлён"})
}

func UpdateUserStatus(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "superadmin" {
		return c.Status(403).JSON(fiber.Map{"error": "Нет доступа"})
	}

	id := c.Params("id")
	var input struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат"})
	}

	_, err := db.DB.Exec("UPDATE users SET is_active = ? WHERE id = ?", input.IsActive, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка обновления пользователя"})
	}

	return c.JSON(fiber.Map{"message": "Статус обновлён"})
}

func GetMyProfile(c *fiber.Ctx) error {
	idRaw := c.Locals("user_id")
	if idRaw == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Нет user_id в токене"})
	}
	userID, ok := idRaw.(int)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Невалидный user_id"})
	}

	var (
		user     User
		fullName sql.NullString
		position sql.NullString
		tabelNum sql.NullString
		phone    sql.NullString
		depoName sql.NullString
	)

	err := db.DB.QueryRow(`
		SELECT u.id, u.username, u.full_name, u.position, d.name, u.tabel_num, u.phone, u.role
		FROM users u
		LEFT JOIN depos d ON u.depo_id = d.id
		WHERE u.id = ?`, userID).
		Scan(&user.ID, &user.Username, &fullName, &position, &depoName, &tabelNum, &phone, &user.Role)

	if err != nil {
		fmt.Println("❌ Ошибка в GetMyProfile:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка загрузки профиля"})
	}

	user.FullName = nullToStr(fullName)
	user.Position = nullToStr(position)
	user.TabelNum = nullToStr(tabelNum)
	user.Phone = nullToStr(phone)
	user.DepoName = nullToStr(depoName)

	return c.JSON(user)
}
