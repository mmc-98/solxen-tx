package logic

import (
	"context"
	"fmt"
	"log"
	"solxen-tx/internal/config"
	"solxen-tx/internal/svc"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type Producer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	all        int
	CounterPda solana.PublicKey
}

func NewProducerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Producer {
	return &Producer{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		all:        0,
		CounterPda: NewCounterPda(svcCtx.Config),
	}
}

func (l *Producer) Start() {
	logx.Infof("start  mint")
	for {
		// 1. CheckAddressBalance
		// err := l.CheckAddressBalance()
		// if err != nil {
		// 	logx.Errorf("%v", err)
		// 	return
		// }
		// todo 2.QueryNetWorkGas
		// err = l.QueryNetWorkGas()
		// if err != nil {
		// 	logx.Errorf("%v", err)
		// 	return
		// }

		// // 3.mint
		err := l.Mint()
		if err != nil {
			logx.Errorf("Mint err:%v", err)
			continue
		}

		time.Sleep(time.Duration(l.svcCtx.Config.Sol.Time) * time.Millisecond)
	}

}

func (l *Producer) SendTxByAddrList() error {
	t := time.Now()
	limit := computebudget.NewSetComputeUnitLimitInstruction(1200000).Build()
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().SetMicroLamports(l.svcCtx.Config.Sol.Fee).Build()

	// var requests jsonrpc.RPCRequests
	var funcs []func() error
	for _, _payerAccount := range l.svcCtx.AddrList {

		payerAccount := _payerAccount
		var accounts solana.AccountMetaSlice
		accounts = solana.AccountMetaSlice{
			&solana.AccountMeta{PublicKey: payerAccount.PublicKey(), IsSigner: true, IsWritable: true},
			&solana.AccountMeta{PublicKey: l.CounterPda, IsSigner: false, IsWritable: true},
			&solana.AccountMeta{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		}

		var fromAddr string
		if common.IsHexAddress(l.svcCtx.Config.Sol.ToAddr) {
			fromAddr = l.svcCtx.Config.Sol.ToAddr[2:]
		}

		data := common.Hex2BytesFixed(fmt.Sprintf("1800000014000000%v", fromAddr), 28)

		programPubKey, _ := solana.PublicKeyFromBase58(l.svcCtx.Config.Sol.ProgramID)

		instruction := solana.NewInstruction(programPubKey, accounts, data)

		recent, err := l.svcCtx.SolCli.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
		hash := recent.Value.Blockhash

		logx.Infof("l.Blockhash :%v", hash)

		tx, err := solana.NewTransactionBuilder().
			AddInstruction(feesInit).
			AddInstruction(limit).
			AddInstruction(instruction).
			SetRecentBlockHash(hash).
			SetFeePayer(payerAccount.PublicKey()).
			Build()
		if err != nil {
			logx.Errorf("err :%+v", err)
		}

		// sign transaction
		signers := []solana.PrivateKey{payerAccount.PrivateKey}

		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				for _, signer := range signers {
					if signer.PublicKey().Equals(key) {
						return &signer
					}
				}

				return nil
			},
		)
		if err != nil {
			logx.Errorf("err :%v", err)
		}

		// tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))
		// send and confirm transaction
		// sig, err := confirm.SendAndConfirmTransaction(
		// 	ctx,
		// 	rpcCli,
		// 	wsCli,
		// 	tx,
		// )
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// logx.Infof("sig :%v", sig.String())
		//
		// _ = wsCli

		funcs = append(funcs, func() error {
			// logx.Infof(" ==== start ==== %v", time.Since(t))
			sig, err := l.svcCtx.SolCli.SendTransaction(l.ctx, tx)
			if err != nil {
				log.Fatalln(err)
			}
			logx.Infof("sig: %v hash: %v t: %v", sig.String(), hash, time.Since(t))
			return nil
		})
		// txData, err := tx.MarshalBinary()
		// if err != nil {
		// 	logx.Errorf("send transaction: encode transaction: %w", err)
		// }

		//
		// opts := rpc.TransactionOpts{
		// 	SkipPreflight:       false,
		// 	PreflightCommitment: "",
		// }
		//
		// obj := opts.ToMap()
		// params := []interface{}{
		// 	base64.StdEncoding.EncodeToString(txData),
		// 	obj,
		// }
		//
		// requests = append(requests, jsonrpc.NewRequest("sendTransaction", params))

	}

	//
	// jsonRpcClient := jsonrpc.NewClient(rpc.DevNet_RPC)
	// client := rpc.NewWithCustomRPCClient(jsonRpcClient)
	//
	// sig, err := client.RPCCallBatch(ctx, requests)
	// if err != nil {
	// 	logx.Errorf("client.Call err:%v", err)
	// }
	//
	// for _, item := range sig {
	// 	var signature solana.Signature
	// 	err := json.Unmarshal(item.Result, &signature)
	// 	if err != nil {
	// 		logx.Errorf("err :%v", err)
	// 	}
	// 	logx.Infof("sig: %v", signature.String())
	//
	// }
	err := mr.Finish(funcs...)
	if err != nil {
		logx.Errorf("err :%v", err)
	}

	return nil
}

func (l *Producer) Stop() {
	logx.Infof("stop Producer \n")
}

func NewCounterPda(config config.Config) solana.PublicKey {
	programPubKey, _ := solana.PublicKeyFromBase58(config.Sol.ProgramID)

	var fromAddr string
	if common.IsHexAddress(config.Sol.ToAddr) {
		fromAddr = config.Sol.ToAddr[2:]
	}

	seed := [][]byte{
		common.FromHex(fromAddr),
	}
	counterPda, _, err := solana.FindProgramAddress(seed, programPubKey)
	if err != nil {
		logx.Errorf("err :%v", err)
	}

	return counterPda
}
