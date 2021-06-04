package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	store "goeth/smart_contracts/02deploying_a_smart_contract/contracts"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatalln(err)
	}

	// Writing transactions requires signing the transaction with our private key.
	privateKey, err := crypto.HexToECDSA("KEYGOESHERE")
	if err != nil {
		log.Fatalln(err)
	}

	publickKey := privateKey.Public()
	publicKeyECDSA, ok := publickKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("can't asser type: publicKey is not of type *ecdsa.PublicKey")
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

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// Create new keyed transactor, taking the private key and set
	// the standard transaction options attached to the keyed transactor.
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalln(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	// Load the smart contract's instance by calling Store's NewStore method.
	address := common.HexToAddress("0x4ccd4Da70E38ec06fC6Eac646E1b2d3e511CC39C")
	instance, err := store.NewStore(address, client)
	if err != nil {
		log.Fatalln(err)
	}

	// Call our contract's external method 'SetItem', passing
	// the auth item we created above, and an array of 32 bytes.
	// This method encodes the function call with it's args:
	// setting it as data of the transaction, and signing with key.
	// It returns a signed transaction object.
	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

	// Verify that the transaction has been sent to the network
	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(result[:]))

}