package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatalln(err)
	}

	// To broadcase the raw transaction created in the last example, use
	// the raw output given and decode to bytes format.
	rawTx := "f86b10843b9aca00825208949deed15f4164a2d0ef37f6c204eb3965d7e28d468802c68af0bb140000802ba0a590dbb93b8a3b234a2a1b70cb5271b71b74e7f100c3aaa473676424d3ee74c3a00db23f9a4f45a5f402b45ddd4688a62276e319abcadda29ed3c19bbd007bf4e4"
	rawTxBytes, err := hex.DecodeString(rawTx)

	// Initialize a new types.Transaction pointer and call DecodeBytes.
	// Pass the raw tx bytes and pointer.
	tx := new(types.Transaction)
	rlp.DecodeBytes(rawTxBytes, &tx)

	// The transaction can now be broadcast
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())
}