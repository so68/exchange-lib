package binance

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
// go test -v ./impl/binance -run "^TestSpotCreateOrder$" -args --symbol=ETHUSDT --amount=20.0
// 卖出 ETHUSDT， amount 为 ETH 数量 没有指定价格, 则使用最新价格 --lastPrice=3200 挂单
// go test -v ./impl/binance -run "^TestSpotCreateOrder$" -args --symbol=ETHUSDT --side=SELL --amount=20.0
func TestSpotCreateOrder(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	binanceExchange := NewBinance(apiKey, secretKey)

	// 如果最新价格不存在, 或 等于 0
	if *lastPrice == "" || *lastPrice == "0" {
		tickers, err := binanceExchange.GetSpotSymbolTickers(ctx, *symbol)
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
	order, err := binanceExchange.CreateSpotOrder(ctx, *symbol, exchange.OrderSide(*side), *lastPrice, quantity)
	if err != nil {
		t.Fatalf("现货下单失败: %v", err)
	}
	fmt.Printf("【Binance】创建现货订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestSpotGetOrder 获取现货订单
// go test -v ./impl/binance -run "^TestSpotGetOrder$" -args --orderID=38636974388
func TestSpotGetOrder(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	order, err := binanceExchange.GetSpotOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("现货获取订单失败: %v", err)
	}
	fmt.Printf("【Binance】获取现货订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestSpotCancelOrder 撤销现货订单
// go test -v ./impl/binance -run "^TestSpotCancelOrder$" -args --orderID=38636974388
func TestSpotCancelOrder(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	order, err := binanceExchange.CancelSpotOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("现货撤销订单失败: %v", err)
	}
	fmt.Printf("【Binance】撤销现货订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesCreateOrder 创建合约订单
// go test -v ./impl/binance -run "^TestFuturesCreateOrder$" -args --symbol=ETHUSDT --side=BUY --leverage=10 --amount=1000
func TestFuturesCreateOrder(t *testing.T) {
	flag.Parse()

	ctx := context.Background()
	ctx = exchange.WithLeverage(ctx, *leverage)
	ctx = exchange.WithMarginType(ctx, exchange.MarginTypeIsolated)
	binanceExchange := NewBinance(apiKey, secretKey)

	// 如果最新价格不存在, 或 等于 0
	if *lastPrice == "" || *lastPrice == "0" {
		tickers, err := binanceExchange.GetFuturesSymbolTickers(ctx, *symbol)
		if err != nil {
			t.Fatalf("获取交易对行情失败: %v", err)
		}
		*lastPrice = tickers.GetTicker(*symbol).LastPrice
	}

	*amount = *amount * float64(*leverage) // 数量乘以杠杆
	quantity := utils.AmountWithPriceToQuantity(*amount, *lastPrice, 8)

	order, err := binanceExchange.CreateFuturesOrder(ctx, *symbol, exchange.OrderSide(*side), *lastPrice, quantity)
	if err != nil {
		t.Fatalf("合约下单失败: %v", err)
	}
	fmt.Printf("【Binance】创建合约订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesGetOrder 获取合约订单
// go test -v ./impl/binance -run "^TestFuturesGetOrder$" -args --symbol=ETHUSDT --orderID=38636974388
func TestFuturesGetOrder(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	order, err := binanceExchange.GetFuturesOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("合约获取订单失败: %v", err)
	}
	fmt.Printf("【Binance】获取合约订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestFuturesGetPositionRisk 获取合约持仓风险
// go test -v ./impl/binance -run "^TestFuturesGetPositionRisk$" -args --symbol=ETHUSDT
func TestFuturesGetPositionRisk(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	positions, err := binanceExchange.GetFuturesPositionRisk(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("合约获取持仓失败: %v", err)
	}
	for _, position := range positions.Data {
		fmt.Printf("【Binance】获取合约持仓|交易对: %s, 方向: %s, 持仓数量: %s, 平均开仓价格: %s, 标记价格: %s, 未实现盈亏: %s, 杠杆倍数: %s, 爆仓价格: %s, 保证金模式: %s, 逐仓保证金金额: %s, 名义价值: %s\n", position.Symbol, position.PositionSide, position.PositionAmt, position.EntryPrice, position.MarkPrice, position.UnRealizedProfit, position.Leverage, position.LiquidationPrice, position.MarginType, position.IsolatedMargin, position.Notional)
	}
}

// TestFuturesSetSLTP 设置合约止损止盈
// go test -v ./impl/binance -run "^TestFuturesSetSLTP$" -args --symbol=ETHUSDT
func TestFuturesSetSLTP(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	err := binanceExchange.SetFuturesSLTP(context.Background(), *symbol, exchange.PositionSideShort, "3500", "3400")
	if err != nil {
		t.Fatalf("合约设置止损止盈失败: %v", err)
	}
	fmt.Printf("【Binance】合约设置止损止盈成功")
}

// TestFuturesCancelSLTP 撤销合约止损止盈
// go test -v ./impl/binance -run "^TestFuturesCancelSLTP$" -args --symbol=ETHUSDT
func TestFuturesCancelSLTP(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	err := binanceExchange.CancelFuturesSLTP(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("合约撤销止损止盈失败: %v", err)
	}
	fmt.Printf("【Binance】合约撤销止损止盈成功")
}

// TestFuturesClosePositionRisk 平仓合约持仓风险
// go test -v ./impl/binance -run "^TestFuturesClosePositionRisk$" -args --symbol=ETHUSDT --side=SHORT
func TestFuturesClosePositionRisk(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	err := binanceExchange.CloseFuturesPositionRisk(context.Background(), *symbol, exchange.PositionSide(*side))
	if err != nil {
		t.Fatalf("合约平仓持仓风险失败: %v", err)
	}
	fmt.Printf("【Binance】合约平仓持仓风险成功")
}

// TestFuturesCancelOrder 撤销合约订单
// go test -v ./impl/binance -run "^TestFuturesCancelOrder$" -args --symbol=ETHUSDT --orderID=38636974388
func TestFuturesCancelOrder(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance(apiKey, secretKey)
	order, err := binanceExchange.CancelFuturesOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("合约撤销订单失败: %v", err)
	}
	fmt.Printf("【Binance】撤销合约订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}
