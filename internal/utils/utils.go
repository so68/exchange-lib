package utils

import (
	"fmt"
	"math/big"
	"regexp"
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
		return 8
	}
	parts := strings.Split(number, ".")
	if len(parts) != 2 {
		return 8
	}
	return len(parts[1])
}

var contractRegex = regexp.MustCompile(`^([A-Z]{1,10})USDT(.*)?$`)

// FormatSymbol 格式化交易对
func FormatSymbol(symbol string, ft string) string {
	if contractRegex.MatchString(symbol) {
		matches := contractRegex.FindStringSubmatch(symbol)
		base := matches[1]
		suffix := matches[2] // 交割合约如 20251227
		return fmt.Sprintf("%s%sUSDT%s", base, ft, suffix)
	}

	return symbol
}

// FormatSymbols 格式化交易对列表
func FormatSymbols(symbols []string, ft string) []string {
	for _, symbol := range symbols {
		symbols = append(symbols, FormatSymbol(symbol, ft))
	}
	return symbols
}
