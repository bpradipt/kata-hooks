package internal

import (
	"encoding/json"
	"os"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func ReadOciConfigJson(configJsonPath string) (specs.Spec, error) {
	var containerConfig specs.Spec

	// Read the config.json file
	ociConfigJsonData, err := os.ReadFile(configJsonPath)
	if err != nil {
		log.Printf("unable to read oci config.json %s\n", err)
		return containerConfig, err
	}

	// Unmarshal the config.json file

	err = json.Unmarshal(ociConfigJsonData, &containerConfig)
	if err != nil {
		log.Printf("unable to parse oci config.json %s\n", err)
		return containerConfig, err
	}
	return containerConfig, nil
}

// Write the config.json file
func WriteOciConfigJson(configJsonPath string, containerConfig specs.Spec) error {
	// Marshal the config.json file
	ociConfigJsonData, err := json.Marshal(containerConfig)
	if err != nil {
		log.Printf("unable to marshal oci config.json %s\n", err)
		return err
	}

	// Write the config.json file
	err = os.WriteFile(configJsonPath, ociConfigJsonData, 0644)
	if err != nil {
		log.Printf("unable to write oci config.json %s\n", err)
		return err
	}
	log.Printf("oci config.json written to %s\n", configJsonPath)
	return nil
}
