package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/pkg/models"
	"github.com/0xbase-Corp/portfolio_svc/pkg/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@BasePath	/api/v1

// PingExample godoc
//
//	@Summary		ping example
//	@Schemes		models.BtcChainAPI
//	@Description	do ping
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string} Helloworld
//	@Router			/portfolio/btc/:btc-address [get]
func BitcoinController(c *gin.Context, db *gorm.DB) {
	btcAddress := c.Param("btc-address")

	if btcAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BTC address"})
		return
	}

	btc, err := fetchBitcoinData(btcAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bitcoinAddressInfo := mapBtcResponseToBitcoinAddressInfoTable(&btc.Data)

	if err := db.Create(&bitcoinAddressInfo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bitcoinAddressInfo})
}

func fetchBitcoinData(btcAddress string) (*types.BtcChainAPI, error) {
	// TODO: move URL to env variable
	url := "https://chain.api.btc.com/v3/address/" + btcAddress

	resp, _ := http.Get(url)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to get the btc information against this address")
	}
	defer resp.Body.Close()

	btc := types.BtcChainAPI{}
	if err := json.NewDecoder(resp.Body).Decode(&btc); err != nil {
		return nil, err
	}

	return &btc, nil
}

func mapBtcResponseToBitcoinAddressInfoTable(btc *types.ChainData) models.BitcoinAddressInfo {
	return models.BitcoinAddressInfo{
		Address:             btc.Address,
		Received:            btc.Received,
		Sent:                btc.Sent,
		Balance:             btc.Balance,
		TxCount:             btc.TxCount,
		UnconfirmedTxCount:  btc.UnconfirmedTxCount,
		UnconfirmedReceived: btc.UnconfirmedReceived,
		UnconfirmedSent:     btc.UnconfirmedSent,
		UnspentTxCount:      btc.UnspentTxCount,
		FirstTx:             btc.FirstTx,
		LastTx:              btc.LastTx,
	}
}
