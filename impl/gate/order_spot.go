package gate

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/gateio/gateapi-go/v6"
	"github.com/so68/exchange-lib/exchange"
)

// CreateSpotOrder 创建现货订单
func (g *gateExchange) CreateSpotOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	spec, err := g.GetSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// 验证交易规则
	quantity, err = g.filtersQuantity(spec, limitPrice, quantity)
	if err != nil {
		return nil, fmt.Errorf("验证交易规则失败: %w", err)
	}

	orderParams := gateapi.Order{
		CurrencyPair: symbol,
		Side:         strings.ToLower(string(side)), // buy 或 sell
		Amount:       quantity,                      // 数量（基础币种）
		Price:        limitPrice,                    // 价格（报价币种）
		Type:         "limit",                       // limit（限价）或 market（市价）
	}

	createdOrder, _, err := g.client.SpotApi.CreateOrder(ctx, orderParams, nil)
	if err != nil {
		return nil, fmt.Errorf("下单失败: %w", err)
	}

	// 计算手续费
	filledAmount := new(big.Float).SetPrec(64)
	if _, ok := filledAmount.SetString(createdOrder.FilledAmount); !ok {
		return nil, fmt.Errorf("无效的已成交数量: %s", createdOrder.FilledAmount)
	}

	feeAmount := new(big.Float).SetPrec(64)
	if _, ok := feeAmount.SetString(createdOrder.Fee); !ok {
		return nil, fmt.Errorf("无效的手续费: %s", createdOrder.Fee)
	}
	actualQty := filledAmount.Sub(filledAmount, feeAmount)

	return &exchange.Order{
		OrderID:       createdOrder.Id,
		Symbol:        createdOrder.CurrencyPair,
		Side:          exchange.OrderSide(strings.ToUpper(createdOrder.Side)),
		Type:          exchange.OrderType(strings.ToUpper(createdOrder.Type)),
		Status:        exchange.OrderStatusNew, // 默认新订单
		Price:         createdOrder.Price,
		Quantity:      createdOrder.Amount,
		ExecutedQty:   createdOrder.FilledAmount,
		ActualQty:     actualQty.Text('f', spec.AmountPrecision),
		QuoteQuantity: createdOrder.FilledTotal,
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(createdOrder.TimeInForce)),
		CreateTime:    createdOrder.CreateTimeMs,
		UpdateTime:    createdOrder.UpdateTimeMs,
	}, nil
}

// GetSpotOrder 获取现货订单
func (g *gateExchange) GetSpotOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	// 获取交易对规格
	spec, err := g.GetSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, err
	}
	singleOrder, _, err := g.client.SpotApi.GetOrder(ctx, orderID, symbol, nil)
	if err != nil {
		return nil, fmt.Errorf("获取单个订单失败: %w", err)
	}

	// 计算手续费
	filledAmount := new(big.Float).SetPrec(64)
	if _, ok := filledAmount.SetString(singleOrder.FilledAmount); !ok {
		return nil, fmt.Errorf("无效的已成交数量: %s", singleOrder.FilledAmount)
	}

	feeAmount := new(big.Float).SetPrec(64)
	if _, ok := feeAmount.SetString(singleOrder.Fee); !ok {
		return nil, fmt.Errorf("无效的手续费: %s", singleOrder.Fee)
	}
	actualQty := filledAmount.Sub(filledAmount, feeAmount)

	status := exchange.OrderStatusNew
	switch singleOrder.FinishAs {
	case "filled":
		status = exchange.OrderStatusFilled
	case "cancelled":
		status = exchange.OrderStatusCanceled
	case "small":
		status = exchange.OrderStatusRejected
	case "depth_not_enough":
		status = exchange.OrderStatusRejected
	case "trader_not_enough":
		status = exchange.OrderStatusRejected
	}

	return &exchange.Order{
		OrderID:       singleOrder.Id,
		Symbol:        singleOrder.CurrencyPair,
		Side:          exchange.OrderSide(strings.ToUpper(singleOrder.Side)),
		Type:          exchange.OrderType(strings.ToUpper(singleOrder.Type)),
		Status:        status,
		Price:         singleOrder.Price,
		Quantity:      singleOrder.Amount,
		ExecutedQty:   singleOrder.FilledAmount,
		ActualQty:     actualQty.Text('f', spec.AmountPrecision),
		QuoteQuantity: singleOrder.FilledTotal,
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(singleOrder.TimeInForce)),
		CreateTime:    singleOrder.CreateTimeMs,
		UpdateTime:    singleOrder.UpdateTimeMs,
	}, nil
}

