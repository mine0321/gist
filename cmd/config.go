package cmd

import (
	"path/filepath"

	"github.com/mine0321/gist/cli"
	"github.com/mine0321/gist/cli/config"
	"github.com/spf13/cobra"
)

var confCmd = &cobra.Command{
	Use:   "config",
	Short: "Config the setting file",
	Long:  "Config the setting file with your editor (default: vim)",
	RunE:  conf,
}

func conf(cmd *cobra.Command, args []string) error {
	editor := config.Conf.Core.Editor
	tomlfile := config.Conf.Core.TomlFile
	if tomlfile == "" {
		dir, _ := config.GetDefaultDir()
		tomlfile = filepath.Join(dir, "config.toml")
	}
	return cli.Run(editor, tomlfile)
}

func init() {
	RootCmd.AddCommand(confCmd)
}
