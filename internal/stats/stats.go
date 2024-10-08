package stats

import (
	"time"

	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/utils"

	"github.com/charmbracelet/log"
)

type Stats struct {
	PreCountImage      int
	PreCountVideo      int
	PreCountUndefiend  int
	PostCountImage     int
	PostCountVideo     int
	PostCountUndefiend int
	StartTime          time.Time
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
		case filebucket.Undefiend:
			stats.PreCountUndefiend += len(files)
		}
	}
	stats.StartTime = time.Now()
	return stats
}

func CountPost(
	stats *Stats,
	outputImageDir string,
	outputVideoDir string,
	outputUndefiendDir string,
) {
	postImageFiles, err := utils.GetFilesFromDir(outputImageDir)
	if err != nil {
		log.Error("Can't get files from output image folder")
	}
	postVideoFiles, err := utils.GetFilesFromDir(outputVideoDir)
	if err != nil {
		log.Error("Can't get files from output video folder")
	}
	postUndefiendFiles, err := utils.GetFilesFromDir(outputUndefiendDir)
	if err != nil {
		log.Error("Can't get files from output undefiend folder")
	}
	stats.PostCountImage = len(postImageFiles)
	stats.PostCountVideo = len(postVideoFiles)
	stats.PostCountUndefiend = len(postUndefiendFiles)
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
	preCount := stats.PreCountImage + stats.PreCountVideo + stats.PreCountUndefiend
	postCount := stats.PostCountImage + stats.PostCountVideo + stats.PreCountUndefiend

  log.Info("Pre-count: ", preCount, "Post-count: ", postCount)
	if preCount == postCount {
    log.Info("Pre-count == Post-count")
	}
	calculateTimeElapsed(stats.StartTime)
}
