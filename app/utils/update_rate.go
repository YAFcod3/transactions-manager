package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type ExchangeRate struct {
	Base      string             `json:"base"`
	Timestamp int64              `json:"timestamp"`
	Rates     map[string]float64 `json:"rates"`
}

var stopChan = make(chan struct{})

func StartExchangeRateUpdater(client *redis.Client, interval time.Duration) {
	ctx := context.Background()
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	baseCurrency := os.Getenv("BASE_CURRENCY")
	if baseCurrency == "" {
		log.Fatal("Base currency not set. Please set the BASE_CURRENCY environment variable.")
		return
	}
	baseURL := os.Getenv("URL_API_EXTERNAL_GET_RATE")
	if baseURL == "" {
		log.Fatal("Base URL not set. Please set the URL_API_EXTERNAL_GET_RATE environment variable.")
		return
	}

	fetchExchangeRates := func() {
		u, err := url.Parse(baseURL)
		if err != nil {
			fmt.Println("Invalid base URL:", err)
			return
		}
		u.Path = "/exchange-rate/api/latest"
		query := u.Query()
		query.Set("base", baseCurrency)
		u.RawQuery = query.Encode()

		resp, err := httpClient.Get(u.String())
		if err != nil {
			fmt.Println("Error fetching exchange rates:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Error: received status code %d\nResponse body: %s\n", resp.StatusCode, string(body))
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var exchangeRate ExchangeRate
		err = json.Unmarshal(body, &exchangeRate)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		key := "exchange_rates"
		pipe := client.TxPipeline()

		for fromCurrency, fromRate := range exchangeRate.Rates {
			if fromCurrency != baseCurrency {
				baseKey := fmt.Sprintf("%s:%s", baseCurrency, fromCurrency)
				pipe.Set(ctx, baseKey, fromRate, 65*time.Minute)

				inverseBaseKey := fmt.Sprintf("%s:%s", fromCurrency, baseCurrency)
				pipe.Set(ctx, inverseBaseKey, 1/fromRate, 65*time.Minute)
			}

			for toCurrency, toRate := range exchangeRate.Rates {
				if fromCurrency != toCurrency {
					relativeKey := fmt.Sprintf("%s:%s", fromCurrency, toCurrency)
					relativeRate := toRate / fromRate
					pipe.Set(ctx, relativeKey, relativeRate, 65*time.Minute)
				}
			}
		}

		pipe.HSet(ctx, key, "base", exchangeRate.Base)
		pipe.HSet(ctx, key, "timestamp", exchangeRate.Timestamp)

		_, err = pipe.Exec(ctx)
		if err != nil {
			fmt.Println("Error storing exchange rates in Redis:", err)
			return
		}

		fmt.Println("Exchange rates updated successfully!")
	}

	// Fetch initial rates
	fetchExchangeRates()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				fetchExchangeRates()
			case <-stopChan:
				fmt.Println("Exchange rate updater stopped.")
				return
			}
		}
	}()
}

func StopExchangeRateUpdater() {
	close(stopChan)
}
