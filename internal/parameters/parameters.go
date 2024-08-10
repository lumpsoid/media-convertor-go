package parameters

import (
	"flag"
	"mediaconvertor/internal/utils"
	"path"
	"strings"

	"github.com/charmbracelet/log"
)

type Parameters struct {
	InputDir           string
	OutputVideoDir     string
	OutputImageDir     string
	OutputUndefiendDir string
	OutputStructured   string
	Extensions         []string
	VideoMinDimension  int
	VideoTargetFps     int
	VideoTargetFormat  string
	ImageTargetFormat  string
	ImageTargetQuality int
	OverrideOutputDir  bool
	FromLogFile        string
	LogFilePath        string
}

func Parse() *Parameters {
	var params Parameters
	var extensions string

	// Define flags for each parameter
	flag.StringVar(
		&params.InputDir,
		"inputDir",
		"",
		"Directory containing input files. Must be flat",
	)
	flag.StringVar(
		&extensions,
		"extensions",
		"",
		"Comma-separated list of file extensions to process (e.g., 'jpg,png,mov,mp4')",
	)
	flag.IntVar(
		&params.VideoMinDimension,
		"videoMinDimension",
		0,
		"Minimum dimension of the video to process",
	)
	flag.IntVar(
		&params.VideoTargetFps,
		"videoTargetFps",
		0,
		"Target frames per second for video processing",
	)
	flag.StringVar(
		&params.VideoTargetFormat,
		"videoTargetFormat",
		"",
		"Target format for processed videos (e.g., 'mp4')",
	)
	flag.StringVar(
		&params.ImageTargetFormat,
		"imageTargetFormat",
		"",
		"Target format for processed images (e.g., 'jpg', 'png')",
	)
	flag.IntVar(
		&params.ImageTargetQuality,
		"imageTargetQuality",
		0,
		"Quality of the processed images (e.g., 80, 57)",
	)
	flag.BoolVar(
		&params.OverrideOutputDir,
		"overrideOutputDir",
		false,
		"Whether to override the output directory",
	)
	flag.StringVar(
		&params.FromLogFile,
		"fromLogFile",
		"",
		"Path to the log file with file paths from a previous run",
	)
	flag.StringVar(
		&params.LogFilePath,
		"logFilePath",
		"",
		"Path to the log file for recording processing details",
	)

	flag.Parse()

	params.Extensions = processExtensions(extensions)

	params.InputDir = utils.ExpandHomeDir(params.InputDir)
	params.LogFilePath = utils.ExpandHomeDir(params.LogFilePath)
	params.FromLogFile = utils.ExpandHomeDir(params.FromLogFile)

	params.OutputVideoDir = path.Join(params.InputDir, "mov")
	params.OutputImageDir = path.Join(params.InputDir, "img")
	params.OutputUndefiendDir = path.Join(params.InputDir, "undef")
	params.OutputStructured = path.Join(params.InputDir, "structured")

	return &params
}

func processExtensions(extensionsString string) []string {
	extensions := strings.Split(extensionsString, ",")
	newExtensions := []string{}
	for _, ext := range extensions {
		newExtensions = append(newExtensions, strings.Trim(ext, "."))
	}
	return newExtensions
}

func Check(params *Parameters) {
	if params.InputDir == "" {
		log.Fatal("Please provide an input dir using the -inputDir flag")
	}
	if !utils.IsPathExist(params.InputDir) {
		log.Fatal("Provided dir with '-inputDir' doesn't exist")
	}
	if params.LogFilePath == "" {
		log.Fatal("Please provide a path to the logfile using the -logFilePath flag")
	}
	if params.FromLogFile != "" && params.LogFilePath != "" && params.LogFilePath == params.FromLogFile {
		log.Fatal("Please use different logfiles with -from_logfile and -logFilePath")
	}
	if utils.IsPathExist(params.LogFilePath) {
		yes, err := utils.ConfirmPrompt("Provided logfile path with '-logFilePath' already exists. Want to override?")
		if err != nil {
			log.Fatal("Error occured while read prompt respose")
		}
		if !yes {
			log.Fatal("Please choose different path for the '-logFilePath' flag")
		}
	}
	if params.VideoMinDimension <= 0 {
		log.Fatal("Please provide a minimum video dimension using the -min_video_dim flag")
	}
	if params.VideoTargetFps <= 0 {
		log.Fatal("Please provide a video target fps using the -videoTargetFps flag")
	}
	if params.VideoTargetFormat == "" {
		log.Fatal("Please provide a video target format using the -videoTargetFormat flag")
	}
	if params.ImageTargetFormat == "" {
		log.Fatal("Please provide a image target format using the -imageTargetFormat flag")
	}
	if params.ImageTargetQuality <= 0 {
		log.Fatal("Please provide a image target quality using the -imageTargetQuality flag")
	}
	if len(params.Extensions) == 0 {
		log.Fatal("Please provide an extensions using the -extensions flag. Like: 'jpg,png,mov,mp4'")
	}

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
}

func LoggingCheckedParams(params *Parameters) {
	log.Infof("Image parameters: target format = '%s', target quality = '%d'", params.ImageTargetFormat, params.ImageTargetQuality)
	log.Infof("Video parameters: target format = '%s', target fps = '%d', min dimension = '%d'", params.VideoTargetFormat, params.VideoTargetFps, params.VideoMinDimension)
}
