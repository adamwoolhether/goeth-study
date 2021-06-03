package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatalln(err)
	}

	// You must load your private key
	privateKey, err := crypto.HexToECDSA("PRIVATE KEY HERE")
	if err != nil {
		log.Fatal(err)
	}

	// Then you need the nonce, a number used only once. New accounts
	// have a none of 0. Every transaction thereafter have a nonce
	// incremented by 1. The ethereum client proves a helper method
	// PendingNonceAt to return the next nonce you should use.
	// It require the public address of the sending account, which can
	// be derived from the private key.
	publicKey := privateKey.Public() // returns an interface containing public key.
	// perform a tye assertion to set type of our publicKey var, assigning it to pubclicKeyECDSA
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("can't asset type: publicKey isn't of type * ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress((*publicKeyECDSA))

	// Now read the nonce to be used for the transaction
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Set the amount of ETH to be transferred. ETH must be
	// converted to wei, in accordance with blockchain. 1 ETH = 1 + 18 zeroes
	// https://etherconverter.netlify.app/
	value := big.NewInt(250000000000000000)

	// Set gas limit. Standard is 21000 units
	gasLimit := uint64(21000)

	// Gas price is also set in wei. 30 gwei is typically fast
	gasPriceHard := big.NewInt(30000000000)
	_ = gasPriceHard

	// Hard-coding as gas price may not be ideal, you can use the
	// SuggestGasPrice function to get average gas based on the
	// x number of previous blocks
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// To determine who eth is being sent to:
	toAddress := common.HexToAddress("0x9DeeD15f4164a2d0ef37f6c204eb3965D7e28D46")

	// Generate an unsigned ethereum transaction with NewTransaction
	// WARNING: NEWTRANSACTION IS DEPRECATED
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction with the private key of the sender. This
	// is done with the SignTx method, which takes in the unsigned tx
	// and private key contstructed earlier. SignTx require the EIP155 signer,
	// derived from the chain ID from the client.
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Finally we broadcast the transactin to the entire network with
	// SendTransaction on the client, which takes the signed transaction.
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
}