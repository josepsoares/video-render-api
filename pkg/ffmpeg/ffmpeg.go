package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/josepsoares/video-render-api/pkg/utils"
)

func Run(command string) {
	s := command

	args := strings.Split(s, " ")

	cmd := exec.Command(args[0], args[1:]...)

	b, err := cmd.CombinedOutput()

	// an error occured
	if err != nil {
		utils.CheckError("running ffmpeg failed", err)
	}

	fmt.Printf("success => %s\n", b)
}

func RunProbe(command string) string {
	s := command

	args := strings.Split(s, " ")

	cmd := exec.Command(args[0], args[1:]...)

	b, err := cmd.CombinedOutput()

	// an error occured
	if err != nil {
		utils.CheckError("running ffprobe failed", err)
	}

	fmt.Printf("%s\n", b)

	return string(b)
}
