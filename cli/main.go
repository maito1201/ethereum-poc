package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/maito1201/ethrium-poc/cli/coin"
)

// Blockchain const
const (
	host             = "http://localhost"
	port             = "8545"
	coinContractAddr = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
	user1PrvKey      = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	user2PrvKey      = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"
)

var amount = big.NewInt(100)

// app config
const listHeight = 14
const defaultWidth = 30

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

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
		SignerPkey:   account,
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

func transferCoin(client *coin.Client) {
	_, err := client.Instance.Transfer(client.Auth, toFriendAddress(string(client.Auth.From.Bytes())), amount)
	if err != nil {
		panic(err)
	}
}

func toFriendAddress(account string) common.Address {
	if account == user1PrvKey {
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
