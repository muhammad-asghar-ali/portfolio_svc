package controllers

import (
	"net/http"
	"strings"
	"sync"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/internal/responses"
	"github.com/0xbase-Corp/portfolio_svc/providers/bitcoin"
	"github.com/0xbase-Corp/portfolio_svc/providers/debank"
	"github.com/0xbase-Corp/portfolio_svc/providers/solana"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/0xbase-Corp/portfolio_svc/shared/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	PortfolioAddresses struct {
		BTC []string `json:"btc"`
		Sol []string `json:"sol"`
		EVM []string `json:"evm"`
	}

	ChannelMap struct {
		btcChs    map[string]chan *models.GlobalWallet
		solanaChs map[string]chan *models.GlobalWallet
		debankChs map[string]chan *models.GlobalWallet
	}
)

//	@BasePath	/api/v1

// AllPortfolioController godoc
//
// AllPortfolioController defines the route and Swagger annotations for fetching all portfolio information.
// @Summary      Fetch all portfolio information
// @Description  Retrieves information for all portfolios including Bitcoin, Solana, and EVM addresses.
// @Tags         portfolio
// @Accept       json
// @Produce      json
// @Param        addresses body PortfolioAddresses true "Portfolio Addresses"
// @Success      200 {object} []responses.PortfolioResponse
// @Failure      400 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /api/v1/all-portfolio [post]
func AllPortfolioController(c *gin.Context, db *gorm.DB) {
	requestBody := PortfolioAddresses{}

	if err := c.BindJSON(&requestBody); err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError("error"+err.Error()))
		return
	}

	btcAddresses := utils.UniqueAddress(requestBody.BTC)
	solanaAddresses := utils.UniqueAddress(requestBody.Sol)
	debankAddresses := utils.UniqueAddress(requestBody.EVM)

	if len(btcAddresses) == 0 && len(solanaAddresses) == 0 && len(debankAddresses) == 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError("empty addresses"))
		return
	}

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	errorCh := make(chan error)
	bitcoinAPIClient := &bitcoin.BitcoinAPI{}
	solanaAPIClient := &solana.SolanaAPI{}
	debankAPIClient := &debank.DebankAPI{}

	channelMap := createChannelMap(btcAddresses, solanaAddresses, debankAddresses, wg, mutex, errorCh, db, bitcoinAPIClient, solanaAPIClient, debankAPIClient)

	go func() {
		wg.Wait()
		closeChannels(channelMap)
		close(errorCh)
	}()

	errs := errors.HandleChannelErrors(errorCh, c)
	if len(errs) > 0 {
		errors.HandleHttpError(c, errors.NewBadRequestError(strings.Join(errs, "; ")))
		return
	}

	allResponses := processResponses(channelMap)

	c.JSON(http.StatusOK, allResponses)
}

func createChannelMap(btcAddresses, solanaAddresses, debankAddresses []string, wg *sync.WaitGroup, mutex *sync.Mutex, errorCh chan error, db *gorm.DB, bitcoinAPIClient *bitcoin.BitcoinAPI, solanaAPIClient *solana.SolanaAPI, debankAPIClient *debank.DebankAPI) ChannelMap {
	channelMap := ChannelMap{
		btcChs:    make(map[string]chan *models.GlobalWallet),
		solanaChs: make(map[string]chan *models.GlobalWallet),
		debankChs: make(map[string]chan *models.GlobalWallet),
	}

	for _, btcAddress := range btcAddresses {
		channelMap.btcChs[btcAddress] = make(chan *models.GlobalWallet, 1)
		wg.Add(1)
		go fetchAndSaveBtc(db, bitcoinAPIClient, btcAddress, channelMap.btcChs[btcAddress], wg, mutex, errorCh)
	}

	for _, solAddress := range solanaAddresses {
		channelMap.solanaChs[solAddress] = make(chan *models.GlobalWallet, 1)
		wg.Add(1)
		go fetchAndSaveSolana(db, solanaAPIClient, solAddress, channelMap.solanaChs[solAddress], wg, mutex, errorCh)
	}

	for _, evmAddress := range debankAddresses {
		channelMap.debankChs[evmAddress] = make(chan *models.GlobalWallet, 1)
		wg.Add(1)
		go fetchAndSaveDebank(db, debankAPIClient, evmAddress, channelMap.debankChs[evmAddress], wg, mutex, errorCh)
	}

	return channelMap
}

func closeChannels(channelMap ChannelMap) {
	for _, ch := range channelMap.btcChs {
		close(ch)
	}
	for _, ch := range channelMap.solanaChs {
		close(ch)
	}
	for _, ch := range channelMap.debankChs {
		close(ch)
	}
}

func processResponses(channelMap ChannelMap) []*responses.PortfolioResponse {
	var allResponses []*responses.PortfolioResponse

	for _, ch := range channelMap.btcChs {
		responses := processBtcResponses(ch)
		allResponses = append(allResponses, responses...)
	}

	for _, ch := range channelMap.solanaChs {
		responses := processSolanaResponses(ch)
		allResponses = append(allResponses, responses...)
	}

	for _, ch := range channelMap.debankChs {
		responses := processDebankResponses(ch)
		allResponses = append(allResponses, responses...)
	}

	return allResponses
}
