package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/responses"
	"github.com/0xbase-Corp/portfolio_svc/providers"
	"github.com/0xbase-Corp/portfolio_svc/providers/solana"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/0xbase-Corp/portfolio_svc/shared/utils"
)

//	@BasePath	/api/v1
//
// SolanaController godoc
// @Summary      Fetch Solana portfolio details for a given Solana address
// @Description  Fetch Solana portfolio details, including tokens and NFTs, for a specific Solana address.
// @Tags         solana
// @Accept       json
// @Produce      json
// @Param        addresses  query      array  true  "Solana Addresses" Format(string)
// @Success      200 {object} []responses.PortfolioResponse
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/solana [get]
func SolanaController(c *gin.Context, db *gorm.DB, apiClient providers.APIClient) {
	addresses := c.Query("addresses")
	solanaAddresses := strings.Split(addresses, ",")

	solanaAddresses = utils.UniqueAddress(solanaAddresses)

	if len(solanaAddresses) == 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty btc addresses"))
		return
	}

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	ch := make(chan *models.GlobalWallet, 1) // 1 specifies the buffer size of the channel
	errorCh := make(chan error, len(solanaAddresses))

	for _, solAddress := range solanaAddresses {
		wg.Add(1)

		go fetchAndSaveSolana(db, apiClient, solAddress, ch, wg, mutex, errorCh)
	}

	// Use a goroutine to close the channel after all goroutines have finished
	go func() {
		wg.Wait()
		close(errorCh)
		close(ch)
	}()

	errs := errors.HandleChannelErrors(errorCh, c)
	if len(errs) > 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError(strings.Join(errs, "; ")))
		return
	}

	// Collect all results from the channel and process to genernic response for solana
	responses := processSolanaResponses(ch)

	c.JSON(http.StatusOK, responses)
}

//	@BasePath	/api/v1
//
// GetSolanaController godoc
// @Summary      Get Solana portfolio for a wallet
// @Description  Retrieve Solana portfolio details, including tokens and NFTs, for a specific wallet.
// @Tags         solana
// @Accept       json
// @Produce      json
// @Param        wallet_id path int true "Wallet ID" Format(int)
// @Param        offset query int false "Pagination offset" Format(int)
// @Param        limit query int false "Pagination limit" Format(int)
// @Success      200 {object} models.GlobalWallet
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /portfolio/solana-wallet/{wallet_id} [get]
func GetSolanaController(c *gin.Context, db *gorm.DB) {
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

	/**
	explain preload query

	SELECT *
	FROM global_wallets
	LEFT JOIN solana_assets_moralis_v1 ON global_wallets.wallet_id = solana_assets_moralis_v1.wallet_id
	LEFT JOIN tokens ON solana_assets_moralis_v1.solana_asset_id = tokens.solana_asset_id
	LEFT JOIN nfts ON solana_assets_moralis_v1.solana_asset_id = nfts.solana_asset_id
	WHERE global_wallets.wallet_id = ?
	LIMIT limit;
	*/

	// Define a reusable function to apply offset and limit for preload
	prlimit := func(query *gorm.DB) *gorm.DB {
		return query.Offset((page - 1) * limit).Limit(limit)
	}

	err = db.
		Preload("SolanaAssetsMoralisV1.Tokens", prlimit).
		Preload("SolanaAssetsMoralisV1.NFTS", prlimit).
		Preload("SolanaAssetsMoralisV1").
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
func saveSolana(db *gorm.DB, solanaAddress string, apiResponse solana.SolanaApiResponse) (*models.GlobalWallet, error) {
	// Begin a new transaction
	tx := db.Begin()

	wallet, err := models.GetOrCreateWallet(tx, solanaAddress, utils.Solana)
	if err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Assuming wallet is the GlobalWallet record found or created
	walletID := wallet.WalletID

	// Initialize solanaAsset and set the WalletID
	solanaAsset := models.SolanaAssetsMoralisV1{}
	solanaAsset.WalletID = walletID

	// Attempt to save the Solana asset data along with the associated tokens and NFTs.
	if err := models.SaveSolanaData(tx, &solanaAsset, apiResponse.Tokens, apiResponse.NFTs); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// save the solana price feed
	// for now hard code the USD -> TODO: change
	if err := HandleCoingeckoPrice(tx, utils.Solana, "usd"); err != nil {
		tx.Rollback()
		return &models.GlobalWallet{}, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return &models.GlobalWallet{}, err
	}

	// get the data
	walletResponse, err := models.GetGlobalWalletWithSolanaInfo(db, solanaAddress)
	if err != nil {
		return &models.GlobalWallet{}, err
	}

	return walletResponse, nil
}

// Fetch and save the data for one address
// TODO: write a common interface in provider which saves the data in database
func fetchAndSaveSolana(db *gorm.DB, apiClient providers.APIClient, address string, ch chan<- *models.GlobalWallet, wg *sync.WaitGroup, mutex *sync.Mutex, errorCh chan<- error) {
	defer wg.Done()

	body, err := apiClient.FetchData(address)
	if err != nil {
		errorCh <- err
		return
	}

	// ignore the Address which returns error or empty response
	resp := solana.SolanaApiResponse{}
	if err := utils.DecodeJSONResponse(body, &resp); err != nil {
		errorCh <- err
		return
	}

	// Save data to the database
	walletResponse, err := saveSolana(db, address, resp)
	if err != nil {
		errorCh <- err
		return
	}

	// Use a mutex to safely append to the channel.
	mutex.Lock()
	ch <- walletResponse
	mutex.Unlock()
}

// processSolanaResponses processes wallet responses and returns a slice of solana responses.
func processSolanaResponses(ch <-chan *models.GlobalWallet) []*responses.PortfolioResponse {
	solanaResponses := make([]*responses.ChainsResponse, 0)
	for walletResponse := range ch {
		if walletResponse == nil || walletResponse.SolanaAssetsMoralisV1 == nil {
			continue
		}

		for _, token := range *walletResponse.SolanaAssetsMoralisV1.Tokens {
			solResponse := &responses.ChainsResponse{}
			solResponse.SolanaTokenResponse(walletResponse, &token)
			solanaResponses = append(solanaResponses, solResponse)
		}
	}

	solanaResponses = responses.AssetTotalsOnChainsAndCalculatePercentages(solanaResponses)

	resp := responses.GroupByAssetSymbolToList(solanaResponses)

	resp = responses.CalculatePortfolioResponse(resp)

	return resp
}
