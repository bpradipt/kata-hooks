package internal

import (
	"os"
	"path/filepath"
	"syscall"

	sysmount "github.com/moby/sys/mount"
)

// Create device nodes using syscall.Mknod
func CreateDevices(rootfsPath string, hookConfig *Config) error {

	log.Printf("Creating devices %v\n", hookConfig.Devices)

	// Loop through the hookConfig.Devices
	for _, device := range hookConfig.Devices {
		// Create the device node
		mode := setDeviceMode(device.Type, *device.FileMode)
		deviceID := device.Major<<8 | device.Minor
		devicePath := filepath.Join(rootfsPath, device.Path)
		err := syscall.Mknod(devicePath, mode, int(deviceID))
		if err != nil {
			log.Printf("unable to create device node %s\n", err)
			return err
		}
		log.Printf("created device node %s\n", devicePath)

	}
	return nil
}

// Method to set the device mode
func setDeviceMode(deviceType string, fileMode os.FileMode) uint32 {
	// Set the device mode
	var mode uint32
	switch deviceType {
	case "c":
		mode = syscall.S_IFCHR
	case "b":
		mode = syscall.S_IFBLK
	default:
		mode = syscall.S_IFCHR
	}
	mode |= uint32(fileMode & os.ModePerm)

	log.Printf("device mode %d, access %d\n", mode, uint32(fileMode&os.ModePerm))
	return mode
}

// Method to mount the hookConfig mounts
func CreateMounts(rootfsPath string, hookConfig *Config) error {

	log.Printf("Creating mounts %v\n", hookConfig.Mounts)

	// Loop through the hookConfig.Mounts
	for _, mount := range hookConfig.Mounts {
		// Create the mount point
		mountPath := filepath.Join(rootfsPath, mount.Destination)
		err := os.MkdirAll(mountPath, 0755)
		if err != nil {
			log.Printf("creating mount point (%s) threw error (%s)\n", mountPath, err)
			return err
		}

		// Mount the mount point
		err = sysmount.Mount(mount.Source, mountPath, mount.Type, ConvertOptionsToString(mount.Options))
		if err != nil {
			log.Printf("mounting (%s) threw error (%s)\n", mountPath, err)
			return err
		}

		log.Printf("mounted %s\n", mountPath)

	}
	return nil
}

// Create method to create the directories
// The input is list of directories and the rootfs path where the directories should be created
func CreateDirs(rootfsPath string, hookConfig *Config) error {

	log.Printf("Creating directories %v\n", hookConfig.Dirs)

	// Loop through the list of directories
	for _, dir := range hookConfig.Dirs {
		// Create the directory
		dirPath := filepath.Join(rootfsPath, dir.Path)
		// if dir.Perm is empty then set it to 0666
		if dir.Perm == 0 {
			dir.Perm = 0666
		}

		if err := os.MkdirAll(dirPath, dir.Perm); err != nil {
			// Let's log and ignore
			log.Printf("creating directory (%s) failed with error (%s)", dirPath, err)
		}
		log.Printf("created directory %s\n", dirPath)
	}

	// Return nil
	return nil
}

// Create method to create the files
// The input is list of files and the rootfs path where the files should be created
func CreateFiles(rootfsPath string, hookConfig *Config) error {

	log.Printf("Creating files %v\n", hookConfig.Files)
	// Loop through the list of files
	for _, file := range hookConfig.Files {
		// Create the file
		filePath := filepath.Join(rootfsPath, file.Path)
		// if file.Perm is empty then set it to 0666
		if file.Perm == 0 {
			file.Perm = 0666
		}
		log.Printf("Creating file %s\n", filePath)
		if _, err := os.OpenFile(filePath, os.O_CREATE, file.Perm); err != nil {
			// Let's log and ignore
			log.Printf("failed to create file %s: %v", filePath, err)
		}
	}

	// Return nil
	return nil
}
