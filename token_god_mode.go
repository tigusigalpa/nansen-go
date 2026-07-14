package nansen

import "context"

// TokenGodModeService provides access to Token God Mode and Screener endpoints.
type TokenGodModeService struct {
	client *Client
}

// TokenScreenerTimeframe enumerates the supported screener time windows.
type TokenScreenerTimeframe string

const (
	Timeframe5M  TokenScreenerTimeframe = "5m"
	Timeframe10M TokenScreenerTimeframe = "10m"
	Timeframe1H  TokenScreenerTimeframe = "1h"
	Timeframe6H  TokenScreenerTimeframe = "6h"
	Timeframe24H TokenScreenerTimeframe = "24h"
	Timeframe7D  TokenScreenerTimeframe = "7d"
	Timeframe30D TokenScreenerTimeframe = "30d"
)

// TokenScreenerSortField enumerates the sortable fields for the token screener.
type TokenScreenerSortField string

const (
	ScreenerSortChain           TokenScreenerSortField = "chain"
	ScreenerSortTokenAddress    TokenScreenerSortField = "token_address"
	ScreenerSortTokenSymbol     TokenScreenerSortField = "token_symbol"
	ScreenerSortMarketCapUSD    TokenScreenerSortField = "market_cap_usd"
	ScreenerSortVolume          TokenScreenerSortField = "volume"
	ScreenerSortLiquidity       TokenScreenerSortField = "liquidity"
	ScreenerSortNofTraders      TokenScreenerSortField = "nof_traders"
	ScreenerSortNofBuyers       TokenScreenerSortField = "nof_buyers"
	ScreenerSortNofSellers      TokenScreenerSortField = "nof_sellers"
	ScreenerSortNofBuys         TokenScreenerSortField = "nof_buys"
	ScreenerSortNofSells        TokenScreenerSortField = "nof_sells"
	ScreenerSortPriceChange     TokenScreenerSortField = "price_change"
	ScreenerSortPriceUSD        TokenScreenerSortField = "price_usd"
	ScreenerSortNetflow         TokenScreenerSortField = "netflow"
	ScreenerSortBuyVolume       TokenScreenerSortField = "buy_volume"
	ScreenerSortSellVolume      TokenScreenerSortField = "sell_volume"
	ScreenerSortFDV             TokenScreenerSortField = "fdv"
	ScreenerSortFDVMCRatio      TokenScreenerSortField = "fdv_mc_ratio"
	ScreenerSortInflowFDVRatio  TokenScreenerSortField = "inflow_fdv_ratio"
	ScreenerSortOutflowFDVRatio TokenScreenerSortField = "outflow_fdv_ratio"
	ScreenerSortTokenAgeDays    TokenScreenerSortField = "token_age_days"
)

