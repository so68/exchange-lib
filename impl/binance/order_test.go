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
	binanceExchange := NewBinance("uAF7jJkOMygW8VK24OEzrrRNYjDFierW04lutVKtPJT9kLUNfHsByLe7lM7dWi4y", "jscrjDQsTpsiL8nDExD8Ty52raja6g8f4d0VOXByQfyKrYAmYtDUUCDSxymhmvgu")

	// 如果最新价格不存在, 或 等于 0
	if *lastPrice == "" || *lastPrice == "0" {
		tickers, err := binanceExchange.GetSymbolTickers(ctx, *symbol)
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
	order, err := binanceExchange.SpotCreateOrder(ctx, *symbol, exchange.OrderSide(*side), *lastPrice, quantity)
	if err != nil {
		t.Fatalf("创建现货订单失败: %v", err)
	}
	fmt.Printf("【Binance】创建现货订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestSpotGetOrder 获取现货订单
// go test -v ./impl/binance -run "^TestSpotGetOrder$" -args --orderID=38636974388
func TestSpotGetOrder(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance("uAF7jJkOMygW8VK24OEzrrRNYjDFierW04lutVKtPJT9kLUNfHsByLe7lM7dWi4y", "jscrjDQsTpsiL8nDExD8Ty52raja6g8f4d0VOXByQfyKrYAmYtDUUCDSxymhmvgu")
	order, err := binanceExchange.SpotGetOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("获取现货订单失败: %v", err)
	}
	fmt.Printf("【Binance】获取现货订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}

// TestSpotCancelOrder 撤销现货订单
// go test -v ./impl/binance -run "^TestSpotCancelOrder$" -args --orderID=38636974388
func TestSpotCancelOrder(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance("uAF7jJkOMygW8VK24OEzrrRNYjDFierW04lutVKtPJT9kLUNfHsByLe7lM7dWi4y", "jscrjDQsTpsiL8nDExD8Ty52raja6g8f4d0VOXByQfyKrYAmYtDUUCDSxymhmvgu")
	order, err := binanceExchange.SpotCancelOrder(context.Background(), *symbol, *orderID)
	if err != nil {
		t.Fatalf("撤销现货订单失败: %v", err)
	}
	fmt.Printf("【Binance】撤销现货订单|订单ID: %d, 交易对: %s, 方向: %s, 类型: %s, 状态: %s, 价格: %s, 数量: %s, 已执行数量: %s, 实际数量: %s, 成交金额: %s, 时间类型: %s, 创建时间: %d, 更新时间: %d\n", order.OrderID, order.Symbol, order.Side, order.Type, order.Status, order.Price, order.Quantity, order.ExecutedQty, order.ActualQty, order.QuoteQuantity, order.TimeInForce, order.CreateTime, order.UpdateTime)
}
