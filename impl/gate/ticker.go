package gate

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/so68/exchange-lib/exchange"
	"github.com/so68/exchange-lib/internal/utils"
)

// GetSpotSymbolTickers 获取现货交易对行情
func (g *gateExchange) GetSpotSymbolTickers(ctx context.Context, symbols ...string) (*exchange.Tickers, error) {
	tickers, _, err := g.client.SpotApi.ListTickers(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("获取现货交易对行情失败: %w", err)
	}

	symbols = utils.FormatSymbols(symbols)
	var data []*exchange.Ticker
	for _, ticker := range tickers {
		// 如果传入了 symbols，并且当前交易对不在 symbols 中，则跳过
		if len(symbols) > 0 && !slices.Contains(symbols, ticker.CurrencyPair) {
			continue
		}
		openPrice, priceChange := calculateOpenAndChangePrice(ticker.Last, ticker.ChangePercentage)
		data = append(data, &exchange.Ticker{
			Symbol:             ticker.CurrencyPair,
			PriceChange:        priceChange,
			PriceChangePercent: ticker.ChangePercentage,
			WeightedAvgPrice:   "", // Gate API 不提供加权平均价
			LastPrice:          ticker.Last,
			LastQty:            "", // Gate API 不提供最新成交量
			OpenPrice:          openPrice,
			HighPrice:          ticker.High24h,
			LowPrice:           ticker.Low24h,
			Volume:             ticker.BaseVolume,
			QuoteVolume:        ticker.QuoteVolume,
			Count:              0, // Gate API 不提供成交笔数
		})
	}

	return &exchange.Tickers{
		Tickers: data,
	}, nil
}

// GetFuturesSymbolTickers 获取合约交易对行情
func (g *gateExchange) GetFuturesSymbolTickers(ctx context.Context, symbols ...string) (*exchange.Tickers, error) {
	tickers, _, err := g.client.FuturesApi.ListFuturesTickers(ctx, strings.ToLower(Settle), nil)
	if err != nil {
		return nil, fmt.Errorf("获取合约交易对行情失败: %w", err)
	}

	symbols = utils.FormatSymbols(symbols)
	var data []*exchange.Ticker
	for _, ticker := range tickers {
		// 如果传入了 symbols，并且当前交易对不在 symbols 中，则跳过
		if len(symbols) > 0 && !slices.Contains(symbols, ticker.Contract) {
			continue
		}

		// 优先使用 Volume24hBase 作为 Volume，如果没有则使用 Volume24h
		volume := ticker.Volume24hBase
		if volume == "" {
			volume = ticker.Volume24h
		}
		// 优先使用 Volume24hQuote 作为 QuoteVolume，如果没有则使用 Volume24hSettle
		quoteVolume := ticker.Volume24hQuote
		if quoteVolume == "" {
			quoteVolume = ticker.Volume24hSettle
		}
		openPrice, priceChange := calculateOpenAndChangePrice(ticker.Last, ticker.ChangePercentage)
		data = append(data, &exchange.Ticker{
			Symbol:             ticker.Contract,
			PriceChange:        priceChange,
			PriceChangePercent: ticker.ChangePercentage,
			WeightedAvgPrice:   "", // Gate API 不提供加权平均价
			LastPrice:          ticker.Last,
			LastQty:            "", // Gate API 不提供最新成交量
			OpenPrice:          openPrice,
			HighPrice:          ticker.High24h,
			LowPrice:           ticker.Low24h,
			Volume:             volume,
			QuoteVolume:        quoteVolume,
			Count:              0, // Gate API 不提供成交笔数
		})
	}

	return &exchange.Tickers{
		Tickers: data,
	}, nil
}

// calculateOpenAndChangePrice 根据最新价和涨跌幅百分比计算开盘价和价格变动
// 公式：ChangePercentage = (Last - OpenPrice) / OpenPrice * 100
// 因此：OpenPrice = Last / (1 + ChangePercentage / 100)
// PriceChange = Last - OpenPrice
func calculateOpenAndChangePrice(lastPrice, changePercentage string) (openPrice, priceChange string) {
	lastFloat := new(big.Float).SetPrec(128)
	changeFloat := new(big.Float).SetPrec(128)

	// 解析最新价
	if _, ok := lastFloat.SetString(lastPrice); !ok {
		return "", ""
	}

	// 解析涨跌幅百分比
	if _, ok := changeFloat.SetString(changePercentage); !ok {
		return "", ""
	}

	// 计算：OpenPrice = Last / (1 + ChangePercentage / 100)
	// 先计算 1 + ChangePercentage / 100
	oneHundred := big.NewFloat(100)
	changeRatio := new(big.Float).Quo(changeFloat, oneHundred)
	denominator := new(big.Float).Add(big.NewFloat(1), changeRatio)

	// 计算开盘价
	openPriceFloat := new(big.Float).Quo(lastFloat, denominator)

	// 计算价格变动：PriceChange = Last - OpenPrice
	priceChangeFloat := new(big.Float).Sub(lastFloat, openPriceFloat)

	// 获取精度（使用 lastPrice 的精度）
	precision := utils.GetNumberPrecision(lastPrice)

	// 格式化为字符串
	openPrice = openPriceFloat.Text('f', precision)
	priceChange = priceChangeFloat.Text('f', precision)
	return openPrice, priceChange
}
