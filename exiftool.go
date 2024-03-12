package main

import (
	"os/exec"

	"github.com/charmbracelet/log"
)

func runExiftool(msg string, args ...string) {
	log.Infof(msg)
	cmd := exec.Command("exiftool", args...)
	err := cmd.Run()
	if err != nil {
		log.Errorf("run problem with exiftool: %v", err)
	}
}
