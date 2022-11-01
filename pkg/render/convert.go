package render

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/josepsoares/video-render-api/pkg/ffmpeg"
	"github.com/josepsoares/video-render-api/pkg/logger"
	"github.com/josepsoares/video-render-api/pkg/utils"
	"go.uber.org/zap"
)

// create a static video from an image
func ConvertImgToVideo(img string, res int, duration string, bgMusic *string) string {
	logger.RenderLogger.Info("starting to convert image to video ⏳")

	outputPath := fmt.Sprintf("/tmp/%s-%s.mp4", strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String())

	ffmpegCmd := fmt.Sprintf(
		"ffmpeg -i %s %s -vf scale=1:%s %s -t %s",
		img, utils.GetResolutionFilters(res), strconv.Itoa(res), outputPath, duration,
	)
	ffmpeg.Run(ffmpegCmd)

	logger.RenderLogger.Info("generated video successfully ✔️")

	// check if the bgMusic param is not nil, if not add the param value in a ffmpeg command to
	// add the audio to the video
	if bgMusic != nil {
		outputPathWithAudio := fmt.Sprintf("/tmp/%s-%s.mp4", strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String())
		bgMusicAsset := *bgMusic

		// todo: might be able to do this (convert img to video and add audio) in only one ffmpeg command
		ffmpegCmd := fmt.Sprintf(
			"ffmpeg -i %s -i %s -c copy -map 0:v:0 -map 1:a:0 %s",
			outputPath, bgMusicAsset, outputPathWithAudio,
		)
		ffmpeg.Run(ffmpegCmd)

		logger.RenderLogger.Info("background music detected, added audio to the generated video ✔️")

		// remove file of the video without audio
		err := os.Remove(outputPath)
		utils.CheckError("error removing file", err)

		outputPath = outputPathWithAudio
	}

	logger.RenderLogger.Info("converted video successfully ✔️", zap.String("result", outputPath))

	return outputPath
}
