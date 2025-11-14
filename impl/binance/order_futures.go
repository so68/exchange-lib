package binance

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/so68/exchange-lib/exchange"
)

// CreateFuturesOrder 合约下单
func (b *binanceExchange) CreateFuturesOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	// 获取交易规则
	spec, err := b.getFuturesSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("获取交易规则失败: %w", err)
	}

	// 验证交易规则
	quantity, err = b.filtersQuantity(spec, limitPrice, quantity)
	if err != nil {
		return nil, fmt.Errorf("验证交易规则失败: %w", err)
	}

	service := b.futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Side(futures.SideType(string(side))).
		Quantity(quantity)
	if side == exchange.OrderSideBuy {
		service.PositionSide(futures.PositionSideTypeLong)
	} else {
		service.PositionSide(futures.PositionSideTypeShort)
	}

	// 市价单
	if limitPrice == "" || limitPrice == "0" {
		service.Type(futures.OrderTypeMarket)
	} else {
		service.Type(futures.OrderTypeLimit).Price(limitPrice).TimeInForce(futures.TimeInForceTypeGTC)
	}

	resp, err := service.Do(ctx)
	if err != nil {
		return nil, err
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
		QuoteQuantity: resp.CumQuote,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.UpdateTime,
		UpdateTime:    resp.UpdateTime,
	}, nil
}

// GetFuturesOrder 获取合约订单
func (b *binanceExchange) GetFuturesOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}
	resp, err := b.futuresClient.NewGetOrderService().Symbol(symbol).OrderID(orderIDInt).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance futures get order: %w", err)
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
		QuoteQuantity: resp.CumQuote,
		ActualQty:     resp.CumQuantity,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.UpdateTime,
		UpdateTime:    resp.UpdateTime,
	}, nil
}

// GetFuturesPositionRisk 获取合约持仓风险
func (b *binanceExchange) GetFuturesPositionRisk(ctx context.Context, symbol string) (*exchange.SymbolPositionRisk, error) {
	positions, err := b.futuresClient.NewGetPositionRiskService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取持仓风险失败: %w", err)
	}

	data := &exchange.SymbolPositionRisk{}
	for _, p := range positions {
		data.Data = append(data.Data, &exchange.PositionRisk{
			Symbol:           p.Symbol,
			PositionSide:     exchange.PositionSide(p.PositionSide),
			PositionAmt:      p.PositionAmt,
			EntryPrice:       p.EntryPrice,
			MarkPrice:        p.MarkPrice,
			UnRealizedProfit: p.UnRealizedProfit,
			Leverage:         p.Leverage,
			LiquidationPrice: p.LiquidationPrice,
			MarginType:       p.MarginType,
			IsolatedMargin:   p.IsolatedMargin,
			Notional:         p.Notional,
		})
	}
	return data, nil
}

// SetFuturesSLTP 设置合约止损止盈
func (b *binanceExchange) SetFuturesSLTP(ctx context.Context, symbol string, positionSide exchange.PositionSide, stopPrice string, takeProfitPrice string) error {
	// 获取合约持仓风险
	positionRisk, err := b.GetFuturesPositionRisk(ctx, symbol)
	if err != nil {
		return err
	}

	// 获取指定方向持仓风险
	sidePositionRisk := positionRisk.GetSidePositionRisk(positionSide)
	if sidePositionRisk == nil {
		return fmt.Errorf("获取指定方向 %s 持仓风险失败: 未找到该方向的持仓", positionSide)
	}

	qtyFloat, _ := strconv.ParseFloat(sidePositionRisk.PositionAmt, 64)
	qtyAbs := fmt.Sprintf("%.8f", math.Abs(qtyFloat))

	// 设置止损(STOP_MARKET: 市价止损)
	if stopPrice != "" {
		if err := b.setFuturesStopLoss(ctx, symbol, positionSide, qtyAbs, stopPrice); err != nil {
			return err
		}
	}

	// 设置止盈(TAKE_PROFIT_MARKET: 市价止盈)
	if takeProfitPrice != "" {
		if err := b.setFuturesTakeProfit(ctx, symbol, positionSide, qtyAbs, takeProfitPrice); err != nil {
			return err
		}
	}
	return nil
}

// SetFuturesLeverage 设置合约杠杆
func (b *binanceExchange) SetFuturesLeverage(ctx context.Context, symbol string, leverage int) error {
	if _, err := b.futuresClient.NewChangeLeverageService().
		Symbol(symbol).
		Leverage(leverage).
		Do(ctx); err != nil {
		return fmt.Errorf("设置杠杆失败: %w", err)
	}
	return nil
}

// SetFuturesMarginMode 设置合约保证金模式
func (b *binanceExchange) SetFuturesMarginMode(ctx context.Context, symbol string, marginMode exchange.MarginMode) error {
	if err := b.futuresClient.NewChangeMarginTypeService().
		Symbol(symbol).
		MarginType(futures.MarginType(string(marginMode))).
		Do(ctx); err != nil {
		return fmt.Errorf("设置保证金模式失败: %w", err)
	}
	return nil
}

// CancelFuturesSLTP 撤销合约止损止盈
func (b *binanceExchange) CancelFuturesSLTP(ctx context.Context, symbol string) error {
	// 获取该 symbol 的所有开放订单
	openOrders, err := b.futuresClient.NewListOpenOrdersService().Symbol(symbol).Do(ctx)
	if err != nil {
		return fmt.Errorf("获取 %s 开放订单失败: %w", symbol, err)
	}

	// 3. 遍历并取消除已成交止损之外的 TP/SL
	for _, o := range openOrders {
		// 检查订单类型是否为止损或止盈
		if o.Type == futures.OrderTypeTakeProfit || o.Type == futures.OrderTypeTakeProfitMarket || o.Type == futures.OrderTypeStop || o.Type == futures.OrderTypeStopMarket {
			// 如果是 STOP_MARKET 且已成交（Status = FILLED），跳过
			if o.Type == futures.OrderTypeStopMarket && o.Status == futures.OrderStatusTypeFilled {
				continue
			}

			// 取消订单
			_, err = b.futuresClient.NewCancelOrderService().
				Symbol(symbol).
				OrderID(o.OrderID).
				Do(ctx)
			if err != nil {
				return fmt.Errorf("取消订单 %d 失败: %w", o.OrderID, err)
			}
		}
	}
	return nil
}

