package render

import (
	"flag"
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

func TestConcatStress(t *testing.T) {
	defer utils.TimeTrack(time.Now(), "Test - StressConcatTest")

	fmt.Println("starting stress test of concat of assets to see the performance of the algorithm")

	err := godotenv.Load("../../.env")
	utils.CheckError("couldn't load .env file", err)

	xfadeOffsets := []float64{}

	assets := []string{
		fmt.Sprintf("%s/assets/videos/ua_intro.mp4", os.Getenv("LOCAL_BASE_PATH")),
		fmt.Sprintf("%s/assets/videos/ua_cp.mp4", os.Getenv("LOCAL_BASE_PATH")),
		fmt.Sprintf("%s/assets/videos/ua_engcivil.mp4", os.Getenv("LOCAL_BASE_PATH")),
		fmt.Sprintf("%s/assets/videos/biblioteca.mp4", os.Getenv("LOCAL_BASE_PATH")),
		fmt.Sprintf("%s/assets/videos/ua_outro.mp4", os.Getenv("LOCAL_BASE_PATH")),
	}

	var (
		resFlagValue     int
		xfadeVideoFilter string
		xfadeAudioFilter string
		cmdAssets        string
	)

	flag.IntVar(&resFlagValue, "res", 1080, "resolution desired for the execution of the test")

	for i := 1; i < 25; i++ {
		outputPath := fmt.Sprintf(
			"%s/tmp/%s-%s-%dp-concat-stress.mp4",
			os.Getenv("LOCAL_BASE_PATH"), strconv.FormatInt(time.Now().Unix(), 10), uuid.New().String(), resFlagValue,
		)

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
	}

	fmt.Println("ended concat stress test, files are in /tmp folder")
}
