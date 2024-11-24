package middleware

import (
	"fmt"
	"net/http"
	"time"
	"transactions-manager/app/database"
	"transactions-manager/app/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

const transactionTTL = 20 * time.Second

func VerifyTransactionDuplicated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var transaction models.TransactionRequest
		if err := c.BodyParser(&transaction); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Invalid transaction data.",
			})
		}

		userID := c.Locals("userId").(string)

		uniqueKey := fmt.Sprintf("%s:%s:%s:%.2f:%s", userID, transaction.FromCurrency, transaction.ToCurrency, transaction.Amount, transaction.TransactionType)

		ctx := context.Background()
		exists, err := database.RedisClient.Exists(ctx, uniqueKey).Result()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    "REDIS_ERROR",
				"message": "Error checking duplicate transaction in Redis.",
			})
		}

		if exists > 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    "DUPLICATE_TRANSACTION",
				"message": "A similar transaction was already processed within the last 20 seconds. Please try again later.",
			})
		}

		err = database.RedisClient.Set(ctx, uniqueKey, "1", transactionTTL).Err()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    "REDIS_ERROR",
				"message": "Error setting transaction limit in Redis.",
			})
		}

		return c.Next()
	}
}
