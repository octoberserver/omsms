package crud

import (
	"omsms/db"
	"omsms/util"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmdRunning bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "\033[34m查看所有伺服器\033[0m",
	Long:  "\033[34m查看所有伺服器\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		var servers []db.Server
		db.DB.Find(&servers)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "名稱", "Java", "備份策略", "運行狀態"})
		table.SetAutoFormatHeaders(false)

		ctx, cli := util.InitDockerClient()

		for _, server := range servers {
			containerName := util.GetServerName(server.ID)
			isRunning := util.DoesContainerExist(containerName, cli, ctx) && util.IsContainerRunning(containerName, cli, ctx)
			if listCmdRunning && !isRunning {
				continue
			}

			status := "未啟動"
			if isRunning {
				status = "運行中"
			}

			table.Append([]string{
				strconv.FormatInt(int64(server.ID), 10),
				server.Name,
				strconv.FormatInt(int64(server.Java), 10),
				server.Backup.String(),
				status,
			})
		}

		table.Render()
	},
}

func RegisterListCmd(parent *cobra.Command) {
	listCmd.Flags().BoolVarP(&listCmdRunning, "running", "r", false, "只顯示運行中")
	listCmd.MarkFlagRequired("id")
	parent.AddCommand(listCmd)
}
