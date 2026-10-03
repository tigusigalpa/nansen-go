package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tigusigalpa/nansen-go"
)

func main() {
	apiKey := os.Getenv("NANSEN_API_KEY")
	if apiKey == "" {
		log.Fatal("NANSEN_API_KEY is required")
	}

	client, err := nansen.New(apiKey,
		nansen.WithTimeout(20*time.Second),
		nansen.WithRetry(3, 500*time.Millisecond, 5*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Netflows
	netflowReq := &nansen.SmartMoneyNetflowRequest{
		Chains: []nansen.Chain{nansen.ChainSolana},
		Filters: &nansen.SmartMoneyNetflowFilters{
			IncludeSmartMoneyLabels: []nansen.SmartMoneyLabel{nansen.LabelFund},
			IncludeStablecoins:      nansen.BoolPtr(false),
		},
		Pagination: &nansen.PaginationRequest{
			Page:    nansen.IntPtr(1),
			PerPage: nansen.IntPtr(10),
		},
		OrderBy: []nansen.SortOrder{
			{Field: string(nansen.NetflowSortNetFlow7DUSD), Direction: nansen.SortDesc},
		},
	}

	netflows, err := client.SmartMoney.Netflow(ctx, netflowReq)
	if err != nil {
		log.Fatalf("netflow request failed: %v", err)
	}
	fmt.Printf("Netflows (page %d/%d, last=%v)\n", netflows.Pagination.Page, netflows.Pagination.PerPage, netflows.Pagination.IsLastPage)
	for _, n := range netflows.Data {
		fmt.Printf("  %s: 1h=%v 24h=%v 7d=%v 30d=%v\n", n.TokenSymbol,
			ptrFloat(n.NetFlow1HUsd), ptrFloat(n.NetFlow24HUsd), ptrFloat(n.NetFlow7DUsd), ptrFloat(n.NetFlow30DUsd))
	}

	// Holdings
	ctx2, cancel2 := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel2()

	holdingsReq := &nansen.SmartMoneyHoldingsRequest{
		Chains: []nansen.Chain{nansen.ChainEthereum},
		Filters: &nansen.SmartMoneyHoldingsFilters{
			IncludeStablecoins: nansen.BoolPtr(false),
		},
		Pagination: &nansen.PaginationRequest{
			Page:    nansen.IntPtr(1),
			PerPage: nansen.IntPtr(10),
		},
		OrderBy: []nansen.SortOrder{
			{Field: string(nansen.HoldingsSortValueUSD), Direction: nansen.SortDesc},
		},
	}

	holdings, err := client.SmartMoney.Holdings(ctx2, holdingsReq)
	if err != nil {
		log.Fatalf("holdings request failed: %v", err)
	}
	fmt.Printf("\nHoldings (page %d/%d, last=%v)\n", holdings.Pagination.Page, holdings.Pagination.PerPage, holdings.Pagination.IsLastPage)
	for _, h := range holdings.Data {
		fmt.Printf("  %s: valueUSD=%v holders=%d\n", h.TokenSymbol, ptrFloat(h.ValueUsd), h.HoldersCount)
	}

	// Daily historical holdings preserve provider numeric lexemes in ExactNumber.
	ctxHistorical, cancelHistorical := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelHistorical()
	historical, err := client.SmartMoney.HistoricalHoldings(ctxHistorical, &nansen.SmartMoneyHistoricalHoldingsRequest{
		DateRange: nansen.DateOnlyRange{
			From: time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
			To:   time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		},
		Chains: []nansen.SmartMoneyHistoricalHoldingsChain{nansen.HistoricalHoldingsChainEthereum},
		Pagination: &nansen.PaginationRequest{
			Page:    nansen.IntPtr(1),
			PerPage: nansen.IntPtr(10),
		},
	})
	if err != nil {
		log.Fatalf("historical holdings request failed: %v", err)
	}

	fmt.Println("\nHistorical holdings")
	for _, h := range historical.Data {
		fmt.Printf("  %s %s: balance=%s valueUSD=%s\n", h.Date, h.TokenSymbol, h.Balance.Lexeme, h.ValueUSD.Lexeme)
	}

	// DEX trades
	ctx3, cancel3 := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel3()

	tradesReq := &nansen.SmartMoneyDexTradesRequest{
		Chains: []nansen.Chain{nansen.ChainEthereum, nansen.ChainSolana},
		Pagination: &nansen.PaginationRequest{
			Page:    nansen.IntPtr(1),
			PerPage: nansen.IntPtr(5),
		},
		OrderBy: []nansen.SortOrder{
			{Field: string(nansen.DexTradesSortTradeValueUSD), Direction: nansen.SortDesc},
		},
	}

	trades, err := client.SmartMoney.DEXTrades(ctx3, tradesReq)
	if err != nil {
		log.Fatalf("DEX trades request failed: %v", err)
	}
	fmt.Printf("\nDEX trades (page %d/%d, last=%v)\n", trades.Pagination.Page, trades.Pagination.PerPage, trades.Pagination.IsLastPage)
	for _, t := range trades.Data {
		fmt.Printf("  %s -> %s valueUSD=%v\n", t.TokenSoldSymbol, t.TokenBoughtSymbol, ptrFloat(t.TradeValueUSD))
	}
}

func ptrFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
