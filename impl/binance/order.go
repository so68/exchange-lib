package binance

import (
	"context"
	"fmt"
	"math/big"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/so68/exchange-lib/exchange"
	"github.com/so68/exchange-lib/internal/utils"
)

// SpotCreateOrder 创建现货订单
func (b *binanceExchange) SpotCreateOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	spec, err := b.getSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("获取交易规则失败: %w", err)
	}

	quantity, err = b.filtersQuantity(spec, quantity)
	if err != nil {
		return nil, fmt.Errorf("验证交易规则失败: %w", err)
	}

	// 创建订单服务
	service := b.client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideType(string(side))).
		Quantity(quantity)

	// 市价单
	if limitPrice == "" || limitPrice == "0" {
		service.Type(binance.OrderTypeMarket)
	} else {
		service.Type(binance.OrderTypeLimit).Price(limitPrice).TimeInForce(binance.TimeInForceTypeGTC)
	}

	// 执行订单
	orderResp, err := service.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建现货订单失败: %w", err)
	}

	return &exchange.Order{
		OrderID:       orderResp.OrderID,
		Symbol:        orderResp.Symbol,
		Side:          exchange.OrderSide(orderResp.Side),
		Type:          exchange.OrderType(orderResp.Type),
		Status:        exchange.OrderStatus(string(orderResp.Status)),
		Price:         orderResp.Price,
		Quantity:      orderResp.OrigQuantity,
		ExecutedQty:   orderResp.ExecutedQuantity,
		QuoteQuantity: orderResp.CummulativeQuoteQuantity,
		TimeInForce:   exchange.OrderTimeInForce(orderResp.TimeInForce),
		CreateTime:    orderResp.TransactTime,
		UpdateTime:    orderResp.TransactTime,
	}, nil
}

// FuturesCreateOrder 合约下单
func (b *binanceExchange) FuturesCreateOrder(ctx context.Context, symbol string, side exchange.OrderSide, orderType exchange.OrderType, quantity, limitPrice string) (*exchange.Order, error) {
	service := b.futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Side(futures.SideType(side)).
		Type(futures.OrderType(orderType)).
		Quantity(quantity)

	if orderType == "LIMIT" {
		service = service.Price(limitPrice).TimeInForce(futures.TimeInForceType(exchange.OrderTimeInForceGTC))
	}

	resp, err := service.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance futures create order: %w", err)
	}

	return &exchange.Order{
		OrderID:       resp.OrderID,
		Symbol:        resp.Symbol,
		Side:          exchange.OrderSide(resp.Side),
		Type:          exchange.OrderType(resp.Type),
		Status:        exchange.OrderStatus(string(resp.Status)),
		Price:         resp.Price,
		Quantity:      resp.OrigQuantity,
		ExecutedQty:   resp.ExecutedQuantity,
		QuoteQuantity: resp.CumQuote,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.UpdateTime,
		UpdateTime:    resp.UpdateTime,
	}, nil
}

