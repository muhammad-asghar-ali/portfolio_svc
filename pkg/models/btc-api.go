package models

import (
	"time"
)

type (
	BitcoinAddressInfo struct {
		Address             string    `json:"address" gorm:"type:varchar(255);uniqueIndex"`
		Received            int64     `json:"received" gorm:"type:bigint"`
		Sent                int64     `json:"sent" gorm:"type:bigint"`
		Balance             int32     `json:"balance" gorm:"type:int"`
		TxCount             int16     `json:"tx_count" gorm:"type:smallint"`
		UnconfirmedTxCount  int16     `json:"unconformed_tx_count" gorm:"type:smallint"`
		UnconfirmedReceived int32     `json:"unconfirmed_received" gorm:"type:int"`
		UnconfirmedSent     int32     `json:"unconfirmed_sent" gorm:"type:int"`
		UnspentTxCount      int32     `json:"unspend_tx_count" gorm:"type:int"`
		FirstTx             string    `json:"first_tx" gorm:"type:text"`
		LastTx              string    `json:"last_tx" gorm:"type:text"`
		CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	}
)

func (BitcoinAddressInfo) TableName() string {
	return "bitcoin_address_info"
}
