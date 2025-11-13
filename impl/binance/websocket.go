package binance

import (
	"encoding/json"

	"github.com/so68/exchange-lib/exchange"
	"github.com/so68/exchange-lib/internal/socket/client"
)

const (
	SpotWebsocketURL    = "wss://stream.binance.com:9443/ws"
	FuturesWebsocketURL = "wss://fstream.binance.com/ws/!ticker@arr"
)

// binanceWebsocket Binance Websocket实例
type binanceWebsocket struct {
	spotWs    *client.Websocket
	futuresWs *client.Websocket
}

// SubscribeParams 订阅参数
type SubscribeParams struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

// NewBinanceWebsocket 创建Binance Websocket实例
func NewBinanceWebsocket() exchange.Websocket {
	return &binanceWebsocket{}
}

// StartListenSpotTickers 开始监听现货交易对行情
func (b *binanceWebsocket) StartListenSpotTickers(handler exchange.WebsocketSpotTickerHandler) error {
	b.spotWs = client.NewWebsocket(SpotWebsocketURL+"/!ticker@arr", func(message []byte) {
		var event []*WsAllTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			return
		}

		for _, event := range event {
			spotTicker := &exchange.Ticker{
				Symbol:             event.Symbol,
				PriceChange:        event.PriceChange,
				PriceChangePercent: event.PriceChangePercent,
				WeightedAvgPrice:   event.WeightedAvgPrice,
				LastPrice:          event.LastPrice,
				LastQty:            event.LastQty,
				OpenPrice:          event.OpenPrice,
				HighPrice:          event.HighPrice,
				LowPrice:           event.LowPrice,
				Volume:             event.TotalVolume,
				QuoteVolume:        event.TotalQuoteVolume,
				Count:              event.TradeCount,
			}
			handler(spotTicker)
		}
	})
	return b.spotWs.Start()
}

// StartListenFuturesTickers 开始监听合约交易对行情
func (b *binanceWebsocket) StartListenFuturesTickers(handler exchange.WebsocketFuturesTickerHandler) error {
	b.futuresWs = client.NewWebsocket(FuturesWebsocketURL, func(message []byte) {
		var event []*WsAllTickerEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			return
		}

		for _, event := range event {
			futuresTicker := &exchange.Ticker{
				Symbol:             event.Symbol,
				PriceChange:        event.PriceChange,
				PriceChangePercent: event.PriceChangePercent,
				WeightedAvgPrice:   event.WeightedAvgPrice,
				LastPrice:          event.LastPrice,
				LastQty:            event.LastQty,
				OpenPrice:          event.OpenPrice,
				HighPrice:          event.HighPrice,
				LowPrice:           event.LowPrice,
				Volume:             event.TotalVolume,
				QuoteVolume:        event.TotalQuoteVolume,
				Count:              event.TradeCount,
			}
			handler(futuresTicker)
		}
	})
	return b.futuresWs.Start()
}
