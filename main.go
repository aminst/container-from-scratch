package main

import (
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	switch os.Args[1] {
	case "run":
		runSelf("/proc/self/exe", append([]string{"rerun"}, os.Args[2:]...))
	case "rerun":
		runCommand(os.Args[2], os.Args[3:])
	default:
		panic("Use the following structure: run <command> <arguments>")
	}
}

func runSelf(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
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
		Unshareflags: syscall.CLONE_NEWNS,
	}
	cmd.Run()
}

func generateRandomHostname(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func runCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	must(syscall.Mount("newroot", "newroot", "", syscall.MS_BIND|syscall.MS_REC, ""))
	must(os.MkdirAll("newroot/oldfs", 0700))
	must(syscall.PivotRoot("newroot", "newroot/oldfs"))
	must(os.Chdir("/"))
	must(syscall.Unmount("/oldfs", syscall.MNT_DETACH))
	must(os.RemoveAll("/oldfs"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(syscall.Sethostname([]byte(generateRandomHostname(10))))
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
