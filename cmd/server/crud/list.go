package crud

import (
	"omsms/db"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "查看所有伺服器",
	Long:  `查看所有伺服器`,
	Run: func(cmd *cobra.Command, args []string) {
		var servers []db.Server
		db.DB.Find(&servers)

		// TODO: Add is container created, running?
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "名稱", "Java", "備份策略"})
		table.SetAutoFormatHeaders(false)

		for _, server := range servers {
			table.Append([]string{
				strconv.FormatInt(int64(server.ID), 10),
				server.Name,
				strconv.FormatInt(int64(server.Java), 10),
				server.Backup.String(),
			})
		}

		table.Render()
	},
}

func RegisterListCmd(parent *cobra.Command) {
	parent.AddCommand(listCmd)

	// TODO: list running servers flag
}
