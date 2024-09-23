package lifecycle

import (
	"errors"
	"fmt"
	"log"
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
	Short: "關閉伺服器",
	Long:  `關閉伺服器`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			os.Exit(0)
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("ID必須是數字")
			os.Exit(0)
		}

		var server db.Server
		if errors.Is(db.DB.First(&server, id).Error, gorm.ErrRecordNotFound) {
			fmt.Println("伺服器不存在: " + strconv.FormatUint(uint64(id), 10))
			os.Exit(0)
		}

		ctx, cli := util.InitDockerClient()
		defer util.CloseDockerClient(cli)

		if !util.DoesContainerExist(server.ID, cli, ctx) {
			fmt.Println("伺服器並沒有在運行中")
			os.Exit(0)
		}

		if util.IsContainerRunning(server.ID, cli, ctx) {
			timeout := int(30 * time.Second)

			fmt.Println("正在關閉舊容器")
			err := cli.ContainerStop(ctx, util.GetServerName(server.ID), container.StopOptions{Timeout: &timeout})
			if err != nil {
				log.Fatalf("無法關閉容器: %v", err)
			}
		}

		fmt.Println("正在移除舊容器")
		err = cli.ContainerRemove(ctx, util.GetServerName(server.ID), container.RemoveOptions{})
		if err != nil {
			fmt.Printf("無法移除容器: %v\n", err)
			os.Exit(0)
		}

		fmt.Println("伺服器成功關閉")
	},
}

func RegisterStopCmd(parent *cobra.Command) {
	// stopCmd.Flags().Uint32VarP(&stopCmdId, "id", "i", 0, "伺服器ID")
	// stopCmd.MarkFlagRequired("id")
	parent.AddCommand(stopCmd)
}
