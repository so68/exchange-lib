package okx

import "flag"

const (
	apiKey     = "c035a2cd-7c2d-4d5f-a91b-ee51a0807bd5"
	secretKey  = "50BF980FA9F02F16DF30A0CDDFBA2ADE"
	passphrase = "Aa123098.."
)

var (
	symbol    = flag.String("symbol", "", "交易对")
	lastPrice = flag.String("lastPrice", "", "最新价格")
	amount    = flag.Float64("amount", 0.0, "数量")
	side      = flag.String("side", "BUY", "方向")
	leverage  = flag.Int("leverage", 1, "杠杆")
	orderID   = flag.String("orderID", "", "订单ID")
)
