package nansen

import (
	"context"
	"encoding/json"
	"fmt"
)

// ProfilerLabelsChain identifies chains supported by the address-label endpoints.
type ProfilerLabelsChain string

// ProfilerLabelsChainAll through ProfilerLabelsChainTron are supported by address labels.
const (
	ProfilerLabelsChainAll         ProfilerLabelsChain = "all"
	ProfilerLabelsChainArbitrum    ProfilerLabelsChain = "arbitrum"
	ProfilerLabelsChainArc         ProfilerLabelsChain = "arc"
	ProfilerLabelsChainAvalanche   ProfilerLabelsChain = "avalanche"
	ProfilerLabelsChainBase        ProfilerLabelsChain = "base"
	ProfilerLabelsChainBNB         ProfilerLabelsChain = "bnb"
	ProfilerLabelsChainEthereum    ProfilerLabelsChain = "ethereum"
	ProfilerLabelsChainHyperEVM    ProfilerLabelsChain = "hyperevm"
	ProfilerLabelsChainHyperliquid ProfilerLabelsChain = "hyperliquid"
	ProfilerLabelsChainIotaEVM     ProfilerLabelsChain = "iotaevm"
	ProfilerLabelsChainLinea       ProfilerLabelsChain = "linea"
	ProfilerLabelsChainMantle      ProfilerLabelsChain = "mantle"
	ProfilerLabelsChainMonad       ProfilerLabelsChain = "monad"
	ProfilerLabelsChainOptimism    ProfilerLabelsChain = "optimism"
	ProfilerLabelsChainPlasma      ProfilerLabelsChain = "plasma"
	ProfilerLabelsChainPolygon     ProfilerLabelsChain = "polygon"
	ProfilerLabelsChainRobinhood   ProfilerLabelsChain = "robinhood"
	ProfilerLabelsChainSei         ProfilerLabelsChain = "sei"
	ProfilerLabelsChainSolana      ProfilerLabelsChain = "solana"
	ProfilerLabelsChainSonic       ProfilerLabelsChain = "sonic"
	ProfilerLabelsChainTron        ProfilerLabelsChain = "tron"
)

// ProfilerAddressLabelsRequest specifies a non-premium address-label query.
type ProfilerAddressLabelsRequest struct {
	Address    string              `json:"address"`
	Chain      ProfilerLabelsChain `json:"chain"`
	Pagination *PaginationRequest  `json:"pagination,omitempty"`
}

// ProfilerAddressPremiumLabelsRequest specifies a premium address-label query.
type ProfilerAddressPremiumLabelsRequest struct {
	Address    string              `json:"address"`
	Chain      ProfilerLabelsChain `json:"chain"`
	Pagination *PaginationRequest  `json:"pagination,omitempty"`
}

// ProfilerAddressLabel is a documented label row. It does not imply taxonomy version or identity.
type ProfilerAddressLabel struct {
	Label    string   `json:"label"`
	Category *string  `json:"category,omitempty"`
	Kind     []string `json:"kind,omitempty"`
}

// ProfilerAddressLabelsResponse contains non-premium labels and the exact provider receipt.
type ProfilerAddressLabelsResponse struct {
	Data       []ProfilerAddressLabel     `json:"data"`
	Pagination PaginationInfo             `json:"pagination"`
	Raw        json.RawMessage            `json:"-"`
	Unknown    map[string]json.RawMessage `json:"-"`
}

// ProfilerAddressPremiumLabelsResponse contains premium labels and the exact provider receipt.
type ProfilerAddressPremiumLabelsResponse ProfilerAddressLabelsResponse

func isProfilerLabelsChain(chain ProfilerLabelsChain) bool {
	switch chain {
	case ProfilerLabelsChainAll, ProfilerLabelsChainArbitrum, ProfilerLabelsChainArc,
		ProfilerLabelsChainAvalanche, ProfilerLabelsChainBase, ProfilerLabelsChainBNB,
		ProfilerLabelsChainEthereum, ProfilerLabelsChainHyperEVM, ProfilerLabelsChainHyperliquid,
		ProfilerLabelsChainIotaEVM, ProfilerLabelsChainLinea, ProfilerLabelsChainMantle,
		ProfilerLabelsChainMonad, ProfilerLabelsChainOptimism, ProfilerLabelsChainPlasma,
		ProfilerLabelsChainPolygon, ProfilerLabelsChainRobinhood, ProfilerLabelsChainSei,
		ProfilerLabelsChainSolana, ProfilerLabelsChainSonic, ProfilerLabelsChainTron:
		return true
	}
	return false
}

