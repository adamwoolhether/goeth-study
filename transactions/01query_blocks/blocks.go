package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strconv"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatalln(err)
	}

	// There are two ways to query block information
	// 1. Call the client's HeaderByNumber to return info about a block,
	// passing nil will return the latest block header.
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("latest block header: %v\n", header.Number.String())
	headerInt, _ := strconv.Atoi(header.Number.String())

	// 2. You can also call the client's BlockByNumber method
	// to get the full block. This allows you to read all metadata
	// content such as block #, timestamp, hash, difficulty,
	// as well as list of transactions and much more
	blockNumber := big.NewInt(int64(headerInt))
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("block number: %v\n", block.Number().Uint64())
	fmt.Printf("block time: %v\n", block.Time())
	fmt.Printf("block difficulty: %v\n", block.Difficulty().Uint64())
	fmt.Printf("block Hash(hex): %v\n", block.Hash().Hex())
	fmt.Printf("block transactions: %v\n", len(block.Transactions()))
	//fmt.Println(block.Number().Uint64())
	//fmt.Println(block.Time().Uint64())
	//fmt.Println(block.Difficulty().Uint64())
	//fmt.Println(block.Hash().Hex())
	//fmt.Println(len(block.Transactions()))

	// You can also simply call Transaction count to return
	// just the count of the transactions in a block.
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("block transactions(simple call method): %v\n", count)
}

