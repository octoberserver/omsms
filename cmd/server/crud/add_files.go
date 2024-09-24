package crud

import (
	"errors"
	"fmt"
	"omsms/db"
	"omsms/util"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var addFilesCmd = &cobra.Command{
	Use:   "add-files",
	Short: "\033[32m新增檔案\033[0m",
	Long:  "\033[32m新增檔案\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("\033[31m使用方式: omsms server add-files [id] [路徑]\033[0m")
			os.Exit(1)
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("\033[31mID必須是數字\033[0m")
			os.Exit(1)
		}

		path := args[1]
		if path == "" {
			fmt.Println("\033[31m請提供路徑\033[0m")
			os.Exit(1)
		}

		var server db.Server
		if errors.Is(db.DB.First(&server, id).Error, gorm.ErrRecordNotFound) {
			fmt.Println("\033[31m伺服器不存在:", id, "\033[0m")
			os.Exit(1)
		}

		serverFolderPath := util.GetServerFolderPath(server.ID)

		if !util.IsFolderEmpty(serverFolderPath) {
			fmt.Println("\033[31m伺服器資料夾內已經有檔案了\033[0m")
			os.Exit(1)
		}

		util.Unzip(path, serverFolderPath)

		util.RemoveTopLevelFolderIfExists(serverFolderPath)
	},
}

func RegisterAddFilesCmd(parent *cobra.Command) {
	parent.AddCommand(addFilesCmd)
}
