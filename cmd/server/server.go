package server

import (
	"fmt"
	"omsms/cmd/server/crud"
	"omsms/cmd/server/lifecycle"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "管理伺服器",
	Long:  `管理伺服器`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("管理伺服器")
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
}
