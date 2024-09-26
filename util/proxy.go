package util

import (
	"context"
	"errors"
	"fmt"
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
	if len(server.HostNames) < 1 {
		return errors.New("伺服器未設定反向代理域名")
	}

	f, err := os.Open(PROXY_CONFIG_FILE)
	if err != nil {
		fmt.Println("Error opening file for reading")
		fmt.Println(err)
		return err
	}

	data := make(map[string]string)
	err = yaml.NewDecoder(f).Decode(&data)
	if err != nil {
		fmt.Println("Error on unmarshal")
		fmt.Println(err)
	}

	serverAddr := fmt.Sprintf("%s:25565", GetServerName(server.ID))

	newData := data
	for k, v := range data {
		if v == serverAddr {
			delete(newData, k)
		}
	}

	if !del {
		for _, hostName := range server.HostNames {
			newData[hostName] = serverAddr
		}
	}

	f.Close()

	f, err = os.Create(PROXY_CONFIG_FILE)
	err = yaml.NewEncoder(f).Encode(newData)
	if err != nil {
		fmt.Println("Error opening file for writing")
		fmt.Println(err)
		return err
	}
	f.Close()

	return nil
}
