package server

import (
	"omsms/cmd/server/crud"
	"omsms/cmd/server/lifecycle"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "\033[35m管理伺服器\033[0m",
	Long:  "\033[35m管理伺服器\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func RegisterServerCmd(parent *cobra.Command) {
	parent.AddCommand(serverCmd)
	crud.RegisterCreateCmd(serverCmd)
	crud.RegisterListCmd(serverCmd)
	crud.RegisterDeleteCmd(serverCmd)
	lifecycle.RegisterStartCmd(serverCmd)
	lifecycle.RegisterStopCmd(serverCmd)
	lifecycle.RegisterAttachCmd(serverCmd)
	crud.RegisterAddFilesCmd(serverCmd)
	crud.RegisterUpdateCmd(serverCmd)
}
