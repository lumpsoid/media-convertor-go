package processing

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

func genImage(imgPath string, width, height int) {
	// Create a new RGBA image with the specified width and height
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Set all pixels to black
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// Create a new PNG file
	file, err := os.Create(imgPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image as PNG and write it to the file
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func TestCheckFFmpeg(t *testing.T) {
	err := checkFFmpeg()
	if err != nil {
		t.Fatalf("FFmpeg is not installed. %s", err)
	}
}

func TestGetVideoDimensions(t *testing.T) {
	filePath := "test/TEST.PNG"
	widthTest := 800
	heightTest := 600
	genImage(filePath, widthTest, heightTest)

	width, height := getVideoDimensions(filePath)
	if width != widthTest {
		t.Errorf("Width of the images don't equal. Test image: %d. Got: %d", widthTest, width)
	}
	if height != heightTest {
		t.Errorf("Height of the images don't equal. Test image: %d. Got: %d", heightTest, height)
	}
	os.Remove(filePath)
}

func TestGetVideoAspectRation(t *testing.T) {
	widthTest := 800
	heightTest := 600

	aspectRation := getVideoAspectRation(widthTest, heightTest)
	if aspectRation != Horizontal {
		t.Error("Must be horizontal, got vertical")
	}
}
