package processing

import (
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/utils"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func processImage(
	params *parameters.Parameters,
	inputFilepath string,
) error {
	outputPath := utils.GenOutputPath(
		params.OutputImageDir,
		uuid.NewString(),
		params.ImageTargetFormat,
	)

	cmd := exec.Command(
		"magick",
		"-quality",
		strconv.Itoa(params.ImageTargetQuality),
		inputFilepath,
		outputPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Errorf("Error while processing image: %s", inputFilepath)
		errAppend := utils.AppendToLogfile(params.LogFilePath, inputFilepath)
		if errAppend != nil {
			return err
		}
	}

	err = utils.TransferModificationTime(inputFilepath, outputPath)
	if err != nil {
		log.Warnf("Error while transfering modification time from: '%s' to: '%s'", inputFilepath, outputPath)
		return err
	}

	return nil
}
