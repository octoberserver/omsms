package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
)

func InitDockerClient() (context.Context, *client.Client) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
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
			panic(fmt.Sprintf("檢測容器失敗: %v", err))
		}
	}
	return true
}

func IsContainerRunning(serverId uint, cli *client.Client, ctx context.Context) bool {
	container, err := cli.ContainerInspect(ctx, GetServerName(serverId))
	if err != nil {
		panic(fmt.Sprintf("檢測容器失敗: %v", err))
	}

	return container.State.Status == "running"
}
