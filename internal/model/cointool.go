package model

import (
	"math/big"

	"github.com/lmittmann/w3"
	"github.com/zeromicro/go-zero/core/logx"
)

type CoinToolContractModel struct {
	customContractModel
	signature string
	returns   string
}

func NewCoinToolContractModel() *CoinToolContractModel {
	return &CoinToolContractModel{
		signature: "t(uint256 arg0, bytes arg1, bytes arg2)",
		returns:   "",
	}
}

func (l CoinToolContractModel) EncodeArgs(args ...any) ([]byte, error) {
	logx.Infof("CoinTool...")
	return w3.MustNewFunc(l.signature, l.returns).EncodeArgs(big.NewInt(10), big.NewInt(1))
}
