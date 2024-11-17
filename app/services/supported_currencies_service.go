package services

import "errors"

type SupportedCurrenciesService struct {
	supportedCurrencies map[string]bool
}

func NewSupportedCurrenciesService() *SupportedCurrenciesService {
	return &SupportedCurrenciesService{
		supportedCurrencies: map[string]bool{
			"USD": true,
			"EUR": true,
			"GBP": true,
			"JPY": true,
			"CAD": true,
		},
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
