package main

import (
	"mediaconvertor/internal/converter"
	"mediaconvertor/internal/stats"
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
	)

	converter.StructureOutputLayout(
		params.ImageTargetFormat,
		params.VideoTargetFormat,
		params.OutputImageDir,
		params.OutputVideoDir,
		params.OutputStructured,
	)

	converter.CleaningUp(
		params.OutputImageDir,
		params.OutputVideoDir,
	)

	stats.Process(statistics)
}
