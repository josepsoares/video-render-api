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

// todo - need to add the template to this, also need to add template files to the assets
func TestConcatTemplate(t *testing.T) {
	defer utils.TimeTrack(time.Now(), "Test - ConcatPresetAssetsWithTemplate")

	fmt.Println("starting concat of preset assets to create a simple video")

	err := godotenv.Load("../../.env")
	utils.CheckError("couldn't load .env file", err)

	xfadeOffsets := []float64{}
	assets := []string{}
	// template :=

	var (
		xfadeVideoFilter string
		xfadeAudioFilter string
		cmdAssets        string
	)

	outputPath := fmt.Sprintf(
		"%s/tmp/%s-%s-concat-p-assets-template.mp4",
		os.Getenv("LOCAL_BASE_PATH"), strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String(),
	)

	// concat assets to create sample UA visit video
	for index, item := range assets {
		incrementIndex := index + 1
		cmdAssets += fmt.Sprintf("-i %s ", item)

		itemDurationString := ffmpeg.RunProbe(fmt.Sprintf(
			"ffprobe -i %s -v quiet -show_entries format=duration -hide_banner -of default=noprint_wrappers=1:nokey=1", item),
		)
		itemDurationFloat, _ := strconv.ParseFloat(itemDurationString, 64)

		if index != 0 {
			xfadeVideoFilter += fmt.Sprintf(
				"[vfade%d][%d:v]xfade=transition=fadeblack:duration=1:offset=%f[vfade%d];",
				index, incrementIndex, (itemDurationFloat+xfadeOffsets[index-1])-1.0, incrementIndex,
			)
			xfadeAudioFilter += fmt.Sprintf(
				"[afade%d][%d:a]acrossfade=d=1[afade%d];",
				index, incrementIndex, incrementIndex,
			)
		} else {
			xfadeVideoFilter += fmt.Sprintf(
				"[%d][%d:v]xfade=transition=fadeblack:duration=1:offset%d[vfade%d];",
				index, incrementIndex, 0, incrementIndex,
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
		cmdAssets, xfadeVideoFilter+xfadeAudioFilter, outputPath,
	)
	ffmpeg.Run(ffmpegTransitionsCmd)

	fmt.Println("ended test video, output path:" + outputPath)
}
