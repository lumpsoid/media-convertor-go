package main

import (
	"mediaconvertor/internal/orginizer"
	"mediaconvertor/internal/stats"

	"github.com/charmbracelet/log"
)

func main() {
	statistics, params, fb := orginizer.Initialize()

	log.Info("Starting orginizing files")
  orginizer.CopyFiles(fb, statistics, params.OutputStructured)

	log.Info("Starting orginizing files")
	orginizer.StructureOutputLayout(params.OutputStructured)

  // no difference PostCountImage or PostCountVideo
  // we are processing them together
  statistics.PostCountImage = stats.CountFilesRecursive(params.OutputStructured)
  
	stats.Process(statistics)
}
