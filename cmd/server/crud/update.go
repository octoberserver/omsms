package crud

import (
	"errors"
	"fmt"
	"omsms/db"
	"omsms/util"
	"omsms/util/enums"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var updateCmdName string
var updateCmdJava uint32
var updateCmdBackup enums.BackupStrat = enums.BACKUP_NULL
var updateCmdProxy string

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "\033[34m更新伺服器\033[0m",
	Long:  "\033[34m更新伺服器\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("\033[31m使用方式: omsms server attach [id]\033[0m")
			os.Exit(1)
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("\033[31mID必須是數字\033[0m")
			os.Exit(1)
		}

		var server db.Server
		if errors.Is(db.DB.First(&server, id).Error, gorm.ErrRecordNotFound) {
			fmt.Println("\033[31m伺服器不存在:", id, "\033[0m")
			os.Exit(1)
		}

		if updateCmdName != "" {
			server.Name = updateCmdName
		}
		if updateCmdJava != 0 {
			server.Java = uint(updateCmdJava)
		}
		if updateCmdBackup != enums.BACKUP_NULL {
			server.Backup = updateCmdBackup
		}
		if updateCmdProxy != "" {
			server.ProxyHost = updateCmdProxy
		}

		ctx, cli := util.InitDockerClient()
		err = util.SetProxyHost(cli, ctx, &server)
		if err == nil {
			fmt.Println("\033[32m成功設定反向代理\033[0m")
		}

		db.DB.Save(server)

		path := util.GetServerFolderPath(server.ID)

		fmt.Println("\033[1;32m----修改伺服器----\033[0m")
		fmt.Println("\033[1;32m伺服器ID:\033[0m", server.ID)
		fmt.Println("\033[32m創建時間:\033[0m", server.CreatedAt)
		fmt.Println("\033[33m伺服器名稱:\033[0m", server.Name)
		fmt.Println("\033[33mJava版本:\033[0m", server.Java)
		fmt.Println("\033[33m備份策略:\033[0m", server.Backup)
		fmt.Println("\033[33m檔案路徑:\033[0m", path)
		fmt.Println("\033[1;32m------------------\033[0m")
	},
}

func RegisterUpdateCmd(parent *cobra.Command) {
	updateCmd.Flags().StringVarP(&updateCmdName, "name", "n", "", "伺服器名稱")
	updateCmd.Flags().Uint32VarP(&updateCmdJava, "java", "j", 0, "Java版本")
	updateCmd.Flags().VarP(&updateCmdBackup, "backup", "b", `備份策略："FULL_SERVER", "WORLD", "CUSTOM" 或 "NONE"`)
	updateCmd.Flags().StringVarP(&updateCmdProxy, "proxy", "p", "", "反向代理域名")
	updateCmd.MarkFlagRequired("proxy")
	parent.AddCommand(updateCmd)
}
