package models

type StatisticsResponse struct {
	TotalTransactions              int                `json:"totalTransactions"`
	TransactionsByType             map[string]int     `json:"transactionsByType"`
	TotalAmountConvertedByCurrency map[string]float64 `json:"totalAmountConvertedByCurrency"`
	AverageAmountByTransactionType map[string]float64 `json:"averageAmountByTransactionType"`
}
