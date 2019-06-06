package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	util "github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"os"
	"strconv"
)

func GetNodes() ([]Node, error) {
	testnetId, err := getPreviousBuildId()
	if err != nil {
		return nil, err
	}
	res, err := jsonRpcCall("nodes", []string{testnetId})
	if err != nil {
		return nil, err
	}
	tmp := res.([]interface{})
	nodes := []map[string]interface{}{}
	for _, t := range tmp {
		nodes = append(nodes, t.(map[string]interface{}))
	}

	out := []Node{}
	for _, node := range nodes {
		out = append(out, Node{
			LocalID:   int(node["localId"].(float64)),
			Server:    int(node["server"].(float64)),
			TestNetID: node["testnetId"].(string),
			ID:        node["id"].(string),
			IP:        node["ip"].(string),
			Label:     node["label"].(string),
		})
	}
	return out, nil
}

func readContractsFile() ([]byte, error) {
	cwd := os.Getenv("HOME")
	return ioutil.ReadFile(cwd + "/smart-contracts/whiteblock/contracts.json")
}

var getCmd = &cobra.Command{
	Use:   "get <command>",
	Short: "Get server and network information.",
	Long:  "\nGet will output server and network information and statistics.\n",
	Run:   util.PartialCommand,
}

var getServerCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"servers"},
	Short:   "Get server information.",
	Long:    "\nServer will output server information.\n",
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("get_servers", []string{})
	},
}

var getTestnetIDCmd = &cobra.Command{
	Use:     "testnetid",
	Aliases: []string{"id"},
	Short:   "Get the last stored testnet id",
	Long:    "\nGet the last stored testnet id.\n",
	Run: func(cmd *cobra.Command, args []string) {
		testnetID, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		fmt.Println(testnetID)
	},
}

var getSupportedCmd = &cobra.Command{
	Use:     "supported",
	Aliases: []string{"blockchains"},
	Short:   "Get the currently supported blockchains",
	Long:    "Fetches the blockchains which whiteblock is currently able build by default",
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("get_supported_blockchains", []string{})
	},
}

var getNodesCmd = &cobra.Command{
	Use:     "nodes",
	Aliases: []string{"node"},
	Short:   "Nodes will show all nodes in the network.",
	Long:    "\nNodes will output all of the nodes in the current network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		testnetID, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if all {
			jsonRpcCallAndPrint("status_nodes", []string{testnetID})
			return
		}
		res, err := jsonRpcCall("status_nodes", []string{testnetID})
		if err != nil {
			util.PrintErrorFatal(err)
		}

		rawNodes := res.([]interface{})
		out := []interface{}{}
		for _, rawNode := range rawNodes {
			if rawNode.(map[string]interface{})["up"].(bool) {
				out = append(out, rawNode)
			}
		}
		fmt.Println(prettypi(out))
	},
}

var getRunningCmd = &cobra.Command{
	Use:   "running",
	Short: "Running will check if a test is running.",
	Long: `
Running will check whether or not there is a test running and get the name of the currently running test.

Response: true or false, on whether or not a test is running; The name of the test or nothing if there is not a test running.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		jsonRpcCallAndPrint("state::is_running", []string{})
		jsonRpcCallAndPrint("state::what_is_running", []string{})
	},
}

var getDefaultsCmd = &cobra.Command{
	Use:     "default <blockchain>",
	Aliases: []string{"defaults"},
	Short:   "Default gets the blockchain params.",
	Long: `
Get the blockchain specific parameters for a deployed blockchain.

Format: The blockchain to get the build params of

Response: The params as a list of key value params, of name and type respectively
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		jsonRpcCallAndPrint("get_defaults", args)
	},
}

