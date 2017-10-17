package main

import (
	"github.com/itimofeev/go-waves/waves"
	"github.com/sanity-io/litter"
)

const seedString = "unaware club online glance evil prize piano oil beyond oak sell wreck beauty lonely milk"
const nodeAddr = "52.30.47.67:6869"
func main() {
	client := waves.NewClient(waves.ChainIDTest, "52.30.47.67:6869")
	account := client.GenerateAccount(seedString)
	litter.Dump(account)


}
