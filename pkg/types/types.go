// these structs are used for 3rd party APIS or for local purposes
// not related with the database models
package types

type (
	BtcChainAPI struct {
		Data    ChainData `json:"data"`
		ErrCode int16     `json:"error_code"`
		ErrNo   int16     `json:"err_no"`
		Message string    `json:"message"`
		Status  string    `json:"status"`
	}

	ChainData struct {
		Address             string `json:"address"`
		Received            int64  `json:"received"`
		Sent                int64  `json:"sent"`
		Balance             int32  `json:"balance"`
		TxCount             int16  `json:"tx_count"`             //max value 32768
		UnconfirmedTxCount  int16  `json:"unconformed_tx_count"` //max value 32768
		UnconfirmedReceived int32  `json:"unconfirmed_received"`
		UnconfirmedSent     int32  `json:"unconfirmed_sent"`
		UnspentTxCount      int32  `json:"unspend_tx_count"`
		FirstTx             string `json:"first_tx"`
		LastTx              string `json:"last_tx"`
	}
)
