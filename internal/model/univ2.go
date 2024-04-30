package model

import (
	"github.com/lmittmann/w3"
	"github.com/zeromicro/go-zero/core/logx"
)

type SwapContractModel struct {
	*customContractModel
	signature string
	returns   string
}

func NewSwapContractModel() *SwapContractModel {
	return &SwapContractModel{
		signature: "swapExactETHForTokens(uint amountOutMin, address[]   path, address to, uint deadline)",
		returns:   "uint[]   amounts",
	}
}

func (l SwapContractModel) EncodeArgs(args ...any) ([]byte, error) {
	logx.Infof("swap...")
	// funcSwap := w3.MustNewFunc(`swapExactETHForTokens(uint amountOutMin, address[]   path, address to, uint deadline)`,
	//
	//	"uint[]   amounts")
	//
	// amountOutMin := w3.I("1")
	// var path []common.Address
	// path = append(path, w3.A("0xa15bb66138824a1c7167f5e85b957d04dd34e468"))
	// path = append(path, w3.A("0xf7cd8fa9b94db2aa972023b379c7f72c65e4de9d"))
	// deadline := big.NewInt(time.Now().Unix() + 600).String()
	return w3.MustNewFunc(l.signature, l.returns).EncodeArgs()
}
