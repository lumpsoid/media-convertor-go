package processing

import (
	"fmt"
	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/stats"
	"sync"

	"github.com/charmbracelet/log"
)

func Files(
	fileBucket *filebucket.FileBucket,
	params *parameters.Parameters,
	stats *stats.Stats,
) {
	log.Info("Starting processing files")
	var wg sync.WaitGroup
	resultCh := make(chan error, 1)

	counterLoop := 0
	counterMax := stats.PreCountImage + stats.PreCountVideo

	for ext, files := range fileBucket.Files {
		fileType := filebucket.GetFileType(ext)
		switch fileType {
		case filebucket.Image:
			for _, filePath := range files {
				counterLoop++
				fmt.Printf("\r%d/%d", counterLoop, counterMax)
				Image(filePath, params, &wg, resultCh)
			}
		case filebucket.Video:
			for _, filePath := range files {
				counterLoop++
				fmt.Printf("\r%d/%d", counterLoop, counterMax)
				Video(filePath, params, &wg, resultCh)
			}
		}
	}
	fmt.Print("\n")
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	// Process results from goroutines
	for err := range resultCh {
		if err != nil {
			log.Errorf("Error appending to file: %v\n", err)
		} else {
			log.Debug("File append successful")
		}
	}
}
