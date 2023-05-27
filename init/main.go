package main

import (
	"os"
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
