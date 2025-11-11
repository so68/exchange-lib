package binance

import (
	"context"
	"fmt"

	"github.com/so68/exchange-lib/exchange"
)

// GetSymbolTickers 获取交易对行情
func (b *binanceExchange) GetSymbolTickers(ctx context.Context, symbol ...string) (*exchange.Tickers, error) {
	resp, err := b.client.NewListSymbolTickerService().Symbols(symbol).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("binance get ticker: %w", err)
	}
	var res []*exchange.Ticker
	for _, ticker := range resp {
		res = append(res, &exchange.Ticker{
			Symbol:      ticker.Symbol,
			OpenPrice:   ticker.OpenPrice,
			HighPrice:   ticker.HighPrice,
			LowPrice:    ticker.LowPrice,
			LastPrice:   ticker.LastPrice,
			Volume:      ticker.Volume,
			QuoteVolume: ticker.QuoteVolume,
			Count:       ticker.Count,
		})
	}
	return &exchange.Tickers{
		Tickers: res,
	}, nil
}
