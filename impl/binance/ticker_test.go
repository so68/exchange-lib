package binance

import (
	"context"
	"flag"
	"fmt"
	"testing"
)

// TestGetSpotSymbolTickers 获取现货交易对行情
// go test -v ./impl/binance -run "^TestGetSpotSymbolTickers$" -args --symbol=BTCUSDT
func TestGetSpotSymbolTickers(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance("", "")
	tickers, err := binanceExchange.GetSpotSymbolTickers(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("获取交易对行情失败: %v", err)
	}
	fmt.Println("tickers", tickers.GetTicker(*symbol))
}

// TestGetFuturesSymbolTickers 获取合约交易对行情
// go test -v ./impl/binance -run "^TestGetFuturesSymbolTickers$" -args --symbol=BTCUSDT
func TestGetFuturesSymbolTickers(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance("", "")
	tickers, err := binanceExchange.GetFuturesSymbolTickers(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("获取交易对行情失败: %v", err)
	}
	fmt.Println("tickers", tickers.GetTicker(*symbol))
}