var getConfigsCmd = &cobra.Command{
	Use:     "configs <blockchain> [file]",
	Aliases: []string{"config"},
	Short:   "Get the resources for a blockchain",
	Long: `
Get the resources for a blockchain. With one argument, lists what is available. With two
	arguments, get the contents of the file

Params: The blockchain to get the resources of, the resource/file name 

Response: The resoures as a list of key value params, of name and type respectively
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		jsonRpcCallAndPrint("get_resources", args)
	},
}

var getStatsCmd = &cobra.Command{
	Use:   "stats <command>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats will allow the user to get statistics regarding the network.

Response: JSON representation of network statistics
	`,
	Run: util.PartialCommand,
}

var statsByTimeCmd = &cobra.Command{
	Use:   "time <start time> <end time>",
	Short: "Get stastics by time",
	Long: `
Stats time will allow the user to get statistics by specifying a start time and stop time (unix time stamp).

Params: start unix timestamp, end unix timestamp

Response: JSON representation of network statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		jsonRpcCallAndPrint("stats", map[string]int64{
			"startTime":  util.CheckAndConvertInt64(args[0], "start unix timestamp"),
			"endTime":    util.CheckAndConvertInt64(args[1], "end unix timestamp"),
			"startBlock": 0,
			"endBlock":   0,
		})
	},
}

var statsByBlockCmd = &cobra.Command{
	Use:   "block <start block> <end block>",
	Short: "Get stastics of a blockchain",
	Long: `
Stats block will allow the user to get statistics regarding the network.

Params: start block number end block number

Response: JSON representation of statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		jsonRpcCallAndPrint("stats", map[string]int64{
			"startTime":  0,
			"endTime":    0,
			"startBlock": util.CheckAndConvertInt64(args[0], "start block number"),
			"endBlock":   util.CheckAndConvertInt64(args[1], "end block number"),
		})
	},
}

var statsPastBlocksCmd = &cobra.Command{
	Use:   "past <blocks> ",
	Short: "Get stastics of a blockchain from the past x blocks",
	Long: `
Stats block will allow the user to get statistics regarding the network.

Params: Number of blocks 

Response: JSON representation of statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		jsonRpcCallAndPrint("stats", map[string]int64{
			"startTime":  0,
			"endTime":    0,
			"startBlock": util.CheckAndConvertInt64(args[0], "blocks") * -1, //Negative number signals past
			"endBlock":   0,
		})
	},
}

var statsAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Get all stastics of a blockchain",
	Long: `
Stats all will allow the user to get all the statistics regarding the network.

Response: JSON representation of network statistics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("all_stats", []string{})
	},
}

/*
Work underway on generalized use commands to consolidate all the different
commands separated by blockchains.
*/

func getBlockCobra(cmd *cobra.Command, args []string) {
	util.CheckArguments(cmd, args, 1, 1)
	blockNum := 0
	var err error
	if len(args) > 0 {
		blockNum, err = strconv.Atoi(args[0])
		if err != nil {
			util.PrintStringError("Invalid block number formatting.")
			return
		}
	}
	if blockNum < 1 && len(args) > 0 {
		util.PrintStringError("Unable to get block information from block 0. Please provide a block number greater than 0.")
		return
	} else {
		res, err := jsonRpcCall("get_block_number", []string{})
		if err != nil {
			util.PrintErrorFatal(err)
		}
		blocknum := int(res.(float64))
		if blocknum < 1 {
			util.PrintStringError("Unable to get block information because no blocks have been created. Please use the command 'whiteblock miner start' to start generating blocks.")
			return
		}
	}
	jsonRpcCallAndPrint("get_block", args)
}

var getBlockCmd = &cobra.Command{
	// Hidden: true,
	Use:   "block <command>",
	Short: "Get information regarding blocks",
	Run:   getBlockCobra,
}

var getBlockNumCmd = &cobra.Command{
	// Hidden: true,
	Use:   "number",
	Short: "Get the block number",
	Long: `
Gets the most recent block number that had been added to the blockchain.

Response: block number
	`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("get_block_number", []string{})
	},
}

