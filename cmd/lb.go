package cmd

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/cszczepaniak/git-tools/lib/git"
	"github.com/cszczepaniak/git-tools/lib/git/client"
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
		latest, err := git.LatestBranches(c, *reflogLimit)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		displayOpts := make([]string, 0, len(latest))
		branches := make([]string, 0, len(latest))
		for _, l := range latest {
			displayOpts = append(displayOpts, l.String())
			branches = append(branches, l.Name())
		}

		var selected int
		err = survey.AskOne(&survey.Select{
			Message:  `Select a branch to switch to:`,
			PageSize: *count,
			Options:  displayOpts,
		}, &selected)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		err = c.Checkout(branches[selected])
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
	},
}

var (
	count       *int
	reflogLimit *int
)

func init() {
	rootCmd.AddCommand(lbCmd)

	count = lbCmd.Flags().IntP(`count`, `n`, 25, `Number of branches displayed per page`)
	reflogLimit = lbCmd.Flags().IntP(`limit`, `limit`, 1000, `Number of reflog lines to parse`)
}
