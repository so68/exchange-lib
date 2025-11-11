package binance

import "flag"

var (
	symbol    = flag.String("symbol", "ETHUSDT", "交易对")
	lastPrice = flag.String("lastPrice", "", "最新价格")
	amount    = flag.Float64("amount", 0.0, "数量")
	side      = flag.String("side", "BUY", "方向")

	orderID = flag.Int64("orderID", 0, "订单ID")
)
