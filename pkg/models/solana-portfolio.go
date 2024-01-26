package models

import (
	"time"

	"gorm.io/gorm"
)

// SolanaAssetsMoralisV1 represents the solana_assets_moralis_v1 table.
type SolanaAssetsMoralisV1 struct {
	SolanaAssetID    uint      `gorm:"primaryKey;autoIncrement"`
	WalletID         int       `gorm:"not null"`
	Lamports         string    `gorm:"type:varchar(255)"`
	Solana           string    `gorm:"type:varchar(255)"`
	TotalTokensCount int       `gorm:"type:integer"`
	TotalNftsCount   int       `gorm:"type:integer"`
	LastUpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Token represents the tokens table.
type Token struct {
	TokenID                int       `gorm:"primary_key"`
	SolanaAssetID          int       `gorm:"not null"`
	AssociatedTokenAddress string    `gorm:"type:varchar(255)"`
	Mint                   string    `gorm:"type:varchar(255):"`
	AmountRaw              string    `gorm:"type:varchar(255)"`
	Amount                 string    `gorm:"type:varchar(255)"`
	Decimals               string    `gorm:"type:varchar(255)"`
	Name                   string    `gorm:"type:varchar(255)"`
	Symbol                 string    `gorm:"type:varchar(50)"`
	UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// NFT represents the nfts table.
type NFT struct {
	NFTID                  int       `gorm:"primary_key"`
	SolanaAssetID          int       `gorm:"not null"`
	AssociatedTokenAddress string    `gorm:"type:varchar(255)"`
	Mint                   string    `gorm:"type:varchar(255)"`
	AmountRaw              string    `gorm:"type:varchar(255)"`
	Decimals               string    `gorm:"type:varchar(255)"`
	Name                   string    `gorm:"type:varchar(255)"`
	Symbol                 string    `gorm:"type:varchar(50)"`
	UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// SaveSolanaData saves a SolanaAssetsMoralisV1 record along with its associated tokens and NFTs.
func SaveSolanaData(tx *gorm.DB, solanaAsset *SolanaAssetsMoralisV1, tokens []Token, nfts []NFT) error {
	// Save the SolanaAssetsMoralisV1 record.
	if result := tx.Create(solanaAsset); result.Error != nil {
		return result.Error
	}

	// Set SolanaAssetID for each Token and NFT.
	for i := range tokens {
		tokens[i].SolanaAssetID = int(solanaAsset.SolanaAssetID)
	}
	for i := range nfts {
		nfts[i].SolanaAssetID = int(solanaAsset.SolanaAssetID)
	}

	// Save each Token
	for _, token := range tokens {
		if result := tx.Create(&token); result.Error != nil {
			return result.Error
		}
	}

	// Save each NFT
	for _, nft := range nfts {
		if result := tx.Create(&nft); result.Error != nil {
			return result.Error
		}
	}

	return nil
}
