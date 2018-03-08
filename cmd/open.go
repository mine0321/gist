package cmd

import (
	"github.com/mine0321/gist/cli"
	"github.com/mine0321/gist/cli/config"
	"github.com/mine0321/gist/cli/screen"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open user's gist",
	Long:  "Open user's gist",
	RunE:  open,
}

func open(cmd *cobra.Command, args []string) (err error) {
	if config.Conf.Flag.OpenBaseURL {
		return cli.Open(config.Conf.Gist.BaseURL)
	}

	s, err := screen.Open()
	if err != nil {
		return err
	}

	rows, err := s.Select()
	if err != nil {
		return err
	}

	return cli.Open(rows[0].URL)
}

func init() {
	RootCmd.AddCommand(openCmd)
	openCmd.Flags().BoolVarP(&config.Conf.Flag.OpenBaseURL, "no-select", "", false, "Open only gist base URL without selecting")
	openCmd.Flags().BoolVarP(&config.Conf.Flag.StarredItems, "starred", "s", false, "Open your starred gist")
}
