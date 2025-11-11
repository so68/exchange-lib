package utils

import (
	"math/big"
	"strings"
)

// AmountWithPriceToQuantity 金额与价格转换为数量
func AmountWithPriceToQuantity(amount float64, price string, prec int) string {
	quantity := big.NewFloat(0)
	if _, ok := quantity.SetString(price); !ok {
		return "0"
	}
	quantity = quantity.Quo(big.NewFloat(1), quantity).Mul(quantity, big.NewFloat(amount))
	return quantity.Text('f', prec)
}

// GetNumberPrecision 从 number 字符串中提取小数位数
func GetNumberPrecision(number string) int {
	if !strings.Contains(number, ".") {
		return 0
	}
	parts := strings.Split(number, ".")
	if len(parts) != 2 {
		return 0
	}
	return len(parts[1])
}
