package binance

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/so68/exchange-lib/exchange"
)

// CreateSpotOrder 创建现货订单
func (b *binanceExchange) CreateSpotOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	spec, err := b.getSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, err
	}

	quantity, err = b.filtersQuantity(spec, limitPrice, quantity)
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
		return nil, err
	}

	return &exchange.Order{
		OrderID:       strconv.FormatInt(orderResp.OrderID, 10),
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

// GetSpotOrder 获取现货订单
func (b *binanceExchange) GetSpotOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}
	resp, err := b.client.NewGetOrderService().
		Symbol(symbol).
		OrderID(orderIDInt).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance spot get order failed: %w", err)
	}

	// 获取订单的成交记录（包含手续费信息）
	trades, err := b.client.NewListTradesService().
		Symbol(symbol).
		OrderId(orderIDInt).
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
		OrderID:       strconv.FormatInt(resp.OrderID, 10),
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

// CancelSpotOrder 撤销现货订单
func (b *binanceExchange) CancelSpotOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}
	resp, err := b.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(orderIDInt).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance spot cancel order failed: %w", err)
	}

	return &exchange.Order{
		OrderID:       orderID,
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

	if spec == nil {
		return nil, fmt.Errorf("交易对规格不存在: %s", symbol)
	}

	return spec, nil
}
