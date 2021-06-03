package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	// An Eth provider that supports RPC over websockets is needed.
	client, err := ethclient.Dial("wss://ropsten.infura.io/ws/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new channel to receive the latest block headers
	headers := make(chan *types.Header)

	// Call the client's SubscribeNewHead method. It takes in the
	// headers channel and returns a subscription object.
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	// The subscription pushes new block headers to the channel, a
	// select statement listens for new messages. An error channel
	// sends new messages in case of a failure with subscription.
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex())

			// To get the block's full contents, we pass the block's
			// header hash to the client's BlockByHash func.
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			// We can read the block's entire metadata fields, some examples:
			fmt.Printf("block hashhex: %v\n", block.Hash().Hex())
			fmt.Printf("block #(uint64): %v\n", block.Number().Uint64())
			fmt.Printf("block time: %v\n", block.Time())
			fmt.Printf("block nonce: %v\n", block.Nonce())
			fmt.Printf("block transactions: %v\n", len(block.Transactions()))
		}
	}

}
