package gate

import (
	"context"

	"github.com/so68/exchange-lib/exchange"
)

func (g *gateSpot) GetSpotSymbolTickers(ctx context.Context, symbol ...string) (*exchange.Tickers, error) {
	return nil, nil
}

func (g *gateSpot) GetFuturesSymbolTickers(ctx context.Context, symbol ...string) (*exchange.Tickers, error) {
	return nil, nil
}
