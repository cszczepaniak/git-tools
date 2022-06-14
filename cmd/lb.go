package cmd

import (
	"fmt"
	"os"

	"github.com/cszczepaniak/git-tools/lib/git"
	"github.com/cszczepaniak/git-tools/lib/git/client"
	"github.com/spf13/cobra"
)

// lbCmd represents the lb command
var lbCmd = &cobra.Command{
	Use:   `lb`,
	Short: `Print a selectable list of latest checked-out git branches`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		c := client.NewClient(wd)
		latest, err := git.LatestBranches(c, 1000)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		choices := make([]string, 0, len(latest))
		for branch, timestamp := range latest {
			choices = append(choices, fmt.Sprintf(`%s (%s)`, branch, timestamp))
		}

		ct := *count
		if len(choices) < ct {
			ct = len(choices)
		}

		fmt.Println(choices[:ct])
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
