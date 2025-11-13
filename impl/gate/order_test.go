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
