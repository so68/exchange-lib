package gate

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/gateio/gateapi-go/v6"
	"github.com/so68/exchange-lib/exchange"
)

// CreateSpotOrder 创建现货订单
func (g *gateExchange) CreateSpotOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	spec, err := g.GetSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// 验证交易规则
	quantity, err = g.filtersQuantity(spec, limitPrice, quantity)
	if err != nil {
		return nil, fmt.Errorf("验证交易规则失败: %w", err)
	}

	order := gateapi.Order{
		CurrencyPair: symbol,
		Side:         strings.ToLower(string(side)), // buy 或 sell
		Amount:       quantity,                      // 数量（基础币种）
		Price:        limitPrice,                    // 价格（报价币种）
		Type:         "limit",                       // limit（限价）或 market（市价）
	}

	createdOrder, _, err := g.client.SpotApi.CreateOrder(ctx, order, nil)
	if err != nil {
		return nil, fmt.Errorf("下单失败: %w", err)
	}

	// 计算手续费
	filledAmount := new(big.Float).SetPrec(64)
	if _, ok := filledAmount.SetString(createdOrder.FilledAmount); !ok {
		return nil, fmt.Errorf("无效的已成交数量: %s", createdOrder.FilledAmount)
	}

	feeAmount := new(big.Float).SetPrec(64)
	if _, ok := feeAmount.SetString(createdOrder.Fee); !ok {
		return nil, fmt.Errorf("无效的手续费: %s", createdOrder.Fee)
	}
	actualQty := filledAmount.Sub(filledAmount, feeAmount)

	return &exchange.Order{
		OrderID:       createdOrder.Id,
		Symbol:        createdOrder.CurrencyPair,
		Side:          exchange.OrderSide(strings.ToUpper(createdOrder.Side)),
		Type:          exchange.OrderType(strings.ToUpper(createdOrder.Type)),
		Status:        exchange.OrderStatusNew, // 默认新订单
		Price:         createdOrder.Price,
		Quantity:      createdOrder.Amount,
		ExecutedQty:   createdOrder.FilledAmount,
		ActualQty:     actualQty.Text('f', spec.AmountPrecision),
		QuoteQuantity: createdOrder.FilledTotal,
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(createdOrder.TimeInForce)),
		CreateTime:    createdOrder.CreateTimeMs,
		UpdateTime:    createdOrder.UpdateTimeMs,
	}, nil
}

// GetSpotOrder 获取现货订单
func (g *gateExchange) GetSpotOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	// 获取交易对规格
	spec, err := g.GetSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, err
	}
	singleOrder, _, err := g.client.SpotApi.GetOrder(ctx, orderID, symbol, nil)
	if err != nil {
		return nil, fmt.Errorf("获取单个订单失败: %w", err)
	}

	// 计算手续费
	filledAmount := new(big.Float).SetPrec(64)
	if _, ok := filledAmount.SetString(singleOrder.FilledAmount); !ok {
		return nil, fmt.Errorf("无效的已成交数量: %s", singleOrder.FilledAmount)
	}

	feeAmount := new(big.Float).SetPrec(64)
	if _, ok := feeAmount.SetString(singleOrder.Fee); !ok {
		return nil, fmt.Errorf("无效的手续费: %s", singleOrder.Fee)
	}
	actualQty := filledAmount.Sub(filledAmount, feeAmount)

	status := exchange.OrderStatusNew
	switch singleOrder.FinishAs {
	case "filled":
		status = exchange.OrderStatusFilled
	case "cancelled":
		status = exchange.OrderStatusCanceled
	case "small":
		status = exchange.OrderStatusRejected
	case "depth_not_enough":
		status = exchange.OrderStatusRejected
	case "trader_not_enough":
		status = exchange.OrderStatusRejected
	}

	return &exchange.Order{
		OrderID:       singleOrder.Id,
		Symbol:        singleOrder.CurrencyPair,
		Side:          exchange.OrderSide(strings.ToUpper(singleOrder.Side)),
		Type:          exchange.OrderType(strings.ToUpper(singleOrder.Type)),
		Status:        status,
		Price:         singleOrder.Price,
		Quantity:      singleOrder.Amount,
		ExecutedQty:   singleOrder.FilledAmount,
		ActualQty:     actualQty.Text('f', spec.AmountPrecision),
		QuoteQuantity: singleOrder.FilledTotal,
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(singleOrder.TimeInForce)),
		CreateTime:    singleOrder.CreateTimeMs,
		UpdateTime:    singleOrder.UpdateTimeMs,
	}, nil
}

