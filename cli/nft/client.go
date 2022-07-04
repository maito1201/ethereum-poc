package nft

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Config struct {
	Host         string
	Port         string
	SignerPrvkey string
	ContractAddr string
}

type Client struct {
	Auth     *bind.TransactOpts
	Instance *Nft
}

func InitClient(config Config) (*Client, error) {
	client, err := ethclient.Dial(fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(config.ContractAddr)
	instance, err := NewNft(address, client)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(config.SignerPrvkey)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("casting public key to ECDSA error")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPrice
	return &Client{Instance: instance, Auth: auth}, nil
}
