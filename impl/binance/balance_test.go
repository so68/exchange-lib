package binance

import (
	"context"
	"fmt"
	"testing"
)

// TestSpotBalance 获取现货余额
// go test -v ./impl/binance -run "^TestSpotBalance$"
func TestSpotBalance(t *testing.T) {
	binanceExchange := NewBinance("uAF7jJkOMygW8VK24OEzrrRNYjDFierW04lutVKtPJT9kLUNfHsByLe7lM7dWi4y", "jscrjDQsTpsiL8nDExD8Ty52raja6g8f4d0VOXByQfyKrYAmYtDUUCDSxymhmvgu")
	balances, err := binanceExchange.SpotBalance(context.Background())
	if err != nil {
		t.Fatalf("获取现货余额失败: %v", err)
	}

	for _, balance := range balances {
		fmt.Printf("币种: %s, 可用余额: %s, 锁定余额: %s, 总余额: %s\n", balance.Symbol, balance.Free, balance.Locked, balance.Total)
	}
}

// TestFuturesBalance 获取合约余额 go test -v ./impl/binance -run "^TestFuturesBalance$"
func TestFuturesBalance(t *testing.T) {
	binanceExchange := NewBinance("uAF7jJkOMygW8VK24OEzrrRNYjDFierW04lutVKtPJT9kLUNfHsByLe7lM7dWi4y", "jscrjDQsTpsiL8nDExD8Ty52raja6g8f4d0VOXByQfyKrYAmYtDUUCDSxymhmvgu")
	balances, err := binanceExchange.FuturesBalance(context.Background())
	if err != nil {
		t.Fatalf("获取合约余额失败: %v", err)
	}

	for _, balance := range balances {
		fmt.Printf("币种: %s, 可用余额: %s, 锁定余额: %s, 总余额: %s\n", balance.Symbol, balance.Free, balance.Locked, balance.Total)
	}
}
