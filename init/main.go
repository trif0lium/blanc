package main

import (
	"os"
	"os/exec"
	"syscall"
)

type ImageConfig struct {
	Entrypoint string
	Cmd        string
	Env        []string
	WorkingDir string
	User       string
}

func main() {
	imageConfig := ImageConfig{
		Entrypoint: "/whoami",
		Cmd:        "",
		Env: []string{
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"WHOAMI_NAME=/blanc/init.json",
		},
		WorkingDir: "/",
		User:       "",
	}

	mount("none", "/proc", "proc", 0)
	mount("none", "/dev/pts", "devpts", 0)
	mount("none", "/dev/mqueue", "mqueue", 0)
	mount("none", "/dev/shm", "tmpfs", 0)
	mount("none", "/sys", "sysfs", 0)
	mount("none", "/sys/fs/cgroup", "cgroup", 0)

	setHostname("blanc")

	cmd := exec.Command(imageConfig.Entrypoint)
	cmd.Env = imageConfig.Env
	cmd.Dir = imageConfig.WorkingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func mount(source, target, fileSystemType string, flags uintptr) {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		if err := os.MkdirAll(target, 0755); err != nil {
			panic(err)
		}
	}

	if err := syscall.Mount(source, target, fileSystemType, flags, ""); err != nil {
		panic(err)
	}
}

func setHostname(hostname string) {
	if err := syscall.Sethostname([]byte(hostname)); err != nil {
		panic(err)
	}

	if err := os.WriteFile("/etc/hostname", []byte(hostname), 0755); err != nil {
		panic(err)
	}
}
