package main

import (
	"fmt"
	"github.com/eris-ltd/mint-client/Godeps/_workspace/src/github.com/spf13/cobra"
	"os"
	"os/user"
	"path"
)

var (
	usr, _ = user.Current()

	DefaultKeyStore = path.Join(usr.HomeDir, ".eris", "keys")
)

func main() {
	var erisToMintCmd = &cobra.Command{
		Use:   "mint",
		Short: "Convert an eris-keys key to a priv_validator.json",
		Long:  "mintkey mint <address>",
		Run:   cliConvertErisKeyToPrivValidator,
	}

	var mintToErisCmd = &cobra.Command{
		Use:   "eris",
		Short: "Convert a priv_validator.json to an eris-keys key",
		Long:  "mintkey ers <path/to/priv_validator.json>",
		Run:   cliConvertPrivValidatorToErisKey,
	}

	var rootCmd = &cobra.Command{Use: "mintkey"}
	rootCmd.AddCommand(mintToErisCmd, erisToMintCmd)
	rootCmd.Execute()
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func ifExit(err error) {
	if err != nil {
		exit(err)
	}
}