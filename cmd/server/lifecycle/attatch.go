package lifecycle

import (
	"errors"
	"fmt"
	"omsms/db"
	"omsms/util"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// var attachCmdId uint32
var attachCmdDirect bool

var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "\033[34m打開伺服器終端\033[0m",
	Long:  "\033[34m打開伺服器終端\033[0m",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("\033[31m使用方式: omsms server attach [id]\033[0m")
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
		srvName := util.GetServerName(server.ID)

		if !util.DoesContainerExist(srvName, cli, ctx) {
			fmt.Println("\033[31m伺服器並沒有在運行中\033[0m")
			os.Exit(1)
		}

		if !util.IsContainerRunning(srvName, cli, ctx) {
			fmt.Println("\033[31m伺服器並沒有在運行中\033[0m")
			os.Exit(1)
		}

		if attachCmdDirect {
			dockerCmd := exec.Command("docker", "attach", util.GetServerName(server.ID))
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr
			dockerCmd.Run()
			os.Exit(0)
		}

		sessionName := fmt.Sprintf("omsms_%s", util.GetServerName(server.ID))
		tmuxCmd := exec.Command("bash", "-c", fmt.Sprintf("tmux attach -t %s", sessionName))
		tmuxCmd.Stdin = os.Stdin
		tmuxCmd.Stdout = os.Stdout
		tmuxCmd.Stderr = os.Stderr
		err = tmuxCmd.Run()
		if err != nil {
			fmt.Println("\033[31m無法連接到Tmux視窗: ", err, "\033[0m")
			os.Exit(1)
		}
	},
}

func RegisterAttachCmd(parent *cobra.Command) {
	// attachCmd.Flags().Uint32VarP(&attachCmdId, "id", "i", 0, "伺服器ID")
	// attachCmd.MarkFlagRequired("id")
	attachCmd.Flags().BoolVarP(&attachCmdDirect, "direct", "d", false, "直接使用Docker Attach(不使用Tmux)")
	parent.AddCommand(attachCmd)
}
