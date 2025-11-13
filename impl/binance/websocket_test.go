package binance

import (
	"fmt"
	"testing"

	"github.com/so68/exchange-lib/exchange"
)

// TestListenSpotTickers 监听现货交易对行情
// go test -v ./impl/binance -run "^TestListenSpotTickers$"
func TestListenSpotTickers(t *testing.T) {
	binanceWebsocket := NewBinanceWebsocket()
	err := binanceWebsocket.StartListenSpotTickers(func(tickers *exchange.Tickers) {
		for _, ticker := range tickers.Tickers {
			fmt.Println(ticker)
		}
	})
	if err != nil {
		t.Fatalf("Failed to start listen spot tickers: %v", err)
	}

	// 阻塞主线程
	select {}
}

// TestListenFuturesTickers 监听合约交易对行情
// go test -v ./impl/binance -run "^TestListenFuturesTickers$"
func TestListenFuturesTickers(t *testing.T) {
	binanceWebsocket := NewBinanceWebsocket()
	err := binanceWebsocket.StartListenFuturesTickers(func(tickers *exchange.Tickers) {
		for _, ticker := range tickers.Tickers {
			fmt.Println(ticker)
		}
	})
	if err != nil {
		t.Fatalf("Failed to start listen futures tickers: %v", err)
	}

	// 阻塞主线程
	select {}
}
