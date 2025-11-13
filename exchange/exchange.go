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
	///////////////////////////////// 现货 /////////////////////////////////////////
	// GetSpotSymbolTickers 获取现货交易对行情
	GetSpotSymbolTickers(ctx context.Context, symbol ...string) (*Tickers, error)
	// GetSpotBalance 获取现货余额
	GetSpotBalance(ctx context.Context) ([]Balance, error)
	// CreateSpotOrder 现货下单
	CreateSpotOrder(ctx context.Context, symbol string, side OrderSide, limitPrice, quantity string) (*Order, error)
	// GetSpotOrder 获取现货订单
	GetSpotOrder(ctx context.Context, symbol string, orderID int64) (*Order, error)
	// CancelSpotOrder 现货取消订单
	CancelSpotOrder(ctx context.Context, symbol string, orderID int64) (*Order, error)

	///////////////////////////////// 合约 /////////////////////////////////////////
	// GetFuturesSymbolTickers 获取合约交易对行情
	GetFuturesSymbolTickers(ctx context.Context, symbol ...string) (*Tickers, error)
	// GetFuturesBalance 获取合约余额
	GetFuturesBalance(ctx context.Context) ([]Balance, error)
	// CreateFuturesOrder 合约下单
	CreateFuturesOrder(ctx context.Context, symbol string, side OrderSide, limitPrice, quantity string) (*Order, error)
	// GetFuturesOrder 获取合约订单
	GetFuturesOrder(ctx context.Context, symbol string, orderID int64) (*Order, error)
	// GetFuturesPositionRisk 获取合约持仓风险
	GetFuturesPositionRisk(ctx context.Context, symbol string) (*SymbolPositionRisk, error)
	// CloseFuturesPositionRisk 平仓合约持仓风险
	CloseFuturesPositionRisk(ctx context.Context, symbol string, positionSide PositionSide) error
	// SetFuturesSLTP 设置合约止损止盈
	SetFuturesSLTP(ctx context.Context, symbol string, positionSide PositionSide, stopPrice string, takeProfitPrice string) error
	// CancelFuturesSLTP 撤销合约止损止盈
	CancelFuturesSLTP(ctx context.Context, symbol string) error
	// CancelFuturesOrder 撤销合约订单
	CancelFuturesOrder(ctx context.Context, symbol string, orderID int64) (*Order, error)
}
