package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"regexp"
)

func main() {

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatal(err)
	}

	// Use regex to check if the eth address is valid
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	fmt.Printf("is valid: %v\n", re.MatchString("0x323b5d4c32345ced77393b3530b1eed0f346429d"))
	fmt.Printf("is valid: %v\n", re.MatchString("0xZYXb5d4c32345ced77393b3530b1eed0f346429d"))

	// Check if the address is an account or a smart contract.
	// Smart contracts have bytecode stored at the address.
	// Example of an 0x Protol Token (ZRX) smart contract address:
	address := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	bytecode, err := client.CodeAt(context.Background(), address, nil) // nil is the latest block
	if err != nil {
		log.Fatal(err)
	}
	isContract := len(bytecode) > 0

	fmt.Printf("is contract: %v\n", isContract) // is contract: true


	// Example of a random account address(non contract)
	address = common.HexToAddress("0x8e215d06ea7ec1fdb4fc5fd21768f4b34ee92ef4")
	bytecode, err = client.CodeAt(context.Background(), address, nil) // nil is the latest block
	if err != nil {
		log.Fatal(err)
	}
	isContract = len(bytecode) > 0

	fmt.Printf("is contract: %v\n", isContract) // is contract: true

}