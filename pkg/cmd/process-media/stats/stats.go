package stats

import (
	"fmt"

	"github.com/charmbracelet/log"
)

type Stats struct {
	PreCountImage  int
	PreCountVideo  int
	PostCountImage int
	PostCountVideo int
}

func statCalcFromFileBucket(fileBucket *FileBucket) Stats {
	stats := Stats{}
	for ext, files := range fileBucket.files {
		fileType := getFileType(ext)
		switch fileType {
		case Image:
			stats.PreCountImage += len(files)
		case Video:
			stats.PreCountVideo += len(files)
		}
	}
	return stats
}

func statCalcPost(stats *Stats, params *Parameters) {
	postImageFiles, err := getFilesFromDir(params.OutputImageDir)
	if err != nil {
		log.Error("Can't get files from output image folder")
	}
	postVideoFiles, err := getFilesFromDir(params.OutputVideoDir)
	if err != nil {
		log.Error("Can't get files from output video folder")
	}
	stats.PostCountImage = len(postImageFiles)
	stats.PostCountVideo = len(postVideoFiles)
}

func processStats(stats *Stats) {
	preCount := stats.PreCountImage + stats.PreCountVideo
	postCount := stats.PostCountImage + stats.PostCountVideo

	fmt.Printf("Pre-count: %d\n", preCount)
	fmt.Printf("Post-count: %d\n", postCount)
	if preCount == postCount {
		fmt.Print("Pre-count == Post-count\n")
	}
}
