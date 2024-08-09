package filebucket

import (
	"bufio"
	"fmt"
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/utils"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type FileBucket struct {
	Files map[string][]string
}

type FileType int

const (
	Image     FileType = 0
	Video     FileType = 1
	Undefiend FileType = 2
)

func ReadLines(filePath string) ([]string, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read all lines into a slice
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func GetFileType(extension string) FileType {
  extString := strings.ToLower(extension)
	switch extString {
	case ".jpeg", ".jpg", ".png", ".heic", ".webp", ".avif":
		return Image
	case ".mov", ".mp4", ".avi", ".mkv":
		return Video
	default:
		log.Infof("Don't know what to do with: %s", extension)
		return Undefiend
	}
}

func FileBucketFromExtensions(extensions []string) *FileBucket {
	fileBucket := &FileBucket{
		Files: make(map[string][]string),
	}
	for _, ext := range extensions {
		extWithDot := fmt.Sprintf(".%s", ext)
		fileType := GetFileType(extWithDot)
		switch fileType {
		case Image:
			fileBucket.Files[extWithDot] = []string{}
		case Video:
			fileBucket.Files[extWithDot] = []string{}
		}
	}
	return fileBucket
}

func FileBucketFromLogFile(pathToLog string) *FileBucket {
	fileBucket := &FileBucket{
		Files: make(map[string][]string),
	}

	lines, err := ReadLines(pathToLog)
	if err != nil {
		log.Errorf("Can't read lines from log file: %s", pathToLog)
		os.Exit(1)
	}
	log.Infof("Read %d lines from log file: %s", len(lines), pathToLog)

  // TODO refator to something more meaningfull
	confirm, err := utils.ConfirmPrompt("With proceeding, logfile will be deleted. It's already read, but if you will cancel this run, this information will be lost. Do you want to proceed?")
	if err != nil {
		log.Errorf("Can't get confirmation from user")
	}
	if !confirm {
		os.Exit(1)
	}

	for _, pathToFile := range lines {
		ext := filepath.Ext(pathToFile)
		fileType := GetFileType(ext)
		switch fileType {
		case Image:
			fileBucket.Files[ext] = append(fileBucket.Files[ext], pathToFile)
		case Video:
			fileBucket.Files[ext] = append(fileBucket.Files[ext], pathToFile)
		}
	}

	err = os.Remove(pathToLog)
	if err != nil {
		log.Errorf("Can't delete log file: %s.\nThis will not affect the run.", pathToLog)
	}
	return fileBucket
}

func PopulateFileBucket(params *parameters.Parameters, fileBucket *FileBucket) {
	files, err := utils.GetFilesFromDir(params.InputDir)
	if err != nil {
		log.Errorf("Can't get files from input folder: %s", params.InputDir)
		os.Exit(1)
	}
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if _, ok := fileBucket.Files[ext]; ok {
			fileBucket.Files[ext] = append(
				fileBucket.Files[ext],
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


