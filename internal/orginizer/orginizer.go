package orginizer

import (
	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/processing"
	"mediaconvertor/internal/stats"
	"mediaconvertor/internal/utils"
)

func SetUpFiles(params *parameters.Orginize) (*stats.Stats, *filebucket.FileBucket) {
	if len(params.FromLogFile) > 0 {
		fileBucket := filebucket.FileBucketFromLogFile(params.FromLogFile)
		stats := stats.FromFileBucket(fileBucket)
		return &stats, fileBucket
	}

	fileBucket := filebucket.FileBucketFromExtensions(params.Extensions)

	if params.PopulateRecursively {
		filebucket.PopulateFileBuketRecursive(fileBucket, params.InputDir)
	} else {
		filebucket.PopulateFileBucket(fileBucket, params.InputDir)
	}

	stats := stats.FromFileBucket(fileBucket)

	utils.CheckAndClearDir(
		params.OutputStructured,
		params.OverrideOutputDir,
		"Output structured",
	)

	return &stats, fileBucket
}

func Initialize() (*stats.Stats, *parameters.Orginize, *filebucket.FileBucket) {
	p := parameters.ParseForOrginizing()
	parameters.CheckForOrginizing(p)

	utils.CheckProgramAvailability("exiftool")

	s, f := SetUpFiles(p)
	return s, p, f
}

func CopyFiles(
	fileBucket *filebucket.FileBucket,
	statistics *stats.Stats,
	outputStructured string,
	saveFileName bool,
) {
	processing.CopyFiles(fileBucket, statistics, outputStructured, saveFileName)
	return
}

func StructureOutputLayout(outputStructuredDir string, saveFileName bool) {
	processing.OrginizeWithExiftool(outputStructuredDir, saveFileName)
}
