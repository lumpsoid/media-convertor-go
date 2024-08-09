package processing

import (
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/utils"
	"os"
	"os/exec"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func Image(
	filePath string,
	params *parameters.Parameters,
	wg *sync.WaitGroup,
	resultCh chan error,
) {
	outputPath := utils.GenOutputPath(params.OutputImageDir, uuid.NewString(), "avif")

	cmd := exec.Command("magick", "-quality", "57", filePath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("Error while processing image: %s", filePath)
		wg.Add(1)
		go utils.AppendToFileAsync(params.LogFilePath, filePath, wg, resultCh)
	}
	err = utils.TransferModificationTime(filePath, outputPath)
	if err != nil {
		log.Warnf("Error while transfering modification time from: '%s' to: '%s'", filePath, outputPath)
	}
}
