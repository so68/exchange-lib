package gate

import (
	"context"
	"fmt"
	"testing"
)

// TestSpotBalance 获取现货余额
// go test -v ./impl/gate -run "^TestSpotBalance$"
func TestSpotBalance(t *testing.T) {
	gateExchange := NewGateExchange(apiKey, secretKey)
	balances, err := gateExchange.GetSpotBalance(context.Background())
	if err != nil {
		t.Fatalf("获取现货余额失败: %v", err)
	}

	for _, balance := range balances {
		fmt.Printf("币种: %s, 可用余额: %s, 锁定余额: %s, 总余额: %s\n", balance.Symbol, balance.Free, balance.Locked, balance.Total)
	}
}

// TestFuturesBalance 获取合约余额
// go test -v ./impl/gate -run "^TestFuturesBalance$"
func TestFuturesBalance(t *testing.T) {
	gateExchange := NewGateExchange(apiKey, secretKey)
	balances, err := gateExchange.GetFuturesBalance(context.Background())
	if err != nil {
		t.Fatalf("获取合约余额失败: %v", err)
	}

	for _, balance := range balances {
		fmt.Printf("币种: %s, 可用余额: %s, 锁定余额: %s, 总余额: %s\n", balance.Symbol, balance.Free, balance.Locked, balance.Total)
	}
}
