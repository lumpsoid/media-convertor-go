package stats

import (
	"fmt"
	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/parameters"
	"os"
	"testing"
)

func genFile(t *testing.T, filePath string) {
	err := os.WriteFile(filePath, []byte{}, 0644) // Set permissions to read-only for others (mode 0644)
	if err != nil {
		t.Errorf("Error creating file: %v", err)
	}
}

func getExtensionWithDot(ext string) string {
	return fmt.Sprintf(".%s", ext)
}

func TestFileBucketFromExtensions(t *testing.T) {
	extensions := []string{"jpg", "mp4"}
	fileBucket := filebucket.FileBucketFromExtensions(extensions)

	if _, ok := fileBucket.Files[getExtensionWithDot(extensions[0])]; !ok {
		t.Errorf("Expected true, got %t", ok)
	}
	if _, ok := fileBucket.Files[".mov"]; ok {
		t.Errorf("Expected false, got %t", ok)
	}
}

func TestPopulateFileBucket(t *testing.T) {
	filePaths := []string{
		"./test/test1.png",
		"./test/test2.png",
		"./test/test3.jpg",
		"./test/test4.jpg",
		"./test/test5.mov",
	}
	for _, path := range filePaths {
		genFile(t, path)
	}
	params := parameters.Parameters{InputDir: "test"}
	fileBucket := filebucket.FileBucketFromExtensions([]string{"jpg"})

	// test
	filebucket.PopulateFileBucket(&params, fileBucket)

	if len(fileBucket.Files[".jpg"]) != 2 {
		t.Errorf("Expected 2, got %d", len(fileBucket.Files[".jpg"]))
	}

	// test
	fileBucket = filebucket.FileBucketFromExtensions([]string{"jpg", "png"})
	filebucket.PopulateFileBucket(&params, fileBucket)
	if len(fileBucket.Files[".jpg"]) != 2 {
		t.Errorf("Expected 2, got %d", len(fileBucket.Files[".jpg"]))
	}
	if len(fileBucket.Files[".png"]) != 2 {
		t.Errorf("Expected 2, got %d", len(fileBucket.Files[".png"]))
	}

	for _, path := range filePaths {
		os.Remove(path)
	}
}

func TestStatCalcFromFileBucket(t *testing.T) {
	filePaths := []string{
		"./test/test1.png",
		"./test/test2.png",
		"./test/test3.jpg",
		"./test/test4.jpg",
		"./test/test5.mov",
	}
	for _, path := range filePaths {
		genFile(t, path)
	}
	params := parameters.Parameters{InputDir: "test"}
	fileBucket := filebucket.FileBucketFromExtensions([]string{"jpg", "png"})
	filebucket.PopulateFileBucket(&params, fileBucket)

	// test
	stats := FromFileBucket(fileBucket)
	if stats.PreCountImage != len(filePaths)-1 {
		t.Errorf("Expected %d, got %d", len(filePaths)-1, stats.PreCountImage)
	}

	for _, path := range filePaths {
		os.Remove(path)
	}
}

func TestStatCalcPost(t *testing.T) {
	filePaths := []string{
		"./test/test1.png",
		"./test/test2.png",
		"./test/test3.jpg",
		"./test/test4.jpg",
		"./test/test5.mov",
	}
	for _, path := range filePaths {
		genFile(t, path)
	}
	stats := Stats{}

	// test
	CountPost(
		&stats,
		"test", 
    "test",
	)
	if stats.PostCountImage != len(filePaths) {
		t.Errorf("Expected %d, got %d", len(filePaths), stats.PostCountImage)
	}
	if stats.PostCountVideo != len(filePaths) {
		t.Errorf("Expected %d, got %d", len(filePaths), stats.PostCountImage)
	}

	for _, path := range filePaths {
		os.Remove(path)
	}
}
