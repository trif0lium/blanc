package main

import (
	"os"

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

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
