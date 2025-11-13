package exchange

type WebsocketSpotTickerHandler func(ticker *Ticker)

type WebsocketFuturesTickerHandler func(ticker *Ticker)

// Websocket 接口
type Websocket interface {
	// StartListenSpotTickers 开始监听现货交易对行情
	StartListenSpotTickers(handler WebsocketSpotTickerHandler) error
	// StartListenFuturesTickers 开始监听合约交易对行情
	StartListenFuturesTickers(handler WebsocketFuturesTickerHandler) error
}
