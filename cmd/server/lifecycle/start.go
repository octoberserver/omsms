package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"omsms/db"
	"omsms/util"
	"os"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var startCmdId uint32

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "啟動伺服器",
	Long:  `啟動伺服器`,
	Run: func(cmd *cobra.Command, args []string) {
		var server db.Server
		if errors.Is(db.DB.First(&server, startCmdId).Error, gorm.ErrRecordNotFound) {
			log.Fatalln("伺服器不存在: " + strconv.FormatUint(uint64(startCmdId), 10))
			return
		}

		ctx, cli := util.InitDockerClient()
		defer util.CloseDockerClient(cli)

		if doesContainerExist(server.ID, cli, ctx) {
			if isContainerRunning(server.ID, cli, ctx) {
				log.Fatalln("伺服器正在運行中，請先關閉再啟動")
				return
			}

			fmt.Println("正在移除舊容器")
			err := cli.ContainerRemove(ctx, util.GetServerName(server.ID), container.RemoveOptions{})
			if err != nil {
				log.Fatalf("無法移除容器: %v", err)
				return
			}
		}

		runContainer(cli, ctx, &server)
		fmt.Println("伺服器成功啟動")
	},
}

func runContainer(cli *client.Client, ctx context.Context, server *db.Server) {
	server_name := util.GetServerName(server.ID)

	// TODO: Check if docker container is already created for this server

	path := util.GetServerFolderPath(server.ID)
	image_name := fmt.Sprintf("docker.io/library/eclipse-temurin:%d", server.Java)

	reader, err := cli.ImagePull(ctx, image_name, image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image_name,
		Cmd:          []string{"/mc_server/start.sh"},
		Tty:          true,
		OpenStdin:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{{
			Type:   mount.TypeBind,
			Source: path,
			Target: "/mc_server",
		}},
	}, nil, nil, server_name)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
}

func RegisterStartCmd(parent *cobra.Command) {
	startCmd.Flags().Uint32VarP(&startCmdId, "id", "i", 0, "伺服器ID")
	startCmd.MarkFlagRequired("id")
	parent.AddCommand(startCmd)
}
