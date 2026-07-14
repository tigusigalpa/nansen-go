package nansen

import "context"

// ProfilerService provides access to wallet-level Profiler endpoints.
type ProfilerService struct {
	client *Client
}

// ProfilerAddressBalancesSortField enumerates the sortable balance fields.
type ProfilerAddressBalancesSortField string

const (
	BalanceSortValueUSD    ProfilerAddressBalancesSortField = "value_usd"
	BalanceSortTokenSymbol ProfilerAddressBalancesSortField = "token_symbol"
)

type ProfilerAddressBalancesFilters struct {
	ValueUSD     *NumericRangeFilter `json:"value_usd,omitempty"`
	PriceUSD     *NumericRangeFilter `json:"price_usd,omitempty"`
	TokenAmount  *IntegerRangeFilter `json:"token_amount,omitempty"`
	TokenSymbol  interface{}         `json:"token_symbol,omitempty"`
	TokenAddress interface{}         `json:"token_address,omitempty"`
	TokenName    interface{}         `json:"token_name,omitempty"`
}

type ProfilerAddressBalancesRequest struct {
	Address       string                          `json:"address,omitempty"`
	EntityName    string                          `json:"entity_name,omitempty"`
	Chain         Chain                           `json:"chain"`
	HideSpamToken *bool                           `json:"hide_spam_token,omitempty"`
	Filters       *ProfilerAddressBalancesFilters `json:"filters,omitempty"`
	Pagination    *PaginationRequest              `json:"pagination,omitempty"`
	OrderBy       []SortOrder                     `json:"order_by,omitempty"`
}

type ProfilerBalance struct {
	Chain        string   `json:"chain"`
	Address      string   `json:"address"`
	TokenAddress string   `json:"token_address"`
	TokenSymbol  string   `json:"token_symbol"`
	TokenName    *string  `json:"token_name,omitempty"`
	TokenAmount  *float64 `json:"token_amount,omitempty"`
	PriceUsd     *float64 `json:"price_usd,omitempty"`
	ValueUsd     *float64 `json:"value_usd,omitempty"`
}

type ProfilerAddressBalancesResponse struct {
	Data       []ProfilerBalance `json:"data"`
	Pagination PaginationInfo    `json:"pagination"`
}

