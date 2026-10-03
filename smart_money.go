package nansen

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SmartMoneyService provides access to Smart Money endpoints.
type SmartMoneyService struct {
	client *Client
}

// SmartMoneyHistoricalHoldingsChain identifies chains supported by historical holdings.
type SmartMoneyHistoricalHoldingsChain string

// HistoricalHoldingsChainArc through HistoricalHoldingsChainSolana are supported by historical holdings.
const (
	HistoricalHoldingsChainArc       SmartMoneyHistoricalHoldingsChain = "arc"
	HistoricalHoldingsChainBase      SmartMoneyHistoricalHoldingsChain = "base"
	HistoricalHoldingsChainBNB       SmartMoneyHistoricalHoldingsChain = "bnb"
	HistoricalHoldingsChainEthereum  SmartMoneyHistoricalHoldingsChain = "ethereum"
	HistoricalHoldingsChainMonad     SmartMoneyHistoricalHoldingsChain = "monad"
	HistoricalHoldingsChainRobinhood SmartMoneyHistoricalHoldingsChain = "robinhood"
	HistoricalHoldingsChainSolana    SmartMoneyHistoricalHoldingsChain = "solana"
)

// DateOnlyRange is a YYYY-MM-DD date range. To may be omitted and defaults to today upstream.
type DateOnlyRange struct {
	From string `json:"from"`
	To   string `json:"to,omitempty"`
}

