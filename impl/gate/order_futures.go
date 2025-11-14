package gate

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"github.com/so68/exchange-lib/exchange"
)

// CreateFuturesOrder 创建合约订单 - amount 金额 * 杠杆
func (g *gateExchange) CreateFuturesOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, amount string) (*exchange.Order, error) {
	// 获取交易规则
	spec, err := g.GetFuturesSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("获取交易规则失败: %w", err)
	}

	// 验证交易规则
	size, err := g.filtersFuturesSize(spec, limitPrice, amount)
	if err != nil {
		return nil, fmt.Errorf("验证交易规则失败: %w", err)
	}

	// 如果方向为卖出，则取反
	if side == exchange.OrderSideSell {
		size = -size
	}

	// 创建订单参数
	orderParams := gateapi.FuturesOrder{
		Contract: symbol,                                                // 合约符号
		Size:     size,                                                  // 数量，正数=买多/开多，负数=卖空/开空（例如 -10）
		Price:    limitPrice,                                            // 价格（限价单，字符串类型）
		Tif:      strings.ToLower(string(exchange.OrderTimeInForceGTC)), // 时间有效性
	}

	// 创建订单
	createdOrder, _, err := g.client.FuturesApi.CreateFuturesOrder(ctx, strings.ToLower(Settle), orderParams, nil)
	if err != nil {
		return nil, fmt.Errorf("合约下单失败: %w", err)
	}

	orderType := exchange.OrderTypeLimit
	if limitPrice == "0" || limitPrice == "" {
		orderType = exchange.OrderTypeMarket
	}

	// 转换数量为字符串
	sizeStr := strconv.FormatInt(createdOrder.Size, 10)
	return &exchange.Order{
		OrderID:       strconv.FormatInt(createdOrder.Id, 10),
		Symbol:        createdOrder.Contract,
		Side:          side,
		Type:          orderType,
		Status:        exchange.OrderStatusNew, // 默认新订单
		Price:         createdOrder.FillPrice,
		Quantity:      sizeStr,
		ExecutedQty:   sizeStr,
		ActualQty:     sizeStr,
		QuoteQuantity: "0",
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(createdOrder.Tif)),
		CreateTime:    int64(createdOrder.CreateTime),
		UpdateTime:    int64(createdOrder.FinishTime),
	}, nil
}

// GetFuturesOrder 获取合约订单
func (g *gateExchange) GetFuturesOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	order, _, err := g.client.FuturesApi.GetFuturesOrder(ctx, strings.ToLower(Settle), orderID)
	if err != nil {
		return nil, fmt.Errorf("获取合约订单失败: %w", err)
	}

	status := exchange.OrderStatusNew
	switch order.FinishAs {
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

	side := exchange.OrderSideBuy
	if order.Size < 0 {
		side = exchange.OrderSideSell
	}

	// 转换数量为字符串
	sizeStr := strconv.FormatInt(order.Size, 10)
	return &exchange.Order{
		OrderID:       strconv.FormatInt(order.Id, 10),
		Symbol:        order.Contract,
		Side:          side,
		Type:          exchange.OrderTypeLimit,
		Status:        status,
		Price:         order.FillPrice,
		Quantity:      sizeStr,
		ExecutedQty:   sizeStr,
		ActualQty:     sizeStr,
		QuoteQuantity: "0",
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(order.Tif)),
		CreateTime:    int64(order.CreateTime),
		UpdateTime:    int64(order.FinishTime),
	}, nil
}

// CancelFuturesOrder 取消合约订单
func (g *gateExchange) CancelFuturesOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	canceledOrder, _, err := g.client.FuturesApi.CancelFuturesOrder(ctx, strings.ToLower(Settle), orderID, nil)
	if err != nil {
		return nil, fmt.Errorf("取消合约订单失败: %w", err)
	}

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

	side := exchange.OrderSideBuy
	if canceledOrder.Size < 0 {
		side = exchange.OrderSideSell
	}

	// 转换数量为字符串
	sizeStr := strconv.FormatInt(canceledOrder.Size, 10)
	return &exchange.Order{
		OrderID:       strconv.FormatInt(canceledOrder.Id, 10),
		Symbol:        canceledOrder.Contract,
		Side:          side,
		Type:          exchange.OrderTypeLimit,
		Status:        status,
		Price:         canceledOrder.FillPrice,
		Quantity:      sizeStr,
		ExecutedQty:   sizeStr,
		ActualQty:     sizeStr,
		QuoteQuantity: "0",
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(canceledOrder.Tif)),
		CreateTime:    int64(canceledOrder.CreateTime),
		UpdateTime:    int64(canceledOrder.FinishTime),
	}, nil
}

