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

// SymbolPositionRisk 交易对持仓风险
type SymbolPositionRisk struct {
	Data []*PositionRisk `json:"data"`
}

// GetSidePositionRisk 获取指定方向持仓风险
func (s *SymbolPositionRisk) GetSidePositionRisk(side PositionSide) *PositionRisk {
	for _, risk := range s.Data {
		if risk.PositionSide == side {
			return risk
		}
	}
	return nil
}

// PositionRisk 持仓风险
type PositionRisk struct {
	Symbol           string       `json:"symbol"`            // 交易对符号（如 "ETHUSDT"）
	PositionSide     PositionSide `json:"position_side"`     // 持仓方向："LONG"（多头）、"SHORT"（空头）、"BOTH"（单向模式，币安/欧易/芝麻支持）
	PositionAmt      string       `json:"position_amt"`      // 持仓数量（多头为正，空头为负；币安：PositionAmt，欧易：Pos，芝麻：Size）
	EntryPrice       string       `json:"entry_price"`       // 平均开仓价格（币安：EntryPrice，欧易：AvgPx，芝麻：EntryPrice）
	MarkPrice        string       `json:"mark_price"`        // 标记价格（用于公平估值、止损爆仓；币安/欧易/芝麻：MarkPx/MarkPrice）
	UnRealizedProfit string       `json:"unrealized_profit"` // 未实现盈亏（浮动盈亏，单位 USDT；币安：UnRealizedProfit，欧易：Upl，芝麻：UnrealisedPnl）
	Leverage         string       `json:"leverage"`          // 杠杆倍数（如 "20"；币安/欧易/芝麻：Lever/Leverage）
	LiquidationPrice string       `json:"liquidation_price"` // 爆仓价格（标记价格触及即强平；币安：LiquidationPrice，欧易：LiqPx，芝麻：LiqPrice）
	MarginType       string       `json:"margin_type"`       // 保证金模式："cross"（全仓）或 "isolated"（逐仓）；币安：MarginType，欧易：MgnMode，芝麻：MarginMode
	IsolatedMargin   string       `json:"isolated_margin"`   // 逐仓保证金金额（USDT；币安：IsolatedMargin，欧易：Margin，芝麻：Margin）
	Notional         string       `json:"notional"`          // 名义价值（持仓总价值，单位 USDT；币安：Notional，欧易：NotionalUsd，芝麻：Value）
}