// SpotGetOrder 获取现货订单
func (b *binanceExchange) SpotGetOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	resp, err := b.client.NewGetOrderService().
		Symbol(symbol).
		OrderID(orderID).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance spot get order failed: %w", err)
	}

	// 获取订单的成交记录（包含手续费信息）
	trades, err := b.client.NewListTradesService().
		Symbol(symbol).
		OrderId(orderID).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance spot get order trades failed: %w", err)
	}

	// 获取交易对规格，用于计算实际数量
	spec, err := b.getSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("获取交易规则失败: %w", err)
	}

	var actualQty string
	side := exchange.OrderSide(resp.Side)
	totalCommission := new(big.Float).SetPrec(64)

	// 将成交记录转换为 Fill
	for _, trade := range trades {
		// 累计手续费
		commission := new(big.Float).SetPrec(64)
		if _, ok := commission.SetString(trade.Commission); !ok {
			return nil, fmt.Errorf("无效的手续费: %s", trade.Commission)
		}

		if side == exchange.OrderSideBuy && trade.CommissionAsset == spec.BaseAsset {
			totalCommission = totalCommission.Add(totalCommission, commission)
		}

		if side == exchange.OrderSideSell && trade.CommissionAsset == spec.QuoteAsset {
			totalCommission = totalCommission.Add(totalCommission, commission)
		}
	}

	if side == exchange.OrderSideBuy {
		bigActualQty := new(big.Float).SetPrec(64)
		if _, ok := bigActualQty.SetString(resp.ExecutedQuantity); !ok {
			return nil, fmt.Errorf("无效的数量: %s", resp.ExecutedQuantity)
		}
		actualQty = bigActualQty.Sub(bigActualQty, totalCommission).Text('f', spec.BasePrecision)
	} else {
		bigActualQty := new(big.Float).SetPrec(64)
		if _, ok := bigActualQty.SetString(resp.CummulativeQuoteQuantity); !ok {
			return nil, fmt.Errorf("无效的成交金额: %s", resp.CummulativeQuoteQuantity)
		}
		actualQty = bigActualQty.Sub(bigActualQty, totalCommission).Text('f', spec.QuotePrecision)
	}

	return &exchange.Order{
		OrderID:       resp.OrderID,
		Symbol:        resp.Symbol,
		Side:          exchange.OrderSide(resp.Side),
		Type:          exchange.OrderType(resp.Type),
		Status:        exchange.OrderStatus(string(resp.Status)),
		Price:         resp.Price,
		Quantity:      resp.OrigQuantity,
		ExecutedQty:   resp.ExecutedQuantity,
		ActualQty:     actualQty,
		QuoteQuantity: resp.CummulativeQuoteQuantity,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.UpdateTime,
		UpdateTime:    resp.UpdateTime,
	}, nil
}

// FuturesGetOrder 获取合约订单
func (b *binanceExchange) FuturesGetOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	resp, err := b.futuresClient.NewGetOrderService().Symbol(symbol).OrderID(orderID).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance futures get order: %w", err)
	}

	return &exchange.Order{
		OrderID:       resp.OrderID,
		Symbol:        resp.Symbol,
		Side:          exchange.OrderSide(resp.Side),
		Type:          exchange.OrderType(resp.Type),
		Status:        exchange.OrderStatus(string(resp.Status)),
		Price:         resp.Price,
		Quantity:      resp.OrigQuantity,
		ExecutedQty:   resp.ExecutedQuantity,
		QuoteQuantity: resp.CumQuote,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.UpdateTime,
		UpdateTime:    resp.UpdateTime,
	}, nil
}

// SpotCancelOrder 撤销现货订单
func (b *binanceExchange) SpotCancelOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	resp, err := b.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderID).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance spot cancel order failed: %w", err)
	}

	return &exchange.Order{
		OrderID:       resp.OrderID,
		Symbol:        resp.Symbol,
		Side:          exchange.OrderSide(resp.Side),
		Type:          exchange.OrderType(resp.Type),
		Status:        exchange.OrderStatus(string(resp.Status)),
		Price:         resp.Price,
		Quantity:      resp.OrigQuantity,
		ExecutedQty:   resp.ExecutedQuantity,
		QuoteQuantity: resp.CummulativeQuoteQuantity,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.TransactTime,
		UpdateTime:    resp.TransactTime,
	}, nil
}

// FuturesCancelOrder 撤销合约订单
func (b *binanceExchange) FuturesCancelOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	resp, err := b.futuresClient.NewCancelOrderService().Symbol(symbol).OrderID(orderID).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance futures cancel order: %w", err)
	}
	return &exchange.Order{
		OrderID:       resp.OrderID,
		Symbol:        resp.Symbol,
		Side:          exchange.OrderSide(resp.Side),
		Type:          exchange.OrderType(resp.Type),
		Status:        exchange.OrderStatus(string(resp.Status)),
		Price:         resp.Price,
		Quantity:      resp.OrigQuantity,
		ExecutedQty:   resp.ExecutedQuantity,
		QuoteQuantity: resp.CumQuote,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.UpdateTime,
		UpdateTime:    resp.UpdateTime,
	}, nil
}

