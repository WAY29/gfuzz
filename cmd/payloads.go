package cmd

import (
	"fmt"
	"gfuzz/payloads"
	"gfuzz/utils"

	"github.com/spf13/cobra"
)

var payloadsCmd = &cobra.Command{
	Use:   "payloads",
	Short: "Show payloads",
	Long:  `Show all available payloads for --payload`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintInfo("PlaceHolders:")
		fmt.Println(" -  FUZZ FUZnZ")
		utils.PrintInfo("Available payloads:")
		for k := range payloads.PayloadsArray {
			fmt.Printf(" -  %-7s %s\n", k, payloads.GetPayloadInfo(k))
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(payloadsCmd)
}
