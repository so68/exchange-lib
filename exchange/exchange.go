package exchange

import "context"

// 余额
type Balance struct {
	Symbol string `json:"symbol"` // 币种符号
	Free   string `json:"free"`   // 可用余额
	Locked string `json:"locked"` // 锁定余额
	Total  string `json:"total"`  // 总余额
}

type Exchange interface {
	// GetSymbolTickers 获取交易对行情
	GetSymbolTickers(ctx context.Context, symbol ...string) (*Tickers, error)
	// SpotBalance 获取现货余额
	SpotBalance(ctx context.Context) ([]Balance, error)
	// FuturesBalance 获取合约余额
	FuturesBalance(ctx context.Context) ([]Balance, error)
	// SpotCreateOrder 现货下单
	SpotCreateOrder(ctx context.Context, symbol string, side OrderSide, limitPrice, quantity string) (*Order, error)
	// SpotGetOrder 获取现货订单
	SpotGetOrder(ctx context.Context, symbol string, orderID int64) (*Order, error)
	// SpotCancelOrder 现货取消订单
	SpotCancelOrder(ctx context.Context, symbol string, orderID int64) (*Order, error)

	// // 合约
	// FuturesBalance(ctx context.Context) ([]Balance, error)
	// FuturesCreateOrder(ctx context.Context, symbol, side, orderType, quantity, price string) (*Order, error)
	// FuturesGetOrder(ctx context.Context, symbol string, orderID int64) (*Order, error)
}
