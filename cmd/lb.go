package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lbCmd represents the lb command
var lbCmd = &cobra.Command{
	Use:   `lb`,
	Short: `Print a selectable list of latest checked-out git branches`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`lb called`, *count)
	},
}

var count *int

func init() {
	rootCmd.AddCommand(lbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	count = lbCmd.Flags().IntP(`count`, `n`, 25, `Number of branches displayed per page`)
}
