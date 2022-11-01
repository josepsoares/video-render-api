package render

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/josepsoares/video-render-api/pkg/ffmpeg"
	"github.com/josepsoares/video-render-api/pkg/logger"

	"go.uber.org/zap"
)

// transcode the video to a different codec
func Transcode(video string, codec string) string {
	logger.RenderLogger.Info("starting to transcode video ⏳", zap.String("codec", codec))

	outputPath := fmt.Sprintf("/tmp/%s-%s.mp4", strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String())

	ffmpegCmd := fmt.Sprintf(
		"ffmpeg -i %s -c:v %s %s",
		video, codec, outputPath,
	)
	ffmpeg.Run(ffmpegCmd)

	logger.RenderLogger.Info("transcoded video successfully ✔️", zap.String("result", outputPath))

	return outputPath
}
