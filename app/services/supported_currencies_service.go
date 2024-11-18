package services

import (
	"errors"
	"os"
	"strings"
)

type SupportedCurrenciesService struct {
	supportedCurrencies map[string]bool
}

func NewSupportedCurrenciesService() *SupportedCurrenciesService {
	currencies := os.Getenv("SUPPORTED_CURRENCIES")
	if currencies == "" {
		currencies = "USD,EUR,GBP,JPY,CAD,AUD"
	}

	supportedCurrencies := make(map[string]bool)
	for _, currency := range strings.Split(currencies, ",") {
		supportedCurrencies[strings.TrimSpace(currency)] = true
	}

	return &SupportedCurrenciesService{
		supportedCurrencies: supportedCurrencies,
	}
}

func (s *SupportedCurrenciesService) GetSupportedCurrencies() []string {
	currencies := make([]string, 0, len(s.supportedCurrencies))
	for currency := range s.supportedCurrencies {
		currencies = append(currencies, currency)
	}
	return currencies
}

func (s *SupportedCurrenciesService) IsCurrencySupported(currency string) error {
	if _, exists := s.supportedCurrencies[currency]; !exists {
		return errors.New("currency not supported")
	}
	return nil
}
