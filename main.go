package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transactions-manager/app/config"
	"transactions-manager/app/database"
	"transactions-manager/app/routes"
	"transactions-manager/app/utils"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.Init()
	defer database.Close()

	// Start exchange rate updater
	utils.StartExchangeRateUpdater(database.RedisClient, 1*time.Hour)

	// Initialize transaction code generator
	codeGen := &generate_transaction_code.CodeGenerator{Client: database.RedisClient}
	codeGen.LoadLastCounter()

	// Create Fiber app
	app := fiber.New()

	// Register routes
	routes.RegisterRoutes(app, codeGen)

	// Log server start message
	log.Printf("Starting server on port %s", cfg.Port)

	// Channel to capture system signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal
	<-quit
	log.Println("Shutting down server...")

	// Shut down Fiber server
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error while shutting down server: %v", err)
	}

	// Stop exchange rate updater
	utils.StopExchangeRateUpdater()

	log.Println("Server shut down gracefully.")
}
