package utils

import (
	"fmt"
	"time"
)

func PrependInt(x []int, y int) []int {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
}

func GetResolutionFilters(res int) string {
	var filters string

	if res == 1080 {
		filters = "-b:v 4500k -minrate 4500k -maxrate 9000k -bufsize 9000k"
	} else if res == 720 {
		filters = "-b:v 2500k -minrate 1500k -maxrate 4000k -bufsize 5000k"
	} else if res == 480 {
		filters = "-b:v 1000k -minrate 500k -maxrate 2000k -bufsize 2000k"
	} else {
		filters = "-b:v 750k -minrate 400k -maxrate 1000k -bufsize 1500k"
	}

	return filters
}
