package render

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/josepsoares/video-render-api/pkg/ffmpeg"
	"github.com/josepsoares/video-render-api/pkg/logger"
	"github.com/josepsoares/video-render-api/pkg/utils"
	"go.uber.org/zap"
)

// downscale a video to other resolutions (from 1080 to 360p)
func Downscale(video string, res int) string {
	logger.RenderLogger.Info("starting to downscale video ⏳", zap.Int("resolution", res))

	outputPath := fmt.Sprintf("/tmp/%s-%s.mp4", strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String())

	ffmpegCmd := fmt.Sprintf(
		"ffmpeg -i %s %s -vf scale=1:%s %s",
		video, utils.GetResolutionFilters(res), strconv.Itoa(res), outputPath,
	)
	ffmpeg.Run(ffmpegCmd)

	logger.RenderLogger.Info("downscaled video successfully ✔️", zap.String("result", outputPath))

	return outputPath
}