// SmartMoneyHistoricalHoldingsFilters limits historical smart-money holdings results.
type SmartMoneyHistoricalHoldingsFilters struct {
	IncludeSmartMoneyLabels []HistoricalSmartMoneyFilterType `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []HistoricalSmartMoneyFilterType `json:"exclude_smart_money_labels,omitempty"`
	IncludeStablecoins      *bool                            `json:"include_stablecoins,omitempty"`
	IncludeNativeTokens     *bool                            `json:"include_native_tokens,omitempty"`
	Balance                 *NumericRangeFilter              `json:"balance,omitempty"`
	ValueUSD                *NumericRangeFilter              `json:"value_usd,omitempty"`
	Balance24HPercentChange *NumericRangeFilter              `json:"balance_24h_percent_change,omitempty"`
	HoldersCount            *IntegerRangeFilter              `json:"holders_count,omitempty"`
	ShareOfHoldingsPercent  *NumericRangeFilter              `json:"share_of_holdings_percent,omitempty"`
	TokenAgeDays            *NumericRangeFilter              `json:"token_age_days,omitempty"`
	MarketCapUSD            *NumericRangeFilter              `json:"market_cap_usd,omitempty"`
	TokenAddress            interface{}                      `json:"token_address,omitempty"`
	TokenSymbol             interface{}                      `json:"token_symbol,omitempty"`
}

// SmartMoneyHistoricalHoldingsRequest specifies a v1 daily historical holdings query.
type SmartMoneyHistoricalHoldingsRequest struct {
	DateRange  DateOnlyRange                        `json:"date_range"`
	Chains     []SmartMoneyHistoricalHoldingsChain  `json:"chains"`
	Filters    *SmartMoneyHistoricalHoldingsFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest                   `json:"pagination,omitempty"`
	OrderBy    []SortOrder                          `json:"order_by,omitempty"`
}

// SmartMoneyHistoricalHolding contains one provider-dated daily holding snapshot.
type SmartMoneyHistoricalHolding struct {
	Date                    string      `json:"date"`
	Chain                   string      `json:"chain"`
	TokenAddress            string      `json:"token_address"`
	TokenSymbol             string      `json:"token_symbol"`
	TokenSectors            []string    `json:"token_sectors"`
	SmartMoneyLabels        []string    `json:"smart_money_labels"`
	Balance                 ExactNumber `json:"balance"`
	ValueUSD                ExactNumber `json:"value_usd"`
	Balance24HPercentChange ExactNumber `json:"balance_24h_percent_change"`
	HoldersCount            int         `json:"holders_count"`
	ShareOfHoldingsPercent  ExactNumber `json:"share_of_holdings_percent"`
	TokenAgeDays            int         `json:"token_age_days"`
	MarketCapUSD            ExactNumber `json:"market_cap_usd"`
}

// SmartMoneyHistoricalHoldingsResponse contains daily historical holdings and the exact provider receipt.
type SmartMoneyHistoricalHoldingsResponse struct {
	Data       []SmartMoneyHistoricalHolding `json:"data"`
	Pagination PaginationInfo                `json:"pagination"`
	Raw        json.RawMessage               `json:"-"`
	Unknown    map[string]json.RawMessage    `json:"-"`
}

// UnmarshalJSON preserves the exact response and unknown outer fields.
func (r *SmartMoneyHistoricalHoldingsResponse) UnmarshalJSON(data []byte) error {
	type response struct {
		Data       []SmartMoneyHistoricalHolding `json:"data"`
		Pagination PaginationInfo                `json:"pagination"`
	}
	var decoded response
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	delete(fields, "data")
	delete(fields, "pagination")
	r.Data, r.Pagination, r.Raw, r.Unknown = decoded.Data, decoded.Pagination, append(r.Raw[:0], data...), fields
	return nil
}

// HistoricalHoldings returns daily smart-money holdings for the requested date range.
func (s *SmartMoneyService) HistoricalHoldings(ctx context.Context, req *SmartMoneyHistoricalHoldingsRequest) (*SmartMoneyHistoricalHoldingsResponse, error) {
	if req == nil || req.DateRange.From == "" || len(req.Chains) == 0 {
		return nil, fmt.Errorf("nansen: historical holdings require date_range.from and chains")
	}
	if _, err := time.Parse("2006-01-02", req.DateRange.From); err != nil {
		return nil, fmt.Errorf("nansen: date_range.from must use YYYY-MM-DD: %w", err)
	}
	if req.DateRange.To != "" {
		if _, err := time.Parse("2006-01-02", req.DateRange.To); err != nil {
			return nil, fmt.Errorf("nansen: date_range.to must use YYYY-MM-DD: %w", err)
		}
	}
	for _, chain := range req.Chains {
		switch chain {
		case HistoricalHoldingsChainArc, HistoricalHoldingsChainBase, HistoricalHoldingsChainBNB,
			HistoricalHoldingsChainEthereum, HistoricalHoldingsChainMonad,
			HistoricalHoldingsChainRobinhood, HistoricalHoldingsChainSolana:
		default:
			return nil, fmt.Errorf("nansen: chain %q is not supported by historical holdings", chain)
		}
	}
	resp := &SmartMoneyHistoricalHoldingsResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/smart-money/historical-holdings", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SmartMoneyNetflowSortField enumerates the sortable fields for netflows.
type SmartMoneyNetflowSortField string

// NetflowSortChain through NetflowSortMarketCapUSD are fields sortable by netflow.
const (
	NetflowSortChain         SmartMoneyNetflowSortField = "chain"
	NetflowSortTokenAddress  SmartMoneyNetflowSortField = "token_address"
	NetflowSortTokenSymbol   SmartMoneyNetflowSortField = "token_symbol"
	NetflowSortNetFlow1HUSD  SmartMoneyNetflowSortField = "net_flow_1h_usd"
	NetflowSortNetFlow24HUSD SmartMoneyNetflowSortField = "net_flow_24h_usd"
	NetflowSortNetFlow7DUSD  SmartMoneyNetflowSortField = "net_flow_7d_usd"
	NetflowSortNetFlow30DUSD SmartMoneyNetflowSortField = "net_flow_30d_usd"
	NetflowSortTokenSectors  SmartMoneyNetflowSortField = "token_sectors"
	NetflowSortTraderCount   SmartMoneyNetflowSortField = "trader_count"
	NetflowSortTokenAgeDays  SmartMoneyNetflowSortField = "token_age_days"
	NetflowSortMarketCapUSD  SmartMoneyNetflowSortField = "market_cap_usd"
)

// SmartMoneyNetflowFilters limits smart-money netflow results.
type SmartMoneyNetflowFilters struct {
	IncludeSmartMoneyLabels []SmartMoneyLabel   `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []SmartMoneyLabel   `json:"exclude_smart_money_labels,omitempty"`
	TokenAddress            interface{}         `json:"token_address,omitempty"`
	IncludeStablecoins      *bool               `json:"include_stablecoins,omitempty"`
	IncludeNativeTokens     *bool               `json:"include_native_tokens,omitempty"`
	TokenSector             []string            `json:"token_sector,omitempty"`
	TraderCount             *IntegerRangeFilter `json:"trader_count,omitempty"`
	TokenAgeDays            *NumericRangeFilter `json:"token_age_days,omitempty"`
	MarketCapUSD            *NumericRangeFilter `json:"market_cap_usd,omitempty"`
}

