package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/urfave/cli/v2"
)

var (
	WORKING_DIRECTORY = "/var/lib/blanc"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "run",
				Action: func(cCtx *cli.Context) error {
					imageRef := cCtx.Args().First()

					vmID, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz", 17)
					if err != nil {
						return err
					}

					workingDir := filepath.Join(WORKING_DIRECTORY, "vms/vm_"+vmID)

					if err := os.MkdirAll(workingDir, 0755); err != nil {
						return err
					}

					if err := os.MkdirAll(filepath.Join(workingDir, "rootfs"), 0755); err != nil {
						return err
					}

					if err := os.MkdirAll(filepath.Join(workingDir, "container/unpacked"), 0755); err != nil {
						return err
					}

					cmd := exec.CommandContext(
						cCtx.Context,
						"skopeo",
						"copy",
						"docker://"+imageRef,
						"oci:"+filepath.Join(workingDir, "container")+":latest",
					)

					if err := cmd.Run(); err != nil {
						return err
					}

					cmd = exec.CommandContext(
						cCtx.Context,
						"umoci",
						"unpack",
						"--image", filepath.Join(workingDir, "container")+":latest",
						filepath.Join(workingDir, "container/unpacked"),
					)

					if err := cmd.Run(); err != nil {
						return err
					}

					if err := exec.CommandContext(cCtx.Context, "fallocate", "-l", "5G", filepath.Join(workingDir, "rootfs.img")).Run(); err != nil {
						return err
					}

					if err := exec.CommandContext(cCtx.Context, "mkfs.ext4", filepath.Join(workingDir, "rootfs.img")).Run(); err != nil {
						return err
					}

					if err := exec.CommandContext(cCtx.Context, "mount", "-o", "loop", filepath.Join(workingDir, "rootfs.img"), filepath.Join(workingDir, "rootfs")).Run(); err != nil {
						return err
					}

					if err := exec.CommandContext(cCtx.Context, "cp", "-R", filepath.Join(workingDir, "container/unpacked/rootfs"), filepath.Join(workingDir, "rootfs")).Run(); err != nil {
						return err
					}

					if err := exec.CommandContext(cCtx.Context, "umount", filepath.Join(workingDir, "rootfs")).Run(); err != nil {
						return err
					}

					if err := exec.CommandContext(cCtx.Context, "firecracker", "--api-sock", filepath.Join(workingDir, "firecracker.sock")).Run(); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func httpRequest(client *http.Client, path string, body map[string]any) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, "http://"+filepath.Join("localhost", path), bytes.NewBuffer(jsonData))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
