package nansen

import (
	"context"
	"fmt"
	"time"
)

// ProfilerAddressBalancesExactFilters is the exact-number variant of
// ProfilerAddressBalancesFilters. Contract: /api/profiler/address-current-balances.
type ProfilerAddressBalancesExactFilters struct {
	ValueUSD     *ExactNumericRangeFilter `json:"value_usd,omitempty"`
	PriceUSD     *ExactNumericRangeFilter `json:"price_usd,omitempty"`
	TokenAmount  *IntegerRangeFilter      `json:"token_amount,omitempty"`
	TokenSymbol  interface{}              `json:"token_symbol,omitempty"`
	TokenAddress interface{}              `json:"token_address,omitempty"`
	TokenName    interface{}              `json:"token_name,omitempty"`
}

// ProfilerAddressBalancesExactRequest is the exact-number variant of the current-balance request.
type ProfilerAddressBalancesExactRequest struct {
	Address       string                               `json:"address,omitempty"`
	EntityName    string                               `json:"entity_name,omitempty"`
	Chain         Chain                                `json:"chain"`
	HideSpamToken *bool                                `json:"hide_spam_token,omitempty"`
	Filters       *ProfilerAddressBalancesExactFilters `json:"filters,omitempty"`
	Pagination    *PaginationRequest                   `json:"pagination,omitempty"`
	OrderBy       []SortOrder                          `json:"order_by,omitempty"`
}

// AddressCurrentBalanceExact executes a current-balance request with exact numeric filters.
func (s *ProfilerService) AddressCurrentBalanceExact(ctx context.Context, req *ProfilerAddressBalancesExactRequest) (*ProfilerAddressBalancesResponse, error) {
	resp := &ProfilerAddressBalancesResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/profiler/address/current-balance", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SmartMoneyNetflowExactFilters is the exact-number variant of SmartMoneyNetflowFilters.
type SmartMoneyNetflowExactFilters struct {
	IncludeSmartMoneyLabels []SmartMoneyLabel        `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []SmartMoneyLabel        `json:"exclude_smart_money_labels,omitempty"`
	TokenAddress            interface{}              `json:"token_address,omitempty"`
	IncludeStablecoins      *bool                    `json:"include_stablecoins,omitempty"`
	IncludeNativeTokens     *bool                    `json:"include_native_tokens,omitempty"`
	TokenSector             []string                 `json:"token_sector,omitempty"`
	TraderCount             *IntegerRangeFilter      `json:"trader_count,omitempty"`
	TokenAgeDays            *ExactNumericRangeFilter `json:"token_age_days,omitempty"`
	MarketCapUSD            *ExactNumericRangeFilter `json:"market_cap_usd,omitempty"`
}

// SmartMoneyNetflowExactRequest is the exact-number variant of SmartMoneyNetflowRequest.
type SmartMoneyNetflowExactRequest struct {
	Chains     []Chain                        `json:"chains"`
	Filters    *SmartMoneyNetflowExactFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest             `json:"pagination,omitempty"`
	OrderBy    []SortOrder                    `json:"order_by,omitempty"`
}

// NetflowExact executes a netflow request with exact numeric filters.
func (s *SmartMoneyService) NetflowExact(ctx context.Context, req *SmartMoneyNetflowExactRequest) (*SmartMoneyNetflowResponse, error) {
	resp := &SmartMoneyNetflowResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/smart-money/netflow", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SmartMoneyHoldingsExactFilters is the exact-number variant of SmartMoneyHoldingsFilters.
type SmartMoneyHoldingsExactFilters struct {
	IncludeSmartMoneyLabels []SmartMoneyLabel        `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []SmartMoneyLabel        `json:"exclude_smart_money_labels,omitempty"`
	IncludeStablecoins      *bool                    `json:"include_stablecoins,omitempty"`
	IncludeNativeTokens     *bool                    `json:"include_native_tokens,omitempty"`
	ValueUSD                *ExactNumericRangeFilter `json:"value_usd,omitempty"`
	Balance24HPercentChange *ExactNumericRangeFilter `json:"balance_24h_percent_change,omitempty"`
	HoldersCount            *IntegerRangeFilter      `json:"holders_count,omitempty"`
	ShareOfHoldingsPercent  *ExactNumericRangeFilter `json:"share_of_holdings_percent,omitempty"`
	TokenAgeDays            *ExactNumericRangeFilter `json:"token_age_days,omitempty"`
	MarketCapUSD            *ExactNumericRangeFilter `json:"market_cap_usd,omitempty"`
	TokenAddress            interface{}              `json:"token_address,omitempty"`
	TokenSymbol             interface{}              `json:"token_symbol,omitempty"`
	TokenSectors            []string                 `json:"token_sectors,omitempty"`
}

// SmartMoneyHoldingsExactRequest is the exact-number variant of SmartMoneyHoldingsRequest.
type SmartMoneyHoldingsExactRequest struct {
	Chains     []Chain                         `json:"chains"`
	Filters    *SmartMoneyHoldingsExactFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest              `json:"pagination,omitempty"`
	OrderBy    []SortOrder                     `json:"order_by,omitempty"`
}

// HoldingsExact executes a current holdings request with exact numeric filters.
func (s *SmartMoneyService) HoldingsExact(ctx context.Context, req *SmartMoneyHoldingsExactRequest) (*SmartMoneyHoldingsResponse, error) {
	resp := &SmartMoneyHoldingsResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/smart-money/holdings", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SmartMoneyHistoricalHoldingsExactFilters is the exact-number variant of historical holdings filters.
type SmartMoneyHistoricalHoldingsExactFilters struct {
	IncludeSmartMoneyLabels []HistoricalSmartMoneyFilterType `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []HistoricalSmartMoneyFilterType `json:"exclude_smart_money_labels,omitempty"`
	IncludeStablecoins      *bool                            `json:"include_stablecoins,omitempty"`
	IncludeNativeTokens     *bool                            `json:"include_native_tokens,omitempty"`
	Balance                 *ExactNumericRangeFilter         `json:"balance,omitempty"`
	ValueUSD                *ExactNumericRangeFilter         `json:"value_usd,omitempty"`
	Balance24HPercentChange *ExactNumericRangeFilter         `json:"balance_24h_percent_change,omitempty"`
	HoldersCount            *IntegerRangeFilter              `json:"holders_count,omitempty"`
	ShareOfHoldingsPercent  *ExactNumericRangeFilter         `json:"share_of_holdings_percent,omitempty"`
	TokenAgeDays            *ExactNumericRangeFilter         `json:"token_age_days,omitempty"`
	MarketCapUSD            *ExactNumericRangeFilter         `json:"market_cap_usd,omitempty"`
	TokenAddress            interface{}                      `json:"token_address,omitempty"`
	TokenSymbol             interface{}                      `json:"token_symbol,omitempty"`
}

// SmartMoneyHistoricalHoldingsExactRequest is the exact-number variant of the v1 daily holdings request.
type SmartMoneyHistoricalHoldingsExactRequest struct {
	DateRange  DateOnlyRange                             `json:"date_range"`
	Chains     []SmartMoneyHistoricalHoldingsChain       `json:"chains"`
	Filters    *SmartMoneyHistoricalHoldingsExactFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest                        `json:"pagination,omitempty"`
	OrderBy    []SortOrder                               `json:"order_by,omitempty"`
}

// HistoricalHoldingsExact executes a v1 historical holdings request with exact numeric filters.
func (s *SmartMoneyService) HistoricalHoldingsExact(ctx context.Context, req *SmartMoneyHistoricalHoldingsExactRequest) (*SmartMoneyHistoricalHoldingsResponse, error) {
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
