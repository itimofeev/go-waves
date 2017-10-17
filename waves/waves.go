package waves

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

var client = http.DefaultClient

func NewClient(id ChainID, nodeAddr string) Client {
	return &wavesClient{
		ChainID: id,
		NodeAddr:nodeAddr,
	}
}

type Client interface {
	GenerateAccount(seedString string) *Account
}

type wavesClient struct {
	ChainID ChainID
	NodeAddr string
	debug bool
}

type Tx struct {
	Type            int    `json:"-"`
	AssetID         string `json:"assetId"`
	SenderPublicKey string `json:"senderPublicKey"`
	Recipient       string `json:"recipient"`
	Fee             uint64 `json:"fee"`
	Amount          uint64 `json:"amount"`
	Attachment      string `json:"attachment"`
	Timestamp       uint64 `json:"timestamp"`
	Signature       string `json:"signature"`
	AmountAssetID   string `json:"amountAssetId"`
	FeeAssetID      string `json:"feeAssetId"`
}

func (t *Tx) TxData() string {
	txBuf := &bytes.Buffer{}

	txBuf.WriteByte(byte(t.Type))

	txBuf.Write(decodeBase58(t.SenderPublicKey))

	//Amount's asset flag (0-Waves, 1-Asset)
	if len(t.AmountAssetID) > 0 {
		txBuf.WriteByte(1)
		txBuf.Write(decodeBase58(t.AmountAssetID))
	} else {
		txBuf.WriteByte(0)
	}

	// Fee's asset flag (0-Waves, 1-Asset)
	if len(t.FeeAssetID) > 0 {
		txBuf.WriteByte(1)
		txBuf.Write(decodeBase58(t.FeeAssetID))
	} else {
		txBuf.WriteByte(0)
	}

	txBuf.Write(uint64ToBytes(t.Timestamp)) //Fee's asset flag (0-Waves, 1-Asset)

	txBuf.Write(uint64ToBytes(t.Amount)) //Amount

	txBuf.Write(uint64ToBytes(t.Fee)) //Fee

	txBuf.Write(decodeBase58(t.Recipient)) //Recipient's address

	txBuf.Write(uint16ToBytes(uint16(len(decodeBase58(t.Attachment)))))

	txBuf.Write(decodeBase58(t.Attachment))

	return encodeBase58(txBuf.Bytes()) //Timestamp
}

type SignResponse struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func signMessage(message, privateKey string) string {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://52.30.47.67:6869/utils/sign/%s", privateKey), strings.NewReader(message))
	must(err)

	resp, err := client.Do(req)
	must(err)

	defer resp.Body.Close()

	httputil.DumpResponse(resp, true)

	body, err := ioutil.ReadAll(resp.Body)
	must(err)

	signResponse := &SignResponse{}

	err = json.Unmarshal(body, signResponse)
	must(err)

	return signResponse.Signature
}

const wavesCoef = 100000000

func sendTx(acc *Account) {
	tx := Tx{
		Type: 4,
		//AssetID:         "E9yZC4cVhCDfbjFJCc9CqkAtkoFy5KaCe64iaxHM2adG",
		//FeeAssetID:      "E9yZC4cVhCDfbjFJCc9CqkAtkoFy5KaCe64iaxHM2adG",
		SenderPublicKey: acc.Public,
		Recipient:       "3Mx2afTZ2KbRrLNbytyzTtXukZvqEB8SkW7",
		Fee:             1 * wavesCoef,
		Amount:          1 * wavesCoef,
		Attachment:      encodeBase58([]byte{1, 2, 3, 4}),
		Timestamp:       uint64(time.Now().UnixNano() / int64(time.Millisecond)),
	}

	tx.Signature = signMessage(tx.TxData(), acc.Private)

	txData, err := json.Marshal(tx)
	must(err)
	req, err := http.NewRequest(http.MethodPost, "http://52.30.47.67:6869/assets/broadcast/transfer", bytes.NewReader(txData))
	req.Header.Set("api-key-hash", "BASE58APIKEYHASH")
	req.Header.Set("Content-Type", "application/json")
	must(err)

	resp, err := client.Do(req)
	must(err)

	respDump, err := httputil.DumpResponse(resp, true)
	must(err)
	fmt.Println("!!!", string(respDump))
}
