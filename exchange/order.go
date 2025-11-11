package exchange

// 订单
type Order struct {
	OrderID       int64            `json:"orderId"`       // 订单ID
	Symbol        string           `json:"symbol"`        // 交易对
	Side          OrderSide        `json:"side"`          // 方向
	Type          OrderType        `json:"type"`          // 类型
	Status        OrderStatus      `json:"status"`        // 状态
	Price         string           `json:"price"`         // 价格
	Quantity      string           `json:"quantity"`      // 数量
	ExecutedQty   string           `json:"executedQty"`   // 已执行数量
	ActualQty     string           `json:"actualQty"`     // 实际数量（扣除手续费后的数量）
	QuoteQuantity string           `json:"quoteQuantity"` // 成交金额
	TimeInForce   OrderTimeInForce `json:"timeInForce"`   // 时间类型
	CreateTime    int64            `json:"createTime"`    // 创建时间
	UpdateTime    int64            `json:"updateTime"`    // 更新时间
}

type Fill struct {
	Price           string `json:"price"`           // 价格
	Quantity        string `json:"quantity"`        // 数量
	Commission      string `json:"commission"`      // 手续费
	CommissionAsset string `json:"commissionAsset"` // 手续费币种
}
