package gate

import (
	"sync"
	"time"
)

// exchangeSpec 交易所规格
type exchangeSpec struct {
	Symbols []*symbolSpec `json:"symbols"`
	sync.RWMutex
	UpdateTime time.Time // 更新时间
}

// symbolSpec 交易对规格
type symbolSpec struct {
	Id              string // 交易对
	Base            string // 交易货币
	BaseName        string // 交易货币名称
	Quote           string // 计价货币
	QuoteName       string // 计价货币名称
	MinBaseAmount   string // 交易货币最低交易数量
	MinQuoteAmount  string // 计价货币最低交易数量
	MaxBaseAmount   string // 交易货币最大交易数量
	MaxQuoteAmount  string // 计价货币最大交易数量
	AmountPrecision int    // 数量精度
	Precision       int    // 价格精度
	TradeStatus     string // 交易状态 untradable: 不可交易 buyable: 可买 sellable: 可卖 tradable: 买卖均可交易
	SellStart       int64  // 允许卖出时间
	BuyStart        int64  // 允许买入时间
	DelistingTime   int64  // 预计下架时间
	TradeUrl        string // 交易链接
	StTag           bool   // 是否在ST风险评估中
}

// SetSymbolSpec 设置交易对规格
func (e *exchangeSpec) SetSymbolSpec(symbol string, spec *symbolSpec) error {
	symbolSpec, index := e.GetSymbolSpec(symbol)
	if symbolSpec == nil {
		e.Symbols = append(e.Symbols, spec)
		return nil
	}

	e.Lock()
	defer e.Unlock()

	// 更新符号规格
	e.Symbols[index] = spec
	return nil
}

// GetSymbolSpec 获取交易对规格
func (e *exchangeSpec) GetSymbolSpec(symbol string) (*symbolSpec, int) {
	e.RLock()
	defer e.RUnlock()
	for i, s := range e.Symbols {
		if s.Id == symbol {
			return s, i
		}
	}
	return nil, -1
}

// DeleteSymbolsSpec 删除所有交易对规格
func (e *exchangeSpec) DeleteSymbolsSpec() {
	e.RLock()
	defer e.RUnlock()

	e.Symbols = make([]*symbolSpec, 0)
	e.UpdateTime = time.Now()
}

type exchangeFuturesSpec struct {
	Contracts []*futuresSpec `json:"contracts"`
	sync.RWMutex
	UpdateTime time.Time // 更新时间
}

// futuresSpec 合约规格
type futuresSpec struct {
	// 合约名称
	Name string `json:"name,omitempty"`
	// 合约类型: inverse - 反向合约, direct - 正向合约
	Type string `json:"type,omitempty"`
	// 转换结算货币的乘数
	QuantoMultiplier string `json:"quanto_multiplier,omitempty"`
	// 最小杠杆
	LeverageMin string `json:"leverage_min,omitempty"`
	// 最大杠杆
	LeverageMax string `json:"leverage_max,omitempty"`
	// 维持保证金率
	MaintenanceRate string `json:"maintenance_rate,omitempty"`
	// 标记价格类型: internal - 内部交易价格, index - 外部指数价格
	MarkType string `json:"mark_type,omitempty"`
	// 当前标记价格
	MarkPrice string `json:"mark_price,omitempty"`
	// 当前指数价格
	IndexPrice string `json:"index_price,omitempty"`
	// 最新交易价格
	LastPrice string `json:"last_price,omitempty"`
	// 做市商手续费率, 负值表示返佣
	MakerFeeRate string `json:"maker_fee_rate,omitempty"`
	// 吃单手续费率
	TakerFeeRate string `json:"taker_fee_rate,omitempty"`
	// 最小订单价格增量
	OrderPriceRound string `json:"order_price_round,omitempty"`
	// 最小标记价格增量
	MarkPriceRound string `json:"mark_price_round,omitempty"`
	// 当前资金费率
	FundingRate string `json:"funding_rate,omitempty"`
	// 资金应用间隔, 单位: 秒
	FundingInterval int32 `json:"funding_interval,omitempty"`
	// 下次资金应用时间
	FundingNextApply float64 `json:"funding_next_apply,omitempty"`
	// 合约允许的最小订单数量
	OrderSizeMin int64 `json:"order_size_min,omitempty"`
	// 合约允许的最大订单数量
	OrderSizeMax int64 `json:"order_size_max,omitempty"`
	// 订单价格与当前标记价格的最大允许偏差. 订单价格 `order_price` 必须满足以下条件:      abs(order_price - mark_price) <= mark_price * order_price_deviate
	OrderPriceDeviate string `json:"order_price_deviate,omitempty"`
	// 推荐用户交易手续费折扣
	RefDiscountRate string `json:"ref_discount_rate,omitempty"`
	// 推荐用户手续费率
	RefRebateRate string `json:"ref_rebate_rate,omitempty"`
	// 订单簿更新ID
	OrderbookId int64 `json:"orderbook_id,omitempty"`
	// 当前交易ID
	TradeId int64 `json:"trade_id,omitempty"`
	// 历史累计交易量
	TradeSize int64 `json:"trade_size,omitempty"`
	// 当前总多头持仓量
	PositionSize int64 `json:"position_size,omitempty"`
	// 最后配置更新时间
	ConfigChangeTime float64 `json:"config_change_time,omitempty"`
	// `in_delisting=true` 且 position_size>0 表示合约处于下架过渡期 `in_delisting=true` 且 position_size=0 表示合约已下架
	InDelisting bool `json:"in_delisting,omitempty"`
	// 最大待处理订单数量
	OrdersLimit int32 `json:"orders_limit,omitempty"`
	// 是否启用奖金
	EnableBonus bool `json:"enable_bonus,omitempty"`
	// 是否启用保证金账户
	EnableCredit bool `json:"enable_credit,omitempty"`
	// 合约创建时间
	CreateTime float64 `json:"create_time,omitempty"`
	// 资金费率最大值的因子. 资金费率最大值 = (1/市场最大杠杆 - 维持保证金率) * funding_cap_ratio
	FundingCapRatio string `json:"funding_cap_ratio,omitempty"`
	// 合约状态类型: prelaunch (预发布), trading (交易中), delisting (下架中), delisted (已下架)
	Status string `json:"status,omitempty"`
	// 合约到期时间
	LaunchTime int64 `json:"launch_time,omitempty"`
}

// SetFuturesSpec 设置合约规格
func (e *exchangeFuturesSpec) SetFuturesSpec(symbol string, spec *futuresSpec) error {
	futuresSpec, index := e.GetFuturesSpec(symbol)
	if futuresSpec == nil {
		e.Contracts = append(e.Contracts, spec)
		return nil
	}

	e.Lock()
	defer e.Unlock()

	// 更新符号规格
	e.Contracts[index] = spec
	return nil
}

// GetFuturesSpec 获取合约规格
func (e *exchangeFuturesSpec) GetFuturesSpec(symbol string) (*futuresSpec, int) {
	e.RLock()
	defer e.RUnlock()
	for i, s := range e.Contracts {
		if s.Name == symbol {
			return s, i
		}
	}
	return nil, -1
}

// DeleteFuturesSpec 删除所有合约规格
func (e *exchangeFuturesSpec) DeleteFuturesSpec() {
	e.RLock()
	defer e.RUnlock()

	e.Contracts = make([]*futuresSpec, 0)
	e.UpdateTime = time.Now()
}
