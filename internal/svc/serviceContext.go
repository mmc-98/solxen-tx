package svc

import (
	"net"
	"net/http"
	"solxen-tx/internal/config"
	pb "solxen-tx/internal/svc/proto"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/klauspost/compress/gzhttp"
)

var (
	defaultMaxIdleConnsPerHost = 9
	defaultTimeout             = 5 * time.Minute
	defaultKeepAlive           = 180 * time.Second
)

type ServiceContext struct {
	Config     config.Config
	Lock       sync.RWMutex
	AddrList   []*solana.Wallet
	SolCli     *rpc.Client
	TxnCli     *rpc.Client
	GrpcCli    pb.GeyserClient
	Blockhash  chan solana.Hash
	HTTPClient *http.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:   c,
		Lock:     sync.RWMutex{},
		AddrList: make([]*solana.Wallet, 0),
		SolCli:   rpc.New(c.Sol.Url),
		// TxnCli:    rpc.New(c.Sol.TxnUrl),
		// GrpcCli:   NewGrpcCli(c.Sol.GrpcUrl, true),
		Blockhash:  make(chan solana.Hash, 10),
		HTTPClient: newHTTP(),
	}
	s.GenKeyByWord()
	return s
}

func newHTTP() *http.Client {
	tr := newHTTPTransport()

	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: gzhttp.Transport(tr),
	}
}

func newHTTPTransport() *http.Transport {
	return &http.Transport{
		IdleConnTimeout:     defaultTimeout,
		MaxConnsPerHost:     defaultMaxIdleConnsPerHost,
		MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
		Proxy:               http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   defaultTimeout,
			KeepAlive: defaultKeepAlive,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2: true,
		// MaxIdleConns:          100,
		TLSHandshakeTimeout: 10 * time.Second,
		// ExpectContinueTimeout: 1 * time.Second,
	}
}
