package models

type (
	SolanaPortfolio struct {
		TotalNativeBalance TotalNativeBalance `json:"nativeBalance"`
		Tokens             []TokensAndNfts    `json:"tokens"`
		Nfts               []TokensAndNfts    `json:"nfts"`
	}

	TotalNativeBalance struct {
		Lamports string `json:"lamports"`
		Solana   string `json:"solana"`
	}

	TokensAndNfts struct {
		AssociatedTokenAddress string `json:"associatedTokenAddress"`
		Mint                   string `json:"mint"`
		AmountRaw              string `json:"amountRaw"`
		Amount                 string `json:"amount"`
		Decimals               string `json:"decimals"`
	}
)
