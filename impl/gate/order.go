package gate

import (
	"fmt"
	"math/big"
)

// filtersQuantity 获取交易对数量精度
func (g *gateExchange) filtersQuantity(spec *symbolSpec, price, quantity string) (string, error) {
	// 使用 big.Float 进行精确比较
	quantityFloat := new(big.Float).SetPrec(64)
	minQtyFloat := new(big.Float).SetPrec(64)
	maxQtyFloat := new(big.Float).SetPrec(64)
	priceFloat := new(big.Float).SetPrec(64)
	minPriceFloat := new(big.Float).SetPrec(64)
	maxPriceFloat := new(big.Float).SetPrec(64)

	if _, ok := quantityFloat.SetString(quantity); !ok {
		return "", fmt.Errorf("无效的数量: %s", quantity)
	}
	if _, ok := minQtyFloat.SetString(spec.MinBaseAmount); !ok {
		return "", fmt.Errorf("无效的最小数量: %s", spec.MinBaseAmount)
	}
	if _, ok := maxQtyFloat.SetString(spec.MaxBaseAmount); !ok {
		return "", fmt.Errorf("无效的最大数量: %s", spec.MaxBaseAmount)
	}
	if _, ok := priceFloat.SetString(price); !ok {
		return "", fmt.Errorf("无效的价格: %s", price)
	}
	if _, ok := minPriceFloat.SetString(spec.MinQuoteAmount); !ok {
		return "", fmt.Errorf("无效的最小价格: %s", spec.MinQuoteAmount)
	}
	if _, ok := maxPriceFloat.SetString(spec.MaxQuoteAmount); !ok {
		return "", fmt.Errorf("无效的最大价格: %s", spec.MaxQuoteAmount)
	}
	// 如果价格小于最小价格，或大于最大价格，返回错误
	if priceFloat.Cmp(minPriceFloat) < 0 || priceFloat.Cmp(maxPriceFloat) > 0 {
		return "", fmt.Errorf("价格 %s 小于最小价格 %s 或大于最大价格 %s", price, spec.MinQuoteAmount, spec.MaxQuoteAmount)
	}
	// 如果 quantity 小于 minQty，或大于 maxQty，返回错误
	if quantityFloat.Cmp(minQtyFloat) < 0 || quantityFloat.Cmp(maxQtyFloat) > 0 {
		return "", fmt.Errorf("数量 %s 小于最小数量 %s 或大于最大数量 %s", quantity, spec.MinBaseAmount, spec.MaxBaseAmount)
	}

	// 逻辑：floor(quantity * 10^precision) / 10^precision
	multiplier := new(big.Float).SetPrec(64).SetInt(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(int64(spec.AmountPrecision)), nil))
	// 乘以 10^precision
	quantityFloat = quantityFloat.Mul(quantityFloat, multiplier)
	// 向下取整
	quantityInt, _ := quantityFloat.Int(nil)
	// 除以 10^precision
	quantityFloat = new(big.Float).SetInt(quantityInt)
	quantityFloat = quantityFloat.Quo(quantityFloat, multiplier)

	return quantityFloat.Text('f', spec.AmountPrecision), nil
}
