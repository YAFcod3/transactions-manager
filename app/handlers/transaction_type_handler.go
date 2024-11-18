package handlers

import (
	"transactions-manager/app/models"
	"transactions-manager/app/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionTypeHandler struct {
	Service *services.TransactionTypeService
}

func NewTransactionTypeHandler(service *services.TransactionTypeService) *TransactionTypeHandler {
	return &TransactionTypeHandler{Service: service}
}

func (h *TransactionTypeHandler) CreateTransactionType(c *fiber.Ctx) error {
	var req models.TransactionType
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request body.",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "MISSING_NAME",
			"message": "Transaction type name is required.",
		})
	}

	req.ID = primitive.NewObjectID()
	if err := h.Service.CreateTransactionType(req); err != nil {
		if err.Error() == "transaction type already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"code":    "DUPLICATE_NAME",
				"message": "A transaction type with this name already exists.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "CREATE_FAILED",
			"message": "Failed to create transaction type.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaction type created successfully.",
		"data":    req,
	})
}

func (h *TransactionTypeHandler) GetTransactionTypes(c *fiber.Ctx) error {
	types, err := h.Service.GetTransactionTypes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "FETCH_FAILED",
			"message": "Failed to fetch transaction types.",
		})
	}
	return c.JSON(fiber.Map{"data": types})
}

func (h *TransactionTypeHandler) GetTransactionTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	transactionType, err := h.Service.GetTransactionTypeByID(id)
	if err != nil {
		if err.Error() == "transaction type not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Transaction type not found.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "FETCH_FAILED",
			"message": "Failed to fetch transaction type.",
		})
	}
	return c.JSON(fiber.Map{"data": transactionType})
}

func (h *TransactionTypeHandler) UpdateTransactionType(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.TransactionType

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request body.",
		})
	}

	if req.Name == "" && req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "MISSING_FIELDS",
			"message": "At least one field must be provided.",
		})
	}

	if err := h.Service.UpdateTransactionType(id, req); err != nil {
		if err.Error() == "transaction type name already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"code":    "DUPLICATE_NAME",
				"message": "A transaction type with this name already exists.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "UPDATE_FAILED",
			"message": "Failed to update transaction type.",
		})
	}

	return c.JSON(fiber.Map{"message": "Transaction type updated successfully."})
}

func (h *TransactionTypeHandler) DeleteTransactionType(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.Service.DeleteTransactionType(id); err != nil {
		if err.Error() == "invalid ID" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_ID",
				"message": "The provided ID is invalid.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "DELETE_FAILED",
			"message": "Failed to delete transaction type.",
		})
	}

	return c.JSON(fiber.Map{"message": "Transaction type deleted successfully."})
}
