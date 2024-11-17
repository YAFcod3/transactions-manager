package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupStatisticsRoutes(app *fiber.App, statisticsService *services.StatisticsService) {
	statisticsGroup := app.Group("/exchange/api/statistics/")
	statisticsHandler := handlers.NewStatisticsHandler(statisticsService)
	statisticsGroup.Get("/", statisticsHandler.GetStatisticsHandler)

}
