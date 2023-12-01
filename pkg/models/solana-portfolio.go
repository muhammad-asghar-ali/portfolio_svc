package models

type SolanaPortfolio struct {
	TotalNativeBalance totalNativeBalance `json:"nativeBalance"`
	Tokens             []tokensAndNfts    `json:"tokens"`
	Nfts               []tokensAndNfts    `json:"nfts"`
}

type totalNativeBalance struct {
	Lamports string `json:"lamports"`
	Solana   string `json:"solana"`
}

type tokensAndNfts struct {
	AssociatedTokenAddress string `json:"associatedTokenAddress"`
	Mint                   string `json:"mint"`
	AmountRaw              string `json:"amount raw"`
	Amount                 string `json:"amount"`
	Decimals               string `json:"decimals"`
}
