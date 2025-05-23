package routes

import (
	"github.com/gofiber/fiber/v2"

	"note-api/app/controllers"
	"note-api/app/repositories"
	"note-api/app/services"
	"note-api/core/database"
)

func RegisterRoutes(app *fiber.App, dbPath string) {
	db, err := database.ConnectDB(dbPath)
	if err != nil {
		panic("❌ Failed to connect to database: " + err.Error())
	}

	// Init repository, service, controller
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Group route API
	api := app.Group("/api")

	// Users route group
	user := api.Group("/users")
	user.Post("/", userController.Create)
}
