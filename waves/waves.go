package waves

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

var client = http.DefaultClient

func NewClient(id ChainID, nodeAddr, apiKeyHash string) Client {
	return &wavesClient{
		ChainID:    id,
		NodeAddr:   nodeAddr,
		APIKeyHash: apiKeyHash,
	}
}

type Client interface {
	GenerateAccount(seedString string) *Account
	SendTransferTx(from *Account, toAddress string, amount, fee int) (*TxResponse, error)
}

type wavesClient struct {
	ChainID    ChainID
	NodeAddr   string
	APIKeyHash string
	debug      bool
}

// Golang doesn't have method to sign message by Curve25519 algorithm
// Use node api to sign for now as a workaround
func (c *wavesClient) signMessage(message, privateKey string) string {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/utils/sign/%s", c.NodeAddr, privateKey), strings.NewReader(message))
	must(err)

	resp, err := client.Do(req)
	must(err)

	defer resp.Body.Close()

	httputil.DumpResponse(resp, true)

	body, err := ioutil.ReadAll(resp.Body)
	must(err)

	type SignResponse struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}
	signResponse := &SignResponse{}

	err = json.Unmarshal(body, signResponse)
	must(err)

	return signResponse.Signature
}

const wavesCoef = 100000000

type TxResponse struct {
	Type int `json:"type"`
	ID string `json:"id"`
	Sender string `json:"sender"`
	SenderPublicKey string `json:"senderPublicKey"`
	Fee uint64 `json:"fee"`
	Timestamp int64 `json:"timestamp"`
	Signature string `json:"signature"`
	Recipient string `json:"recipient"`
	AssetID string `json:"assetId"`
	Amount uint64 `json:"amount"`
	FeeAssetID string `json:"feeAsset"`
	Attachment string `json:"attachment"`
}

func (c *wavesClient) SendTransferTx(from *Account, toAddress string, amount, fee int) (*TxResponse, error) {
	if amount <= 0 || fee <= 0 {
		return nil, errors.New("amount and fee must be positive")
	}
	tx := Tx{
		Type:            4,
		SenderPublicKey: from.Public,
		Recipient:       toAddress,
		Fee:             uint64(amount) * wavesCoef,
		Amount:          uint64(fee) * wavesCoef,
		Attachment:      EncodeBase58([]byte{1, 2, 3, 4}),
		Timestamp:       uint64(time.Now().UnixNano() / int64(time.Millisecond)),
	}

	tx.Signature = c.signMessage(tx.TxData(), from.Private)

	txData, err := json.Marshal(tx)
	must(err)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/assets/broadcast/transfer", c.NodeAddr), bytes.NewReader(txData))
	req.Header.Set("api-key-hash", c.APIKeyHash)
	req.Header.Set("Content-Type", "application/json")
	must(err)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	must(err)

	if c.debug {
		respDump, err := httputil.DumpResponse(resp, true)
		must(err)
		fmt.Println("!!!", string(respDump))
	}

	body, err := ioutil.ReadAll(resp.Body)
	must(err)

	txResponse := &TxResponse{}
	err = json.Unmarshal(body, txResponse)
	must(err)

	return txResponse, nil
}
