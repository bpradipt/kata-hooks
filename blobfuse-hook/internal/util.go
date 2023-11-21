package internal

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
