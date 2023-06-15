package utils

import (
	"axis/ecommerce-backend/pkg/entities"
	"math"
)

func Empty() entities.E {
	return entities.E{}
}

func RoundOff(value float64) float64 {
	var decimalFactor float64 = 100
	return math.Ceil(value*decimalFactor) / decimalFactor
}

func String(v string) *string {
	return &v
}
