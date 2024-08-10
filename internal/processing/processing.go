package processing

import (
	"fmt"
	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/stats"
	"mediaconvertor/internal/utils"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func Files(
	fileBucket *filebucket.FileBucket,
	params *parameters.Parameters,
	stats *stats.Stats,
) {
	log.Info("Starting processing files")

	counterLoop := 0
	counterMax := stats.PreCountImage + stats.PreCountVideo + stats.PreCountUndefiend

	for ext, files := range fileBucket.Files {
		var errInLoop error
		fileType := filebucket.GetFileType(ext)
		switch fileType {
		case filebucket.Image:
			for _, filePath := range files {
				counterLoop++
				fmt.Printf("\r%d/%d", counterLoop, counterMax)
				errInLoop = processImage(params, filePath)
			}
		case filebucket.Video:
			for _, filePath := range files {
				counterLoop++
				fmt.Printf("\r%d/%d", counterLoop, counterMax)
				errInLoop = processVideo(params, filePath)
			}
		case filebucket.Undefiend:
			for _, filePath := range files {
				counterLoop++
				fmt.Printf("\r%d/%d", counterLoop, counterMax)
        extNew := utils.GetFileExtension(filePath)
				outputFilePath := utils.GenOutputPath(
					params.OutputUndefiendDir,
					uuid.NewString(),
          extNew,
				)
				utils.CopyFile(filePath, outputFilePath)
			}
		}
		if errInLoop != nil {
			// TODO error in two cases
			// failed to append to the file
			// failed to migrate time modification from input file to output file
			// probably better to just delete output file
			// and try luck in appending to the log file second time
			// maybe create class for handling writing to the log file?
			// then we can queue into it
			// if will be able to async images
			// then we can just unload this task to him
			// and don't think too much here
			// make him running in a separate goroutine
			log.Error("Failed in most ugly way")
		}
	}
	fmt.Print("\n")
}

func CopyFiles(
	fileBucket *filebucket.FileBucket,
	stats *stats.Stats,
	outputStructured string,
	saveFileName bool,
) {

	counterLoop := 0
	counterMax := stats.PreCountImage + stats.PreCountVideo + stats.PreCountUndefiend

	for ext, files := range fileBucket.Files {
		for _, filePath := range files {
			counterLoop++
			fmt.Printf("\r%d/%d", counterLoop, counterMax)

			var outputPath string
			if saveFileName {
				baseName := filepath.Base(filePath)
				//ext := filepath.Ext(baseName)
				// Remove the extension from the base name
				//nameWithoutExt := baseName[:len(baseName)-len(ext)]
				outputPath = filepath.Join(outputStructured, baseName)
			} else {
				outputPath = utils.GenOutputPath(
					outputStructured,
					uuid.NewString(),
					ext,
				)
			}

			utils.CopyFile(filePath, outputPath)
		}
	}
	fmt.Print("\n")
}
