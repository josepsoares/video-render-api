package render

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/josepsoares/video-render-api/pkg/ffmpeg"
	"github.com/josepsoares/video-render-api/pkg/utils"
)

func TestThumbnail(t *testing.T) {
	defer utils.TimeTrack(time.Now(), "Test - ImageToVideo")

	fmt.Println("starting test of turning image into video")

	err := godotenv.Load("../../.env")
	utils.CheckError("couldn't load .env file", err)

	outputPath := fmt.Sprintf(
		"%s/tmp/%s-%s-gen-thumbnail.png",
		os.Getenv("LOCAL_BASE_PATH"), strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String(),
	)

	ffmpegCmd := fmt.Sprintf(
		"ffmpeg -i %s/assets/videos/ua_biblioteca.mp4 -vf \"thumbnail\" -frames:v 1 %s",
		os.Getenv("LOCAL_BASE_PATH"), outputPath,
	)
	ffmpeg.Run(ffmpegCmd)

	fmt.Println("ended generateThumbnailTest concluded")
}
