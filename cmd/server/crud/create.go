package crud

import (
	"fmt"
	"omsms/db"
	"omsms/util"
	"omsms/util/enums"
	"os"

	"github.com/spf13/cobra"
)

var createCmdName string
var createCmdJava uint32
var createCmdBackup enums.BackupStrat = enums.BACKUP_NONE

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "創建伺服器",
	Long:  `創建伺服器`,
	Run: func(cmd *cobra.Command, args []string) {
		server := &db.Server{
			Name:   createCmdName,
			Java:   uint(createCmdJava),
			Backup: createCmdBackup,
		}

		db.DB.Create(server)

		path := util.GetServerFolderPath(server.ID)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("\033[32m無法創建資料夾: ", err)
			return
		}

		fmt.Println("\033[1;32m----創建伺服器----\033[0m")
		fmt.Println("\033[1;32m伺服器ID:\033[0m", server.ID)
		fmt.Println("\033[32m創建時間:\033[0m", server.CreatedAt)
		fmt.Println("\033[33m伺服器名稱:\033[0m", server.Name)
		fmt.Println("\033[33mJava版本:\033[0m", server.Java)
		fmt.Println("\033[33m備份策略:\033[0m", server.Backup)
		fmt.Println("\033[33m檔案路徑:\033[0m", path)
		fmt.Println("\033[1;32m------------------\033[0m")
	},
}

func RegisterCreateCmd(parent *cobra.Command) {
	createCmd.Flags().StringVarP(&createCmdName, "name", "n", "", "伺服器名稱")
	createCmd.MarkFlagRequired("name")
	createCmd.Flags().Uint32VarP(&createCmdJava, "java", "j", 0, "Java版本")
	createCmd.MarkFlagRequired("java")
	createCmd.Flags().VarP(&createCmdBackup, "backup", "b", `備份策略："FULL_SERVER", "WORLD", "CUSTOM" 或 "NONE"`)
	parent.AddCommand(createCmd)
}
