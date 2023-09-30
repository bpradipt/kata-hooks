package internal

import (
	"encoding/json"
	"io/fs"
	"os"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

// Add a logger to the configuration
var log *logrus.Logger

// Create methods to handle the configuration for OCI hook
// The OCI hook configuration consists of
// list of devices
// list of directories
// list of files
// list of mounts

// Create a struct to hold the configuration
type Config struct {

	// Add an activation flag to the configuration
	// This flag will be used to determine if the hook should be activated
	// or not

	// Activation flag needs to be container specific and not pod specific.
	// So best is to use container environment variable to activate it.
	// Check if the hookConfig.ActivationFlag* is present in containerConfig.Process.Env to activate the hook

	ActivationFlagAll     string `json:"activation_flag_all"`
	ActivationFlagFiles   string `json:"activation_flag_files"`
	ActivationFlagDirs    string `json:"activation_flag_dirs"`
	ActivationFlagMounts  string `json:"activation_flag_mounts"`
	ActivationFlagDevices string `json:"activation_flag_devices"`

	// Example devices
	/*
			   [

			        {
		                "path": "/dev/fuse",
		                "type": "c",
		                "major": 10,
		                "minor": 229,
		                "fileMode": 438,
		                "uid": 0,
		                "gid": 0
		            },
		            {
		                "path": "/dev/sda",
		                "type": "b",
		                "major": 8,
		                "minor": 0,
		                "fileMode": 432,
		                "uid": 0,
		                "gid": 0
		            }
				]

	*/
	Devices []specs.LinuxDevice `json:"devices"`
	Dirs    []Dir               `json:"dirs"`
	Files   []File              `json:"files"`
	// Example mount
	/*
		{
			"destination": "/etc/resolv.conf",
			"type": "bind",
			"source": "/run/kata-containers/shared/containers/49a54568371eb19597d8ba20394d29f849d420df22942e438070cffa2a242fcf-46e40207515f27d3-resolv.conf",
			"options": [
			  "ro",
			  "bind",
			  "nodev",
			  "nosuid",
			  "noexec"
			]
		   },
	*/
	Mounts []specs.Mount `json:"mounts"`
}

// Create a struct to hold the directory configuration
type Dir struct {
	Path string `json:"path"`
	// Add permissions to the directory configuration
	// This will be used to set the permissions on the directory
	// Default should be 0666 if not specified
	Perm fs.FileMode `json:"perm"`
}

// Create a struct to hold the file configuration
type File struct {
	Path string `json:"path"`
	// Add permissions to the file configuration
	// This will be used to set the permissions on the file
	// Default should be 0666 if not specified
	Perm fs.FileMode `json:"perm"`
}

// Create a method to read the configuration file
func ReadConfig(configFile string) (*Config, error) {
	// Read the configuration file
	jsonData, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("unable to read configuration file %s\n", err)
		return nil, err
	}

	// Create a variable to hold the configuration
	var config Config

	// Unmarshal the configuration
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		log.Printf("unable to unmarshal configuration file %s\n", err)
		return nil, err
	}

	// Return the configuration
	return &config, nil
}

// Set the logger
func SetLogger(logger *logrus.Logger) {
	log = logger
}

// Method to check which all activation flags are present in the env string slice
// Return a bit mask with the flags that are present set to 1
func GetActivationFlags(env []string, activationFlagAll string, activationFlagFiles string, activationFlagDirs string,
	activationFlagMounts string, activationFlagDevices string) int {

	log.Printf("Searching for activation flags %s, %s, %s, %s, %s\n", activationFlagAll, activationFlagFiles,
		activationFlagDirs, activationFlagMounts, activationFlagDevices)
	var activationFlags int = 0
	if IsActivationFlagPresent(env, activationFlagAll) {
		activationFlags = activationFlags | 1
	}
	if IsActivationFlagPresent(env, activationFlagFiles) {
		activationFlags = activationFlags | 2
	}
	if IsActivationFlagPresent(env, activationFlagDirs) {
		activationFlags = activationFlags | 4
	}
	if IsActivationFlagPresent(env, activationFlagMounts) {
		activationFlags = activationFlags | 8
	}
	if IsActivationFlagPresent(env, activationFlagDevices) {
		activationFlags = activationFlags | 16
	}

	// Print which activation flags are present
	if activationFlags&1 == 1 {
		log.Printf("Activation flag %s is present\n", activationFlagAll)
	}
	if activationFlags&2 == 2 {
		log.Printf("Activation flag %s is present\n", activationFlagFiles)
	}
	if activationFlags&4 == 4 {
		log.Printf("Activation flag %s is present\n", activationFlagDirs)
	}
	if activationFlags&8 == 8 {
		log.Printf("Activation flag %s is present\n", activationFlagMounts)
	}
	if activationFlags&16 == 16 {
		log.Printf("Activation flag %s is present\n", activationFlagDevices)
	}

	return activationFlags
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
