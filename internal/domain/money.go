package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
)

var ErrAmountDecimalPlaces = errors.New("Amount cannot have more than 2 decimal places")

// Avoid floating-point precision issues
type Money int64

func (m *Money) UnmarshalJSON(data []byte) error {
	strValue := string(data)
	parts := strings.Split(strValue, ".")

	// Ensure at most 2 decimal places
	if len(parts) == 2 && len(parts[1]) > 2 {
		return ErrAmountDecimalPlaces
	}

	var floatAmount float64
	if err := json.Unmarshal(data, &floatAmount); err != nil {
		return err
	}

	*m = Money(math.Round(floatAmount * 100))
	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	floatAmount := float64(m) / 100.0
	formattedAmount := fmt.Sprintf("%.2f", floatAmount)
	return []byte(formattedAmount), nil
}
