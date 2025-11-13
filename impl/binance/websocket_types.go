package binance

// 24hr Ticker 事件结构（Spot）
type WsAllTickerEvent struct {
	EventType          string `json:"e"` // "24hrTicker"
	EventTime          int64  `json:"E"` // 事件时间
	Symbol             string `json:"s"` // 交易对
	PriceChange        string `json:"p"` // 价格变动
	PriceChangePercent string `json:"P"` // 24h 涨跌幅
	WeightedAvgPrice   string `json:"w"` // 加权平均价
	LastPrice          string `json:"c"` // 最新价
	LastQty            string `json:"Q"` // 最新成交量
	OpenPrice          string `json:"o"` // 开盘价
	HighPrice          string `json:"h"` // 最高价
	LowPrice           string `json:"l"` // 最低价
	TotalVolume        string `json:"v"` // 总成交量
	TotalQuoteVolume   string `json:"q"` // 总成交额
	OpenTime           int64  `json:"O"` // 开盘时间
	CloseTime          int64  `json:"C"` // 收盘时间
	FirstTradeID       int64  `json:"F"` // 第一笔成交ID
	LastTradeID        int64  `json:"L"` // 最后一笔成交ID
	TradeCount         int64  `json:"n"` // 成交笔数
}
