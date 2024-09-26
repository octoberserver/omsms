package util

import (
	"fmt"
	"os"
)

func getRootPath() string {
	wd, _ := os.Getwd()
	return wd + "/omsms_data"
}

func GetServerName(serverId uint) string {
	return fmt.Sprintf("server_%d", serverId)
}

func GetServerFolderPath(serverId uint) string {
	return fmt.Sprintf("%s/server_files/%s", getRootPath(), GetServerName(serverId))
}

func GetProxyContainerName() string {
	return "SalmonProxy"
}
