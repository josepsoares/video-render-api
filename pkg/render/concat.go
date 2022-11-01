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

// create a video from a concatenation of list of videos
func ConcatVideos(assets []string, intro *string, outro *string, separator *string, watermark *string, transition string) string {
	logger.RenderLogger.Info("starting to concat videos ⏳")

	// initiliaze vars
	items := []string{}
	finalItems := []string{}
	xfadeOffsets := []float64{}

	var (
		xfadeVideoFilter string
		xfadeAudioFilter string
		treatedAssets    string
	)

	tmpFolder := strconv.FormatInt(time.Now().Unix(), 10) + "-" + uuid.New().String()
	concatOutputPath := fmt.Sprintf(
		"/tmp/%s-%s.mp4",
		strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String(),
	)

	// create folder in tmp to save rendered files for now
	dname, err := os.MkdirTemp("", tmpFolder)
	defer os.RemoveAll(dname)
	// check if folder in tmp was created successfully
	utils.CheckError("error creating temporary folder in /tmp", err)

	items = append(items, assets...)

	// verify if there is a template to follow
	// if so check if the template has a intro, outro or separator to add them to the items array
	if intro != nil {
		introAsset := *intro
		items = append(items, "")
		copy(items[1:], items)
		items[0] = introAsset

		logger.RenderLogger.Info("intro detected, added intro ✔️")
	}

	if outro != nil {
		outroAsset := *outro
		items = append(items, (outroAsset))

		logger.RenderLogger.Info("outro detected, added outro to the assets ✔️")
	}

	if separator != nil {
		// add separator asset url between the assets starting from the second asset (including intro)
		// and finishing before the last asset (including outro)
		for index := range items {
			finalItems = append(finalItems, items[index])
			if (index + 1) != len(items) {
				separatorAsset := *separator
				finalItems = append(finalItems, separatorAsset)
			}
		}

		items = finalItems

		logger.RenderLogger.Info("separator detected, added separators between the existing assets ✔️")
	}

	// before concat all the items together its needed to loop through all the items
	// to create the string for the xfade filter, its required to probe all the assets
	// to get their duration time
	// the xfadeVideoFilter and xfadeAudioFilter will be a string like the one bellow
	// [0][1:v]xfade=transition=fade:duration=1:offset=3[vfade1];[vfade1][2:v]xfade=transition=fade:duration=1:offset=10[vfade2]; \
	// [0:a][1:a]acrossfade=d=1[afade1];[afade1][2:a]acrossfade=d=1[afade2]; \
	for index, item := range items {
		incrementIndex := index + 1
		treatedAssets += fmt.Sprintf("-i %s ", item)

		itemDurationString := ffmpeg.RunProbe(fmt.Sprintf(
			"ffprobe -i %s -v quiet -show_entries format=duration -hide_banner -of default=noprint_wrappers=1:nokey=1", item),
		)
		itemDurationFloat, _ := strconv.ParseFloat(itemDurationString, 64)

		if index != 0 {
			xfadeVideoFilter += fmt.Sprintf(
				"[vfade%d][%d:v]xfade=transition=%s:duration=1:offset=%f[vfade%d];",
				index, incrementIndex, transition, (itemDurationFloat+xfadeOffsets[index-1])-1.0, incrementIndex,
			)
			xfadeAudioFilter += fmt.Sprintf(
				"[afade%d][%d:a]acrossfade=d=1[afade%d];",
				index, incrementIndex, incrementIndex,
			)
		} else {
			xfadeVideoFilter += fmt.Sprintf(
				"[%d][%d:v]xfade=transition=%s:duration=1:offset%d[vfade%d];",
				index, incrementIndex, transition, 0, incrementIndex,
			)
			xfadeAudioFilter += fmt.Sprintf(
				"[%d:a][%d:a]acrossfade=d=1[afade%d];",
				index, incrementIndex, incrementIndex,
			)
		}

		xfadeOffsets = append(xfadeOffsets, itemDurationFloat-1.0)
	}

	// concat all the processed items with the animations betweens inputs filter
	ffmpegTransitionsCmd := fmt.Sprintf(
		"ffmpeg %s -filter_complex \"%s\" -movflags +faststart %s",
		treatedAssets, xfadeVideoFilter+xfadeAudioFilter, concatOutputPath,
	)
	ffmpeg.Run(ffmpegTransitionsCmd)

	logger.RenderLogger.Info("concated videos successfully with xfade filters ✔️")

	// check if there is a required watermark asset
	if watermark != nil {
		logger.RenderLogger.Info("template watermark detected, generating final")

		finalOutputPath := fmt.Sprintf("/tmp/%s-%s.mp4", strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String())
		watermarkAsset := *watermark

		ffmpegOverlayCmd := fmt.Sprintf(
			"ffmpeg -i %s -i %s -filter_complex \"[0:v][1:v] overlay=W-w-10:H-h-10:enable='gte(t,1)'\" -pix_fmt yuv420p -c:a copy -vf scale %s",
			watermarkAsset, concatOutputPath, finalOutputPath,
		)
		ffmpeg.Run(ffmpegOverlayCmd)

		// remove file of the video without audio
		err := os.Remove(concatOutputPath)
		utils.CheckError("error deleting file", err)

		concatOutputPath = finalOutputPath

		logger.RenderLogger.Info("watermark detected, placed watermark image in the already concat video ✔️")
	}

	logger.RenderLogger.Info("concated videos successfully ✔️", zap.String("result", concatOutputPath))

	// finally return the final video
	return concatOutputPath
}
