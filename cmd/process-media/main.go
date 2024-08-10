package main

import (
	"mediaconvertor/internal/converter"
	"mediaconvertor/internal/stats"

	"github.com/charmbracelet/log"
)

func main() {
	statistics, params, fileBucket := converter.Initialize()

	converter.Run(
		statistics,
		params,
		fileBucket,
	)

	stats.CountPost(
		statistics,
		params.OutputImageDir,
		params.OutputVideoDir,
    params.OutputUndefiendDir,
	)

	log.Info("Structuring output folder")
	converter.StructureOutputLayout(params)

	log.Info("Cleaning")
	converter.CleaningUp(
		params.OutputImageDir,
		params.OutputVideoDir,
    params.OutputUndefiendDir,
	)

	stats.Process(statistics)
}
