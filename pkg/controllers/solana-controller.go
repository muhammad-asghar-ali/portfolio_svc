package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/0xbase-Corp/portfolio_svc/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SolanaController handles requests for Solana portfolio information.
func SolanaController(c *gin.Context, db *gorm.DB) {
	log.Println("SolanaController invoked")
	// Extract the Solana address from the request parameter.
	solAddress := c.Param("sol-address")
	log.Println("for address:", solAddress)
	// Extract the Moralis API key from the request header.
	moralisAccessKey := c.GetHeader("x-api-key")

	// Prepare the Moralis API request URL.
	url := "https://solana-gateway.moralis.io/account/mainnet/" + solAddress + "/portfolio"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add the Moralis API key to the request header.
	req.Header.Add("x-api-key", moralisAccessKey)

	// Create an HTTP client and execute the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Define a struct to match the JSON response structure from the Moralis API.
	var response struct {
		Tokens        []models.Token `json:"tokens"`
		NFTs          []models.NFT   `json:"nfts"`
		NativeBalance struct {
			Lamports string `json:"lamports"`
			Solana   string `json:"solana"`
		} `json:"nativeBalance"`
	}

	// Parse the JSON response into the defined struct.
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("Error parsing JSON response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON response: " + err.Error()})
		return
	}
	log.Println("Received response from Moralis API")
	// Check if a wallet with the given Solana address exists in the global_wallets table.
	var wallet models.GlobalWallet
	err = db.Where("wallet_address = ?", solAddress).First(&wallet).Error
	if err != nil {
		// If the wallet doesn't exist, create a new one.
		if err == gorm.ErrRecordNotFound {
			wallet = models.GlobalWallet{
				PortfolioID:    0, // or nil, depending on the data type of PortfolioID
				WalletAddress:  solAddress,
				BlockchainType: "Solana",
			}
			if err := db.Create(&wallet).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet: " + err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error: " + err.Error()})
			return
		}
	}

	// Prepare a SolanaAssetsMoralisV1 object with the response data.
	solanaAsset := models.SolanaAssetsMoralisV1{
		WalletID:         wallet.WalletID,
		Lamports:         response.NativeBalance.Lamports,
		Solana:           response.NativeBalance.Solana,
		TotalTokensCount: len(response.Tokens),
		TotalNftsCount:   len(response.NFTs),
		LastUpdatedAt:    time.Now(),
	}

	// Start a new database transaction.
	tx := db.Begin()

	// Update or create the SolanaAssetsMoralisV1 record.
	if err := updateOrCreateSolanaAsset(tx, &solanaAsset); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Solana assets: " + err.Error()})
		return
	}
	log.Println("Attempting to save data to the database")
	// Save the Solana asset data along with the associated tokens and NFTs.
	// The SaveSolanaData function will handle the creation of all these records.
	if err := models.SaveSolanaData(tx, &solanaAsset, response.Tokens, response.NFTs); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to the database: " + err.Error()})
		return
	}

	// Commit the transaction if everything is successful.
	if err := tx.Commit().Error; err != nil {
		log.Println("Error committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// Send a success response.
	c.JSON(http.StatusOK, gin.H{"message": "Data saved successfully"})

}

// updateOrCreateSolanaAsset updates an existing SolanaAssetsMoralisV1 record or creates a new one.
func updateOrCreateSolanaAsset(tx *gorm.DB, asset *models.SolanaAssetsMoralisV1) error {
	// Check if the asset record already exists for the given SolanaAssetID.
	var existingAsset models.SolanaAssetsMoralisV1
	result := tx.Where("solana_asset_id = ?", asset.SolanaAssetID).First(&existingAsset)

	// If the record does not exist, create a new one.
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create a new record without specifying SolanaAssetID.
			newAsset := models.SolanaAssetsMoralisV1{
				WalletID:         asset.WalletID,
				Lamports:         asset.Lamports,
				Solana:           asset.Solana,
				TotalTokensCount: asset.TotalTokensCount,
				TotalNftsCount:   asset.TotalNftsCount,
				LastUpdatedAt:    asset.LastUpdatedAt,
			}
			if err := tx.Create(&newAsset).Error; err != nil {
				return err
			}
			return nil
		}
		// Return any other error encountered during the query.
		return result.Error
	}

	// If the record exists, update it with the new information.
	return tx.Model(&existingAsset).Updates(asset).Error
}
