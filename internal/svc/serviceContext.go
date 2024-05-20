package svc

import (
	"solxen-tx/internal/config"
	pb "solxen-tx/internal/svc/proto"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type ServiceContext struct {
	Config    config.Config
	Lock      sync.RWMutex
	AddrList  []*solana.Wallet
	SolCli    *rpc.Client
	TxnCli    *rpc.Client
	GrpcCli   pb.GeyserClient
	Blockhash chan solana.Hash
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:   c,
		Lock:     sync.RWMutex{},
		AddrList: make([]*solana.Wallet, 0),
		SolCli:   rpc.New(c.Sol.Url),
		// TxnCli:    rpc.New(c.Sol.TxnUrl),
		// GrpcCli:   NewGrpcCli(c.Sol.GrpcUrl, true),
		Blockhash: make(chan solana.Hash, 10),
	}
	s.GenKeyByWord()
	return s
}
