package cmd

import (
	"fmt"
	"omsms/cmd/server"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "omsms",
	Short: "\033[35m十月模組伺服器管理系統\033[0m",
	Long:  "\033[35m十月模組伺服器管理系統\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		println("\033[35m" + "歡迎來到十月模組伺服器管理系統")
		cmd.Help()
	},
}

func Execute() {
	server.RegisterServerCmd(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
