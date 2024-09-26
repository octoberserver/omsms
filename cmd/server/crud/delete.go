package crud

import (
	"errors"
	"fmt"
	"omsms/db"
	"omsms/util"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var deleteCmdId uint32

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "\033[31m刪除伺服器\033[0m",
	Long:  "\033[31m刪除伺服器\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		var server db.Server
		if errors.Is(db.DB.First(&server, deleteCmdId).Error, gorm.ErrRecordNotFound) {
			fmt.Println("\033[31m伺服器不存在:", deleteCmdId, "\033[0m")
			os.Exit(1)
		}

		ctx, cli := util.InitDockerClient()
		util.DeleteProxyHost(cli, ctx, &server)

		path := util.GetServerFolderPath(server.ID)
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println("\033[31m無法刪除資料夾: ", err)
		}

		db.DB.Delete(&server)

		fmt.Println("\033[1;31m----刪除伺服器----")
		fmt.Println("\033[1;31m伺服器ID:\033[0m", server.ID)
		fmt.Println("\033[31m創建時間:\033[0m", server.CreatedAt)
		fmt.Println("\033[31m伺服器名稱:\033[0m", server.Name)
		fmt.Println("\033[31mJava版本:\033[0m", server.Java)
		fmt.Println("\033[31m備份策略:\033[0m", server.Backup)
		fmt.Println("\033[31m檔案路徑:\033[0m", path)
		fmt.Println("\033[1;31m------------------\033[0m")
	},
}

func RegisterDeleteCmd(parent *cobra.Command) {
	deleteCmd.Flags().Uint32VarP(&deleteCmdId, "id", "i", 0, "伺服器ID")
	deleteCmd.MarkFlagRequired("id")
	parent.AddCommand(deleteCmd)
}
