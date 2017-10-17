package main

import (
	"github.com/itimofeev/go-waves/waves"
	"github.com/sanity-io/litter"
)

const seedString = "unaware club online glance evil prize piano oil beyond oak sell wreck beauty lonely milk"
const nodeAddr = "52.30.47.67:6869"

func main() {
	client := waves.NewClient(waves.ChainIDTest, nodeAddr, "BASE58APIKEYHASH")
	account := client.GenerateAccount(seedString)
	litter.Dump(account)



	tx, _ := client.SendTransferTx(account, "3Mx2afTZ2KbRrLNbytyzTtXukZvqEB8SkW7", 1, 1)
	litter.Dump(tx)
}