// CancelSpotOrder 撤销现货订单
func (g *gateExchange) CancelSpotOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	// 获取交易对规格
	spec, err := g.GetSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, err
	}

	canceledOrder, _, err := g.client.SpotApi.CancelOrder(ctx, orderID, symbol, nil)
	if err != nil {
		return nil, fmt.Errorf("取消单个订单失败: %w", err)
	}
	// 计算手续费
	filledAmount := new(big.Float).SetPrec(64)
	if _, ok := filledAmount.SetString(canceledOrder.FilledAmount); !ok {
		return nil, fmt.Errorf("无效的已成交数量: %s", canceledOrder.FilledAmount)
	}

	feeAmount := new(big.Float).SetPrec(64)
	if _, ok := feeAmount.SetString(canceledOrder.Fee); !ok {
		return nil, fmt.Errorf("无效的手续费: %s", canceledOrder.Fee)
	}
	actualQty := filledAmount.Sub(filledAmount, feeAmount)

	status := exchange.OrderStatusNew
	switch canceledOrder.FinishAs {
	case "filled":
		status = exchange.OrderStatusFilled
	case "cancelled":
		status = exchange.OrderStatusCanceled
	case "small":
		status = exchange.OrderStatusRejected
	case "depth_not_enough":
		status = exchange.OrderStatusRejected
	case "trader_not_enough":
		status = exchange.OrderStatusRejected
	}

	return &exchange.Order{
		OrderID:       canceledOrder.Id,
		Symbol:        canceledOrder.CurrencyPair,
		Side:          exchange.OrderSide(strings.ToUpper(canceledOrder.Side)),
		Type:          exchange.OrderType(strings.ToUpper(canceledOrder.Type)),
		Status:        status,
		Price:         canceledOrder.Price,
		Quantity:      canceledOrder.Amount,
		ExecutedQty:   canceledOrder.FilledAmount,
		ActualQty:     actualQty.Text('f', spec.AmountPrecision),
		QuoteQuantity: canceledOrder.FilledTotal,
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(canceledOrder.TimeInForce)),
		CreateTime:    canceledOrder.CreateTimeMs,
		UpdateTime:    canceledOrder.UpdateTimeMs,
	}, nil
}

// GetSpotSymbolSpec 获取现货交易对规格
func (g *gateExchange) GetSpotSymbolSpec(ctx context.Context, symbol string) (*symbolSpec, error) {
	var spec *symbolSpec

	// 从缓存中获取交易对规格
	spec, _ = gateSpotSpec.GetSymbolSpec(symbol)

	if spec == nil {
		// 获取所有现货交易对规则
		pairs, _, err := g.client.SpotApi.ListCurrencyPairs(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取现货交易对规则失败: %w", err)
		}

		for _, pair := range pairs {
			// 只获取结算货币为 USDT 的交易对
			if pair.Quote != Settle {
				continue
			}
			specTmp := &symbolSpec{
				Id:              pair.Id,
				Base:            pair.Base,
				Quote:           pair.Quote,
				MinBaseAmount:   pair.MinBaseAmount,
				MinQuoteAmount:  pair.MinQuoteAmount,
				MaxBaseAmount:   pair.MaxBaseAmount,
				MaxQuoteAmount:  pair.MaxQuoteAmount,
				AmountPrecision: int(pair.AmountPrecision),
				Precision:       int(pair.Precision),
				TradeStatus:     pair.TradeStatus,
				SellStart:       pair.SellStart,
				BuyStart:        pair.BuyStart,
				DelistingTime:   pair.DelistingTime,
				TradeUrl:        pair.TradeUrl,
				StTag:           pair.StTag,
			}

			// 如果交易对匹配，则设置交易对规格
			if symbol == pair.Id {
				spec = specTmp
			}
			gateSpotSpec.SetSymbolSpec(pair.Id, specTmp)
		}
	}

	if spec == nil {
		return nil, fmt.Errorf("交易对规格不存在: %s", symbol)
	}
	return spec, nil
}

// GetSpotSymbols 获取现货交易对列表
func (g *gateExchange) GetSpotSymbols() []string {
	var symbols []string
	pairs, _, err := g.client.SpotApi.ListCurrencyPairs(context.Background())
	if err != nil {
		return symbols
	}
	for _, pair := range pairs {
		// 只获取结算货币为 USDT 的交易对
		if pair.Quote != Settle {
			continue
		}
		symbols = append(symbols, pair.Id)
	}
	return symbols
}
