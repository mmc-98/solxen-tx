package model

import (
	"github.com/zeromicro/go-zero/core/logx"
)

type TxContractModel struct {
	customContractModel
	signature string
	returns   string
}

func NewTxContractModel() *TxContractModel {
	return &TxContractModel{
		signature: "",
		returns:   "",
	}
}

func (l TxContractModel) EncodeArgs(args ...any) ([]byte, error) {
	logx.Infof("Tx...")
	return []byte{}, nil
}
