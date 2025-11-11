package binance

import (
	"context"
	"flag"
	"fmt"
	"testing"
)

// TestGetSymbolTickers 获取交易对行情
// go test -v ./impl/binance -run "^TestGetSymbolTickers$" -args --symbol=BTCUSDT
func TestGetSymbolTickers(t *testing.T) {
	flag.Parse()

	binanceExchange := NewBinance("", "")
	tickers, err := binanceExchange.GetSymbolTickers(context.Background(), *symbol)
	if err != nil {
		t.Fatalf("获取交易对行情失败: %v", err)
	}
	fmt.Println("tickers", tickers.GetTicker(*symbol))
}
