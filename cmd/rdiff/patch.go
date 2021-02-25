package main

import (
	_ "io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/balena-os/librsync-go"
)

func CommandPatch(c *cli.Context) {
	if len(c.Args()) > 3 {
		logrus.Warnf("%d additional arguments passed are ignored", len(c.Args())-2)
	}

	if c.Args().Get(0) == "" {
		logrus.Fatalf("Missing basis file")
	}

	if c.Args().Get(1) == "" {
		logrus.Fatalf("Missing delta file")
	}
	if c.Args().Get(2) == "" {
		logrus.Fatalf("Missing newfile file")
	}

	basis, err := os.Open(c.Args().Get(0))
	if err != nil {
		logrus.Fatal(err)
	}
	defer basis.Close()

	delta, err := os.Open(c.Args().Get(1))
	if err != nil {
		logrus.Fatal(err)
	}
	defer delta.Close()

	newfile, err := os.OpenFile(c.Args().Get(2), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		logrus.Fatal(err)
	}
	defer newfile.Close()

	if err := librsync.Patch(basis, delta, newfile); err != nil {
		logrus.Fatal(err)
	}
}
