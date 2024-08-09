package parameters

import (
	"flag"
	"mediaconvertor/internal/utils"
	"path"
	"strings"

	"github.com/charmbracelet/log"
)

type Parameters struct {
	InputDir          string
	OutputVideoDir    string
	OutputImageDir    string
	Extensions        []string
	MinVideoDimension int
  OverrideOutputDir bool
	FromLogFile       string
	LogFilePath       string
}

func Parse() *Parameters {
	var params Parameters
	var extensions string

	// Define and parse command-line parameters
	// flag.StringVar(&params.OutputDir, "output_dir", "", "Output directory")
	flag.StringVar(&params.InputDir, "input_dir", "", "Input directory")
	flag.StringVar(&extensions, "ext", "", "File extension to process. In format: 'jpg,png,mov,mp4'")
	flag.StringVar(&params.LogFilePath, "logfile", "", "Logfile for media file paths in case of a error")
	flag.StringVar(&params.FromLogFile, "from_logfile", "", "Process files from previous failed run.")
	flag.IntVar(&params.MinVideoDimension, "min_video_dim", 0, "Minimum video dimension convert to.")
	flag.BoolVar(&params.OverrideOutputDir, "override_output_dir", false, "Override output dir ('-input_dir' + 'img' or 'mov').")
	// TODO add target extensions for video and photo
	// just copy them if they appear in stream
	// TODO add flag for forcefully process appeared files which extension same as target
	flag.Parse()

	params.Extensions = processExtensions(extensions)

	params.OutputVideoDir = path.Join(params.InputDir, "mov")
	params.OutputImageDir = path.Join(params.InputDir, "img")

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
		log.Fatal("Please provide an input dir using the -input_dir flag.")
	}
	if params.LogFilePath == "" {
		log.Fatal("Please provide a path to the logfile using the -logfile flag. In case of an error in process of converting, path to the media file would be recorded in it, then you will be able to try process them again after first run using -from_logfile flag")
	}
  if params.FromLogFile != "" && params.LogFilePath != "" && params.LogFilePath == params.FromLogFile {
		log.Fatal("Please use different logfiles with -from_logfile and -logfile")
  }
	if utils.IsPathExist(params.LogFilePath) {
    yes, err := utils.ConfirmPrompt("Provided logfile path with '-logfile' already exists. Want to override?")
    if err != nil {
      log.Fatal("Error occured while read prompt respose")
    }
    if !yes {
      log.Fatal("Please choose different path for the '-logfile' flag")
    }
	}
	if !utils.IsPathExist(params.InputDir) {
		log.Fatal("Provided dir with '-input_dir' doesn't exist")
	}
	if params.MinVideoDimension == 0 {
		log.Fatal("Please provide a minimum video dimension using the -min_video_dim flag.")
	}

	if len(params.Extensions) == 0 {
		log.Fatal("Please provide an extensions using the -ext flag. Like: 'jpg,png,mov,mp4'")
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
}
