/*
Download a mini linux filesystem like https://alpinelinux.org/downloads/
wget https://dl-cdn.alpinelinux.org/alpine/v3.18/releases/x86_64/alpine-minirootfs-3.18.3-x86_64.tar.gz

If cgroups are used then you should execute this as root. In any other case you can execute it a regular user.
*/

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"	
	"syscall"
	
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Expected at least one argument")
	}

	switch os.Args[1] {
	case "run":
		run(os.Args[2:]...)
	case "child":
		child(os.Args[2:]...)
	default:
		log.Fatal("Unknown command. Use run <command_name>, like `run /bin/sh` or `run /bin/echo hello` ")
	}
}

func run(command ...string) {
	log.Println("Executing", command, "from run")
	log.Printf("Running %v as user %d in process %d\n", os.Args[2:], os.Geteuid(), os.Getpid())


	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, command...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Add user namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS, // Hide the mounts of container from the host
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				// Maps to a regular user
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	// Run child using namespaces
	must(cmd.Run())
}

func child(command ...string) {
	log.Println("Executing", command, "from child")
	log.Printf("Running %v as user %d in process %d\n", os.Args[2:], os.Geteuid(), os.Getpid())

	// Create cgroup	
	cgPids()

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("./alpine_fs"))
	// Change directory after chroot
	must(os.Chdir("/"))

	// Mount /proc inside container so that `ps` command works
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	defer syscall.Unmount("proc", 0)

	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to run command: %v", err)
	}
}


// Working with cgroupv2
// Requires root permissions to write on cgroups dir
// Restricts the number of processes to 10
func cgPids() {
	cgroups := "/sys/fs/cgroup/"
	containerDir := filepath.Join(cgroups, "go_container/")	
	os.Mkdir(containerDir, 0755)
	must(os.WriteFile(filepath.Join(containerDir, "pids.max"), []byte("10"), 0700))	

	// Adding the process to the cgroup
	must(os.WriteFile(filepath.Join(containerDir, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

	// proper clean up is needed
}


func must(err error) {
	if err != nil {
		panic(err)
	}
}
