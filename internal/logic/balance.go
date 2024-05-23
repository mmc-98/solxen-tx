package logic

import (
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/logx"
	"solxen-tx/internal/logic/generated/sol_xen_minter"
)

func (l *Producer) Balance() {

	for _, account := range l.svcCtx.AddrList {
		// 获取账户余额
		balance, err := l.svcCtx.SolCli.GetBalance(l.ctx, account.PublicKey(), rpc.CommitmentConfirmed)
		if err != nil {
			logx.Errorf("failed to get balance for account %v: %v", account.PublicKey(), err)
			continue
		}
		balanceInSOL := float64(balance.Value) / 1_000_000_000

		// 获取tokens
		var user_balance_data_raw sol_xen_minter.UserTokensRecord

		// 打印账户余额和 tokens 信息
		logx.Infof("account: %v  balance: %.7f SOL  tokens: %v",
			account.PublicKey(),
			balanceInSOL,
			user_balance_data_raw.TokensMinted,
		)
	}
}
