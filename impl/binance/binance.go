package binance

import (
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/so68/exchange-lib/exchange"
)

// 现货实例
type binanceExchange struct {
	client        *binance.Client
	futuresClient *futures.Client
}

// 创建现货实例
func NewBinance(apiKey, secretKey string) exchange.Exchange {
	return &binanceExchange{
		client:        binance.NewClient(apiKey, secretKey),
		futuresClient: futures.NewClient(apiKey, secretKey),
	}
}
