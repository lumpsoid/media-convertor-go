package parameters

import (
	"flag"
	"mediaconvertor/internal/utils"
	"path"

	"github.com/charmbracelet/log"
)

type Orginize struct {
	InputDir            string
	OutputStructured    string
	Extensions          []string
	OverrideOutputDir   bool
	FromLogFile         string
	LogFilePath         string
	PopulateRecursively bool
	SaveFileName        bool
	NoProcessing        bool
}

func ParseForOrginizing() *Orginize {
	var params Orginize
	var extensions string

	// Define flags for each parameter
	flag.StringVar(
		&extensions,
		"extensions",
		"",
		"Comma-separated list of file extensions to process (e.g., 'jpg,png,mov,mp4')",
	)
	flag.BoolVar(
		&params.OverrideOutputDir,
		"overrideOutputDir",
		false,
		"Whether to override the output directory",
	)
	flag.BoolVar(
		&params.SaveFileName,
		"saveFileName",
		false,
		"Don't process file name into date pattern",
	)
	flag.BoolVar(
		&params.PopulateRecursively,
		"recursive",
		false,
		"Collect files recursively from input dir",
	)
	flag.BoolVar(
		&params.NoProcessing,
		"noProcess",
		false,
		"Don't process file name and structure. Creating folder with copies.",
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

	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("Provide path to the directory with media")
	}

	if len(args) > 1 {
		log.Fatal("Only one non-flag argument is supported - path to the directory with media")
	}

	params.Extensions = processExtensions(extensions)

	params.InputDir = utils.ExpandHomeDir(args[0])
	params.LogFilePath = utils.ExpandHomeDir(params.LogFilePath)
	params.FromLogFile = utils.ExpandHomeDir(params.FromLogFile)

	params.OutputStructured = path.Join(params.InputDir, "structured")

	return &params
}

func CheckForOrginizing(params *Orginize) {
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
	if len(params.Extensions) == 0 {
		log.Fatal("Please provide an extensions using the -extensions flag. Like: 'jpg,png,mov,mp4'")
	}
}
