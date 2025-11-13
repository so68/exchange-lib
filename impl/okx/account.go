package okx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/so68/exchange-lib/exchange"
)

// SpotBalance 获取现货余额
func (o *okx) GetSpotBalance(ctx context.Context) ([]exchange.Balance, error) {
	resp, err := o.authRequest("GET", "/api/v5/account/balance", nil)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal spot balance data error: %w", err)
	}

	balances := make([]exchange.Balance, 0)
	for _, balance := range data {
		if _, ok := balance["details"]; ok {
			details := balance["details"].([]interface{})
			for _, detail := range details {
				detailMap := detail.(map[string]interface{})
				balances = append(balances, exchange.Balance{
					Symbol: detailMap["ccy"].(string),
					Free:   detailMap["availBal"].(string),
					Locked: detailMap["frozenBal"].(string),
					Total:  detailMap["cashBal"].(string),
				})
			}
		}
	}
	return balances, nil
}

// FuturesBalance 获取合约余额
func (o *okx) GetFuturesBalance(ctx context.Context) ([]exchange.Balance, error) {
	return o.GetSpotBalance(ctx)
}
