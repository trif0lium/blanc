package main

import (
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

					if err := os.MkdirAll(filepath.Join(workingDir, "container/unpacked", 0755)); err != nil {
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

					cmd := exec.CommandContext(
						cCtx.Context,
						"umoci",
						"unpack",
						"--image", filepath.Join(workingDir, "container")+":latest",
						filepath.Join(workingDir, "container/unpacked")
					)

					if err := cmd.Run(); err != nil {
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
