package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"{{.ModuleName}}/app/routes"
	"{{.ModuleName}}/core/config"
	"{{.ModuleName}}/core/database"

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
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.GetAllowOrigins(),
		AllowMethods: config.GetAllowMethods(),
		AllowHeaders: config.GetAllowHeaders(),
	}))

	// Register routes
	routes.RegisterRoutes(app, configPath)

	port := config.GlobalConfig.App.Port
	log.Println("🚀 Server running on port " + port)
	log.Fatal(app.Listen(":" + port))
}
