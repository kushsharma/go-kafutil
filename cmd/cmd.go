package cmd

import (
	"github.com/kushsharma/go-kafutil/config"
	"github.com/spf13/cobra"
)

// InitCommands initializes application cli interface
func InitCommands(appname, version string, conf config.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     appname,
		Version: version,
	}
	rootCmd.AddCommand(initWriter(conf))
	rootCmd.AddCommand(initReader(conf))
	return rootCmd
}
