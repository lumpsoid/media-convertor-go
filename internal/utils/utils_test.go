package utils

import (
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

func genFile(t *testing.T, filePath string) {
	err := os.WriteFile(filePath, []byte{}, 0644) // Set permissions to read-only for others (mode 0644)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
}

func genTestFolder(t *testing.T) (time.Time, string) {
	timeNow := time.Now()
	folderPath := fmt.Sprintf("../../test/%d", timeNow.UnixNano())
  err := CreateFolder(folderPath)
  if err != nil {
    t.Fatal(err)
  }
	return timeNow, folderPath
}

func TestGetModificationTime(t *testing.T) {
	// set up
	timeNow, folderPath := genTestFolder(t)
	fileName := "data.txt"
	filePath := path.Join(folderPath, fileName)
	genFile(t, filePath)
	err := os.Chtimes(filePath, timeNow, timeNow)
	if err != nil {
		t.Errorf("Error changing modification time: %v", err)
	}

	// test
	timeMod, err := GetModificationTime(filePath)
	if err != nil {
		t.Errorf("Error get modification time: %v", err)
	}
	if !timeMod.Equal(timeNow) {
		t.Fatalf("modified time of the file: %d, expected modified time: %d", timeMod.Unix(), timeNow.Unix())
	}
	// clean up
	os.RemoveAll(folderPath)
}

func TestPathExist(t *testing.T) {
	_, folderPath := genTestFolder(t)
	fileName := "test.999"
	filePath := path.Join(folderPath, fileName)
	genFile(t, filePath)

	pathExist := IsPathExist(filePath)
	if pathExist != true {
		t.Fatalf("Path: %s don't exist. Check the file, or function is broken.", filePath)
	}

	os.RemoveAll(folderPath)
}

func TestGenOutputPath(t *testing.T) {
	outDir := "test"
	baseName := "TESTING"
	extName := "EXT"
	path := GenOutputPath(outDir, baseName, extName)
	testPath := "test/TESTING.EXT"
	if path != testPath {
		t.Fatalf("Paths are not equal.\n%s != %s", path, testPath)
	}
}

func TestCopyFile(t *testing.T) {
	// set up
	_, folderPath := genTestFolder(t)
	filePath := GenOutputPath(folderPath, "test", "file")
	outputPath := GenOutputPath(folderPath, "testTEST", "file")
	genFile(t, filePath)

	// test
	CopyFile(filePath, outputPath)
	timeModInit, err := GetModificationTime(filePath)
	if err != nil {
		t.Fatalf("Can't get modification time of the file: %s", filePath)
	}
	timeModOut, err := GetModificationTime(outputPath)
	if err != nil {
		t.Fatalf("Can't get modification time of another file: %s", outputPath)
	}
	pathExist := IsPathExist(outputPath)
	if pathExist == false {
		t.Fatalf("File exist: %t", pathExist)
	}
	timeModEqual := timeModInit.Equal(timeModOut)
	if timeModEqual == false {
		t.Fatalf("Time modification equality is %t", timeModEqual)
	}

	// clean up
	os.RemoveAll(folderPath)
}

func TestCreateFolder(t *testing.T) {
	folderPath := "test/testFolder"
	err := CreateFolder(folderPath)
	if err != nil {
		t.Fatalf("Can't create folder: %s", err)
	}
	pathExist := IsPathExist(folderPath)
	if pathExist == false {
		t.Fatalf("Folder doesn't exist")
	}
	err = CreateFolder(folderPath)
	if err == nil {
		t.Fatal("Function must fail, because already exist")
	}
	os.Remove(folderPath)
}

func TestGetFilesFromDir(t *testing.T) {
	filePath := "test/lenfunc.txt"
	genFile(t, filePath)
	files, err := GetFilesFromDir("test")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if len(files) != 1 {
		t.Fatalf("Files count: %d", len(files))
	}

	// clean up
	os.Remove(filePath)
}

func TestAppendToFileAsync(t *testing.T) {
	// set up
	_, folderPath := genTestFolder(t)
	fileName := "TEST"
	fileExt := "text"
	filePath := GenOutputPath(folderPath, fileName, fileExt)
	genFile(t, filePath)
	content := "test"

	// test
	var wg sync.WaitGroup
	resultCh := make(chan error, 1)

	wg.Add(1)
	go AppendToFileAsync(filePath, content, &wg, resultCh)

	go func() {
		wg.Wait()
		close(resultCh)
		os.Remove(filePath)
	}()
	// Process results from goroutines
	for err := range resultCh {
		if err != nil {
			t.Fatalf("Error appending to file: %v\n", err)
		}
	}

	// clean up
	os.RemoveAll(folderPath)
}

func TestSpecGetModificationTime(t *testing.T) {
	// set up
	_, folderPath := genTestFolder(t)
	filePathOne := GenOutputPath(folderPath, "test1", "txt")
	filePathTwo := GenOutputPath(folderPath, "test2", "txt")
	genFile(t, filePathOne)
	time.Sleep(time.Second)
	genFile(t, filePathTwo)

	// test
	TransferModificationTime(filePathOne, filePathTwo)
	timeModOne, err := GetModificationTime(filePathOne)
	if err != nil {
		t.Errorf("Error while getting modification time: %v", err)
	}
	timeModTwo, err := GetModificationTime(filePathTwo)
	if err != nil {
		t.Errorf("Error while getting modification time: %v", err)
	}
	if !timeModOne.Equal(timeModTwo) {
		t.Errorf("Mod time different. One: %d. Two: %d.", timeModOne.Unix(), timeModTwo.Unix())
	}

	// clean up
	os.RemoveAll(folderPath)
}

func TestGetDiskSpace(t *testing.T) {
	free, total, used, err  := GetDiskSpace("/")
  // Convert bytes to megabytes
	const bytesInMB = 1024 * 1024
  freeMB := free / bytesInMB
  totalMB := total / bytesInMB
  usedMB := used / bytesInMB
  t.Fatal(freeMB, totalMB, usedMB, err)
}
