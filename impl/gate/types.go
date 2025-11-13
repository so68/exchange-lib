package gate

import "time"

const (
	Settle = "USDT" // 默认结算货币
)

var gateSpotSpec *exchangeSpec
var gateFuturesSpec *exchangeFuturesSpec

func init() {
	gateSpotSpec = &exchangeSpec{
		Symbols:    make([]*symbolSpec, 0),
		UpdateTime: time.Now(),
	}
	gateFuturesSpec = &exchangeFuturesSpec{
		Contracts:  make([]*futuresSpec, 0),
		UpdateTime: time.Now(),
	}

	// 定时更新交易对规格
	initSymbolsSpec()
}

// reloadSymbolsSpec 重新加载交易对规格
func initSymbolsSpec() {
	gateSpotSpec.DeleteSymbolsSpec()
	gateFuturesSpec.DeleteFuturesSpec()

	time.AfterFunc(time.Hour, initSymbolsSpec)
}
