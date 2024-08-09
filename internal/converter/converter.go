package converter

import (
	"fmt"
	"mediaconvertor/internal/filebucket"
	"mediaconvertor/internal/parameters"
	"mediaconvertor/internal/processing"
	"mediaconvertor/internal/stats"
	"mediaconvertor/internal/utils"
	"os/exec"

	"github.com/charmbracelet/log"
)

// Check if specified program is available in the system
//
// Will call os.Fatal(1) on failure
func CheckProgramAvailability(program string) {
	// Run "which" command to check program availability
	cmd := exec.Command("which", program)
	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Fatal(fmt.Sprintf("Program %s is not present on the system", program))
	}
}

func SetUpFiles(params *parameters.Parameters) (*stats.Stats, *filebucket.FileBucket) {
	if len(params.FromLogFile) > 0 {
		fileBucket := filebucket.FileBucketFromLogFile(params.FromLogFile)
		stats := stats.FromFileBucket(fileBucket)
		return &stats, fileBucket
	}
	fileBucket := filebucket.FileBucketFromExtensions(params.Extensions)
	filebucket.PopulateFileBucket(params, fileBucket)
	stats := stats.FromFileBucket(fileBucket)
	return &stats, fileBucket
}

func Initialize() (*stats.Stats, *parameters.Parameters, *filebucket.FileBucket) {
	// Parse parameters
	p := parameters.Parse()
	parameters.Check(p)

	CheckProgramAvailability("exiftool")
	CheckProgramAvailability("ffmpeg")

	parameters.LoggingCheckedParams(p)

	s, f := SetUpFiles(p)
	return s, p, f
}

func Run(statistics *stats.Stats, params *parameters.Parameters, fileBucket *filebucket.FileBucket) {
	processing.Files(fileBucket, params, statistics)
	return
}

func StructureOutputLayout(
	imageExt string,
	videoExt string,
	outputImageDir string,
	outputVideoDir string,
	outputStructuredDir string,
) {
	log.Info("Structuring output folder")
	processing.PostWithExiftool(
		imageExt,
		videoExt,
		outputImageDir,
		outputVideoDir,
		outputStructuredDir,
	)
}

func CleaningUp(outputImageDir string, outputVideoDir string) {
  log.Info("Cleaning")
  utils.RemoveEmptyDir(outputImageDir)
  utils.RemoveEmptyDir(outputVideoDir)
}
