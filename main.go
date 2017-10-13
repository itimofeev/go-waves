package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

func main() {
	sendTx()
}

type Tx struct {
	AssetID         string `json:"assetId"`
	SenderPublicKey string `json:"senderPublicKey"`
	Recipient       string `json:"recipient"`
	Fee             int    `json:"fee"`
	Amount          int    `json:"amount"`
	Attachment      string `json:"attachment"`
	Timestamp       int64  `json:"timestamp"`
	Signature       string `json:"signature"`
}

func sendTx() {
	tx := Tx{
		AssetID:         "E9yZC4cVhCDfbjFJCc9CqkAtkoFy5KaCe64iaxHM2adG",
		SenderPublicKey: "CRxqEuxhdZBEHX42MU4FfyJxuHmbDBTaHMhM3Uki7pLw",
		Recipient:       "3Mx2afTZ2KbRrLNbytyzTtXukZvqEB8SkW7",
		Fee:             100000,
		Amount:          5500000000,
		Attachment:      "BJa6cfyGUmzBFTj3vvvaew",
		Timestamp:       int64(1479222433704),
		Signature:       "2TyN8pNS7mS9gfCbX2ktpkWVYckoAmRmDZzKH3K35DKs6sUoXHArzukV5hvveK9t79uzT3cA8CYZ9z3Utj6CnCEo",
	}

	txData, err := json.Marshal(tx)
	must(err)
	req, err := http.NewRequest(http.MethodPost, "http://52.30.47.67:6869/assets/broadcast/transfer", bytes.NewReader(txData))
	req.Header.Set("api-key-hash", "BASE58APIKEYHASH")
	req.Header.Set("Content-Type", "application/json")
	must(err)

	client := http.Client{}

	resp, err := client.Do(req)
	must(err)

	respDump, err := httputil.DumpResponse(resp, true)
	must(err)
	fmt.Println("!!!", string(respDump))
}
