package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// First establish a client connection.
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to infura mainnet..")
	_ = client

	// When interacting with account addresses, we muct first convert them to the
	// common.Address type
	address := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	fmt.Printf("address.Hex() %v\n", address.Hex())
	fmt.Printf("address.Hash().Hex() %v\n", address.Hash().Hex())
	fmt.Printf("address.Bytes() %v\n\n", address.Bytes())


	// ACCOUNT BALANCES
	// Use the client-balances's BalanceAt method to see balances. Passing
	// 'nil' to the block number will return the most recent balance.
	account := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Most recent balance: %v\n\n", balance)

	//// Pass a block number to read balance at the time of the block. It must be a 'big.Int'
	//blockNumber := big.NewInt(12474947)
	//balance, err = client-balances.BalanceAt(context.Background(), account, blockNumber)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("\n\nValue at block 12474947: %v\n\n", balance)

	// Numbers are given in smallest possible unit. They're fixed point precision.
	// In eth, the unit is 'wei'.  Here's how to convert the numbers: wei / 10^18
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Printf("Current balance in ETH: %v\n\n", ethValue)

	// You ca also see the pending account balance with PendingBalanceAt
	pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	fmt.Printf("Pending balance: %v\n\n", pendingBalance)

}
