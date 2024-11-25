package utils

import "math"

func RoundOperations(value float64) float64 {
	return math.Round(value*100) / 100
}
