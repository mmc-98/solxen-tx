package logic

import (
	"context"
	"solxen-tx/internal/logic/generated/sol_xen_minter"

	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

func (l *Producer) Mint() error {

	var (
		fns            []func() error
		programIdMiner = solana.MustPublicKeyFromBase58(l.svcCtx.Config.Sol.ProgramIdMiner)
		programId      = solana.MustPublicKeyFromBase58(l.svcCtx.Config.Sol.ProgramId)
	)

	mint_pda, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("mint"),
		},
		programId,
	)
	if err != nil {
		return errorx.Wrap(err, "mint_pda")
	}
	// spew.Dump(mint_pda)

	// limit := computebudget.NewSetComputeUnitLimitInstruction(1400000).Build()
	feesInit := computebudget.NewSetComputeUnitPriceInstructionBuilder().SetMicroLamports(l.svcCtx.Config.Sol.Fee).Build()

	for _, _account := range l.svcCtx.AddrList {
		account := _account
		fns = append(fns, func() error {
			// t := time.Now()
			user_sol_xn_record_pda, _, err := solana.FindProgramAddress(
				[][]byte{
					[]byte("xn-by-sol"),
					account.PublicKey().Bytes(),
					[]uint8{0},
					programIdMiner.Bytes(),
				},
				programIdMiner,
			)
			if err != nil {
				return errorx.Wrap(err, "userSolXnRecordPda")
			}
			// spew.Dump(user_sol_xn_record_pda)

			user_token_record_pda, _, err := solana.FindProgramAddress(
				[][]byte{
					[]byte("sol-xen-minted"),
					account.PublicKey().Bytes(),
				},
				programId,
			)
			if err != nil {
				return errorx.Wrap(err, "user_eth_xn_record_pda")
			}

			// spew.Dump(user_token_record_pda)

			associate_token_program := solana.SPLAssociatedTokenAccountProgramID
			user_token_account, _, err := solana.FindAssociatedTokenAddress(account.PublicKey(), mint_pda)
			// spew.Dump(user_token_account)

			mintToken := sol_xen_minter.NewMintTokensInstruction(
				0,
				user_sol_xn_record_pda,
				user_token_record_pda,
				user_token_account,
				account.PublicKey(),
				mint_pda,
				solana.TokenProgramID,
				solana.SystemProgramID,
				associate_token_program,
				programIdMiner,
			).Build()

			sol_xen_minter.SetProgramID(solana.MustPublicKeyFromBase58(l.svcCtx.Config.Sol.ProgramId))

			data, _ := mintToken.Data()
			// spew.Dump(data)
			// logx.Infof("data :%v", data)

			instruction := solana.NewInstruction(mintToken.ProgramID(), mintToken.Accounts(), data)

			// spew.Dump(instruction)

			signers := []solana.PrivateKey{account.PrivateKey}
			recent, err := l.svcCtx.SolCli.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
			rent := recent.Value.Blockhash
			tx, err := solana.NewTransactionBuilder().
				AddInstruction(feesInit).
				// AddInstruction(limit).
				AddInstruction(instruction).
				SetRecentBlockHash(rent).
				SetFeePayer(account.PublicKey()).
				Build()
			if err != nil {
				return errorx.Wrap(err, "tx")
			}

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

			// logx.Infof("tx :%v", tx)

			var (
				user_balance_data_raw sol_xen_minter.UserTokensRecord
			)

			_, err = l.svcCtx.SolCli.SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{
				SkipPreflight: false,
			})
			if err != nil {
				return errorx.Wrap(err, "sig")
			}
			err = mr.Finish(
				func() error {
					err = l.svcCtx.SolCli.GetAccountDataInto(l.ctx, user_token_record_pda, &user_balance_data_raw)
					return nil
				},
			)
			if err != nil {
				return err
			}

			logx.Infof("account:%v tokens:%v ",
				account.PublicKey(),
				user_balance_data_raw.TokensMinted,
			)
			return nil

		})
	}

	err = mr.Finish(fns...)
	if err != nil {
		logx.Errorf("err: %v", err)
	}
	return nil

}
