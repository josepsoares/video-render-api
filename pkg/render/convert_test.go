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

func TestConvert(t *testing.T) {
	fmt.Println("starting test of turning image into video")

	err := godotenv.Load("../../.env")
	utils.FailOnError("couldn't load .env file", err)

	res := 1080
	duration := 10

	outputPath := fmt.Sprintf(
		"%s/tmp/%s-%s-%dp-img-to-video.mp4",
		os.Getenv("LOCAL_BASE_PATH"), strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String(), res,
	)
	// turn preset image to video
	ffmpegCmd := fmt.Sprintf(
		"ffmpeg -i %s/assets/imgs/cows-oleksandr-kurchev.jpg -codec:v libx264 %s -vf scale=-2:%d %s -t %d",
		os.Getenv("LOCAL_BASE_PATH"), utils.GetResolutionFilters(res), res, outputPath, duration,
	)
	ffmpeg.Run(ffmpegCmd)

	fmt.Println("ended ImageToVideo test - output path where is the result:" + outputPath)
}
