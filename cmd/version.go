package cmd

import (
	"fmt"

	"github.com/jsdelivr/globalping-cli/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Globalping CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Globalping CLI v" + version.Version)
	},
}
