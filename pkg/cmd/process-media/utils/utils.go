package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

func getModificationTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}

func transferModificationTime(filePathOne, filePathTwo string) error {
	timeModOne, err := getModificationTime(filePathOne)
	if err != nil {
		return err
	}
	os.Chtimes(filePathTwo, timeModOne, timeModOne)

	return nil
}

func copyFile(srcPath, destPath string) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}

	// Get source file info to obtain modification time
	srcFileInfo, err := os.Stat(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	// Set the modification time of the destination file
	err = os.Chtimes(destPath, srcFileInfo.ModTime(), srcFileInfo.ModTime())
	if err != nil {
		log.Fatal(err)
	}
}

// func deleteExt(fileName string) string {
// 	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
// }

// func checkFileExist(filePath string) bool {
// 	_, err := os.Stat(filePath)
// 	return err == nil
// }

// func joinOutputPath(baseName string, ext string, suffix string) string {
// 	return fmt.Sprintf("%s%s.%s", baseName, suffix, ext)
// }

func genOutputPath(outDir, baseName, ext string) string {
	filename := fmt.Sprintf("%s.%s", baseName, ext)
	return path.Join(outDir, filename)
}

func isPathExist(dirPath string) bool {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Debug("Directory don't exist: %s", dirPath)
		return false
	}
	return true
}

func createFolder(dirPath string) error {
	err := os.Mkdir(dirPath, 0755)
	return err
}

func getFilesFromDir(dirPath string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func appendToFileAsync(filePath string, content string, wg *sync.WaitGroup, resultCh chan error) {
	defer wg.Done()
	if len(content) > 0 && content[len(content)-1] != '\n' {
		content += "\n"
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		resultCh <- err
		return
	}
	defer file.Close()
	_, err = file.WriteString(content)
	resultCh <- err
}

func readLines(filePath string) ([]string, error) {
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

func calcTimeElapsed(startTime time.Time) {
	endTime := time.Now()
	log.Infof("Time elapsed: %s", endTime.Sub(startTime))
}

func cleanDir(dirPath string) {
	files, err := getFilesFromDir(dirPath)
	if err != nil {
		log.Errorf("Can't get files from folder: %s", dirPath)
		os.Exit(1)
	}

	if len(files) > 0 {
		confirm, err := confirmPrompt("Output image folder has files, do you want to delete them?")
		if err != nil {
			log.Error("Can't get confirmation from user")
			os.Exit(1)
		}
		if confirm {
			err = os.RemoveAll(dirPath)
			if err != nil {
				log.Errorf("Can't delete directory: %s. Error: %s", dirPath, err)
				os.Exit(1)
			}
			err := createFolder(dirPath)
			if err != nil {
				log.Errorf("Can't create folder: %s", dirPath)
				os.Exit(1)
			}
		} else {
			log.Infof("Move files from: '%s' and try again.", dirPath)
			os.Exit(1)
		}
	}
}

func checkAndCleanDir(path string, no_clean bool, messagePrefix string) {
	// not existing
	if !isPathExist(path) {
		err := createFolder(path)
		if err != nil {
			log.Errorf("Can't create %s folder: %v", path, err)
			os.Exit(1)
		}
		return
	}

	if no_clean {
		return
	}
	log.Infof("%s folder already exists", messagePrefix)
	cleanDir(path)

}
