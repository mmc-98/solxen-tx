package model

var _ ContractModel = (*customContractModel)(nil)

type (
	ContractModel interface {
		EncodeArgs(args ...any) ([]byte, error)
	}

	customContractModel struct {
	}
)

func (c *customContractModel) EncodeArgs(args ...any) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func NewBaseContractModel() *SwapContractModel {
	return &SwapContractModel{
		signature: "",
		returns:   "",
	}
}
