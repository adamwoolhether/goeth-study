package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	createKs()
	importKs()
}

// createKs invokes NewKeyStore to generate an encrypted wallet private key,
// telling it where to save the keystores.
// NewAccount is caleld to generate the wallet, passing the password.
// go-eth can only generate one wallet key par file.
func createKs() {
	ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())
}

// importKs imports the keystore, also by calling NewKeyStore
// and Uses the Import method, which takes the keystore JSDON
// data as bytes. Import method can take a third arg to change
// the password if desired.
func importKs() {
	file, err := filepath.Glob("./wallets/UTC--*")
	if err != nil {
		panic(err)
	}
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file[0])
	if err != nil {
		log.Fatal(err)
	}

	password := "secret"
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	if err := os.Remove(file[0]); err != nil {
		log.Fatal(err)
	}
}