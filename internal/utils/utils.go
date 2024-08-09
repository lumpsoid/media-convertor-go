package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

func GetModificationTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}

func TransferModificationTime(filePathOne, filePathTwo string) error {
	timeModOne, err := GetModificationTime(filePathOne)
	if err != nil {
		return err
	}
	os.Chtimes(filePathTwo, timeModOne, timeModOne)

	return nil
}

func CopyFile(srcPath, destPath string) {
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

func GenOutputPath(outDir, baseName, ext string) string {
	filename := fmt.Sprintf("%s.%s", baseName, ext)
	return path.Join(outDir, filename)
}

func IsPathExist(dirPath string) bool {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Debug("Directory don't exist: %s", dirPath)
		return false
	}
	return true
}

func CreateFolder(dirPath string) error {
	err := os.Mkdir(dirPath, 0755)
	return err
}

func GetFilesFromDir(dirPath string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func AppendToLogfile(logfilePath string, mediaFilepath string) error {
  if len(mediaFilepath) == 0 {
    return errors.New("Media filepath is empty")
  }
	if mediaFilepath[len(mediaFilepath)-1] != '\n' {
		mediaFilepath += "\n"
	}

	file, err := os.OpenFile(logfilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(mediaFilepath)
  if err != nil {
    return err
  }
  return nil
}

func AppendToFileAsync(filePath string, content string, wg *sync.WaitGroup, resultCh chan error) {
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

// Clear the dirPath
// then Create the empty folder
func ClearCreate(dirPath string) {
	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Errorf("Can't delete directory: %s. Error: %s", dirPath, err)
		os.Exit(1)
	}
	err = CreateFolder(dirPath)
	if err != nil {
		log.Errorf("Can't create folder: %s", dirPath)
		os.Exit(1)
	}
}

func CheckAndClearDir(dirPath string, overrideOutputDir bool, messagePrefix string) {
	// not existing
	if !IsPathExist(dirPath) {
		err := CreateFolder(dirPath)
		if err != nil {
			log.Errorf("Can't create %s folder: %v", dirPath, err)
			os.Exit(1)
		}
		return
	}
	log.Infof("%s already exists", dirPath)

	files, err := GetFilesFromDir(dirPath)
	if err != nil {
		log.Errorf("Can't get files from folder: %s", dirPath)
		os.Exit(1)
	}

	// the folder is clean
	// we safe to continue
	if len(files) == 0 {
		return
	}

	if overrideOutputDir {
		ClearCreate(dirPath)
		return
	}

	confirm, err := ConfirmPrompt(fmt.Sprintf("%s has files, do you want to delete them?", messagePrefix))
	if err != nil {
		log.Error("Can't get confirmation from user")
		os.Exit(1)
	}

	if confirm {
		ClearCreate(dirPath)
		return
	} else {
		log.Infof("Move files from: '%s' and try again.", dirPath)
		os.Exit(1)
	}
}
