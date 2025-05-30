package main
import "github.com/AxrorbekDev93/go-apps/db"
import (
	"go-api/handlers"
	"go-api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Разрешить CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Подключение к БД
	db.Connect()

	// Главная страница
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Сервер Go работает!")
	})

	// 🔓 Открытая регистрация и логин
	app.Post("/register", handlers.RegisterUser)
	app.Post("/login", handlers.Login)

	// 🔐 Пользователи (только авторизованные)
	app.Get("/users", middleware.Protect(), handlers.GetUsers)
	app.Patch("/users/:id", middleware.Protect(), handlers.UpdateUserBySuperAdmin)
	app.Patch("/users/:id/status", middleware.Protect(), handlers.UpdateUserStatus)
	app.Get("/users/me", middleware.Protect(), handlers.GetMyProfile)

	app.Get("/locomotives", middleware.Protect(), handlers.GetLocomotives)
	app.Post("/locomotives", middleware.Protect(), handlers.AddLocomotive)
	app.Delete("/locomotives/:id", middleware.Protect(), handlers.DeleteLocomotive)

	app.Get("/diesel-oil", middleware.Protect(), handlers.GetDieselOil)
	app.Post("/diesel-oil", middleware.Protect(), handlers.AddDieselOil)
	app.Delete("/diesel-oil/:id", middleware.Protect(), handlers.DeleteDieselOil)

	// 📋 Депо (получить и добавить)
	app.Get("/depos", handlers.GetDepos)
	app.Post("/depos", middleware.Protect(), handlers.CreateDepo)

	// Запуск сервера
	app.Listen(":4000")
}
