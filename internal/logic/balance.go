package logic

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"math/big"
	"solxen-tx/internal/logic/generated/sol_xen_miner"
	"solxen-tx/internal/logic/generated/sol_xen_minter"
)

func (l *Producer) Balance() error {
	for _index, _account := range l.svcCtx.AddrList {
		account := _account
		index := _index
		kind := index % 4
		// 获取账户余额
		balance, err := l.svcCtx.SolCli.GetBalance(l.ctx, account.PublicKey(), rpc.CommitmentConfirmed)
		if err != nil {
			logx.Errorf("failed to get balance for account %v: %v", account.PublicKey(), err)
			continue
		}
		balanceInSOL := float64(balance.Value) / 1_000_000_000

		// 获取points,tokens
		var (
			userSolXnRecordPda    solana.PublicKey
			user_balance_data_raw sol_xen_minter.UserTokensRecord
			userSolAccountDataRaw sol_xen_miner.UserSolXnRecord
		)
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
		err = l.svcCtx.SolCli.GetAccountDataInto(
			l.ctx,
			userSolXnRecordPda,
			&userSolAccountDataRaw,
		)
		if err != nil {
			logx.Infof("failed to get account data: %v", err)
			return err
		}
		// 打印账户余额和 tokens 信息
		logx.Infof("account: %v  balance: %.7f SOL  points : %v   tokens: %v",
			account.PublicKey(),
			balanceInSOL,
			big.NewInt(0).Div(userSolAccountDataRaw.Points.BigInt(), big.NewInt(1_000_000_000)),
			user_balance_data_raw.TokensMinted,
		)
	}
	return nil
}
