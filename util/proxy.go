package util

import (
	"context"
	"fmt"
	"log"
	"omsms/db"
	"os"

	"github.com/docker/docker/client"
	"gopkg.in/yaml.v2"
)

const PROXY_CONFIG_FILE = "/home/october/minecraft/salmon-proxy/data/config.yml"

func SetProxyHost(_ *client.Client, _ context.Context, server *db.Server) error {
	return setOrDeleteProxyHost(server, false)
}

func DeleteProxyHost(_ *client.Client, _ context.Context, server *db.Server) error {
	return setOrDeleteProxyHost(server, true)
}

func setOrDeleteProxyHost(server *db.Server, del bool) error {
	f, err := os.Open(PROXY_CONFIG_FILE)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	data := make(map[string]string)
	err = yaml.NewDecoder(f).Decode(&data)
	if err != nil {
		log.Fatal(err)
		return err
	}

	serverAddr := fmt.Sprintf("%s:25565", GetServerName(server.ID))

	newData := data
	for k, v := range data {
		if v == serverAddr {
			delete(newData, k)
		}
	}

	if !del {
		newData[server.ProxyHost] = serverAddr
	}

	return nil
}
