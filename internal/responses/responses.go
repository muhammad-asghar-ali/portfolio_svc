package responses

import (
	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/shared/utils"
)

type (
	PortfolioResponse struct {
		AssetSymbol         string            `json:"asset_symbol"`
		UnitPrice           float64           `json:"unit_price"`
		Quantity            float64           `json:"quantity"`
		PortfolioPercentage float64           `json:"portfolio_percentage"`
		TotalPrice          float64           `json:"total_price"`
		ChainsInfo          []*ChainsResponse `json:"chains_info"`
	}

	ChainsResponse struct {
		WalletID        int     `json:"wallet_id"`
		AssetSymbol     string  `json:"asset_symbol"`
		AssetID         string  `json:"asset_id"`
		Chain           string  `json:"chain"`
		UnitPrice       float64 `json:"unit_price"`
		Quantity        float64 `json:"quantity"`
		AssetPercentage float64 `json:"asset_percentage"`
		TotalPrice      float64 `json:"total_price"`
		IsVerified      bool    `json:"is_verified"`
	}
)

// Handle bitcoin related responses
//
// BTCResponse updates the Response struct with Bitcoin-related information from the wallet.
func (r *ChainsResponse) BTCResponse(wallet *models.GlobalWallet) {
	r.WalletID = wallet.WalletID
	r.AssetSymbol = "token"
	r.Chain = wallet.BlockchainType
	r.UnitPrice = wallet.BitcoinBtcComV1.CoingeckoPriceFeed.Price
	r.Quantity = wallet.BitcoinBtcComV1.BitcoinAddressInfo.Balance
	r.TotalPrice = r.UnitPrice * r.Quantity
	r.IsVerified = true
}

// Handle solana related responses
//
// SolanaResponse updates the Response struct with Solana-token-related information from the wallet and token.
func (r *ChainsResponse) SolanaTokenResponse(wallet *models.GlobalWallet, token *models.Token) {
	quantity, _ := utils.StrToFloat64(wallet.SolanaAssetsMoralisV1.Solana)

	r.WalletID = wallet.WalletID
	r.AssetSymbol = token.Name
	r.Chain = wallet.BlockchainType
	r.UnitPrice = wallet.SolanaAssetsMoralisV1.CoingeckoPriceFeed.Price
	r.Quantity = quantity
	r.TotalPrice = r.UnitPrice * r.Quantity
	r.IsVerified = true
}

// NOTE: ignore for now
//
// SolanaResponse updates the Response struct with Solana-nft-related information from the wallet and nft.
func (r *ChainsResponse) SolanaNFTResponse(wallet *models.GlobalWallet, nft *models.NFT) {
	quantity, _ := utils.StrToFloat64(wallet.SolanaAssetsMoralisV1.Solana)

	r.WalletID = wallet.WalletID
	r.AssetSymbol = nft.Name
	r.Chain = wallet.BlockchainType
	r.UnitPrice = wallet.SolanaAssetsMoralisV1.CoingeckoPriceFeed.Price
	r.Quantity = quantity
	r.TotalPrice = r.UnitPrice * r.Quantity
	r.IsVerified = true
}

// Handle debank related responses
//
// DebankTokenResponse updates the Response struct with Debank-token-related information from the wallet and token list.
func (r *ChainsResponse) DebankTokenResponse(wallet *models.GlobalWallet, token *models.TokenList) {
	r.WalletID = wallet.WalletID
	r.AssetSymbol = token.Symbol
	r.AssetID = token.ID
	r.Chain = token.Chain
	r.UnitPrice = token.Price
	r.Quantity = token.Amount
	r.TotalPrice = r.UnitPrice * r.Quantity
	r.IsVerified = token.IsVerified
}

// NOTE: ignore for now
//
// DebankNFTResponse updates the Response struct with Debank-nft-related information from the wallet and nft list.
func (r *ChainsResponse) DebankNFTResponse(wallet *models.GlobalWallet, nft *models.NFTList) {
	r.WalletID = wallet.WalletID
	r.AssetSymbol = nft.Name
	r.Chain = nft.Chain
	r.UnitPrice = nft.USDPrice
	r.Quantity = float64(nft.Amount)
	r.TotalPrice = r.UnitPrice * r.Quantity
	r.IsVerified = true
}

func (p *PortfolioResponse) BitcoinPortfolioResponse(wallet *models.GlobalWallet) {
	p.AssetSymbol = "btc"
	p.UnitPrice = wallet.BitcoinBtcComV1.CoingeckoPriceFeed.Price
	p.Quantity = wallet.BitcoinBtcComV1.BitcoinAddressInfo.Balance
	p.TotalPrice = p.UnitPrice * p.Quantity
}