// filtersQuantity 获取交易对数量精度
func (b *binanceExchange) filtersQuantity(spec *symbolSpec, quantity string) (string, error) {
	// 使用 big.Float 进行精确比较
	quantityFloat := new(big.Float).SetPrec(64)
	minQtyFloat := new(big.Float).SetPrec(64)

	if _, ok := quantityFloat.SetString(quantity); !ok {
		return "", fmt.Errorf("无效的数量: %s", quantity)
	}
	if _, ok := minQtyFloat.SetString(spec.MinQty); !ok {
		return "", fmt.Errorf("无效的最小数量: %s", spec.MinQty)
	}
	// 如果 quantity 小于 minQty，返回错误
	if quantityFloat.Cmp(minQtyFloat) < 0 {
		return "", fmt.Errorf("数量 %s 小于最小数量 %s", quantity, spec.MinQty)
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

// getSpotSymbolSpec 获取现货交易对规格
func (b *binanceExchange) getSpotSymbolSpec(ctx context.Context, symbol string) (*symbolSpec, error) {
	var spec *symbolSpec

	// 从缓存中获取交易对规格
	spec, _ = binanceSpotSpec.GetSymbolSpec(symbol)

	// 如果缓存中没有，则获取最新交易对规格
	if spec == nil {
		info, err := b.client.NewExchangeInfoService().Symbol(symbol).Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取交易规则失败: %w", err)
		}

		symbolInfo := info.Symbols[0]
		spec = &symbolSpec{
			Symbol:         symbol,
			BaseAsset:      symbolInfo.BaseAsset,
			QuoteAsset:     symbolInfo.QuoteAsset,
			BasePrecision:  symbolInfo.BaseAssetPrecision,
			QuotePrecision: symbolInfo.QuoteAssetPrecision,
			Status:         symbolInfo.Status,
		}
		for _, f := range symbolInfo.Filters {
			if f["filterType"].(string) == "PRICE_FILTER" {
				spec.MinPrice = f["minPrice"].(string)
				spec.MaxPrice = f["maxPrice"].(string)
				spec.TickSize = f["tickSize"].(string)
			}
			if f["filterType"].(string) == "LOT_SIZE" {
				spec.MinQty = f["minQty"].(string)
				spec.MaxQty = f["maxQty"].(string)
				spec.StepSize = f["stepSize"].(string)
			}
			if f["filterType"].(string) == "MIN_NOTIONAL" {
				spec.MinNotional = f["minNotional"].(string)
			}
		}

		binanceSpotSpec.SetSymbolSpec(symbol, spec)
	}

	return spec, nil
}

// getFuturesSymbolSpec 获取合约交易对规格
func (b *binanceExchange) getFuturesSymbolSpec(ctx context.Context, symbol string) (*symbolSpec, error) {
	var spec *symbolSpec

	// 从缓存中获取交易对规格
	spec, _ = binanceFuturesSpec.GetSymbolSpec(symbol)

	// 如果缓存中没有，则获取最新交易对规格
	if spec == nil {
		info, err := b.futuresClient.NewExchangeInfoService().Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取交易规则失败: %w", err)
		}

		for _, s := range info.Symbols {
			if s.Status != "TRADING" {
				continue
			}

			specTmp := &symbolSpec{
				Symbol:         s.Symbol,
				BaseAsset:      s.BaseAsset,
				QuoteAsset:     s.QuoteAsset,
				BasePrecision:  s.BaseAssetPrecision,
				QuotePrecision: s.QuantityPrecision,
				Status:         s.Status,
			}

			for _, f := range s.Filters {
				if f["filterType"].(string) == "PRICE_FILTER" {
					specTmp.MinPrice = f["minPrice"].(string)
					specTmp.MaxPrice = f["maxPrice"].(string)
					specTmp.TickSize = f["tickSize"].(string)
				}
				if f["filterType"].(string) == "LOT_SIZE" {
					specTmp.MinQty = f["minQty"].(string)
					specTmp.MaxQty = f["maxQty"].(string)
					specTmp.StepSize = f["stepSize"].(string)
				}
				if f["filterType"].(string) == "MIN_NOTIONAL" {
					specTmp.MinNotional = f["minNotional"].(string)
				}
			}

			// 找到对应的交易对规格
			if s.Symbol == symbol {
				spec = specTmp
			}

			binanceFuturesSpec.SetSymbolSpec(s.Symbol, specTmp)
		}
	}

	return spec, nil
}
