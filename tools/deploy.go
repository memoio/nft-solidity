// $go run tools/deploy-contracts/deploy.go --sk=xxx --sk1=xxx --sk2=xxx --sk3=xxx --sk4=xxx --sk5=xxx
// Run this command to deploy accountDid and fileDid contract.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/memoio/nft-solidity/go-contracts/token"

	"github.com/ethereum/go-ethereum/ethclient"

	com "github.com/memoio/contractsv2/common"
)

var (
	endPoint string
	hexSk    string

	allGas = make([]uint64, 0)
)

// deploy all contracts by admin with specified chain
func main() {
	inputeth := flag.String("eth", "dev", "eth api Address;") //dev test or product
	sk := flag.String("sk", "", "signature for sending transaction")

	flag.Parse()

	chain := *inputeth
	_, endPoint = com.GetInsEndPointByChain(chain)
	fmt.Println()
	fmt.Println("endPoint:", endPoint)
	hexSk = *sk

	// get client
	client, err := ethclient.DialContext(context.TODO(), endPoint)
	if err != nil {
		log.Fatal(err)
	}

	// make auth to send transaction
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		chainID = big.NewInt(666)
	}
	txAuth, err := com.MakeAuth(chainID, hexSk)
	if err != nil {
		log.Fatal(err)
	}

	// deploy Proxy
	erc721Addr, tx, erc721Ins, err := token.DeployERC721(txAuth, client, "xspace", "uni")
	if err != nil {
		log.Fatal("deploy erc721 err:", err)
	}
	fmt.Println("erc721Addr: ", erc721Addr.Hex())
	go com.PrintGasUsed(endPoint, tx.Hash(), "deploy erc721 gas:", &allGas)
	_ = erc721Ins

	time.Sleep(time.Second)

	// print gas used information
	fmt.Println("tx num:", len(allGas), "\t", "all gasUsed:", allGas)
	var totalGas uint64
	for _, gas := range allGas {
		totalGas += gas
	}
	fmt.Println("totalGas:", totalGas)

	fmt.Println()
}
