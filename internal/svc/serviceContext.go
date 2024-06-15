package svc

import (
	"net/http"
	"solxen-tx/internal/config"
	pb "solxen-tx/internal/svc/proto"
	httpclient "solxen-tx/pkg/http"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	defaultMaxIdleConnsPerHost = 9
	defaultTimeout             = 5 * time.Minute
	defaultKeepAlive           = 180 * time.Second
)

type ServiceContext struct {
	Config      config.Config
	Lock        sync.RWMutex
	AddrList    []*solana.Wallet
	SolReadCli  *rpc.Client
	SolWriteCli *rpc.Client
	// TxnCli     *rpc.Client
	GrpcCli    pb.GeyserClient
	Blockhash  chan solana.Hash
	HTTPClient *http.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:      c,
		Lock:        sync.RWMutex{},
		AddrList:    make([]*solana.Wallet, 0),
		SolReadCli:  rpc.New(c.Sol.Url),
		SolWriteCli: rpc.New(c.Sol.Url),
		// TxnCli:   rpc.New(c.Sol.TxnUrl),
		// GrpcCli:   NewGrpcCli(c.Sol.GrpcUrl, true),
		Blockhash:  make(chan solana.Hash, 10),
		HTTPClient: httpclient.NewHTTP(),
	}
	s.GenKeyByWord()
	return s
}
