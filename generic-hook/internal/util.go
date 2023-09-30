package internal

import (
	"os"
)

// Convert options []string to comma separated string
func ConvertOptionsToString(options []string) string {
	// Create a variable to hold the options
	var optionsString string

	// Loop through the options
	for _, option := range options {
		optionsString = optionsString + "," + option
	}

	// Remove the first comma
	optionsString = optionsString[1:]

	log.Printf("options string %s\n", optionsString)
	// Return the options string
	return optionsString
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isWritable(path string) bool {
	_, err := os.OpenFile(path, os.O_WRONLY, 0666)
	return err == nil
}

func isReadable(path string) bool {
	_, err := os.Open(path)
	return err == nil
}
