package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"
	"os"
	"solxen-tx/internal/config"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/okx/go-wallet-sdk/example"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config       config.Config
	W3Cli        *w3.Client
	Lock         sync.RWMutex
	AddressKey   map[common.Address]*ecdsa.PrivateKey
	AddrList     []common.Address
	FromNocesMap map[common.Address]uint64
	chainID      uint64
}

func (s *ServiceContext) SetAddrKeyAndAddrList() {

	for i := 0; i < s.Config.Eth.Num; i++ {
		hdPath := example.GetDerivedPath(i)
		derivePrivateKey, _ := example.GetDerivedPrivateKey(s.Config.Eth.Key, hdPath)
		address := example.GetNewAddress(derivePrivateKey)

		privateKeyByte, err := hexutil.Decode(fmt.Sprintf("0x%v", derivePrivateKey))
		if err != nil {
			logx.Errorf("err:%v", err)
		}

		privateKey, err := crypto.ToECDSA(privateKeyByte)
		// logx.Infof("addr: %v", address)

		s.AddressKey[common.HexToAddress(address)] = privateKey
		s.AddrList = append(s.AddrList, common.HexToAddress(address))
	}

}

func NewServiceContext(c config.Config) *ServiceContext {

	s := &ServiceContext{
		Config:       c,
		W3Cli:        w3.MustDial(c.Eth.Url),
		Lock:         sync.RWMutex{},
		AddressKey:   make(map[common.Address]*ecdsa.PrivateKey),
		AddrList:     make([]common.Address, 0),
		FromNocesMap: make(map[common.Address]uint64),
	}
	s.SetAddrKeyAndAddrList()
	return s
}

func (s *ServiceContext) SetFromAddresMap() {

	var nonce uint64
	fromAddress := s.AddrList[0]
	err := s.W3Cli.Call(eth.Nonce(fromAddress, nil).Returns(&nonce),
		eth.ChainID().Returns(&s.chainID))
	if err != nil {
		logx.Errorf("client.NonceAt err:%v", err)
	}

	s.FromNocesMap[fromAddress] = nonce
}

func (s *ServiceContext) SendEthtoAll() {

	value := w3.I(s.Config.Eth.Value)
	tipGasBigInt := w3.I("10 gwei")
	freeGasBigInt := w3.I("30 gwei")
	signKey := s.AddressKey[s.AddrList[0]]
	fromAddr := s.AddrList[0]

	for _, addr := range s.AddrList[1:] {
		toAddr := addr
		noce := s.FromNocesMap[s.AddrList[0]]
		signer := types.LatestSignerForChainID(big.NewInt(int64(s.chainID)))
		tx := types.MustSignNewTx(signKey, signer, &types.DynamicFeeTx{
			Nonce:     noce,
			GasFeeCap: freeGasBigInt,
			GasTipCap: tipGasBigInt,
			Gas:       21000,
			To:        &toAddr,
			Value:     value,
		})
		var txHash common.Hash
		err := s.W3Cli.Call(eth.SendTx(tx).Returns(&txHash))
		if err != nil {
			logx.Errorf("err:%v fromAddr:%v toAddr:%v", err, fromAddr, toAddr)
			os.Exit(-1)
		}
		s.FromNocesMap[fromAddr]++
		logx.Infof("fromAddr:%v toAddr:%v noce:%v value: %v", fromAddr, toAddr, noce, s.Config.Eth.Value)

		time.Sleep(1 * time.Second)
	}

}

var configFile = flag.String("f", "build/etc/solxen-tx.yaml", "the config file")

func main() {
	flag.Parse()
	logx.DisableStat()
	// 配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	s := NewServiceContext(c)
	s.SetFromAddresMap()
	s.SendEthtoAll()

}
