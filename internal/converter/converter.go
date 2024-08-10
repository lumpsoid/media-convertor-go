package converter

import (
	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/processing"
	"mediaconvertor/internal/stats"
	"mediaconvertor/internal/utils"
)

func SetUpFiles(params *parameters.Parameters) (*stats.Stats, *filebucket.FileBucket) {
	if len(params.FromLogFile) > 0 {
		fileBucket := filebucket.FileBucketFromLogFile(params.FromLogFile)
		stats := stats.FromFileBucket(fileBucket)
		return &stats, fileBucket
	}
	fileBucket := filebucket.FileBucketFromExtensions(params.Extensions)
	filebucket.PopulateFileBucket(fileBucket, params.InputDir)
	stats := stats.FromFileBucket(fileBucket)

	utils.CheckAndClearDir(
		params.OutputImageDir,
		params.OverrideOutputDir,
		"Output image",
	)
	utils.CheckAndClearDir(
		params.OutputVideoDir,
		params.OverrideOutputDir,
		"Output video",
	)
	utils.CheckAndClearDir(
		params.OutputUndefiendDir,
		params.OverrideOutputDir,
		"Output undefiend",
	)
	utils.CheckAndClearDir(
		params.OutputStructured,
		params.OverrideOutputDir,
		"Output structured",
	)

	return &stats, fileBucket
}

func Initialize() (*stats.Stats, *parameters.Parameters, *filebucket.FileBucket) {
	// Parse parameters
	p := parameters.Parse()
	parameters.Check(p)

	utils.CheckProgramAvailability("exiftool")
	utils.CheckProgramAvailability("ffmpeg")

	parameters.LoggingCheckedParams(p)

	s, f := SetUpFiles(p)
	return s, p, f
}

func Run(statistics *stats.Stats, params *parameters.Parameters, fileBucket *filebucket.FileBucket) {
	processing.Files(fileBucket, params, statistics)
	return
}

func StructureOutputLayout(params *parameters.Parameters) {
	processing.PostWithExiftool(params)
}

func CleaningUp(
	outputImageDir string,
	outputVideoDir string,
	outputUndefiendDir string,
) {
	utils.RemoveEmptyDir(outputImageDir)
	utils.RemoveEmptyDir(outputVideoDir)
	utils.RemoveEmptyDir(outputUndefiendDir)
}
