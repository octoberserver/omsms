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

var attachCmdId uint32
var attachCmdDirect bool

var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "打開伺服器終端",
	Long:  `打開伺服器終端`,
	Run: func(cmd *cobra.Command, args []string) {
		var server db.Server
		if errors.Is(db.DB.First(&server, attachCmdId).Error, gorm.ErrRecordNotFound) {
			fmt.Println("伺服器不存在: " + strconv.FormatUint(uint64(attachCmdId), 10))
			return
		}

		ctx, cli := util.InitDockerClient()
		defer util.CloseDockerClient(cli)

		if !util.DoesContainerExist(server.ID, cli, ctx) {
			fmt.Println("伺服器並沒有在運行中")
			return
		}

		if !util.IsContainerRunning(server.ID, cli, ctx) {
			fmt.Println("伺服器並沒有在運行中")
			return
		}

		if attachCmdDirect {
			dockerCmd := exec.Command("docker", "attach", util.GetServerName(server.ID))
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr
			dockerCmd.Run()
			return
		}

		sessionName := fmt.Sprintf("omsms_%s", util.GetServerName(server.ID))
		tmuxCmd := exec.Command("bash", "-c", fmt.Sprintf("tmux attach -t %s", sessionName))
		tmuxCmd.Stdin = os.Stdin
		tmuxCmd.Stdout = os.Stdout
		tmuxCmd.Stderr = os.Stderr
		err := tmuxCmd.Run()
		if err != nil {
			fmt.Println("Error attaching to tmux session:", err)
			return
		}
	},
}

func RegisterAttachCmd(parent *cobra.Command) {
	attachCmd.Flags().Uint32VarP(&attachCmdId, "id", "i", 0, "伺服器ID")
	attachCmd.MarkFlagRequired("id")
	attachCmd.Flags().BoolVarP(&attachCmdDirect, "direct", "d", false, "直接使用Docker Attach(不使用Tmux)")
	parent.AddCommand(attachCmd)
}
