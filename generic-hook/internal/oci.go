package internal

import (
	"encoding/json"
	"os"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func ReadOciConfigJson(configJsonPath string) (*specs.Spec, error) {
	// Read the config.json file
	ociConfigJsonData, err := os.ReadFile(configJsonPath)
	if err != nil {
		log.Printf("unable to read oci config.json %s\n", err)
		return nil, err
	}

	// Unmarshal the config.json file
	var containerConfig specs.Spec
	err = json.Unmarshal(ociConfigJsonData, &containerConfig)
	if err != nil {
		log.Printf("unable to parse oci config.json %s\n", err)
		return nil, err
	}
	return &containerConfig, nil
}

// Write the config.json file
func WriteOciConfigJson(configJsonPath string, containerConfig *specs.Spec) error {
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

// Method to add hookConfig mounts to the containerConfig mounts
func AddMountsToOciSpec(containerConfig *specs.Spec, hookConfig *Config) error {
	// Add the hookConfig mounts to the containerConfig mounts
	containerConfig.Mounts = append(containerConfig.Mounts, hookConfig.Mounts...)

	log.Printf("containerConfig.Mounts: %v\n", containerConfig.Mounts)
	return nil
}

// Method to add hookConfig devices to the containerConfig devices
func AddDevicesToOciSpec(containerConfig *specs.Spec, hookConfig *Config) error {
	// Add the hookConfig devices to the containerConfig devices
	containerConfig.Linux.Devices = append(containerConfig.Linux.Devices, hookConfig.Devices...)

	log.Printf("containerConfig.Linux.Devices: %v\n", containerConfig.Linux.Devices)
	return nil
}

// Method to whitelist the hookConfig devices
// Only works for cgroupv1
func AddDeviceWhitelistToOciSpec(containerConfig *specs.Spec, hookConfig *Config) error {

	/* "resources": {
		 "devices": [
	                {
	                    "allow": false,
	                    "access": "rwm"
	                },
	                {
	                    "allow": true,
	                    "type": "c",
	                    "major": 10,
	                    "minor": 229,
	                    "access": "rw"
	                },
	                {
	                    "allow": true,
	                    "type": "b",
	                    "major": 8,
	                    "minor": 0,
	                    "access": "r"
	                }
	            ]
			}
	*/

	// Loop through the hookConfig.Devices
	for _, device := range hookConfig.Devices {

		// Populate the deviceCgroup struct members from the device members
		deviceCgroup := specs.LinuxDeviceCgroup{
			Allow:  true,
			Access: "rwm",
			Type:   device.Type,
			Major:  &device.Major,
			Minor:  &device.Minor,
		}

		// Append the deviceCgroup to the containerConfig.Linux.Resources.Devices
		containerConfig.Linux.Resources.Devices = append(containerConfig.Linux.Resources.Devices, deviceCgroup)

	}

	log.Printf("containerConfig.Linux.Resources.Devices: %v\n", containerConfig.Linux.Resources.Devices)
	return nil
}
