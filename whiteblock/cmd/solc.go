package cmd

import (
	"fmt"
	"os"
	"os/exec"

	solc "github.com/ethereum/go-ethereum/common/compiler"
	"github.com/spf13/cobra"
)

/*
should figure out how to identify which blockchain for compile
- have them input the type of blockchain as args
- have program figure out what blockchain the smart contract corresponds to
*/

func compile(path, filename string) {
	out, err := solc.CompileSolidity("solc", path+"/"+filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

var solcCMD = &cobra.Command{
	Use:   "contract",
	Short: "Add and compile a smart contract.",
	Long: `
Contract allows the user to add and compile a smart contract.
`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock contract -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var addSolcCMD = &cobra.Command{
	Use:   "add <path> <filename>",
	Short: "Add a smart contract.",
	Long: `
Adds the specified smart contract into the /Downloads folder.
`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			out, err := exec.Command("bash", "-c", "./whiteblock contract add -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		cp := "cp " + args[0] + "/" + args[1] + " ~/Downloads/"

		out, err := exec.Command("bash", "-c", cp).Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
	},
}

var compileSolcCMD = &cobra.Command{
	Use:   "compile <path> <filename>",
	Short: "Smart contract compiler.",
	Long: `
Compiles the specified smart contract.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			out, err := exec.Command("bash", "-c", "./whiteblock contract compile -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}

		compile(args[0], args[1])
	},
}

func init() {
	solcCMD.AddCommand(addSolcCMD, compileSolcCMD)

	// RootCmd.AddCommand(solcCMD)
}
