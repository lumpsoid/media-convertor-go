package processing

import (
	"fmt"
	"mediaconvertor/internal/parameters"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

func RunExiftool(args ...string) {
	cmd := exec.Command("exiftool", args...)
	err := cmd.Run()
	if err != nil {
		log.Errorf("run problem with exiftool: %v", err)
	}
}

func renameImages(imageExt string, outputImageDir string) {
	dateFormat := "IMG_%Y-%m-%d_%H-%M-%S"
	imageExtUpper := strings.ToUpper(imageExt)
	dateTags := []string{
		"FileModifyDate",
		"CreateDate",
		"DateCreated",
	}
	for _, dateTag := range dateTags {
		RunExiftool(
			"-ext", imageExtUpper,
			fmt.Sprint("-FileName<${", dateTag, "}_${ImageSize}%-c.${FileTypeExtension}"),
			"-d", dateFormat,
			outputImageDir,
		)

	}
}

func structureImages(imageExt string, outputImageDir string, outputStructured string) {
	directoryStructure := "%Y/%m"
	imageExtUpper := strings.ToUpper(imageExt)
	dateTags := []string{
		"DateCreated",
		"CreateDate",
		"FileModifyDate",
	}
	for _, dateTag := range dateTags {
		RunExiftool(
			"-ext", imageExtUpper,
			fmt.Sprint("-Directory<", outputStructured, "/${", dateTag, "}"),
			"-d", directoryStructure,
			outputImageDir,
		)
	}
}

func processImages(imageExt string, outputImageDir string, outputStructured string) {
	renameImages(imageExt, outputImageDir)
	structureImages(imageExt, outputImageDir, outputStructured)
}

func renameVideos(videoExt string, outputVideoDir string) {
	dateFormat := "MOV_%Y-%m-%d_%H-%M-%S"
	videoExtUpper := strings.ToUpper(videoExt)
	dateTags := []string{
		"FileModifyDate",
		"CreateDate",
	}
	for _, dateTag := range dateTags {
		RunExiftool(
			"-ext", videoExtUpper,
			fmt.Sprint("-FileName<${", dateTag, "}_${ImageSize}%-c.${FileTypeExtension}"),
			"-d", dateFormat,
			outputVideoDir,
		)
	}
}

func structureVideos(videoExt string, outputVideoDir string, outputStructured string) {
	directoryStructure := "%Y/%m"
	videoExtUpper := strings.ToUpper(videoExt)
	dateTags := []string{
		"FileModifyDate",
		"CreateDate",
	}
	for _, dateTag := range dateTags {
		RunExiftool(
			"-ext", videoExtUpper,
			fmt.Sprint("-Directory<", outputStructured, "/${", dateTag, "}"),
			"-d", directoryStructure,
			outputVideoDir,
		)
	}
}

func processVideos(videoExt string, outputVideoDir string, outputStructuredDir string) {
	renameVideos(videoExt, outputVideoDir)
	structureVideos(videoExt, outputVideoDir, outputStructuredDir)
}

func renameUndefiend(outputUndefiendDir string) {
	dateFormat := "%Y-%m-%d_%H-%M-%S"
	dateTags := []string{
		"FileModifyDate",
		"CreateDate",
		"DateCreated",
	}
	for _, dateTag := range dateTags {
		RunExiftool(
			fmt.Sprint("-FileName<${", dateTag, "}_${ImageSize}%-c.${FileTypeExtension}"),
			"-d", dateFormat,
			outputUndefiendDir,
		)
	}
}

func structureUndefiend(outputUndefiendDir string, outputStructured string) {
	directoryStructure := "%Y/%m"
	dateTags := []string{
		"DateCreated",
		"CreateDate",
		"FileModifyDate",
	}
	for _, dateTag := range dateTags {
		RunExiftool(
			fmt.Sprint("-Directory<", outputStructured, "/${", dateTag, "}"),
			"-d", directoryStructure,
			outputUndefiendDir,
		)
	}
}

func processUndefiend(outputUndefiendDir string, outputStructuredDir string) {
  renameUndefiend(outputUndefiendDir)
  structureUndefiend(outputUndefiendDir, outputStructuredDir)
}

func PostWithExiftool(params *parameters.Parameters) {

	processImages(params.ImageTargetFormat, params.OutputImageDir, params.OutputStructured)
	processVideos(params.VideoTargetFormat, params.OutputVideoDir, params.OutputStructured)
  processUndefiend(params.OutputUndefiendDir, params.OutputStructured)
}

func OrginizeWithExiftool(
	outputStructuredDir string,
	saveFileName bool,
) {
	// renaming
	if !saveFileName {
    log.Info("Renaming files")
		dateFormat := "IMG_%Y-%m-%d_%H-%M-%S"
		dateTags := []string{
			"FileModifyDate",
			"CreateDate",
			"DateCreated",
		}
		for _, dateTag := range dateTags {
			RunExiftool(
				fmt.Sprint("-FileName<${", dateTag, "}_${ImageSize}%-c.${FileTypeExtension}"),
				"-d", dateFormat,
				outputStructuredDir,
			)
		}
	}
	// structuring
	//
	// probably, because there is not -r recursive tag
	// it should be okay
  log.Info("Structuring files")
	directoryStructure := "%Y/%m"
	dateTagsForStructuring := []string{
		"DateCreated",
		"CreateDate",
		"FileModifyDate",
	}
	for _, dateTag := range dateTagsForStructuring {
		RunExiftool(
			fmt.Sprint("-Directory<", outputStructuredDir, "/${", dateTag, "}"),
			"-d", directoryStructure,
			outputStructuredDir,
		)
	}
}
