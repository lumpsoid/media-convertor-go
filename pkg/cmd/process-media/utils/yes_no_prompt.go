package utils

import "github.com/charmbracelet/huh"

func confirmPrompt(title string) (bool, error) {
	var confirm bool
	err := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&confirm).Run()
	if err != nil {
		return false, err
	}
	return confirm, nil
}
