package responses

import (
	"github.com/0xbase-Corp/portfolio_svc/shared/utils"
)

// filder only IsVerified = true
func FilterVerifiedResponses(responses []*ChainsResponse) []*ChainsResponse {
	filteredResponses := []*ChainsResponse{}

	for _, response := range responses {
		if response.IsVerified {
			filteredResponses = append(filteredResponses, response)
		}
	}

	return filteredResponses
}

// Add this function to calculate chains totals
func AssetTotalsOnChainsAndCalculatePercentages(responses []*ChainsResponse) []*ChainsResponse {
	totalAssetAmount := make(map[string]float64)

	// Calculate the total value per chain
	for _, response := range responses {
		totalAssetAmount[response.AssetSymbol] += response.Quantity
	}

	// Calculate and update the percentage of each asset within its chain
	for _, response := range responses {
		if total, ok := totalAssetAmount[response.AssetSymbol]; ok && total > 0 {
			response.AssetPercentage = utils.Ternary(total != 0, (response.Quantity/total)*100, 0)
		}
	}

	return responses
}

func GroupByAssetSymbol(responses []*ChainsResponse) map[string][]*ChainsResponse {
	groupedResponses := make(map[string][]*ChainsResponse)

	for _, response := range responses {
		if _, ok := groupedResponses[response.AssetSymbol]; !ok {
			groupedResponses[response.AssetSymbol] = []*ChainsResponse{response}
		} else {
			groupedResponses[response.AssetSymbol] = append(groupedResponses[response.AssetSymbol], response)
		}
	}

	return groupedResponses
}

func GroupByAssetSymbolToList(responses []*ChainsResponse) []*PortfolioResponse {
	groupedResponses := GroupByAssetSymbol(responses)

	result := make([]*PortfolioResponse, 0)

	for assetSymbol, responses := range groupedResponses {
		var assetData PortfolioResponse
		assetData.AssetSymbol = assetSymbol

		if len(responses) > 1 {
			data := make([]*ChainsResponse, 0)
			data = append(data, responses...)
			assetData.ChainsInfo = data
		} else {
			response := responses[0]
			assetData.ChainsInfo = []*ChainsResponse{response}
		}

		result = append(result, &assetData)
	}

	return result
}

func CalculatePortfolioResponse(portfolioResponse []*PortfolioResponse) []*PortfolioResponse {
	var totalPortfolioPrice float64
	for _, portfolio := range portfolioResponse {
		var quantity, total_price float64

		// Iterate over the chain data to accumulate values
		for _, chain := range portfolio.ChainsInfo {
			quantity += chain.Quantity
			total_price += chain.TotalPrice
		}

		// Update the PortfolioResponse struct with the accumulated values
		portfolio.Quantity = quantity
		portfolio.TotalPrice = total_price
		portfolio.UnitPrice = MostCommonUnitPrice(portfolio)

		// Add the TotalPrice of each PortfolioResponse to the totalPortfolioPrice
		totalPortfolioPrice += portfolio.TotalPrice
	}

	// Loop through portfolioResponse and update the PortfolioPercentage
	for _, portfolio := range portfolioResponse {
		portfolio.PortfolioPercentage = utils.Ternary(totalPortfolioPrice != 0, (portfolio.TotalPrice/totalPortfolioPrice)*100, 0)
	}

	return portfolioResponse
}

func MostCommonUnitPrice(portfolio *PortfolioResponse) float64 {
	unitPriceCount := make(map[float64]int)
	var mostCommonUnitPrice float64
	maxCount := 0

	// Count occurrences of unit_price values
	for _, chain := range portfolio.ChainsInfo {
		unitPriceCount[chain.UnitPrice]++
		if unitPriceCount[chain.UnitPrice] > maxCount {
			maxCount = unitPriceCount[chain.UnitPrice]
			mostCommonUnitPrice = chain.UnitPrice
		}
	}

	return mostCommonUnitPrice
}
