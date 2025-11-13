package binance

import (
	"context"
	"fmt"

	"github.com/so68/exchange-lib/exchange"
)

// GetSymbolTickers 获取交易对行情
func (b *binanceExchange) GetSpotSymbolTickers(ctx context.Context, symbols ...string) (*exchange.Tickers, error) {
	resp, err := b.client.NewListPriceChangeStatsService().Symbols(symbols).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance get ticker: %w", err)
	}
	var res []*exchange.Ticker
	for _, ticker := range resp {
		res = append(res, &exchange.Ticker{
			Symbol:             ticker.Symbol,
			PriceChange:        ticker.PriceChange,
			PriceChangePercent: ticker.PriceChangePercent,
			WeightedAvgPrice:   ticker.WeightedAvgPrice,
			LastPrice:          ticker.LastPrice,
			LastQty:            ticker.LastQty,
			OpenPrice:          ticker.OpenPrice,
			HighPrice:          ticker.HighPrice,
			LowPrice:           ticker.LowPrice,
			Volume:             ticker.Volume,
			QuoteVolume:        ticker.QuoteVolume,
			Count:              ticker.Count,
		})
	}
	return &exchange.Tickers{
		Tickers: res,
	}, nil
}

// GetFuturesSymbolTickers 获取合约交易对行情
func (b *binanceExchange) GetFuturesSymbolTickers(ctx context.Context, symbols ...string) (*exchange.Tickers, error) {
	var res []*exchange.Ticker
	for _, symbol := range symbols {
		resp, err := b.futuresClient.NewListPriceChangeStatsService().Symbol(symbol).Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("binance futures get ticker: %w", err)
		}

		for _, t := range resp {
			res = append(res, &exchange.Ticker{
				Symbol:             t.Symbol,
				PriceChange:        t.PriceChange,
				PriceChangePercent: t.PriceChangePercent,
				WeightedAvgPrice:   t.WeightedAvgPrice,
				LastPrice:          t.LastPrice,
				LastQty:            t.LastQuantity,
				OpenPrice:          t.OpenPrice,
				HighPrice:          t.HighPrice,
				LowPrice:           t.LowPrice,
				Volume:             t.Volume,
				QuoteVolume:        t.QuoteVolume,
				Count:              t.Count,
			})
		}
	}
	return &exchange.Tickers{Tickers: res}, nil
}
