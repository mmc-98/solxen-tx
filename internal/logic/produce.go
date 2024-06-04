package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"solxen-tx/internal/svc"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type resp struct {
	Jsonrpc string   `json:"jsonrpc"`
	Result  []string `json:"result"`
	Id      int64    `json:"id"`
}

type Producer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	all            int
	mux            sync.RWMutex
	ProgramIdMiner solana.PublicKeySlice
	Respd          *resp
}

func NewProducerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Producer {

	s := &Producer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		mux:    sync.RWMutex{},
		ProgramIdMiner: solana.PublicKeySlice{
			solana.MustPublicKeyFromBase58("B8HwMYCk1o7EaJhooM4P43BHSk5M8zZHsTeJixqw7LMN"),
			solana.MustPublicKeyFromBase58("2Ewuie2KnTvMLwGqKWvEM1S2gUStHzDUfrANdJfu45QJ"),
			solana.MustPublicKeyFromBase58("5dxcK28nyAJdK9fSFuReRREeKnmAGVRpXPhwkZxAxFtJ"),
			solana.MustPublicKeyFromBase58("DdVCjv7fsPPm64HnepYy5MBfh2bNfkd84Rawey9rdt5S"),

			// eat ,,,
			// solana.MustPublicKeyFromBase58("CFRDmC2xPN7K2D8GadHKpcwSAC5YvPzPjbjYA6v439oi"),
			// solana.MustPublicKeyFromBase58("7vQ9pG7MUjkswNkL96XiiYbz3swM9dkqgMEAbgDaLggi"),
			// solana.MustPublicKeyFromBase58("DpLx72BXVhZN6hkA6LKKres3EUKvK36mmh5JaKyaVSYU"),
			// solana.MustPublicKeyFromBase58("7u5D7qPHGZHXQ3nQTeZu5eFKtKGKQWKhJCdM1B3T4Ly4"),
		},
		Respd: &resp{},
	}
	s.GenJitoAddr()
	return s
}

func (l *Producer) GenJitoAddr() error {
	type req struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int64  `json:"id"`
		Method  string `json:"method"`
		Params  string `json:"params"`
	}

	// reqData, err := json.Marshal(req{Jsonrpc: "2.0", Id: 1, Method: "getTipAccounts", Params: ""})
	reqData, err := json.Marshal(&req{Jsonrpc: "2.0", Id: 1, Method: "getTipAccounts", Params: ""})
	if err != nil {
		log.Fatal(err)
	}
	respData, err := l.svcCtx.HTTPClient.Post(
		JitoBundleUrl[0],
		"application/json",
		bytes.NewBuffer(reqData))
	if err != nil {
		return errorx.Wrap(err, "HTTPClient.")
	}
	defer respData.Body.Close()
	BodyData, err := ioutil.ReadAll(respData.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(BodyData, &l.Respd)
	if err != nil {
		return errorx.Wrap(err, "HTTPClient.")
	}
	// if len(respd.Result) == 0 {
	// 	return nil
	// }

	return nil
}

func (l *Producer) Start() {
	logx.Infof("start  miner")

	// var subscription pb.SubscribeRequest
	// subscription.Blocks = make(map[string]*pb.SubscribeRequestFilterBlocks)
	// subscription.Blocks["blocks"] = &pb.SubscribeRequestFilterBlocks{}
	//
	// subscriptionJson, err := json.Marshal(&subscription)
	// if err != nil {
	// 	logx.Errorf("Failed to marshal subscription request: %v", subscriptionJson)
	// }
	// logx.Infof("Subscription request: %s", string(subscriptionJson))
	//
	// stream, err := l.svcCtx.GrpcCli.Subscribe(context.Background())
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// err = stream.Send(&subscription)
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// for {
	// 	resp, err := stream.Recv()
	// 	// timestamp := time.Now().UnixNano()
	// 	if err == io.EOF {
	// 		return
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("Error occurred in receiving update: %v", err)
	// 	}
	// 	block := resp.GetBlock()
	// 	// log.Printf("timestamp: %v block:%v %v %v", timestamp, block.GetBlockhash(), block.GetSlot(), block.GetBlockHeight())
	// 	Blockhash := block.GetBlockhash()
	// 	slot := block.GetSlot()
	// 	if Blockhash == "" || slot == 0 {
	// 		continue
	// 	}
	//
	// 	err = l.Mint()
	// 	if err != nil {
	// 		logx.Errorf("Mint err:%v", err)
	// 		continue
	// 	}
	//
	// }

	for {
		// 1. CheckAddressBalance
		err := l.CheckAddressBalance()
		if err != nil {
			logx.Errorf("%v", err)
			return
		}
		// todo 2.QueryNetWorkGas
		// err = l.QueryNetWorkGas()
		// if err != nil {
		// 	logx.Errorf("%v", err)
		// 	return
		// }

		// 3.Miner

		err = l.Miner()
		if err != nil {
			logx.Errorf("Mint err:%v", err)
			continue
		}

		time.Sleep(time.Duration(l.svcCtx.Config.Sol.Time) * time.Millisecond)
	}

}

func (l *Producer) Stop() {
	logx.Infof("stop Producer \n")
}
