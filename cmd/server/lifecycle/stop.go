package lifecycle

import (
	"errors"
	"fmt"
	"log"
	"omsms/db"
	"omsms/util"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var stopCmdId uint32

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "關閉伺服器",
	Long:  `關閉伺服器`,
	Run: func(cmd *cobra.Command, args []string) {
		var server db.Server
		if errors.Is(db.DB.First(&server, stopCmdId).Error, gorm.ErrRecordNotFound) {
			fmt.Println("伺服器不存在: " + strconv.FormatUint(uint64(stopCmdId), 10))
			return
		}

		ctx, cli := util.InitDockerClient()

		if !doesContainerExist(server.ID, cli, ctx) {
			fmt.Println("伺服器並沒有在運行中")
			return
		}

		if isContainerRunning(server.ID, cli, ctx) {
			timeout := int(30 * time.Second) // Set the desired timeout for graceful shutdown

			fmt.Println("正在關閉舊容器")
			err := cli.ContainerStop(ctx, util.GetServerName(server.ID), container.StopOptions{Timeout: &timeout})
			if err != nil {
				log.Fatalf("無法關閉容器: %v", err)
			}
		}

		fmt.Println("正在移除舊容器")
		err := cli.ContainerRemove(ctx, util.GetServerName(server.ID), container.RemoveOptions{})
		if err != nil {
			fmt.Printf("無法移除容器: %v\n", err)
			return
		}

		fmt.Println("伺服器成功關閉")
	},
}

func RegisterStopCmd(parent *cobra.Command) {
	stopCmd.Flags().Uint32VarP(&stopCmdId, "id", "i", 0, "伺服器ID")
	stopCmd.MarkFlagRequired("id")
	parent.AddCommand(stopCmd)
}
