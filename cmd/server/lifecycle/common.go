package lifecycle

import (
	"context"
	"fmt"
	"omsms/util"
	"strings"

	"github.com/docker/docker/client"
)

func doesContainerExist(serverId uint, cli *client.Client, ctx context.Context) bool {
	_, err := cli.ContainerInspect(ctx, util.GetServerName(serverId))
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

func isContainerRunning(serverId uint, cli *client.Client, ctx context.Context) bool {
	container, err := cli.ContainerInspect(ctx, util.GetServerName(serverId))
	if err != nil {
		panic(fmt.Sprintf("檢測容器失敗: %v", err))
	}

	return container.State.Status == "running"
}
