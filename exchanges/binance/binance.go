package binance

// BinanceExchange 交易所实现
type BinanceExchange struct {
}

// NewBinanceExchange 创建交易所实例
func NewBinanceExchange() *BinanceExchange {
	return &BinanceExchange{}
}

// Name 返回交易所名称
func (e *BinanceExchange) Name() string {
	return "binance"
}
