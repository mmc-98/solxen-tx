package logic

import (
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *Producer) Airdrop() {
	for _, account := range l.svcCtx.AddrList {
		out, err := l.svcCtx.SolWriteCli.RequestAirdrop(
			l.ctx,
			account.PublicKey(),
			solana.LAMPORTS_PER_SOL*100,
			rpc.CommitmentFinalized,
		)
		if err != nil {
			logx.Errorf("err :%v", err)
		}
		balance, err := l.svcCtx.SolReadCli.GetBalance(l.ctx, account.PublicKey(), rpc.CommitmentConfirmed)

		logx.Infof("signature: %v account:%v  amount:%v    before:%v", out.String(), account.PublicKey(), 100, balance.Value)
		time.Sleep(1)
	}

	for _, accout := range l.svcCtx.AddrList {
		balance, err := l.svcCtx.SolReadCli.GetBalance(l.ctx, accout.PublicKey(), rpc.CommitmentConfirmed)
		if err != nil {
			return
		}
		logx.Infof("account :%v amount:%v    ", accout.PublicKey(), balance.Value)
	}

}
