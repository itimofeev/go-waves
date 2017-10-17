package waves

import "bytes"

// Tx is an transaction in waves blockchain
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

// TxData serializes data from transaction into binary message for sign
func (t *Tx) TxData() string {
	txBuf := &bytes.Buffer{}

	txBuf.WriteByte(byte(t.Type))

	txBuf.Write(DecodeBase58(t.SenderPublicKey))

	//Amount's asset flag (0-Waves, 1-Asset)
	if len(t.AmountAssetID) > 0 {
		txBuf.WriteByte(1)
		txBuf.Write(DecodeBase58(t.AmountAssetID))
	} else {
		txBuf.WriteByte(0)
	}

	// Fee's asset flag (0-Waves, 1-Asset)
	if len(t.FeeAssetID) > 0 {
		txBuf.WriteByte(1)
		txBuf.Write(DecodeBase58(t.FeeAssetID))
	} else {
		txBuf.WriteByte(0)
	}

	txBuf.Write(uint64ToBytes(t.Timestamp)) //Fee's asset flag (0-Waves, 1-Asset)

	txBuf.Write(uint64ToBytes(t.Amount)) //Amount

	txBuf.Write(uint64ToBytes(t.Fee)) //Fee

	txBuf.Write(DecodeBase58(t.Recipient)) //Recipient's address

	txBuf.Write(uint16ToBytes(uint16(len(DecodeBase58(t.Attachment)))))

	txBuf.Write(DecodeBase58(t.Attachment))

	return EncodeBase58(txBuf.Bytes()) //Timestamp
}

