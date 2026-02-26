package cmd

//import (
//	"docmate/client/conn"
//	"docmate/client/logger"
//	"docmate/config"
//	"fmt"
//	"os"
//
//	"github.com/spf13/cobra"
//)
//
//var (
//	RootCmd = &cobra.Command{
//		Use: "docmate",
//	}
//)
//
//// Execute executes the root command.
//func Execute() {
//	RootCmd.AddCommand(serveCmd)
//	if err := config.Load(); err != nil {
//		os.Exit(1)
//	}
//
//	appConfig := config.App()
//	logger.Set(*appConfig)
//
//	conn.ConnectDB()
//
//	if err := RootCmd.Execute(); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//}

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
