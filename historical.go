package nansen

import "context"

// HistoricalService provides access to backtesting (v1beta1) endpoints.
type HistoricalService struct {
	client *Client
}

type TGMHistoricalTokenFlowSummaryRequest struct {
	Chain                Chain     `json:"chain"`
	TokenAddress         string    `json:"token_address"`
	DateRange            DateRange `json:"date_range"`
	ApplyBlacklistFilter *bool     `json:"apply_blacklist_filter,omitempty"`
}

type TGMHistoricalTokenFlowSummary struct {
	TokenSymbol             *string  `json:"token_symbol,omitempty"`
	PublicFigureNetFlowUSD  *float64 `json:"public_figure_net_flow_usd,omitempty"`
	PublicFigureAvgFlowUSD  *float64 `json:"public_figure_avg_flow_usd,omitempty"`
	PublicFigureWalletCount *int     `json:"public_figure_wallet_count,omitempty"`
	TopPnLNetFlowUSD        *float64 `json:"top_pnl_net_flow_usd,omitempty"`
	TopPnLAvgFlowUSD        *float64 `json:"top_pnl_avg_flow_usd,omitempty"`
	TopPnLWalletCount       *int     `json:"top_pnl_wallet_count,omitempty"`
	WhaleNetFlowUSD         *float64 `json:"whale_net_flow_usd,omitempty"`
	WhaleAvgFlowUSD         *float64 `json:"whale_avg_flow_usd,omitempty"`
	WhaleWalletCount        *int     `json:"whale_wallet_count,omitempty"`
	ExchangeNetFlowUSD      *float64 `json:"exchange_net_flow_usd,omitempty"`
	ExchangeAvgFlowUSD      *float64 `json:"exchange_avg_flow_usd,omitempty"`
	ExchangeWalletCount     *int     `json:"exchange_wallet_count,omitempty"`
	SmartTraderNetFlowUSD   *float64 `json:"smart_trader_net_flow_usd,omitempty"`
	SmartTraderAvgFlowUSD   *float64 `json:"smart_trader_avg_flow_usd,omitempty"`
	SmartTraderWalletCount  *int     `json:"smart_trader_wallet_count,omitempty"`
	FreshWalletsNetFlowUSD  *float64 `json:"fresh_wallets_net_flow_usd,omitempty"`
	FreshWalletsAvgFlowUSD  *float64 `json:"fresh_wallets_avg_flow_usd,omitempty"`
	FreshWalletsWalletCount *int     `json:"fresh_wallets_wallet_count,omitempty"`
}

type TGMHistoricalTokenFlowSummaryResponse struct {
	Data     []TGMHistoricalTokenFlowSummary `json:"data"`
	Warnings []string                        `json:"warnings,omitempty"`
}

// HistoricalTokenFlowSummary returns temporally-correct token flow intelligence.
func (s *HistoricalService) HistoricalTokenFlowSummary(ctx context.Context, req *TGMHistoricalTokenFlowSummaryRequest) (*TGMHistoricalTokenFlowSummaryResponse, error) {
	resp := &TGMHistoricalTokenFlowSummaryResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1beta1/tgm/historical-token-flow-summary", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// HistoricalSmartMoneyFilterType enumerates the label types available for
// historical smart-money balances.
type HistoricalSmartMoneyFilterType string

const (
	HistoricalLabelFund               HistoricalSmartMoneyFilterType = "Fund"
	HistoricalLabelSmartTrader        HistoricalSmartMoneyFilterType = "Smart Trader"
	HistoricalLabel30DSmartTrader     HistoricalSmartMoneyFilterType = "30D Smart Trader"
	HistoricalLabel90DSmartTrader     HistoricalSmartMoneyFilterType = "90D Smart Trader"
	HistoricalLabel180DSmartTrader    HistoricalSmartMoneyFilterType = "180D Smart Trader"
	HistoricalLabelSmartDexTrader     HistoricalSmartMoneyFilterType = "Smart Dex Trader"
	HistoricalLabel30DSmartDexTrader  HistoricalSmartMoneyFilterType = "30D Smart Dex Trader"
	HistoricalLabel90DSmartDexTrader  HistoricalSmartMoneyFilterType = "90D Smart Dex Trader"
	HistoricalLabel180DSmartDexTrader HistoricalSmartMoneyFilterType = "180D Smart Dex Trader"
	HistoricalLabelSmartHLPerpsTrader HistoricalSmartMoneyFilterType = "Smart HL Perps Trader"
)

type SmartMoneyHistoricalTokenBalancesFilters struct {
	SMFilter            []HistoricalSmartMoneyFilterType `json:"sm_filter,omitempty"`
	IncludeStablecoins  *bool                            `json:"include_stablecoins,omitempty"`
	IncludeNativeTokens *bool                            `json:"include_native_tokens,omitempty"`
	HoldersCount        *IntegerRangeFilter              `json:"holders_count,omitempty"`
}

type SmartMoneyHistoricalTokenBalancesRequest struct {
	AsOfDate             string                                    `json:"as_of_date"`
	Chains               []Chain                                   `json:"chains,omitempty"`
	Filters              *SmartMoneyHistoricalTokenBalancesFilters `json:"filters,omitempty"`
	ApplyBlacklistFilter *bool                                     `json:"apply_blacklist_filter,omitempty"`
	Pagination           *PaginationRequest                        `json:"pagination,omitempty"`
}

type SmartMoneyHistoricalTokenBalance struct {
	Chain                   string   `json:"chain"`
	TokenAddress            string   `json:"token_address"`
	TokenSymbol             string   `json:"token_symbol"`
	TokenSectors            []string `json:"token_sectors,omitempty"`
	ValueUsd                *float64 `json:"value_usd,omitempty"`
	Balance24HPercentChange *float64 `json:"balance_24h_percent_change,omitempty"`
	HoldersCount            *int     `json:"holders_count,omitempty"`
	ShareOfHoldingsPercent  *float64 `json:"share_of_holdings_percent,omitempty"`
	TokenAgeDays            *int     `json:"token_age_days,omitempty"`
	MarketCapUsd            *float64 `json:"market_cap_usd,omitempty"`
}

type SmartMoneyHistoricalTokenBalancesResponse struct {
	Data       []SmartMoneyHistoricalTokenBalance `json:"data"`
	Pagination PaginationInfo                     `json:"pagination"`
}

// HistoricalSmartMoneyTokenBalances returns point-in-time smart-money token balances.
func (s *HistoricalService) HistoricalSmartMoneyTokenBalances(ctx context.Context, req *SmartMoneyHistoricalTokenBalancesRequest) (*SmartMoneyHistoricalTokenBalancesResponse, error) {
	resp := &SmartMoneyHistoricalTokenBalancesResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1beta1/smart-money/historical-token-balances", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
