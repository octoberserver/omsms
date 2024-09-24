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
	"os/exec"
	"path"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// var startCmdId uint32

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "\033[34m啟動伺服器\033[0m",
	Long:  "\033[34m啟動伺服器\033[0m",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("\033[31m使用方式: omsms server start [id]\033[0m")
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

		if util.DoesContainerExist(server.ID, cli, ctx) {
			if util.IsContainerRunning(server.ID, cli, ctx) {
				fmt.Println("\033[31m伺服器正在運行中，請先關閉再啟動\033[0m")
				os.Exit(1)
			}

			fmt.Println("\033[34m正在移除舊容器\033[0m")
			err := cli.ContainerRemove(ctx, util.GetServerName(server.ID), container.RemoveOptions{})
			if err != nil {
				log.Printf("\033[31m無法移除容器: %v\033[0m", err)
				os.Exit(1)
			}
		}

		util.GiveExecutePermission(path.Join(util.GetServerFolderPath(server.ID), "start.sh"))

		runContainer(cli, ctx, &server)
		fmt.Println("\033[32m伺服器成功啟動\033[0m")
		createTmuxSession(server.ID)
		fmt.Println("\033[32m成功創建Tmux視窗\033[0m")
	},
}

func runContainer(cli *client.Client, ctx context.Context, server *db.Server) {
	server_name := util.GetServerName(server.ID)

	path := util.GetServerFolderPath(server.ID)
	image_name := fmt.Sprintf("docker.io/library/eclipse-temurin:%d", server.Java)

	reader, err := cli.ImagePull(ctx, image_name, image.PullOptions{})
	if err != nil {
		fmt.Println("\033[31m無法獲取鏡像: ", err, "\033[0m")
		os.Exit(1)
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
		fmt.Println("\033[31m無法創建容器: ", err, "\033[0m")
		os.Exit(1)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Println("\033[31m無法啟動容器: ", err, "\033[0m")
		os.Exit(1)
	}
}

func createTmuxSession(serverId uint) {
	sessionName := fmt.Sprintf("omsms_%s", util.GetServerName(serverId))
	fmt.Println(sessionName)

	dockerAttachCmd := fmt.Sprintf("\"docker attach %s\"", util.GetServerName(serverId))
	fmt.Println(dockerAttachCmd)
	tmuxCmd := exec.Command("bash", "-c", fmt.Sprintf("(tmux new-session -d -s %s -n 十月模組伺服器 %s)&", sessionName, dockerAttachCmd))
	tmuxCmd.Stdin = os.Stdin
	tmuxCmd.Stdout = os.Stdout
	tmuxCmd.Stderr = os.Stderr
	err := tmuxCmd.Run()
	if err != nil {
		fmt.Println("\033[31m無法創建Tmux視窗: ", err, "\033[0m")
		os.Exit(1)
	}
}
func RegisterStartCmd(parent *cobra.Command) {
	// startCmd.Flags().Uint32VarP(&startCmdId, "id", "i", 0, "伺服器ID")
	// startCmd.MarkFlagRequired("id")
	parent.AddCommand(startCmd)
}
