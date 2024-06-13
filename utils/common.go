package utils

import (
	"math"
	"time"
)

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func InTimeSpan(start, end, now time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}