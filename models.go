package nansen

// StringPtr returns a pointer to s.
func StringPtr(s string) *string { return &s }

// IntPtr returns a pointer to i.
func IntPtr(i int) *int { return &i }

// BoolPtr returns a pointer to b.
func BoolPtr(b bool) *bool { return &b }

// Float64Ptr returns a pointer to f.
func Float64Ptr(f float64) *float64 { return &f }

// Chain represents a blockchain network supported across Nansen endpoints.
type Chain string

// ChainAll through ChainTron identify blockchain networks supported by Nansen.
const (
	ChainAll         Chain = "all"
	ChainArbitrum    Chain = "arbitrum"
	ChainAvalanche   Chain = "avalanche"
	ChainBase        Chain = "base"
	ChainBitcoin     Chain = "bitcoin"
	ChainBNB         Chain = "bnb"
	ChainCitrea      Chain = "citrea"
	ChainEthereum    Chain = "ethereum"
	ChainHyperEVM    Chain = "hyperevm"
	ChainHyperliquid Chain = "hyperliquid"
	ChainInjective   Chain = "injective"
	ChainIotaEVM     Chain = "iotaevm"
	ChainLinea       Chain = "linea"
	ChainMantle      Chain = "mantle"
	ChainMantra      Chain = "mantra"
	ChainMonad       Chain = "monad"
	ChainNear        Chain = "near"
	ChainOptimism    Chain = "optimism"
	ChainPlasma      Chain = "plasma"
	ChainPolygon     Chain = "polygon"
	ChainRonin       Chain = "ronin"
	ChainScroll      Chain = "scroll"
	ChainSei         Chain = "sei"
	ChainSolana      Chain = "solana"
	ChainSonic       Chain = "sonic"
	ChainStarknet    Chain = "starknet"
	ChainSui         Chain = "sui"
	ChainTon         Chain = "ton"
	ChainTron        Chain = "tron"
)

// SortDirection controls ordering of result sets.
type SortDirection string

// SortAsc and SortDesc define result ordering directions.
const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

// SortOrder describes a single field/direction pair.
type SortOrder struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

// PaginationRequest controls paging. Pointer fields are omitted when nil so
// that the API applies its own defaults.
type PaginationRequest struct {
	Page    *int `json:"page,omitempty"`
	PerPage *int `json:"per_page,omitempty"`
}

// PaginationInfo is returned with paginated responses.
type PaginationInfo struct {
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	IsLastPage bool `json:"is_last_page"`
}

// DateRange is an explicit ISO-8601 interval.
type DateRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// NumericRangeFilter filters numeric values by inclusive min/max bounds.
type NumericRangeFilter struct {
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
}

// IntegerRangeFilter filters integer values by inclusive min/max bounds.
type IntegerRangeFilter struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

// SmartMoneyLabel filters smart-money cohorts.
type SmartMoneyLabel string

// LabelFund through LabelSmartHLPerpsTrader identify smart-money cohorts.
const (
	LabelFund               SmartMoneyLabel = "Fund"
	LabelSmartTrader        SmartMoneyLabel = "Smart Trader"
	Label30DSmartTrader     SmartMoneyLabel = "30D Smart Trader"
	Label90DSmartTrader     SmartMoneyLabel = "90D Smart Trader"
	Label180DSmartTrader    SmartMoneyLabel = "180D Smart Trader"
	LabelSmartHLPerpsTrader SmartMoneyLabel = "Smart HL Perps Trader"
)

// TraderType filters the cohort of traders for the token screener.
type TraderType string

// TraderAll through TraderPredictedWinner identify token-screener trader cohorts.
const (
	TraderAll                   TraderType = "all"
	TraderSmartMoney            TraderType = "sm"
	TraderWhale                 TraderType = "whale"
	TraderPublicFigure          TraderType = "public_figure"
	TraderTrending              TraderType = "trending"
	TraderConsistentPerpsWinner TraderType = "consistent_perps_winner"
	TraderHighWinrateHLPerps    TraderType = "high_winrate_hl_perps_trader"
	TraderPredictedWinner       TraderType = "predicted_winner"
)

// BuyOrSell controls the trade direction for who-bought-sold queries.
type BuyOrSell string

// Buy and Sell identify the direction of a trade.
const (
	Buy  BuyOrSell = "BUY"
	Sell BuyOrSell = "SELL"
)

// PositionType describes the role of a DeFi token position.
type PositionType string

// PositionDeposit through PositionMixed identify DeFi position roles.
const (
	PositionDeposit PositionType = "deposit"
	PositionStake   PositionType = "stake"
	PositionBorrow  PositionType = "borrow"
	PositionMixed   PositionType = "mixed"
)
