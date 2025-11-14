package gate

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/so68/exchange-lib/exchange"
	"github.com/so68/exchange-lib/internal/utils"
)

// TestSpotCreateOrder 创建现货订单
// 买入 ETHUSDT 20u  没有指定价格, 则使用最新价格 --lastPrice=3200 挂单
// go test -v ./impl/gate -run "^TestSpotCreateOrder$" -args --symbol=ETHUSDT --amount=20.0
// 卖出 ETHUSDT， amount 为 ETH 数量 没有指定价格, 则使用最新价格 --lastPrice=3200 挂单
// go test -v ./impl/gate -run "^TestSpotCreateOrder$" -args --symbol=ETHUSDT --side=SELL --amount=20.0
func TestSpotCreateOrder(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	gateExchange := NewGateExchange(apiKey, secretKey)
	*symbol = utils.FormatSymbol(*symbol)

	// 如果最新价格不存在, 或 等于 0
	if *lastPrice == "" || *lastPrice == "0" {
		tickers, err := gateExchange.GetSpotSymbolTickers(ctx, *symbol)
		if err != nil {
			t.Fatalf("获取交易对行情失败: %v", err)
		}
		*lastPrice = tickers.GetTicker(*symbol).LastPrice
	}
	quantity := utils.AmountWithPriceToQuantity(*amount, *lastPrice, 8)

	// 如果方向不存在, 则默认买入
	*side = strings.ToUpper(*side)
	if *side != "BUY" {
		*side = "SELL"
		quantity = fmt.Sprintf("%f", *amount)
	}
	order, err := gateExchange.CreateSpotOrder(ctx, *symbol, exchange.OrderSide(*side), *lastPrice, quantity)
	if err != nil {
		t.Fatalf("现货下单失败: %v", err)
	}
	fmt.Printf("【Gate】创建现货订单|订单ID: %s, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestSpotGetOrder 获取现货订单
// go test -v ./impl/gate -run "^TestSpotGetOrder$" -args --symbol=ETHUSDT --orderID=38636974388
func TestSpotGetOrder(t *testing.T) {
	flag.Parse()

	*symbol = utils.FormatSymbol(*symbol)

	ctx := context.Background()
	gateExchange := NewGateExchange(apiKey, secretKey)
	order, err := gateExchange.GetSpotOrder(ctx, *symbol, *orderID)
	if err != nil {
		t.Fatalf("现货获取订单失败: %v", err)
	}
	fmt.Printf("【Binance】获取现货订单|订单ID: %s, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesGetPositionRisk 获取合约持仓风险
// go test -v ./impl/gate -run "^TestFuturesGetPositionRisk$" -args --symbol=ETHUSDT
func TestFuturesGetPositionRisk(t *testing.T) {
	flag.Parse()

	*symbol = utils.FormatSymbol(*symbol)
	gateExchange := NewGateExchange(apiKey, secretKey)

	positions, err := gateExchange.GetFuturesPositionRisk(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("合约获取持仓失败: %v", err)
	}
	for _, position := range positions.Data {
		fmt.Printf("【Binance】获取合约持仓|交易对: %s, 方向: %s, 持仓数量: %s, 平均开仓价格: %s, 标记价格: %s, 未实现盈亏: %s, 杠杆倍数: %s, 爆仓价格: %s, 保证金模式: %s, 逐仓保证金金额: %s, 名义价值: %s\n", position.Symbol, position.PositionSide, position.PositionAmt, position.EntryPrice, position.MarkPrice, position.UnRealizedProfit, position.Leverage, position.LiquidationPrice, position.MarginType, position.IsolatedMargin, position.Notional)
	}
}

// TestSpotCancelOrder 撤销现货订单
// go test -v ./impl/gate -run "^TestSpotCancelOrder$" -args --symbol=ETHUSDT --orderID=38636974388
func TestSpotCancelOrder(t *testing.T) {
	flag.Parse()

	*symbol = utils.FormatSymbol(*symbol)

	ctx := context.Background()
	gateExchange := NewGateExchange(apiKey, secretKey)
	order, err := gateExchange.CancelSpotOrder(ctx, *symbol, *orderID)
	if err != nil {
		t.Fatalf("现货撤销订单失败: %v", err)
	}
	fmt.Printf("【Binance】撤销现货订单|订单ID: %s, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesCreateOrder 创建合约订单
// go test -v ./impl/gate -run "^TestFuturesCreateOrder$" -args --symbol=ETHUSDT --side=BUY --leverage=10 --amount=10
func TestFuturesCreateOrder(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	gateExchange := NewGateExchange(apiKey, secretKey)

	*symbol = utils.FormatSymbol(*symbol)

	// 如果最新价格不存在, 或 等于 0
	if *lastPrice == "" || *lastPrice == "0" {
		tickers, err := gateExchange.GetFuturesSymbolTickers(ctx, *symbol)
		if err != nil {
			t.Fatalf("获取交易对行情失败: %v", err)
		}
		*lastPrice = tickers.GetTicker(*symbol).LastPrice
	}

	*amount = *amount * float64(*leverage) // 数量乘以杠杆
	amountStr := fmt.Sprintf("%f", *amount)

	// 设置杠杆倍数
	err := gateExchange.SetFuturesLeverage(ctx, *symbol, *leverage)
	if err != nil {
		t.Fatalf("设置杠杆倍数失败: %v", err)
	}

	order, err := gateExchange.CreateFuturesOrder(ctx, *symbol, exchange.OrderSide(*side), *lastPrice, amountStr)
	if err != nil {
		t.Fatalf("合约下单失败: %v", err)
	}
	fmt.Printf("【Binance】创建合约订单|订单ID: %s, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesGetOrder 获取合约订单
// go test -v ./impl/gate -run "^TestFuturesGetOrder$" -args --symbol=ETHUSDT --orderID=38636974388
func TestFuturesGetOrder(t *testing.T) {
	flag.Parse()

	gateExchange := NewGateExchange(apiKey, secretKey)
	order, err := gateExchange.GetFuturesOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("合约获取订单失败: %v", err)
	}
	fmt.Printf("【Binance】获取合约订单|订单ID: %s, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesCancelOrder 撤销合约订单
// go test -v ./impl/gate -run "^TestFuturesCancelOrder$" -args --symbol=ETHUSDT --orderID=38636974388
func TestFuturesCancelOrder(t *testing.T) {
	flag.Parse()

	gateExchange := NewGateExchange(apiKey, secretKey)
	order, err := gateExchange.CancelFuturesOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("合约撤销订单失败: %v", err)
	}
	fmt.Printf("【Gate】撤销合约订单|订单ID: %s, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesSetSLTP 设置合约止损止盈
// go test -v ./impl/gate -run "^TestFuturesSetSLTP$" -args --symbol=ETHUSDT
func TestFuturesSetSLTP(t *testing.T) {
	flag.Parse()

	*symbol = utils.FormatSymbol(*symbol)
	gateExchange := NewGateExchange(apiKey, secretKey)
	err := gateExchange.SetFuturesSLTP(context.Background(), *symbol, exchange.PositionSideShort, "3200", "3000")
	if err != nil {
		t.Fatalf("合约设置止损止盈失败: %v", err)
	}
	fmt.Printf("【Gate】合约设置止损止盈成功")
}

// TestFuturesCancelSLTP 撤销合约止损止盈
// go test -v ./impl/gate -run "^TestFuturesCancelSLTP$" -args --symbol=ETHUSDT
func TestFuturesCancelSLTP(t *testing.T) {
	flag.Parse()

	*symbol = utils.FormatSymbol(*symbol)
	gateExchange := NewGateExchange(apiKey, secretKey)
	err := gateExchange.CancelFuturesSLTP(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("合约撤销止损止盈失败: %v", err)
	}
	fmt.Printf("【Gate】合约撤销止损止盈成功")
}
