package main

import (
	"flag"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type Parameters struct {
	InputDir          string
	OutputVideoDir    string
	OutputImageDir    string
	Extensions        []string
	MinVideoDimension int
	NoClean           bool
	FromLogFile       string
	startTime         time.Time
}

func parseParameters() Parameters {
	var params Parameters
	var extensions string

	params.startTime = time.Now()
	// Define and parse command-line parameters
	// flag.StringVar(&params.OutputDir, "output_dir", "", "Output directory")
	flag.StringVar(&params.InputDir, "input_dir", "", "Input directory")
	flag.StringVar(&extensions, "ext", "", "File extension to process. In format: 'jpg,png,mov,mp4'")
	flag.StringVar(&params.FromLogFile, "from_log_file", "", "Process files from previous failed run.")
	flag.IntVar(&params.MinVideoDimension, "min_video_dim", 0, "Minimum video dimension convert to.")
	flag.BoolVar(&params.NoClean, "no_clean", false, "Clean output directory before processing.")
	// TODO add target extensions for video and photo
	// just copy them if they appear in stream
	// TODO add flag for forcefully process appeared files which extension same as target
	flag.Parse()

	if params.InputDir == "" {
		log.Fatal("Please provide an input dir using the -input_dir flag.")
	}
	if extensions == "" {
		log.Fatal("Please provide an extensions using the -ext flag. Like: 'jpg,png,mov,mp4'")
	}
	if params.MinVideoDimension == 0 {
		log.Fatal("Please provide a minimum video dimension using the -min_video_dim flag.")
	}

	params.Extensions = strings.Split(extensions, ",")
	newExtensions := []string{}
	for _, ext := range params.Extensions {
		newExtensions = append(newExtensions, strings.Trim(ext, "."))
	}
	params.Extensions = newExtensions

	params.OutputVideoDir = path.Join(params.InputDir, "mov")
	params.OutputImageDir = path.Join(params.InputDir, "img")

	return params
}
