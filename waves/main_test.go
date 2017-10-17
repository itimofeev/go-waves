package waves

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	tx := Tx{
		Type:            4,
		SenderPublicKey: "EENPV1mRhUD9gSKbcWt84cqnfSGQP5LkCu5gMBfAanYH",
		AmountAssetID:   "BG39cCNUFWPQYeyLnu7tjKHaiUGRxYwJjvntt9gdDPxG",
		FeeAssetID:      "BG39cCNUFWPQYeyLnu7tjKHaiUGRxYwJjvntt9gdDPxG",
		Timestamp:       1479287120875,
		Amount:          1,
		Fee:             1,
		Recipient:       "3NBVqYXrapgJP9atQccdBPAgJPwHDKkh6A8",
		Attachment:      EncodeBase58([]byte{1, 2, 3, 4}),
	}

	wikiValue := "Ht7FtLJBrnukwWtywum4o1PbQSNyDWMgb4nXR5ZkV78krj9qVt17jz74XYSrKSTQe6wXuPdt3aCvmnF5hfjhnd1gyij36hN1zSDaiDg3TFi7c7RbXTHDDUbRgGajXci8PJB3iJM1tZvh8AL5wD4o4DCo1VJoKk2PUWX3cUydB7brxWGUxC6mPxKMdXefXwHeB4khwugbvcsPgk8F6YB"

	//assert.Equal(t, wikiValue, tx.TxData())

	original := DecodeBase58(wikiValue)
	my := DecodeBase58(tx.TxData())

	assert.Equal(t, original[0:1], my[0:1], "TxType")
	assert.Equal(t, original[1:33], my[1:33], "Sender")
	assert.Equal(t, original[33:34], my[33:34], "AmmountAssetId")
	assert.Equal(t, original[33:34], my[33:34], "AmmountAssetId Flag")
	assert.Equal(t, original[34:66], my[34:66], "AmmountAssetId")
	assert.Equal(t, original[66:67], my[66:67], "AmmountFeeId Flag")
	assert.Equal(t, original[67:99], my[67:99], "AmmountFeeId")

	assert.Equal(t, original[99:107], my[99:107], "Timestamp")
	assert.Equal(t, original[107:115], my[107:115], "Amount")
	assert.Equal(t, original[115:123], my[115:123], "Fee")

	assert.Equal(t, original[123:149], my[123:149], "Recipient")
	assert.Equal(t, original[149:151], my[149:151], "Attachments len")
	assert.Equal(t, original[151:], my[151:], "Attachments")
}
