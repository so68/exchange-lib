package gate

import (
	"context"
	"math/big"

	"github.com/so68/exchange-lib/exchange"
)

// GetSpotBalance 获取现货余额
func (g *gateSpot) GetSpotBalance(ctx context.Context) ([]exchange.Balance, error) {
	bal, _, err := g.client.SpotApi.ListSpotAccounts(ctx, nil)
	if err != nil {
		return nil, err
	}
	var res []exchange.Balance
	for _, account := range bal {
		if account.Available == "0" && account.Locked == "0" {
			continue
		}
		availableFloat := new(big.Float).SetPrec(64)
		lockedFloat := new(big.Float).SetPrec(64)
		if _, ok := availableFloat.SetString(account.Available); !ok {
			continue
		}
		if _, ok := lockedFloat.SetString(account.Locked); !ok {
			continue
		}
		total := new(big.Float).Add(availableFloat, lockedFloat).Text('f', -1)
		res = append(res, exchange.Balance{
			Symbol: account.Currency,
			Free:   account.Available,
			Locked: account.Locked,
			Total:  total,
		})
	}
	return res, nil
}

// GetFuturesBalance 获取合约余额
func (g *gateSpot) GetFuturesBalance(ctx context.Context) ([]exchange.Balance, error) {
	// Gate.io 支持多种结算货币，通常使用 USDT
	settle := "usdt"
	account, _, err := g.client.FuturesApi.ListFuturesAccounts(ctx, settle)
	if err != nil {
		return nil, err
	}

	var res []exchange.Balance
	// 跳过余额为 0 的账户
	availableFloat := new(big.Float).SetPrec(64)
	orderMarginFloat := new(big.Float).SetPrec(64)
	if _, ok := availableFloat.SetString(account.Available); !ok {
		return res, nil
	}
	if _, ok := orderMarginFloat.SetString(account.OrderMargin); !ok {
		return res, nil
	}

	if availableFloat.Cmp(big.NewFloat(0)) == 0 && orderMarginFloat.Cmp(big.NewFloat(0)) == 0 {
		return res, nil
	}

	// 计算总余额：可用余额 + 订单保证金
	total := new(big.Float).Add(availableFloat, orderMarginFloat).Text('f', -1)

	res = append(res, exchange.Balance{
		Symbol: account.Currency,
		Free:   account.Available,
		Locked: account.OrderMargin,
		Total:  total,
	})

	return res, nil
}