// GetFuturesPositionRisk 获取合约持仓风险
func (g *gateExchange) GetFuturesPositionRisk(ctx context.Context, symbol string) (*exchange.SymbolPositionRisk, error) {
	positions, _, err := g.client.FuturesApi.GetDualModePosition(ctx, strings.ToLower(Settle), symbol)
	if err != nil {
		return nil, fmt.Errorf("获取合约持仓风险失败: %w", err)
	}

	exchangePositionRisk := &exchange.SymbolPositionRisk{
		Data: []*exchange.PositionRisk{},
	}
	for _, position := range positions {
		fmt.Printf("position: %+v\n", position)
		if position.Size != 0 && position.Contract == symbol {
			side := exchange.PositionSideLong
			if position.Size < 0 {
				side = exchange.PositionSideShort
			}

			marginMode := exchange.MarginModeIsolated
			if position.Leverage == "0" {
				marginMode = exchange.MarginModeCrossed
			}

			exchangePositionRisk.Data = append(exchangePositionRisk.Data, &exchange.PositionRisk{
				Symbol:           position.Contract,
				PositionSide:     side,
				PositionAmt:      strconv.FormatInt(position.Size, 10),
				EntryPrice:       position.EntryPrice,
				MarkPrice:        position.MarkPrice,
				UnRealizedProfit: position.UnrealisedPnl,
				Leverage:         position.Leverage,
				LiquidationPrice: position.LiqPrice,
				MarginType:       string(marginMode),
				IsolatedMargin:   position.Margin,
				Notional:         position.Value,
			})
		}
	}
	return exchangePositionRisk, nil
}

func (g *gateExchange) CloseFuturesPositionRisk(ctx context.Context, symbol string, positionSide exchange.PositionSide) error {
	return nil
}

// SetFuturesSLTP 设置合约止损止盈
func (g *gateExchange) SetFuturesSLTP(ctx context.Context, symbol string, positionSide exchange.PositionSide, stopPrice, takeProfitPrice string) error {
	// 获取合约持仓风险
	positionRisk, err := g.GetFuturesPositionRisk(ctx, symbol)
	if err != nil {
		return err
	}

	// 获取指定方向持仓风险
	sidePositionRisk := positionRisk.GetSidePositionRisk(positionSide)
	if sidePositionRisk == nil {
		return fmt.Errorf("获取指定方向 %s 持仓风险失败: 未找到该方向的持仓", positionSide)
	}

	// 设置止损(STOP_MARKET: 市价止损)
	if stopPrice != "" {
		if err := g.setFuturesStopLoss(ctx, symbol, positionSide, sidePositionRisk.PositionAmt, stopPrice); err != nil {
			return err
		}
	}

	// 设置止盈(TAKE_PROFIT_MARKET: 市价止盈)
	if takeProfitPrice != "" {
		if err := g.setFuturesTakeProfit(ctx, symbol, positionSide, sidePositionRisk.PositionAmt, takeProfitPrice); err != nil {
			return err
		}
	}
	return nil
}

// SetFuturesLeverage 设置合约杠杆
func (g *gateExchange) SetFuturesLeverage(ctx context.Context, symbol string, leverage int) error {
	_, _, err := g.client.FuturesApi.UpdateDualModePositionLeverage(ctx, strings.ToLower(Settle), symbol, strconv.Itoa(leverage), nil)
	if err != nil {
		return fmt.Errorf("更新杠杆失败: %w", err)
	}
	return nil
}

// SetFuturesMarginMode 设置合约保证金模式
func (g *gateExchange) SetFuturesMarginMode(ctx context.Context, symbol string, marginMode exchange.MarginMode) error {
	mode := "ISOLATED"
	if marginMode == exchange.MarginModeCrossed {
		mode = "CROSS"
	}
	_, _, err := g.client.FuturesApi.UpdateDualCompPositionCrossMode(ctx, strings.ToLower(Settle), gateapi.InlineObject{
		Mode:     mode,
		Contract: symbol,
	})
	if err != nil {
		return fmt.Errorf("设置保证金模式失败: %w", err)
	}
	return nil
}

// SetFuturesDualMode 设置持仓模式
func (g *gateExchange) SetFuturesDualMode(ctx context.Context, dualMode bool) error {
	_, _, err := g.client.FuturesApi.SetDualMode(ctx, strings.ToLower(Settle), dualMode)
	if err != nil {
		return fmt.Errorf("设置持仓模式失败: %w", err)
	}
	return nil
}

// CancelFuturesSLTP 撤销合约止损止盈
func (g *gateExchange) CancelFuturesSLTP(ctx context.Context, symbol string) error {
	opts := &gateapi.CancelPriceTriggeredOrderListOpts{
		Contract: optional.NewString(symbol),
	}
	_, _, err := g.client.FuturesApi.CancelPriceTriggeredOrderList(ctx, strings.ToLower(Settle), opts)
	if err != nil {
		return fmt.Errorf("撤销合约止损止盈失败: %w", err)
	}
	return nil
}

