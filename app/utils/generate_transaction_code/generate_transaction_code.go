package generate_transaction_code

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type CodeGenerator struct {
	counter  int
	Client   *redis.Client
	mu       sync.Mutex
	lastHour string
}

func (cg *CodeGenerator) LoadLastCounter() {
	ctx := context.Background()
	vals, err := cg.Client.HGetAll(ctx, "transaction_data").Result()
	if err == nil {
		cg.counter, _ = strconv.Atoi(vals["counter"])
		cg.lastHour = vals["last_hour"]
	}
}

func (cg *CodeGenerator) GenerateCode() (string, error) {
	cg.mu.Lock()
	defer cg.mu.Unlock()

	location, err := time.LoadLocation("America/Havana")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	now := time.Now().In(location)
	currentHour := now.Format("06010215")

	ctx := context.Background()
	if currentHour != cg.lastHour {
		cg.counter = 0
		cg.lastHour = currentHour

		err = cg.Client.HSet(ctx, "transaction_data", "counter", cg.counter, "last_hour", cg.lastHour).Err()
		if err != nil {
			log.Println("Error al actualizar datos en Redis:", err)
		}
	} else if cg.counter == 100000000 {
		return "", fmt.Errorf("límite de códigos alcanzado para la hora actual")
	}

	code := fmt.Sprintf("T%s%08d", currentHour, cg.counter)
	cg.counter++

	err = cg.Client.HSet(ctx, "transaction_data", "counter", cg.counter).Err()
	if err != nil {
		log.Println("Error al actualizar counter:", err)
	}

	return code, nil
}
