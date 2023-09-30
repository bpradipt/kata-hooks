package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kata-hooks/generic-hook/internal"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// version is the version string of the hook. Set at build time.
	Version = "0.1"
	log     = logrus.New()
)

const (
// kataContainersPath = "/run/kata-containers"
)

func startOciHook(hookConfig *internal.Config, debug bool) error {
	//Hook receives container State in Stdin
	//https://github.com/opencontainers/runtime-spec/blob/master/config.md#posix-platform-hooks
	//https://github.com/opencontainers/runtime-spec/blob/master/runtime.md#state

	var s spec.State
	reader := bufio.NewReader(os.Stdin)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&s)
	if err != nil {
		return err
	}

	//log spec State to file
	log.Infof("spec.State is %v", s)

	bundlePath := s.Bundle
	containerPid := s.Pid
	containerState := s.Status

	log.Infof("container pid (%d): state (%s): bundle location (%s)\n", containerPid, containerState, bundlePath)

	configJsonPath := filepath.Join(bundlePath, "config.json")

	log.Infof("Config.json location: %s\n", configJsonPath)

	containerConfig, err := internal.ReadOciConfigJson(configJsonPath)
	if err != nil {
		log.Errorf("unable to read config.json %s", err)
		return err
	}

	if debug {
		log.Debugf("containerConfig contents: %v", containerConfig)
	}

	rootfsPath := filepath.Join(bundlePath, "rootfs")

	log.Printf("rootfsPath is %s\n", rootfsPath)

	// Get all the activation flags
	activationFlags := internal.GetActivationFlags(containerConfig.Process.Env, hookConfig.ActivationFlagAll,
		hookConfig.ActivationFlagFiles, hookConfig.ActivationFlagDirs, hookConfig.ActivationFlagMounts, hookConfig.ActivationFlagDevices)

	// Check if ActivationFlagMounts is present in activationFlags
	if activationFlags&8 == 8 {
		log.Infof("Activation flag %s is set to true. Creating mounts specified in hookConfig", hookConfig.ActivationFlagMounts)
		// Create the hookConfig mounts
		err = internal.CreateMounts(rootfsPath, hookConfig)
		if err != nil {
			log.Printf("unable to create mounts defined in hook config %s\n", err)
			return err
		}
	}

	// Check if ActivationFlagDevices is present in activationFlags
	if activationFlags&16 == 16 {
		log.Infof("Activation flag %s is set to true. Creating devices specified in hookConfig", hookConfig.ActivationFlagDevices)
		// Create the hookConfig devices
		err = internal.CreateDevices(rootfsPath, hookConfig)
		if err != nil {
			log.Printf("unable to create devices defined in hook config %s\n", err)
			return err
		}
	}

	if debug {
		log.Debugf("updated containerConfig contents: %v", containerConfig)
	}

	// Write the config.json file
	if err := internal.WriteOciConfigJson(configJsonPath, containerConfig); err != nil {
		log.Printf("unable to write config.json %s\n", err)
		return err
	}

	return nil
}

func main() {
	var hookConfigFile string
	var debug, version bool
	var logFile string

	// Create a cmd line parser based on "github.com/spf13/cobra" package
	rootCmd := &cobra.Command{
		Use:   "hook",
		Short: "OCI hook",
		Long:  "OCI hook to use with Kata Containers",
		Run: func(cmd *cobra.Command, args []string) {

			// if version flag is set, print the version and exit
			if version {
				fmt.Printf("OCI hook version %s\n", Version)
				os.Exit(0)
			}

			log.Out = os.Stdout

			// Check if log file is specified, otherwise create a temp file
			if logFile != "" {
				file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err == nil {
					log.Infof("Log file: %s\n", logFile)
					log.Out = file
				}
			} else {
				log.Info("No log file. Using temp file")

				dname, err := os.MkdirTemp("", "guesthooklog")
				if err != nil {
					log.Fatal(err)
				}

				fname := filepath.Join(dname, "guesthooklog.log")
				file, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err == nil {
					log.Infof("Log file: %s\n", fname)
					log.Out = file
				} else {
					log.Info("Failed to log to file, using default stderr\n")
				}
			}

			log.Infof("Started OCI hook version %s\n", Version)

			if debug {
				log.SetLevel(logrus.DebugLevel)
			}

			// set logger for internal package
			internal.SetLogger(log)

			// Parse hook config file
			hookConfig, err := internal.ReadConfig(hookConfigFile)
			if err != nil {
				log.Fatal(err)
			}

			log.Info("Starting OCI hook\n")
			if err := startOciHook(hookConfig, debug); err != nil {
				//Hook should not fail
				//log.Fatal(err)
				log.Info(err)
				return
			}
		},
	}

	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode (default is false)")
	rootCmd.Flags().BoolP("start", "s", true, "Start the OCI hook")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print the version")
	rootCmd.Flags().StringVarP(&hookConfigFile, "config", "c", "/usr/share/oci/hooks/hookconfig.json", "Path to the hook config file (default /usr/share/oci/hooks/hookconfig.json))")
	// Log file or create a temp file
	rootCmd.Flags().StringVarP(&logFile, "log", "l", "", "Path to the log file (default is temp file)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
