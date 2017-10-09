package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wsddn/go-ecdh"
	"golang.org/x/crypto/blake2b"
	"crypto/sha256"
)

const seedString = "manage manual recall harvest series desert melt police rose hollow moral pledge kitten position add"

func main() {
	seedBytes := []byte(seedString)

	fmt.Println("seed string in bas58", encodeBase58(seedBytes))
	fmt.Println("seed string with nonce", encodeBase58(prependNonce(0, seedBytes)))

	accountSeed := makeSeedHash(0, seedBytes)
	fmt.Println("account seed with nonce 0", encodeBase58(accountSeed))

	hashedAccountSeed := sha256.Sum256(accountSeed)
	fmt.Println("hashed account seed ", encodeBase58(hashedAccountSeed[:]))

	g := ecdh.NewCurve25519ECDH()
	priv, pub, err := g.GenerateKey(bytes.NewReader(hashedAccountSeed[:]))
	must(err)

	fmt.Println("Private key", encodeBase58(g.Marshal(priv)))
	fmt.Println("Public key", encodeBase58(g.Marshal((pub))))
}

func encodeBase58(src []byte) string {
	return base58.Encode(src)
}

func prependNonce(nonce uint32, seedBytes []byte) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, nonce)

	return append(buf, seedBytes...)
}

func makeSeedHash(nonce uint32, seedBytes []byte) []byte {
	sum := blake2b.Sum256(prependNonce(nonce, seedBytes))
	return crypto.Keccak256(sum[:])
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
