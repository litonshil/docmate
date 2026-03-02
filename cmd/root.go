package cmd

import (
	"docmate/client/conn"
	"docmate/client/logger"
	"docmate/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "doc-mate service entrypoint",
}

// RegisterSubCommand allows subcommands to be added to rootCmd.
func RegisterSubCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func Execute() error {
	RegisterServeCommand()
	RegisterMigrateCommand()

	if err := config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	appConfig := config.App()
	logger.Set(*appConfig)

	conn.ConnectDB()

	return rootCmd.Execute()
}