type TokenScreenerFilters struct {
	TokenAddress            interface{}         `json:"token_address,omitempty"`
	TokenSymbol             interface{}         `json:"token_symbol,omitempty"`
	OnlySmartMoney          *bool               `json:"only_smart_money,omitempty"`
	TraderType              *TraderType         `json:"trader_type,omitempty"`
	Sectors                 []string            `json:"sectors,omitempty"`
	ExcludeSectors          []string            `json:"exclude_sectors,omitempty"`
	TokenAgeDays            *NumericRangeFilter `json:"token_age_days,omitempty"`
	MarketCapUSD            *NumericRangeFilter `json:"market_cap_usd,omitempty"`
	Liquidity               *NumericRangeFilter `json:"liquidity,omitempty"`
	PriceUSD                *NumericRangeFilter `json:"price_usd,omitempty"`
	PriceChange             *NumericRangeFilter `json:"price_change,omitempty"`
	FDV                     *NumericRangeFilter `json:"fdv,omitempty"`
	FDVMCRatio              *NumericRangeFilter `json:"fdv_mc_ratio,omitempty"`
	NofBuyers               *IntegerRangeFilter `json:"nof_buyers,omitempty"`
	NofTraders              *IntegerRangeFilter `json:"nof_traders,omitempty"`
	NofSellers              *IntegerRangeFilter `json:"nof_sellers,omitempty"`
	NofBuys                 *IntegerRangeFilter `json:"nof_buys,omitempty"`
	NofSells                *IntegerRangeFilter `json:"nof_sells,omitempty"`
	BuyVolume               *NumericRangeFilter `json:"buy_volume,omitempty"`
	SellVolume              *NumericRangeFilter `json:"sell_volume,omitempty"`
	Volume                  *NumericRangeFilter `json:"volume,omitempty"`
	Netflow                 *NumericRangeFilter `json:"netflow,omitempty"`
	InflowFDVRatio          *NumericRangeFilter `json:"inflow_fdv_ratio,omitempty"`
	OutflowFDVRatio         *NumericRangeFilter `json:"outflow_fdv_ratio,omitempty"`
	IncludeStablecoins      *bool               `json:"include_stablecoins,omitempty"`
	IncludeSmartMoneyLabels []SmartMoneyLabel   `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []SmartMoneyLabel   `json:"exclude_smart_money_labels,omitempty"`
}

type TokenScreenerRequest struct {
	Chains     []Chain                 `json:"chains"`
	Timeframe  *TokenScreenerTimeframe `json:"timeframe,omitempty"`
	Date       *DateRange              `json:"date,omitempty"`
	Pagination *PaginationRequest      `json:"pagination,omitempty"`
	Filters    *TokenScreenerFilters   `json:"filters,omitempty"`
	OrderBy    []SortOrder             `json:"order_by,omitempty"`
}

type TokenScreenerResult struct {
	Chain               string   `json:"chain"`
	TokenAddress        string   `json:"token_address"`
	TokenSymbol         string   `json:"token_symbol"`
	TokenAgeDays        *float64 `json:"token_age_days,omitempty"`
	TokenAgeHours       *float64 `json:"token_age_hours,omitempty"`
	TokenDeploymentDate *string  `json:"token_deployment_date,omitempty"`
	MarketCapUSD        *float64 `json:"market_cap_usd,omitempty"`
	Liquidity           *float64 `json:"liquidity,omitempty"`
	PriceUSD            *float64 `json:"price_usd,omitempty"`
	PriceChange         *float64 `json:"price_change,omitempty"`
	FDV                 *float64 `json:"fdv,omitempty"`
	FDVMCRatio          *float64 `json:"fdv_mc_ratio,omitempty"`
	NofTraders          *int     `json:"nof_traders,omitempty"`
	NofBuyers           *int     `json:"nof_buyers,omitempty"`
	NofSellers          *int     `json:"nof_sellers,omitempty"`
	NofBuys             *int     `json:"nof_buys,omitempty"`
	NofSells            *int     `json:"nof_sells,omitempty"`
	BuyVolume           *float64 `json:"buy_volume,omitempty"`
	SellVolume          *float64 `json:"sell_volume,omitempty"`
	Volume              *float64 `json:"volume,omitempty"`
	Netflow             *float64 `json:"netflow,omitempty"`
	InflowFDVRatio      *float64 `json:"inflow_fdv_ratio,omitempty"`
	OutflowFDVRatio     *float64 `json:"outflow_fdv_ratio,omitempty"`
}

type TokenScreenerResponse struct {
	Data       []TokenScreenerResult `json:"data"`
	Pagination PaginationInfo        `json:"pagination"`
}

// TokenScreener discovers and screens tokens across multiple blockchains.
func (s *TokenGodModeService) TokenScreener(ctx context.Context, req *TokenScreenerRequest) (*TokenScreenerResponse, error) {
	resp := &TokenScreenerResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/token-screener", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TGMFlowIntelligenceTimeframe enumerates the supported flow-intelligence windows.
type TGMFlowIntelligenceTimeframe string

const (
	FlowTimeframe5M  TGMFlowIntelligenceTimeframe = "5m"
	FlowTimeframe1H  TGMFlowIntelligenceTimeframe = "1h"
	FlowTimeframe6H  TGMFlowIntelligenceTimeframe = "6h"
	FlowTimeframe12H TGMFlowIntelligenceTimeframe = "12h"
	FlowTimeframe1D  TGMFlowIntelligenceTimeframe = "1d"
	FlowTimeframe7D  TGMFlowIntelligenceTimeframe = "7d"
)

type TGMFlowIntelligenceFilters struct {
	PublicFigureNetFlowUSD  *NumericRangeFilter `json:"public_figure_net_flow_usd,omitempty"`
	PublicFigureAvgFlowUSD  *NumericRangeFilter `json:"public_figure_avg_flow_usd,omitempty"`
	PublicFigureWalletCount *IntegerRangeFilter `json:"public_figure_wallet_count,omitempty"`
	TopPnLNetFlowUSD        *NumericRangeFilter `json:"top_pnl_net_flow_usd,omitempty"`
	TopPnLAvgFlowUSD        *NumericRangeFilter `json:"top_pnl_avg_flow_usd,omitempty"`
	TopPnLWalletCount       *IntegerRangeFilter `json:"top_pnl_wallet_count,omitempty"`
	WhaleNetFlowUSD         *NumericRangeFilter `json:"whale_net_flow_usd,omitempty"`
	WhaleAvgFlowUSD         *NumericRangeFilter `json:"whale_avg_flow_usd,omitempty"`
	WhaleWalletCount        *IntegerRangeFilter `json:"whale_wallet_count,omitempty"`
	SmartTraderNetFlowUSD   *NumericRangeFilter `json:"smart_trader_net_flow_usd,omitempty"`
	SmartTraderAvgFlowUSD   *NumericRangeFilter `json:"smart_trader_avg_flow_usd,omitempty"`
	SmartTraderWalletCount  *IntegerRangeFilter `json:"smart_trader_wallet_count,omitempty"`
	ExchangeNetFlowUSD      *NumericRangeFilter `json:"exchange_net_flow_usd,omitempty"`
	ExchangeAvgFlowUSD      *NumericRangeFilter `json:"exchange_avg_flow_usd,omitempty"`
	ExchangeWalletCount     *IntegerRangeFilter `json:"exchange_wallet_count,omitempty"`
	FreshWalletsNetFlowUSD  *NumericRangeFilter `json:"fresh_wallets_net_flow_usd,omitempty"`
	FreshWalletsAvgFlowUSD  *NumericRangeFilter `json:"fresh_wallets_avg_flow_usd,omitempty"`
	FreshWalletsWalletCount *IntegerRangeFilter `json:"fresh_wallets_wallet_count,omitempty"`
}

type TGMFlowIntelligenceRequest struct {
	Chain        Chain                         `json:"chain"`
	TokenAddress string                        `json:"token_address"`
	Timeframe    *TGMFlowIntelligenceTimeframe `json:"timeframe,omitempty"`
	Filters      *TGMFlowIntelligenceFilters   `json:"filters,omitempty"`
}

type TGMFlowIntelligence struct {
	PublicFigureNetFlowUSD  *float64 `json:"public_figure_net_flow_usd,omitempty"`
	PublicFigureAvgFlowUSD  *float64 `json:"public_figure_avg_flow_usd,omitempty"`
	PublicFigureWalletCount *int     `json:"public_figure_wallet_count,omitempty"`
	TopPnLNetFlowUSD        *float64 `json:"top_pnl_net_flow_usd,omitempty"`
	TopPnLAvgFlowUSD        *float64 `json:"top_pnl_avg_flow_usd,omitempty"`
	TopPnLWalletCount       *int     `json:"top_pnl_wallet_count,omitempty"`
	WhaleNetFlowUSD         *float64 `json:"whale_net_flow_usd,omitempty"`
	WhaleAvgFlowUSD         *float64 `json:"whale_avg_flow_usd,omitempty"`
	WhaleWalletCount        *int     `json:"whale_wallet_count,omitempty"`
	SmartTraderNetFlowUSD   *float64 `json:"smart_trader_net_flow_usd,omitempty"`
	SmartTraderAvgFlowUSD   *float64 `json:"smart_trader_avg_flow_usd,omitempty"`
	SmartTraderWalletCount  *int     `json:"smart_trader_wallet_count,omitempty"`
	ExchangeNetFlowUSD      *float64 `json:"exchange_net_flow_usd,omitempty"`
	ExchangeAvgFlowUSD      *float64 `json:"exchange_avg_flow_usd,omitempty"`
	ExchangeWalletCount     *int     `json:"exchange_wallet_count,omitempty"`
	FreshWalletsNetFlowUSD  *float64 `json:"fresh_wallets_net_flow_usd,omitempty"`
	FreshWalletsAvgFlowUSD  *float64 `json:"fresh_wallets_avg_flow_usd,omitempty"`
	FreshWalletsWalletCount *int     `json:"fresh_wallets_wallet_count,omitempty"`
}

type TGMFlowIntelligenceResponse struct {
	Data     []TGMFlowIntelligence `json:"data"`
	Warnings []string              `json:"warnings,omitempty"`
}

// FlowIntelligence returns token flow summaries broken down by holder segments.
func (s *TokenGodModeService) FlowIntelligence(ctx context.Context, req *TGMFlowIntelligenceRequest) (*TGMFlowIntelligenceResponse, error) {
	resp := &TGMFlowIntelligenceResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/tgm/flow-intelligence", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// TGMWhoBoughtSoldSortField enumerates the sortable fields.
type TGMWhoBoughtSoldSortField string

const (
	WhoBoughtSoldSortBoughtVolumeUSD   TGMWhoBoughtSoldSortField = "bought_volume_usd"
	WhoBoughtSoldSortSoldVolumeUSD     TGMWhoBoughtSoldSortField = "sold_volume_usd"
	WhoBoughtSoldSortTokenTradeVolume  TGMWhoBoughtSoldSortField = "token_trade_volume"
	WhoBoughtSoldSortTradeVolumeUSD    TGMWhoBoughtSoldSortField = "trade_volume_usd"
	WhoBoughtSoldSortBoughtTokenVolume TGMWhoBoughtSoldSortField = "bought_token_volume"
	WhoBoughtSoldSortSoldTokenVolume   TGMWhoBoughtSoldSortField = "sold_token_volume"
)

type TGMWhoBoughtSoldFilters struct {
	IncludeSmartMoneyLabels []SmartMoneyLabel   `json:"include_smart_money_labels,omitempty"`
	ExcludeSmartMoneyLabels []SmartMoneyLabel   `json:"exclude_smart_money_labels,omitempty"`
	Address                 interface{}         `json:"address,omitempty"`
	AddressLabel            interface{}         `json:"address_label,omitempty"`
	BoughtTokenVolume       *NumericRangeFilter `json:"bought_token_volume,omitempty"`
	SoldTokenVolume         *NumericRangeFilter `json:"sold_token_volume,omitempty"`
	TokenTradeVolume        *NumericRangeFilter `json:"token_trade_volume,omitempty"`
	BoughtVolumeUSD         *NumericRangeFilter `json:"bought_volume_usd,omitempty"`
	SoldVolumeUSD           *NumericRangeFilter `json:"sold_volume_usd,omitempty"`
	TradeVolumeUSD          *NumericRangeFilter `json:"trade_volume_usd,omitempty"`
}

type TGMWhoBoughtSoldRequest struct {
	Chain        Chain                    `json:"chain"`
	TokenAddress string                   `json:"token_address"`
	BuyOrSell    *BuyOrSell               `json:"buy_or_sell,omitempty"`
	Date         DateRange                `json:"date"`
	Pagination   *PaginationRequest       `json:"pagination,omitempty"`
	Filters      *TGMWhoBoughtSoldFilters `json:"filters,omitempty"`
	OrderBy      []SortOrder              `json:"order_by,omitempty"`
}

type TGMWhoBoughtSold struct {
	Address           string   `json:"address"`
	AddressLabel      *string  `json:"address_label,omitempty"`
	BoughtTokenVolume *float64 `json:"bought_token_volume,omitempty"`
	SoldTokenVolume   *float64 `json:"sold_token_volume,omitempty"`
	TokenTradeVolume  *float64 `json:"token_trade_volume,omitempty"`
	BoughtVolumeUSD   *float64 `json:"bought_volume_usd,omitempty"`
	SoldVolumeUSD     *float64 `json:"sold_volume_usd,omitempty"`
	TradeVolumeUSD    *float64 `json:"trade_volume_usd,omitempty"`
}

type TGMWhoBoughtSoldResponse struct {
	Data       []TGMWhoBoughtSold `json:"data"`
	Pagination PaginationInfo     `json:"pagination"`
}

// WhoBoughtSold returns aggregated summaries of token buyers or sellers.
func (s *TokenGodModeService) WhoBoughtSold(ctx context.Context, req *TGMWhoBoughtSoldRequest) (*TGMWhoBoughtSoldResponse, error) {
	resp := &TGMWhoBoughtSoldResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/tgm/who-bought-sold", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
