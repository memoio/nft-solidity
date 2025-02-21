package test

import (
	"context"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/nft-solidity/go-contracts/token"
	"golang.org/x/xerrors"
)

func TestNFT(t *testing.T) {
	client, err := ethclient.Dial("https://devchain.metamemo.one:8501")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ercIns, err := token.NewERC721(common.HexToAddress("0xa75150D716423c069529A3B2908Eb454e0a00Dfc"), client)
	if err != nil {
		t.Fatal(err)
	}

	sk, err := crypto.HexToECDSA("0a95533a110ee10bdaa902fed92e56f3f7709a532e22b5974c03c0251648a5d4")
	if err != nil {
		t.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		chainID = big.NewInt(985)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(sk, chainID)
	if err != nil {
		t.Fatal(err)
	}
	auth.Value = big.NewInt(0) // in wei

	tokenId, err := ercIns.Id(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokenId)

	tx, err := ercIns.Mint(auth, crypto.PubkeyToAddress(sk.PublicKey), "xspace.com/test")
	if err != nil {
		t.Fatal(err)
	}

	err = CheckTx("https://devchain.metamemo.one:8501", tx.Hash(), "mint")
	if err != nil {
		t.Fatal(err)
	}

	newTokenId, err := ercIns.Id(&bind.CallOpts{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(newTokenId)

	address, err := ercIns.OwnerOf(&bind.CallOpts{}, tokenId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(address)
	t.Log(crypto.PubkeyToAddress(sk.PublicKey))

	uri, err := ercIns.TokenURI(&bind.CallOpts{}, tokenId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(uri)
}

func CheckTx(endpoint string, txHash common.Hash, name string) error {
	var receipt *types.Receipt

	t := 6
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(t) * time.Second)
		receipt = GetTransactionReceipt(endpoint, txHash)
		if receipt != nil {
			break
		}
		t = 5
	}

	if receipt == nil {
		return xerrors.Errorf("%s: cann't get transaction(%s) receipt, not packaged", name, txHash)
	}

	// 0 means fail
	if receipt.Status == 0 {
		if receipt.GasUsed != receipt.CumulativeGasUsed {
			return xerrors.Errorf("%s: transaction(%s) exceed gas limit", name, txHash)
		}

		return xerrors.Errorf("%s: transaction(%s) mined but execution failed, please check your tx input", name, txHash)
	}
	return nil
}

func GetTransactionReceipt(endPoint string, hash common.Hash) *types.Receipt {
	client, err := ethclient.Dial(endPoint)
	if err != nil {
		log.Fatal("rpc.Dial err", err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	receipt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		log.Println("get transaction receipt: ", err)
	}
	return receipt
}
