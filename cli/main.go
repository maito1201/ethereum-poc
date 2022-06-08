package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/maito1201/ethrium-poc/cli/poll"
)

const addr = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
const url = "http://127.0.0.1:8545/"

// private key of first user
const pkey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// 1 is agree, 2 is disagree
var voteNum = big.NewInt(2)

var instance *poll.Poll
var fromAddress common.Address
var auth *bind.TransactOpts

func init() {
	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	address := common.HexToAddress(addr)
	instance, err = poll.NewPoll(address, client)
	if err != nil {
		panic(err)
	}

	privateKey, err := crypto.HexToECDSA(pkey)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}
	fromAddress = crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	auth.GasPrice = gasPrice
}

func main() {
	// execute contract
	title, err := instance.GetPoll(nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Get Poll: %v\n", title)

	// check vote result
	opts := bind.CallOpts{From: fromAddress, Context: context.Background()}
	v, err := instance.CheckPoll(&opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Check Poll: %v\n", v)

	// vote if result is 0
	if v.Cmp(big.NewInt(0)) == 0 {
		tx, err := instance.Vote(auth, voteNum)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Vote %+v\n", tx)

		// recheck vote result
		v, err := instance.CheckPoll(&opts)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Check Poll: %v\n", v)
	}
}
