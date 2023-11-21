package internal

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Add a logger to the configuration
var log *logrus.Logger

// Create a struct to hold the configuration
type Config struct {

	// Add an activation flag to the configuration
	// This flag will be used to determine if the hook should be activated
	// or not

	// Activation flag needs to be container specific and not pod specific.
	// So best is to use container environment variable to activate it.
	// Check if the hookConfig.ActivationFlag* is present in containerConfig.Process.Env to activate the hook

	// If the ActivationFlag is not present in containerConfig.Process.Env, then the hook will not be activated
	// If the ActivationFlag is present in containerConfig.Process.Env, then the hook will be activated
	ActivationFlag string `json:"activation_flag,omitempty"`

	// Blobfuse program path
	ProgramPath string `json:"program_path"`

	// Host mountpoint
	HostMountPoint string `json:"host_mountpoint"`

	// Container mountpoint
	ContainerMountPoint string `json:"container_mountpoint"`
}

// Create a method to read the configuration file
func ReadConfig(configFile string) (Config, error) {
	// Create a variable to hold the configuration
	var config Config

	// Read the configuration file
	jsonData, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("unable to read configuration file %s\n", err)
		return config, err
	}
	// Unmarshal the configuration
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		log.Printf("unable to unmarshal configuration file %s\n", err)
		return config, err
	}

	// Return the configuration
	return config, nil
}

// Set the logger
func SetLogger(logger *logrus.Logger) {
	log = logger
}

// Method to check if ActivationFlag is present in a slice of strings
func IsActivationFlagPresent(env []string, activationFlag string) bool {
	log.Printf("Searching for activation flag %s\n", activationFlag)
	for _, val := range env {
		// env strings are of the form key=value
		// Match key with activationFlag
		if strings.Contains(val, activationFlag) {
			log.Printf("Activation flag %s is present\n", activationFlag)
			return true
		}
	}
	return false
}
