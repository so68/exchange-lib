package gate

import (
	"context"
	"fmt"
	"strings"

	"github.com/so68/exchange-lib/exchange"
)

func (g *gateExchange) CreateFuturesOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateExchange) GetFuturesOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateExchange) CancelFuturesOrder(ctx context.Context, symbol string, orderID string) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateExchange) GetFuturesPositionRisk(ctx context.Context, symbol string) (*exchange.SymbolPositionRisk, error) {
	return nil, nil
}

func (g *gateExchange) CloseFuturesPositionRisk(ctx context.Context, symbol string, positionSide exchange.PositionSide) error {
	return nil
}

func (g *gateExchange) SetFuturesSLTP(ctx context.Context, symbol string, positionSide exchange.PositionSide, stopPrice, takeProfitPrice string) error {
	return nil
}

func (g *gateExchange) CancelFuturesSLTP(ctx context.Context, symbol string) error {
	return nil
}

// GetFuturesSymbolSpec 获取合约交易对规格
func (g *gateExchange) GetFuturesSymbolSpec(ctx context.Context, symbol string) (*futuresSpec, error) {
	var spec *futuresSpec

	// 从缓存中获取交易对规格
	spec, _ = gateFuturesSpec.GetFuturesSpec(symbol)

	if spec == nil {
		contracts, _, err := g.client.FuturesApi.ListFuturesContracts(context.Background(), strings.ToLower(Settle), nil)
		if err != nil {
			return nil, fmt.Errorf("获取合约交易对规则失败: %w", err)
		}

		for _, contract := range contracts {
			specTmp := &futuresSpec{
				Name:              contract.Name,
				Type:              contract.Type,
				QuantoMultiplier:  contract.QuantoMultiplier,
				LeverageMin:       contract.LeverageMin,
				LeverageMax:       contract.LeverageMax,
				MaintenanceRate:   contract.MaintenanceRate,
				MarkType:          contract.MarkType,
				MarkPrice:         contract.MarkPrice,
				IndexPrice:        contract.IndexPrice,
				LastPrice:         contract.LastPrice,
				MakerFeeRate:      contract.MakerFeeRate,
				TakerFeeRate:      contract.TakerFeeRate,
				OrderPriceRound:   contract.OrderPriceRound,
				MarkPriceRound:    contract.MarkPriceRound,
				FundingRate:       contract.FundingRate,
				FundingInterval:   contract.FundingInterval,
				FundingNextApply:  contract.FundingNextApply,
				OrderSizeMin:      contract.OrderSizeMin,
				OrderSizeMax:      contract.OrderSizeMax,
				OrderPriceDeviate: contract.OrderPriceDeviate,
				RefDiscountRate:   contract.RefDiscountRate,
				RefRebateRate:     contract.RefRebateRate,
				OrderbookId:       contract.OrderbookId,
				TradeId:           contract.TradeId,
				TradeSize:         contract.TradeSize,
				PositionSize:      contract.PositionSize,
				ConfigChangeTime:  contract.ConfigChangeTime,
				InDelisting:       contract.InDelisting,
				OrdersLimit:       contract.OrdersLimit,
				EnableBonus:       contract.EnableBonus,
				EnableCredit:      contract.EnableCredit,
				CreateTime:        contract.CreateTime,
				FundingCapRatio:   contract.FundingCapRatio,
				Status:            contract.Status,
				LaunchTime:        contract.LaunchTime,
			}

			// 如果合约匹配，则设置合约规格
			if symbol == contract.Name {
				spec = specTmp
			}
			gateFuturesSpec.SetFuturesSpec(contract.Name, specTmp)
		}
	}

	if spec == nil {
		return nil, fmt.Errorf("合约规格不存在: %s", symbol)
	}

	return spec, nil
}

// GetFuturesSymbols 获取合约交易对列表
func (g *gateExchange) GetFuturesSymbols() []string {
	var symbols []string
	contracts, _, err := g.client.FuturesApi.ListFuturesContracts(context.Background(), strings.ToLower(Settle), nil)
	if err != nil {
		return symbols
	}
	for _, contract := range contracts {
		symbols = append(symbols, contract.Name)
	}
	return symbols
}