// AddressCurrentBalance returns current token balances for an address or entity.
func (s *ProfilerService) AddressCurrentBalance(ctx context.Context, req *ProfilerAddressBalancesRequest) (*ProfilerAddressBalancesResponse, error) {
	resp := &ProfilerAddressBalancesResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/profiler/address/current-balance", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ProfilerDexTradeSortField enumerates the sortable trade fields.
type ProfilerDexTradeSortField string

const (
	ProfilerDexSortChain                ProfilerDexTradeSortField = "chain"
	ProfilerDexSortBlockTimestamp       ProfilerDexTradeSortField = "block_timestamp"
	ProfilerDexSortTransactionHash      ProfilerDexTradeSortField = "transaction_hash"
	ProfilerDexSortTokenBoughtAmount    ProfilerDexTradeSortField = "token_bought_amount"
	ProfilerDexSortTokenSoldAmount      ProfilerDexTradeSortField = "token_sold_amount"
	ProfilerDexSortTokenBoughtSymbol    ProfilerDexTradeSortField = "token_bought_symbol"
	ProfilerDexSortTokenSoldSymbol      ProfilerDexTradeSortField = "token_sold_symbol"
	ProfilerDexSortTokenBoughtAgeDays   ProfilerDexTradeSortField = "token_bought_age_days"
	ProfilerDexSortTokenSoldAgeDays     ProfilerDexTradeSortField = "token_sold_age_days"
	ProfilerDexSortTokenBoughtMarketCap ProfilerDexTradeSortField = "token_bought_market_cap"
	ProfilerDexSortTokenSoldMarketCap   ProfilerDexTradeSortField = "token_sold_market_cap"
	ProfilerDexSortTokenBoughtFDV       ProfilerDexTradeSortField = "token_bought_fdv"
	ProfilerDexSortTokenSoldFDV         ProfilerDexTradeSortField = "token_sold_fdv"
	ProfilerDexSortTradeValueUSD        ProfilerDexTradeSortField = "trade_value_usd"
)

type ProfilerDexTradeFilters struct {
	TokenBoughtAddress   string              `json:"token_bought_address,omitempty"`
	TokenSoldAddress     string              `json:"token_sold_address,omitempty"`
	TokenBoughtSymbol    string              `json:"token_bought_symbol,omitempty"`
	TokenSoldSymbol      string              `json:"token_sold_symbol,omitempty"`
	TokenBoughtAmount    *NumericRangeFilter `json:"token_bought_amount,omitempty"`
	TokenSoldAmount      *NumericRangeFilter `json:"token_sold_amount,omitempty"`
	TokenBoughtAgeDays   *IntegerRangeFilter `json:"token_bought_age_days,omitempty"`
	TokenSoldAgeDays     *IntegerRangeFilter `json:"token_sold_age_days,omitempty"`
	TokenBoughtMarketCap *NumericRangeFilter `json:"token_bought_market_cap,omitempty"`
	TokenSoldMarketCap   *NumericRangeFilter `json:"token_sold_market_cap,omitempty"`
	TokenBoughtFDV       *NumericRangeFilter `json:"token_bought_fdv,omitempty"`
	TokenSoldFDV         *NumericRangeFilter `json:"token_sold_fdv,omitempty"`
	TradeValueUSD        *NumericRangeFilter `json:"trade_value_usd,omitempty"`
}

type ProfilerDexTradeRequest struct {
	Address    string                   `json:"address"`
	Chain      Chain                    `json:"chain"`
	Date       DateRange                `json:"date"`
	Filters    *ProfilerDexTradeFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest       `json:"pagination,omitempty"`
	OrderBy    []SortOrder              `json:"order_by,omitempty"`
}

type ProfilerDexTrade struct {
	Chain                string   `json:"chain"`
	BlockTimestamp       string   `json:"block_timestamp"`
	TransactionHash      string   `json:"transaction_hash"`
	TraderAddress        string   `json:"trader_address"`
	TraderAddressLabel   *string  `json:"trader_address_label,omitempty"`
	TokenBoughtAddress   string   `json:"token_bought_address"`
	TokenSoldAddress     string   `json:"token_sold_address"`
	TokenBoughtAmount    *float64 `json:"token_bought_amount,omitempty"`
	TokenSoldAmount      *float64 `json:"token_sold_amount,omitempty"`
	TokenBoughtSymbol    *string  `json:"token_bought_symbol,omitempty"`
	TokenSoldSymbol      *string  `json:"token_sold_symbol,omitempty"`
	TokenBoughtAgeDays   *int     `json:"token_bought_age_days,omitempty"`
	TokenSoldAgeDays     *int     `json:"token_sold_age_days,omitempty"`
	TokenBoughtMarketCap *float64 `json:"token_bought_market_cap,omitempty"`
	TokenSoldMarketCap   *float64 `json:"token_sold_market_cap,omitempty"`
	TokenBoughtFDV       *float64 `json:"token_bought_fdv,omitempty"`
	TokenSoldFDV         *float64 `json:"token_sold_fdv,omitempty"`
	TradeValueUSD        *float64 `json:"trade_value_usd,omitempty"`
}

type ProfilerDexTradeResponse struct {
	Data       []ProfilerDexTrade `json:"data"`
	Pagination PaginationInfo     `json:"pagination"`
}

// AddressDEXTrades returns DEX trade history for a wallet address on a specific chain.
func (s *ProfilerService) AddressDEXTrades(ctx context.Context, req *ProfilerDexTradeRequest) (*ProfilerDexTradeResponse, error) {
	resp := &ProfilerDexTradeResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/profiler/dex-trades", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
