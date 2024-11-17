package routes

import (
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, codeGen *generate_transaction_code.CodeGenerator) {
	// SetupConversionRoutes(app, mongoClient, redisClient, codeGen)
	SetupConversionRoutes(app, codeGen)

}
