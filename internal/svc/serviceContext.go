package svc

import (
	"crypto/ecdsa"
	"solxen-tx/internal/config"
	"solxen-tx/internal/model"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/lmittmann/w3"
)

type ServiceContext struct {
	Config        config.Config
	W3Cli         *w3.Client
	Lock          sync.RWMutex
	AddressKey    map[common.Address]*ecdsa.PrivateKey
	AddrList      []common.Address
	ContractModel model.ContractModel
	SolCli        *rpc.Client
}

func NewServiceContext(c config.Config) *ServiceContext {

	s := &ServiceContext{
		Config:        c,
		W3Cli:         w3.MustDial(c.Sol.Url),
		Lock:          sync.RWMutex{},
		AddressKey:    make(map[common.Address]*ecdsa.PrivateKey),
		AddrList:      make([]common.Address, 0),
		ContractModel: model.NewBaseContractModel(),
		SolCli:        rpc.New(c.Sol.Url),
	}
	return s
}
