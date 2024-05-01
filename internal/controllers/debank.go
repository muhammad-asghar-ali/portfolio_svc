package controllers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/responses"
	"github.com/0xbase-Corp/portfolio_svc/providers"
	"github.com/0xbase-Corp/portfolio_svc/providers/debank"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/0xbase-Corp/portfolio_svc/shared/utils"
)

//	@BasePath	/api/v1

// DebankController godoc
//
// @Summary      Fetch Debank Wallet Information
// @Description  Retrieves information for a given Debank address using the BTC.com API.
// @Tags         debank
// @Accept       json
// @Produce      json
// @Param        addresses  query      array  true  "Debank Address" Format(string)
// @Success      200 {object} []responses.PortfolioResponse
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/debank [get]
func DebankController(c *gin.Context, db *gorm.DB, apiClient providers.APIClient) {
	addresses := c.Query("addresses")
	debankAddresses := strings.Split(addresses, ",")

	debankAddresses = utils.UniqueAddress(debankAddresses)

	if len(debankAddresses) == 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty btc addresses"))
		return
	}

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	ch := make(chan *models.GlobalWallet, len(debankAddresses)) // len(debankAddresses) specifies the buffer size of the channel
	errorCh := make(chan error, len(debankAddresses))

	for _, btcAddress := range debankAddresses {
		wg.Add(1)

		go fetchAndSaveDebank(db, apiClient, btcAddress, ch, wg, mutex, errorCh)
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

	// Collect all results from the channel and process to genernic response for debank
	responses := processDebankResponses(ch)

	c.JSON(http.StatusOK, responses)
}

// Fetch and save the data for one address
// TODO: write a common interface in provider which saves the data in database
func fetchAndSaveDebank(db *gorm.DB, apiClient providers.APIClient, address string, ch chan<- *models.GlobalWallet, wg *sync.WaitGroup, mutex *sync.Mutex, errorCh chan<- error) {
	defer wg.Done()

	wallet, _ := models.GetWallet(db, address)

	now, _ := utils.GetDBTime()

	if wallet != nil {
		// Calculate the duration between the timestamps
		duration := now.Sub(wallet.LastUpdatedAt)

		if duration.Hours() > 24 {
			updateWalletAndSend(db, wallet, address, ch, mutex, errorCh)
		} else {
			retrieveWalletAndSend(db, address, ch, mutex, errorCh)
		}
	} else {
		fetchFromAPIAndSave(db, apiClient, address, ch, mutex, errorCh)
	}
}

// fetch data from database and send in response
func retrieveWalletAndSend(db *gorm.DB, address string, ch chan<- *models.GlobalWallet, mutex *sync.Mutex, errorCh chan<- error) {
	walletResponse, err := models.GetGlobalWalletWithEvmDebankInfo(db, address)
	if err != nil {
		errorCh <- err
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}

// fetch data from api and save data to database
func fetchFromAPIAndSave(db *gorm.DB, apiClient providers.APIClient, address string, ch chan<- *models.GlobalWallet, mutex *sync.Mutex, errorCh chan<- error) {
	body, err := apiClient.FetchData(address)
	if err != nil {
		errorCh <- err
		return
	}

	// ignore the Address which returns error or empty response
	resp := debank.EvmDebankTotalBalanceApiResponse{}
	if err := utils.DecodeJSONResponse(body, &resp); err != nil {
		errorCh <- err
		return
	}

	walletResponse, err := saveDebank(db, address, resp)
	if err != nil {
		errorCh <- err
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}

// update the fetch the data from database
func updateWalletAndSend(db *gorm.DB, wallet *models.GlobalWallet, address string, ch chan<- *models.GlobalWallet, mutex *sync.Mutex, errorCh chan<- error) {
	walletResponse, err := models.GetGlobalWalletWithEvmDebankInfo(db, address)
	if err != nil {
		errorCh <- err
		return
	}

	if _, err := models.UpdateWalletLastUpdateAt(db, wallet); err != nil {
		errorCh <- err
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}

// helper to save data
// TODO: write a common interface in provider which saves the data in database
func saveDebank(db *gorm.DB, address string, apiResponse debank.EvmDebankTotalBalanceApiResponse) (*models.GlobalWallet, error) {
	// Begin a new transaction
	tx := db.Begin()

	wallet, err := models.GetOrCreateWallet(tx, address, utils.Debank)
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Initialize EvmAssetsDebankV1 and set the WalletID
	evmAssetsDebankV1 := models.EvmAssetsDebankV1{
		WalletID:      wallet.WalletID,
		TotalUsdValue: apiResponse.TotalUsdValue,
	}

	// Save EvmAssetsDebankV1
	if err = models.CreateOrUpdateEvmAssetsDebankV1(tx, &evmAssetsDebankV1); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Save Chain
	if err = models.SaveChainDetails(tx, wallet.WalletID, apiResponse.ChainList); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Save token list
	err = models.SaveTokenListByEvmAssetsDebankV1ID(tx, evmAssetsDebankV1.EvmAssetID, apiResponse.TokensList)
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Save nft list
	err = models.SaveNFTSListByEvmAssetsDebankV1ID(tx, evmAssetsDebankV1.EvmAssetID, apiResponse.NFTList)
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return &models.GlobalWallet{}, err
	}

	// get the data
	walletResponse, err := models.GetGlobalWalletWithEvmDebankInfo(db, address)
	if err != nil {
		return &models.GlobalWallet{}, err
	}

	return walletResponse, nil
}

// processDebankResponses processes wallet responses and returns a slice of debank responses.
func processDebankResponses(ch <-chan *models.GlobalWallet) []*responses.PortfolioResponse {
	debankResponses := make([]*responses.ChainsResponse, 0)
	for walletResponse := range ch {
		if walletResponse == nil || walletResponse.EvmAssetsDebankV1 == nil {
			continue
		}

		for _, token := range *walletResponse.EvmAssetsDebankV1.TokenList {
			debankResponse := &responses.ChainsResponse{}
			debankResponse.DebankTokenResponse(walletResponse, &token)
			debankResponses = append(debankResponses, debankResponse)
		}
	}

	debankResponses = responses.FilterVerifiedResponses(debankResponses)

	debankResponses = responses.AssetTotalsOnChainsAndCalculatePercentages(debankResponses)

	resp := responses.GroupByAssetSymbolToList(debankResponses)

	resp = responses.CalculatePortfolioResponse(resp)

	return resp
}
