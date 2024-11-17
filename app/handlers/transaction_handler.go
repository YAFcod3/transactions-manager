package handlers

import (
	"transactions-manager/app/models"
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	Service *services.ConversionService
}

func NewTransactionHandler(codeGen *generate_transaction_code.CodeGenerator, supportedCurrenciesService *services.SupportedCurrenciesService) *TransactionHandler {
	return &TransactionHandler{
		Service: services.NewConversionService(codeGen, supportedCurrenciesService),
		// Service: services.NewConversionService(mongoClient, redisClient, codeGen),

	}
}

func (h *TransactionHandler) HandleTransaction(c *fiber.Ctx) error {
	var req models.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// userID := c.Locals("userId").(string)
	userID := "1234"

	// Llama al servicio para procesar la transacción
	response, err := h.Service.ProcessTransaction(req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}
