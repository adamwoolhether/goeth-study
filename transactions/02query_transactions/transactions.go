package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatalln(err)
	}

	// Using the block number we retrieved in the previous example,
	// we can also read all the transactions in a block by
	// calling the Transactions method, which returns a list
	// of Transaction tye.
	blockNumber := big.NewInt(12527961)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range block.Transactions() {
		fmt.Printf("block Hash(hex): %v\n", tx.Hash().Hex())
		fmt.Printf("tx value: %v\n", tx.Value().String)
		fmt.Printf("gas: %v\n", tx.Gas())
		fmt.Printf("gas price: %v\n", tx.GasPrice().Uint64())
		fmt.Printf("nonce: %v\n", tx.Nonce())
		fmt.Printf("data: %v\n", tx.Data())
		fmt.Printf("to hex: %v\n", tx.To().Hex())


		// To read the sender's address, you need to call AsMessage, which returns a Message
		// type and contains a function to return the from address. It requires an EIP55 signer,
		//	derived from the chain ID from client.
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err != nil {
			fmt.Printf("sender: %v", msg.From().Hex())
		}

		// Each transaction contains the result of the transaction execution.
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("recipt status: %v", receipt.Status)
		fmt.Printf("recipt logs: %v", receipt.Logs)
	}

	// To iterate over a transaction without fetching the block, call the
	// client's TransactionBlock method, which only accepts the block hash
	// and index of the transaction within the block. Here, TransactionCount
	// tells us how many transactions are in the block
	blockHash := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < count; idx ++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("transaction in block: %v\n", tx.Hash().Hex())
	}

	// Single transactions can be queries directly given the transaction hash
	// with TransactionByHash
	txHash := common.HexToHash("0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("transaction hex hash: %v\n", tx.Hash().Hex())
	fmt.Println(isPending)
}
