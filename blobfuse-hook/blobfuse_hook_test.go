package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bpradipt/kata-hooks/blobfuse-hook/internal"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

// Set the logger
func init() {
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	internal.SetLogger(log)
}

// Create a test helper function which
// Read state.json file and populate s (spec.State)
// Read hookConfig.json file and populate hookConfig (internal.Config)
// Copy config.json file to bundle directory

func testHelper(stateJsonFile string, configJsonFile string, hookConfigJsonFile string) (spec.State, internal.Config, error) {

	var s spec.State
	var hookConfig internal.Config

	// Read the test-data/state.json file and copy it to jsonBytes
	jsonBytes, err := os.ReadFile(stateJsonFile)
	if err != nil {
		log.Errorf("unable to read state.json %s", err)
		return s, hookConfig, err
	}

	// Decode the json
	err = json.Unmarshal(jsonBytes, &s)
	if err != nil {
		log.Errorf("unable to decode state.json %s", err)
		return s, hookConfig, err
	}

	// Create bundle directory
	err = os.MkdirAll(s.Bundle, 0755)
	if err != nil {
		log.Errorf("unable to create bundle directory %s", err)
		return s, hookConfig, err
	}
	// Create rootfs directory
	err = os.MkdirAll(filepath.Join(s.Bundle, "rootfs"), 0755)
	if err != nil {
		log.Errorf("unable to create rootfs directory %s", err)
		return s, hookConfig, err
	}

	// Read the file test-data/config.json
	// Copy the file to bundle directory
	configJsonBytes, err := os.ReadFile(configJsonFile)
	if err != nil {
		log.Errorf("unable to read config.json %s", err)
		return s, hookConfig, err
	}

	// Write the file to bundle directory
	err = os.WriteFile(filepath.Join(s.Bundle, "config.json"), configJsonBytes, 0644)
	if err != nil {
		log.Errorf("unable to write config.json %s", err)
		return s, hookConfig, err
	}

	// Read the hookConfigJsonFile file and copy it to hookConfigJsonBytes
	hookConfigJsonBytes, err := os.ReadFile(hookConfigJsonFile)
	if err != nil {
		log.Errorf("unable to read hookConfig.json %s", err)
		return s, hookConfig, err
	}

	// Decode the json
	err = json.Unmarshal(hookConfigJsonBytes, &hookConfig)
	if err != nil {
		log.Errorf("unable to decode hookConfig.json %s", err)
		return s, hookConfig, err
	}

	return s, hookConfig, nil
}

func Test_doWork(t *testing.T) {
	type args struct {
		stateJsonFile      string
		configJsonFile     string
		hookConfigJsonFile string
		debug              bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Test valid stateJsonFile, configJsonFile and hookConfigJsonFile
		{
			name: "valid stateJsonFile, configJsonFile and hookConfigJsonFile",
			args: args{
				stateJsonFile:      "test-data/state.json",
				configJsonFile:     "test-data/config.json",
				hookConfigJsonFile: "test-data/hookconfig.json",
				debug:              false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// Execute testHelper
		s, hookConfig, err := testHelper(tt.args.stateJsonFile, tt.args.configJsonFile, tt.args.hookConfigJsonFile)
		if err != nil {
			// If testHelper returns an error, fail the test
			t.Errorf("testHelper() error = %v", err)
			return
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := doWork(s, hookConfig, tt.args.debug); (err != nil) != tt.wantErr {
				t.Errorf("doWork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
