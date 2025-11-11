package gate

import (
	"context"

	"github.com/so68/exchange-lib/exchange"
)

func (g *gateSpot) SpotCreateOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateSpot) SpotGetOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateSpot) SpotCancelOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateSpot) FuturesCreateOrder(ctx context.Context, symbol string, side exchange.OrderSide, orderType exchange.OrderType, quantity, price string) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateSpot) FuturesGetOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}

func (g *gateSpot) FuturesCancelOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}
