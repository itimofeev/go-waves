package waves

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type ChainID byte

//noinspection GoUnusedConst
const (
	ChainIDMain = 'W'
	ChainIDTest = 'T'
)

type Account struct {
	Private string
	Public  string
	Address string
	Seed    string
}

func (c *wavesClient) GenerateAccount(seedString string) *Account {
	seedBytes := []byte(seedString)

	if c.debug {
		fmt.Println("seed string in bas58", EncodeBase58(seedBytes))
		fmt.Println("seed string with nonce", EncodeBase58(prependNonce(0, seedBytes)))
	}

	accountSeed := makeSeedHash(0, seedBytes)
	if c.debug {
		fmt.Println("account seed with nonce 0", EncodeBase58(accountSeed))
	}

	hashedAccountSeed := sha256.Sum256(accountSeed)
	if c.debug {
		fmt.Println("hashed account seed", EncodeBase58(hashedAccountSeed[:]))
	}

	priv, pub := generateCurve25519Keys(hashedAccountSeed[:])
	if c.debug {
		fmt.Println("Private key", EncodeBase58(priv))
		fmt.Println("Public key", EncodeBase58(pub))
	}

	addr := generateAddr(pub, c.ChainID)
	if c.debug {
		fmt.Println("Address", EncodeBase58(addr))
	}

	return &Account{
		Private: EncodeBase58(priv),
		Public:  EncodeBase58(pub),
		Address: EncodeBase58(addr),
		Seed:    seedString,
	}
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

func prependNonce(nonce uint32, seedBytes []byte) []byte {
	return append(uint32ToBytes(nonce), seedBytes...)
}
