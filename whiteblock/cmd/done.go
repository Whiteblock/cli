package cmd

import (
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
)

var doneCmd = &cobra.Command{
	Aliases: []string{"die", "stop", "teardown", "purge"},
	Use:     "done",
	Short:   "Tears down the testnet",
	Long: `
	Tears down the nodes, and frees up any resources which they are using.
`,

	Run: func(cmd *cobra.Command, args []string) {
		util.JsonRpcCallAndPrint("delete_testnet", []interface{}{build.GetPreviousBuildID()})
	},
}

func init() {
	RootCmd.AddCommand(doneCmd)
}
