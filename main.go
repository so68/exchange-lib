package main

import (
	"fmt"
	"math/big"
	"time"
)

func main() {
	reloadFunc()
	select {}
}

func reloadFunc() {
	quantity := "0.00823197123123"
	quantityFloat := new(big.Float).SetPrec(64)
	if _, ok := quantityFloat.SetString(quantity); !ok {
		panic("111")
	}

	precision := 6

	// 向下取整，如果保留 6位小数
	// 逻辑：floor(quantity * 10^precision) / 10^precision
	multiplier := new(big.Float).SetPrec(64).SetInt(big.NewInt(10).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil))
	// 乘以 10^precision
	quantityFloat = quantityFloat.Mul(quantityFloat, multiplier)
	// 向下取整
	quantityInt, _ := quantityFloat.Int(nil)
	// 除以 10^precision
	quantityFloat = new(big.Float).SetInt(quantityInt)
	quantityFloat = quantityFloat.Quo(quantityFloat, multiplier)
	fmt.Println("=====>", quantityFloat.Text('f', precision))
	time.AfterFunc(2*time.Second, reloadFunc)
}
