package main

import (
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/blake2b"
)

const seedString = "manage manual recall harvest series desert melt police rose hollow moral pledge kitten position add"

func main() {
	seedBytes := []byte(seedString)

	fmt.Println("seed string in bas58", encodeBase58(seedBytes))
	fmt.Println("seed string with nonce", encodeBase58(prependNonce(0, seedBytes)))

	accountSeed := makeSeedHash(0, seedBytes)
	fmt.Println("account seed with nonce 0", encodeBase58(accountSeed))

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
	blake2bSum := make([]byte, 32)
	sum := blake2b.Sum256(prependNonce(nonce, seedBytes))
	copy(blake2bSum, sum[:])

	return crypto.Keccak256(blake2bSum)
}
