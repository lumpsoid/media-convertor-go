package processing

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type AspectRation int

const (
	Horizontal AspectRation = 0
	Vertical   AspectRation = 1
)

const (
	UnknownDirection = iota
	North
	South
	East
	West
)

// Video struct
// type Video struct {
// 	tempFileName      uuid.UUID
// 	InputFile         string
// 	CreationDate      string
// 	LastModification  string
// 	MinVideoDimension int
// }

// ConvertParametersToVideo converts Parameters to Video
// func ConvertParametersToVideo(params Parameters) Video {
// 	return Video{
// 		InputDir:          params.InputDir,
// 		OutputDir:         params.OutputDir,
// 		MinVideoDimension: params.MinVideoDimension,
// 	}
// }

// func newVideo(inputFile string) *Video {
// 	getVideoDimensions(inputFile)

// 	video := &Video{
// 		tempFileName: uuid.New(),
// 	}
// 	return video
// }

func checkFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func runFFmpeg(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func getVideoDimensions(filePath string) (width, height int) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Can't get video dimensions of the file: %s. error=%v", filePath, err)
	}

	fmt.Sscanf(string(output), "%dx%d", &width, &height)
	return width, height
}

func getVideoAspectRation(width, height int) AspectRation {
	if width > height {
		return Horizontal
	}
	return Vertical
}

func processVideo(filePath string, params *Parameters, wg *sync.WaitGroup, resultCh chan error) {
	outputFile := utils.genOutputPath(params.OutputVideoDir, uuid.NewString(), "mp4")
	width, height := getVideoDimensions(filePath)
	if width < params.MinVideoDimension || height < params.MinVideoDimension {
		copyFile(filePath, outputFile)
	}

	aspectRation := getVideoAspectRation(width, height)
	modTime, err := getModificationTime(filePath)
	if err != nil {
		log.Errorf("Can't get modification time of the file: %s", filePath)
	}

	switch aspectRation {
	case Vertical:
		scaleOption := fmt.Sprintf("scale=%d:%d", params.MinVideoDimension, -1)
		err = runFFmpeg("-loglevel", "quiet", "-i", filePath, "-map_metadata", "0", "-c:v", "h264_nvenc", "-vf", scaleOption, "-r", "30", "-y", outputFile)
	case Horizontal:
		scaleOption := fmt.Sprintf("scale=%d:%d", -1, params.MinVideoDimension)
		err = runFFmpeg("-loglevel", "quiet", "-i", filePath, "-map_metadata", "0", "-c:v", "h264_nvenc", "-vf", scaleOption, "-r", "30", "-y", outputFile)
	}
	if err != nil {
		log.Errorf("Error while processing video: %s", filePath)
		wg.Add(1)
		logFileName := fmt.Sprintf(".processError-%s", params.startTime.Local().Format("2006-01-02-15-04-05"))
		go appendToFileAsync(path.Join(params.InputDir, logFileName), filePath, wg, resultCh)
	}
	os.Chtimes(outputFile, time.Time{}, modTime)
}
