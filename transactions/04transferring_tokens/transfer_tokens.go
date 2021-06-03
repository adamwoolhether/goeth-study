package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func main() {
	// Same as before, you must connect to a client, load account private key
	// and configure gas price to use for your transaction.
	client, err := ethclient.Dial("https://rinkeby.infura.io/v3/95040d1a4c7748b595f8afb6949bb88d")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("PRIVATE KEY")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Set the ETH value in wei (0 eth)
	value := big.NewInt(0)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x9DeeD15f4164a2d0ef37f6c204eb3965D7e28D46")

	// Sending tokens requires invoking a function on the smart contract.
	// 1. Find the fucntion of the signature of transfer() smart contract function.
	// 2. Figure out inputs for the function (receiient address, tokens value)
	// 3. Get first 8 chars(4 bytes) of Keccak256 hash of that function signature.
	// 4. Zero-pad the left of our function call - the address and value.
	// **this must be 256 bits long.
	tokenAddress := common.HexToAddress("0xaFF4481D10270F50f203E0763e2597776068CBc5") //WEENUS testnet

	// Form the smart contract function call. We pass the receiver's address
	// and the second arg's type (the amount of tokens to send).
	// This function signature needs to be a byte slice.
	transferFnSignature := []byte("transfer(address,uint256)")

	// To get methodID of our function, we use cryto/sha3 to generate the
	// Keccak256 hash of the function signature, taking first 8 chars
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID))

	// Now zero pad the account address we're sending tokens to.
	// The resulting byte slice must be 32 bytes long.
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Printf("paddedAddress: %s\n", hexutil.Encode(paddedAddress))

	// The tokens' value is set as a *big.Int. Denomination is determined by
	// the token contract that you're interacting with. The WEENUS token we
	// use here uses 18 decimals, which is standard ERC-20 practice. To
	// represent 1 token: amt * 10^18.
	// (If the decimal90 value was 0, then big.NewInt(1000) equals 1000 tokens.)
	// Here we send 50 tokens:
	amount := new(big.Int)
	amount.SetString("50000000000000000000", 10)

	// Left padding to 32 bytes also required for the amount per EVM reqs
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Printf("paddedAmount: %s\n", hexutil.Encode(paddedAmount))

	// We then concatenate the method ID, padded address, and padded amt
	// into a byte slice to create the data field.
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// Setting gas limit with EstimateGas function may result in repeated failure.
	//gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	//	To:   &tokenAddress,
	//	Data: data,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	gasLimit := uint64(42000)
	fmt.Printf("gasLimit: %d\n", gasLimit)

	// Creating the transaction is similar to the transferring ETH
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
