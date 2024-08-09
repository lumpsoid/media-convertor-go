package stats

import (
	"fmt"
	"time"

	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/utils"

	"github.com/charmbracelet/log"
)

type Stats struct {
	PreCountImage  int
	PreCountVideo  int
	PostCountImage int
	PostCountVideo int
	StartTime      time.Time
}

func FromFileBucket(fileBucket *filebucket.FileBucket) Stats {
	stats := Stats{}
	for ext, files := range fileBucket.Files {
		fileType := filebucket.GetFileType(ext)
		switch fileType {
		case filebucket.Image:
			stats.PreCountImage += len(files)
		case filebucket.Video:
			stats.PreCountVideo += len(files)
		}
	}
	stats.StartTime = time.Now()
	return stats
}

func CountPost(
  stats *Stats, 
  outputImageDir string,
  outputVideoDir string,
) {
	postImageFiles, err := utils.GetFilesFromDir(outputImageDir)
	if err != nil {
		log.Error("Can't get files from output image folder")
	}
	postVideoFiles, err := utils.GetFilesFromDir(outputVideoDir)
	if err != nil {
		log.Error("Can't get files from output video folder")
	}
	stats.PostCountImage = len(postImageFiles)
	stats.PostCountVideo = len(postVideoFiles)
}

func CountFilesRecursive(dirPath string) int {
  counter, err := utils.CountFiles(dirPath)
  if err != nil {
    log.Errorf("Error while counting files: %v", err)
    return 0
  }
  return counter
}

func calculateTimeElapsed(startTime time.Time) {
	endTime := time.Now()
	log.Infof("Time elapsed: %s", endTime.Sub(startTime))
}

func Process(stats *Stats) {
	preCount := stats.PreCountImage + stats.PreCountVideo
	postCount := stats.PostCountImage + stats.PostCountVideo

	fmt.Printf("Pre-count: %d\n", preCount)
	fmt.Printf("Post-count: %d\n", postCount)
	if preCount == postCount {
		fmt.Print("Pre-count == Post-count\n")
	}
	calculateTimeElapsed(stats.StartTime)
}

