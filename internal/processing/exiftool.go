package processing

import (
	"os/exec"

	"github.com/charmbracelet/log"
)


func RunExiftool(msg string, args ...string) {
	log.Infof(msg)
	cmd := exec.Command("exiftool", args...)
	err := cmd.Run()
	if err != nil {
		log.Errorf("run problem with exiftool: %v", err)
	}
}

func PostWithExiftool(outputImageDir string, outputVideoDir string)  {
	// -r -- recursive
	RunExiftool("Run exiftool for images", "-ext", "AVIF", "-FileName<${FileModifyDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", outputImageDir)
	RunExiftool("Run exiftool for images, trying CreateDate", "-ext", "AVIF", "-FileName<${CreateDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", outputImageDir)
	RunExiftool("Run exiftool for images, trying DateCreated", "-ext", "AVIF", "-FileName<${DateCreated}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", outputImageDir)

	RunExiftool("Run exiftool for videos", "-ext", "MP4", "-FileName<${FileModifyDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "MOV_%Y-%m-%d_%H-%M-%S", outputVideoDir)
	RunExiftool("Run exiftool for videos, trying CreateDate", "-ext", "MP4", "-FileName<${CreateDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "MOV_%Y-%m-%d_%H-%M-%S", outputVideoDir)
  
}
