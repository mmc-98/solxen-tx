package root

import (
	"solxen-tx/core/airdrop"
	"solxen-tx/core/miner"
	minter "solxen-tx/core/mint"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/api/apigen"
)

var (
	// Miner Cmd describes an api command.
	Miner = &cobra.Command{
		Use:   "miner",
		Short: "miner",
		RunE:  miner.Miner,
	}

	Minter = &cobra.Command{
		Use:   "minter",
		Short: "minter",
		RunE:  minter.Minter,
	}

	Balance = &cobra.Command{
		Use:   "balance",
		Short: "balance",
		RunE:  apigen.CreateApiTemplate,
	}
	Airdrop = &cobra.Command{
		Use:   "airdrop",
		Short: "airdrop",
		RunE:  airdrop.Airdrop,
	}
)

func init() {
	Miner.Flags().StringVar(miner.ConfigFile, "f", "solxen-tx.yaml", "the config file")
	Minter.Flags().StringVar(minter.ConfigFile, "f2", "solxen-tx.yaml", "the config file")
	Airdrop.Flags().StringVar(airdrop.ConfigFile, "f3", "solxen-tx.yaml", "the config file")
}
