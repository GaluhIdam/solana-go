package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"note-api/app/routes"
	"note-api/core/config"
	"note-api/core/database"
	"note-api/core/middleware"
)

func main() {
	configPath := "application.yml"
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	_, err := database.ConnectDB(configPath)
	if err != nil {
		log.Fatalf("❌ Error initializing database: %v", err)
	}
	app := fiber.New()
	app.Use(middleware.RecoverMiddleware())
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.GetAllowOrigins(),
		AllowMethods: config.GetAllowMethods(),
		AllowHeaders: config.GetAllowHeaders(),
	}))
	routes.RegisterRoutes(app, configPath)
	port := config.GlobalConfig.App.Port
	log.Println("🚀 Server running on port " + port)
	log.Fatal(app.Listen(":" + port))
}
