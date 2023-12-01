package controllers

import (
	// "encoding/json"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SolanaTotalBalance struct {
	Lamports string `json:"lamports"`
	Solana   string `json:"solana"`
}

func SolanaController(c *gin.Context, db *gorm.DB) {
	solAddress := c.Param("sol-address")
	moralisAccessKey := c.GetHeader("x-api-key")

	url := "https://solana-gateway.moralis.io/account/mainnet/" + solAddress + "/balance"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("x-api-key", moralisAccessKey)

	//send the request
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var solResponse SolanaTotalBalance
	if err := json.Unmarshal(body, &solResponse); err != nil {
		log.Fatal(err)
	}

	//TODO: add the data into database

	c.JSON(http.StatusOK, gin.H{
		"message": solResponse,
	})
}
