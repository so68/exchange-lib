package gate

import (
	"fmt"
	"testing"

	"github.com/so68/exchange-lib/exchange"
)

// TestListenSpotTickers 监听现货交易对行情
// go test -v ./impl/gate -run "^TestListenSpotTickers$"
func TestListenSpotTickers(t *testing.T) {
	gateWebsocket := NewGateWebsocket()
	err := gateWebsocket.StartListenSpotTickers(func(ticker *exchange.Ticker) {
		fmt.Println("===>", ticker)
	})
	if err != nil {
		t.Fatalf("Failed to start listen spot tickers: %v", err)
	}

	// 阻塞主线程
	select {}
}

// TestListenFuturesTickers 监听合约交易对行情
// go test -v ./impl/gate -run "^TestListenFuturesTickers$"
func TestListenFuturesTickers(t *testing.T) {
	gateWebsocket := NewGateWebsocket()
	err := gateWebsocket.StartListenFuturesTickers(func(ticker *exchange.Ticker) {
		fmt.Println("===>", ticker)
	})
	if err != nil {
		t.Fatalf("Failed to start listen futures tickers: %v", err)
	}

	// 阻塞主线程
	select {}
}
