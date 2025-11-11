package exchange

import "context"

// 上下文键【测试网】
type ctxKey string

// 订单时间类型
type OrderTimeInForce string

// 订单类型
type OrderType string

// 订单方向
type OrderSide string

// 订单状态
type OrderStatus string

const (
	CtxKeyTestnet ctxKey = "testnet" // 测试网

	OrderTimeInForceGTC OrderTimeInForce = "GTC" // 一直有效，直到手动取消或完全成交
	OrderTimeInForceIOC OrderTimeInForce = "IOC" // 立即成交，否则取消
	OrderTimeInForceFOK OrderTimeInForce = "FOK" // 全部立即成交，否则整单取消
	OrderTimeInForceGTX OrderTimeInForce = "GTX" // 只做挂单，不立即成交

	OrderTypeLimit  OrderType = "LIMIT"  // 限价单
	OrderTypeMarket OrderType = "MARKET" // 市价单

	OrderSideBuy  OrderSide = "BUY"  // 买入
	OrderSideSell OrderSide = "SELL" // 卖出

	OrderStatusNew             OrderStatus = "NEW"              // 新订单
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED" // 部分成交
	OrderStatusFilled          OrderStatus = "FILLED"           // 完全成交
	OrderStatusCanceled        OrderStatus = "CANCELED"         // 已取消
	OrderStatusPendingCancel   OrderStatus = "PENDING_CANCEL"   // 待取消
	OrderStatusRejected        OrderStatus = "REJECTED"         // 已拒绝
	OrderStatusExpired         OrderStatus = "EXPIRED"          // 已过期
)

func WithTestnet(parent context.Context) context.Context {
	return context.WithValue(parent, CtxKeyTestnet, true)
}
