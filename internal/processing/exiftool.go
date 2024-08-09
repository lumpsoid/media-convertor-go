package processing

import (
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

func PostWithExiftool(imageExt string, videoExt string, outputImageDir string, outputVideoDir string)  {
	// -r -- recursive
	RunExiftool("-ext", strings.ToUpper(imageExt), "-FileName<${FileModifyDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", outputImageDir)
	RunExiftool("-ext", strings.ToUpper(imageExt), "-FileName<${CreateDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", outputImageDir)
	RunExiftool("-ext", strings.ToUpper(imageExt), "-FileName<${DateCreated}_${ImageSize}%-c.${FileTypeExtension}", "-d", "IMG_%Y-%m-%d_%H-%M-%S", outputImageDir)

	RunExiftool("-ext", strings.ToUpper(videoExt), "-FileName<${FileModifyDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "MOV_%Y-%m-%d_%H-%M-%S", outputVideoDir)
	RunExiftool("-ext", strings.ToUpper(videoExt), "-FileName<${CreateDate}_${ImageSize}%-c.${FileTypeExtension}", "-d", "MOV_%Y-%m-%d_%H-%M-%S", outputVideoDir)
  
}
