package exchange

// Tickers 行情列表
type Tickers struct {
	Tickers []*Ticker `json:"tickers"`
}

// GetTicker 获取指定交易对的行情
func (t *Tickers) GetTicker(symbol string) *Ticker {
	for _, ticker := range t.Tickers {
		if ticker.Symbol == symbol {
			return ticker
		}
	}
	return nil
}

// Ticker 行情
type Ticker struct {
	Symbol      string `json:"symbol"`      // 交易对
	OpenPrice   string `json:"openPrice"`   // 开盘价
	HighPrice   string `json:"highPrice"`   // 最高价
	LowPrice    string `json:"lowPrice"`    // 最低价
	LastPrice   string `json:"lastPrice"`   // 最新价
	Volume      string `json:"volume"`      // 成交量
	QuoteVolume string `json:"quoteVolume"` // 成交额，单位：USDT
	Count       int64  `json:"count"`       // 成交笔数
}
