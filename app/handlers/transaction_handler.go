package handlers

import (
	"fmt"
	"strings"
	"transactions-manager/app/models"
	"transactions-manager/app/services"
	"transactions-manager/app/utils/generate_transaction_code"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	Service *services.ConversionService
}

func NewTransactionHandler(codeGen *generate_transaction_code.CodeGenerator, transactionTypeService *services.TransactionTypeService) *TransactionHandler {
	return &TransactionHandler{
		Service: services.NewConversionService(codeGen, transactionTypeService),
	}
}

func (h *TransactionHandler) HandleTransaction(c *fiber.Ctx) error {
	var req models.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid JSON request payload",
		})
	}

	missingFields := []string{}
	if req.Amount <= 0 {
		missingFields = append(missingFields, "amount")
	}
	if req.FromCurrency == "" {
		missingFields = append(missingFields, "fromCurrency")
	}
	if req.ToCurrency == "" {
		missingFields = append(missingFields, "toCurrency")
	}
	if req.TransactionType == "" {
		missingFields = append(missingFields, "transactionType")
	}
	if len(missingFields) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "MISSING_REQUIRED_FIELDS",
			"message": fmt.Sprintf("The following required fields are missing: %s. Please provide valid values for all required fields", strings.Join(missingFields, ", ")),
		})
	}

	userID := c.Locals("userId").(string)

	response, err := h.Service.ProcessTransaction(req, userID)
	if err != nil {
		parts := strings.SplitN(err.Error(), ": ", 2)
		code := "UNKNOWN_ERROR"
		message := err.Error()
		if len(parts) == 2 {
			code = parts[0]
			message = parts[1]
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    code,
			"message": message,
		})
	}

	return c.JSON(response)
}