// CancelSpotOrder 撤销现货订单
func (g *gateExchange) CancelSpotOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	// 获取交易对规格
	spec, err := g.GetSpotSymbolSpec(ctx, symbol)
	if err != nil {
		return nil, err
	}

	canceledOrder, _, err := g.client.SpotApi.CancelOrder(ctx, orderID, symbol, nil)
	if err != nil {
		return nil, fmt.Errorf("取消单个订单失败: %w", err)
	}
	// 计算手续费
	filledAmount := new(big.Float).SetPrec(64)
	if _, ok := filledAmount.SetString(canceledOrder.FilledAmount); !ok {
		return nil, fmt.Errorf("无效的已成交数量: %s", canceledOrder.FilledAmount)
	}

	feeAmount := new(big.Float).SetPrec(64)
	if _, ok := feeAmount.SetString(canceledOrder.Fee); !ok {
		return nil, fmt.Errorf("无效的手续费: %s", canceledOrder.Fee)
	}
	actualQty := filledAmount.Sub(filledAmount, feeAmount)

	status := exchange.OrderStatusNew
	switch canceledOrder.FinishAs {
	case "filled":
		status = exchange.OrderStatusFilled
	case "cancelled":
		status = exchange.OrderStatusCanceled
	case "small":
		status = exchange.OrderStatusRejected
	case "depth_not_enough":
		status = exchange.OrderStatusRejected
	case "trader_not_enough":
		status = exchange.OrderStatusRejected
	}

	return &exchange.Order{
		OrderID:       canceledOrder.Id,
		Symbol:        canceledOrder.CurrencyPair,
		Side:          exchange.OrderSide(strings.ToUpper(canceledOrder.Side)),
		Type:          exchange.OrderType(strings.ToUpper(canceledOrder.Type)),
		Status:        status,
		Price:         canceledOrder.Price,
		Quantity:      canceledOrder.Amount,
		ExecutedQty:   canceledOrder.FilledAmount,
		ActualQty:     actualQty.Text('f', spec.AmountPrecision),
		QuoteQuantity: canceledOrder.FilledTotal,
		TimeInForce:   exchange.OrderTimeInForce(strings.ToUpper(canceledOrder.TimeInForce)),
		CreateTime:    canceledOrder.CreateTimeMs,
		UpdateTime:    canceledOrder.UpdateTimeMs,
	}, nil
}