// CloseFuturesPositionRisk 平仓合约持仓风险
func (b *binanceExchange) CloseFuturesPositionRisk(ctx context.Context, symbol string, positionSide exchange.PositionSide) error {
	// 获取合约持仓风险
	positionRisk, err := b.GetFuturesPositionRisk(ctx, symbol)
	if err != nil {
		return err
	}

	// 获取指定方向持仓风险
	sidePositionRisk := positionRisk.GetSidePositionRisk(positionSide)
	if sidePositionRisk == nil {
		return fmt.Errorf("获取指定方向 %s 持仓风险失败: 未找到该方向的持仓", positionSide)
	}

	// 如果持仓风险为0，则不进行平仓
	if sidePositionRisk.PositionAmt == "0" {
		return nil
	}

	// 确定平仓方向：LONG 持仓用 SELL 平仓，SHORT 持仓用 BUY 平仓
	amt, _ := strconv.ParseFloat(sidePositionRisk.PositionAmt, 64)
	side := futures.SideTypeSell
	if positionSide == exchange.PositionSideShort {
		side = futures.SideTypeBuy
		amt = -amt
	}

	// 市价平仓
	_, err = b.futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(futures.OrderTypeMarket).
		Quantity(fmt.Sprintf("%.8f", amt)).                           // 必须传持仓数量
		PositionSide(futures.PositionSideType(string(positionSide))). // 持仓方向
		Do(ctx)

	if err != nil {
		return fmt.Errorf("平仓失败: %w", err)
	}
	return nil
}

// CancelFuturesOrder 撤销合约订单
func (b *binanceExchange) CancelFuturesOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	orderIDInt, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的订单ID: %w", err)
	}
	resp, err := b.futuresClient.NewCancelOrderService().Symbol(symbol).OrderID(orderIDInt).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance futures cancel order: %w", err)
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
		QuoteQuantity: resp.CumQuote,
		ActualQty:     resp.CumQuantity,
		TimeInForce:   exchange.OrderTimeInForce(resp.TimeInForce),
		CreateTime:    resp.UpdateTime,
		UpdateTime:    resp.UpdateTime,
	}, nil
}

// SetFuturesDualMode 设置持仓模式
func (b *binanceExchange) SetFuturesDualMode(ctx context.Context, dualMode bool) error {
	if err := b.futuresClient.NewChangePositionModeService().
		DualSide(dualMode).
		Do(ctx); err != nil {
		return fmt.Errorf("设置持仓模式失败: %w", err)
	}
	return nil
}

// setFuturesStopLoss 设置止损
func (b *binanceExchange) setFuturesStopLoss(ctx context.Context, symbol string, positionSide exchange.PositionSide, quantity string, stopPriceSL string) error {
	// 确定平仓方向：LONG 持仓用 SELL 平仓，SHORT 持仓用 BUY 平仓
	side := futures.SideTypeSell
	if positionSide == exchange.PositionSideShort {
		side = futures.SideTypeBuy
	}
	// 设置止损 (STOP_MARKET: 市价止损)
	_, err := b.futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Side(side).                        // 平仓方向
		Type(futures.OrderTypeStopMarket). // 市价止损类型
		// ClosePosition(true).                                          // 关闭整个持仓
		Quantity(quantity).                                           // 平仓数量
		StopPrice(stopPriceSL).                                       // 触发价
		PositionSide(futures.PositionSideType(string(positionSide))). // 持仓方向
		WorkingType(futures.WorkingTypeMarkPrice).                    // 使用 Mark Price
		Do(ctx)

	if err != nil {
		return fmt.Errorf("设置止损失败: %w", err)
	}
	return nil
}

// setFuturesTakeProfit 设置止盈
func (b *binanceExchange) setFuturesTakeProfit(ctx context.Context, symbol string, positionSide exchange.PositionSide, quantity string, stopPriceTP string) error {
	// 确定平仓方向：多头持仓用 SELL 平仓，空头持仓用 BUY 平仓
	side := futures.SideTypeSell
	if positionSide == exchange.PositionSideShort {
		side = futures.SideTypeBuy
	}

	_, err := b.futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Side(side).                              // 平仓方向
		Type(futures.OrderTypeTakeProfitMarket). // 市价止盈类型
		// ClosePosition(true).                                          // 关闭整个持仓
		Quantity(quantity).                                           // 平仓数量
		StopPrice(stopPriceTP).                                       // 触发价
		PositionSide(futures.PositionSideType(string(positionSide))). // 持仓方向
		WorkingType(futures.WorkingTypeMarkPrice).                    // 使用 Mark Price
		Do(ctx)

	if err != nil {
		return fmt.Errorf("设置止盈失败: %w", err)
	}
	return nil
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
			return nil, err
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
					specTmp.MinNotional = f["notional"].(string)
				}
			}

			// 找到对应的交易对规格
			if s.Symbol == symbol {
				spec = specTmp
			}

			binanceFuturesSpec.SetSymbolSpec(s.Symbol, specTmp)
		}
	}

	if spec == nil {
		return nil, fmt.Errorf("合约规格不存在: %s", symbol)
	}

	return spec, nil
}
