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

	page := nansen.IntPtr(1)
	perPage := nansen.IntPtr(20)
	minMarketCap := nansen.Float64Ptr(1_000_000)
	maxMarketCap := nansen.Float64Ptr(50_000_000)

	req := &nansen.TokenScreenerRequest{
		Chains: []nansen.Chain{nansen.ChainEthereum, nansen.ChainSolana},
		Timeframe: func() *nansen.TokenScreenerTimeframe {
			t := nansen.Timeframe24H
			return &t
		}(),
		Pagination: &nansen.PaginationRequest{
			Page:    page,
			PerPage: perPage,
		},
		Filters: &nansen.TokenScreenerFilters{
			MarketCapUSD: &nansen.NumericRangeFilter{
				Min: minMarketCap,
				Max: maxMarketCap,
			},
			IncludeStablecoins: nansen.BoolPtr(false),
		},
		OrderBy: []nansen.SortOrder{
			{Field: string(nansen.ScreenerSortVolume), Direction: nansen.SortDesc},
		},
	}

	resp, err := client.TokenGodMode.TokenScreener(ctx, req)
	if err != nil {
		if apiErr, ok := err.(*nansen.APIError); ok {
			log.Fatalf("API error %d: %s", apiErr.StatusCode, apiErr.Message)
		}
		log.Fatalf("token screener request failed: %v", err)
	}

	fmt.Printf("Page %d of %d per page; last=%v\n", resp.Pagination.Page, resp.Pagination.PerPage, resp.Pagination.IsLastPage)
	for _, t := range resp.Data {
		fmt.Printf("%s (%s) on %s - price=%v marketCap=%v volume=%v\n",
			t.TokenSymbol, t.TokenAddress, t.Chain,
			ptrFloat(t.PriceUSD), ptrFloat(t.MarketCapUSD), ptrFloat(t.Volume))
	}
}

func ptrFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
