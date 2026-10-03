package nansen

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddressLabelsOfflineAcceptance(t *testing.T) {
	fixture := []byte(`{"data":[{"label":"Known Entity","category":"cefi","kind":["entity"]}],"pagination":{"page":1,"per_page":10,"is_last_page":true},"future_field":"retained"}`)
	for _, tt := range []struct {
		name, path, body string
		call             func(*Client) error
	}{
		{"common", "/api/v1/profiler/address/labels", `{"address":"0xabc","chain":"ethereum","pagination":{"page":1,"per_page":10}}`, func(c *Client) error {
			r, err := c.Profiler.AddressLabels(context.Background(), &ProfilerAddressLabelsRequest{Address: "0xabc", Chain: ProfilerLabelsChainEthereum, Pagination: &PaginationRequest{Page: IntPtr(1), PerPage: IntPtr(10)}})
			if err == nil && (!bytes.Equal(r.Raw, fixture) || r.Unknown["future_field"] == nil) {
				t.Error("raw receipt or unknown field was not retained")
			}
			return err
		}},
		{"premium", "/api/v1/profiler/address/premium-labels", `{"address":"0xabc","chain":"ethereum","pagination":{"page":1,"per_page":10}}`, func(c *Client) error {
			r, err := c.Profiler.AddressPremiumLabels(context.Background(), &ProfilerAddressPremiumLabelsRequest{Address: "0xabc", Chain: ProfilerLabelsChainEthereum, Pagination: &PaginationRequest{Page: IntPtr(1), PerPage: IntPtr(10)}})
			if err == nil && !bytes.Equal(r.Raw, fixture) {
				t.Error("raw receipt was not retained")
			}
			return err
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != tt.path {
					t.Errorf("request = %s %s", r.Method, r.URL.Path)
				}
				body, _ := io.ReadAll(r.Body)
				if string(body) != tt.body {
					t.Errorf("body = %s, want %s", body, tt.body)
				}
				_, _ = w.Write(fixture)
			}))
			defer srv.Close()
			if err := tt.call(newTestClient(t, srv)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestHistoricalHoldingsOfflineAcceptance(t *testing.T) {
	fixture := []byte(`{"data":[{"date":"2026-10-01","chain":"ethereum","token_address":"0x1","token_symbol":"TOK","token_sectors":["defi"],"smart_money_labels":["Smart Trader"],"balance":123456789.123456789,"value_usd":null,"balance_24h_percent_change":0,"holders_count":0,"share_of_holdings_percent":1.20e-3,"token_age_days":0,"market_cap_usd":999999999999999999.01}],"pagination":{"page":1,"per_page":10,"is_last_page":true},"provider_revision":"unknown"}`)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/smart-money/historical-holdings" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		want := `{"date_range":{"from":"2026-10-01","to":"2026-10-02"},"chains":["ethereum"],"pagination":{"page":1,"per_page":10}}`
		if string(body) != want {
			t.Errorf("body = %s, want %s", body, want)
		}
		_, _ = w.Write(fixture)
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	resp, err := c.SmartMoney.HistoricalHoldings(context.Background(), &SmartMoneyHistoricalHoldingsRequest{DateRange: DateOnlyRange{From: "2026-10-01", To: "2026-10-02"}, Chains: []SmartMoneyHistoricalHoldingsChain{HistoricalHoldingsChainEthereum}, Pagination: &PaginationRequest{Page: IntPtr(1), PerPage: IntPtr(10)}})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(resp.Raw, fixture) || resp.Unknown["provider_revision"] == nil {
		t.Fatal("raw receipt or unknown field was not retained")
	}
	row := resp.Data[0]
	if row.Balance.Lexeme != "123456789.123456789" || !row.ValueUSD.Present || !row.ValueUSD.Null || row.Balance24HPercentChange.Lexeme != "0" || !row.MarketCapUSD.Present {
		t.Fatalf("exact numeric presence was not retained: %#v", row)
	}
}

func TestLabelsValidation(t *testing.T) {
	c, _ := New("key")
	if _, err := c.Profiler.AddressLabels(context.Background(), &ProfilerAddressLabelsRequest{Chain: ProfilerLabelsChainEthereum}); err == nil {
		t.Error("missing address accepted")
	}
	if _, err := c.Profiler.AddressLabels(context.Background(), &ProfilerAddressLabelsRequest{Address: "0xabc", Chain: "bitcoin"}); err == nil {
		t.Error("unsupported chain accepted")
	}
}
