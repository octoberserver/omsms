package util

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/client"
)

func InitDockerClient() (context.Context, *client.Client) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("無法初始化Docker引擎連線")
		os.Exit(0)
	}
	return ctx, cli
}

func CloseDockerClient(cli *client.Client) {
	cli.Close()
}

func DoesContainerExist(serverId uint, cli *client.Client, ctx context.Context) bool {
	_, err := cli.ContainerInspect(ctx, GetServerName(serverId))
	if err != nil {
		// Handle the error if the container doesn't exist
		if strings.Contains(err.Error(), "No such container") {
			return false
		} else {
			fmt.Printf("檢測容器失敗: %v\n", err)
			fmt.Println("失敗函式: DoesContainerExist")
			os.Exit(0)
		}
	}
	return true
}

func IsContainerRunning(serverId uint, cli *client.Client, ctx context.Context) bool {
	container, err := cli.ContainerInspect(ctx, GetServerName(serverId))
	if err != nil {
		fmt.Printf("檢測容器失敗: %v\n", err)
		fmt.Println("失敗函式: IsContainerRunning")
		os.Exit(0)
	}

	return container.State.Status == "running"
}
