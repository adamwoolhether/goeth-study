// A smart contract's Application Binary Interface (ABI)
// must be generated to interact with a smart contract
// in a Go application. It must be compiled and imported
// into the Go application. This requires the Solidity compiler.

// Here we use 0.4.24 to create a very simple smart contract,
// a key-value store with only 1 external method to set a k/v
// pair by anyone.
pragma solidity ^0.4.24;

contract Store {
    event ItemSet(bytes32 key, bytes32 value);

    string public version;
    mapping (bytes32 => bytes32) public items;

    constructor(string _version) public {
        version = _version;
    }

    function setItem(bytes32 key, bytes32 value) external {
        items[key] = value;
        emit ItemSet(key, value);
    }
}

/*
// Generate the ABI from this .sol file,
// then convert to a Go file to import:

// To deploy a smart contract from Go, first compile the
// solidity smart contract to EVM bytecode. A require bin
// file will be generated for generating deploy methods
// for the Go contract file.

// The compile the Go contract file:

solc --abi Store.sol -o build
abigen --abi=./build/Store.abi --pkg=store --out=Store.go

solc --bin Store.sol -o build
abigen --bin=./build/Store.bin --abi=./build/Store.abi --pkg=store --out=Store.go
*/


/*
// if system solidity compiler different than above....

docker pull ethereum/solc:0.4.24
docker run --rm -v $(pwd):/root ethereum/solc:0.4.24 --abi /root/Store.sol -o /root/build
docker run --rm -v $(pwd):/root ethereum/solc:0.4.24 --bin /root/Store.sol -o /root/build
*/