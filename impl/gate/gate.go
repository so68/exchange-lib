package gate

import (
	"github.com/gateio/gateapi-go/v6"
	"github.com/so68/exchange-lib/exchange"
)

// 现货实例
type gateExchange struct {
	client *gateapi.APIClient
}

// 创建现货实例
func NewGateExchange(apiKey, secretKey string) exchange.Exchange {
	cfg := gateapi.NewConfiguration()
	cfg.Key = apiKey
	cfg.Secret = secretKey
	client := gateapi.NewAPIClient(cfg)
	return &gateExchange{client: client}
}

// 创建现货实例
func newGateExchange(apiKey, secretKey string) *gateExchange {
	cfg := gateapi.NewConfiguration()
	cfg.Key = apiKey
	cfg.Secret = secretKey
	client := gateapi.NewAPIClient(cfg)
	return &gateExchange{client: client}
}
