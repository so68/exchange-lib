package gate

import (
	"context"
	"fmt"

	"github.com/so68/exchange-lib/exchange"
)

func (g *gateExchange) CreateSpotOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateExchange) GetSpotOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateExchange) CancelSpotOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
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
			if symbol == pair.Base+pair.Quote {
				spec = specTmp
			}
			gateSpotSpec.SetSymbolSpec(pair.Base+pair.Quote, specTmp)
		}
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
