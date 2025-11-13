package gate

import (
	"context"
	"flag"
	"fmt"
	"testing"
)

// TestGetSpotSymbolTickers 获取现货交易对行情
// go test -v ./impl/gate -run "^TestGetSpotSymbolTickers$" -args --symbol=BTCUSDT
func TestGetSpotSymbolTickers(t *testing.T) {
	flag.Parse()

	gateExchange := NewGateExchange("", "")
	tickers, err := gateExchange.GetSpotSymbolTickers(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("获取交易对行情失败: %v", err)
	}
	for _, ticker := range tickers.Tickers {
		fmt.Println(ticker.Symbol, ticker.LastPrice, ticker.PriceChangePercent)
	}
}

// TestGetFuturesSymbolTickers 获取合约交易对行情
// go test -v ./impl/gate -run "^TestGetFuturesSymbolTickers$" -args --symbol=BTCUSDT
func TestGetFuturesSymbolTickers(t *testing.T) {
	flag.Parse()

	gateExchange := NewGateExchange("", "")
	tickers, err := gateExchange.GetFuturesSymbolTickers(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("获取交易对行情失败: %v", err)
	}
	for _, ticker := range tickers.Tickers {
		fmt.Println(ticker.Symbol, ticker.LastPrice, ticker.PriceChangePercent)
	}
}
