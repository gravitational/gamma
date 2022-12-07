package utils

import (
	"fmt"
	"os"
	"path"
)

func NormalizeDirectories(workingDirectory, outputDirectory string) (string, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("could not get current working directory: %v", err)
	}

	if workingDirectory == "" {
		workingDirectory = wd
	} else {
		if !path.IsAbs(workingDirectory) {
			workingDirectory = path.Join(wd, workingDirectory)
		}
	}

	if !path.IsAbs(outputDirectory) {
		outputDirectory = path.Join(workingDirectory, outputDirectory)
	}

	return workingDirectory, outputDirectory, nil
}
