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
		nansen.WithBaseURL("https://api.nansen.ai"),
		nansen.WithTimeout(30*time.Second),
		nansen.WithRetry(3, 500*time.Millisecond, 5*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	address := "0x0000000000000000000000000000000000000000"

	// Current balances
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	balances, err := client.Profiler.AddressCurrentBalance(ctx, &nansen.ProfilerAddressBalancesRequest{
		Address:       address,
		Chain:         nansen.ChainEthereum,
		HideSpamToken: nansen.BoolPtr(true),
		Pagination: &nansen.PaginationRequest{
			Page:    nansen.IntPtr(1),
			PerPage: nansen.IntPtr(20),
		},
	})
	if err != nil {
		log.Fatalf("current balance request failed: %v", err)
	}

	fmt.Printf("Balances for %s (page %d/%d)\n", address, balances.Pagination.Page, balances.Pagination.PerPage)
	for _, b := range balances.Data {
		fmt.Printf("  %s: amount=%v valueUSD=%v\n", b.TokenSymbol, ptrFloat(b.TokenAmount), ptrFloat(b.ValueUsd))
	}

	// DEX trades
	ctx2, cancel2 := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel2()

	trades, err := client.Profiler.AddressDEXTrades(ctx2, &nansen.ProfilerDexTradeRequest{
		Address: address,
		Chain:   nansen.ChainEthereum,
		Date: nansen.DateRange{
			From: time.Now().AddDate(0, 0, -7).Format(time.RFC3339),
			To:   time.Now().Format(time.RFC3339),
		},
		Pagination: &nansen.PaginationRequest{
			Page:    nansen.IntPtr(1),
			PerPage: nansen.IntPtr(10),
		},
	})
	if err != nil {
		log.Fatalf("DEX trades request failed: %v", err)
	}

	fmt.Printf("\nDEX trades for %s (page %d/%d)\n", address, trades.Pagination.Page, trades.Pagination.PerPage)
	for _, t := range trades.Data {
		fmt.Printf("  %s bought %s sold %s valueUSD=%v\n",
			t.BlockTimestamp, ptrString(t.TokenBoughtSymbol), ptrString(t.TokenSoldSymbol), ptrFloat(t.TradeValueUSD))
	}
}

func ptrFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func ptrString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
