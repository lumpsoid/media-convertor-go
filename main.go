package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/log"
)

func checkStartup(params Parameters) {
	if !isPathExist(params.InputDir) {
		log.Info("-input_dir is not exist")
		os.Exit(1)
	}

	checkAndCleanDir(params.OutputImageDir, params.NoClean, "Output image")
	checkAndCleanDir(params.OutputVideoDir, params.NoClean, "Output video")
}

func setUpFiles(params *Parameters) (Stats, *FileBucket) {
	if len(params.FromLogFile) > 0 {
		fileBucket := fileBucketFromLogFile(params.FromLogFile)
		stats := statCalcFromFileBucket(fileBucket)
		return stats, fileBucket
	}
	fileBucket := fileBucketFromExtensions(params.Extensions)
	populateFileBucket(params, fileBucket)
	stats := statCalcFromFileBucket(fileBucket)
	return stats, fileBucket
}

func processFiles(fileBucket *FileBucket, params *Parameters, stats *Stats) {
	log.Info("Starting processing files")
	var wg sync.WaitGroup
	resultCh := make(chan error, 1)

	counterLoop := 0
	counterMax := stats.PreCountImage + stats.PreCountVideo

	for ext, files := range fileBucket.files {
		fileType := getFileType(ext)
		switch fileType {
		case Image:
			for _, filePath := range files {
				counterLoop++
				fmt.Printf("\r%d/%d", counterLoop, counterMax)
				processImage(filePath, params, &wg, resultCh)
			}
		case Video:
			for _, filePath := range files {
				counterLoop++
				fmt.Printf("\r%d/%d", counterLoop, counterMax)
				processVideo(filePath, params, &wg, resultCh)
			}
		}
	}
	fmt.Print("\n")
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	// Process results from goroutines
	for err := range resultCh {
		if err != nil {
			log.Errorf("Error appending to file: %v\n", err)
		} else {
			log.Debug("File append successful")
		}
	}
}

func main() {
	// Parse parameters
	params := parseParameters()

	// params.OutputDir := "./post/"
	checkStartup(params)

	stats, fileBucket := setUpFiles(&params)

	processFiles(fileBucket, &params, &stats)

	statCalcPost(&stats, &params)

	// -r -- recursive
	runExiftool("Run exiftool for images", "-ext", "AVIF", "-FileName<${FileModifyDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", params.OutputImageDir)
	runExiftool("Run exiftool for images, trying CreateDate", "-ext", "AVIF", "-FileName<${CreateDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", params.OutputImageDir)
	runExiftool("Run exiftool for images, trying DateCreated", "-ext", "AVIF", "-FileName<${DateCreated}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", params.OutputImageDir)

	runExiftool("Run exiftool for videos", "-ext", "MP4", "-FileName<${FileModifyDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "MOV_%Y-%m-%d_%H-%M-%S", params.OutputVideoDir)
	runExiftool("Run exiftool for videos, trying CreateDate", "-ext", "MP4", "-FileName<${CreateDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "MOV_%Y-%m-%d_%H-%M-%S", params.OutputVideoDir)

	processStats(&stats)
	calcTimeElapsed(params.startTime)
}
