package controllers

import (
	er "errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/responses"
	"github.com/0xbase-Corp/portfolio_svc/providers"
	"github.com/0xbase-Corp/portfolio_svc/providers/bitcoin"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/0xbase-Corp/portfolio_svc/shared/utils"
)

//	@BasePath	/api/v1

// BitcoinController godoc
//
// @Summary      Fetch Bitcoin Wallet Information
// @Description  Retrieves information for a given Bitcoin address using the BTC.com API.
// @Tags         bitcoin
// @Accept       json
// @Produce      json
// @Param        addresses  query      array  true  "Bitcoin Addresses" Format(string)
// @Success      200 {object} []responses.PortfolioResponse
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/btc [get]
func BitcoinController(c *gin.Context, db *gorm.DB, apiClient providers.APIClient) {
	addresses := c.Query("addresses")
	btcAddresses := strings.Split(addresses, ",")

	btcAddresses = utils.UniqueAddress(btcAddresses)

	if len(btcAddresses) == 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty btc addresses"))
		return
	}

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	ch := make(chan *models.GlobalWallet, len(btcAddresses)) // 1 specifies the buffer size of the channel
	errorCh := make(chan error, len(btcAddresses))

	for _, btcAddress := range btcAddresses {
		wg.Add(1)

		go fetchAndSaveBtc(db, apiClient, btcAddress, ch, wg, mutex, errorCh)
	}

	// Use a goroutine to close the channel after all goroutines have finished
	go func() {
		wg.Wait()
		close(ch)
		close(errorCh)
	}()

	errs := errors.HandleChannelErrors(errorCh, c)
	if len(errs) > 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError(strings.Join(errs, "; ")))
		return
	}

	// Collect all results from the channel and process to genernic response for btc
	responses := processBtcResponses(ch)

	c.JSON(http.StatusOK, responses)
}

//	@BasePath	/api/v1
//
// GetBtcDataController godoc
// @Summary      Get BTC portfolio for a wallet
// @Description  Retrieve BTC portfolio details, including BitcoinAddressInfo, for a specific wallet.
// @Tags         bitcoin
// @Accept       json
// @Produce      json
// @Param        wallet_id path int true "Wallet ID" Format(int)
// @Param        offset query int false "Pagination offset" Format(int)
// @Param        limit query int false "Pagination limit" Format(int)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/btc-wallet/{wallet_id} [get]
func GetBtcDataController(c *gin.Context, db *gorm.DB) {
	wallet := models.GlobalWallet{}
	walletID, err := strconv.Atoi(c.Param("wallet-id"))

	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("invalid wallet id"))
		return
	}

	// Parse optional query parameters
	page, err := strconv.Atoi(c.DefaultQuery("offset", "1"))
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("invalid offset"))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("invalid limit"))
		return
	}

	// Define a reusable function to apply offset and limit for preload on BitcoinAddressInfo
	prlimit := func(query *gorm.DB) *gorm.DB {
		return query.Offset((page - 1) * limit).Limit(limit)
	}

	err = db.
		Preload("BitcoinBtcComV1.BitcoinAddressInfo", prlimit).
		Preload("BitcoinBtcComV1").
		Where("wallet_id = ?", walletID).
		First(&wallet).Error

	if err != nil {
		errors.HandleHttpError(c, errors.NewNotFoundError("wallet not found"))
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// helper to save data
// TODO: write a common interface in provider which saves the data in database
func saveBtc(db *gorm.DB, btcAddress string, apiResponse bitcoin.BtcApiResponse) (*models.GlobalWallet, error) {
	// Begin a new transaction
	tx := db.Begin()

	wallet, err := models.GetOrCreateWallet(tx, btcAddress, utils.Bitcoin)
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Assuming wallet is the GlobalWallet record found or created
	walletID := wallet.WalletID

	// Initialize btcComV1 and set the WalletID
	btcComV1 := models.BitcoinBtcComV1{}
	btcComV1.WalletID = uint(walletID)

	// Proceed to call SaveBitcoinData with the updated BitcoinAddressInfo and btcComV1
	if err := models.SaveBitcoinData(tx, &apiResponse.Data, &btcComV1); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// save the bitcoin price feed
	// for now hard code the USD -> TODO: change
	if err := HandleCoingeckoPrice(tx, utils.Bitcoin, "usd"); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return &models.GlobalWallet{}, err
	}

	// get the data
	walletResponse, err := models.GetGlobalWalletWithBitcoinInfo(db, btcAddress)
	if err != nil {
		return &models.GlobalWallet{}, err
	}

	return walletResponse, nil
}

// Fetch and save the data for one address
// TODO: write a common interface in provider which saves the data in database
func fetchAndSaveBtc(db *gorm.DB, apiClient providers.APIClient, address string, ch chan<- *models.GlobalWallet, wg *sync.WaitGroup, mutex *sync.Mutex, errorCh chan<- error) {
	defer wg.Done()

	body, err := apiClient.FetchData(address)
	if err != nil {
		errorCh <- err
		return
	}

	// ignore the Address which returns error or empty response
	resp := bitcoin.BtcApiResponse{}
	if err := utils.DecodeJSONResponse(body, &resp); err != nil {
		errorCh <- err
		return
	}

	if resp.Status == "fail" {
		errorCh <- er.New("API error: " + resp.Message)
		return
	}

	walletResponse, err := saveBtc(db, address, resp)
	if err != nil {
		errorCh <- err
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}

// processBtcResponses processes wallet responses and returns a slice of btc responses.
func processBtcResponses(ch <-chan *models.GlobalWallet) []*responses.PortfolioResponse {
	btcResponses := make([]*responses.PortfolioResponse, 0)
	for walletResponse := range ch {
		if walletResponse == nil || walletResponse.BitcoinBtcComV1 == nil {
			continue
		}

		btcResponse := &responses.PortfolioResponse{}
		btcResponse.BitcoinPortfolioResponse(walletResponse)
		btcResponses = append(btcResponses, btcResponse)
	}

	btcResponses = responses.CalculatePortfolioResponse(btcResponses)

	return btcResponses
}
