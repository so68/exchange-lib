package gate

import (
	"encoding/json"
	"time"

	"github.com/gateio/gateapi-go/v6"
	"github.com/so68/exchange-lib/exchange"
	"github.com/so68/exchange-lib/internal/socket/client"
)

const (
	SpotWebsocketURL    = "wss://api.gateio.ws"
	FuturesWebsocketURL = "wss://fx-ws.gateio.ws"
)

// gateWebsocket Gate Websocket实例
type gateWebsocket struct {
	spotWs    *client.Websocket
	futuresWs *client.Websocket
}

// SubscribeParams 订阅参数
type SubscribeParams struct {
	Time    int64       `json:"time"`
	Channel string      `json:"channel"`
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

// SubscribeResult 订阅结果
type SubscribeResult struct {
	Time    int64           `json:"time"`
	Channel string          `json:"channel"`
	Event   string          `json:"event"`
	Result  json.RawMessage `json:"result"`
}

// NewGateWebsocket 创建Gate Websocket实例
func NewGateWebsocket() exchange.Websocket {
	return &gateWebsocket{}
}

// StartListenSpotTickers 开始监听现货交易对行情
func (g *gateWebsocket) StartListenSpotTickers(handler exchange.WebsocketSpotTickerHandler) error {
	g.spotWs = client.NewWebsocket(SpotWebsocketURL+"/ws/v4/", func(message []byte) {
		resp := &SubscribeResult{}
		err := json.Unmarshal(message, resp)
		if err != nil {
			return
		}

		// 解析对应的结果
		switch resp.Channel {
		case "spot.tickers":
			spotTicker := &gateapi.Ticker{}
			err := json.Unmarshal(resp.Result, spotTicker)
			if err != nil || spotTicker.CurrencyPair == "" {
				return
			}

			// 计算开盘价和价格变动
			openPrice, priceChange := calculateOpenAndChangePrice(spotTicker.Last, spotTicker.ChangePercentage)
			handler(&exchange.Ticker{
				Symbol:             spotTicker.CurrencyPair,
				PriceChange:        priceChange,
				PriceChangePercent: spotTicker.ChangePercentage,
				WeightedAvgPrice:   "", // Gate API 不提供加权平均价
				LastPrice:          spotTicker.Last,
				LastQty:            "", // Gate API 不提供最新成交量
				OpenPrice:          openPrice,
				HighPrice:          spotTicker.High24h,
				LowPrice:           spotTicker.Low24h,
				Volume:             spotTicker.BaseVolume,
				QuoteVolume:        spotTicker.QuoteVolume,
				Count:              0, // Gate API 不提供成交笔数
			})
		}

	})
	// 设置连接成功后的回调处理器
	g.spotWs.SetAfterConnectionHandler(func() error {
		gateExchange := newGateExchange("", "")
		symbols := gateExchange.GetSpotSymbols()
		subscribeParams := SubscribeParams{
			Time:    time.Now().Unix(),
			Channel: "spot.tickers",
			Event:   "subscribe",
			Payload: symbols,
		}
		subscribeBytes, err := json.Marshal(subscribeParams)
		if err != nil {
			return err
		}
		return g.spotWs.WriteMessage(subscribeBytes)
	})
	return g.spotWs.Start()
}

// StartListenFuturesTickers 开始监听合约交易对行情
func (g *gateWebsocket) StartListenFuturesTickers(handler exchange.WebsocketFuturesTickerHandler) error {
	g.futuresWs = client.NewWebsocket(FuturesWebsocketURL+"/v4/ws/usdt", func(message []byte) {
		resp := &SubscribeResult{}
		err := json.Unmarshal(message, resp)
		if err != nil {
			return
		}

		// 解析对应的结果
		switch resp.Channel {
		case "futures.tickers":
			futuresTickers := make([]*gateapi.FuturesTicker, 0)
			err := json.Unmarshal(resp.Result, &futuresTickers)
			if err != nil || len(futuresTickers) == 0 {
				return
			}

			futuresTicker := futuresTickers[0]
			// 优先使用 Volume24hBase 作为 Volume，如果没有则使用 Volume24h
			volume := futuresTicker.Volume24hBase
			if volume == "" {
				volume = futuresTicker.Volume24h
			}
			// 优先使用 Volume24hQuote 作为 QuoteVolume，如果没有则使用 Volume24hSettle
			quoteVolume := futuresTicker.Volume24hQuote
			if quoteVolume == "" {
				quoteVolume = futuresTicker.Volume24hSettle
			}
			openPrice, priceChange := calculateOpenAndChangePrice(futuresTicker.Last, futuresTicker.ChangePercentage)
			handler(&exchange.Ticker{
				Symbol:             futuresTicker.Contract,
				PriceChange:        priceChange,
				PriceChangePercent: futuresTicker.ChangePercentage,
				WeightedAvgPrice:   "", // Gate API 不提供加权平均价
				LastPrice:          futuresTicker.Last,
				LastQty:            "", // Gate API 不提供最新成交量
				OpenPrice:          openPrice,
				HighPrice:          futuresTicker.High24h,
				LowPrice:           futuresTicker.Low24h,
				Volume:             volume,
				QuoteVolume:        quoteVolume,
				Count:              0, // Gate API 不提供成交笔数
			})
		}
	})

	g.futuresWs.SetAfterConnectionHandler(func() error {
		gateExchange := newGateExchange("", "")
		symbols := gateExchange.GetFuturesSymbols()
		subscribeParams := SubscribeParams{
			Time:    time.Now().Unix(),
			Channel: "futures.tickers",
			Event:   "subscribe",
			Payload: symbols,
		}
		subscribeBytes, err := json.Marshal(subscribeParams)
		if err != nil {
			return err
		}
		return g.futuresWs.WriteMessage(subscribeBytes)
	})
	return g.futuresWs.Start()
}
