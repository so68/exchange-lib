package binance

import (
	"sync"
	"time"
)

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

// exchangeSpec 交易所规格
type exchangeSpec struct {
	Symbols []*symbolSpec `json:"symbols"`
	sync.RWMutex
	UpdateTime time.Time // 更新时间
}

// SymbolSpec 交易对规格
type symbolSpec struct {
	Symbol         string // 交易对
	BaseAsset      string // 基础资产
	QuoteAsset     string // 计价资产
	BasePrecision  int    // 基础资产精度
	QuotePrecision int    // 计价资产精度
	MinQty         string // 最小数量
	MaxQty         string // 最大数量
	StepSize       string // 步长
	MinPrice       string // 最小价格
	MaxPrice       string // 最大价格
	TickSize       string // 最小价格变动
	MinNotional    string // 最小交易金额
	Status         string // 状态
}

// SetSymbolSpec 设置交易对规格
func (e *exchangeSpec) SetSymbolSpec(symbol string, spec *symbolSpec) error {
	symbolSpec, index := e.GetSymbolSpec(symbol)
	if symbolSpec == nil {
		e.Symbols = append(e.Symbols, spec)
		return nil
	}

	e.Lock()
	defer e.Unlock()

	// 更新符号规格
	e.Symbols[index] = spec
	return nil
}

// GetSymbolSpec 获取交易对规格
func (e *exchangeSpec) GetSymbolSpec(symbol string) (*symbolSpec, int) {
	e.RLock()
	defer e.RUnlock()
	for i, s := range e.Symbols {
		if s.Symbol == symbol {
			return s, i
		}
	}
	return nil, -1
}

// DeleteSymbolsSpec 删除所有交易对规格
func (e *exchangeSpec) DeleteSymbolsSpec() {
	e.RLock()
	defer e.RUnlock()

	e.Symbols = make([]*symbolSpec, 0)
	e.UpdateTime = time.Now()
}
