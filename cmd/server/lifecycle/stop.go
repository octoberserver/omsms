package lifecycle

import (
	"errors"
	"fmt"
	"omsms/db"
	"omsms/util"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// var stopCmdId uint32

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "\033[33m關閉伺服器\033[0m",
	Long:  "\033[33m關閉伺服器\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("\033[31m使用方式: omsms server stop [id]\033[0m")
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

		ctx, cli := util.InitDockerClient()
		defer util.CloseDockerClient(cli)

		srvName := util.GetServerName(server.ID)
		if !util.DoesContainerExist(srvName, cli, ctx) {
			fmt.Println("\033[31m伺服器並沒有在運行中\033[0m")
			os.Exit(1)
		}

		if util.IsContainerRunning(srvName, cli, ctx) {
			timeout := int(30 * time.Second)

			fmt.Println("\033[34m正在關閉舊容器\033[0m")
			err := cli.ContainerStop(ctx, util.GetServerName(server.ID), container.StopOptions{Timeout: &timeout})
			if err != nil {
				fmt.Printf("\033[31m無法關閉容器: %v\033[0m\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("\033[34m正在移除舊容器\033[0m")
		err = cli.ContainerRemove(ctx, util.GetServerName(server.ID), container.RemoveOptions{})
		if err != nil {
			fmt.Printf("\033[31m無法移除容器: %v\033[0m\n", err)
			os.Exit(1)
		}

		util.DeleteProxyHost(cli, ctx, &server)

		fmt.Println("\033[32m伺服器成功關閉\033[0m")
	},
}

func RegisterStopCmd(parent *cobra.Command) {
	// stopCmd.Flags().Uint32VarP(&stopCmdId, "id", "i", 0, "伺服器ID")
	// stopCmd.MarkFlagRequired("id")
	parent.AddCommand(stopCmd)
}