// GetSpotSymbolSpec 获取现货交易对规格
func (g *gateExchange) GetSpotSymbolSpec(ctx context.Context, symbol string) (*symbolSpec, error) {
	var spec *symbolSpec

	// 从缓存中获取交易对规格
	spec, _ = gateSpotSpec.GetSymbolSpec(symbol)

	if spec == nil {
		// 获取所有现货交易对规则
		pairs, _, err := g.client.SpotApi.ListCurrencyPairs(ctx)
		if err != nil {
			return nil, fmt.Errorf("获取现货交易对规则失败: %w", err)
		}

		for _, pair := range pairs {
			// 只获取结算货币为 USDT 的交易对
			if pair.Quote != Settle {
				continue
			}
			specTmp := &symbolSpec{
				Id:              pair.Id,
				Base:            pair.Base,
				Quote:           pair.Quote,
				MinBaseAmount:   pair.MinBaseAmount,
				MinQuoteAmount:  pair.MinQuoteAmount,
				MaxBaseAmount:   pair.MaxBaseAmount,
				MaxQuoteAmount:  pair.MaxQuoteAmount,
				AmountPrecision: int(pair.AmountPrecision),
				Precision:       int(pair.Precision),
				TradeStatus:     pair.TradeStatus,
				SellStart:       pair.SellStart,
				BuyStart:        pair.BuyStart,
				DelistingTime:   pair.DelistingTime,
				TradeUrl:        pair.TradeUrl,
				StTag:           pair.StTag,
			}

			// 如果交易对匹配，则设置交易对规格
			if symbol == pair.Id {
				spec = specTmp
			}
			gateSpotSpec.SetSymbolSpec(pair.Id, specTmp)
		}
	}

	if spec == nil {
		return nil, fmt.Errorf("交易对规格不存在: %s", symbol)
	}
	return spec, nil
}

// GetSpotSymbols 获取现货交易对列表
func (g *gateExchange) GetSpotSymbols() []string {
	var symbols []string
	pairs, _, err := g.client.SpotApi.ListCurrencyPairs(context.Background())
	if err != nil {
		return symbols
	}
	for _, pair := range pairs {
		// 只获取结算货币为 USDT 的交易对
		if pair.Quote != Settle {
			continue
		}
		symbols = append(symbols, pair.Id)
	}
	return symbols
}