// SmartMoneyNetflowRequest specifies a smart-money netflow query.
type SmartMoneyNetflowRequest struct {
	Chains     []Chain                   `json:"chains"`
	Filters    *SmartMoneyNetflowFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest        `json:"pagination,omitempty"`
	OrderBy    []SortOrder               `json:"order_by,omitempty"`
}

// SmartMoneyNetflow contains netflow metrics for one token and chain.
type SmartMoneyNetflow struct {
	TokenAddress  string   `json:"token_address"`
	TokenSymbol   string   `json:"token_symbol"`
	NetFlow1HUsd  *float64 `json:"net_flow_1h_usd,omitempty"`
	NetFlow24HUsd *float64 `json:"net_flow_24h_usd,omitempty"`
	NetFlow7DUsd  *float64 `json:"net_flow_7d_usd,omitempty"`
	NetFlow30DUsd *float64 `json:"net_flow_30d_usd,omitempty"`
	Chain         string   `json:"chain"`
	TokenSectors  []string `json:"token_sectors"`
	TraderCount   int      `json:"trader_count"`
	TokenAgeDays  int      `json:"token_age_days"`
	MarketCapUsd  *float64 `json:"market_cap_usd,omitempty"`
}

// SmartMoneyNetflowResponse contains paginated netflow results.
type SmartMoneyNetflowResponse struct {
	Data       []SmartMoneyNetflow `json:"data"`
	Pagination PaginationInfo      `json:"pagination"`
}

