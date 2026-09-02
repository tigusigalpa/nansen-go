package nansen

import "context"

// PortfolioService provides access to portfolio-level endpoints.
type PortfolioService struct {
	client *Client
}

// PortfolioDefiHoldingsRequest specifies a wallet portfolio query.
type PortfolioDefiHoldingsRequest struct {
	WalletAddress string `json:"wallet_address"`
}

// HoldingsSummary contains aggregate values for a wallet's DeFi holdings.
type HoldingsSummary struct {
	TotalValueUSD   float64 `json:"total_value_usd"`
	TotalAssetsUSD  float64 `json:"total_assets_usd"`
	TotalDebtsUSD   float64 `json:"total_debts_usd"`
	TotalRewardsUSD float64 `json:"total_rewards_usd"`
	TokenCount      int     `json:"token_count"`
	ProtocolCount   int     `json:"protocol_count"`
}

// ProtocolToken describes a token position within a DeFi protocol.
type ProtocolToken struct {
	Address      *string      `json:"address,omitempty"`
	Symbol       *string      `json:"symbol,omitempty"`
	Amount       *float64     `json:"amount,omitempty"`
	ValueUSD     *float64     `json:"value_usd,omitempty"`
	PositionType PositionType `json:"position_type"`
}

// ProtocolHolding describes all positions in a DeFi protocol.
type ProtocolHolding struct {
	ProtocolName    string          `json:"protocol_name"`
	Chain           string          `json:"chain"`
	TotalValueUSD   float64         `json:"total_value_usd"`
	TotalAssetsUSD  float64         `json:"total_assets_usd"`
	TotalDebtsUSD   float64         `json:"total_debts_usd"`
	TotalRewardsUSD float64         `json:"total_rewards_usd"`
	Tokens          []ProtocolToken `json:"tokens"`
}

// PortfolioDefiHoldingsResponse contains a wallet's DeFi holdings.
type PortfolioDefiHoldingsResponse struct {
	Summary   HoldingsSummary   `json:"summary"`
	Protocols []ProtocolHolding `json:"protocols"`
}

// DeFiHoldings returns simplified DeFi holdings for a wallet address.
func (s *PortfolioService) DeFiHoldings(ctx context.Context, req *PortfolioDefiHoldingsRequest) (*PortfolioDefiHoldingsResponse, error) {
	resp := &PortfolioDefiHoldingsResponse{}
	if err := s.client.doRequest(ctx, "POST", "/api/v1/portfolio/defi-holdings", req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
