package filebucket

import (
	"bufio"
	"io/fs"
	"mediaconvertor/internal/utils"
	"os"
	"path"
	"path/filepath"

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

func (fb *FileBucket) InputFile(inputDir string, dirEntry fs.DirEntry) {
	ext := utils.GetFileExtension(dirEntry.Name())

	// if extension is in the bucket
	// then we append it for processing
	if _, ok := fb.Files[ext]; ok {
		fb.Files[ext] = append(
			fb.Files[ext],
			path.Join(inputDir, dirEntry.Name()),
		)
	}
}

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
	switch extension {
	case "jpeg", "jpg", "png", "heic", "webp", "avif", "jxl":
		return Image
	case "mov", "mp4", "avi", "mkv":
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
		fileType := GetFileType(ext)
		switch fileType {
		case Image:
			fileBucket.Files[ext] = []string{}
		case Video:
			fileBucket.Files[ext] = []string{}
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

func PopulateFileBuketRecursive(fb *FileBucket, inputDir string) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal("Error reading dir: ", err)
	}

	// Iterate through each file and directory
	for _, file := range files {
		fullPath := filepath.Join(inputDir, file.Name())

		if file.IsDir() {
			PopulateFileBuketRecursive(fb, fullPath)
		} else {
			fb.InputFile(inputDir, file)
		}
	}
	return 
}

func PopulateFileBucket(fileBucket *FileBucket, inputDir string) {
	files, err := utils.GetFilesFromDir(inputDir)
	if err != nil {
		log.Errorf("Can't get files from input folder: %s", inputDir)
		os.Exit(1)
	}
	for _, file := range files {
		fileBucket.InputFile(inputDir, file)
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
