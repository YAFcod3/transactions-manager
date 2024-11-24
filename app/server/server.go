package server

import (
	"log"
	"os"
	"transactions-manager/app/middleware"
	"transactions-manager/app/routes"
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CreateServer(appServices *services.AppServices, codeGen *generate_transaction_code.CodeGenerator, port string) *fiber.App {
	app := fiber.New()

	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "*"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowCredentials: false,
		AllowMethods:     "GET, POST, PUT,PATCH, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Authorization, Accept",
	}))

	app.Use("/", middleware.JWTMiddleware())
	apiGroup := app.Group("/exchange/api")
	routes.RegisterRoutes(apiGroup, codeGen, appServices)
	log.Printf("Server configured and ready on port %s", port)
	return app
}
