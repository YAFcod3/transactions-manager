package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transactions-manager/app/database"
	"transactions-manager/app/server"
	"transactions-manager/app/services"
	"transactions-manager/app/utils"
	"transactions-manager/app/utils/generate_transaction_code"
)

func main() {
	mongoDBName := utils.GetEnvOrPanic("MONGO_DB_NAME")
	port := utils.GetEnvOrPanic("PORT")
	database.Init()
	defer database.Close()

	utils.StartExchangeRateUpdater(database.RedisClient, 1*time.Hour)
	codeGen := &generate_transaction_code.CodeGenerator{Client: database.RedisClient}
	codeGen.LoadLastCounter()
	appServices := services.InitServices(database.MongoClient.Database(mongoDBName), database.RedisClient, codeGen)
	app := server.CreateServer(appServices, codeGen, port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Listen and server on port on the /  goroutine to prevent blocking the main goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
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
