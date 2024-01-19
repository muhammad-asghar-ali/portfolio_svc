package models

import (
	"time"

	"gorm.io/gorm"
)

// SolanaAssetsMoralisV1 represents the solana_assets_moralis_v1 table.
type SolanaAssetsMoralisV1 struct {
	SolanaAssetID    int     `gorm:"primary_key"`
	WalletID         int     `gorm:"not null"`
	NativeBalance    float64 // DECIMAL in SQL
	TotalTokensCount int
	TotalNFTsCount   int
	LastUpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Token represents the tokens table.
type Token struct {
	TokenID                int `gorm:"primary_key"`
	SolanaAssetID          int `gorm:"not null"`
	AssociatedTokenAddress string
	Mint                   string `gorm:"unique"`
	AmountRaw              int64
	Amount                 float64 // DECIMAL in SQL
	Decimals               int
	UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// NFT represents the nfts table.
type NFT struct {
	NFTID                  int `gorm:"primary_key"`
	SolanaAssetID          int `gorm:"not null"`
	AssociatedTokenAddress string
	Mint                   string
	AmountRaw              int64
	Decimals               int
	UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TokenPrice represents the token_prices table.
type TokenPrice struct {
	PriceID          int     `gorm:"primary_key"`
	TokenMint        string  `gorm:"not null"`
	USDPrice         float64 // DECIMAL in SQL
	ExchangeName     string
	ExchangeAddress  string
	NativePriceValue int64
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func SaveSolanaData(tx *gorm.DB, solanaAsset SolanaAssetsMoralisV1, tokens []Token, nfts []NFT, tokenPrices []TokenPrice) error {
	if result := tx.Create(&solanaAsset); result.Error != nil {
		return result.Error
	}

	for _, token := range tokens {
		if result := tx.Create(&token); result.Error != nil {
			return result.Error
		}
	}

	for _, nft := range nfts {
		if result := tx.Create(&nft); result.Error != nil {
			return result.Error
		}
	}

	for _, tokenPrice := range tokenPrices {
		if result := tx.Create(&tokenPrice); result.Error != nil {
			return result.Error
		}
	}

	return nil
}
