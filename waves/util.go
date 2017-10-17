package waves

import (
	"bytes"
	"encoding/binary"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wsddn/go-ecdh"
	"golang.org/x/crypto/blake2b"
)

func uint32ToBytes(n uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, n)

	return buf
}

func uint16ToBytes(n uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, n)

	return buf
}

func uint64ToBytes(n uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, n)

	return buf
}

func makeSeedHash(nonce uint32, seedBytes []byte) []byte {
	return secureHash(prependNonce(nonce, seedBytes))
}

func secureHash(data []byte) []byte {
	sum := blake2b.Sum256(data)
	return crypto.Keccak256(sum[:])
}

func generateCurve25519Keys(accountSeed []byte) (priv []byte, pub []byte) {
	g := ecdh.NewCurve25519ECDH()
	privI, pubI, err := g.GenerateKey(bytes.NewReader(accountSeed))
	must(err)

	return g.Marshal(privI), g.Marshal(pubI)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func EncodeBase58(src []byte) string {
	return base58.Encode(src)
}

func DecodeBase58(src string) []byte {
	return base58.Decode(src)
}
