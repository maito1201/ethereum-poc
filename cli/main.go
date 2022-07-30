package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/maito1201/ethrium-poc/cli/coin"
	"github.com/maito1201/ethrium-poc/cli/nft"
)

// Blockchain const
const (
	host             = "http://localhost"
	port             = "8545"
	coinContractAddr = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	nftContractAddr  = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
	user1PrvKey      = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	user2PrvKey      = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
	user1PubKey      = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	nftURI           = "http://localhost:3000/nft.json"
)

var amount = big.NewInt(100)

func main() {
	if err := tea.NewProgram(newModel()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func initCoinClient(account string) *coin.Client {
	config := coin.Config{
		Host:         host,
		Port:         port,
		ContractAddr: coinContractAddr,
		SignerPrvkey: account,
	}
	cl, err := coin.InitClient(config)
	if err != nil {
		panic(err)
	}
	return cl
}

func mintCoin(client *coin.Client) {
	_, err := client.Instance.Mint(client.Auth, client.Auth.From, amount)
	if err != nil {
		panic(err)
	}
}

func transferCoinToFrined(client *coin.Client) {
	_, err := client.Instance.Transfer(client.Auth, toFriendAddress(client.Auth.From.String()), amount)
	if err != nil {
		panic(err)
	}
}

func transferCoinToMe(client *coin.Client) {
	_, err := client.Instance.TransferFrom(client.Auth, toFriendAddress(client.Auth.From.String()), client.Auth.From, amount)
	if err != nil {
		panic(err)
	}
}

func initNftClient(account string) *nft.Client {
	config := nft.Config{
		Host:         host,
		Port:         port,
		ContractAddr: nftContractAddr,
		SignerPrvkey: account,
	}
	cl, err := nft.InitClient(config)
	if err != nil {
		panic(err)
	}
	return cl
}

func mintNft(client *nft.Client) {
	_, err := client.Instance.SafeMint(client.Auth, client.Auth.From, nftURI)
	if err != nil {
		panic(err)
	}
}

func transferNftToFrined(client *nft.Client, tokenID *big.Int) {
	_, err := client.Instance.TransferFrom(client.Auth, client.Auth.From, toFriendAddress(client.Auth.From.String()), tokenID)
	if err != nil {
		panic(err)
	}
}

func transferNftToMe(client *nft.Client, tokenID *big.Int) {
	_, err := client.Instance.TransferFrom(client.Auth, toFriendAddress(client.Auth.From.String()), client.Auth.From, tokenID)
	if err != nil {
		panic(err)
	}
}

func toFriendAddress(account string) common.Address {
	if account == user1PubKey {
		account = user2PrvKey
	} else {
		account = user1PrvKey
	}
	privateKey, err := crypto.HexToECDSA(account)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("casting public key to ECDSA error")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
