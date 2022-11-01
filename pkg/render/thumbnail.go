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

// generate a .png thumbnail of a video
func GenerateThumbnail(video string, frameTime string) string {
	logger.RenderLogger.Info("starting to generate thumbnail from video ⏳", zap.String("frame", frameTime))

	outputPath := fmt.Sprintf("/tmp/%s-%s.mp4", strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String())

	ffmpegCmd := fmt.Sprintf(
		"ffmpeg -i %s -vf \"thumbnail\" -frames:v 1 %s",
		video, outputPath,
	)
	ffmpeg.Run(ffmpegCmd)

	logger.RenderLogger.Info("generated thumbnail from video successfully ✔️", zap.String("result", outputPath))

	return outputPath
}
