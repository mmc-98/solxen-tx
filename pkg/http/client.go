package httpclient

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/klauspost/compress/gzhttp"
)

var (
	defaultMaxIdleConnsPerHost = 9
	defaultTimeout             = 5 * time.Minute
	defaultKeepAlive           = 180 * time.Second
)

func NewWithProxy(rpcEndpoint string, proxyUrl string) *rpc.Client {
	// 解析代理URL
	proxyUrll, err := url.Parse(proxyUrl)
	if err != nil {
		panic(err)
	}
	opts := &jsonrpc.RPCClientOpts{
		HTTPClient: NewHTTPWithProxy(proxyUrll),
	}
	rpcClient := jsonrpc.NewClientWithOpts(rpcEndpoint, opts)
	return rpc.NewWithCustomRPCClient(rpcClient)
}

func NewHTTP() *http.Client {
	tr := newHTTPTransport()

	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: gzhttp.Transport(tr),
	}
}

func NewHTTPWithProxy(proxy *url.URL) *http.Client {
	tr := newHTTPTransportWithProxy(proxy)

	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: gzhttp.Transport(tr),
	}
}

func newHTTPTransportWithProxy(proxy *url.URL) *http.Transport {
	return &http.Transport{
		IdleConnTimeout:     defaultTimeout,
		MaxConnsPerHost:     defaultMaxIdleConnsPerHost,
		MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
		Proxy:               http.ProxyURL(proxy),
		DialContext: (&net.Dialer{
			Timeout:   defaultTimeout,
			KeepAlive: defaultKeepAlive,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2: true,
		// MaxIdleConns:          100,
		TLSHandshakeTimeout: 10 * time.Second,
		// ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives: true,
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
		// MaxIdleConns:        100,
		TLSHandshakeTimeout: 10 * time.Second,
		// ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives: true,
	}
}
