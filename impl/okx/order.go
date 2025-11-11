package okx

import (
	"context"

	"github.com/so68/exchange-lib/exchange"
)

func (o *okx) SpotCreateOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	return nil, nil
}

func (o *okx) SpotGetOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}

func (o *okx) SpotCancelOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}

func (o *okx) FuturesCreateOrder(ctx context.Context, symbol string, side exchange.OrderSide, limitPrice, quantity string) (*exchange.Order, error) {
	return nil, nil
}

func (o *okx) FuturesGetOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}

func (o *okx) FuturesCancelOrder(ctx context.Context, symbol string, orderID int64) (*exchange.Order, error) {
	return nil, nil
}
