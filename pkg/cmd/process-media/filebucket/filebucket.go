package filebucket

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/charmbracelet/log"
)

type FileBucket struct {
	files map[string][]string
}

type FileType int

const (
	Image     FileType = 0
	Video     FileType = 1
	Undefiend FileType = 2
)

func getFileType(extension string) FileType {
	switch extension {
	case ".jpeg", ".jpg", ".png", ".heic", ".webp", ".avif":
		return Image
	case ".mov", ".mp4", ".avi", ".mkv":
		return Video
	default:
		log.Infof("Don't know what to do with: %s", extension)
		return Undefiend
	}
}

func fileBucketFromExtensions(extensions []string) *FileBucket {
	fileBucket := &FileBucket{
		files: make(map[string][]string),
	}
	for _, ext := range extensions {
		extWithDot := fmt.Sprintf(".%s", ext)
		fileType := getFileType(extWithDot)
		switch fileType {
		case Image:
			fileBucket.files[extWithDot] = []string{}
		case Video:
			fileBucket.files[extWithDot] = []string{}
		}
	}
	return fileBucket
}

func fileBucketFromLogFile(pathToLog string) *FileBucket {
	fileBucket := &FileBucket{
		files: make(map[string][]string),
	}

	lines, err := utils.readLines(pathToLog)
	if err != nil {
		log.Errorf("Can't read lines from log file: %s", pathToLog)
		os.Exit(1)
	}
	log.Infof("Read %d lines from log file: %s", len(lines), pathToLog)

	confirm, err := utils.confirmPrompt("With proceeding log file would be deleted. It's already read. But if you will cancel this run, this information will be lost. Do you want to proceed?")
	if err != nil {
		log.Errorf("Can't get confirmation from user")
	}
	if !confirm {
		os.Exit(1)
	}

	for _, pathToFile := range lines {
		ext := filepath.Ext(pathToFile)
		fileType := getFileType(ext)
		switch fileType {
		case Image:
			fileBucket.files[ext] = append(fileBucket.files[ext], pathToFile)
		case Video:
			fileBucket.files[ext] = append(fileBucket.files[ext], pathToFile)
		}
	}

	err = os.Remove(pathToLog)
	if err != nil {
		log.Errorf("Can't delete log file: %s.\nThis will not affect the run.", pathToLog)
	}
	return fileBucket
}

func populateFileBucket(params *Parameters, fileBucket *FileBucket) {
	files, err := utils.getFilesFromDir(params.InputDir)
	if err != nil {
		log.Errorf("Can't get files from input folder: %s", params.InputDir)
		os.Exit(1)
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if _, ok := fileBucket.files[ext]; ok {
			fileBucket.files[ext] = append(
				fileBucket.files[ext],
				path.Join(params.InputDir, file.Name()),
			)
		}
	}
}

// func countFilesWithExtensions(extensions ...string) int {
// 	count := 0
// 	for _, ext := range extensions {
// 		files, err := filepath.Glob("*." + ext)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		count += len(files)
// 	}
// 	return count
// }

// func getFilesWithExtensions(extensions ...string) []string {
// 	var files []string
// 	for _, ext := range extensions {
// 		filesWithType, err := filepath.Glob("*." + ext)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		files = append(files, filesWithType...)
// 	}
// 	return files
// }
