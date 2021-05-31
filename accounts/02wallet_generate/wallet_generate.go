package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
)

func main() {
	// Generate a new random private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(privateKey)

	// Convert the key into bytes
	privateKeyBytes := crypto.FromECDSA(privateKey)

	// Now convert it to hexadecimal with the hexutil package's Encode method.
	// It takes a byteslice. We strip off '0x' at beginning
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

	// The public key is derived from the private key
	publicKey := privateKey.Public()


	// Convert to hex and strip 0x and first two chars
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot asset type: publicKey is not of type *ecdsaPublicKey")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])

	// Generate the public address using the PubkeyToAddress
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	// Public address is the Keccak-256 hash of the public key.
	// We take take last 40 characters (20 bytes) and prefix it with 0x
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))

}