package exchange

type WebsocketTickerHandler func(tickers *Tickers)

// Websocket 接口
type Websocket interface {
	// StartListenSpotTickers 开始监听现货交易对行情
	StartListenSpotTickers(handler WebsocketTickerHandler) error
	// StartListenFuturesTickers 开始监听合约交易对行情
	StartListenFuturesTickers(handler WebsocketTickerHandler) error
}
