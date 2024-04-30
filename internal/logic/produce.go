package logic

import (
	"context"
	"fmt"
	"log"
	"solxen-tx/internal/svc"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/zeromicro/go-zero/core/logx"
)

type Producer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	all int
}

func NewProducerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Producer {
	return &Producer{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		all:    0,
	}
}

func (l *Producer) Start() {
	logx.Infof("start  Producer \n")

	for {
		// // 1.查询余额
		// err := l.CheckAddressBalance()
		// if err != nil {
		// 	logx.Errorf("%v", err)
		// 	os.Exit(-1)
		// }
		// // 2.获取gas和tipGas
		// err = l.QueryFreeGasAndTipGas()
		// if err != nil {
		// 	logx.Errorf("%v", err)
		// 	os.Exit(-2)
		// }
		// // 3.获取将发送tx的地址列表
		// err = l.ListTxpoolPendding()
		// if err != nil {
		// 	logx.Errorf("ListTxpoolPendding err:%v", err)
		// 	continue
		// }
		// // 4.获取nonce
		// err = l.BatchListNoceByAddr()
		// if err != nil {
		// 	continue
		// }
		// // 5.发送tx
		// err = l.SendTxByAddrList()
		// if err != nil {
		// 	logx.Errorf("ListTxpoolPendding err:%v", err)
		// 	continue
		// }
		l.Do()
		time.Sleep(time.Duration(l.svcCtx.Config.Sol.Time) * time.Millisecond)
	}

}

func (l *Producer) Do() {
	ctx := context.Background()
	rpcCli := rpc.New(rpc.DevNet_RPC)
	wsCli, _ := ws.Connect(ctx, rpc.DevNet_WS)
	payerPrivateKey, _ := solana.PrivateKeyFromBase58(l.svcCtx.Config.Sol.Key) // don't worry, I only use this private key in my local computer
	payerAccount, _ := solana.WalletFromPrivateKeyBase58(payerPrivateKey.String())
	programID := l.svcCtx.Config.Sol.ProgramID
	programPubKey, _ := solana.PublicKeyFromBase58(programID)
	// logx.Infof("payerAccount : %+v ", payerAccount.PublicKey())

	var fromAddr string
	if common.IsHexAddress(l.svcCtx.Config.Sol.ToAddr) {
		fromAddr = l.svcCtx.Config.Sol.ToAddr
	}

	seed := [][]byte{
		common.FromHex(fromAddr),
	}
	counterPda, _, err := solana.FindProgramAddress(seed, programPubKey)
	if err != nil {
		logx.Infof("err %v", err)
	}
	// logx.Infof("counter_pda :%+v %+v", counterPda, as)

	// create recent blockhash
	recent, err := rpcCli.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		logx.Errorf("err :%v", err)
	}

	limit := computebudget.NewSetComputeUnitLimitInstruction(1200000).Build()
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().SetMicroLamports(l.svcCtx.Config.Sol.Fee).Build()

	var accounts solana.AccountMetaSlice
	accounts = solana.AccountMetaSlice{
		&solana.AccountMeta{PublicKey: payerAccount.PublicKey(), IsSigner: true, IsWritable: true},
		&solana.AccountMeta{PublicKey: counterPda, IsSigner: false, IsWritable: true},
		&solana.AccountMeta{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
	}
	data := common.Hex2BytesFixed(fmt.Sprintf("1800000014000000%v", fromAddr), 28)

	instruction := solana.NewInstruction(programPubKey, accounts, data)

	tx, err := solana.NewTransactionBuilder().
		AddInstruction(feesInit).
		AddInstruction(limit).
		AddInstruction(instruction).
		SetRecentBlockHash(recent.Value.Blockhash).
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
		log.Fatalln(err)
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
	_ = wsCli
	sig, err := l.svcCtx.SolCli.SendTransaction(ctx, tx)
	if err != nil {
		log.Fatalln(err)
	}
	logx.Infof("sig :%v", sig.String())
}

func (l *Producer) Stop() {
	logx.Infof("stop Producer \n")
}