var getBlockInfoCmd = &cobra.Command{
	// Hidden: true,
	Use:   "info <block number>",
	Short: "Get the information of a block",
	Long: `
Gets the information inside a block including transactions and other information relevant to the currently connected blockchain.

Params: Block number

Response: JSON representation of the block
	`,
	Run: getBlockCobra,
}

var getTxCmd = &cobra.Command{
	// Hidden: true,
	Use:   "tx <command>",
	Short: "Get information regarding transactions",
	Run:   util.PartialCommand,
}

var getTxInfoCmd = &cobra.Command{
	// Hidden: true,
	Use:   "info <tx hash>",
	Short: "Get transaction information",
	Long: `
Get a transaction by its hash. The user can find the transaction hash by viewing block information. To view block information, the command 'get block info <block number>' can be used.

Params: The transaction hash

Response: JSON representation of the transaction.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		jsonRpcCallAndPrint("get_transaction", args)
	},
}

// eth::get_transaction_receipt does not work.
/*
var getTxReceiptCmd = &cobra.Command{
	// Hidden: true,
	Use:   "receipt <tx hash>",
	Short: "Get the transaction receipt",
	Long: `
Get the transaction receipt by the tx hash. The user can find the transaction hash by viewing block information. To view block information, the command 'get block info <block number>' can be used.

Format: <hash>
Params: The transaction hash

Response: JSON representation of the transaction receipt.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		command := ""
		switch blockchain {
		case "ethereum":
			command = "eth::get_transaction_receipt"
		case "eos":
			fmt.Println("This function is not supported for the syscoin client.")
		case "syscoin":
			fmt.Println("This function is not supported for the syscoin client.")
		default:
			fmt.Println("No blockchain found. Please use the build function to create one")
			return
		}
		jsonRpcCallAndPrint(command, args)
	},
}
*/

var getAccountCmd = &cobra.Command{
	// Hidden: true,
	Use:   "account <command>",
	Short: "Get account information",
	Run:   util.PartialCommand,
}

var getAccountInfoCmd = &cobra.Command{
	// Hidden: true,
	Use:   "info",
	Short: "Get account information",
	Long: `
Gets the account information relevant to the currently connected blockchain.

Response: JSON representation of the accounts information.
`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonRpcCallAndPrint("accounts_status", []string{})
	},
}

var getContractsCmd = &cobra.Command{
	// Hidden: true,
	Use:   "contracts",
	Short: "Get contracts deployed to network.",
	Long: `
Gets the list of contracts that were deployed to the network. The information includes the address that deployed the contract, the contract name, and the contract's address.

Response: JSON representation of the contract information.
`,
	Run: func(cmd *cobra.Command, args []string) {

		contracts, err := readContractsFile()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(contracts) == 0 {
			util.PrintStringError("No smart contract has been deployed yet. Please use the command 'whiteblock geth solc deploy <smart contract> to deploy a smart contract.")
			os.Exit(1)
		} else {
			fmt.Println(prettyp(string(contracts)))
		}
	},
}

func init() {
	getNodesCmd.Flags().Bool("all", false, "output all of the nodes, even if they are no longer running")
	getCmd.AddCommand(getServerCmd, getNodesCmd, getStatsCmd, getDefaultsCmd, getSupportedCmd, getRunningCmd, getConfigsCmd, getTestnetIDCmd)

	getStatsCmd.AddCommand(statsByTimeCmd, statsByBlockCmd, statsPastBlocksCmd, statsAllCmd)

	// dev commands that are currently being implemented
	getCmd.AddCommand(getBlockCmd, getTxCmd, getAccountCmd, getContractsCmd)
	getBlockCmd.AddCommand(getBlockNumCmd, getBlockInfoCmd)
	getTxCmd.AddCommand(getTxInfoCmd)
	getAccountCmd.AddCommand(getAccountInfoCmd)

	RootCmd.AddCommand(getCmd)
}
