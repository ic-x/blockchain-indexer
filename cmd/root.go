package cmd

import (
	"fmt"
	"os"

	"github.com/ic-x/blockchain-indexer/cmd/run"
	"github.com/ic-x/blockchain-indexer/internal/config"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "indexer",
	Short: "Blockchain Indexer CLI",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config.LoadConfig()

	RootCmd.AddCommand(run.RunCmd)
}
