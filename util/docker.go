package util

import (
	"context"

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

func CloseDockerClient(cli client.Client) {
	cli.Close()
}
