package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	// As before, first establish connect, sign with private key, & define receiver and gas.

	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("PRIVATEKEYHERE")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: pulicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(200000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x9DeeD15f4164a2d0ef37f6c204eb3965D7e28D46")
	var data []byte

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Here we learn how to get raw transaction data to
	// broadcast it at a later time. First, construct the
	// object and sign it:
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a types.Transactions type with the signed
	// transaction as the first value. This is done because
	// the Transactions type provides a GetRlp method for
	// returning the transaction in RLP encoded format, a
	// special encoding method Eth uses for serializing objects.
	//ts := types.Transactions{signedTx}
	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}
	rawTxHex := hex.EncodeToString(rawTxBytes)
	fmt.Printf("raw tx hex: %v", rawTxHex)

}
