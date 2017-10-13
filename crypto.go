package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wsddn/go-ecdh"
	"golang.org/x/crypto/blake2b"
)

const seedString = "manage manual recall harvest series desert melt police rose hollow moral pledge kitten position add"

type ChainID byte

//noinspection GoUnusedConst
const (
	ChainIDMain = 'W'
	ChainIDTest = 'T'
)

func someFunc() {
	seedBytes := []byte(seedString)

	fmt.Println("seed string in bas58", encodeBase58(seedBytes))
	fmt.Println("seed string with nonce", encodeBase58(prependNonce(0, seedBytes)))

	accountSeed := makeSeedHash(0, seedBytes)
	fmt.Println("account seed with nonce 0", encodeBase58(accountSeed))

	hashedAccountSeed := sha256.Sum256(accountSeed)
	fmt.Println("hashed account seed", encodeBase58(hashedAccountSeed[:]))

	priv, pub := generateCurve25519Keys(hashedAccountSeed[:])
	fmt.Println("Private key", encodeBase58(priv))
	fmt.Println("Public key", encodeBase58(pub))

	addr := generateAddr(pub, ChainIDMain)
	fmt.Println("Address", encodeBase58(addr))

	txBuf := bytes.Buffer{}

	txBuf.WriteByte(byte(4)) // TX type = 4 (transfer)

}

func generateAddr(pub []byte, chainID ChainID) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 26))
	version := byte(1)
	buf.WriteByte(version)
	buf.WriteByte(byte(chainID))
	buf.Write(secureHash(pub)[:20])

	if len(buf.Bytes()) != 22 {
		panic(len(buf.Bytes()))
	}

	buf.Write(secureHash(buf.Bytes())[:4])

	return buf.Bytes()
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
