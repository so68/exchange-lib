package gate

import (
	"github.com/gateio/gateapi-go/v6"
	"github.com/so68/exchange-lib/exchange"
)

// 现货实例
type gateSpot struct {
	client *gateapi.APIClient
}

// 创建现货实例
func NewGateSpot(apiKey, secretKey string) exchange.Exchange {
	cfg := gateapi.NewConfiguration()
	cfg.Key = apiKey
	cfg.Secret = secretKey
	client := gateapi.NewAPIClient(cfg)
	return &gateSpot{client: client}
}