// GetFuturesSymbolSpec 获取合约交易对规格
func (g *gateExchange) GetFuturesSymbolSpec(ctx context.Context, symbol string) (*futuresSpec, error) {
	var spec *futuresSpec

	// 从缓存中获取交易对规格
	spec, _ = gateFuturesSpec.GetFuturesSpec(symbol)

	if spec == nil {
		contracts, _, err := g.client.FuturesApi.ListFuturesContracts(context.Background(), strings.ToLower(Settle), nil)
		if err != nil {
			return nil, fmt.Errorf("获取合约交易对规则失败: %w", err)
		}

		for _, contract := range contracts {
			specTmp := &futuresSpec{
				Name:              contract.Name,
				Type:              contract.Type,
				QuantoMultiplier:  contract.QuantoMultiplier,
				LeverageMin:       contract.LeverageMin,
				LeverageMax:       contract.LeverageMax,
				MaintenanceRate:   contract.MaintenanceRate,
				MarkType:          contract.MarkType,
				MarkPrice:         contract.MarkPrice,
				IndexPrice:        contract.IndexPrice,
				LastPrice:         contract.LastPrice,
				MakerFeeRate:      contract.MakerFeeRate,
				TakerFeeRate:      contract.TakerFeeRate,
				OrderPriceRound:   contract.OrderPriceRound,
				MarkPriceRound:    contract.MarkPriceRound,
				FundingRate:       contract.FundingRate,
				FundingInterval:   contract.FundingInterval,
				FundingNextApply:  contract.FundingNextApply,
				OrderSizeMin:      contract.OrderSizeMin,
				OrderSizeMax:      contract.OrderSizeMax,
				OrderPriceDeviate: contract.OrderPriceDeviate,
				RefDiscountRate:   contract.RefDiscountRate,
				RefRebateRate:     contract.RefRebateRate,
				OrderbookId:       contract.OrderbookId,
				TradeId:           contract.TradeId,
				TradeSize:         contract.TradeSize,
				PositionSize:      contract.PositionSize,
				ConfigChangeTime:  contract.ConfigChangeTime,
				InDelisting:       contract.InDelisting,
				OrdersLimit:       contract.OrdersLimit,
				EnableBonus:       contract.EnableBonus,
				EnableCredit:      contract.EnableCredit,
				CreateTime:        contract.CreateTime,
				FundingCapRatio:   contract.FundingCapRatio,
				Status:            contract.Status,
				LaunchTime:        contract.LaunchTime,
			}

			// 如果合约匹配，则设置合约规格
			if symbol == contract.Name {
				spec = specTmp
			}
			gateFuturesSpec.SetFuturesSpec(contract.Name, specTmp)
		}
	}

	if spec == nil {
		return nil, fmt.Errorf("合约规格不存在: %s", symbol)
	}

	return spec, nil
}

// GetFuturesSymbols 获取合约交易对列表
func (g *gateExchange) GetFuturesSymbols() []string {
	var symbols []string
	contracts, _, err := g.client.FuturesApi.ListFuturesContracts(context.Background(), strings.ToLower(Settle), nil)
	if err != nil {
		return symbols
	}
	for _, contract := range contracts {
		symbols = append(symbols, contract.Name)
	}
	return symbols
}

// setFuturesStopLoss 设置止损
func (g *gateExchange) setFuturesStopLoss(ctx context.Context, symbol string, positionSide exchange.PositionSide, quantity string, stopPriceSL string) error {
	var rule int32 = 2
	if positionSide == exchange.PositionSideShort {
		rule = 1
	}
	quantityInt, _ := strconv.ParseInt(quantity, 10, 64)
	slInitial := gateapi.FuturesInitialOrder{
		Contract:   symbol,
		Size:       quantityInt,
		Price:      "0", // 市价
		Tif:        strings.ToLower(string(exchange.OrderTimeInForceIOC)),
		ReduceOnly: true,
	}
	slTrigger := gateapi.FuturesPriceTrigger{
		PriceType: 2,           // 标记价格触发
		Price:     stopPriceSL, // 触发价
		Rule:      rule,
	}
	slOrder := gateapi.FuturesPriceTriggeredOrder{
		Initial: slInitial,
		Trigger: slTrigger,
	}

	// 设置止损
	_, _, err := g.client.FuturesApi.CreatePriceTriggeredOrder(ctx, strings.ToLower(Settle), slOrder)
	if err != nil {
		return fmt.Errorf("设置止损失败: %w", err)
	}
	return err
}

// setFuturesTakeProfit 设置止盈
func (g *gateExchange) setFuturesTakeProfit(ctx context.Context, symbol string, positionSide exchange.PositionSide, quantity string, stopPriceTP string) error {
	var rule int32 = 1
	if positionSide == exchange.PositionSideShort {
		rule = 2
	}

	quantityInt, _ := strconv.ParseInt(quantity, 10, 64)
	tpInitial := gateapi.FuturesInitialOrder{
		Contract:   symbol,
		Size:       quantityInt,
		Price:      "0",
		Tif:        strings.ToLower(string(exchange.OrderTimeInForceIOC)),
		ReduceOnly: true,
	}
	tpTrigger := gateapi.FuturesPriceTrigger{
		PriceType: 2,           // 标记价格触发
		Price:     stopPriceTP, // 触发价
		Rule:      rule,
	}
	tpOrder := gateapi.FuturesPriceTriggeredOrder{
		Initial: tpInitial,
		Trigger: tpTrigger,
	}

	// 设置止盈
	_, _, err := g.client.FuturesApi.CreatePriceTriggeredOrder(ctx, strings.ToLower(Settle), tpOrder)
	if err != nil {
		return fmt.Errorf("设置止盈失败: %w", err)
	}
	return err
}