// Netflow returns net capital flows for smart money wallets.
func (s *SmartMoneyService) Netflow(ctx context.Context, req *SmartMoneyNetflowRequest) (*SmartMoneyNetflowResponse, error) {
	resp := &SmartMoneyNetflowResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/smart-money/netflow", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SmartMoneyHoldingsSortField enumerates the sortable fields for holdings.
type SmartMoneyHoldingsSortField string

// HoldingsSortChain through HoldingsSortMarketCapUSD are fields sortable by holdings.
const (
	HoldingsSortChain              SmartMoneyHoldingsSortField = "chain"
	HoldingsSortTokenAddress       SmartMoneyHoldingsSortField = "token_address"
	HoldingsSortTokenSymbol        SmartMoneyHoldingsSortField = "token_symbol"
	HoldingsSortValueUSD           SmartMoneyHoldingsSortField = "value_usd"
	HoldingsSortBalance24HChange   SmartMoneyHoldingsSortField = "balance_24h_percent_change"
	HoldingsSortHoldersCount       SmartMoneyHoldingsSortField = "holders_count"
	HoldingsSortShareOfHoldingsPct SmartMoneyHoldingsSortField = "share_of_holdings_percent"
	HoldingsSortTokenAgeDays       SmartMoneyHoldingsSortField = "token_age_days"
	HoldingsSortMarketCapUSD       SmartMoneyHoldingsSortField = "market_cap_usd"
)

// SmartMoneyHoldingsFilters limits smart-money holdings results.
type SmartMoneyHoldingsFilters struct {
	IncludeSmartMoneyLabels []SmartMoneyLabel   `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []SmartMoneyLabel   `json:"exclude_smart_money_labels,omitempty"`
	IncludeStablecoins      *bool               `json:"include_stablecoins,omitempty"`
	IncludeNativeTokens     *bool               `json:"include_native_tokens,omitempty"`
	ValueUSD                *NumericRangeFilter `json:"value_usd,omitempty"`
	Balance24HPercentChange *NumericRangeFilter `json:"balance_24h_percent_change,omitempty"`
	HoldersCount            *IntegerRangeFilter `json:"holders_count,omitempty"`
	ShareOfHoldingsPercent  *NumericRangeFilter `json:"share_of_holdings_percent,omitempty"`
	TokenAgeDays            *NumericRangeFilter `json:"token_age_days,omitempty"`
	MarketCapUSD            *NumericRangeFilter `json:"market_cap_usd,omitempty"`
	TokenAddress            interface{}         `json:"token_address,omitempty"`
	TokenSymbol             interface{}         `json:"token_symbol,omitempty"`
	TokenSectors            []string            `json:"token_sectors,omitempty"`
}

// SmartMoneyHoldingsRequest specifies a smart-money holdings query.
type SmartMoneyHoldingsRequest struct {
	Chains     []Chain                    `json:"chains"`
	Filters    *SmartMoneyHoldingsFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest         `json:"pagination,omitempty"`
	OrderBy    []SortOrder                `json:"order_by,omitempty"`
}

// SmartMoneyHolding contains holdings metrics for one token and chain.
type SmartMoneyHolding struct {
	Chain                   string   `json:"chain"`
	TokenAddress            string   `json:"token_address"`
	TokenSymbol             string   `json:"token_symbol"`
	TokenSectors            []string `json:"token_sectors"`
	ValueUsd                *float64 `json:"value_usd,omitempty"`
	Balance24HPercentChange *float64 `json:"balance_24h_percent_change,omitempty"`
	HoldersCount            int      `json:"holders_count"`
	ShareOfHoldingsPercent  *float64 `json:"share_of_holdings_percent,omitempty"`
	TokenAgeDays            int      `json:"token_age_days"`
	MarketCapUsd            *float64 `json:"market_cap_usd,omitempty"`
}

// SmartMoneyHoldingsResponse contains paginated holdings results.
type SmartMoneyHoldingsResponse struct {
	Data       []SmartMoneyHolding `json:"data"`
	Pagination PaginationInfo      `json:"pagination"`
}

// Holdings returns aggregated token balances held by smart money.
func (s *SmartMoneyService) Holdings(ctx context.Context, req *SmartMoneyHoldingsRequest) (*SmartMoneyHoldingsResponse, error) {
	resp := &SmartMoneyHoldingsResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/smart-money/holdings", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SmartMoneyDexTradesSortField enumerates the sortable fields for DEX trades.
type SmartMoneyDexTradesSortField string

// DexTradesSortChain through DexTradesSortTradeValueUSD are fields sortable by DEX trades.
const (
	DexTradesSortChain                SmartMoneyDexTradesSortField = "chain"
	DexTradesSortBlockTimestamp       SmartMoneyDexTradesSortField = "block_timestamp"
	DexTradesSortTransactionHash      SmartMoneyDexTradesSortField = "transaction_hash"
	DexTradesSortTraderAddress        SmartMoneyDexTradesSortField = "trader_address"
	DexTradesSortTraderAddressLabel   SmartMoneyDexTradesSortField = "trader_address_label"
	DexTradesSortTokenBoughtAddress   SmartMoneyDexTradesSortField = "token_bought_address"
	DexTradesSortTokenSoldAddress     SmartMoneyDexTradesSortField = "token_sold_address"
	DexTradesSortTokenBoughtAmount    SmartMoneyDexTradesSortField = "token_bought_amount"
	DexTradesSortTokenSoldAmount      SmartMoneyDexTradesSortField = "token_sold_amount"
	DexTradesSortTokenBoughtSymbol    SmartMoneyDexTradesSortField = "token_bought_symbol"
	DexTradesSortTokenSoldSymbol      SmartMoneyDexTradesSortField = "token_sold_symbol"
	DexTradesSortTokenBoughtAgeDays   SmartMoneyDexTradesSortField = "token_bought_age_days"
	DexTradesSortTokenSoldAgeDays     SmartMoneyDexTradesSortField = "token_sold_age_days"
	DexTradesSortTokenBoughtMarketCap SmartMoneyDexTradesSortField = "token_bought_market_cap"
	DexTradesSortTokenSoldMarketCap   SmartMoneyDexTradesSortField = "token_sold_market_cap"
	DexTradesSortTokenBoughtFDV       SmartMoneyDexTradesSortField = "token_bought_fdv"
	DexTradesSortTokenSoldFDV         SmartMoneyDexTradesSortField = "token_sold_fdv"
	DexTradesSortTradeValueUSD        SmartMoneyDexTradesSortField = "trade_value_usd"
)

// SmartMoneyDexTradesFilters limits smart-money DEX trade results.
type SmartMoneyDexTradesFilters struct {
	IncludeSmartMoneyLabels []SmartMoneyLabel   `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []SmartMoneyLabel   `json:"exclude_smart_money_labels,omitempty"`
	Chain                   interface{}         `json:"chain,omitempty"`
	TransactionHash         interface{}         `json:"transaction_hash,omitempty"`
	TraderAddress           interface{}         `json:"trader_address,omitempty"`
	TraderAddressLabel      interface{}         `json:"trader_address_label,omitempty"`
	TokenBoughtAddress      interface{}         `json:"token_bought_address,omitempty"`
	TokenSoldAddress        interface{}         `json:"token_sold_address,omitempty"`
	TokenBoughtAmount       *NumericRangeFilter `json:"token_bought_amount,omitempty"`
	TokenSoldAmount         *NumericRangeFilter `json:"token_sold_amount,omitempty"`
	TokenBoughtSymbol       interface{}         `json:"token_bought_symbol,omitempty"`
	TokenSoldSymbol         interface{}         `json:"token_sold_symbol,omitempty"`
	TokenBoughtAgeDays      *NumericRangeFilter `json:"token_bought_age_days,omitempty"`
	TokenSoldAgeDays        *NumericRangeFilter `json:"token_sold_age_days,omitempty"`
	TokenBoughtMarketCap    *NumericRangeFilter `json:"token_bought_market_cap,omitempty"`
	TokenSoldMarketCap      *NumericRangeFilter `json:"token_sold_market_cap,omitempty"`
	TokenBoughtFDV          *NumericRangeFilter `json:"token_bought_fdv,omitempty"`
	TokenSoldFDV            *NumericRangeFilter `json:"token_sold_fdv,omitempty"`
	TradeValueUSD           *NumericRangeFilter `json:"trade_value_usd,omitempty"`
}

// SmartMoneyDexTradesRequest specifies a smart-money DEX trades query.
type SmartMoneyDexTradesRequest struct {
	Chains     []Chain                     `json:"chains"`
	Filters    *SmartMoneyDexTradesFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest          `json:"pagination,omitempty"`
	OrderBy    []SortOrder                 `json:"order_by,omitempty"`
}

// SmartMoneyDexTrade contains a smart-money DEX trade.
type SmartMoneyDexTrade struct {
	Chain                string   `json:"chain"`
	BlockTimestamp       string   `json:"block_timestamp"`
	TransactionHash      string   `json:"transaction_hash"`
	TraderAddress        string   `json:"trader_address"`
	TraderAddressLabel   string   `json:"trader_address_label"`
	TokenBoughtAddress   string   `json:"token_bought_address"`
	TokenSoldAddress     string   `json:"token_sold_address"`
	TokenBoughtAmount    *float64 `json:"token_bought_amount,omitempty"`
	TokenSoldAmount      *float64 `json:"token_sold_amount,omitempty"`
	TokenBoughtSymbol    string   `json:"token_bought_symbol"`
	TokenSoldSymbol      string   `json:"token_sold_symbol"`
	TokenBoughtAgeDays   int      `json:"token_bought_age_days"`
	TokenSoldAgeDays     int      `json:"token_sold_age_days"`
	TokenBoughtMarketCap *float64 `json:"token_bought_market_cap,omitempty"`
	TokenSoldMarketCap   *float64 `json:"token_sold_market_cap,omitempty"`
	TokenBoughtFDV       *float64 `json:"token_bought_fdv,omitempty"`
	TokenSoldFDV         *float64 `json:"token_sold_fdv,omitempty"`
	TradeValueUSD        *float64 `json:"trade_value_usd,omitempty"`
}

// SmartMoneyDexTradesResponse contains paginated DEX trade results.
type SmartMoneyDexTradesResponse struct {
	Data       []SmartMoneyDexTrade `json:"data"`
	Pagination PaginationInfo       `json:"pagination"`
}

// DEXTrades returns real-time DEX trading activity from smart money.
func (s *SmartMoneyService) DEXTrades(ctx context.Context, req *SmartMoneyDexTradesRequest) (*SmartMoneyDexTradesResponse, error) {
	resp := &SmartMoneyDexTradesResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/smart-money/dex-trades", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
