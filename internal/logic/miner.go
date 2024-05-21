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
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

func (l *Producer) Miner() error {

	// out := make([]rpc.PriorizationFeeResult, 0)
	// feeAccount := []solana.PublicKey{
	// 	// solana.MustPublicKeyFromBase58("41twqNJmPHv8a5AW32if2CcGRcPzaetwErXaNggGWu1q"),
	// 	solana.MustPublicKeyFromBase58("7LBe4g8Q6hq8Sk1nT8tQUiz2mCHjsoQJbmZ7zCQtutuT"),
	// }
	// fee := l.svcCtx.Config.Sol.Fee
	// if fee == 0 {
	// 	out, _ = l.svcCtx.SolCli.GetRecentPrioritizationFees(l.ctx, feeAccount)
	// }
	// _ = out
	// spew.Dump(out)
	// return nil

	var (
		fns []func() error
		// ProgramIdMiner = solana.MustPublicKeyFromBase58(l.svcCtx.Config.Sol.ProgramIdMiner)
		// ProgramIdMiner = solana.PublicKeySlice{
		// 	solana.MustPublicKeyFromBase58("H4Nk2SDQncEv5Cc6GAbradB4WLrHn7pi9VByFL9zYZcA"),
		// 	solana.MustPublicKeyFromBase58("58UESDt7K7GqutuHBYRuskSgX6XoFe8HXjwrAtyeDULM"),
		// 	solana.MustPublicKeyFromBase58("B1Dw79PE8dzpHPKjiQ8HYUBZ995hL1U32bUTRdNVtRbr"),
		// 	solana.MustPublicKeyFromBase58("7ukQWD7UqoC61eATrBMrdfMrJMUuY1wuPTk4m4noZpsH"),
		// }
		limit    = computebudget.NewSetComputeUnitLimitInstruction(1150000).Build()
		feesInit = computebudget.NewSetComputeUnitPriceInstructionBuilder().SetMicroLamports(l.svcCtx.Config.Sol.Fee).Build()
	)
	ethAccount := common.HexToAddress(l.svcCtx.Config.Sol.ToAddr)
	var uint8Array [20]uint8
	copy(uint8Array[:], ethAccount[:])
	eth := sol_xen_miner.EthAccount{}
	eth.Address = uint8Array
	eth.AddressStr = ethAccount.String()

	for _index, _account := range l.svcCtx.AddrList {
		account := _account
		index := _index
		kind := index % 4
		fns = append(fns, func() error {

			// logx.Infof("account:%v kind:%v ProgramIdMiner[kind]:%v", account.PublicKey(), kind, ProgramIdMiner[kind])

			t := time.Now()
			// global_xn_record_pda
			globalXnRecordPda, _, err := solana.FindProgramAddress(
				[][]byte{
					[]byte("xn-miner-global"),
					{uint8(kind)},
				},
				l.ProgramIdMiner[kind])
			if err != nil {
				return errorx.Wrap(err, "global_xn_record_pda")
			}
			// user_eth_xn_record_pda
			var (
				fromAddr string
			)
			if common.IsHexAddress(l.svcCtx.Config.Sol.ToAddr) {
				fromAddr = l.svcCtx.Config.Sol.ToAddr[2:]
			}

			userEthXnRecordPda, _, err := solana.FindProgramAddress(
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

			userSolXnRecordPda, _, err := solana.FindProgramAddress(
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
			rent := recent.Value.Blockhash
			tx, err := solana.NewTransactionBuilder().
				AddInstruction(feesInit).
				AddInstruction(limit).
				AddInstruction(instruction).
				SetRecentBlockHash(rent).
				SetFeePayer(account.PublicKey()).
				Build()
			if err != nil {
				return errorx.Wrap(err, "tx")
			}

			signers := []solana.PrivateKey{account.PrivateKey}
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
				return errorx.Wrap(err, "Sign")
			}

			var (
				userAccountDataRaw    sol_xen_miner.UserEthXnRecord
				userSolAccountDataRaw sol_xen_miner.UserSolXnRecord
				signature             solana.Signature
			)
			err = mr.Finish(
				func() error {
					signature, err = l.svcCtx.SolCli.SendTransactionWithOpts(context.TODO(), tx, rpc.TransactionOpts{
						SkipPreflight: false,
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
						logx.Infof("userAccountDataRaw:%v", err)
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
						logx.Infof("userSolAccountDataRaw:%v", err)
						return nil
					}
					return nil
				},
			)
			if err != nil {
				return err
			}

			logx.Infof("account:%v slot:%v kind:%v hashs:%v superhashes:%v Points:%v t:%v",
				// signature,
				account.PublicKey(),
				recent.Context.Slot,
				index,
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
