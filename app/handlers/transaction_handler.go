package handlers

import (
	"errors"
	"transactions-manager/app/models"
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	Service *services.ConversionService
}

func NewTransactionHandler(codeGen *generate_transaction_code.CodeGenerator, supportedCurrenciesService *services.SupportedCurrenciesService, transactionTypeService *services.TransactionTypeService) *TransactionHandler {
	return &TransactionHandler{
		Service: services.NewConversionService(codeGen, supportedCurrenciesService, transactionTypeService),
	}
}

func (h *TransactionHandler) HandleTransaction(c *fiber.Ctx) error {
	var req models.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Amount <= 0 {
		return errors.New("'amount' must be greater than zero")
	}
	if req.FromCurrency == "" {
		return errors.New("'fromCurrency' is required")
	}
	if req.ToCurrency == "" {
		return errors.New("'toCurrency' is required")
	}
	if req.TransactionType == "" {
		return errors.New("'transactionType' is required")
	}

	userID := c.Locals("userId").(string)
	if userID == "" {
		return errors.New("userID is required")
	}

	response, err := h.Service.ProcessTransaction(req, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(response)
}
