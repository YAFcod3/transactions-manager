package routes

import (
	"transactions-manager/app/handlers"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
)

func SetupStatisticsRoutes(apiGroup fiber.Router, statisticsService *services.StatisticsService) {
	statisticsHandler := handlers.NewStatisticsHandler(statisticsService)
	apiGroup.Get("/statistics", statisticsHandler.GetStatisticsHandler)

}
