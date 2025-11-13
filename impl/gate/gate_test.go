package gate

import "flag"

const (
	apiKey    = "6bd608e28a45f98ebbea8e7d03f98924"
	secretKey = "741e2093ef86bf5e4c0032a5d41afdaaa899def585b87db0033bd89a1cb02803"
)

var (
	symbol    = flag.String("symbol", "", "交易对")
	lastPrice = flag.String("lastPrice", "", "最新价格")
	amount    = flag.Float64("amount", 0.0, "数量")
	side      = flag.String("side", "BUY", "方向")
	leverage  = flag.Int("leverage", 1, "杠杆")
	orderID   = flag.Int64("orderID", 0, "订单ID")
)
