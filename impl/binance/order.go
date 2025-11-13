package binance

import (
	"fmt"
	"math/big"

	"github.com/so68/exchange-lib/internal/utils"
)

// filtersQuantity 获取交易对数量精度
func (b *binanceExchange) filtersQuantity(spec *symbolSpec, price, quantity string) (string, error) {
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
	if _, ok := minQtyFloat.SetString(spec.MinQty); !ok {
		return "", fmt.Errorf("无效的最小数量: %s", spec.MinQty)
	}
	if _, ok := maxQtyFloat.SetString(spec.MaxQty); !ok {
		return "", fmt.Errorf("无效的最大数量: %s", spec.MaxQty)
	}
	if _, ok := priceFloat.SetString(price); !ok {
		return "", fmt.Errorf("无效的价格: %s", price)
	}
	if _, ok := minPriceFloat.SetString(spec.MinPrice); !ok {
		return "", fmt.Errorf("无效的最小价格: %s", spec.MinPrice)
	}
	if _, ok := maxPriceFloat.SetString(spec.MaxPrice); !ok {
		return "", fmt.Errorf("无效的最大价格: %s", spec.MaxPrice)
	}
	// 如果价格小于最小价格，或大于最大价格，返回错误
	if priceFloat.Cmp(minPriceFloat) < 0 || priceFloat.Cmp(maxPriceFloat) > 0 {
		return "", fmt.Errorf("价格 %s 小于最小价格 %s 或大于最大价格 %s", price, spec.MinPrice, spec.MaxPrice)
	}
	// 如果 quantity 小于 minQty，或大于 maxQty，返回错误
	if quantityFloat.Cmp(minQtyFloat) < 0 || quantityFloat.Cmp(maxQtyFloat) > 0 {
		return "", fmt.Errorf("数量 %s 小于最小数量 %s 或大于最大数量 %s", quantity, spec.MinQty, spec.MaxQty)
	}

	// 按照 stepSize 的倍数向下取整 quantity
	// 逻辑：floor(quantity / stepSize) * stepSize
	stepSizeFloat := new(big.Float).SetPrec(64)
	if _, ok := stepSizeFloat.SetString(spec.StepSize); !ok {
		return "", fmt.Errorf("无效的步长: %s", spec.StepSize)
	}

	// quantity / stepSize
	ratio := new(big.Float).Quo(quantityFloat, stepSizeFloat)

	// 向下取整
	ratioInt, _ := ratio.Int(nil)

	// 乘以 stepSize: floor(quantity / stepSize) * stepSize
	quantityFloat = new(big.Float).SetInt(ratioInt)
	quantityFloat = quantityFloat.Mul(quantityFloat, stepSizeFloat)

	// 再次检查处理后的 quantity 是否大于等于 minQty
	if quantityFloat.Cmp(minQtyFloat) < 0 {
		precision := utils.GetNumberPrecision(spec.StepSize)
		return "", fmt.Errorf("处理后的数量 %s 小于最小数量 %s", quantityFloat.Text('f', precision), spec.MinQty)
	}

	precision := utils.GetNumberPrecision(spec.StepSize)
	return quantityFloat.Text('f', precision), nil
}
