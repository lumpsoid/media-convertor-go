package stats

import (
	"fmt"
	"os"
	"testing"
)

func getExtensionWithDot(ext string) string {
	return fmt.Sprintf(".%s", ext)
}

func TestFileBucketFromExtensions(t *testing.T) {
	extensions := []string{"jpg", "mp4"}
	fileBucket := fileBucketFromExtensions(extensions)

	if _, ok := fileBucket.files[getExtensionWithDot(extensions[0])]; !ok {
		t.Errorf("Expected true, got %t", ok)
	}
	if _, ok := fileBucket.files[".mov"]; ok {
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
	params := Parameters{InputDir: "test"}
	fileBucket := fileBucketFromExtensions([]string{"jpg"})

	// test
	populateFileBucket(&params, fileBucket)

	if len(fileBucket.files[".jpg"]) != 2 {
		t.Errorf("Expected 2, got %d", len(fileBucket.files[".jpg"]))
	}

	// test
	fileBucket = fileBucketFromExtensions([]string{"jpg", "png"})
	populateFileBucket(&params, fileBucket)
	if len(fileBucket.files[".jpg"]) != 2 {
		t.Errorf("Expected 2, got %d", len(fileBucket.files[".jpg"]))
	}
	if len(fileBucket.files[".png"]) != 2 {
		t.Errorf("Expected 2, got %d", len(fileBucket.files[".png"]))
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
	params := Parameters{InputDir: "test"}
	fileBucket := fileBucketFromExtensions([]string{"jpg", "png"})
	populateFileBucket(&params, fileBucket)

	// test
	stats := statCalcFromFileBucket(fileBucket)
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
	statCalcPost(
		&stats,
		&Parameters{OutputImageDir: "test", OutputVideoDir: "test"},
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