func validateProfilerLabelsRequest(address string, chain ProfilerLabelsChain) error {
	if address == "" {
		return fmt.Errorf("nansen: address is required")
	}
	if !isProfilerLabelsChain(chain) {
		return fmt.Errorf("nansen: chain %q is not supported by profiler address labels", chain)
	}
	return nil
}

// UnmarshalJSON preserves the complete provider response alongside documented fields.
func (r *ProfilerAddressLabelsResponse) UnmarshalJSON(data []byte) error {
	type response struct {
		Data       []ProfilerAddressLabel `json:"data"`
		Pagination PaginationInfo         `json:"pagination"`
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

// UnmarshalJSON preserves the complete premium-label response alongside documented fields.
func (r *ProfilerAddressPremiumLabelsResponse) UnmarshalJSON(data []byte) error {
	var decoded ProfilerAddressLabelsResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*r = ProfilerAddressPremiumLabelsResponse(decoded)
	return nil
}

// ProfilerService provides access to wallet-level Profiler endpoints.
type ProfilerService struct {
	client *Client
}

// AddressLabels returns non-premium labels for an address on a supported chain.
func (s *ProfilerService) AddressLabels(ctx context.Context, req *ProfilerAddressLabelsRequest) (*ProfilerAddressLabelsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("nansen: address labels request is required")
	}
	if err := validateProfilerLabelsRequest(req.Address, req.Chain); err != nil {
		return nil, err
	}
	resp := &ProfilerAddressLabelsResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/profiler/address/labels", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AddressPremiumLabels returns all labels, including premium labels, for an address.
func (s *ProfilerService) AddressPremiumLabels(ctx context.Context, req *ProfilerAddressPremiumLabelsRequest) (*ProfilerAddressPremiumLabelsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("nansen: premium address labels request is required")
	}
	if err := validateProfilerLabelsRequest(req.Address, req.Chain); err != nil {
		return nil, err
	}
	resp := &ProfilerAddressPremiumLabelsResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/profiler/address/premium-labels", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ProfilerAddressBalancesSortField enumerates the sortable balance fields.
type ProfilerAddressBalancesSortField string

// BalanceSortValueUSD and BalanceSortTokenSymbol are fields sortable by address balances.
const (
	BalanceSortValueUSD    ProfilerAddressBalancesSortField = "value_usd"
	BalanceSortTokenSymbol ProfilerAddressBalancesSortField = "token_symbol"
)

// ProfilerAddressBalancesFilters limits profiler address balance results.
type ProfilerAddressBalancesFilters struct {
	ValueUSD     *NumericRangeFilter `json:"value_usd,omitempty"`
	PriceUSD     *NumericRangeFilter `json:"price_usd,omitempty"`
	TokenAmount  *IntegerRangeFilter `json:"token_amount,omitempty"`
	TokenSymbol  interface{}         `json:"token_symbol,omitempty"`
	TokenAddress interface{}         `json:"token_address,omitempty"`
	TokenName    interface{}         `json:"token_name,omitempty"`
}

// ProfilerAddressBalancesRequest specifies an address balance query.
type ProfilerAddressBalancesRequest struct {
	Address       string                          `json:"address,omitempty"`
	EntityName    string                          `json:"entity_name,omitempty"`
	Chain         Chain                           `json:"chain"`
	HideSpamToken *bool                           `json:"hide_spam_token,omitempty"`
	Filters       *ProfilerAddressBalancesFilters `json:"filters,omitempty"`
	Pagination    *PaginationRequest              `json:"pagination,omitempty"`
	OrderBy       []SortOrder                     `json:"order_by,omitempty"`
}

// ProfilerBalance contains a token balance held by an address.
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

// ProfilerAddressBalancesResponse contains paginated address balance results.
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

// ProfilerDexSortChain through ProfilerDexSortTradeValueUSD are fields sortable by profiler DEX trades.
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

// ProfilerDexTradeFilters limits profiler DEX trade results.
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

// ProfilerDexTradeRequest specifies a profiler DEX trades query.
type ProfilerDexTradeRequest struct {
	Address    string                   `json:"address"`
	Chain      Chain                    `json:"chain"`
	Date       DateRange                `json:"date"`
	Filters    *ProfilerDexTradeFilters `json:"filters,omitempty"`
	Pagination *PaginationRequest       `json:"pagination,omitempty"`
	OrderBy    []SortOrder              `json:"order_by,omitempty"`
}

// ProfilerDexTrade contains a DEX trade associated with an address.
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

// ProfilerDexTradeResponse contains paginated profiler DEX trade results.
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
