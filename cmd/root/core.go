package root

import (
	"solxen-tx/cmd/airdrop"
	"solxen-tx/cmd/miner"
	"solxen-tx/cmd/balance"
	minter "solxen-tx/cmd/mint"

	"github.com/spf13/cobra"
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
		RunE:  balance.Balance,
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
	Balance.Flags().StringVar(balance.ConfigFile, "f4", "solxen-tx.yaml", "the config file")
}
