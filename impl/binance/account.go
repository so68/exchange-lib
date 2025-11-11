package binance

import (
	"context"
	"math/big"

	"github.com/so68/exchange-lib/exchange"
)

// SpotBalance 获取现货余额
func (b *binanceExchange) SpotBalance(ctx context.Context) ([]exchange.Balance, error) {
	acc, err := b.client.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, err
	}
	var res []exchange.Balance
	for _, bal := range acc.Balances {
		if bal.Free == "0.00000000" && bal.Locked == "0.00000000" {
			continue
		}
		freeFloat := new(big.Float).SetPrec(64)
		lockedFloat := new(big.Float).SetPrec(64)
		if _, ok := freeFloat.SetString(bal.Free); !ok {
			continue
		}
		if _, ok := lockedFloat.SetString(bal.Locked); !ok {
			continue
		}
		total := new(big.Float).Add(freeFloat, lockedFloat).Text('f', -1)
		res = append(res, exchange.Balance{
			Symbol: bal.Asset,
			Free:   bal.Free,
			Locked: bal.Locked,
			Total:  total,
		})
	}
	return res, nil
}

// FuturesBalance 获取合约余额
func (b *binanceExchange) FuturesBalance(ctx context.Context) ([]exchange.Balance, error) {
	acc, err := b.futuresClient.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, err
	}

	var res []exchange.Balance
	for _, asset := range acc.Assets {
		// 解析余额
		availableFloat := new(big.Float).SetPrec(64)
		orderMarginFloat := new(big.Float).SetPrec(64)

		if _, ok := availableFloat.SetString(asset.AvailableBalance); !ok {
			continue
		}
		if _, ok := orderMarginFloat.SetString(asset.OpenOrderInitialMargin); !ok {
			continue
		}

		// 跳过余额为 0 的资产
		if availableFloat.Cmp(big.NewFloat(0)) == 0 && orderMarginFloat.Cmp(big.NewFloat(0)) == 0 {
			continue
		}

		// 总余额使用钱包余额（WalletBalance）
		res = append(res, exchange.Balance{
			Symbol: asset.Asset,
			Free:   asset.AvailableBalance,
			Locked: asset.OpenOrderInitialMargin,
			Total:  asset.WalletBalance,
		})
	}

	return res, nil
}
