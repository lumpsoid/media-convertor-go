package processing

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func processImage(filePath string, params *Parameters, wg *sync.WaitGroup, resultCh chan error) {
	outputPath := genOutputPath(params.OutputImageDir, uuid.NewString(), "avif")

	cmd := exec.Command("convert", "-quality", "57", filePath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Errorf("Error while processing image: %s", filePath)
		wg.Add(1)
		logFileName := fmt.Sprintf(".processError-%s", params.startTime.Local().Format("2006-01-02-15-04-05"))
		go appendToFileAsync(path.Join(params.InputDir, logFileName), filePath, wg, resultCh)
	}
	err = transferModificationTime(filePath, outputPath)
	if err != nil {
		log.Warnf("Error while transfering modification time from: '%s' to: '%s'", filePath, outputPath)
	}
}
