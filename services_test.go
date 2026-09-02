package nansen

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServicesUseExpectedEndpoints(t *testing.T) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTestClient(t, server)
	tests := []struct {
		name string
		path string
		call func(context.Context) error
	}{
		{"smart money netflow", "/api/v1/smart-money/netflow", func(ctx context.Context) error { _, err := client.SmartMoney.Netflow(ctx, nil); return err }},
		{"smart money holdings", "/api/v1/smart-money/holdings", func(ctx context.Context) error { _, err := client.SmartMoney.Holdings(ctx, nil); return err }},
		{"smart money DEX trades", "/api/v1/smart-money/dex-trades", func(ctx context.Context) error { _, err := client.SmartMoney.DEXTrades(ctx, nil); return err }},
		{"token screener", "/api/v1/token-screener", func(ctx context.Context) error { _, err := client.TokenGodMode.TokenScreener(ctx, nil); return err }},
		{"flow intelligence", "/api/v1/tgm/flow-intelligence", func(ctx context.Context) error { _, err := client.TokenGodMode.FlowIntelligence(ctx, nil); return err }},
		{"who bought sold", "/api/v1/tgm/who-bought-sold", func(ctx context.Context) error { _, err := client.TokenGodMode.WhoBoughtSold(ctx, nil); return err }},
		{"address balance", "/api/v1/profiler/address/current-balance", func(ctx context.Context) error { _, err := client.Profiler.AddressCurrentBalance(ctx, nil); return err }},
		{"profiler DEX trades", "/api/v1/profiler/dex-trades", func(ctx context.Context) error { _, err := client.Profiler.AddressDEXTrades(ctx, nil); return err }},
		{"DeFi holdings", "/api/v1/portfolio/defi-holdings", func(ctx context.Context) error { _, err := client.Portfolio.DeFiHoldings(ctx, nil); return err }},
		{"historical token flow", "/api/v1beta1/tgm/historical-token-flow-summary", func(ctx context.Context) error {
			_, err := client.Historical.HistoricalTokenFlowSummary(ctx, nil)
			return err
		}},
		{"historical token balances", "/api/v1beta1/smart-money/historical-token-balances", func(ctx context.Context) error {
			_, err := client.Historical.HistoricalSmartMoneyTokenBalances(ctx, nil)
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.baseURL = server.URL
			server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Errorf("path = %q, want %q", r.URL.Path, tt.path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{}`))
			})
			if err := tt.call(context.Background()); err != nil {
				t.Fatalf("call() error = %v", err)
			}
		})
	}
}
