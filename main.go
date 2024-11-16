package main

import (
	"log"
	"transactions-manager/app/config"
	"transactions-manager/app/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.LoadConfig()

	database.Init()
	defer database.Close()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("¡Hello, world from Fiber!")
	})

	log.Printf("Iniciando servidor en el puerto %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