// Spot order details
type Order struct {
	// 订单ID
	Id string `json:"id,omitempty"`
	// 用户定义的信息. 如果非空，必须遵循以下规则：  1. 以 `t-` 开头 2. 不超过 28 个字节，不包含 `t-` 前缀 3. 只能包含 0-9, A-Z, a-z, 下划线(_), 连字符(-) 或 点(.) 除了用户定义的信息，保留内容列表如下，表示订单的创建方式：  - 101: 来自 android - 102: 来自 IOS - 103: 来自 IPAD - 104: 来自 webapp - 3: 来自 web - 2: 来自 apiv2 - apiv4: 来自 apiv4
	Text string `json:"text,omitempty"`
	// 用户在修改订单时备注的自定义数据
	AmendText string `json:"amend_text,omitempty"`
	// 订单创建时间
	CreateTime string `json:"create_time,omitempty"`
	// 订单最后修改时间
	UpdateTime string `json:"update_time,omitempty"`
	// 订单创建时间 (毫秒)
	CreateTimeMs int64 `json:"create_time_ms,omitempty"`
	// 订单最后修改时间 (毫秒)
	UpdateTimeMs int64 `json:"update_time_ms,omitempty"`
	// 订单状态  - `open`: 待填充 - `closed`: 已填充 - `cancelled`: 已取消
	Status string `json:"status,omitempty"`
	// 交易对
	CurrencyPair string `json:"currency_pair"`
	// 订单类型   - limit : Limit Order - market : Market Order
	Type string `json:"type,omitempty"`
	// 账户类型, spot - spot account, margin - leveraged account, unified - unified account
	Account string `json:"account,omitempty"`
	// 买入或卖出订单
	Side string `json:"side"`
	// 交易数量 When `type` is `limit`, it refers to the base currency (the currency being traded), such as `BTC` in `BTC_USDT` When `type` is `market`, it refers to different currencies based on the side: - `side`: `buy` refers to quote currency, `BTC_USDT` means `USDT` - `side`: `sell` refers to base currency, `BTC_USDT` means `BTC`
	Amount string `json:"amount"`
	// 交易价格, required when `type`=`limit`
	Price string `json:"price,omitempty"`
	// 时间类型  - gtc: GoodTillCancelled - ioc: ImmediateOrCancelled, taker only - poc: PendingOrCancelled, makes a post-only order that always enjoys a maker fee - fok: FillOrKill, fill either completely or none Only `ioc` and `fok` are supported when `type`=`market`
	TimeInForce string `json:"time_in_force,omitempty"`
	// 冰山订单数量. Null or 0 for normal orders. Hiding all amount is not supported
	Iceberg string `json:"iceberg,omitempty"`
	// 在保证金或跨保证金交易中允许自动借入不足的金额如果余额不足
	AutoBorrow bool `json:"auto_borrow,omitempty"`
	// 启用或禁用自动还款用于自动借入订单生成的保证金订单。默认禁用。注意：  1. 此字段仅适用于跨保证金订单。保证金账户不支持设置自动还款用于订单。 2. `auto_borrow` 和 `auto_repay` 可以同时设置为 true 在一笔订单中
	AutoRepay bool `json:"auto_repay,omitempty"`
	// 剩余数量
	Left string `json:"left,omitempty"`
	// 已成交数量
	FilledAmount string `json:"filled_amount,omitempty"`
	// 成交金额. Deprecated in favor of `filled_total`
	FillPrice string `json:"fill_price,omitempty"`
	// 成交金额. Deprecated in favor of `filled_total`
	FilledTotal string `json:"filled_total,omitempty"`
	// 平均成交价格
	AvgDealPrice string `json:"avg_deal_price,omitempty"`
	// 手续费
	Fee string `json:"fee,omitempty"`
	// 手续费币种
	FeeCurrency string `json:"fee_currency,omitempty"`
	// 点数手续费
	PointFee string `json:"point_fee,omitempty"`
	// GT手续费
	GtFee string `json:"gt_fee,omitempty"`
	// GT maker手续费
	GtMakerFee string `json:"gt_maker_fee,omitempty"`
	// GT taker手续费
	GtTakerFee string `json:"gt_taker_fee,omitempty"`
	// GT 手续费是否启用
	GtDiscount bool `json:"gt_discount,omitempty"`
	// 返现手续费
	RebatedFee string `json:"rebated_fee,omitempty"`
	// 返现手续费币种
	RebatedFeeCurrency string `json:"rebated_fee_currency,omitempty"`
	// 用户在同一 `stp_id` 组内的订单不允许自交易  1. 如果两个订单的 `stp_id` 非零且相等，则不会执行。相反，将根据取单的 `stp_act` 执行相应的策略。 2. `stp_id` 默认返回 `0` 用于未设置 `STP 组`的订单
	StpId int32 `json:"stp_id,omitempty"`
	// 自交易预防行动. 用户可以使用此字段设置自交易预防策略  1. 用户加入 `STP 组`后，可以传递 `stp_act` 限制用户自交易预防策略。如果未传递 `stp_act`，则默认使用 `cn` 策略。 2. 当用户未加入 `STP 组`时，传递 `stp_act` 参数时会返回错误。 3. 如果用户在下单时未使用 `stp_act`，则 `stp_act` 将返回 '-'  - cn: 取消最新订单，取消新订单并保留旧订单 - co: 取消最旧订单，取消旧订单并保留新订单 - cb: 取消两者，两者旧订单和新订单都会被取消
	StpAct string `json:"stp_act,omitempty"`
	// 订单完成状态包括：  - open: 待处理 - filled: 完全成交 - cancelled: 用户取消 - liquidate_cancelled: 因清算取消 - small: 订单数量太小 - depth_not_enough: 因市场深度不足取消 - trader_not_enough: 因对手方不足取消 - ioc: 因 tif 设置为 ioc 未立即成交 - poc: 因 poc 未立即成交 - fok: 因 tif 设置为 fok 未立即成交 - stp: 因自交易预防取消 - unknown: 未知
	FinishAs string `json:"finish_as,omitempty"`
	// 处理模式: 在下单时，根据 action_mode 返回不同的字段。此字段仅在请求期间有效，不会包含在响应结果中 ACK: 异步模式，仅返回关键订单字段 RESULT: 无清算信息 FULL: 全模式 (默认)
	ActionMode string `json:"action_mode,omitempty"`
}
