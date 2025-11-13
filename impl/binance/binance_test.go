package binance

import (
	"flag"
)

const (
	apiKey    = "uAF7jJkOMygW8VK24OEzrrRNYjDFierW04lutVKtPJT9kLUNfHsByLe7lM7dWi4y"
	secretKey = "jscrjDQsTpsiL8nDExD8Ty52raja6g8f4d0VOXByQfyKrYAmYtDUUCDSxymhmvgu"
)

var (
	symbol    = flag.String("symbol", "", "交易对")
	lastPrice = flag.String("lastPrice", "", "最新价格")
	amount    = flag.Float64("amount", 0.0, "数量")
	side      = flag.String("side", "BUY", "方向")
	leverage  = flag.Int("leverage", 1, "杠杆")
	orderID   = flag.Int64("orderID", 0, "订单ID")
)
