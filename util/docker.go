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
		fmt.Println("\033[31m無法初始化Docker引擎連線\033[0m")
		os.Exit(1)
	}
	return ctx, cli
}

func CloseDockerClient(cli *client.Client) {
	cli.Close()
}

func DoesContainerExist(name string, cli *client.Client, ctx context.Context) bool {
	_, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		// Handle the error if the container doesn't exist
		if strings.Contains(err.Error(), "No such container") {
			return false
		} else {
			fmt.Printf("\033[31m檢測容器失敗: %v\n", err)
			fmt.Println("失敗函式: DoesContainerExist\033[0m")
			os.Exit(1)
		}
	}
	return true
}

func IsContainerRunning(name string, cli *client.Client, ctx context.Context) bool {
	container, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		fmt.Printf("\033[31m檢測容器失敗: %v\n", err)
		fmt.Println("失敗函式: IsContainerRunning\033[0m")
		os.Exit(1)
	}

	return container.State.Status == "running"
}
