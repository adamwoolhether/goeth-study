package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	store "goeth/smart_contracts/01smart_contract_compilation"
	"log"
	"math/big"
)

//import (
//	//store "goeth/smart_contracts/02deploying_a_smart_contract/contracts"
//)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatalln(err)
	}

	privateKey, err := crypto.HexToECDSA("KEYHERE")
	if err != nil {
		log.Fatalln(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("can't assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalln(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// First create the keyed transactor with NewKeyedTransactorWithChainID,
	// passing in the private key. Then set usual properties of nonce,
	// value, gas price, gas limit, etc.
	//auth := bind.NewKeyedTransactor(privateKey) // This is deprecated
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	// The go file generated from the Store contract contains a
	// deploy method: Deploy<contractName>. It takes a keyed transactor,
	// ethclient, and any input needed by the smart contract construtor.
	// Our contract takes in a string for the version. It returns the
	// contract's new Eth address, transaction object, contract
	// instance and an error.
	input := "1.0"
	address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("address.Hex: %s\n", address.Hex())
	fmt.Printf("tx.Hash: %s\n", tx.Hash().Hex())

	_ = instance //to be used in the next part.
}
