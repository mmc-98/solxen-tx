package model

import (
	"github.com/lmittmann/w3"
	"github.com/zeromicro/go-zero/core/logx"
)

type XenContractModel struct {
	customContractModel
	signature string
	returns   string
}

func NewXenContractModel() *XenContractModel {
	return &XenContractModel{
		// signature: "bulkClaimRank(uint256 arg0, uint256 arg1)",
		signature: "bulkClaimMintReward(uint256 arg0, address arg1)",
		returns:   "",
	}
}

func (l XenContractModel) EncodeArgs(args ...any) ([]byte, error) {
	logx.Infof("xen...")
	return w3.MustNewFunc(l.signature, l.returns).EncodeArgs(args...)
}
