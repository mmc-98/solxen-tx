package logic

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"solxen-tx/internal/logic/generated/sol_xen_miner"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

var (
	JitoBundleUrl = []string{"https://ny.mainnet.block-engine.jito.wtf/api/v1/bundles"} // "https://mainnet.block-engine.jito.wtf/api/v1/bundles",
	// "https://amsterdam.mainnet.block-engine.jito.wtf/api/v1/bundles",
	// "https://frankfurt.mainnet.block-engine.jito.wtf/api/v1/bundles",
	// "https://ny.mainnet.block-engine.jito.wtf/api/v1/bundles",
	// "https://tokyo.mainnet.block-engine.jito.wtf/api/v1/bundles",

	JitoTxUrl = []string{"https://ny.mainnet.block-engine.jito.wtf/api/v1/transactions"} // "https://mainnet.block-engine.jito.wtf/api/v1/transactions",
	// "https://amsterdam.mainnet.block-engine.jito.wtf/api/v1/transactions",
	// "https://frankfurt.mainnet.block-engine.jito.wtf/api/v1/transactions",
	// "https://ny.mainnet.block-engine.jito.wtf/api/v1/transactions",
	// "https://tokyo.mainnet.block-engine.jito.wtf/api/v1/transactions",

)

func (l *Producer) Miner() error {

	var (
		fns   []func() error
		limit = computebudget.NewSetComputeUnitLimitInstruction(1150000).Build()
	)
	ethAccount := common.HexToAddress(l.svcCtx.Config.Sol.ToAddr)
	var uint8Array [20]uint8
	copy(uint8Array[:], ethAccount[:])
	eth := sol_xen_miner.EthAccount{}
	eth.Address = uint8Array
	eth.AddressStr = ethAccount.String()

	// out := make([]rpc.PriorizationFeeResult, 0)
	// feeAccount := []solana.PublicKey{
	// 	solana.MustPublicKeyFromBase58(l.svcCtx.Config.Sol.ProgramId),
	// }

	// fee := l.svcCtx.Config.Sol.Fee
	// if fee == 0 {
	// 	out, _ = l.svcCtx.SolCli.GetRecentPrioritizationFees(l.ctx, feeAccount)
	// 	var feeFata []float64
	// 	for _, item := range out {
	// 		feeFata = append(feeFata, float64(item.PrioritizationFee))
	// 	}
	// 	_fee, _ := stats.Mean(feeFata)
	// 	fee = uint64(_fee) * 1_000_000
	// }
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().SetMicroLamports(l.svcCtx.Config.Sol.Fee).Build()

	for _index, _account := range l.svcCtx.AddrList {
		account := _account
		index := _index
		kind := index % 4

		kind = l.svcCtx.Config.Sol.Kind
		if kind == -1 {
			kind = index % 4
		} else {
			account = l.svcCtx.AddrList[0]
		}

		fns = append(fns, func() error {

			t := time.Now()
			var (
				err                error
				globalXnRecordPda  solana.PublicKey
				userEthXnRecordPda solana.PublicKey
				userSolXnRecordPda solana.PublicKey
			)
			mr.Finish(
				func() error {
					globalXnRecordPda, _, err = solana.FindProgramAddress(
						[][]byte{
							[]byte("xn-miner-global"),
							{uint8(kind)},
						},
						l.ProgramIdMiner[kind])
					if err != nil {
						return errorx.Wrap(err, "global_xn_record_pda")
					}
					return nil
				},
				func() error {
					var (
						fromAddr string
					)
					if common.IsHexAddress(l.svcCtx.Config.Sol.ToAddr) {
						fromAddr = l.svcCtx.Config.Sol.ToAddr[2:]
					}

					userEthXnRecordPda, _, err = solana.FindProgramAddress(
						[][]byte{
							[]byte("xn-by-eth"),
							common.FromHex(fromAddr),
							{uint8(kind)},
							l.ProgramIdMiner[kind].Bytes(),
						},
						l.ProgramIdMiner[kind])
					if err != nil {
						return errorx.Wrap(err, "userEthXnRecordAccount")
					}
					return nil
				},
				func() error {

					userSolXnRecordPda, _, err = solana.FindProgramAddress(
						[][]byte{
							[]byte("xn-by-sol"),
							account.PublicKey().Bytes(),
							{uint8(kind)},
							l.ProgramIdMiner[kind].Bytes(),
						},
						l.ProgramIdMiner[kind])
					if err != nil {
						return errorx.Wrap(err, "global_xn_record_pda")
					}

					return nil
				},
			)

			mintToken := sol_xen_miner.NewMineHashesInstruction(
				eth,
				uint8(kind),
				globalXnRecordPda,
				userEthXnRecordPda,
				userSolXnRecordPda,
				account.PublicKey(),
				solana.SystemProgramID,
			).Build()

			// l.svcCtx.Lock.Lock()
			// sol_xen_miner.SetProgramID(ProgramIdMiner[kind])
			data, _ := mintToken.Data()
			instruction := solana.NewInstruction(l.ProgramIdMiner[kind], mintToken.Accounts(), data)
			// l.svcCtx.Lock.Unlock()

			recent, err := l.svcCtx.SolCli.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
			if err != nil {
				return errorx.Wrap(err, "network.")
			}
			rent := recent.Value.Blockhash

			tx := solana.NewTransactionBuilder().
				// AddInstruction(feesInit).
				AddInstruction(limit).
				AddInstruction(instruction).
				SetRecentBlockHash(rent).
				SetFeePayer(account.PublicKey())
			if err != nil {
				return errorx.Wrap(err, "tx")
			}
			//
			// rand.Seed(time.Now().UnixNano())
			// n := rand.Intn(5-0+1) + 0
			var rpcClient *rpc.Client
			if l.svcCtx.Config.Sol.JitoTip != 0 {

				rpcClient = rpc.New(JitoTxUrl[0])
				// time.Sleep(1000 * time.Millisecond)

				if len(l.Respd.Result) == 0 {
					err := l.GenJitoAddr()
					if err != nil {
						return err
					}
					return nil
				}
				// logx.Infof("JitoBundleUrl:%v", respd)
				accountTo := solana.MustPublicKeyFromBase58(l.Respd.Result[0])
				transferInstruction := system.NewTransferInstruction(
					l.svcCtx.Config.Sol.JitoTip,
					account.PublicKey(),
					accountTo,
				).Build()
				tx.AddInstruction(transferInstruction)
			} else {
				tx.AddInstruction(feesInit)
				rpcClient = l.svcCtx.SolCli
			}
			txData, err := tx.Build()

			signers := []solana.PrivateKey{account.PrivateKey}
			_, err = txData.Sign(
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
				return errorx.Wrap(err, "Sign")
			}
			var (
				userAccountDataRaw    sol_xen_miner.UserEthXnRecord
				userSolAccountDataRaw sol_xen_miner.UserSolXnRecord
				signature             solana.Signature
			)
			err = mr.Finish(
				func() error {
					signature, err = rpcClient.SendTransactionWithOpts(context.TODO(), txData, rpc.TransactionOpts{
						// _ = rpcClient
						// signature, err = l.svcCtx.TxnCli.SendTransactionWithOpts(context.TODO(), txData, rpc.TransactionOpts{
						SkipPreflight: true,
						MaxRetries:    new(uint),
					})
					_ = signature
					if err != nil {
						return errorx.Wrap(err, "sig")
					}

					return nil
				},

				func() error {
					err = l.svcCtx.SolCli.GetAccountDataInto(
						l.ctx,
						userEthXnRecordPda,
						&userAccountDataRaw,
					)
					if err != nil {
						// logx.Infof("userAccountDataRaw:%v", err)
						return nil
					}
					return nil
				},

				func() error {
					err = l.svcCtx.SolCli.GetAccountDataInto(
						l.ctx,
						userSolXnRecordPda,
						&userSolAccountDataRaw,
					)
					if err != nil {
						// logx.Infof("userSolAccountDataRaw:%v", err)
						return nil
					}
					return nil
				},
			)
			if err != nil {
				return err
			}

			logx.Infof("account:%v fee:%v jito:%v slot:%v kind:%v hashs:%v superhashes:%v Points:%v t:%v",
				account.PublicKey(),
				l.svcCtx.Config.Sol.Fee,
				l.svcCtx.Config.Sol.JitoTip,
				recent.Context.Slot,
				kind,
				// common.Bytes2Hex(maybe_user_account_data_raw.Nonce[:]),
				userAccountDataRaw.Hashes,
				userAccountDataRaw.Superhashes,
				big.NewInt(0).Div(userSolAccountDataRaw.Points.BigInt(), big.NewInt(1_000_000_000)),
				time.Since(t))

			return nil

		})
	}
	err := mr.Finish(fns...)
	if err != nil {
		logx.Errorf("err: %v", err)
	}
	return nil

}

func (l *Producer) CheckAddressBalance() error {

	var (
		fns []func() error
	)
	for _, addr := range l.svcCtx.AddrList {
		fns = append(fns, func() error {
			balance, err := l.svcCtx.SolCli.GetBalance(l.ctx, addr.PublicKey(), rpc.CommitmentFinalized)
			if err != nil {
				return err
			}
			if (balance.Value) < 1_000_000 {
				return errorx.Wrap(err, fmt.Sprintf("%v Balance less than 0.01, please recharge.余额小于0.01请充值", addr.PublicKey()))
			}
			return nil
		})

	}
	err := mr.Finish(
		fns...,
	)
	if err != nil {
		logx.Errorf("err %v", err)
	}
	return nil
}

func (l *Producer) QueryNetWorkGas() error {
	return nil
}
