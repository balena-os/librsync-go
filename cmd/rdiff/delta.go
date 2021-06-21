package main

import (
	"os"

	"github.com/balena-os/librsync-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func CommandDelta(c *cli.Context) {
	if len(c.Args()) > 3 {
		logrus.Warnf("%d additional arguments passed are ignored", len(c.Args())-2)
	}

	if c.Args().Get(0) == "" {
		logrus.Fatalf("Missing signature file")
	}

	if c.Args().Get(1) == "" {
		logrus.Fatalf("Missing newfile file")
	}

	if c.Args().Get(2) == "" {
		logrus.Fatalf("Missing delta file")
	}

	signature, err := librsync.ReadSignatureFile(c.Args().Get(0))
	if err != nil {
		logrus.Fatal(err)
	}

	newfile, err := os.Open(c.Args().Get(1))
	if err != nil {
		logrus.Fatal(err)
	}
	defer newfile.Close()

	delta, err := os.OpenFile(c.Args().Get(2), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		logrus.Fatal(err)
	}
	defer newfile.Close()

	err = librsync.Delta(signature, newfile, delta)
	if err != nil {
		logrus.Fatal(err)
	}
}
