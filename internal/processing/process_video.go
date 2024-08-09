package processing

import (
	"fmt"
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/utils"
	"os"
	"os/exec"

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

func CheckFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func RunFFmpeg(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func GetVideoDimensions(filePath string) (width, height int) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Can't get video dimensions of the file: %s. error=%v", filePath, err)
	}

	fmt.Sscanf(string(output), "%dx%d", &width, &height)
	return width, height
}

func GetVideoAspectRation(width, height int) AspectRation {
	if width > height {
		return Horizontal
	}
	return Vertical
}

func rerenderVideo(
	logfilePath string,
	inputPath string,
	outputPath string,
	minDimension int,
	width int,
	height int,
) error {
	var err error

	aspectRation := GetVideoAspectRation(width, height)

	switch aspectRation {

	case Vertical:
		scaleOption := fmt.Sprintf("scale=%d:%d", minDimension, -1)
		err = RunFFmpeg("-loglevel", "quiet", "-i", inputPath, "-map_metadata", "0", "-c:v", "h264_nvenc", "-vf", scaleOption, "-r", "30", "-y", outputPath)

	case Horizontal:
		scaleOption := fmt.Sprintf("scale=%d:%d", -1, minDimension)
		err = RunFFmpeg("-loglevel", "quiet", "-i", inputPath, "-map_metadata", "0", "-c:v", "h264_nvenc", "-vf", scaleOption, "-r", "30", "-y", outputPath)
	}
	if err != nil {
		log.Errorf("Error while processing video: %s", inputPath)
		errAppend := utils.AppendToLogfile(logfilePath, inputPath)
    // TODO if appending resulted in error is there a gracefull end?
		if errAppend != nil {
			log.Errorf(
				"Error while appending to the logfile: %s filepath: %s",
				logfilePath,
				inputPath,
			)
      return errAppend
		}
	}
  return nil
}

// Can't be async, always sync
// ffmpeg running on all cors by default?
func processVideo(
	params *parameters.Parameters,
	inputFilePath string,
) error {
	var err error

	outputFilePath := utils.GenOutputPath(params.OutputVideoDir, uuid.NewString(), "mp4")
	width, height := GetVideoDimensions(inputFilePath)

	// don't need to process
	// dimension is smaller than provided parameter MinVideoDimension
	if width < params.MinVideoDimension || height < params.MinVideoDimension {
		utils.CopyFile(inputFilePath, outputFilePath)
		return nil
	}

	err = rerenderVideo(
		params.LogFilePath,
		inputFilePath,
		outputFilePath,
		params.MinVideoDimension,
		width,
		height,
	)
	if err != nil {
    // TODO probably should remove processed file?
    return err
	}

	err = utils.TransferModificationTime(inputFilePath, outputFilePath)
	if err != nil {
    // TODO probably should remove processed file?
    // try luck next time?
		log.Warnf(
			"Error while transfering modification time from '%s' to '%s'",
			inputFilePath,
			outputFilePath,
		)
    return err
	}

  return nil
}
