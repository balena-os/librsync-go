package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rdiff"
	app.Usage = "binary diff utility based on librsync-go"
	app.Version = "v0.0.1"
	app.Author = "Petros Angelatos"
	app.Email = "petrosagg@gmail.com"
	app.Action = cli.ShowAppHelp
	app.Commands = []cli.Command{
		{
			Name:      "signature",
			Usage:     "creates a signature of the input file",
			ArgsUsage: "BASIS SIGNATURE",
			Action:    CommandSignature,
			Flags: []cli.Flag{
				cli.UintFlag{
					Name:  "block-size, b",
					Value: 2048,
					Usage: "Signature block size",
				},
				cli.UintFlag{
					Name:  "sum-size, S",
					Value: 32,
					Usage: "Set signature strength",
				},
				cli.StringFlag{
					Name:  "hash, H",
					Value: "blake2",
					Usage: "Hash algorithm: blake2, md4",
				},
			},
		},
		{
			Name:      "delta",
			Usage:     "calculates the binary diff between old and new files",
			ArgsUsage: "SIGNATURE NEWFILE DELTA",
			Action:    CommandDelta,
		},
		{
			Name:      "patch",
			Usage:     "uses the delta file and old file to produce the new file",
			ArgsUsage: "BASIS DELTA NEWFILE",
			Action:    CommandPatch,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
