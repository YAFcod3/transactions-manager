package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transactions-manager/app/config"
	"transactions-manager/app/database"
	"transactions-manager/app/middleware"
	"transactions-manager/app/routes"
	"transactions-manager/app/services"
	"transactions-manager/app/utils"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.LoadConfig()
	database.Init()
	defer database.Close()

	utils.StartExchangeRateUpdater(database.RedisClient, 1*time.Hour)
	codeGen := &generate_transaction_code.CodeGenerator{Client: database.RedisClient}
	codeGen.LoadLastCounter()
	appServices := services.InitServices(database.MongoClient.Database(os.Getenv("MONGO_DB_NAME")), database.RedisClient, codeGen)

	app := fiber.New()
	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "*"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowCredentials: false,
		// AllowCredentials: true,
		AllowMethods: "GET, POST, PUT,PATCH, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Authorization, Accept",
	}))

	app.Use("/", middleware.JWTMiddleware())
	routes.RegisterRoutes(app, codeGen, appServices)

	log.Printf("Starting server on port %s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Listen and server on port on the /  goroutine to prevent blocking the main goroutine
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error while shutting down server: %v", err)
	}

	utils.StopExchangeRateUpdater()

	log.Println("Server shut down gracefully.")
}
