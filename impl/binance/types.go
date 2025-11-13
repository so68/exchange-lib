package binance

import "time"

var binanceSpotSpec *exchangeSpec
var binanceFuturesSpec *exchangeSpec

func init() {
	binanceSpotSpec = &exchangeSpec{
		Symbols:    make([]*symbolSpec, 0),
		UpdateTime: time.Now(),
	}
	binanceFuturesSpec = &exchangeSpec{
		Symbols:    make([]*symbolSpec, 0),
		UpdateTime: time.Now(),
	}

	// 定时更新交易对规格
	initSymbolsSpec()
}

// reloadSymbolsSpec 重新加载交易对规格
func initSymbolsSpec() {
	binanceSpotSpec.DeleteSymbolsSpec()
	binanceFuturesSpec.DeleteSymbolsSpec()

	time.AfterFunc(time.Hour, initSymbolsSpec)
}
