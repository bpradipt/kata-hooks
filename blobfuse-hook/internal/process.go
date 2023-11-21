package internal

import (
	"os"
	"os/exec"
	"strings"

	sysmount "github.com/moby/sys/mount"
)

// Execute process using syscall.Exec
// The blobfuse program path will be in hookConfig.ProgramPath
// The host mount point will be in hookConfig.HostMountPoint
// The container mount point will be in hookConfig.ContainerMountPoint
// Also use the environment variables from the containerConfig.Process.Env to execute the process

func ExecuteBlobFuseProcess(env []string, hookConfig Config) error {
	// Create the host mount point directory path
	err := os.MkdirAll(hookConfig.HostMountPoint, 0755)
	if err != nil {
		log.Printf("unable to create host mount point directory %s\n", err)
		return err
	}

	log.Printf("Executing program %s\n", hookConfig.ProgramPath)

	// Build the arguments for the process
	// The arguments will be the host mount point and other required
	arguments := []string{
		"mount",
		hookConfig.HostMountPoint,
		"--config-file=/etc/blobfuseconfig.yaml"}

	// Create a new command with the program path and arguments
	cmd := exec.Command(hookConfig.ProgramPath, arguments...)

	// Set the environment variables for the command
	cmd.Env = env

	// Run the command
	err = cmd.Run()
	if err != nil {
		log.Printf("unable to execute process %s\n", err)
		return err
	}

	return nil
}

// Bind mount src to dst
// The src will be the host mount point and dst will be the container mount point

func BindMount(srcMountPoint string, dstMountPoint string) error {

	log.Printf("Bind mounting host mount point %s to container mount point %s\n",
		srcMountPoint, dstMountPoint)

	// Create the dst mount point directory path
	err := os.MkdirAll(dstMountPoint, 0755)
	if err != nil {
		log.Printf("create container mount point directory returned err: %s\n", err)
		return err
	}

	// Bind mount the host mount point to container mount point
	err = sysmount.Mount(srcMountPoint, dstMountPoint, "none", "bind,rw")
	if err != nil {
		log.Printf("bind mount srcMountPoint (%s) dstMountPoint (%s) returned err: %s\n", srcMountPoint, dstMountPoint, err)
		return err
	}

	return nil
}

// Get CONTAINER_MOUNT_POINT value from containerConfig.Process.Env

func GetContainerMountPoint(env []string) string {

	for _, envVar := range env {
		if envVar == "CONTAINER_MOUNT_POINT" {
			// Split the envVar on "="
			// The second part will be the value of CONTAINER_MOUNT_POINT
			containerMountPoint := strings.Split(envVar, "=")[1]
			return containerMountPoint
		}
	}
	return ""
}
