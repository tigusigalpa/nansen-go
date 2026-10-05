package nansen

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func exact(t *testing.T, lexeme string) ExactNumber {
	t.Helper()
	n, err := ExactNumberFromLexeme(lexeme)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func TestExactRequestFiltersPreserveNumericLexemes(t *testing.T) {
	cases := []struct {
		name, path string
		call       func(*Client, *ExactNumericRangeFilter) error
	}{
		{"balances", "/api/v1/profiler/address/current-balance", func(c *Client, f *ExactNumericRangeFilter) error {
			_, e := c.Profiler.AddressCurrentBalanceExact(context.Background(), &ProfilerAddressBalancesExactRequest{Address: "0xabc", Chain: ChainEthereum, Filters: &ProfilerAddressBalancesExactFilters{ValueUSD: f}})
			return e
		}},
		{"netflow", "/api/v1/smart-money/netflow", func(c *Client, f *ExactNumericRangeFilter) error {
			_, e := c.SmartMoney.NetflowExact(context.Background(), &SmartMoneyNetflowExactRequest{Chains: []Chain{ChainEthereum}, Filters: &SmartMoneyNetflowExactFilters{MarketCapUSD: f}})
			return e
		}},
		{"holdings", "/api/v1/smart-money/holdings", func(c *Client, f *ExactNumericRangeFilter) error {
			_, e := c.SmartMoney.HoldingsExact(context.Background(), &SmartMoneyHoldingsExactRequest{Chains: []Chain{ChainEthereum}, Filters: &SmartMoneyHoldingsExactFilters{ValueUSD: f}})
			return e
		}},
		{"historical", "/api/v1/smart-money/historical-holdings", func(c *Client, f *ExactNumericRangeFilter) error {
			_, e := c.SmartMoney.HistoricalHoldingsExact(context.Background(), &SmartMoneyHistoricalHoldingsExactRequest{DateRange: DateOnlyRange{From: "2026-01-01"}, Chains: []SmartMoneyHistoricalHoldingsChain{HistoricalHoldingsChainEthereum}, Filters: &SmartMoneyHistoricalHoldingsExactFilters{Balance: f}})
			return e
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var body string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tc.path {
					t.Errorf("path = %s", r.URL.Path)
				}
				raw, _ := io.ReadAll(r.Body)
				body = string(raw)
				_, _ = w.Write([]byte(`{"data":[],"pagination":{"page":1,"per_page":10,"is_last_page":true}}`))
			}))
			defer srv.Close()
			filter := &ExactNumericRangeFilter{Min: exact(t, "9007199254740993"), Max: exact(t, "0.12345678901234567890123456789")}
			if err := tc.call(newTestClient(t, srv), filter); err != nil {
				t.Fatal(err)
			}
			if !contains(body, "9007199254740993") || !contains(body, "0.12345678901234567890123456789") {
				t.Fatalf("exact lexemes missing: %s", body)
			}
		})
	}
}

func TestExactRangeRejectsInvalidBoundsBeforeHTTP(t *testing.T) {
	invalid := []string{"", `"1"`, "[]", "{}", "true", "1 2", "NaN", "Infinity"}
	for _, lexeme := range invalid {
		t.Run(lexeme, func(t *testing.T) {
			var attempts int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { atomic.AddInt32(&attempts, 1) }))
			defer srv.Close()
			c := newTestClient(t, srv)
			_, err := c.SmartMoney.HoldingsExact(context.Background(), &SmartMoneyHoldingsExactRequest{Chains: []Chain{ChainEthereum}, Filters: &SmartMoneyHoldingsExactFilters{ValueUSD: &ExactNumericRangeFilter{Min: ExactNumber{Lexeme: lexeme, Present: true}}}})
			if err == nil || atomic.LoadInt32(&attempts) != 0 {
				t.Errorf("err=%v attempts=%d", err, attempts)
			}
		})
	}
}

func contains(s, part string) bool {
	return len(part) == 0 || (len(s) >= len(part) && index(s, part) >= 0)
}
func index(s, part string) int {
	for i := 0; i+len(part) <= len(s); i++ {
		if s[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
