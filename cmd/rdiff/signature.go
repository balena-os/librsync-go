package main

import (
	_ "io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/balena-os/librsync-go"
)

func CommandSignature(c *cli.Context) {
	if len(c.Args()) > 2 {
		logrus.Warnf("%d additional arguments passed are ignored", len(c.Args())-2)
	}

	if c.Args().Get(0) == "" {
		logrus.Fatalf("Missing basis file")
	}

	if c.Args().Get(1) == "" {
		logrus.Fatalf("Missing signature file")
	}

	var sigType librsync.MagicNumber

	switch c.String("hash") {
	case "blake2":
		sigType = librsync.BLAKE2_SIG_MAGIC
	case "md4":
		sigType = librsync.MD4_SIG_MAGIC
	default:
		logrus.Fatalf("Invalid hash type: %v", c.String("hash"))
	}

	basis, err := os.Open(c.Args().Get(0))
	if err != nil {
		logrus.Fatal(err)
	}
	defer basis.Close()

	signature, err := os.OpenFile(c.Args().Get(1), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		logrus.Fatal(err)
	}
	defer signature.Close()

	_, err = librsync.Signature(basis, signature, uint32(c.Uint("block-size")), uint32(c.Uint("sum-size")), sigType)
	if err != nil {
		logrus.Fatal(err)
	}
}
